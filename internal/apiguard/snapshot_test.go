package apiguard

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"testing"
)

func snapshot(t *testing.T, source string) Snapshot {
	t.Helper()
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "fixture.go", "package fixture\n"+source, 0)
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := new(types.Config).Check("example.test/fixture", fs, []*ast.File{file}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s, err := Capture([]*types.Package{pkg})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSnapshotDetectsConsumerVisibleChanges(t *testing.T) {
	for _, tc := range []struct{ name, before, after, contains string }{
		{"removed function", `func Draw(int) {}`, ``, "Draw"},
		{"parameter type", `func Draw(int) {}`, `func Draw(string) {}`, "Draw"},
		{"variadic", `func Draw([]int) {}`, `func Draw(...int) {}`, "Draw"},
		{"field removal", `type Props struct { Label string; Ready bool }`, `type Props struct { Label string }`, "Props"},
		{"field addition needs review", `type Props struct { Label string }`, `type Props struct { Label string; Ready bool }`, "Props"},
		{"field reorder", `type Props struct { Label string; Ready bool }`, `type Props struct { Ready bool; Label string }`, "Props"},
		{"field tag", "type Props struct { Label string `json:\"label\"` }", "type Props struct { Label string `json:\"name\"` }", "Props"},
		{"constant value", `const Ready = 1`, `const Ready = 2`, "Ready"},
		{"constant type", `const Ready = 1`, `const Ready int = 1`, "Ready"},
		{"resolved iota", `const (_ = iota; Ready)`, `const (_ = iota; _; Ready)`, "Ready"},
		{"alias", `type Props = string`, `type Props string`, "Props"},
		{"interface method", `type Renderer interface { Draw() }`, `type Renderer interface { Draw() error }`, "Renderer"},
		{"method pointer receiver", `type Props int;func (Props) Draw(){}`, `type Props int;func (*Props) Draw(){}`, "Props.Draw"},
		{"pointer method result", `type Props int;func (*Props) Draw() int{return 1}`, `type Props int;func (*Props) Draw() string{return "x"}`, "(*Props).Draw"},
		{"promoted method", `type inner int;func(inner) Draw(){};type Props struct{inner}`, `type inner int;type Props struct{inner}`, "Props.Draw"},
		{"generic constraint", `type Box[T ~int] struct{Value T}`, `type Box[T ~string] struct{Value T}`, "Box"},
		{"generic function", `func Draw[T ~int](v T){}`, `func Draw[T ~string](v T){}`, "Draw"},
		{"variable inferred type", `var Default = 1`, `var Default = "x"`, "Default"},
		{"reachable private type", `type result struct{Value string};func Draw() result{return result{}}`, `type result struct{Value int};func Draw() result{return result{}}`, "reachable result"},
		{"addition", ``, `func NewFeature(){}`, "NewFeature"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changes, err := Compare(snapshot(t, tc.before), snapshot(t, tc.after))
			if err != nil {
				t.Fatal(err)
			}
			for _, change := range changes {
				if change.Symbol == tc.contains {
					return
				}
			}
			t.Fatalf("missed %s: %+v", tc.contains, changes)
		})
	}
}

func TestSnapshotIgnoresBodiesCommentsAndEquivalentConstantExpressions(t *testing.T) {
	before := snapshot(t, `const Answer = 2*3;func Draw() int{return 1};type recursive struct{Next *recursive};func Node() *recursive{return nil}`)
	after := snapshot(t, `// Changed documentation.
 const Answer = 6;func Draw() int{return 42};type recursive struct{Next *recursive};func Node() *recursive{return &recursive{}}`)
	changes, err := Compare(before, after)
	if err != nil || len(changes) != 0 {
		t.Fatalf("implementation drifted API: %v %+v", err, changes)
	}
	first, _ := json.Marshal(before)
	second, _ := json.Marshal(after)
	if string(first) != string(second) {
		t.Fatal("snapshot serialization is not deterministic")
	}
}

func TestInvalidSnapshotsAndPackageChanges(t *testing.T) {
	valid := snapshot(t, `type Props struct{Label string}`)
	for _, s := range []Snapshot{{}, {Format: 2, Packages: valid.Packages}, {Format: Format, Packages: map[string]map[string]string{"": {}}}, {Format: Format, Packages: map[string]map[string]string{"p": {"": ""}}}} {
		if _, err := Compare(s, valid); err == nil {
			t.Fatal("accepted invalid baseline")
		}
	}
	next := Snapshot{Format: Format, Packages: map[string]map[string]string{"new/package": {}}}
	changes, err := Compare(valid, next)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 3 || changes[0].Kind != "package removed" || changes[2].Kind != "package added" {
		t.Fatalf("package changes lost: %+v", changes)
	}
	again, _ := Compare(valid, next)
	if !reflect.DeepEqual(changes, again) {
		t.Fatal("nondeterministic change order")
	}
	if _, err := Capture(nil); err == nil {
		t.Fatal("accepted no packages")
	}
	if _, err := Capture([]*types.Package{types.NewPackage("incomplete", "incomplete")}); err == nil {
		t.Fatal("accepted incomplete package")
	}
	for _, decl := range valid.Packages["example.test/fixture"] {
		if !strings.Contains(decl, "Label string") {
			t.Fatal(decl)
		}
	}
}

func TestGenericInstancesDoNotOverwriteDeclaredMethods(t *testing.T) {
	s := snapshot(t, `type Box[T any] struct{ Value T };func(b Box[T]) Get() T {return b.Value};func Values()(Box[int],Box[string]){panic("fixture")};type Alias = Box[int]`)
	objects := s.Packages["example.test/fixture"]
	if objects["Box.Get"] != "func() T" || objects["(*Box).Get"] != "func() T" || objects["Alias.Get"] != "func() int" {
		t.Fatalf("generic method sets overwritten: %+v", objects)
	}
}

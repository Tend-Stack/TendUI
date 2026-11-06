// Package apiguard records compiler-resolved public Go declarations for review.
// It is a conservative change gate, not a semantic compatibility solver.
package apiguard

import (
	"fmt"
	"go/types"
	"sort"
)

const Format = 1

type Snapshot struct {
	Format   int                          `json:"format"`
	Packages map[string]map[string]string `json:"packages"`
}

type Change struct {
	Package string `json:"package"`
	Symbol  string `json:"symbol"`
	Kind    string `json:"kind"`
	Before  string `json:"before,omitempty"`
	After   string `json:"after,omitempty"`
}

func qualified(p *types.Package) string { return p.Path() }

// Capture includes resolved constant values, type definitions, struct field order
// and tags, function signatures, aliases and both value and pointer method sets.
// Local private types reachable through declarations are recorded too. Imported
// types remain qualified references: their implementations are not copied here.
func Capture(packages []*types.Package) (Snapshot, error) {
	s := Snapshot{Format: Format, Packages: map[string]map[string]string{}}
	if len(packages) == 0 {
		return s, fmt.Errorf("API snapshot needs at least one package")
	}
	for _, pkg := range packages {
		if pkg == nil || !pkg.Complete() {
			return s, fmt.Errorf("incomplete API package")
		}
		if _, ok := s.Packages[pkg.Path()]; ok {
			return s, fmt.Errorf("duplicate API package %s", pkg.Path())
		}
		objects := map[string]string{}
		s.Packages[pkg.Path()] = objects
		seen := map[types.Type]bool{}
		var walk func(types.Type)
		methods := func(name string, t types.Type) {
			for _, m := range []struct {
				prefix string
				set    *types.MethodSet
			}{{name, types.NewMethodSet(t)}, {"(*" + name + ")", types.NewMethodSet(types.NewPointer(t))}} {
				for i := 0; i < m.set.Len(); i++ {
					selection := m.set.At(i)
					if selection.Obj().Exported() {
						objects[m.prefix+"."+selection.Obj().Name()] = types.TypeString(selection.Type(), qualified)
						walk(selection.Type())
					}
				}
			}
		}
		walk = func(t types.Type) {
			if t == nil || seen[t] {
				return
			}
			seen[t] = true
			params := func(list *types.TypeParamList) {
				if list != nil {
					for i := 0; i < list.Len(); i++ {
						walk(list.At(i))
					}
				}
			}
			switch t := t.(type) {
			case *types.Named:
				for i := 0; i < t.TypeArgs().Len(); i++ {
					walk(t.TypeArgs().At(i))
				}
				if t != t.Origin() {
					walk(t.Origin())
					return
				}
				if t.Obj().Pkg() != pkg {
					return
				}
				if !t.Obj().Exported() {
					objects["reachable "+t.Obj().Name()] = types.ObjectString(t.Obj(), qualified)
				}
				params(t.TypeParams())
				methods(t.Obj().Name(), t)
				walk(t.Underlying())
			case *types.Alias:
				for i := 0; i < t.TypeArgs().Len(); i++ {
					walk(t.TypeArgs().At(i))
				}
				if t != t.Origin() {
					walk(t.Origin())
					return
				}
				if t.Obj().Pkg() != pkg {
					return
				}
				if !t.Obj().Exported() {
					objects["reachable "+t.Obj().Name()] = types.ObjectString(t.Obj(), qualified)
				}
				params(t.TypeParams())
				walk(t.Rhs())
			case *types.Pointer:
				walk(t.Elem())
			case *types.Slice:
				walk(t.Elem())
			case *types.Array:
				walk(t.Elem())
			case *types.Map:
				walk(t.Key())
				walk(t.Elem())
			case *types.Chan:
				walk(t.Elem())
			case *types.Struct:
				for i := 0; i < t.NumFields(); i++ {
					walk(t.Field(i).Type())
				}
			case *types.Interface:
				t.Complete()
				for i := 0; i < t.NumMethods(); i++ {
					walk(t.Method(i).Type())
				}
				for i := 0; i < t.NumEmbeddeds(); i++ {
					walk(t.EmbeddedType(i))
				}
			case *types.Signature:
				params(t.TypeParams())
				walk(t.Params())
				walk(t.Results())
			case *types.Tuple:
				for i := 0; i < t.Len(); i++ {
					walk(t.At(i).Type())
				}
			case *types.TypeParam:
				walk(t.Constraint())
			case *types.Union:
				for i := 0; i < t.Len(); i++ {
					walk(t.Term(i).Type())
				}
			}
		}
		for _, name := range pkg.Scope().Names() {
			object := pkg.Scope().Lookup(name)
			if !object.Exported() {
				continue
			}
			declaration := types.ObjectString(object, qualified)
			if value, ok := object.(*types.Const); ok {
				declaration += " = " + value.Val().ExactString()
			}
			objects[name] = declaration
			walk(object.Type())
			if typ, ok := object.(*types.TypeName); ok && typ.IsAlias() {
				methods(name, typ.Type())
			}
		}
	}
	return s, nil
}

func (s Snapshot) Validate() error {
	if s.Format != Format || len(s.Packages) == 0 {
		return fmt.Errorf("invalid API snapshot format or empty packages")
	}
	for pkg, objects := range s.Packages {
		if pkg == "" || objects == nil {
			return fmt.Errorf("invalid API package %q", pkg)
		}
		for name, decl := range objects {
			if name == "" || decl == "" {
				return fmt.Errorf("empty API declaration in %s", pkg)
			}
		}
	}
	return nil
}

// Compare flags every declaration change, including additions requiring review.
// An empty package is tracked as well, so package removal cannot disappear silently.
func Compare(before, after Snapshot) ([]Change, error) {
	if err := before.Validate(); err != nil {
		return nil, err
	}
	if err := after.Validate(); err != nil {
		return nil, err
	}
	changes := []Change{}
	packages := map[string]bool{}
	for p := range before.Packages {
		packages[p] = true
	}
	for p := range after.Packages {
		packages[p] = true
	}
	for p := range packages {
		old, oldExists := before.Packages[p]
		next, nextExists := after.Packages[p]
		if !oldExists {
			changes = append(changes, Change{Package: p, Kind: "package added"})
		}
		if !nextExists {
			changes = append(changes, Change{Package: p, Kind: "package removed"})
		}
		names := map[string]bool{}
		for n := range old {
			names[n] = true
		}
		for n := range next {
			names[n] = true
		}
		for n := range names {
			a, aExists := old[n]
			b, bExists := next[n]
			if aExists && bExists && a == b {
				continue
			}
			kind := "changed"
			if !aExists {
				kind = "added"
			} else if !bExists {
				kind = "removed"
			}
			changes = append(changes, Change{Package: p, Symbol: n, Kind: kind, Before: a, After: b})
		}
	}
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Package != changes[j].Package {
			return changes[i].Package < changes[j].Package
		}
		return changes[i].Symbol < changes[j].Symbol
	})
	return changes, nil
}

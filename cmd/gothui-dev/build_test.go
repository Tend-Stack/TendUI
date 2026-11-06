package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// This executes the real native candidate and replaces it twice. The CI matrix
// runs it on each OS, including Windows where an executable suffix is required.
func TestNativeBuildExecutesBeforeReplacingReviewBinary(t *testing.T) {
	root := filepath.Join(t.TempDir(), "module with spaces")
	for name, body := range map[string]string{
		"go.mod":                 "module fixture\n\ngo 1.26.0\n",
		".gitignore":             "/tests/results/\n",
		"examples/basic/main.go": "package main\nimport \"fmt\"\nvar version string\nfunc main(){fmt.Println(version)}\n",
	} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"commit", "-qm", "Build fixture"}} {
		cmd := exec.Command("git", append([]string{"-c", "user.name=GothUI build test", "-c", "user.email=test@example.invalid"}, args...)...)
		cmd.Dir = root
		if b, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git: %v %s", err, b)
		}
	}
	r := runner{root: root, results: filepath.Join(root, "tests/results/validation"), goTool: goExecutable()}
	if err := os.MkdirAll(r.results, 0755); err != nil {
		t.Fatal(err)
	}
	name := "gothui-demo"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(root, "tests/results/demos", name)
	var original string
	for iteration := 0; iteration < 2; iteration++ {
		if iteration == 1 {
			f, err := os.OpenFile(filepath.Join(root, "examples/basic/main.go"), os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.WriteString("\n// Changed source must produce a dirty build identifier.\n")
			f.Close()
			if err != nil {
				t.Fatal(err)
			}
		}
		if err := r.build(); err != nil {
			t.Fatal(err)
		}
		got, err := r.output(path, "-version")
		if err != nil {
			t.Fatal(err)
		}
		if iteration == 0 {
			original = got
		} else if got != original+"-dirty" {
			t.Fatalf("candidate did not replace original: %s -> %s", original, got)
		}
		data, err := os.ReadFile(filepath.Join(r.results, "build.json"))
		if err != nil {
			t.Fatal(err)
		}
		var report struct {
			Executable string `json:"executable"`
			Version    string `json:"version"`
			Smoke      bool   `json:"version_smoke"`
		}
		if err := json.Unmarshal(data, &report); err != nil {
			t.Fatal(err)
		}
		if report.Executable != name || report.Version != got || !report.Smoke {
			t.Fatalf("wrong build proof: %s", data)
		}
	}
	files, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Name() != name || strings.HasSuffix(original, "-dirty") {
		t.Fatalf("unexpected leftovers or initial version: %v %s", files, original)
	}
}

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

func TestBuildVersionUsesGitAndIncludesUntrackedWork(t *testing.T) {
	dir := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-c", "user.name=GothUI test", "-c", "user.email=test@example.invalid"}, args...)...)
		cmd.Dir = dir
		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, b)
		}
		return strings.TrimSpace(string(b))
	}
	git("init", "-q")
	if err := os.WriteFile(filepath.Join(dir, "source"), []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}
	git("add", "source")
	git("commit", "-qm", "Initial test commit")
	r := runner{root: dir, results: t.TempDir(), goTool: filepath.Join(runtime.GOROOT(), "bin", "go")}
	version, err := r.buildVersion()
	if err != nil || !strings.Contains(version, "-dev-g") || strings.HasSuffix(version, "-dirty") {
		t.Fatal(version, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "untracked"), []byte("new work"), 0644); err != nil {
		t.Fatal(err)
	}
	version, err = r.buildVersion()
	if err != nil || !strings.HasSuffix(version, "-dirty") {
		t.Fatal("untracked state lost", version, err)
	}
	if err := r.report("test.json", map[string]any{"passed": true}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(r.results, "test.json"))
	if err != nil {
		t.Fatal(err)
	}
	var report map[string]any
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report["dirty"] != true || report["revision"] != git("rev-parse", "HEAD") || report["toolchain"] != runtime.Version() {
		t.Fatalf("incomplete provenance: %s", data)
	}
	if err := os.Remove(filepath.Join(dir, "untracked")); err != nil {
		t.Fatal(err)
	}
	git("tag", "v2026.1.0")
	version, err = r.buildVersion()
	if err != nil || !strings.HasPrefix(version, "v2026.1.0-0-g") {
		t.Fatal("tag describe lost", version, err)
	}
	git("tag", "-d", "v2026.1.0")
	git("tag", "v2026.not-a-release")
	if version, err := r.buildVersion(); err == nil {
		t.Fatal("accepted malformed release tag", version)
	}
}

func TestSourceWalkExcludesGeneratedEvidence(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"main.go", "components/part.templ", "tests/results/fixture/ignore.go", ".git/ignore.go"} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			t.Fatal(err)
		}
	}
	paths, err := (runner{root: root}).files(".go")
	if err != nil || len(paths) != 1 || paths[0] != "main.go" {
		t.Fatal(paths, err)
	}
}

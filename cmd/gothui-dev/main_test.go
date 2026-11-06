package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunnerOutputSeparatesDiagnostics(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOTHUI_TEST_OUTPUT_HELPER", "1")
	r := runner{root: t.TempDir(), goTool: goExecutable()}
	output, err := r.output(executable, "-test.run=^TestRunnerOutputHelper$")
	if err != nil || output != `{"package":"example.invalid/demo"}` {
		t.Fatalf("successful diagnostics contaminated machine output: %q, %v", output, err)
	}
	t.Setenv("GOTHUI_TEST_OUTPUT_FAILURE", "1")
	output, err = r.output(executable, "-test.run=^TestRunnerOutputHelper$")
	if output != "" || err == nil || !strings.Contains(err.Error(), "fixture diagnostic") || !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("failed command lost diagnostics: %q, %v", output, err)
	}
}

// A subprocess fixture is portable across native Unix and Windows tests.
func TestRunnerOutputHelper(t *testing.T) {
	if os.Getenv("GOTHUI_TEST_OUTPUT_HELPER") != "1" {
		return
	}
	fmt.Fprintln(os.Stdout, `{"package":"example.invalid/demo"}`)
	fmt.Fprintln(os.Stderr, "fixture diagnostic on stderr")
	if os.Getenv("GOTHUI_TEST_OUTPUT_FAILURE") == "1" {
		os.Exit(7)
	}
	os.Exit(0)
}

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
	r := runner{root: dir, results: t.TempDir(), goTool: goExecutable()}
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

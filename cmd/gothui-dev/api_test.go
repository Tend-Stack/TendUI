package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tend-Stack/GothUI/v2026/internal/apiguard"
)

func TestAPICheckUsesCompiledExportsAndRejectsDrift(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module "+modulePath+"\n\ngo 1.26.0\n")
	write("gothui.go", "package gothui\nconst Version=\"2026.1.0\"\n")
	write("components/components.go", "package components\ntype Props struct{Label string}\n")
	write("render/render.go", "package render\nfunc HTML(){}\n")
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-c", "user.name=GothUI API test", "-c", "user.email=test@example.invalid"}, args...)...)
		cmd.Dir = root
		if b, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git: %v %s", err, b)
		}
	}
	git("init", "-q")
	git("add", ".")
	git("commit", "-qm", "API fixture")
	r := runner{root: root, results: filepath.Join(root, "tests/results/validation"), goTool: goExecutable()}
	if err := os.MkdirAll(r.results, 0755); err != nil {
		t.Fatal(err)
	}
	if err := r.api(); err == nil || !strings.Contains(err.Error(), "read API baseline") {
		t.Fatalf("missing baseline was not rejected: %v", err)
	}
	candidate, err := os.ReadFile(filepath.Join(r.results, "api-current.json"))
	if err != nil {
		t.Fatal(err)
	}
	write("api/public.json", string(candidate))
	if err := r.api(); err != nil {
		t.Fatal(err)
	}
	report, err := os.ReadFile(filepath.Join(r.results, "api.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(report), `"matching": true`) {
		t.Fatalf("missing success report: %s", report)
	}
	write("components/components.go", "package components\ntype Props struct{Label int}\n")
	if err := r.api(); err == nil || !strings.Contains(err.Error(), "API changed") {
		t.Fatalf("changed field was not rejected: %v", err)
	}
	if _, err := os.Stat(filepath.Join(r.results, "api.json")); !os.IsNotExist(err) {
		t.Fatalf("stale success report survives: %v", err)
	}
	baseline, err := os.ReadFile(filepath.Join(root, "api/public.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(baseline) != string(candidate) {
		t.Fatal("check rewrote the reviewed baseline")
	}
	changesJSON, err := os.ReadFile(filepath.Join(r.results, "api-changes.json"))
	if err != nil {
		t.Fatal(err)
	}
	var changes struct {
		Changes  []apiguard.Change `json:"changes"`
		Matching bool              `json:"matching"`
	}
	if err := json.Unmarshal(changesJSON, &changes); err != nil {
		t.Fatal(err)
	}
	if changes.Matching || len(changes.Changes) != 1 || changes.Changes[0].Symbol != "Props" || changes.Changes[0].Kind != "changed" || !strings.Contains(changes.Changes[0].After, "Label int") {
		t.Fatalf("bad drift evidence: %s", changesJSON)
	}
	write("extra/extra.go", "package extra\nfunc Public(){}\n")
	if err := r.api(); err == nil || !strings.Contains(err.Error(), "unreviewed public package") {
		t.Fatalf("new public package bypassed baseline coverage: %v", err)
	}
}

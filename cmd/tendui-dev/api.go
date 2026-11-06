package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Tend-Stack/TendUI/v2026/internal/apiguard"
)

// api reads compiler export data, so aliases, inferred types and iota constants
// are resolved exactly as they are in a consumer's Go build.
func (r runner) api() error {
	for _, name := range []string{"api.json", "api-current.json", "api-changes.json"} {
		if err := os.Remove(filepath.Join(r.results, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	paths := []string{modulePath, modulePath + "/components", modulePath + "/render"}
	known := make(map[string]bool, len(paths))
	for _, path := range paths {
		known[path] = true
	}
	listed, err := r.output(r.goTool, "list", "-f", "{{.ImportPath}}", "./...")
	if err != nil {
		return err
	}
	for _, path := range strings.Fields(listed) {
		rel := strings.TrimPrefix(path, modulePath+"/")
		if rel == "cmd" || strings.HasPrefix(rel, "cmd/") || rel == "examples" || strings.HasPrefix(rel, "examples/") || rel == "tests" || strings.HasPrefix(rel, "tests/") || strings.Contains(path, "/internal/") || strings.HasSuffix(path, "/internal") {
			continue
		}
		if !known[path] {
			return fmt.Errorf("unreviewed public package %s; add it to the API package list and baseline", path)
		}
	}
	args := append([]string{"list", "-deps", "-export", "-json"}, paths...)
	output, err := r.output(r.goTool, args...)
	if err != nil {
		return err
	}
	exports := map[string]string{}
	decoder := json.NewDecoder(bytes.NewBufferString(output))
	for {
		var p struct{ ImportPath, Export string }
		if err := decoder.Decode(&p); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		exports[p.ImportPath] = p.Export
	}
	loader := importer.ForCompiler(token.NewFileSet(), "gc", func(path string) (io.ReadCloser, error) {
		file := exports[path]
		if file == "" {
			return nil, fmt.Errorf("missing compiler export for %s", path)
		}
		return os.Open(file)
	})
	var packages []*types.Package
	for _, path := range paths {
		pkg, err := loader.Import(path)
		if err != nil {
			return err
		}
		packages = append(packages, pkg)
	}
	current, err := apiguard.Capture(packages)
	if err != nil {
		return err
	}
	candidate, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(r.results, "api-current.json"), append(candidate, '\n'), 0644); err != nil {
		return err
	}
	baseline, err := os.Open(filepath.Join(r.root, "api/public.json"))
	if err != nil {
		return fmt.Errorf("read API baseline: %w; candidate is tests/results/validation/api-current.json", err)
	}
	defer baseline.Close()
	var previous apiguard.Snapshot
	decoder = json.NewDecoder(baseline)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&previous); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("API baseline must contain one JSON object")
	}
	changes, err := apiguard.Compare(previous, current)
	if err != nil {
		return err
	}
	if err := r.report("api-changes.json", map[string]any{"changes": changes, "matching": len(changes) == 0}); err != nil {
		return err
	}
	if len(changes) != 0 {
		return fmt.Errorf("API changed (%d entries); review tests/results/validation/api-changes.json and docs/api-compatibility.md before updating api/public.json", len(changes))
	}
	declarations := 0
	for _, p := range current.Packages {
		declarations += len(p)
	}
	target, err := r.output(r.goTool, "env", "GOOS", "GOARCH")
	if err != nil {
		return err
	}
	platform := strings.Fields(target)
	if len(platform) != 2 {
		return fmt.Errorf("unexpected Go target %q", target)
	}
	return r.report("api.json", map[string]any{"packages": paths, "declarations": declarations, "matching": true, "baseline": "api/public.json", "goos": platform[0], "goarch": platform[1]})
}

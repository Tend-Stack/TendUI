// tendui-dev provides repeatable local checks without frontend tooling.
package main

import (
	"bytes"
	"crypto/sha256"
	"debug/buildinfo"
	"debug/elf"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Tend-Stack/TendUI/v2026"
	"github.com/Tend-Stack/TendUI/v2026/internal/recipes"
)

const modulePath = "github.com/Tend-Stack/TendUI/v2026"
const stylesheetBudget = 64 << 10
const demoBudget = 16 << 20

type runner struct{ root, results, goTool, templTool string }

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/tendui-dev check|api|generated|recipes|consumer|bench|build")
		os.Exit(2)
	}
	root, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	mod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil || !bytes.Contains(mod, []byte("module "+modulePath+"\n")) {
		fatal(errors.New("run tendui-dev from the repository root"))
	}
	r := runner{root: root, results: filepath.Join(root, "tests/results/validation"), goTool: goExecutable(), templTool: os.Getenv("TEMPL")}
	if r.templTool == "" {
		r.templTool = "templ"
	}
	if err := os.MkdirAll(r.results, 0755); err != nil {
		fatal(err)
	}
	// A failed rerun must not leave an older success report beside new logs.
	reports := map[string]string{"check": "checks.json", "api": "api.json", "generated": "generated.json", "recipes": "recipes.json", "consumer": "consumer.json", "bench": "benchmarks.json", "build": "build.json"}
	if name := reports[os.Args[1]]; name != "" {
		if err := os.Remove(filepath.Join(r.results, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			fatal(err)
		}
	}
	switch os.Args[1] {
	case "check":
		err = r.check()
	case "api":
		err = r.api()
	case "generated":
		err = r.generated()
	case "recipes":
		err = r.recipes()
	case "consumer":
		err = r.consumer()
	case "bench":
		err = r.run("benchmarks.log", r.root, r.goTool, "test", "-run=^$", "-bench=Benchmark", "-benchmem", "-count=3", "-benchtime=200ms", "-cpu=2", "./...")
		if err == nil {
			err = r.report("benchmarks.json", map[string]any{"samples": 3, "benchtime": "200ms", "gomaxprocs": 2, "browser_measurement": false})
		}
	case "build":
		err = r.build()
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fatal(err)
	}
	fmt.Println(os.Args[1] + ": passed; evidence in tests/results/validation")
}

func fatal(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

func executableName(name, goos string) string {
	if goos == "windows" {
		return name + ".exe"
	}
	return name
}

func goExecutable() string {
	return filepath.Join(runtime.GOROOT(), "bin", executableName("go", runtime.GOOS))
}

func (r runner) command(dir, program string, args ...string) *exec.Cmd {
	cmd := exec.Command(program, args...)
	cmd.Dir = dir
	// Reuse the toolchain that compiled this tool, including when PATH's go is a wrapper.
	cmd.Env = append(os.Environ(), "PATH="+filepath.Dir(r.goTool)+string(os.PathListSeparator)+os.Getenv("PATH"))
	return cmd
}

func (r runner) run(logName, dir, program string, args ...string) error {
	fmt.Println("Running", filepath.Base(program), strings.Join(args, " "))
	f, err := os.Create(filepath.Join(r.results, logName))
	if err != nil {
		return err
	}
	defer f.Close()
	cmd := r.command(dir, program, args...)
	cmd.Stdout, cmd.Stderr = f, f
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w; see %s", filepath.Base(program), err, filepath.Join(r.results, logName))
	}
	return nil
}

func (r runner) output(program string, args ...string) (string, error) {
	cmd := r.command(r.root, program, args...)
	var diagnostics bytes.Buffer
	cmd.Stderr = &diagnostics
	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w: %s", filepath.Base(program), err, strings.TrimSpace(string(b)+"\n"+diagnostics.String()))
	}
	// Go can emit warnings on a successful command. Keep those visible, but
	// never mix them into package paths, JSON, versions or other parsed stdout.
	if diagnostics.Len() > 0 {
		fmt.Fprint(os.Stderr, diagnostics.String())
	}
	return strings.TrimSpace(string(b)), nil
}

func (r runner) files(suffix string) ([]string, error) {
	var result []string
	err := filepath.WalkDir(r.root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(r.root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if rel != "." && (strings.HasPrefix(entry.Name(), ".") || rel == filepath.Join("tests", "results")) {
				return filepath.SkipDir
			}
		} else if strings.HasSuffix(path, suffix) {
			result = append(result, rel)
		}
		return nil
	})
	return result, err
}

func (r runner) check() error {
	files, err := r.files(".go")
	if err != nil {
		return err
	}
	for _, path := range files {
		source, err := os.ReadFile(filepath.Join(r.root, path))
		if err != nil {
			return err
		}
		formatted, err := format.Source(source)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !bytes.Equal(source, formatted) {
			return fmt.Errorf("%s needs gofmt", path)
		}
	}
	css, err := os.Stat(filepath.Join(r.root, "static/tendui.css"))
	if err != nil {
		return err
	}
	if css.Size() > stylesheetBudget {
		return fmt.Errorf("CSS %d exceeds %d-byte budget", css.Size(), stylesheetBudget)
	}
	// Production package dependencies, not the generator's module graph.
	for _, packages := range [][]string{{".", "./render"}, {"./components"}} {
		args := append([]string{"list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}} {{if .Module}}{{.Module.Path}}{{end}}{{end}}"}, packages...)
		result, err := r.output(r.goTool, args...)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(result, "\n") {
			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}
			if len(fields) != 2 || fields[1] != modulePath && !(len(packages) == 1 && fields[1] == "github.com/a-h/templ") {
				return fmt.Errorf("unexpected runtime dependency: %s", line)
			}
		}
	}
	if err := r.api(); err != nil {
		return err
	}
	if err := r.run("tests.log", r.root, r.goTool, "test", "-race", "./..."); err != nil {
		return err
	}
	if err := r.run("vet.log", r.root, r.goTool, "vet", "./..."); err != nil {
		return err
	}
	return r.report("checks.json", map[string]any{"css_bytes": css.Size(), "css_budget_bytes": stylesheetBudget, "checks": []string{"gofmt", "runtime dependencies", "CSS size", "public API", "race tests", "vet"}})
}

func (r runner) generator() error {
	want, err := r.output(r.goTool, "list", "-m", "-f", "{{.Version}}", "github.com/a-h/templ")
	if err != nil {
		return err
	}
	got, err := r.output(r.templTool, "version")
	if err != nil {
		return err
	}
	if strings.TrimPrefix(got, "v") != strings.TrimPrefix(want, "v") {
		return fmt.Errorf("templ generator %q does not match runtime %q", got, want)
	}
	return nil
}

func copyFile(from, to string) error {
	b, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	return os.WriteFile(to, b, 0644)
}

func (r runner) generated() error {
	if err := r.generator(); err != nil {
		return err
	}
	temp, err := os.MkdirTemp("", "tendui-generated-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	files, err := r.files(".templ")
	if err != nil {
		return err
	}
	for _, path := range append(append([]string{}, files...), "go.mod", "go.sum") {
		if err := copyFile(filepath.Join(r.root, path), filepath.Join(temp, path)); err != nil {
			return err
		}
	}
	if err := r.run("generated.log", temp, r.templTool, "generate"); err != nil {
		return err
	}
	for _, path := range files {
		generated := strings.TrimSuffix(path, ".templ") + "_templ.go"
		expected, err := os.ReadFile(filepath.Join(r.root, generated))
		if err != nil {
			return err
		}
		actual, err := os.ReadFile(filepath.Join(temp, generated))
		if err != nil {
			return err
		}
		if !bytes.Equal(expected, actual) {
			return fmt.Errorf("generated source differs: %s; run templ generate", generated)
		}
	}
	outputs, err := r.files("_templ.go")
	if err != nil {
		return err
	}
	if len(outputs) != len(files) {
		return errors.New("orphan generated Go files found; every _templ.go must have a .templ source")
	}
	return r.report("generated.json", map[string]any{"templates": len(files), "reproducible": true})
}

func (r runner) recipes() error {
	if err := r.generator(); err != nil {
		return err
	}
	temp, err := os.MkdirTemp("", "tendui-recipes-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	source, err := os.ReadFile(filepath.Join(r.root, "examples/site/catalog_samples.templ"))
	if err != nil {
		return err
	}
	names := recipes.Names(string(source))
	if len(names) == 0 {
		return errors.New("no catalog recipes found")
	}
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			return fmt.Errorf("duplicate recipe %s", name)
		}
		seen[name] = true
		snippet, err := recipes.Source(string(source), name)
		if err != nil {
			return err
		}
		dir := filepath.Join(temp, strings.ToLower(name))
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "example.templ"), []byte(snippet), 0644); err != nil {
			return err
		}
	}
	templVersion, err := r.output(r.goTool, "list", "-m", "-f", "{{.Version}}", "github.com/a-h/templ")
	if err != nil {
		return err
	}
	mod := fmt.Sprintf("module tendui-recipe-check\n\ngo 1.26.0\n\nrequire (\n %s v%s\n github.com/a-h/templ %s\n)\n\nreplace %s => %q\n", modulePath, tendui.Version, templVersion, modulePath, r.root)
	if err := os.WriteFile(filepath.Join(temp, "go.mod"), []byte(mod), 0644); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(r.root, "go.sum"), filepath.Join(temp, "go.sum")); err != nil {
		return err
	}
	if err := r.run("recipes-generate.log", temp, r.templTool, "generate"); err != nil {
		return err
	}
	if err := r.run("recipes-build.log", temp, r.goTool, "test", "-mod=mod", "./..."); err != nil {
		return err
	}
	return r.report("recipes.json", map[string]any{"independently_compiled": names, "published_install_test": false})
}

func (r runner) buildVersion() (string, error) {
	describe, err := r.output("git", "describe", "--tags", "--long", "--dirty", "--match", "v[0-9][0-9][0-9][0-9].*")
	if err != nil {
		revision, err := r.output("git", "rev-parse", "--short=12", "HEAD")
		if err != nil {
			return "", err
		}
		describe = tendui.Version + "-dev-g" + revision
	}
	dirty, err := r.output("git", "status", "--porcelain", "--untracked-files=normal")
	if err != nil {
		return "", err
	}
	if dirty != "" && !strings.HasSuffix(describe, "-dirty") {
		describe += "-dirty"
	}
	if !regexp.MustCompile(`^v?[0-9]{4}\.[1-9][0-9]*\.[0-9]+(?:-[0-9A-Za-z.-]+)?$`).MatchString(describe) {
		return "", fmt.Errorf("unexpected development version %q", describe)
	}
	return describe, nil
}

func (r runner) build() error {
	version, err := r.buildVersion()
	if err != nil {
		return err
	}
	output := filepath.Join(r.root, "tests/results/demos")
	if err := os.MkdirAll(output, 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(output, executableName(".tendui-build-*", runtime.GOOS))
	if err != nil {
		return err
	}
	temp := f.Name()
	defer os.Remove(temp)
	if err := f.Close(); err != nil {
		return err
	}
	logFile, err := os.Create(filepath.Join(r.results, "build.log"))
	if err != nil {
		return err
	}
	defer logFile.Close()
	cmd := r.command(r.root, r.goTool, "build", "-trimpath", "-ldflags", "-s -w -X main.version="+version, "-o", temp, "./examples/basic")
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Stdout, cmd.Stderr = logFile, logFile
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w; see build.log", err)
	}
	info, err := os.Stat(temp)
	if err != nil {
		return err
	}
	if info.Size() > demoBudget {
		return fmt.Errorf("demo %d exceeds %d-byte budget", info.Size(), demoBudget)
	}
	embedded, err := buildinfo.ReadFile(temp)
	if err != nil {
		return err
	}
	cgoDisabled := false
	for _, setting := range embedded.Settings {
		if setting.Key == "CGO_ENABLED" && setting.Value == "0" {
			cgoDisabled = true
		}
	}
	if !cgoDisabled {
		return errors.New("demo did not record CGO_ENABLED=0")
	}
	if runtime.GOOS == "linux" {
		binary, err := elf.Open(temp)
		if err != nil {
			return err
		}
		defer binary.Close()
		for _, program := range binary.Progs {
			if program.Type == elf.PT_INTERP {
				return errors.New("demo has a dynamic ELF interpreter")
			}
		}
	}
	got, err := r.output(temp, "-version")
	if err != nil || got != version {
		return fmt.Errorf("version smoke check: got %q, want %q: %v", got, version, err)
	}
	artifact, err := os.Open(temp)
	if err != nil {
		return err
	}
	hash := sha256.New()
	_, err = io.Copy(hash, artifact)
	artifact.Close()
	if err != nil {
		return err
	}
	name := executableName("tendui-demo", runtime.GOOS)
	if err := os.Rename(temp, filepath.Join(output, name)); err != nil {
		return err
	}
	return r.report("build.json", map[string]any{"version": version, "executable": name, "bytes": info.Size(), "budget_bytes": demoBudget, "sha256": fmt.Sprintf("%x", hash.Sum(nil)), "go": embedded.GoVersion, "os": runtime.GOOS, "arch": runtime.GOARCH, "cgo": false, "version_smoke": true, "stripped": true})
}

func (r runner) report(name string, v map[string]any) error {
	revision, err := r.output("git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	status, err := r.output("git", "status", "--porcelain", "--untracked-files=normal")
	if err != nil {
		return err
	}
	v["revision"], v["dirty"], v["checked_at"] = revision, status != "", time.Now().UTC()
	v["toolchain"] = runtime.Version()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(r.results, name), append(b, '\n'), 0644)
}

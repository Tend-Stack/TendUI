package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Tend-Stack/TendUI/v2026"
)

func proxyModulePath(path string) string {
	var b strings.Builder
	for _, c := range path {
		if c >= 'A' && c <= 'Z' {
			b.WriteByte('!')
			c += 'a' - 'A'
		}
		b.WriteRune(c)
	}
	return b.String()
}

func fileProxyURL(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path // Absolute Windows drive path.
	}
	return (&url.URL{Scheme: "file", Path: path}).String()
}

// Git includes explicit directories; Go's module ZIP hashes must describe the
// same files as the extracted cache. Preserve file bytes, omit directory entries.
func writeModuleArchive(source, destination string) (resultErr error) {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, archive.Close()) }()
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, file.Close()) }()
	writer := zip.NewWriter(file)
	defer func() { resultErr = errors.Join(resultErr, writer.Close()) }()
	for _, entry := range archive.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		if !entry.Mode().IsRegular() {
			return fmt.Errorf("unsupported module archive entry %q: %s", entry.Name, entry.Mode())
		}
		if err := writer.Copy(entry); err != nil {
			return err
		}
	}
	return nil
}

func (r runner) consumer() error {
	proof, err := r.consumerFixture()
	if err != nil {
		return err
	}
	return r.report("consumer.json", proof)
}

// This is a local versioned-module packaging proof, not a publication check.
// A clean committed tree prevents an archive from silently omitting current work.
func (r runner) consumerFixture() (proof map[string]any, resultErr error) {
	status, err := r.output("git", "status", "--porcelain", "--untracked-files=normal")
	if err != nil {
		return nil, err
	}
	if status != "" {
		return nil, errors.New("consumer requires a clean commit so the proxy archive matches the recorded source; preserve and commit current work first")
	}
	revision, err := r.output("git", "rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}
	cache, err := r.output(r.goTool, "env", "GOMODCACHE")
	if err != nil {
		return nil, err
	}
	temp, err := os.MkdirTemp("", "tendui-consumer-")
	if err != nil {
		return nil, err
	}
	moduleCache := filepath.Join(temp, "modules")
	proxy := filepath.Join(temp, "proxy")
	consumer := filepath.Join(temp, "consumer")
	command := func(args ...string) *exec.Cmd {
		cmd := r.command(consumer, r.goTool, args...)
		cmd.Env = append(cmd.Env,
			"GOMODCACHE="+moduleCache,
			"GOPROXY="+fileProxyURL(proxy)+","+fileProxyURL(filepath.Join(cache, "cache", "download")),
			"GOSUMDB=sum.golang.org", "GONOSUMDB="+modulePath, "GONOPROXY=none",
			"GOWORK=off", "GOFLAGS=", "GOTOOLCHAIN=local", "GOENV=off")
		return cmd
	}
	defer func() {
		// Go module-cache directories can be read-only. Let Go remove its owned cache.
		cmd := command("clean", "-modcache")
		cmd.Dir = temp
		if output, err := cmd.CombinedOutput(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("clean isolated module cache %s: %w: %s", temp, err, output))
			return
		}
		if err := os.RemoveAll(temp); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove owned consumer fixture %s: %w", temp, err))
		}
	}()
	version := "v" + tendui.Version
	versionDir := filepath.Join(proxy, filepath.FromSlash(proxyModulePath(modulePath)), "@v")
	for _, dir := range []string{consumer, versionDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	mod, err := os.ReadFile(filepath.Join(r.root, "go.mod"))
	if err != nil {
		return nil, err
	}
	sums, err := os.ReadFile(filepath.Join(r.root, "go.sum"))
	if err != nil {
		return nil, err
	}
	// Copy committed dependency checksums so cached third-party bytes are verified
	// without querying a checksum server. Only this unpublished module skips sumdb.
	files := map[string][]byte{
		filepath.Join(versionDir, version+".mod"):   mod,
		filepath.Join(versionDir, version+".info"):  []byte(fmt.Sprintf(`{"Version":%q,"Time":"2000-01-01T00:00:00Z"}`, version)),
		filepath.Join(consumer, "go.mod"):           []byte("module example.invalid/tendui-consumer\n\ngo 1.26.0\n"),
		filepath.Join(consumer, "go.sum"):           sums,
		filepath.Join(consumer, "consumer_test.go"): []byte(consumerTestSource),
	}
	// Build the documented onboarding application outside the library module,
	// using only the installed version and checked-in generated Go.
	quickstart := filepath.Join(consumer, "quickstart")
	if err := os.MkdirAll(quickstart, 0755); err != nil {
		return nil, err
	}
	for _, name := range []string{"main.go", "page.templ", "page_templ.go", "main_test.go"} {
		contents, err := os.ReadFile(filepath.Join(r.root, "examples/quickstart", name))
		if err != nil {
			return nil, err
		}
		files[filepath.Join(quickstart, name)] = contents
	}
	for path, contents := range files {
		if err := os.WriteFile(path, contents, 0644); err != nil {
			return nil, err
		}
	}
	zipPath := filepath.Join(versionDir, version+".zip")
	sourceArchive := filepath.Join(temp, "source.zip")
	if err := r.run("consumer-archive.log", r.root, "git", "archive", "--format=zip", "--prefix="+modulePath+"@"+version+"/", "--output="+sourceArchive, "HEAD"); err != nil {
		return nil, err
	}
	if err := writeModuleArchive(sourceArchive, zipPath); err != nil {
		return nil, err
	}
	log, err := os.Create(filepath.Join(r.results, "consumer.log"))
	if err != nil {
		return nil, err
	}
	defer log.Close()
	run := func(args ...string) error {
		fmt.Fprintln(log, "go", strings.Join(args, " "))
		cmd := command(args...)
		cmd.Stdout, cmd.Stderr = log, log
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("consumer %s: %w; see consumer.log", args[0], err)
		}
		return nil
	}
	// Resolve the imported packages and their runtime dependencies without replace.
	// A dependency's own test imports are outside this consumer's build contract.
	if err := run("get", modulePath+"@"+version, modulePath+"/components@"+version, modulePath+"/render@"+version); err != nil {
		return nil, err
	}
	if err := run("test", "-mod=readonly", "-count=1", "./..."); err != nil {
		return nil, err
	}
	if err := run("mod", "verify"); err != nil {
		return nil, err
	}
	data, err := command("list", "-m", "-json", modulePath).Output()
	if err != nil {
		return nil, fmt.Errorf("inspect installed consumer module: %w", err)
	}
	var installed struct {
		Path, Version, Sum string
		Replace            *json.RawMessage
	}
	if err := json.Unmarshal(data, &installed); err != nil {
		return nil, err
	}
	if installed.Path != modulePath || installed.Version != version || installed.Replace != nil || installed.Sum == "" {
		return nil, errors.New("consumer did not resolve the versioned module without replace")
	}
	// Detect a source change during the check before attaching evidence to HEAD.
	afterRevision, err := r.output("git", "rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}
	afterStatus, err := r.output("git", "status", "--porcelain", "--untracked-files=normal")
	if err != nil || afterRevision != revision || afterStatus != "" {
		return nil, errors.New("source changed during consumer verification; rerun from a clean commit")
	}
	zip, err := os.ReadFile(zipPath)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"module": installed.Path, "version": installed.Version, "module_sum": installed.Sum,
		"archive_sha256": fmt.Sprintf("%x", sha256.Sum256(zip)), "archive_bytes": len(zip),
		"source_revision": revision, "replace_used": false, "templ_generator_used": false,
		"quickstart_tested": true, "local_file_proxy": true, "published_install_test": false, "isolated_fixture_removed": true,
	}, nil
}

const consumerTestSource = `package consumer_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/Tend-Stack/TendUI/v2026"
    "github.com/Tend-Stack/TendUI/v2026/components"
    "github.com/Tend-Stack/TendUI/v2026/render"
)

func TestInstalledModule(t *testing.T) {
    mux := http.NewServeMux()
    mux.Handle(tendui.AssetPath, tendui.AssetHandler())
    css := httptest.NewRecorder()
    mux.ServeHTTP(css, httptest.NewRequest("GET", tendui.StylesheetURL, nil))
    if css.Code != 200 || !strings.Contains(css.Body.String(), ":root") || css.Header().Get("ETag") == "" {
        t.Fatal("installed module lost embedded stylesheet")
    }
    page := components.Layout(components.LayoutProps{Title: "Installed <TendUI>"},
        components.Card(components.CardProps{Body: components.Button(components.ButtonProps{Text: "Go <native>"})}))
    response := httptest.NewRecorder()
    if err := render.HTML(response, httptest.NewRequest("GET", "/", nil), 200, page, nil); err != nil { t.Fatal(err) }
    body := response.Body.String()
    for _, want := range []string{"<!doctype html>", "Installed &lt;TendUI&gt;", "Go &lt;native&gt;", tendui.StylesheetURL} {
        if !strings.Contains(body, want) { t.Errorf("installed components missing %q", want) }
    }
    if response.Code != 200 || strings.Contains(body, "<script") { t.Fatal("unexpected installed page") }
}
`

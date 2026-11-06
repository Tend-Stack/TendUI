package main

import (
	"archive/zip"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestModuleArchiveOmitsDirectoriesAndPreservesFiles(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.zip")
	f, err := os.Create(source)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for _, name := range []string{"example@v1.0.0/", "example@v1.0.0/static/", "example@v1.0.0/static/style.css"} {
		entry, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasSuffix(name, "/") {
			if _, err := io.WriteString(entry, ":root { color: red; }\n"); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(dir, "module.zip")
	if err := writeModuleArchive(source, destination); err != nil {
		t.Fatal(err)
	}
	archive, err := zip.OpenReader(destination)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	if len(archive.File) != 1 || archive.File[0].Name != "example@v1.0.0/static/style.css" {
		t.Fatal("unexpected module archive entries", archive.File)
	}
	contents, err := archive.File[0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer contents.Close()
	data, err := io.ReadAll(contents)
	if err != nil || string(data) != ":root { color: red; }\n" {
		t.Fatalf("archive content changed: %q, %v", data, err)
	}
	if err := writeModuleArchive(source, destination); err == nil {
		t.Fatal("overwrote an existing archive")
	}
}

func TestConsumerProxyLocations(t *testing.T) {
	if got := proxyModulePath(modulePath); got != "github.com/!tend-!stack/!tend!u!i/v2026" {
		t.Fatal(got)
	}
	path := filepath.Join(t.TempDir(), "cache with spaces # and %")
	u, err := url.Parse(fileProxyURL(path))
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.ToSlash(path)
	if !strings.HasPrefix(want, "/") {
		want = "/" + want
	}
	if u.Scheme != "file" || u.Host != "" || u.Path != want || u.RawQuery != "" || u.Fragment != "" {
		t.Fatalf("proxy URL changed filesystem location: %s", u)
	}
}

func TestConsumerRejectsDirtySourcesWithoutChangingThem(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		name := "untracked"
		if tracked {
			name = "tracked"
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			git := func(args ...string) {
				t.Helper()
				cmd := exec.Command("git", append([]string{"-c", "user.name=TendUI test", "-c", "user.email=test@example.invalid"}, args...)...)
				cmd.Dir = dir
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git %v: %v: %s", args, err, output)
				}
			}
			git("init", "-q")
			path := filepath.Join(dir, "source")
			if err := os.WriteFile(path, []byte("initial"), 0644); err != nil {
				t.Fatal(err)
			}
			git("add", "source")
			git("commit", "-qm", "Consumer guard fixture")
			if !tracked {
				path = filepath.Join(dir, "new-work")
			}
			const work = "preserve unfinished work"
			if err := os.WriteFile(path, []byte(work), 0644); err != nil {
				t.Fatal(err)
			}
			r := runner{root: dir, results: t.TempDir(), goTool: goExecutable()}
			if err := r.consumer(); err == nil || !strings.Contains(err.Error(), "clean commit") {
				t.Fatalf("dirty source accepted: %v", err)
			}
			contents, err := os.ReadFile(path)
			if err != nil || string(contents) != work {
				t.Fatalf("unfinished work changed: %s, %v", contents, err)
			}
			entries, err := os.ReadDir(r.results)
			if err != nil || len(entries) != 0 {
				t.Fatalf("dirty source produced evidence: %v, %v", entries, err)
			}
		})
	}
}

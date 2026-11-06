package site

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"strings"
	"testing"
)

type uploadPart struct{ name, filename, content string }

func uploadRequest(t *testing.T, parts ...uploadPart) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, part := range parts {
		var dst io.Writer
		var err error
		if part.filename == "" {
			dst, err = writer.CreateFormField(part.name)
		} else {
			dst, err = writer.CreateFormFile(part.name, part.filename)
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(dst, part.content); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/uploads", &body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r
}

func serveUpload(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, r)
	return w
}

func TestUploadNativeFlowAndIsolation(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	start := serveUpload(httptest.NewRequest("GET", "/uploads", nil))
	for _, want := range []string{`enctype="multipart/form-data"`, `name="projects"`, `type="file"`, `href="/uploads/sample.csv"`, `aria-current="page"`} {
		if start.Code != 200 || !strings.Contains(start.Body.String(), want) {
			t.Fatalf("upload page missing %q", want)
		}
	}
	sample := serveUpload(httptest.NewRequest("GET", "/uploads/sample.csv", nil))
	if sample.Code != 200 || sample.Header().Get("Content-Type") != "text/csv; charset=utf-8" || sample.Header().Get("Content-Disposition") != `attachment; filename="projects.csv"` || sample.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("sample download headers")
	}
	preview := serveUpload(uploadRequest(t, uploadPart{"projects", "sample.csv", sample.Body.String()}))
	for _, want := range []string{"Validation passed", "Project rows ready to inspect: 3.", "Validated project rows", `scope="row"`, `>Studio website</th>`, "Community portal"} {
		if preview.Code != 200 || !strings.Contains(preview.Body.String(), want) {
			t.Errorf("preview missing %q; status %d", want, preview.Code)
		}
	}
	if preview.Header().Get("Set-Cookie") != "" || preview.Header().Get("Cache-Control") != "no-store" || !strings.Contains(preview.Header().Get("Content-Security-Policy"), "default-src 'none'") || strings.Contains(preview.Body.String(), "<script") {
		t.Fatal("unexpected persistence or script contract")
	}
	fresh := serveUpload(httptest.NewRequest("GET", "/uploads", nil))
	if strings.Contains(fresh.Body.String(), "Studio website") || strings.Contains(fresh.Body.String(), "Your CSV preview") {
		t.Fatal("upload data survived the request")
	}
	files, err := os.ReadDir(tmp)
	if err != nil || len(files) != 0 {
		t.Fatal("upload left temporary files", err)
	}
}

func TestUploadCSVValidation(t *testing.T) {
	const header = "project,owner,status\n"
	for _, tc := range []struct {
		name, csv, want string
		code            int
	}{
		{"BOM CRLF and unicode", "\xef\xbb\xbfproject,owner,status\r\n  Café  ,Zoë,active\r\n", ">Café</th>", 200},
		{"quoted comma", header + "\"One, two\",Owner,planned\n", "One, two", 200},
		{"HTML escaping", header + "<script>alert(1)</script>,A&B,complete\n", "&lt;script&gt;alert(1)&lt;/script&gt;", 200},
		{"empty", "", "exact header", 422},
		{"header only", header, "at least one", 422},
		{"wrong header", "name,owner,status\nA,B,active\n", "exact header", 422},
		{"short header", "project,owner\n", "exact header", 422},
		{"extra field", header + "A,B,active,extra\n", "Line 2", 422},
		{"missing field", header + "A,B\n", "Line 2", 422},
		{"missing owner", header + "A, ,active\n", "Line 2", 422},
		{"bad status after valid row", header + "PrivateProject,B,active\nC,D,unknown\n", "Line 3", 422},
		{"uppercase status", header + "A,B,ACTIVE\n", "status must be", 422},
		{"bad quote", header + "\"A,B,active\n", "Line 2", 422},
		{"multiline", header + "\"A\nB\",C,active\n", "single-line", 422},
		{"NUL", header + "A\x00,B,active\n", "control characters", 422},
		{"invalid UTF8", header + "\xff,B,active\n", "UTF-8", 422},
		{"long field", header + strings.Repeat("é", 81) + ",B,active\n", "80 characters", 422},
		{"max field", header + strings.Repeat("é", 80) + ",B,active\n", "Validation passed", 200},
		{"50 rows", header + strings.Repeat("A,B,active\n", 50), "Project rows ready to inspect: 50.", 200},
		{"51 rows", header + strings.Repeat("A,B,active\n", 51), "up to 50", 422},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := uploadRequest(t, uploadPart{"projects", `../../<img onerror=alert(1)>.csv`, tc.csv})
			w := serveUpload(r)
			body := w.Body.String()
			if w.Code != tc.code || !strings.Contains(body, tc.want) {
				t.Fatalf("got status %d, wanted %d with %q", w.Code, tc.code, tc.want)
			}
			if strings.Contains(body, "<script>") || strings.Contains(body, "<img onerror") {
				t.Fatal("unescaped upload")
			}
			if tc.code != 200 {
				for _, want := range []string{`href="#project-csv"`, `aria-invalid="true"`, `aria-describedby="upload-help upload-error"`, "choose it again"} {
					if !strings.Contains(body, want) {
						t.Errorf("missing error binding %q", want)
					}
				}
				if strings.Contains(body, "Validated project rows") || strings.Contains(body, "PrivateProject") {
					t.Fatal("invalid upload leaked a partial preview")
				}
			}
		})
	}
}

func TestUploadMultipartAndRequestLimits(t *testing.T) {
	valid := uploadPart{"projects", "projects.csv", sampleProjectsCSV}
	for _, tc := range []struct {
		name  string
		parts []uploadPart
		code  int
	}{
		{"missing", nil, 400},
		{"text field", []uploadPart{{"projects", "", sampleProjectsCSV}}, 400},
		{"wrong name", []uploadPart{{"other", "a.csv", sampleProjectsCSV}}, 400},
		{"duplicate", []uploadPart{valid, valid}, 400},
		{"extra field", []uploadPart{valid, {"other", "", "value"}}, 400},
		{"file too big", []uploadPart{{"projects", "a.csv", strings.Repeat("x", uploadFileLimit+1)}}, 413},
		{"exact file limit", []uploadPart{{"projects", "a.csv", sampleProjectsCSV + strings.Repeat("\n", uploadFileLimit-len(sampleProjectsCSV))}}, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := serveUpload(uploadRequest(t, tc.parts...))
			if w.Code != tc.code {
				t.Fatalf("got %d want %d", w.Code, tc.code)
			}
		})
	}
	for _, size := range []int{uploadBodyLimit, uploadBodyLimit + 1} {
		r := uploadRequest(t, valid)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		body = append(body, bytes.Repeat([]byte("x"), size-len(body))...)
		r.Body = io.NopCloser(bytes.NewReader(body))
		r.ContentLength = -1 // Chunked/unknown sizes get the same enforced bound, including epilogue.
		w := serveUpload(r)
		want := 200
		if size > uploadBodyLimit {
			want = 413
		}
		if w.Code != want {
			t.Errorf("body %d: got %d want %d", size, w.Code, want)
		}
	}
	for _, contentType := range []string{"", "text/csv", "multipart/form-data", "multipart/form-data; boundary=bad"} {
		r := httptest.NewRequest("POST", "/uploads", strings.NewReader(sampleProjectsCSV))
		r.Header.Set("Content-Type", contentType)
		if w := serveUpload(r); w.Code != 400 {
			t.Errorf("content type %q: %d", contentType, w.Code)
		}
	}
	r := uploadRequest(t, valid)
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body[:len(body)-30]))
	if w := serveUpload(r); w.Code != 400 {
		t.Errorf("truncated body: %d", w.Code)
	}

	var encoded bytes.Buffer
	writer := multipart.NewWriter(&encoded)
	part, err := writer.CreatePart(textproto.MIMEHeader{"Content-Disposition": {`form-data; name="projects"; filename="a.csv"`}, "Content-Transfer-Encoding": {"quoted-printable"}})
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, sampleProjectsCSV)
	_ = writer.Close()
	r = httptest.NewRequest("POST", "/uploads", &encoded)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	if w := serveUpload(r); w.Code != 400 {
		t.Errorf("transfer encoding: %d", w.Code)
	}
}

func TestUploadMethodsAndOrigin(t *testing.T) {
	for _, path := range []string{"/uploads", "/uploads/sample.csv"} {
		w := serveUpload(httptest.NewRequest("HEAD", path, nil))
		if w.Code != 200 || w.Body.Len() != 0 {
			t.Fatalf("HEAD %s: %d body %d", path, w.Code, w.Body.Len())
		}
		w = serveUpload(httptest.NewRequest("DELETE", path, nil))
		if w.Code != 405 || w.Header().Get("Allow") == "" {
			t.Fatalf("DELETE %s: %d", path, w.Code)
		}
	}
	r := uploadRequest(t, uploadPart{"projects", "a.csv", sampleProjectsCSV})
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	if w := serveUpload(r); w.Code != 403 {
		t.Errorf("cross-origin upload: %d", w.Code)
	}
	if !demoPath("/uploads") || demoPath("/uploads/sample.csv") {
		t.Fatal("theme navigation allowlist")
	}
}

func TestUploadConcurrentIsolation(t *testing.T) {
	for i := range 8 {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			name := fmt.Sprintf("PrivateRow-%d", i)
			w := serveUpload(uploadRequest(t, uploadPart{"projects", "a.csv", "project,owner,status\n" + name + ",Owner,active\n"}))
			if w.Code != 200 || strings.Count(w.Body.String(), "PrivateRow-") != 1 || !strings.Contains(w.Body.String(), name) {
				t.Fatal("request rows mixed")
			}
		})
	}
}

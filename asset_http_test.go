package gothui

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func assetRequest(handler http.Handler, method string, headers map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, StylesheetURL, nil)
	for name, value := range headers {
		r.Header.Set(name, value)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func TestAssetContentValidatorsAndRevalidation(t *testing.T) {
	h := AssetHandler()
	full := assetRequest(h, "GET", nil)
	body := full.Body.Bytes()
	if full.Code != 200 || len(body) == 0 {
		t.Fatal("missing embedded stylesheet")
	}
	wantETag := fmt.Sprintf(`"%x"`, sha256.Sum256(body))
	if full.Header().Get("ETag") != wantETag {
		t.Fatal("validator is not derived from the exact served bytes")
	}
	for key, want := range map[string]string{
		"Cache-Control": "public, max-age=0, must-revalidate", "Content-Length": strconv.Itoa(len(body)),
		"Content-Type": "text/css; charset=utf-8", "X-Content-Type-Options": "nosniff", "Accept-Ranges": "bytes",
	} {
		if got := full.Header().Get(key); got != want {
			t.Errorf("%s: got %q want %q", key, got, want)
		}
	}
	if full.Header().Get("Last-Modified") != "" {
		t.Error("embedded assets must not invent a source modification date")
	}
	// Independent mounts of the same binary must identify the same content.
	if got := assetRequest(AssetHandler(), "GET", nil); got.Header().Get("ETag") != wantETag || !bytes.Equal(got.Body.Bytes(), body) {
		t.Fatal("handler instances disagree on stylesheet identity")
	}
	for _, tc := range []struct {
		name, validator string
		status          int
	}{
		{"current", wantETag, 304}, {"weak current", "W/" + wantETag, 304}, {"list", `"previous-build", ` + wantETag, 304},
		{"any", "*", 304}, {"old build", `"previous-build"`, 200},
	} {
		for _, method := range []string{"GET", "HEAD"} {
			t.Run(tc.name+"/"+method, func(t *testing.T) {
				w := assetRequest(h, method, map[string]string{"If-None-Match": tc.validator})
				if w.Code != tc.status || w.Header().Get("ETag") != wantETag {
					t.Fatalf("status/validator: %d %q", w.Code, w.Header().Get("ETag"))
				}
				if w.Header().Get("Cache-Control") != full.Header().Get("Cache-Control") {
					t.Error("revalidation policy lost")
				}
				if method == "HEAD" || tc.status == 304 {
					if w.Body.Len() != 0 {
						t.Error("body emitted for HEAD/304")
					}
				} else if !bytes.Equal(w.Body.Bytes(), body) {
					t.Error("stale validator did not receive current complete CSS")
				}
			})
		}
	}
	head := assetRequest(h, "HEAD", nil)
	if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("Content-Length") != strconv.Itoa(len(body)) {
		t.Fatal("HEAD does not describe the complete stylesheet")
	}
}

func TestAssetSingleRangesAndIfRange(t *testing.T) {
	h := AssetHandler()
	full := assetRequest(h, "GET", nil)
	body, etag := full.Body.Bytes(), full.Header().Get("ETag")
	for _, tc := range []struct {
		name, requested, ifRange string
		start, end, status       int
	}{
		{"bounded", "bytes=3-18", "", 3, 19, 206},
		{"open ended", "bytes=32-", "", 32, len(body), 206},
		{"suffix", "bytes=-12", "", len(body) - 12, len(body), 206},
		{"clamped end", fmt.Sprintf("bytes=10-%d", len(body)+100), "", 10, len(body), 206},
		{"current validator", "bytes=3-18", etag, 3, 19, 206},
		{"old validator", "bytes=3-18", `"previous-build"`, 0, len(body), 200},
		{"weak validator", "bytes=3-18", "W/" + etag, 0, len(body), 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := assetRequest(h, "GET", map[string]string{"Range": tc.requested, "If-Range": tc.ifRange})
			if w.Code != tc.status || !bytes.Equal(w.Body.Bytes(), body[tc.start:tc.end]) {
				t.Fatalf("range returned status %d or wrong bytes", w.Code)
			}
			if w.Header().Get("Content-Length") != strconv.Itoa(tc.end-tc.start) || w.Header().Get("ETag") != etag {
				t.Error("range identity or length mismatch")
			}
			wantRange := ""
			if tc.status == 206 {
				wantRange = fmt.Sprintf("bytes %d-%d/%d", tc.start, tc.end-1, len(body))
			}
			if w.Header().Get("Content-Range") != wantRange {
				t.Errorf("Content-Range: %q", w.Header().Get("Content-Range"))
			}
		})
	}
	w := assetRequest(h, "GET", map[string]string{"Range": fmt.Sprintf("bytes=%d-", len(body))})
	if w.Code != 416 || w.Header().Get("Content-Range") != fmt.Sprintf("bytes */%d", len(body)) {
		t.Fatal("unsatisfiable range did not expose current length")
	}
	w = assetRequest(h, "GET", map[string]string{"Range": "bytes=0-7", "If-None-Match": etag})
	if w.Code != 304 || w.Body.Len() != 0 {
		t.Fatal("matching cache validator must take precedence over Range")
	}
	w = assetRequest(h, "GET", map[string]string{"Range": "bytes=0-7", "If-Match": `"previous-build"`})
	if w.Code != 412 {
		t.Fatal("failed precondition must take precedence over Range")
	}
}

func TestAssetMultipartRanges(t *testing.T) {
	h := AssetHandler()
	full := assetRequest(h, "GET", nil)
	body := full.Body.Bytes()
	w := assetRequest(h, "GET", map[string]string{"Range": "bytes=0-7,32-47"})
	mediaType, params, err := mime.ParseMediaType(w.Header().Get("Content-Type"))
	if err != nil || w.Code != 206 || mediaType != "multipart/byteranges" || params["boundary"] == "" {
		t.Fatalf("multipart range metadata: %d %v", w.Code, err)
	}
	reader := multipart.NewReader(bytes.NewReader(w.Body.Bytes()), params["boundary"])
	for _, span := range [][2]int{{0, 8}, {32, 48}} {
		part, err := reader.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		got, err := io.ReadAll(part)
		if err != nil || !bytes.Equal(got, body[span[0]:span[1]]) {
			t.Fatal("wrong multipart bytes", err)
		}
		if part.Header.Get("Content-Range") != fmt.Sprintf("bytes %d-%d/%d", span[0], span[1]-1, len(body)) || part.Header.Get("Content-Type") != "text/css; charset=utf-8" {
			t.Fatal("wrong multipart headers")
		}
	}
	if _, err := reader.NextPart(); err != io.EOF {
		t.Fatal("unexpected extra part or incomplete boundary", err)
	}
	if w.Header().Get("ETag") != full.Header().Get("ETag") || w.Header().Get("Content-Length") != strconv.Itoa(w.Body.Len()) {
		t.Fatal("multipart identity or length mismatch")
	}
}

func TestAssetConcurrentReaders(t *testing.T) {
	h := AssetHandler()
	full := assetRequest(h, "GET", nil)
	body := bytes.Clone(full.Body.Bytes())
	for i := range 32 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			start := i * 19
			w := assetRequest(h, "GET", map[string]string{"Range": fmt.Sprintf("bytes=%d-%d", start, start+31)})
			if w.Code != 206 || !bytes.Equal(w.Body.Bytes(), body[start:start+32]) {
				t.Fatal("concurrent reader position leaked between requests")
			}
		})
	}
}

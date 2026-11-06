package render

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type componentFunc func(context.Context, io.Writer) error

func (f componentFunc) Render(ctx context.Context, w io.Writer) error { return f(ctx, w) }

func content(s string) Component {
	return componentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	})
}

func TestPageFragmentAndHistory(t *testing.T) {
	for _, tc := range []struct {
		name, hx, history, boosted, want string
	}{
		{name: "native", want: "page"},
		{name: "fragment", hx: "true", want: "fragment"},
		{name: "false", hx: "false", want: "page"},
		{name: "history with hx", hx: "true", history: "true", want: "page"},
		{name: "history alone", history: "true", want: "page"},
		{name: "boosted", hx: "true", boosted: "true", want: "page"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("HX-Request", tc.hx)
			r.Header.Set("HX-History-Restore-Request", tc.history)
			r.Header.Set("HX-Boosted", tc.boosted)
			w := httptest.NewRecorder()
			w.Header().Add("Vary", "Accept-Encoding, hx-request")
			w.Header().Set("Cache-Control", "private, no-store")
			if err := HTML(w, r, 201, content("page"), content("fragment")); err != nil {
				t.Fatal(err)
			}
			if w.Code != 201 || w.Body.String() != tc.want || w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
				t.Fatalf("unexpected response: %d %q", w.Code, w.Body.String())
			}
			v := strings.ToLower(strings.Join(w.Header().Values("Vary"), ","))
			for _, key := range []string{"accept-encoding", "hx-request", "hx-history-restore-request", "hx-boosted"} {
				if strings.Count(v, key) != 1 {
					t.Errorf("incorrect Vary %q", v)
				}
			}
			if w.Header().Get("Cache-Control") != "private, no-store" {
				t.Fatal("cache policy lost")
			}
		})
	}
}

func TestHeadFallbackAndVaryWildcard(t *testing.T) {
	for _, method := range []string{"GET", "HEAD"} {
		w := httptest.NewRecorder()
		w.Header().Set("Vary", "*")
		r := httptest.NewRequest(method, "/", nil)
		r.Header.Set("HX-Request", "true")
		if err := HTML(w, r, 422, content("page"), nil); err != nil {
			t.Fatal(err)
		}
		want := "page"
		if method == "HEAD" {
			want = ""
		}
		if w.Code != 422 || w.Body.String() != want || strings.Join(w.Header().Values("Vary"), ",") != "*" {
			t.Fatalf("bad HEAD/fallback/wildcard: %+v", w)
		}
	}
}

func TestErrorsNeverLeakPartialHTML(t *testing.T) {
	fail := componentFunc(func(_ context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "private partial markup")
		return errors.New("private database failure")
	})
	for _, method := range []string{"GET", "HEAD"} {
		w := httptest.NewRecorder()
		w.Header().Set("Content-Length", "100")
		err := HTML(w, httptest.NewRequest(method, "/", nil), 200, fail, nil)
		if err == nil || w.Code != 500 || strings.Contains(w.Body.String(), "private") || w.Header().Get("Content-Length") != "" {
			t.Fatalf("unsafe render failure: %v %+v", err, w)
		}
		if method == "HEAD" && w.Body.Len() != 0 {
			t.Fatal("HEAD error body")
		}
	}
	for _, status := range []int{0, 199, 204, 205, 304, 600} {
		w := httptest.NewRecorder()
		if HTML(w, httptest.NewRequest("GET", "/", nil), status, content("page"), nil) == nil || w.Code != 500 {
			t.Fatalf("accepted invalid status %d", status)
		}
	}
	w := httptest.NewRecorder()
	if HTML(w, httptest.NewRequest("GET", "/", nil), 200, nil, nil) == nil || w.Code != 500 {
		t.Fatal("accepted missing page")
	}
}

type brokenWriter struct {
	header http.Header
	status int
}

func (w *brokenWriter) Header() http.Header       { return w.header }
func (w *brokenWriter) WriteHeader(s int)         { w.status = s }
func (w *brokenWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestContextAndTransportErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	err := HTML(httptest.NewRecorder(), r, 200, componentFunc(func(ctx context.Context, _ io.Writer) error { return ctx.Err() }), nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("context lost: %v", err)
	}
	w := &brokenWriter{header: http.Header{}}
	if err := HTML(w, r, 200, content("page"), nil); !errors.Is(err, io.ErrClosedPipe) || w.status != 200 {
		t.Fatalf("write error lost: %v", err)
	}
}

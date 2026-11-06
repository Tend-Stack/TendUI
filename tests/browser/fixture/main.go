// This consumer fixture is test tooling, never part of the public demo server.
package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Tend-Stack/GothUI/v2026/examples/site"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8089", "test listen address")
	runtime := flag.String("htmx", "tests/browser/node_modules/htmx.org/dist/htmx.min.js", "consumer-supplied test runtime")
	flag.Parse()
	demo := site.Handler()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/__test/htmx.js" {
			w.Header().Set("Content-Type", "text/javascript")
			http.ServeFile(w, r, *runtime)
			return
		}
		result := httptest.NewRecorder()
		demo.ServeHTTP(result, r)
		for k, values := range result.Header() {
			w.Header()[k] = values
		}
		body := result.Body.Bytes()
		if c, err := r.Cookie("gothui_test_htmx"); err == nil && c.Value == "1" {
			w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; connect-src 'self'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
			body = bytes.Replace(body, []byte("</head>"), []byte(`<meta name="htmx-config" content='{"allowEval":false,"allowScriptTags":false,"includeIndicatorStyles":false}'><script src="/__test/htmx.js" defer></script></head>`), 1)
			w.Header().Del("Content-Length")
		}
		w.WriteHeader(result.Code)
		if r.Method != http.MethodHead {
			_, _ = w.Write(body)
		}
	})
	s := &http.Server{Addr: *addr, Handler: h, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(s.ListenAndServe())
}

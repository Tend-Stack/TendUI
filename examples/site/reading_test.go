package site

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadingPatternNativeControlsAndContent(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ query, measure string }{{"", "standard"}, {"&measure=narrow", "narrow"}, {"&measure=wide", "wide"}, {"&measure=%3Cscript%3E", "standard"}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/patterns?section=reading"+tc.query, nil))
		html := w.Body.String()
		for _, want := range []string{`data-measure="` + tc.measure + `"`, `name="section" value="reading"`, `action="/patterns" method="get"`, `name="measure"`, `value="` + tc.measure + `" selected`, `id="reading-title"`, `<blockquote>`, `<dt>Measure</dt>`, `<figcaption>`, `data-prose="off"`, `aria-label="Go route example"`, `href="/catalog/prose"`} {
			if w.Code != 200 || !strings.Contains(html, want) {
				t.Errorf("reading %q: missing %s (status %d)", tc.query, want, w.Code)
			}
		}
		if strings.Contains(html, "<script") || strings.Contains(html, "onclick=") {
			t.Fatal("unexpected reading script")
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("HEAD", "/patterns?section=reading", nil))
	if w.Code != 200 || w.Body.Len() != 0 {
		t.Fatal("reading HEAD must have no body")
	}
}

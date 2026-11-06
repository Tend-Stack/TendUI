package main

import (
	"net/http/httptest"
	"testing"
)

// Includes request handling, rendering, buffering and a captured response, not TCP or a browser.
func BenchmarkDemoHandler(b *testing.B) {
	handler := newHandler()
	for _, path := range []string{"/dashboard", "/catalog", "/catalog/field", "/wizard", "/editing"} {
		b.Run(path, func(b *testing.B) {
			request := httptest.NewRequest("GET", path, nil)
			probe := httptest.NewRecorder()
			handler.ServeHTTP(probe, request)
			b.SetBytes(int64(probe.Body.Len()))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, request)
				if w.Code != 200 {
					b.Fatal(w.Code)
				}
			}
		})
	}
}

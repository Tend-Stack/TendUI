package tendui

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkAssetHandler(b *testing.B) {
	handler := AssetHandler()
	probe := httptest.NewRecorder()
	handler.ServeHTTP(probe, httptest.NewRequest("GET", StylesheetURL, nil))
	for _, name := range []string{"GET", "NotModified"} {
		b.Run(name, func(b *testing.B) {
			request := httptest.NewRequest("GET", StylesheetURL, nil)
			want := 200
			if name == "NotModified" {
				request.Header.Set("If-None-Match", probe.Header().Get("ETag"))
				want = 304
			} else {
				b.SetBytes(int64(probe.Body.Len()))
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, request)
				if w.Code != want {
					b.Fatal(w.Code)
				}
			}
		})
	}
}

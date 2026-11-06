package site

import (
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

// Read the actual hidden inputs emitted by templ to exercise the form-state round trip.
var hiddenInput = regexp.MustCompile(`<input type="hidden" name="([^"]+)" value="([^"]*)"`)

func wizardHidden(body string) url.Values {
	values := url.Values{}
	for _, m := range hiddenInput.FindAllStringSubmatch(body, -1) {
		values.Add(m[1], html.UnescapeString(m[2]))
	}
	return values
}
func wizardPost(h http.Handler, values url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", "/wizard", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func validBrief(stage, action string) url.Values {
	return url.Values{"stage": {stage}, "action": {action}, "name": {"Studio <sample>"}, "slug": {"studio-sample"}, "summary": {"A sample & a plan"}, "visibility": {"private"}, "accent": {"indigo"}}
}

func TestWizardNativeRoundTripAndBackPreservesEdits(t *testing.T) {
	h := newHandler()
	start := httptest.NewRecorder()
	h.ServeHTTP(start, httptest.NewRequest("GET", "/wizard", nil))
	values := wizardHidden(start.Body.String())
	values.Set("name", "Studio <sample>")
	values.Set("slug", "studio-sample")
	values.Set("summary", "A sample & a plan\nSecond line")
	values.Set("action", "next")
	step2 := wizardPost(h, values)
	if step2.Code != 200 || !strings.Contains(step2.Body.String(), "Choose how it should feel") {
		t.Fatalf("details transition: %d %s", step2.Code, step2.Body.String())
	}
	if strings.Index(step2.Body.String(), ">Continue</button>") > strings.Index(step2.Body.String(), "Back — keep my edits") {
		t.Fatal("Enter should default to forward submission, not Back")
	}
	values = wizardHidden(step2.Body.String())
	if values.Get("name") != "Studio <sample>" || values.Get("summary") != "A sample & a plan\nSecond line" {
		t.Fatal("details lost in hidden state")
	}
	values.Set("visibility", "public")
	values.Set("accent", "rose")
	values.Set("updates", "on")
	values.Set("action", "back")
	back := wizardPost(h, values)
	if back.Code != 200 || !strings.Contains(back.Body.String(), `value="Studio &lt;sample&gt;"`) {
		t.Fatal("Back did not preserve text edits")
	}
	values = wizardHidden(back.Body.String())
	if values.Get("accent") != "rose" || values.Get("updates") != "on" {
		t.Fatal("Back lost later preferences")
	}
	values.Set("name", "Revised <sample>")
	values.Set("slug", "revised-sample")
	values.Set("summary", "Revised & saved in form")
	values.Set("action", "next")
	step2 = wizardPost(h, values)
	if !strings.Contains(step2.Body.String(), `value="rose" selected`) || !strings.Contains(step2.Body.String(), " checked") {
		t.Fatal("preferences did not reappear")
	}
	values = wizardHidden(step2.Body.String())
	values.Set("visibility", "public")
	values.Set("accent", "rose")
	values.Set("updates", "on")
	values.Set("action", "next")
	review := wizardPost(h, values)
	if review.Code != 200 || !strings.Contains(review.Body.String(), "Take one last look") {
		t.Fatal("review transition failed")
	}
	values = wizardHidden(review.Body.String())
	values.Set("action", "finish")
	done := wizardPost(h, values)
	if done.Code != 200 || !strings.Contains(done.Body.String(), "Your preview is complete") || !strings.Contains(done.Body.String(), "Revised &lt;sample&gt;") || strings.Contains(done.Body.String(), "Revised <sample>") || len(done.Result().Cookies()) != 0 {
		t.Fatal("completion or escaping failed")
	}
	fresh := httptest.NewRecorder()
	h.ServeHTTP(fresh, httptest.NewRequest("GET", "/wizard?stage=3", nil))
	if !strings.Contains(fresh.Body.String(), `name="stage" value="1"`) || strings.Contains(fresh.Body.String(), "Revised") {
		t.Fatal("draft leaked into a new request")
	}
}

func TestWizardValidationAndHiddenStateTampering(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ key, value, stage, wantID string }{
		{"name", "", "1", "brief-name"}, {"slug", "../../bad", "1", "brief-slug"},
		{"summary", strings.Repeat("x", 401), "3", "brief-summary"},
		{"name", "", "3", "brief-name"}, {"visibility", "unexpected", "3", "brief-visibility"},
		{"accent", "<script>", "3", "brief-accent"},
	} {
		action := "next"
		if tc.stage == "3" {
			action = "finish"
		}
		values := validBrief(tc.stage, action)
		values.Set(tc.key, tc.value)
		w := wizardPost(h, values)
		if w.Code != 422 || !strings.Contains(w.Body.String(), `href="#`+tc.wantID+`"`) || !strings.Contains(w.Body.String(), `id="`+tc.wantID+`-error"`) || !strings.Contains(w.Body.String(), `aria-invalid="true"`) || strings.Contains(w.Body.String(), "Your preview is complete") {
			t.Fatalf("validation failed for %s: %d", tc.key, w.Code)
		}
	}
	// Back must retain incomplete preferences, without forcing validation first.
	values := validBrief("2", "back")
	values.Set("visibility", "")
	w := wizardPost(h, values)
	if w.Code != 200 || wizardHidden(w.Body.String()).Get("stage") != "1" {
		t.Fatal("Back was blocked by invalid current fields")
	}
}

func TestWizardRejectsMalformedAndCrossOriginRequests(t *testing.T) {
	h := newHandler()
	for _, body := range []string{
		"stage=0&action=next", "stage=1&action=finish", "stage=1&action=back", "stage=3&action=next", "stage=2&action=unknown",
		"stage=1&stage=3&action=next", "stage=1&action=next&name=%FF", "stage=1&action=next&updates=invalid", "stage=%ZZ",
	} {
		r := httptest.NewRequest("POST", "/wizard", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 400 {
			t.Errorf("accepted malformed state %q: %d", body, w.Code)
		}
	}
	r := httptest.NewRequest("POST", "/wizard", strings.NewReader("summary="+strings.Repeat("x", 9000)))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 413 {
		t.Fatalf("oversize: %d", w.Code)
	}
	r = httptest.NewRequest("POST", "/wizard", strings.NewReader(validBrief("1", "next").Encode()))
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross-origin request accepted")
	}
	for _, method := range []string{"HEAD", "PUT"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, "/wizard", nil))
		if method == "HEAD" && (w.Code != 200 || w.Body.Len() != 0) {
			t.Fatal("HEAD contract")
		}
		if method == "PUT" && w.Code != 405 {
			t.Fatal("method contract")
		}
	}
}

func TestWizardFieldBindingsAndNoClientRuntime(t *testing.T) {
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, httptest.NewRequest("GET", "/wizard", nil))
	for _, want := range []string{`for="brief-name"`, `id="brief-name"`, `aria-describedby="brief-name-description"`, `id="brief-slug-description"`, `aria-hidden="true"`, "formnovalidate", " required"} {
		if !strings.Contains(w.Body.String(), want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(w.Body.String(), "<script") || !strings.Contains(w.Header().Get("Content-Security-Policy"), "default-src 'none'") {
		t.Fatal("client runtime contract changed")
	}
}

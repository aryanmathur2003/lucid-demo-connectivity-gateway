package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer() http.Handler {
	return NewServer(NewMemoryStore(SeedVehicles())).Routes()
}

func get(t *testing.T, h http.Handler, path, acceptLanguage string) (*httptest.ResponseRecorder, connectivityResponse) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if acceptLanguage != "" {
		req.Header.Set("Accept-Language", acceptLanguage)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var body connectivityResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	return rec, body
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestConnectivityDefaultsToEnglish(t *testing.T) {
	rec, body := get(t, newTestServer(), "/v1/vehicles/5LDEMAA0000000001/connectivity", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if body.Locale != "en-US" || body.Message != "Your vehicle is connected." {
		t.Fatalf("got locale %q message %q", body.Locale, body.Message)
	}
	if body.Link != LinkOnline || body.SignalBars != 4 {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestConnectivityNegotiatesLocale(t *testing.T) {
	cases := []struct {
		header string
		locale string
	}{
		{"fr-CA,fr;q=0.9,en;q=0.5", "fr-CA"},
		{"de-AT", "de"},
		{"ar", "ar"},
		{"ja,ko;q=0.8", "en-US"},
		{"not a header;;", "en-US"},
	}
	for _, c := range cases {
		_, body := get(t, newTestServer(), "/v1/vehicles/5LDEMAA0000000002/connectivity", c.header)
		if body.Locale != c.locale {
			t.Errorf("Accept-Language %q: locale = %q, want %q", c.header, body.Locale, c.locale)
		}
		if body.Message == "" {
			t.Errorf("Accept-Language %q: empty message", c.header)
		}
	}
}

func TestConnectivityErrors(t *testing.T) {
	h := newTestServer()
	cases := []struct {
		path string
		want int
	}{
		{"/v1/vehicles/short/connectivity", http.StatusBadRequest},
		{"/v1/vehicles/5LDEMAA0000000099/connectivity", http.StatusNotFound},
		{"/v1/vehicles/5LDEMAA0000000001/other", http.StatusNotFound},
	}
	for _, c := range cases {
		rec, _ := get(t, h, c.path, "")
		if rec.Code != c.want {
			t.Errorf("%s: status = %d, want %d", c.path, rec.Code, c.want)
		}
	}
}

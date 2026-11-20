package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const goodBody = `{"homePrice":400000,"downPercent":20,"annualRate":6.5,"termYears":30,"taxPercent":1.1,"annualIns":1500}`

func post(body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, req)
	return rec
}

func TestAPICalculateOK(t *testing.T) {
	rec := post(goodBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var r Result
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if !close(r.Monthly.Total, 2514.28) {
		t.Errorf("monthly total = %.2f, want 2514.28", r.Monthly.Total)
	}
}

func TestAPIErrors(t *testing.T) {
	cases := []struct {
		name, body string
		want       int
	}{
		{"malformed JSON", `{"homePrice":`, http.StatusBadRequest},
		{"wrong type", `{"homePrice":"abc"}`, http.StatusBadRequest},
		{"unknown field", `{"price":1}`, http.StatusBadRequest},
		{"invalid values", strings.Replace(goodBody, `"downPercent":20`, `"downPercent":100`, 1), http.StatusUnprocessableEntity},
	}
	for _, c := range cases {
		rec := post(c.body)
		if rec.Code != c.want {
			t.Errorf("%s: status = %d, want %d", c.name, rec.Code, c.want)
		}
		if !strings.Contains(rec.Body.String(), `"error"`) {
			t.Errorf("%s: body should contain an error message, got %s", c.name, rec.Body.String())
		}
	}
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest("GET", "/api/health", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestFrontendFilesServed(t *testing.T) {
	for _, path := range []string{"/", "/style.css", "/app.js"} {
		rec := httptest.NewRecorder()
		newRouter().ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, rec.Code)
		}
	}
}

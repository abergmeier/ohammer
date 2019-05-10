package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testLocationMatch(t *testing.T, resp *http.Response, redirLoc string) {
	loc, ok := resp.Header["Location"]
	if !ok {
		t.Fatalf("Missing Location Header")
	}

	if len(loc) != 1 {
		t.Fatalf("Unexpected count of Location Headers: %v <-> 1", len(loc))
	}

	if loc[0] != redirLoc {
		t.Fatalf(`Incorrect redirect
to       %v
expected %v`, loc[0], redirLoc)
	}
}

func testStatus(t *testing.T, resp *http.Response, status int) {
	if resp.StatusCode != status {
		t.Fatalf("Unexpected status code %v <-> %v", resp.StatusCode, status)
	}
}

func testRedirect(t *testing.T, resp *http.Response, redirLoc string) {
	testStatus(t, resp, http.StatusTemporaryRedirect)
	testLocationMatch(t, resp, redirLoc)
}

func TestV2(t *testing.T) {

	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/v2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code %v", resp.StatusCode)
	}
}

func TestPullImageManifestPatched(t *testing.T) {
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/v2/gcr.io/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	testStatus(t, resp, http.StatusOK)
}

func TestPullImageManifestUnpatched(t *testing.T) {
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/v2/foobar/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	testRedirect(t, resp, "//foobar/v2/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516")
}

func TestExistanceImageManifestPatched(t *testing.T) {
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("HEAD", "/v2/gcr.io/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	testStatus(t, resp, http.StatusOK)
}

func TestExistanceImageManifestUnpatched(t *testing.T) {
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("HEAD", "/v2/spider/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	testRedirect(t, resp, "//spider/v2/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516")
}

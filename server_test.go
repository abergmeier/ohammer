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
		t.Fatalf("Uncorrect redirect to %v <-> %v", loc[0], redirLoc)
	}
}

func testRedirect(t *testing.T, resp *http.Response, redirLoc string) {
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("Unexpected status code %v <-> %v", resp.StatusCode, http.StatusTemporaryRedirect)
	}

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
	testRedirect(t, resp, "//gcr.io/v2/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516")
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
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("Unexpected status code %v", resp.StatusCode)
	}
}

func TestExistanceImageManifest(t *testing.T) {
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("HEAD", "/v2/gcr.io/spinnaker-marketplace/gate/manifests/1.2.1-20181108172516", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("Unexpected status code %v", resp.StatusCode)
	}
}

// gcr.io/spinnaker-marketplace/gate:1.2.1-20181108172516

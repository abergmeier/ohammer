package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	patchedManifestPaths = []string{
		"/v2/golang/manifests/1.12",
		"/v2/index.docker.io/golang/manifests/1.12",
		"/v2/index.docker.io/library/golang/manifests/1.12",
		"/v2/docker.io/golang/manifests/1.12",
		"/v2/docker.io/library/golang/manifests/1.12",
	}
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

func testStatus(t *testing.T, resp *http.Response, status int, id string) {
	if resp.StatusCode != status {
		t.Errorf("Unexpected status code %v <-> %v for %s", resp.StatusCode, status, id)
	}
}

func testRedirect(t *testing.T, resp *http.Response, redirLoc string, id string) {
	testStatus(t, resp, http.StatusTemporaryRedirect, id)
	testLocationMatch(t, resp, redirLoc)
}

func TestV2(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	for _, p := range patchedManifestPaths {
		req := httptest.NewRequest("GET", p, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		resp := w.Result()
		testStatus(t, resp, http.StatusOK, p)
	}
}

func TestPullImageManifestUnpatched(t *testing.T) {
	t.Parallel()
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	p := "/v2/foobar/manifests/1.12"
	req := httptest.NewRequest("GET", p, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	testRedirect(t, resp, "//docker.io/v2/library/foobar/manifests/1.12", p)
}

func TestExistanceImageManifestPatched(t *testing.T) {
	t.Parallel()
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	for _, p := range patchedManifestPaths {
		req := httptest.NewRequest("HEAD", p, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		resp := w.Result()
		testStatus(t, resp, http.StatusOK, p)
	}
}

func TestExistanceImageManifestUnpatched(t *testing.T) {
	t.Parallel()
	r, err := NewRouter()

	if err != nil {
		t.Fatal(err)
	}

	p := "/v2/spider/manifests/1.12"
	req := httptest.NewRequest("HEAD", p, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	testRedirect(t, resp, "//docker.io/v2/library/spider/manifests/1.12", p)
}

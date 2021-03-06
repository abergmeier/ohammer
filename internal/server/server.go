package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/otto-de/ohammer/internal"
	"github.com/otto-de/ohammer/internal/backing"
	"github.com/otto-de/ohammer/internal/build"
	"github.com/otto-de/ohammer/internal/config"

	"github.com/otto-de/ohammer/third_party/docker_distribution/reference"
)

func NewRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	err := r.HandleFunc("/v2", internal.HandleApiVersionCheck).GetError()
	if err != nil {
		return nil, err
	}

	sub := r.PathPrefix("/v2").Subrouter()
	err = sub.HandleFunc("/{origin:.*}/manifests/{reference}", proxyGetManifestHandler).Methods("GET", "HEAD").GetError()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = sub.HandleFunc("/{origin}/blobs/{digest}", proxyGetBlobHandler).Methods("GET").GetError()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewServer(addr string, repository string) (*http.Server, error) {

	r, err := NewRouter()
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}, nil
}

func findMatch(t *config.Target) *config.Patch {

	i := func() int {
		dockerReg := []byte(fmt.Sprintf("%s/%s:%s", t.Host, t.Path, t.Ref))
		for i, patch := range config.Patches {
			ok := patch.Reg.Find(dockerReg)
			if ok != nil {
				return i
			}
		}

		return -1
	}()

	if i == -1 {
		return nil
	}

	return &config.Patches[i]
}

func extractTarget(req *http.Request) *config.Target {

	vars := mux.Vars(req)
	orig, ok := vars["origin"]
	if !ok {
		return nil
	}

	t := &config.Target{}
	t.Host, t.Path = reference.SplitDockerDomain(orig)

	t.Ref, ok = vars["reference"]
	if !ok {
		return nil
	}

	return t
}

func redirectRequestToTarget(req *http.Request, resp http.ResponseWriter, sectionPath string, t *config.Target) {
	redirURL := req.URL
	redirURL.Host = t.Host
	redirURL.Path = fmt.Sprintf("/v2/%s/%s/%s", t.Path, sectionPath, t.Ref)

	http.Redirect(resp, req, redirURL.String(), http.StatusTemporaryRedirect)
}

func proxyGetManifestHandler(resp http.ResponseWriter, req *http.Request) {

	t := extractTarget(req)

	if t == nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	patch := findMatch(t)
	if patch == nil {
		// No match
		redirectRequestToTarget(req, resp, "manifests", t)
		return
	}

	b := backing.NewBacking()

	backingResponse, err := b.Poll(req, "manifests", t)
	if err != nil {
		panic(err)
	}

	defer backingResponse.Body.Close()

	if backingResponse.StatusCode != http.StatusOK {
		switch backingResponse.StatusCode {
		case http.StatusUnauthorized:
			resp.WriteHeader(http.StatusUnauthorized)
			return
		case http.StatusNotFound:
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		panic(backingResponse.StatusCode)
	}

	err = build.ApplyPatch(patch)
	if err != nil {
		panic(err)
	}
	// Recheck
	backingResponse, err = b.Poll(req, "manifests", t)
	if err != nil {
		panic(err)
	}

	if backingResponse.StatusCode != http.StatusOK {
		panic(backingResponse.Body)
	}
}

func proxyGetBlobHandler(resp http.ResponseWriter, req *http.Request) {

}

package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type match struct {
	reg *regexp.Regexp
	patch string
}

var (
	matchers = []match {
		match{
			reg: regexp.MustCompile(`.*gcr\.io/spinnaker-marketplace/gate.*`),
			patch: "mypatch",
		},
	}
)

func NewRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	err := r.HandleFunc("/v2", handleApiVersionCheck).GetError()
	if err != nil {
		return nil, err
	}

	sub := r.PathPrefix("/v2").Subrouter()
	err = sub.HandleFunc("/{originHost}/{originPath:.*?}/manifests/{reference}", proxyGetManifestHandler).Methods("GET").GetError()
	if err != nil {
		return nil, err
	}
	err = sub.HandleFunc("/{originHost}/{originPath:.*?}/manifests/{reference}", proxyHeadManifestHandler).Methods("HEAD").GetError()
	if err != nil {
		return nil, err
	}
	err = sub.HandleFunc("/{originHost}/{originPath:.*?}/blobs/{digest}", proxyGetBlobHandler).Methods("GET").GetError()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewServer(addr string) (*http.Server, error) {

	r, err := NewRouter()
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}, nil
}

type target struct {
	host string
	path string
	ref string
}

func findMatch(t *target) *match {

	i := func() int {
		dockerReg := []byte(fmt.Sprintf("%s/%s:%s", t.host, t.path, t.ref))
		for i, match := range matchers {
			ok := match.reg.Find(dockerReg)
			if ok != nil {
				return i
			}
		}

		return -1
	}()

	if i == -1 {
		return nil
	}

	return &matchers[i]
}

func extractTarget(req *http.Request) *target {

	vars := mux.Vars(req)
	host, ok := vars["originHost"]
	if !ok {
		return nil
	}

	path, ok := vars["originPath"]
	if !ok {
		return nil
	}

	ref, ok := vars["reference"]
	if !ok {
		return nil
	}

	return &target{
		host: host,
		path: path,
		ref:  ref,
	}
}

func redirect(req *http.Request, resp http.ResponseWriter, t *target) {
	redirURL := req.URL
	redirURL.Host = t.host
	redirURL.Path = fmt.Sprintf("/v2/%s/manifests/%s", t.path, t.ref)

	http.Redirect(resp, req, redirURL.String(), http.StatusTemporaryRedirect)
}

func proxyGetManifestHandler(resp http.ResponseWriter, req *http.Request) {

	t := extractTarget(req)

	if t == nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	match := findMatch(t)
	if match == nil {
		// No match
		redirect(req, resp, t)
		return
	}

	resp.WriteHeader(http.StatusNotImplemented)
}

func proxyHeadManifestHandler(resp http.ResponseWriter, req *http.Request) {

	t := extractTarget(req)

	if t == nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	match := findMatch(t)
	if match == nil {
		// No match
		redirect(req, resp, t)
		return
	}

	resp.WriteHeader(http.StatusNotImplemented)
}

func proxyGetBlobHandler(resp http.ResponseWriter, req *http.Request) {

}

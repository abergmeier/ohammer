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

func proxyGetManifestHandler(resp http.ResponseWriter, req *http.Request) {

	redirUrl := req.URL

	vars := mux.Vars(req)
	host, ok := vars["originHost"]
	if !ok {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	path, ok := vars["originPath"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	ref, ok := vars["reference"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	i := func() int {
		dockerReg := []byte(fmt.Sprintf("%s/%s:%s", host, path, ref))
		for i, match := range matchers {
			ok := match.reg.Find(dockerReg)
			if ok != nil {
				return i
			}
		}

		return -1
	}()

	if i == -1 {
		redirUrl.Host = host
		redirUrl.Path = fmt.Sprintf("/v2/%s/manifests/%s", path, ref)

		http.Redirect(resp, req, redirUrl.String(), http.StatusTemporaryRedirect)
		return
	}

	resp.WriteHeader(http.StatusNotImplemented)
}

func proxyGetBlobHandler(resp http.ResponseWriter, req *http.Request) {

}

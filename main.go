package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	source = flag.String("source", os.Getenv("SOURCE"), "Source registry to proxy")
)

func init() {
	flag.Parse()
}

func redirectToSource(resp http.ResponseWriter, req *http.Request) {
	redirUrl := req.URL
	redirUrl.Host = *source
	http.Redirect(resp, req, redirUrl.String(), http.StatusTemporaryRedirect)
}

func pullImageManifestHandler(resp http.ResponseWriter, req *http.Request) {
	redirectToSource(resp, req)
	// GET
	// HEAD
}

func pullLayerHandler(resp http.ResponseWriter, req *http.Request) {
	redirectToSource(resp, req)
}

func main() {

	r := mux.NewRouter()

	err := r.HandleFunc("/v2", handleApiVersionCheck).GetError()
	if err != nil {
		panic(err)
	}

	err = r.HandleFunc("/v2/{name}/manifests/{reference}", pullImageManifestHandler).GetError()
	if err != nil {
		panic(err)
	}

	err = r.HandleFunc("/v2/{name}/blobs/{digest}", pullLayerHandler).GetError()
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

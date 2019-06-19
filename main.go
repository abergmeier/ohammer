package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	selector = flag.String("selector", os.Getenv("SELECTOR"), "Selector to extract origin registry")
	repository  = flag.String("repository", os.Getenv("REPOSITORY"), "")
	version = flag.Bool("version", false, "Prints version information")
)

func init() {
	flag.Parse()
	if *selector == "" {
		defaultSelector := ``
		selector = &defaultSelector
	}
}

func redirectToSource(resp http.ResponseWriter, req *http.Request) {
	redirUrl := req.URL

	redirUrl.Host = "foobar.de"
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

	if *version {
		fmt.Printf("O'Hammer Version v%v\n", "0.0.1")
		return
	}

	s, err := NewServer(":8080", *repository)

	if err != nil {
		panic(err)
	}

	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

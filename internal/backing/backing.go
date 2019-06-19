package backing

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/otto-de/ohammer/internal/config"
	"github.com/docker/cli/cli/config/credentials"
)

// Poll checks whether a certain result is available on the redirection target
func Poll(req *http.Request, sectionPath string, t *config.Target) (*http.Response, error) {
	backingURL := req.URL
	backingURL.Host = t.Host
	backingURL.Path = fmt.Sprintf("/v2/%s/%s/%s", t.Path, sectionPath, t.Ref)

	if backingURL.Scheme == "" {
		// Always fallback to https by default
		backingURL.Scheme = "https"
	}

	store := config.DefaultFile().GetCredentialsStore(backingURL.Host)

	switch req.Method {
	case "GET":
		return pollGet(req, store, backingURL.String())
	case "HEAD":
		return pollHead(req, store, backingURL.String())
	}

	panic("NOT IMPLEMENTED YET")
}

func pollGet(req *http.Request, store credentials.Store, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func pollHead(req *http.Request, store credentials.Store, url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

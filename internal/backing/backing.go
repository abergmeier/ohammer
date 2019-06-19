package backing

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/cli/cli/config/credentials"
	"github.com/otto-de/ohammer/internal/config"
	"golang.org/x/oauth2"
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
		return pollGet(req, store, backingURL)
	case "HEAD":
		return pollHead(req, store, backingURL)
	}

	panic("NOT IMPLEMENTED YET")
}

func authenticate(req *http.Request, store credentials.Store, url *url.URL) error {
	hostAuth, err := store.Get(url.Host)
	if err != nil {
		return err
	}
	if hostAuth.RegistryToken == "" {
		return nil
	}
	token := oauth2.Token{
		AccessToken: hostAuth.RegistryToken,
	}
	token.SetAuthHeader(req)
	return nil
}

func pollGet(req *http.Request, store credentials.Store, url *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", url.String(), bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}
	err = authenticate(req, store, url)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func pollHead(req *http.Request, store credentials.Store, url *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url.String(), bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}
	err = authenticate(req, store, url)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

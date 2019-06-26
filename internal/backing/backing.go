package backing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/otto-de/ohammer/internal/config"
)

var (
	keyChain = authn.DefaultKeychain
)

type Backing struct {
	hubTokens map[string]string
}

func NewBacking() *Backing {
	return &Backing{
		hubTokens: make(map[string]string),
	}
}

func (b *Backing) hubAuth(t *config.Target) (string, error) {

	prevToken, ok := b.hubTokens[t.Path]
	if ok {
		return fmt.Sprintf("Bearer %s", prevToken), nil
	}

	if !strings.HasSuffix(t.Host, "docker.io") {
		return "", nil
	}

	authURL := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", t.Path)

	resp, err := http.Get(authURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	token, ok := data["token"]
	if !ok {
		return "", nil
	}

	b.hubTokens[t.Path] = token.(string)

	return fmt.Sprintf("Bearer %s", token), nil
}

// Poll checks whether a certain result is available on the redirection target
func (b *Backing) Poll(req *http.Request, sectionPath string, t *config.Target) (*http.Response, error) {

	backingURL := req.URL
	backingURL.Host = t.Host
	backingURL.Path = fmt.Sprintf("/v2/%s/%s/%s", t.Path, sectionPath, t.Ref)

	if backingURL.Scheme == "" {
		// Always fallback to https by default
		backingURL.Scheme = "https"
	}

	registry, err := name.NewRegistry(name.DefaultRegistry)
	if err != nil {
		return nil, err
	}
	auth, err := keyChain.Resolve(registry)
	if err != nil {
		return nil, err
	}

	switch req.Method {
	case "GET":
		return b.pollGet(req, auth, backingURL, t)
	case "HEAD":
		return b.pollHead(req, auth, backingURL, t)
	}

	panic("NOT IMPLEMENTED YET")
}

func (b *Backing) authenticate(req *http.Request, auth authn.Authenticator, t *config.Target) error {
	authString, err := auth.Authorization()
	if err != nil {
		return err
	}

	if authString == "" {
		authString, err = b.hubAuth(t)
		if err != nil {
			// TODO: Add retry
			return nil
		}
	}

	if authString != "" {
		req.Header.Set("Authorization", authString)
	}
	return nil
}

func newRequest(httpMethod string, url *url.URL) (*http.Request, error) {
	if strings.HasSuffix(url.Host, "docker.io") {
		// Spec by behavior
		url.Host = "registry-1.docker.io"
	}

	return http.NewRequest(httpMethod, url.String(), bytes.NewBuffer([]byte{}))
}

func (b *Backing) pollGet(req *http.Request, auth authn.Authenticator, url *url.URL, t *config.Target) (*http.Response, error) {

	req, err := newRequest("GET", url)
	if err != nil {
		return nil, err
	}
	err = b.authenticate(req, auth, t)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func (b *Backing) pollHead(req *http.Request, auth authn.Authenticator, url *url.URL, t *config.Target) (*http.Response, error) {

	req, err := newRequest("HEAD", url)
	if err != nil {
		return nil, err
	}
	err = b.authenticate(req, auth, t)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

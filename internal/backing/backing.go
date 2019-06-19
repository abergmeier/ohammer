package backing

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/otto-de/ohammer/internal/config"
)

// Poll checks whether a certain result is available on the redirection target
func Poll(req *http.Request, sectionPath string, t *config.Target) (*http.Response, error) {
	if req.Method != "GET" {
		panic("ONLY IMPLEMENT GATE")
	}

	backingURL := req.URL
	backingURL.Host = t.Host
	backingURL.Path = fmt.Sprintf("/v2/%s/%s/%s", t.Path, sectionPath, t.Ref)

	return pollGet(req, backingURL.String())
}

func pollGet(req *http.Request, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

package main

import (
	"net/http"
)

func handleApiVersionCheck(resp http.ResponseWriter, req *http.Request) {
	redirectToSource(resp, req)
}

package config

import (
	"regexp"
)

// Patch represents changes applicable to OCI images
type Patch struct {
	Target
	Reg     *regexp.Regexp
	Content string
}

var (
	// Patches contains all available patches
	Patches = []Patch {
		Patch{
			Target: Target{
				Host: "",
				Path: "golang",
				Ref:  "",
			},
			Reg: regexp.MustCompile(`docker.io/library/golang:.*`),
			Content: `
RUN chgrp -R 0 /go && \
    chmod -R g=u /go

# Chances are very high that cache differs
RUN mkdir -p /.cache && \
    chgrp -R 0 /.cache && \
    chmod -R g=u /.cache

USER 1001:0
`,
		},
	}
)

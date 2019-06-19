package config

// Target represents where a request should be sent to instead
type Target struct {
	Host string
	Path string
	Ref string
}

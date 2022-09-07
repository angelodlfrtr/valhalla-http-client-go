package client

import (
	"crypto/tls"
)

// ClientConfig is the configuration for the client
type ClientConfig struct {
	CustomHeaders map[string]string `json:"custom_headers" yaml:"custom_headers"`
	Endpoint      string            `json:"endpoint" yaml:"endpoint"`
	TLSConfig     *tls.Config
}

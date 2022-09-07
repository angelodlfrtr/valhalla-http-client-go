package client

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// BeforeRequestFn is a function that can be called before sending a request
// allowing to customize request before it is sent
type BeforeRequestFn func(req *fasthttp.Request) error

// Client is the client for the valhalla service
type Client struct {
	config          *ClientConfig
	httpClient      *fasthttp.Client
	beforeRequestFn BeforeRequestFn
}

// NewClient creates a new client with given config cfg
func NewClient(cfg *ClientConfig) *Client {
	clt := &Client{config: cfg}

	httpClient := &fasthttp.Client{
		Name:      "valhalla-http-client-go",
		TLSConfig: cfg.TLSConfig,
	}
	clt.httpClient = httpClient

	return clt
}

// GetFastHTTPClient returns the fasthttp client, allowing custom configuration
func (client *Client) GetFastHTTPClient() *fasthttp.Client {
	return client.httpClient
}

// BeforeRequest allow caller to customize fasthttp request object (ex: adding headers, ...)
func (client *Client) BeforeRequest(fn BeforeRequestFn) {
	client.beforeRequestFn = fn
}

// buildBaseRequest for given method and path
func (client *Client) buildBaseRequest(
	method, path string,
	body interface{},
) (*fasthttp.Request, error) {
	req := fasthttp.AcquireRequest()

	// Set uri
	if err := req.URI().Parse(nil, []byte(client.config.Endpoint+"/"+path)); err != nil {
		fasthttp.ReleaseRequest(req)
		return nil, fmt.Errorf("unable to build request uri: %w", err)
	}

	if client.beforeRequestFn != nil {
		if err := client.beforeRequestFn(req); err != nil {
			fasthttp.ReleaseRequest(req)
			return nil, fmt.Errorf("error while calling BeforeRequest custom fn: %w", err)
		}
	}

	// Set request body
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			fasthttp.ReleaseRequest(req)
			return nil, fmt.Errorf("error while encoding body to json: %w", err)
		}

		req.SetBody(bodyBytes)
	}

	// We send json
	req.Header.SetContentType("application/json")

	return req, nil
}

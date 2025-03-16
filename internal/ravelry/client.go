package ravelry

import "github.com/rkhx/ravelry-exporter/internal/httpclient"

type Client struct {
	Users    *UserClient
	Bundles  *BundleClient
	Patterns *PatternsClient
}

// NewRavelryClient creates `RavelryClient` and initialize
func NewRavelryClient(baseURL, username, password string) *Client {
	httpClient := httpclient.NewCustomHTTPClient(baseURL, username, password)

	return &Client{
		Users:    &UserClient{httpClient},
		Bundles:  &BundleClient{httpClient},
		Patterns: &PatternsClient{httpClient},
	}
}

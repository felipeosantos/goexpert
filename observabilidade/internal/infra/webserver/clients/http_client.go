package clients

import "net/http"

type HTTPClient interface {
	DoFunc(req *http.Request) (*http.Response, error)
}

type CustomClient struct {
	client *http.Client
}

func NewCustomClient(client *http.Client) *CustomClient {
	return &CustomClient{client: client}
}

func (c *CustomClient) DoFunc(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

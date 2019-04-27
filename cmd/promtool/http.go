package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/prometheus/client_golang/api"
)

const defaultTimeout = 2 * time.Minute

type prometheusHTTPClient struct {
	requestTimeout	time.Duration
	httpClient	api.Client
}

func newPrometheusHTTPClient(serverURL string) (*prometheusHTTPClient, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	hc, err := api.NewClient(api.Config{Address: serverURL})
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP client: %s", err)
	}
	return &prometheusHTTPClient{requestTimeout: defaultTimeout, httpClient: hc}, nil
}
func (c *prometheusHTTPClient) do(req *http.Request) (*http.Response, []byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	return c.httpClient.Do(ctx, req)
}
func (c *prometheusHTTPClient) urlJoin(path string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.httpClient.URL(path, nil).String()
}

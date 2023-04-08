package kubelogin

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type Options struct {
	Context   context.Context
	Timeout   time.Duration
	Transport *http.Transport
	Header    http.Header
}

func DefaultOptions() *Options {
	return &Options{
		Context: context.Background(),
		Timeout: time.Duration(20) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Header: http.Header{},
	}
}

func httpRequest(opts *Options, method, uri string, body io.Reader) (*http.Response, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	client := http.Client{
		Timeout:   opts.Timeout,
		Transport: opts.Transport,
	}

	req, err := http.NewRequestWithContext(opts.Context, method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header = opts.Header

	return client.Do(req)
}

func HttpPost(opts *Options, uri string, body io.Reader) (*http.Response, error) {
	return httpRequest(opts, "POST", uri, body)
}

func HttpGet(opts *Options, uri string) (*http.Response, error) {
	return httpRequest(opts, "GET", uri, nil)
}

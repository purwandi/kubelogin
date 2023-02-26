package kubelogin

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

func httpRequest(ctx context.Context, method, uri string, body io.Reader) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Duration(20) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return client.Do(req)
}

func HttpPost(ctx context.Context, uri string, body io.Reader) (*http.Response, error) {
	return httpRequest(ctx, "POST", uri, body)
}

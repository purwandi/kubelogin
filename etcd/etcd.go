package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/purwandi/kubelogin"
)

type EtcdClient struct {
	Endpoint string
	CaCert   string
	TlsCert  string
	TlsKey   string
}

func NewEtcdClient(endpoint, cacert, tlscert, tlskey string) *EtcdClient {
	return &EtcdClient{
		Endpoint: endpoint,
		CaCert:   cacert,
		TlsCert:  tlscert,
		TlsKey:   tlskey,
	}
}

func (c *EtcdClient) GetMetrics() (*http.Response, error) {
	// TLS parse
	cert, err := tls.LoadX509KeyPair(c.TlsCert, c.TlsKey)
	if err != nil {
		return nil, fmt.Errorf("unable to load certificate %s", err.Error())
	}

	caCert, err := os.ReadFile(c.CaCert)
	if err != nil {
		return nil, fmt.Errorf("unable to load ca certificacte %s", err.Error())
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	opts := kubelogin.DefaultOptions()
	opts.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	res, err := kubelogin.HttpGet(opts, fmt.Sprintf("%s/metrics", c.Endpoint))
	if err != nil {
		return nil, err
	}

	return res, nil
}

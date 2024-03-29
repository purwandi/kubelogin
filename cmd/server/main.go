package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/purwandi/kubelogin/etcd"
	"github.com/purwandi/kubelogin/keycloak"
	"github.com/purwandi/kubelogin/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port                    string
	httpsCertificateFile    string
	httpsCertificateKeyFile string
	k8sApiServer            string
	k8sCaCertificate        string
	oidcIssuerUrl           string
	oidcClientId            string
	oidcClientSecret        string
	etcdEndpoint            string
	etcdCacert              string
	etcdTLSCert             string
	etcdTLSKey              string
)

func main() {
	cmd := &cobra.Command{
		Use: "",
		RunE: func(cmd *cobra.Command, args []string) error {

			bytes, err := os.ReadFile(k8sCaCertificate)
			if err != nil {
				return fmt.Errorf("unable to read ca certificate file %s", err.Error())
			}

			s := server.NewServer(server.ServerConfig{
				Port:                port,
				CertificateFile:     httpsCertificateFile,
				CertiticateKeyFile:  httpsCertificateKeyFile,
				KubernetesApiServer: k8sApiServer,
				Keycloak: keycloak.Config{
					OIDCIssuerUrl:    oidcIssuerUrl,
					OIDCCLientID:     oidcClientId,
					OIDCClientSecret: oidcClientSecret,
					KubeAPICaCert:    base64.StdEncoding.EncodeToString(bytes),
				},
				Etcd: etcd.EtcdClient{
					Endpoint: etcdEndpoint,
					CaCert:   etcdCacert,
					TlsCert:  etcdTLSCert,
					TlsKey:   etcdTLSKey,
				},
			})

			s.Run()

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "listening server port")
	cmd.PersistentFlags().StringVar(&httpsCertificateFile, "https-certificate-file", "", "")
	cmd.PersistentFlags().StringVar(&httpsCertificateKeyFile, "https-certificate-key-file", "", "")
	cmd.PersistentFlags().StringVar(&k8sApiServer, "apiserver-host", "", "kubernetes api server host")
	cmd.PersistentFlags().StringVar(&k8sCaCertificate, "apiserver-cacert", "", "kubernetes api server ca certificate")
	cmd.PersistentFlags().StringVar(&oidcIssuerUrl, "oidc-issuer-url", "", "")
	cmd.PersistentFlags().StringVar(&oidcClientId, "oidc-client-id", "", "")
	cmd.PersistentFlags().StringVar(&oidcClientSecret, "oidc-client-secret", "", "")

	cmd.PersistentFlags().StringVar(&etcdEndpoint, "etcd-endpoint", "https://127.0.0.1:2379", "etcd gRPC endpoints")
	cmd.PersistentFlags().StringVar(&etcdCacert, "etcd-cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle")
	cmd.PersistentFlags().StringVar(&etcdTLSCert, "etcd-tlscert", "", "identify secure client using this TLS certificate file")
	cmd.PersistentFlags().StringVar(&etcdTLSKey, "etcd-tlskey", "", "identify secure client using this TLS key file")

	if err := cmd.Execute(); err != nil {
		logrus.Error(err)
	}
}

package main

import (
	"github.com/purwandi/kubelogin/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port                    string
	httpsCertificateFile    string
	httpsCertificateKeyFile string
	k8sApiServer            string
	oidcIssuerUrl           string
	oidcClientId            string
	oidcClientSecret        string
)

func main() {
	cmd := &cobra.Command{
		Use: "",
		RunE: func(cmd *cobra.Command, args []string) error {

			s := server.NewServer(server.ServerConfig{
				Port:                port,
				CertificateFile:     httpsCertificateFile,
				CertiticateKeyFile:  httpsCertificateKeyFile,
				KubernetesApiServer: k8sApiServer,
				Keycloak: server.KeycloakConfig{
					OIDCIssuerUrl:    oidcIssuerUrl,
					OIDCCLientID:     oidcClientId,
					OIDCClientSecret: oidcClientSecret,
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
	cmd.PersistentFlags().StringVar(&oidcIssuerUrl, "oidc-issuer-url", "", "")
	cmd.PersistentFlags().StringVar(&oidcClientId, "oidc-client-id", "", "")
	cmd.PersistentFlags().StringVar(&oidcClientSecret, "oidc-client-secret", "", "")

	if err := cmd.Execute(); err != nil {
		logrus.Error(err)
	}
}

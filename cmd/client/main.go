package main

import (
	"fmt"

	"github.com/purwandi/kubelogin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	server         string
	username       string
	password       string
	token          string
	insecureVerify bool
)

func main() {
	rootCommand := &cobra.Command{
		Use:           "",
		SilenceErrors: true,
		Short: `
Examples:
  # Log in interactively
  kubectl login --username=myuser

  # Log in to the given server with insecure skip tls verify
  kubectl login localhost:6444 --insecure-skip-tls-verify

  # Log in to the given server with the given credentials (will not prompt interactively)
  kubectl login localhost:6444 --username=myuser --password=mypass
		`,
		Run: func(cmd *cobra.Command, args []string) {
			client := &kubelogin.Client{
				Server:   server,
				Username: username,
				Password: password,
			}

			if err := client.Validate(); err != nil {
				fmt.Println(err)
			} else {
				client.Run()
			}

		},
	}

	rootCommand.PersistentFlags().StringVarP(&server, "server", "s", "", "Api server")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "Username for server")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for server")
	rootCommand.PersistentFlags().StringVar(&token, "token", "", "Bearer token for authentication to the API server")
	rootCommand.PersistentFlags().BoolVar(&insecureVerify, "insecure-skip-tls-verify", false, "If true, the server's certificate will not be checked for validity. \nThis will make your HTTPS connections insecure")

	if err := rootCommand.Execute(); err != nil {
		logrus.Error(err)
	}
}

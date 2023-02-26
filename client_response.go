package kubelogin

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
)

type CLientResponse struct {
	Username     string `json:"username"`
	ApiServer    string `json:"apiserver"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (c CLientResponse) GetHostname() string {
	uri, err := url.ParseRequestURI(c.ApiServer)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port())
}

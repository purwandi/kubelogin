package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/purwandi/kubelogin"
)

type KeycloakResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ClientSecret string `json:"client_secret"`
}

type KeycloakRequest struct {
	ResponseType string `json:"response_type"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func (k KeycloakRequest) ToBytes() []byte {
	byt, err := json.Marshal(k)
	if err != nil {
		return nil
	}

	return byt
}

func (k KeycloakRequest) ToFormData() url.Values {
	form := url.Values{}

	form.Add("response_type", k.ResponseType)
	form.Add("grant_type", k.GrantType)
	form.Add("scope", k.Scope)
	form.Add("client_id", k.ClientID)
	form.Add("client_secret", k.ClientSecret)
	form.Add("username", k.Username)
	form.Add("password", k.Password)

	return form
}

func RequestToken(uri string, k KeycloakRequest) (KeycloakResponse, error) {
	var (
		keyRes KeycloakResponse
		err    error
	)

	body := strings.NewReader(k.ToFormData().Encode())
	res, err := kubelogin.HttpPost(context.Background(), uri, body)
	if err != nil {
		fmt.Println(err.Error())
		return keyRes, errors.New("unable to reach out endpoint")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return keyRes, errors.New("invalid credentials")
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return keyRes, errors.New("unable to read message body")
	}

	if err := json.Unmarshal(content, &keyRes); err != nil {
		fmt.Println(err.Error())
		return keyRes, errors.New("unable decode message body")
	}

	return keyRes, err
}

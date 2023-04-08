package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/purwandi/kubelogin/etcd"
	"github.com/purwandi/kubelogin/keycloak"
)

type Handler struct {
	k8sHost  string
	keycloak keycloak.Config
	etcd     etcd.EtcdClient
}

func (h *Handler) Authenticate(c echo.Context) error {
	var (
		err      error
		username = c.FormValue("username")
		password = c.FormValue("password")
	)

	if username == "" {
		err = errors.New("username is required")
	}

	if password == "" {
		err = errors.New("password is required")
	}

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// perform http request
	req := keycloak.KeycloakRequest{
		ResponseType: "id_token",
		GrantType:    "password",
		Scope:        "profile email openid groups",
		ClientID:     h.keycloak.OIDCCLientID,
		ClientSecret: h.keycloak.OIDCClientSecret,
		Username:     username,
		Password:     password,
	}

	res, err := keycloak.KeycloakRequestToken(fmt.Sprintf("%s/protocol/openid-connect/token", h.keycloak.OIDCIssuerUrl), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"username":      username,
		"apiserver":     h.k8sHost,
		"ca_data":       h.keycloak.KubeAPICaCert,
		"id_token":      res.IDToken,
		"refresh_token": res.RefreshToken,
		"access_token":  res.AccessToken,
	})
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, ".")
}

func (h *Handler) EtcdMetrics(c echo.Context) error {
	res, err := h.etcd.GetMetrics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.HTMLBlob(http.StatusOK, content)
}

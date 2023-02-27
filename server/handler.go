package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	k8sHost  string
	keycloak KeycloakConfig
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
	req := KeycloakRequest{
		ResponseType: "id_token",
		GrantType:    "password",
		Scope:        "profile email openid groups",
		ClientID:     h.keycloak.OIDCCLientID,
		ClientSecret: h.keycloak.OIDCClientSecret,
		Username:     username,
		Password:     password,
	}

	res, err := RequestToken(fmt.Sprintf("%s/protocol/openid-connect/token", h.keycloak.OIDCIssuerUrl), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"username":      username,
		"apiserver":     h.k8sHost,
		"id_token":      res.IDToken,
		"refresh_token": res.RefreshToken,
		"access_token":  res.AccessToken,
	})
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, ".")
}

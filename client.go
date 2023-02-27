package kubelogin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"

	"github.com/purwandi/kubelogin/prompt"
)

type Client struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (c *Client) Validate() error {
	if c.Server == "" {
		c.Server = prompt.StringDefault("Server", "https://localhost:8443")
	}

	if c.Username == "" {
		c.Username = prompt.StringRequired("Username")
	}

	if c.Password == "" {
		c.Password = prompt.Password("Password")
	}

	return nil
}

func (c *Client) ToBytes() []byte {
	byt, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return byt
}

func (c *Client) ToForm() url.Values {
	form := url.Values{}
	form.Add("username", c.Username)
	form.Add("password", c.Password)
	return form
}

func (c *Client) Run() {
	var (
		clr  CLientResponse
		body = strings.NewReader(c.ToForm().Encode())
	)

	res, err := HttpPost(context.Background(), c.Server, body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		fmt.Println("invalid credential")
		return
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := json.Unmarshal(content, &clr); err != nil {
		fmt.Println(err.Error())
		return
	}

	// setter kubectl
	exec.Command("kubectl", "config", "set-cluster", clr.GetHostname(), fmt.Sprintf("--server=%s", clr.ApiServer), "--insecure-skip-tls-verify").Output()
	exec.Command("kubectl", "config", "set-credentials", clr.Username, fmt.Sprintf("--token=%s", clr.IDToken)).Output()
	exec.Command("kubectl", "config", "set-context", fmt.Sprintf("%s@%s", clr.Username, clr.GetHostname()), fmt.Sprintf("--cluster=%s", clr.GetHostname()), fmt.Sprintf("--user=%s", clr.Username)).Output()
	exec.Command("kubectl", "config", "use-context", fmt.Sprintf("%s@%s", clr.Username, clr.GetHostname())).Output()

}

package dockerhub

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	dockerTypes "github.com/docker/docker/api/types"

	"github.com/lieinsun/registrytool/registry"
)

func (c *Client) Schema() string {
	return c.url.Scheme
}

func (c *Client) Host() string {
	return c.url.Host
}

func (c *Client) UserName() string {
	return c.username
}

func (c *Client) Password() string {
	return c.password
}

func (c *Client) Token() string {
	return c.token
}

func (c *Client) Login(ctx context.Context) (string, error) {
	data, _ := json.Marshal(dockerTypes.AuthConfig{
		Username: c.username,
		Password: c.password,
	})
	body := bytes.NewBuffer(data)

	u := c.url
	u.Path = LoginURL
	q := url.Values{}
	q.Add("refresh_token", fmt.Sprintf("%v", true))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	tokenResp := tokenResponse{}
	if err = json.Unmarshal(resp, &tokenResp); err != nil {
		return "", err
	}
	c.token = tokenResp.Token
	return tokenResp.Token, nil
}

// CheckConn 查询rateLimit判断客户端是否能正常访问
// https://docs.docker.com/docker-hub/download-rate-limit/
func (c *Client) CheckConn(ctx context.Context) error {
	if c.token == "" {
		return errors.New("unauthorized")
	}
	u := c.url
	u.Host = RateCheckDomain
	u.Path = RateCheckURL
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, withAuth(c.authRateServer))
	if err != nil {
		return err
	}
	return nil
}

// authRateServer 获取检查连接的token
func (c *Client) authRateServer(ctx context.Context) (string, error) {
	u := c.url
	u.Host = RateAuthDomain
	u.Path = RateAuthURL
	q := url.Values{}
	q.Add("service", "registry.docker.io")
	q.Add("scope", "repository:ratelimitpreview/test:pull")
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	c.token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header["Authorization"] = []string{fmt.Sprintf("Basic %s", c.token)}

	resp, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	tokenResp := tokenResponse{}
	if err = json.Unmarshal(resp, &tokenResp); err != nil {
		return "", err
	}
	return tokenResp.Token, nil
}

// ListProjects dockerhub不能获取组织列表
func (c *Client) ListProjects(_ context.Context, _ url.Values) ([]registry.Project, int, error) {
	return nil, 0, nil
}

func (c Client) ProjectClient(account ...string) registry.ProjectCli {
	if len(account) > 0 {
		c.account = account[0]
	}
	return &c
}

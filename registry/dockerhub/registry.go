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

	"github.com/docker/docker/api/types"

	"github.com/lieinsun/registrytool/registry"
)

func (c *Client) Login(ctx context.Context) (string, error) {
	data, _ := json.Marshal(types.AuthConfig{
		Username: c.username,
		Password: c.password,
	})
	body := bytes.NewBuffer(data)

	// 不可以用URL.String() 转义字符会导致404
	u := fmt.Sprintf("%s://%s%s", c.url.Scheme, c.url.Host, LoginURL)
	req, err := http.NewRequestWithContext(ctx, "POST", u, body)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	q.Add("refresh_token", fmt.Sprintf("%v", true))
	c.url.RawQuery = q.Encode()

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
func (c Client) CheckConn(ctx context.Context) error {
	if c.token == "" {
		return errors.New("unauthorized")
	}
	c.url.Host = RateCheckDomain
	c.url.Path = RateCheckURL
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.url.String(), nil)
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
func (c Client) authRateServer(ctx context.Context) (string, error) {
	// 不可以用URL.String() 转义字符会导致404
	u := fmt.Sprintf("%s://%s%s", c.url.Scheme, RateAuthDomain, RateAuthURL)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
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
func (c Client) ListProjects(_ context.Context, _ url.Values) ([]registry.Project, int, error) {
	return nil, 0, nil
}

func (c Client) ProjectClient(account ...string) registry.ProjectCli {
	if len(account) > 0 {
		c.account = account[0]
	}
	return c
}

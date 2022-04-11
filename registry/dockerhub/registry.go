package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/lie-inthesun/remotescan/registry"
)

func (c Client) Login(ctx context.Context) (string, error) {
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

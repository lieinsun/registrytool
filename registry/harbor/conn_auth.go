package harbor

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

// Login harbor使用Basic token
func (c *Client) Login(ctx context.Context) (string, error) {
	c.token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	// 使用ping 检查token是否有效
	err := c.Ping(ctx)
	if err != nil {
		c.token = ""
		return "", err
	}
	return c.token, nil
}

func (c *Client) Ping(ctx context.Context) error {
	u := c.url
	u.Path = CurrentUserURL
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

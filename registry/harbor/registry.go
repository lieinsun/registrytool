package harbor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang/glog"
	"github.com/lie-inthesun/registrytool/registry"
)

// Login harbor使用Basic token
func (c *Client) Login(ctx context.Context) (string, error) {
	c.token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	// 使用ping 检查token是否有效
	err := c.ping(ctx)
	if err != nil {
		c.token = ""
		return "", err
	}
	return c.token, nil
}

// Ping 使用查询当前登录用户的方法验证登录
func (c Client) ping(ctx context.Context) error {
	u := c.url
	u.Path = CurrentUserURL
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CheckConn(ctx context.Context) error {
	if c.token == "" {
		return errors.New("unauthorized")
	}
	return c.ping(ctx)
}

func (c Client) ProjectClient(project ...string) registry.ProjectCli {
	if len(project) > 0 {
		c.project = project[0]
	}
	return c
}

func (c Client) ListProjects(ctx context.Context, params url.Values) ([]registry.Project, int, error) {
	c.url.Path = ListProjectsURL
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, header, err := c.doRequest(req)
	if err != nil {
		return nil, 0, err
	}
	var projectsResp projectsResponse
	if err = json.Unmarshal(resp, &projectsResp); err != nil {
		return nil, 0, err
	}

	list := make([]registry.Project, 0, len(projectsResp))
	for _, p := range projectsResp {
		project := registry.Project{
			Name:        p.Name,
			Metadata:    p.Metadata,
			OwnerName:   p.OwnerName,
			RepoCount:   p.RepoCount,
			CreatedTime: p.CreationTime.Unix(),
			UpdatedTime: p.UpdateTime.Unix(),
		}
		list = append(list, project)
	}

	total, err := strconv.Atoi(header.Get("X-Total-Count"))
	if err != nil {
		glog.Error(err)
	}
	return list, total, nil
}

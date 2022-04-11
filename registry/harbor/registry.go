package harbor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lie-inthesun/remotescan/registry"
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
func (c *Client) ping(ctx context.Context) error {
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

func (c Client) AccountOrProject(project string) registry.ProjectCli {
	c.project = project
	return c
}

func (c Client) ListProjects(ctx context.Context, params url.Values) ([]registry.Project, int, error) {
	c.url.Path = ListProjectsURL
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, 0, err
	}
	var listProjectsResp projectsResp
	if err = json.Unmarshal(resp, &listProjectsResp); err != nil {
		return nil, 0, err
	}

	list := make([]registry.Project, 0, len(listProjectsResp))
	for _, p := range listProjectsResp {
		project := registry.Project{
			Name:         p.Name,
			Metadata:     p.Metadata,
			OwnerName:    p.OwnerName,
			RepoCount:    p.RepoCount,
			CreationTime: p.CreationTime.Unix(),
			UpdatedTime:  p.UpdateTime.Unix(),
		}
		list = append(list, project)
	}

	// TODO /api/v2.0/statistics 查询project总数
	return list, 0, nil
}

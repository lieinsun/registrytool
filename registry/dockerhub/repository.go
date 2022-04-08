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

func (c Client) AccountOrProject(account string) registry.ProjectCli {
	c.account = account
	return c
}

func (c Client) Image(image string) registry.ImageCli {
	c.image = image
	return c
}

// GetRepositories 查询镜像tag列表
//func (c *Client) GetRepositories(ctx context.Context, username string, page, pageSize int, ordering string) ([]hubReposResult, int, error) {
//	if username == "" {
//		username = c.username
//	}
//
//	c.url.Path = fmt.Sprintf(RepositoriesURL, username)
//	q := url.Values{}
//	q.Add("page", fmt.Sprintf("%v", page))
//	q.Add("page_size", fmt.Sprintf("%v", pageSize))
//	q.Add("ordering", ordering)
//	c.url.RawQuery = q.Encode()
//	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	resp, err := c.doRequest(req)
//	if err != nil {
//		return nil, 0, err
//	}
//	var reposResp hubReposResponse
//	if err = json.Unmarshal(resp, &reposResp); err != nil {
//		return nil, 0, err
//	}
//
//	return reposResp.Results, reposResp.Count, nil
//}
//
//// GetTags 查询镜像tag列表
//func (c *Client) GetTags(ctx context.Context, repository string, page, pageSize int, ordering string) ([]hubTagResult, int, error) {
//	repoPath, err := referencePath(repository)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	c.url.Path = fmt.Sprintf(TagsURL, repoPath)
//	q := url.Values{}
//	q.Add("page", fmt.Sprintf("%v", page))
//	q.Add("page_size", fmt.Sprintf("%v", pageSize))
//	q.Add("ordering", ordering)
//	c.url.RawQuery = q.Encode()
//	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	resp, err := c.doRequest(req)
//	if err != nil {
//		return nil, 0, err
//	}
//	var tagsResp hubTagsResponse
//	if err = json.Unmarshal(resp, &tagsResp); err != nil {
//		return nil, 0, err
//	}
//	return tagsResp.Results, tagsResp.Count, nil
//}

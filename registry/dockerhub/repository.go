package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"remotescan-pkg/registry"
)

// GetImageDetail 查询指定tag的镜像详情
func (c *Client) GetImageDetail(ctx context.Context, account, image, tag string) (registry.Image, error) {
	if account == "" {
		account = c.username
	}
	if tag == "" {
		tag = "latest"
	}
	repoPath, err := referencePath(account, image)
	if err != nil {
		return nil, err
	}

	u := c.url
	u.Path = fmt.Sprintf(ImageDetailURL, repoPath, tag)
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var detailResp imageDetailResp
	if err = json.Unmarshal(resp, &detailResp); err != nil {
		return nil, err
	}

	i := Image{
		Auth: registry.Auth{
			UserName: c.username,
			Password: c.password,
			Token:    c.token,
		},
		Account:     account,
		Name:        image,
		Tag:         tag,
		Size:        detailResp.FullSize,
		LastUpdated: detailResp.LastUpdated.Unix(),
	}
	if len(detailResp.Images) > 0 {
		img := detailResp.Images[0]
		i.Digest = img.Digest
		i.Size = img.Size
		i.Os = img.Os
	}
	return &i, nil
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

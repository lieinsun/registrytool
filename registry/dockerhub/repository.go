package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aquasecurity/fanal/types"
	"github.com/lie-inthesun/remotescan/registry"
)

type Image struct {
	registry.Auth
	Account     string // dockerhub用户
	Name        string
	Tag         string
	Digest      string
	Size        int
	Os          string
	LastUpdated int64
}

func (c Client) ImageDetail(ctx context.Context, tag string) (registry.Image, error) {
	c.tag = tag
	if c.account == "" {
		c.account = c.username
	}
	if tag == "" {
		tag = "latest"
	}
	repoPath, err := referencePath(c.account, c.image)
	if err != nil {
		return nil, err
	}

	c.url.Path = fmt.Sprintf(ImageDetailURL, repoPath, tag)
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
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
		Account:     c.account,
		Name:        c.image,
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

func (i *Image) TrivyReference() (string, types.DockerOption) {
	account := i.Account
	if i.Account == "" {
		account = i.UserName
	}
	if i.Tag == "" {
		i.Tag = "latest"
	}
	ref := fmt.Sprintf("%s/%s:%s", account, i.Name, i.Tag)
	if i.Digest != "" {
		ref = ref + "@" + i.Digest
	}

	dockerOption := types.DockerOption{
		UserName:      i.UserName,
		Password:      i.Password,
		RegistryToken: i.Token,
	}
	return ref, dockerOption
}

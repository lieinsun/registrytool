package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"remotescan-pkg/registry"
)

// GetImageDetail 查询指定tag的镜像详情
// tagOrDigest 镜像tag或者sha256 digest
func (c *Client) GetImageDetail(ctx context.Context, project, image, tagOrDigest string) (registry.Image, error) {
	u := c.url
	u.Path = fmt.Sprintf(ImageDetailURL, project, image, tagOrDigest)
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
			Url:      c.url.Host,
			UserName: c.username,
			Password: c.password,
			Token:    c.token,
		},
		Project:  project,
		Name:     image,
		Digest:   detailResp.Digest,
		Size:     detailResp.Size,
		Os:       detailResp.ExtraAttrs.Os,
		PushTime: detailResp.PushTime.Unix(),
	}
	if len(detailResp.Tags) > 0 {
		tag := detailResp.Tags[0]
		i.Tag = tag.Name
	}
	return &i, nil
}

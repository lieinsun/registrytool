package harbor

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
	Project  string
	Name     string
	Tag      string
	Digest   string
	Size     int
	Os       string
	PushTime int64
}

// Detail 查询指定tag的镜像详情
// tagOrDigest 镜像tag或者sha256 digest
func (c Client) Detail(ctx context.Context, tagOrDigest string) (registry.Image, error) {
	c.url.Path = fmt.Sprintf(ImageDetailURL, c.project, c.image, tagOrDigest)
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
			Url:      c.url.Host,
			UserName: c.username,
			Password: c.password,
			Token:    c.token,
		},
		Project:  c.project,
		Name:     c.image,
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

func (i *Image) TrivyReference() (string, types.DockerOption) {
	ref := fmt.Sprintf("%s/%s/%s", i.Url, i.Project, i.Name)
	if i.Tag != "" {
		ref = ref + ":" + i.Tag
	}
	if i.Digest != "" {
		ref = ref + "@" + i.Digest
	}

	dockerOption := types.DockerOption{
		UserName: i.UserName,
		Password: i.Password,
	}
	return ref, dockerOption
}

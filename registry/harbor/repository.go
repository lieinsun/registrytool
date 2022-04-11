package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aquasecurity/fanal/types"
	"github.com/lie-inthesun/remotescan/registry"
	"github.com/lie-inthesun/remotescan/scanner"
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

func (c Client) ListArtifacts(ctx context.Context, params url.Values) ([]registry.Artifact, int, error) {
	c.url.Path = fmt.Sprintf(ListArtifactsURL, c.project, c.repository)
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, 0, err
	}
	var artifactsResp artifactsResponse
	if err = json.Unmarshal(resp, &artifactsResp); err != nil {
		return nil, 0, err
	}

	list := make([]registry.Artifact, 0, len(artifactsResp))
	for _, a := range artifactsResp {
		artifact := registry.Artifact{
			Digest:      a.Digest,
			Os:          a.ExtraAttrs.Os,
			Size:        a.Size,
			UpdatedTime: a.PushTime.Unix(),
		}

		list = append(list, artifact)
	}

	// TODO /api/v2.0/projects/image/repositories/ccs-build 查询artifacts总数
	return list, 0, nil
}

func (c Client) ListTags(ctx context.Context, params url.Values, reference ...string) ([]registry.Tag, int, error) {
	//if len(reference) == 0 {
	//	reference = []string{"latest"}
	//}
	c.url.Path = fmt.Sprintf(ListTagsURL, c.project, c.repository, reference[0])
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, 0, err
	}
	var tagsResp tagsResponse
	if err = json.Unmarshal(resp, &tagsResp); err != nil {
		return nil, 0, err
	}

	list := make([]registry.Tag, 0, len(tagsResp))
	for _, t := range tagsResp {
		tag := registry.Tag{
			Name:        t.Name,
			UpdatedTime: t.PushTime.Unix(),
		}

		list = append(list, tag)
	}

	// TODO tags总数 需要从response Header X-Total-Count获取
	return list, len(list), nil
}

// ImageDetail
// tagOrDigest 镜像tag或者sha256 digest
func (c Client) ImageDetail(ctx context.Context, tagOrDigest string) (registry.Image, error) {
	c.url.Path = fmt.Sprintf(ImageDetailURL, c.project, c.repository, tagOrDigest)
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
		Name:     c.repository,
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

func (i *Image) TrivyReference() *scanner.ScanReference {
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
	return &scanner.ScanReference{
		ImageName:    ref,
		DockerOption: dockerOption,
	}
}

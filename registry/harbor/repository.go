package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aquasecurity/fanal/types"
	"github.com/lieinsun/registrytool/registry"
	"github.com/lieinsun/registrytool/scanner"
)

func (c Client) ListArtifacts(ctx context.Context, params url.Values) ([]registry.Artifact, int, error) {
	c.url.Path = fmt.Sprintf(ListArtifactsURL, c.project, c.repository)
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, header, err := c.doRequest(req)
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

	total, _ := strconv.Atoi(header.Get("X-Total-Count"))
	return list, total, nil
}

func (c Client) ListTags(ctx context.Context, params url.Values, reference ...string) ([]registry.Tag, int, error) {
	if len(reference) == 0 {
		reference = []string{"latest"}
	}
	c.url.Path = fmt.Sprintf(ListTagsURL, c.project, c.repository, reference[0])
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, header, err := c.doRequest(req)
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

	total, _ := strconv.Atoi(header.Get("X-Total-Count"))
	return list, total, nil
}

// ImageDetail
// tagOrDigest 镜像tag或者sha256 digest
func (c Client) ImageDetail(ctx context.Context, tagOrDigest string) (*registry.Image, error) {
	c.url.Path = fmt.Sprintf(ImageDetailURL, c.project, c.repository, tagOrDigest)
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var detailResp imageDetailResp
	if err = json.Unmarshal(resp, &detailResp); err != nil {
		return nil, err
	}

	i := registry.Image{
		Namespace:   c.project,
		Name:        c.repository,
		Digest:      detailResp.Digest,
		Size:        detailResp.Size,
		Os:          detailResp.ExtraAttrs.Os,
		UpdatedTime: detailResp.PushTime.Unix(),
	}
	if len(detailResp.Tags) > 0 {
		tag := detailResp.Tags[0]
		i.Tag = tag.Name
	}
	return &i, nil
}

func (c Client) ScanReference(tag, digest string) *scanner.ScanReference {
	c.tag = tag
	if c.tag == "" {
		c.tag = "latest"
	}
	ref := fmt.Sprintf("%s/%s/%s:%s", c.url.Host, c.project, c.repository, c.tag)
	if digest != "" {
		ref = ref + "@" + digest
	}

	dockerOption := types.DockerOption{
		UserName:      c.username,
		Password:      c.password,
		RegistryToken: c.token,
	}
	return &scanner.ScanReference{
		ImageName:    ref,
		DockerOption: dockerOption,
	}
}

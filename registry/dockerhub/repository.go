package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aquasecurity/fanal/types"
	"github.com/lie-inthesun/registrytool/registry"
	"github.com/lie-inthesun/registrytool/scanner"
)

func (c Client) ListArtifacts(_ context.Context, _ url.Values) ([]registry.Artifact, int, error) {
	return nil, 0, nil
}

func (c Client) ListTags(ctx context.Context, params url.Values, _ ...string) ([]registry.Tag, int, error) {
	if c.account == "" {
		c.account = c.username
	}
	c.url.Path = fmt.Sprintf(ListTagsURL, c.account, c.repository)
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

	var list []registry.Tag
	if tagsCount := len(tagsResp.Results); tagsCount > 0 {
		list = make([]registry.Tag, 0, tagsCount)
		for _, result := range tagsResp.Results {
			tag := registry.Tag{
				Name:        result.Name,
				Size:        result.FullSize,
				UpdatedTime: result.LastUpdated.Unix(),
			}
			if len(result.Images) > 0 {
				tag.Digest = result.Images[0].Digest
			}
			list = append(list, tag)
		}
	}

	return list, tagsResp.Count, nil
}

func (c Client) ImageDetail(ctx context.Context, tag string) (*registry.Image, error) {
	c.tag = tag
	if c.account == "" {
		c.account = c.username
	}
	if tag == "" {
		tag = "latest"
	}
	repoPath, err := referencePath(c.account, c.repository)
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
	var detailResp imageDetailResponse
	if err = json.Unmarshal(resp, &detailResp); err != nil {
		return nil, err
	}

	i := registry.Image{
		Namespace:   c.account,
		Name:        c.repository,
		Tag:         tag,
		Size:        detailResp.FullSize,
		UpdatedTime: detailResp.LastUpdated.Unix(),
	}
	if len(detailResp.Images) > 0 {
		img := detailResp.Images[0]
		i.Digest = img.Digest
		i.Size = img.Size
		i.Os = img.Os
	}
	return &i, nil
}

func (c Client) ScanReference(tag, digest string) *scanner.ScanReference {
	if c.account == "" {
		c.account = c.username
	}
	c.tag = tag
	if c.tag == "" {
		c.tag = "latest"
	}
	ref := fmt.Sprintf("%s/%s:%s", c.account, c.repository, c.tag)
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

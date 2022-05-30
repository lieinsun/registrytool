package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aquasecurity/fanal/types"

	"github.com/lieinsun/registrytool/registry"
	"github.com/lieinsun/registrytool/scanner"
)

func (c *Client) Repository() string {
	return c.query.repository
}

// ListArtifacts tags按照digest分组
func (c *Client) ListArtifacts(ctx context.Context, params url.Values) ([]registry.Artifact, int, error) {
	u := c.url
	u.Path = fmt.Sprintf(ListTagsURL, c.account, url.PathEscape(c.repository))
	u.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req, withAuth(c.Login))
	if err != nil {
		return nil, 0, err
	}
	var tagsResp tagsResponse
	if err = json.Unmarshal(resp, &tagsResp); err != nil {
		return nil, 0, err
	}

	// map分组tagResult
	temp := make(map[string][]tagResult)
	for _, tag := range tagsResp.Results {
		if len(tag.Images) == 0 {
			continue
		}
		digest := tag.Images[0].Digest
		if arr, ok := temp[digest]; ok {
			arr = append(arr, tag)
			temp[digest] = arr
		} else {
			temp[digest] = []tagResult{tag}
		}
	}

	artifacts := make([]registry.Artifact, 0, len(temp))
	for _, tagArr := range temp {
		// 组装artifact
		var artifact registry.Artifact
		for i, tagRes := range tagArr {
			tagImg := tagRes.Images[0]
			if i == 0 {
				artifact = registry.Artifact{
					Digest:      tagImg.Digest,
					Os:          tagImg.Os,
					Size:        tagImg.Size,
					UpdatedTime: tagRes.LastUpdated.Unix(),
					Tags:        make([]registry.Tag, 0, len(tagArr)),
				}
			}
			tag := registry.Tag{
				Name:        tagRes.Name,
				Digest:      tagImg.Digest,
				Size:        tagImg.Size,
				UpdatedTime: tagImg.LastPushed.Unix(),
			}
			artifact.Tags = append(artifact.Tags, tag)
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, len(artifacts), nil
}

func (c Client) ListTags(ctx context.Context, reference string, params url.Values) ([]registry.Tag, int, error) {
	u := c.url
	u.Path = fmt.Sprintf(ListTagsURL, c.account, c.repository)
	u.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req, withAuth(c.Login))
	if err != nil {
		return nil, 0, err
	}
	var tagsResp tagsResponse
	if err = json.Unmarshal(resp, &tagsResp); err != nil {
		return nil, 0, err
	}

	// 过滤不匹配的tag数据
	if reference != "" && !strings.HasPrefix(reference, "sha256:") {
		for i := range tagsResp.Results {
			if tagsResp.Results[i].Name == reference {
				reference = tagsResp.Results[i].Images[0].Digest
				break
			}
		}
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

			if reference != "" && tag.Digest != reference {
				tagsResp.Count--
				continue
			}
			list = append(list, tag)
		}
	}

	return list, tagsResp.Count, nil
}

func (c *Client) ImageDetail(ctx context.Context, tag string) (*registry.Image, error) {
	if tag == "" {
		tag = "latest"
	}
	c.tag = tag
	repoPath, err := referencePath(c.account, c.repository)
	if err != nil {
		return nil, err
	}

	u := c.url
	u.Path = fmt.Sprintf(ImageDetailURL, repoPath, tag)
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req, withAuth(c.Login))
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

func (c *Client) Reference(tag, digest string) *scanner.RemoteReference {
	if tag == "" && digest == "" {
		tag = "latest"
	}
	c.tag, c.digest = tag, digest
	repoTag := fmt.Sprintf("%s/%s", c.account, c.repository)

	if c.tag != "" {
		repoTag = repoTag + ":" + c.tag
	}
	if digest != "" {
		repoTag = repoTag + "@" + digest
	}

	dockerOption := types.DockerOption{
		UserName:      c.username,
		Password:      c.password,
		RegistryToken: c.token,
	}
	return scanner.NewRemoteReference(repoTag, dockerOption)
}

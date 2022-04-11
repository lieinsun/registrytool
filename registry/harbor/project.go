package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lie-inthesun/remotescan/registry"
)

func (c Client) ListRepositories(ctx context.Context, params url.Values) ([]registry.Repository, int, error) {
	c.url.Path = fmt.Sprintf(ListRepositoriesURL, c.project)
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, 0, err
	}
	var repositoriesResp repositoriesResponse
	if err = json.Unmarshal(resp, &repositoriesResp); err != nil {
		return nil, 0, err
	}

	list := make([]registry.Repository, 0, len(repositoriesResp))
	for _, r := range repositoriesResp {
		repository := registry.Repository{
			Name:        r.Name,
			UpdatedTime: r.UpdateTime.Unix(),
		}
		if index := strings.Index(r.Name, "/"); index > 0 {
			repository.Namespace = r.Name[:index]
			repository.Name = r.Name[index+1:]
		}
		list = append(list, repository)
	}

	// TODO /api/v2.0/projects/34/summary 查询repositories总数
	return list, 0, nil
}

func (c Client) Repository(repository string) registry.RepositoryCli {
	c.repository = repository
	return c
}

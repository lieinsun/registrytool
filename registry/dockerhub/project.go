package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lie-inthesun/remotescan/registry"
)

func (c Client) ListRepositories(ctx context.Context, params url.Values) ([]registry.Repository, int, error) {
	if c.account == "" {
		c.account = c.username
	}
	c.url.Path = fmt.Sprintf(ListRepositoriesURL, c.account)
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, 0, err
	}
	var reposResp repositoriesResponse
	if err = json.Unmarshal(resp, &reposResp); err != nil {
		return nil, 0, err
	}

	var list []registry.Repository
	if reposCount := len(reposResp.Results); reposCount > 0 {
		list = make([]registry.Repository, 0, reposCount)
		for _, result := range reposResp.Results {
			repository := registry.Repository{
				Name:        result.Name,
				Namespace:   result.Namespace,
				UpdatedTime: result.LastUpdated.Unix(),
			}
			list = append(list, repository)
		}
	}

	return list, reposResp.Count, nil
}

func (c Client) Repository(repository string) registry.RepositoryCli {
	c.repository = repository
	return c
}

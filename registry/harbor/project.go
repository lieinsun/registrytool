package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lieinsun/registrytool/registry"
)

func (c Client) ListRepositories(ctx context.Context, params url.Values) ([]registry.Repository, int, error) {
	c.url.Path = fmt.Sprintf(ListRepositoriesURL, c.project)
	c.url.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", c.url.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, header, err := c.doRequest(req)
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
			CreatedTime: r.CreationTime.Unix(),
			UpdatedTime: r.UpdateTime.Unix(),
		}
		if index := strings.Index(r.Name, "/"); index > 0 {
			repository.Namespace = r.Name[:index]
			repository.Name = r.Name[index+1:]
		}
		list = append(list, repository)
	}

	total, _ := strconv.Atoi(header.Get("X-Total-Count"))
	return list, total, nil
}

func (c Client) RepositoryClient(repository string) registry.RepositoryCli {
	c.repository = repository
	return &c
}

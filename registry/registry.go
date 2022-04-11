package registry

import (
	"context"
	"net/url"

	"github.com/aquasecurity/fanal/types"
)

type Registry interface {
	// Login 检验账号密码
	Login(ctx context.Context) (string, error)
	// AccountOrProject 指定仓库account或project
	AccountOrProject(accountOrProject string) ProjectCli
	// ListProjects harbor查询项目列表
	ListProjects(ctx context.Context, params url.Values) ([]Project, int, error)
}

type ProjectCli interface {
	// Image 指定镜像名
	Image(image string) ImageCli
	ListRepositories(ctx context.Context, params url.Values) ([]Repository, int, error)
	//Get() (Project, error)
}

type ImageCli interface {

	ImageDetail(ctx context.Context, tag string) (Image, error)
}

type Image interface {
	TrivyReference() (string, types.DockerOption)
}

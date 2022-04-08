package registry

import (
	"context"

	"github.com/aquasecurity/fanal/types"
)

type Auth struct {
	Url      string
	UserName string
	Password string
	Token    string
}

type Registry interface {
	// Login 检验账号密码
	Login(ctx context.Context) (string, error)
	// GetImageDetail 查询指定tag的镜像详情
	//GetImageDetail(ctx context.Context, accountOrProject, image, tag string) (Image, error)

	// AccountOrProject 指定仓库account或project
	AccountOrProject(accountOrProject string) ProjectCli
}

type ProjectCli interface {
	Image(image string) ImageCli
	//Get() (Project, error)
	//List() ([]Project, int, error)
}

type Project interface {
}

type ImageCli interface {
	//Get(tag string) (Image, error)
	Detail(ctx context.Context, tag string) (Image, error)
}

type Image interface {
	TrivyReference() (string, types.DockerOption)
}

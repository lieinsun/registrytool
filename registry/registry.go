package registry

import (
	"context"

	"github.com/aquasecurity/fanal/types"
)

type Registry interface {
	// Login 检验账号密码
	Login(ctx context.Context) (string, error)
	// GetImageDetail 查询指定tag的镜像详情
	GetImageDetail(ctx context.Context, accountOrProject, image, tag string) (Image, error)
}

type Image interface {
	TrivyReference() (string, types.DockerOption)
}

type Auth struct {
	Url      string
	UserName string
	Password string
	Token    string
}

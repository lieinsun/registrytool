package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/lieinsun/registrytool/scanner"
)

type Registry interface {
	Schema() string
	Host() string
	UserName() string
	Password() string
	Token() string

	// Login 检验账号密码
	Login(ctx context.Context) (string, error)
	// CheckConn 检查客户端连接(仓库是否可访问，登录是否有效) 不可以使用login频繁检测
	CheckConn(ctx context.Context) error
	// ProjectClient 指定仓库account或project(值接收器 实现client并发复用)
	ProjectClient(accountOrProject ...string) ProjectCli
	// ListProjects harbor查询项目列表
	ListProjects(ctx context.Context, params url.Values) ([]Project, int, error)
}

type ProjectCli interface {
	Project() string

	// RepositoryClient 指定镜像名(值接收器 实现client并发复用)
	RepositoryClient(repo string) RepositoryCli
	// ListRepositories 查询project下面的镜像repo列表
	ListRepositories(ctx context.Context, params url.Values) ([]Repository, int, error)

	Registry
}

type RepositoryCli interface {
	Repository() string

	// ListArtifacts harbor查询repo下面的tag列表
	ListArtifacts(ctx context.Context, params url.Values) ([]Artifact, int, error)
	// ListTags 查询repo下面的tag列表
	// reference(tag或digest) dockerhub查询不需要指定
	ListTags(ctx context.Context, params url.Values, reference ...string) ([]Tag, int, error)
	// ImageDetail 指定tag查询镜像详情
	ImageDetail(ctx context.Context, tag string) (*Image, error)
	// Reference 镜像全称 用于拉取/扫描
	Reference(tag, digest string) *scanner.Reference

	ProjectCli
}

func EncodeAuthHeader(username string, password string) string {
	src := fmt.Sprintf("{ \"username\": \"%s\", \"password\": \"%s\" }", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(src))

	return encoded
}

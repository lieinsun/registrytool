package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/lieinsun/registrytool/scanner"
)

type Registry interface {
	// Login 检验账号密码
	Login(ctx context.Context) (string, error)
	// CheckConn 检查客户端连接 不可以使用login频繁检测
	CheckConn(ctx context.Context) error
	// ProjectClient 指定仓库account或project
	ProjectClient(accountOrProject ...string) ProjectCli
	// ListProjects harbor查询项目列表
	ListProjects(ctx context.Context, params url.Values) ([]Project, int, error)
}

type ProjectCli interface {
	// RepositoryClient 指定镜像名
	RepositoryClient(repo string) RepositoryCli
	// ListRepositories 查询project下面的镜像repo列表
	ListRepositories(ctx context.Context, params url.Values) ([]Repository, int, error)
}

type RepositoryCli interface {
	// ListArtifacts harbor查询repo下面的tag列表
	ListArtifacts(ctx context.Context, params url.Values) ([]Artifact, int, error)
	// ListTags 查询repo下面的tag列表
	// reference(tag或digest) dockerhub查询不需要指定
	ListTags(ctx context.Context, params url.Values, reference ...string) ([]Tag, int, error)
	// ImageDetail 指定tag查询镜像详情
	ImageDetail(ctx context.Context, tag string) (*Image, error)
	// Reference 镜像全称 用于拉取/扫描
	Reference(tag, digest string) *scanner.Reference

	ImagePull(ctx context.Context, tag, digest string) error
}

func EncodeAuthHeader(username string, password string) string {
	src := fmt.Sprintf("{ \"username\": \"%s\", \"password\": \"%s\" }", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(src))

	return encoded
}

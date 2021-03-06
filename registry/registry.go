package registry

import (
	"context"
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
	// reference(tagName或digest)
	// 	harbor: 必须指定tagName或digest
	// 	dockerhub: api无法指定查询 只能返回repo下所有tag 通过筛选结果集处理
	ListTags(ctx context.Context, reference string, params url.Values) ([]Tag, int, error)
	// ImageDetail 查询镜像详情
	// reference(tagName或digest)
	// 	harbor: 使用tagName或digest
	// 	dockerhub: 只能使用tagName
	ImageDetail(ctx context.Context, reference string) (*Image, error)
	// Reference 镜像全称 用于拉取/扫描
	Reference(tag, digest string) *scanner.RemoteReference

	ProjectCli
}

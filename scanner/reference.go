package scanner

import (
	"context"
	"fmt"
	"os"

	"github.com/aquasecurity/fanal/image"
	"github.com/aquasecurity/fanal/types"
)

type Reference interface {
	NewFanalImage(ctx context.Context) (types.Image, func(), error)
	ImageName() string
}

type remoteReference struct {
	// repoTag 扫描所需要的镜像全称 digest为可选字段
	// dockerhub	{username}/{repositoryName}:{tagName}@{digest}
	//				datadog/agent:latest@sha256:775b3a472c0581c4bceaeafabab78c104cc8a5ce806e6ce0634aa4fcf11e41ab
	// 自建仓库		{host}/{projectName}/{repositoryName}:{tagName}@{digest}
	//				127.0.0.1/myProject/myRepo:myTag:@sha256:xxx
	repoTag      string
	dockerOption types.DockerOption
}

func NewRemoteReference(repoTag string, dockerOption types.DockerOption) *remoteReference {
	return &remoteReference{
		repoTag:      repoTag,
		dockerOption: dockerOption,
	}
}

func (r *remoteReference) NewFanalImage(ctx context.Context) (types.Image, func(), error) {
	// TODO 考虑替换NewDockerImage 只处理remote逻辑
	img, cleanup, err := image.NewDockerImage(ctx, r.repoTag, r.dockerOption)
	if err != nil {
		return nil, nil, fmt.Errorf("NewDockerImage err: %w", err)
	}

	return img, cleanup, nil
}

func (r *remoteReference) ImageName() string {
	return r.repoTag
}

type archiveReference struct {
	// fileName 镜像压缩包/导出文件
	// 本地压缩包		{path}/{fileName}
	//				/tmp/image/775b3a472c0581c4bceaeafabab78c104cc8a5ce806e6ce0634aa4fcf11e41ab.tar
	fileName string
}

func NewArchiveReference(fileName string) *archiveReference {
	return &archiveReference{fileName: fileName}
}

func (r *archiveReference) NewFanalImage(_ context.Context) (types.Image, func(), error) {
	// TODO 考虑替换NewArchiveImage fanal未处理file.Close()
	img, err := image.NewArchiveImage(r.fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("NewArchiveImage err: %w", err)
	}

	return img, func() {
		_ = os.Remove(r.fileName)
	}, nil
}

func (r *archiveReference) ImageName() string {
	return r.fileName
}

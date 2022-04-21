package scanner

import (
	"context"
	"fmt"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	artifactImg "github.com/aquasecurity/fanal/artifact/image"
	"github.com/aquasecurity/fanal/image"
	dockerTypes "github.com/aquasecurity/fanal/types"
	trivyScanner "github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
)

type Reference struct {
	// ImageName 扫描所需要的镜像全称 digest为可选字段
	// dockerhub	{username}/{repositoryName}:{tagName}@{digest}
	//				datadog/agent:latest@sha256:775b3a472c0581c4bceaeafabab78c104cc8a5ce806e6ce0634aa4fcf11e41ab
	// 自建仓库		{host}/{projectName}/{repositoryName}:{tagName}@{digest}
	//				127.0.0.1/myProject/myRepo:myTag:@sha256:xxx
	ImageName    string
	DockerOption dockerTypes.DockerOption
}

func (t *Trivy) newScanner(ctx context.Context, refer *Reference) (*trivyScanner.Scanner, func(), error) {
	// TODO 考虑替换NewDockerImage 只处理remote逻辑
	dockerImg, cleanup, err := image.NewDockerImage(ctx, refer.ImageName, refer.DockerOption)
	if err != nil {
		return nil, nil, fmt.Errorf("creating dockerImage: %w", err)
	}

	// NoProgress 不显示git相关信息
	artifactOpt := artifact.Option{NoProgress: true}
	af, err := artifactImg.NewArtifact(dockerImg, t.cache, artifactOpt, config.ScannerOption{})
	if err != nil {
		return nil, nil, fmt.Errorf("creating artifact: %w", err)
	}

	sc := trivyScanner.NewScanner(t.scanner, af)
	return &sc, cleanup, nil
}

func (t *Trivy) Scan(ctx context.Context, refer *Reference) (*types.Report, error) {
	scanner, cleanup, err := t.newScanner(ctx, refer)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	scanOptions := types.ScanOptions{
		VulnType:            []string{types.VulnTypeOS, types.VulnTypeLibrary},
		SecurityChecks:      []string{types.SecurityCheckVulnerability},
		ScanRemovedPackages: true,
		ListAllPackages:     true,
	}

	// 扫描拿不到镜像size
	// 集群节点上的镜像信息，是在扫描前由informer拿到的
	report, err := scanner.ScanArtifact(ctx, scanOptions)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

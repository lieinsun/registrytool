package scanner

import (
	"context"
	"fmt"
	"os"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	artifactImg "github.com/aquasecurity/fanal/artifact/image"
	fanalImg "github.com/aquasecurity/fanal/image"
	fanalTypes "github.com/aquasecurity/fanal/types"
	trivyScanner "github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
)

type Reference struct {
	// ImageName 扫描所需要的镜像全称 digest为可选字段
	// dockerhub	{username}/{repositoryName}:{tagName}@{digest}
	//				datadog/agent:latest@sha256:775b3a472c0581c4bceaeafabab78c104cc8a5ce806e6ce0634aa4fcf11e41ab
	// 自建仓库		{host}/{projectName}/{repositoryName}:{tagName}@{digest}
	//				127.0.0.1/myProject/myRepo:myTag:@sha256:xxx
	// 本地压缩包		{path}/{fileName}
	//				/tmp/image/775b3a472c0581c4bceaeafabab78c104cc8a5ce806e6ce0634aa4fcf11e41ab.tar
	ImageName    string
	DockerOption fanalTypes.DockerOption
}

var (
	RemoteConf = &scanConfig{
		NewImage: NewRemoteImage,
		artifactOpt: artifact.Option{
			NoProgress: true,
		},
		scanOpt: config.ScannerOption{},
	}
	ArchiveConf = &scanConfig{
		NewImage:    NewArchiveImage,
		artifactOpt: artifact.Option{},
		scanOpt:     config.ScannerOption{},
	}
)

type scanConfig struct {
	NewImage    func(context.Context, *Reference) (fanalTypes.Image, func(), error)
	artifactOpt artifact.Option
	scanOpt     config.ScannerOption
}

func NewRemoteImage(ctx context.Context, refer *Reference) (fanalTypes.Image, func(), error) {
	// TODO 考虑替换NewDockerImage 只处理remote逻辑
	img, cleanup, err := fanalImg.NewDockerImage(ctx, refer.ImageName, refer.DockerOption)
	if err != nil {
		return nil, nil, fmt.Errorf("NewDockerImage err: %w", err)
	}

	return img, cleanup, err
}

func NewArchiveImage(_ context.Context, refer *Reference) (fanalTypes.Image, func(), error) {
	// TODO 考虑替换NewArchiveImage fanal未处理file.Close()
	img, err := fanalImg.NewArchiveImage(refer.ImageName)
	if err != nil {
		return nil, nil, fmt.Errorf("NewArchiveImage err: %w", err)
	}

	return img, func() {
		_ = os.Remove(refer.ImageName)
	}, nil
}

func (t *Trivy) newScanner(ctx context.Context, refer *Reference, config *scanConfig) (*trivyScanner.Scanner, func(), error) {
	image, cleanup, err := config.NewImage(ctx, refer)
	if err != nil {
		return nil, nil, err
	}

	af, err := artifactImg.NewArtifact(image, t.cache, config.artifactOpt, config.scanOpt)
	if err != nil {
		return nil, nil, fmt.Errorf("NewArtifact err: %w", err)
	}

	sc := trivyScanner.NewScanner(t.scanner, af)
	return &sc, cleanup, nil
}

func (t *Trivy) Scan(ctx context.Context, refer *Reference, config *scanConfig) (*types.Report, error) {
	scanner, cleanup, err := t.newScanner(ctx, refer, config)
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

	report, err := scanner.ScanArtifact(ctx, scanOptions)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

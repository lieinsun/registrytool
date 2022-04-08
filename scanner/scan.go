package scanner

import (
	"context"
	"fmt"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	artifactImg "github.com/aquasecurity/fanal/artifact/image"
	"github.com/aquasecurity/fanal/image"
	trivyScanner "github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/lie-inthesun/remotescan/registry"
)

func (t *Trivy) newScanner(ctx context.Context, img registry.Image) (*trivyScanner.Scanner, func(), error) {
	reference, dockerOption := img.TrivyReference()
	// TODO 考虑替换NewDockerImage 只处理remote逻辑
	dockerImg, cleanup, err := image.NewDockerImage(ctx, reference, dockerOption)
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

func (t *Trivy) Scan(ctx context.Context, img registry.Image) (*types.Report, error) {
	scanner, cleanup, err := t.newScanner(ctx, img)
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

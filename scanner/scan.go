package scanner

import (
	"context"
	"fmt"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/aquasecurity/fanal/artifact/image"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
)

var defaultScanOptions = &ScanOptions{
	Option: artifact.Option{
		NoProgress: true,
	},
	ScannerOption: config.ScannerOption{},
	ScanOptions: types.ScanOptions{
		VulnType:            []string{types.VulnTypeOS, types.VulnTypeLibrary},
		SecurityChecks:      []string{types.SecurityCheckVulnerability},
		ScanRemovedPackages: true,
		ListAllPackages:     true,
	},
}

type ScanOptions struct {
	artifact.Option
	config.ScannerOption
	types.ScanOptions
}

func (t *Trivy) newScanner(ctx context.Context, refer Reference, artifactOpt artifact.Option, scannerOpt config.ScannerOption) (*scanner.Scanner, func(), error) {
	img, cleanup, err := refer.NewFanalImage(ctx)
	if err != nil {
		return nil, nil, err
	}

	af, err := image.NewArtifact(img, t.cache, artifactOpt, scannerOpt)
	if err != nil {
		return nil, nil, fmt.Errorf("NewArtifact err: %w", err)
	}

	sc := scanner.NewScanner(t.scanner, af)
	return &sc, cleanup, nil
}

func (t *Trivy) Scan(ctx context.Context, refer Reference, config *ScanOptions) (*types.Report, error) {
	if config == nil {
		config = defaultScanOptions
	}
	sc, cleanup, err := t.newScanner(ctx, refer, config.Option, config.ScannerOption)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	report, err := sc.ScanArtifact(ctx, config.ScanOptions)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

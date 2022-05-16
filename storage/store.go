package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
)

type ImageConfig struct {
	StorePath  string
	AuthConfig authn.AuthConfig
}

func SaveImage(ctx context.Context, repoTag string, config *ImageConfig) (string, error) {
	ref, err := name.ParseReference(repoTag)
	if err != nil {
		return "", err
	}

	opts := []crane.Option{
		crane.WithContext(ctx),
		crane.WithAuth(authn.FromConfig(config.AuthConfig)),
	}
	img, err := crane.Pull(repoTag, opts...)
	if err != nil {
		return "", err
	}

	// save as tarball
	digest, err := img.Digest()
	if err != nil {
		return "", err
	}
	storePath := fmt.Sprintf("%s/%s", config.StorePath, ref.Context().Registry.Name())
	fileName := fmt.Sprintf("%s/%s.tar", storePath, digest.Hex)
	if _, err = os.Stat(storePath); os.IsNotExist(err) {
		err = os.MkdirAll(storePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	if _, err = os.Stat(fileName); os.IsNotExist(err) {
		_, err = os.Create(fileName)
		if err != nil {
			return "", err
		}
	}
	err = crane.Save(img, repoTag, fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

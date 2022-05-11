package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
)

func (c *client) Pull(_ context.Context, imageName string) (string, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return "", err
	}

	img, err := crane.Pull(imageName, crane.WithAuth(authn.FromConfig(c.authConfig)))
	if err != nil {
		return "", err
	}
	digest, err := img.Digest()
	if err != nil {
		return "", err
	}

	// save as tarball
	storePath := fmt.Sprintf("%s/%s", c.storePath, ref.Context().Registry.Name())
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
	err = crane.Save(img, imageName, fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

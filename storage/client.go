package storage

import "github.com/google/go-containerregistry/pkg/authn"

type client struct {
	storePath  string
	authConfig authn.AuthConfig
}

func NewClient(storePath string, authConfig authn.AuthConfig) *client {
	return &client{
		storePath:  storePath,
		authConfig: authConfig,
	}
}

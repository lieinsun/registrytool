package dockerhub

import (
	"fmt"

	"github.com/aquasecurity/fanal/types"

	"remotescan-pkg/registry"
)

type Image struct {
	registry.Auth
	Account     string // dockerhub用户
	Name        string
	Tag         string
	Digest      string
	Size        int
	Os          string
	LastUpdated int64
}

func (i *Image) TrivyReference() (string, types.DockerOption) {
	account := i.Account
	if i.Account == "" {
		account = i.UserName
	}
	if i.Tag == "" {
		i.Tag = "latest"
	}
	ref := fmt.Sprintf("%s/%s:%s", account, i.Name, i.Tag)
	if i.Digest != "" {
		ref = ref + "@" + i.Digest
	}

	dockerOption := types.DockerOption{
		UserName:      i.UserName,
		Password:      i.Password,
		RegistryToken: i.Token,
	}
	return ref, dockerOption
}

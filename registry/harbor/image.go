package harbor

import (
	"fmt"

	"github.com/aquasecurity/fanal/types"

	"remotescan-pkg/registry"
)

type Image struct {
	registry.Auth
	Project  string
	Name     string
	Tag      string
	Digest   string
	Size     int
	Os       string
	PushTime int64
}

func (i *Image) TrivyReference() (string, types.DockerOption) {
	ref := fmt.Sprintf("%s/%s/%s", i.Url, i.Project, i.Name)
	if i.Tag != "" {
		ref = ref + ":" + i.Tag
	}
	if i.Digest != "" {
		ref = ref + "@" + i.Digest
	}

	dockerOption := types.DockerOption{
		UserName: i.UserName,
		Password: i.Password,
	}
	return ref, dockerOption
}

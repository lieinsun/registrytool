package harbor

import "time"

type imageDetailResp struct {
	Digest            string               `json:"digest"`
	ExtraAttrs        imageDetailExtraAttr `json:"extra_attrs"`
	Icon              string               `json:"icon"`
	Id                int                  `json:"id"`
	ManifestMediaType string               `json:"manifest_media_type"`
	MediaType         string               `json:"media_type"`
	ProjectId         int                  `json:"project_id"`
	PullTime          time.Time            `json:"pull_time"`
	PushTime          time.Time            `json:"push_time"`
	RepositoryId      int                  `json:"repository_id"`
	Size              int                  `json:"size"`
	Tags              []imageDetailTag     `json:"tags"`
	//Labels		[]struct{}
	//References	struct{}
}

type imageDetailExtraAttr struct {
	Architecture string    `json:"architecture"`
	Author       string    `json:"author"`
	Created      time.Time `json:"created"`
	Os           string    `json:"os"`
	//Config	struct{}
}

type imageDetailTag struct {
	ArtifactId   int       `json:"artifact_id"`
	Id           int       `json:"id"`
	Immutable    bool      `json:"immutable"`
	Name         string    `json:"name"`
	PullTime     time.Time `json:"pull_time"`
	PushTime     time.Time `json:"push_time"`
	RepositoryId int       `json:"repository_id"`
	Signed       bool      `json:"signed"`
}

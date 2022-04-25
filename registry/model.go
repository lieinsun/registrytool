package registry

type Auth struct {
	Url      string
	UserName string
	Password string
	Token    string
}

type Project struct {
	Name        string            `json:"name"`
	Metadata    map[string]string `json:"metadata"`
	OwnerName   string            `json:"owner_name"`
	RepoCount   int               `json:"repo_count"`
	CreatedTime int64             `json:"creation_time"`
	UpdatedTime int64             `json:"updated_time"`
}

type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	IsPrivate   bool   `json:"is_private"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
}

type Artifact struct {
	Digest      string `json:"digest"`
	Os          string `json:"os"`
	Size        int    `json:"size"`
	UpdatedTime int64  `json:"updated_time"`
	Tags        []Tag  `json:"tags"` // harbor支持查询artifact列表同时拿到内部tag列表
}

type Tag struct {
	Name        string `json:"name"`
	Digest      string `json:"digest"`
	Size        int    `json:"size"`
	UpdatedTime int64  `json:"updated_time"`
}

type Image struct {
	Namespace   string `json:"namespace"` // dockerhub用户或harbor项目名
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Digest      string `json:"digest"`
	Size        int    `json:"size"`
	Os          string `json:"os"`
	UpdatedTime int64  `json:"updated_time"`
}

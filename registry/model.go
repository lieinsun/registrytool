package registry

type Auth struct {
	Url      string
	UserName string
	Password string
	Token    string
}

type Project struct {
	Name         string            `json:"name"`
	Metadata     map[string]string `json:"metadata"`
	OwnerName    string            `json:"owner_name"`
	RepoCount    int               `json:"repo_count"`
	CreationTime int64             `json:"creation_time"`
	UpdatedTime  int64             `json:"updated_time"`
}

type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	UpdatedTime int64  `json:"updated_time"`
}

type Artifact struct {
	Digest      string `json:"digest"`
	Os          string `json:"os"`
	Size        int    `json:"size"`
	UpdatedTime int64  `json:"updated_time"`
}

type Tag struct {
	Name        string `json:"name"`
	Digest      string `json:"digest"`
	Size        int    `json:"size"`
	UpdatedTime int64  `json:"updated_time"`
}

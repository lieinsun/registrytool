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
	Namespace   string
	Name        string
	UpdatedTime int64
}

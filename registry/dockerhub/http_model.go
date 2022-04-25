package dockerhub

import "time"

type tokenResponse struct {
	Detail       string `json:"detail"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type repositoriesResponse struct {
	Count    int                `json:"count"`
	Next     string             `json:"next,omitempty"`
	Previous string             `json:"previous,omitempty"`
	Results  []repositoryResult `json:"results,omitempty"`
}

type repositoryResult struct {
	Name           string    `json:"name"`
	Namespace      string    `json:"namespace"`
	PullCount      int       `json:"pull_count"`
	StarCount      int       `json:"star_count"`
	RepositoryType string    `json:"repository_type"`
	CanEdit        bool      `json:"can_edit"`
	Description    string    `json:"description,omitempty"`
	IsAutomated    bool      `json:"is_automated"`
	IsMigrated     bool      `json:"is_migrated"`
	IsPrivate      bool      `json:"is_private"`
	LastUpdated    time.Time `json:"last_updated"`
	Status         int       `json:"status"`
	User           string    `json:"user"`
}

type imageDetailResponse struct {
	Creator             int        `json:"creator"`
	Id                  int        `json:"id"`
	ImageId             string     `json:"image_id,omitempty"`
	Images              []hubImage `json:"images"`
	LastUpdated         time.Time  `json:"last_updated"`
	LastUpdater         int        `json:"last_updater"`
	LastUpdaterUserName string     `json:"last_updater_username"`
	Name                string     `json:"name"`
	Repository          int        `json:"repository"`
	FullSize            int        `json:"full_size"`
	V2                  bool       `json:"v2"`
	TagStatus           string     `json:"tag_status,omitempty"`
	TagLastPulled       time.Time  `json:"tag_last_pulled"`
	TagLastPushed       time.Time  `json:"tag_last_pushed"`
}

type hubImage struct {
	Architecture string    `json:"architecture"`
	Os           string    `json:"os"`
	Features     string    `json:"features,omitempty"`
	Variant      string    `json:"variant,omitempty"`
	Digest       string    `json:"digest"`
	OsFeatures   string    `json:"os_features,omitempty"`
	OsVersion    string    `json:"os_version,omitempty"`
	Size         int       `json:"size"`
	LastPulled   time.Time `json:"last_pulled,omitempty"`
	LastPushed   time.Time `json:"last_pushed,omitempty"`
	Status       string    `json:"status,omitempty"`
}

type tagsResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next,omitempty"`
	Previous string      `json:"previous,omitempty"`
	Results  []tagResult `json:"results,omitempty"`
}

type tagResult struct {
	Creator             int        `json:"creator"`
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	ImageId             string     `json:"image_id,omitempty"`
	LastUpdated         time.Time  `json:"last_updated"`
	LastUpdater         int        `json:"last_updater"`
	LastUpdaterUserName string     `json:"last_updater_username"`
	Images              []hubImage `json:"images,omitempty"`
	Repository          int        `json:"repository"`
	FullSize            int        `json:"full_size"`
	V2                  bool       `json:"v2"`
	LastPulled          time.Time  `json:"tag_last_pulled,omitempty"`
	LastPushed          time.Time  `json:"tag_last_pushed,omitempty"`
	Status              string     `json:"tag_status,omitempty"`
}

package asana

import (
	"fmt"
	a_types "github.com/leapforce-libraries/go_asana/types"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// StatusUpdate stores StatusUpdate from Service
type StatusUpdate struct {
	Gid             string                 `json:"gid"`
	ResourceType    string                 `json:"resource_type"`
	Title           string                 `json:"title"`
	ResourceSubtype string                 `json:"resource_subtype"`
	Author          Object                 `json:"author"`
	CreatedAt       a_types.DateTimeString `json:"created_at"`
	CreatedBy       Object                 `json:"created_by"`
	Hearted         bool                   `json:"hearted"`
	Hearts          struct {
		User []Object `json:"user"`
	} `json:"hearts"`
	HtmlText string `json:"html_text"`
	Liked    bool   `json:"liked"`
	Likes    struct {
		User []Object `json:"user"`
	} `json:"likes"`
	ModifiedAt a_types.DateTimeString `json:"modified_at"`
	NumHearts  int64                  `json:"num_hearts"`
	NumLikes   int64                  `json:"num_likes"`
	Parent     Object                 `json:"parent"`
	Path       string                 `json:"path"`
	StatusType string                 `json:"status_type"`
	Text       string                 `json:"text"`
	Uri        string                 `json:"uri"`
}

type GetStatusUpdatesConfig struct {
	Parent       string
	CreatedSince *time.Time
	//Fields       *[]string
}

// GetStatusUpdates returns all statusUpdates
func (service *Service) GetStatusUpdates(config *GetStatusUpdatesConfig) ([]StatusUpdate, *errortools.Error) {
	var statusUpdates []StatusUpdate

	params := url.Values{}
	params.Set("parent", config.Parent)
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", StatusUpdate{}))

	if config.CreatedSince != nil {
		params.Set("created_since", (*config.CreatedSince).Format("2006-01-02T15:04:05.999Z"))

	}

	for {
		var statusUpdates_ []StatusUpdate

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("status_updates?%s", params.Encode())),
			ResponseModel: &statusUpdates_,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		statusUpdates = append(statusUpdates, statusUpdates_...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return statusUpdates, nil
}

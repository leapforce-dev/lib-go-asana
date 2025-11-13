package asana

import (
	"cloud.google.com/go/civil"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/url"
)

// TimeTrackingEntry stores TimeTrackingEntry from Service
type TimeTrackingEntry struct {
	Id              string     `json:"gid"`
	ResourceType    string     `json:"resource_type"`
	DurationMinutes int        `json:"duration_minutes"`
	CreatedBy       Object     `json:"created_by"`
	AttributableTo  Object     `json:"attributable_to"`
	EnteredOn       civil.Date `json:"entered_on"`
}

type GetTimeTrackingEntriesConfig struct {
	WorkspaceID        string
	EnteredOnStartDate civil.Date
	EnteredOnEndDate   civil.Date
}

// GetTimeTrackingEntries returns all timeTrackingEntries
func (service *Service) GetTimeTrackingEntries(config *GetTimeTrackingEntriesConfig) ([]TimeTrackingEntry, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be nil")
	}

	var timeTrackingEntries []TimeTrackingEntry

	params := url.Values{}
	params.Set("workspace", config.WorkspaceID)
	params.Set("entered_on_start_date", config.EnteredOnStartDate.String())
	params.Set("entered_on_end_date", config.EnteredOnEndDate.String())

	for {
		var _timeTrackingEntries []TimeTrackingEntry

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("time_tracking_entries?%s", params.Encode())),
			ResponseModel: &_timeTrackingEntries,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		timeTrackingEntries = append(timeTrackingEntries, _timeTrackingEntries...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return timeTrackingEntries, nil
}

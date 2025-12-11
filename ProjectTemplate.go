package asana

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// ProjectTemplate stores ProjectTemplate from Service
type ProjectTemplate struct {
	Id              string          `json:"gid"`
	ResourceType    string          `json:"resource_type"`
	Description     string          `json:"description"`
	HtmlDescription string          `json:"html_description"`
	Team            Object          `json:"team"`
	Public          bool            `json:"public"`
	Color           string          `json:"color"`
	Name            string          `json:"name"`
	Owner           Object          `json:"owner"`
	RequestedRoles  []RequestedRole `json:"requested_roles"`
	RequestedDates  []RequestedDate `json:"requested_dates"`
}

type RequestedRole struct {
	Id   string `json:"gid"`
	Name string `json:"name"`
}

type RequestedDate struct {
	Id          string `json:"gid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetProjectTemplatesConfig struct {
	TeamID      *string
	WorkspaceID *string
	Fields      *string
}

// GetProjectTemplates returns all projectTemplates
func (service *Service) GetProjectTemplates(config *GetProjectTemplatesConfig) ([]ProjectTemplate, *errortools.Error) {
	projectTemplates := []ProjectTemplate{}

	params := url.Values{}
	params.Set("opt_fields", utilities.GetTaggedTagNames("json", ProjectTemplate{}))

	if config != nil {
		if config.WorkspaceID != nil {
			params.Set("workspace", *config.WorkspaceID)
			params.Set("limit", fmt.Sprintf("%v", limitDefault)) // pagination only if workspace is specified
		} else if config.TeamID != nil {
			params.Set("team", fmt.Sprintf("%v", *config.TeamID))
			params.Set("limit", fmt.Sprintf("%v", limitDefault)) // pagination only if workspace is specified
		}

		if config.Fields != nil {
			params.Set("opt_fields", *config.Fields)
		}
	}

	for {
		_projectTemplates := []ProjectTemplate{}

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("project_templates?%s", params.Encode())),
			ResponseModel: &_projectTemplates,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		projectTemplates = append(projectTemplates, _projectTemplates...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return projectTemplates, nil
}

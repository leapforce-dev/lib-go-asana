package asana

import (
	"cloud.google.com/go/civil"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Allocation stores Allocation from Service
type Allocation struct {
	Gid             string `json:"gid"`
	ResourceType    string `json:"resource_type"`
	ResourceSubtype string `json:"resource_subtype"`
	Assignee        struct {
		Gid          string `json:"gid"`
		Name         string `json:"name"`
		ResourceType string `json:"resource_type"`
	} `json:"assignee"`
	Parent struct {
		Gid          string `json:"gid"`
		Name         string `json:"name"`
		ResourceType string `json:"resource_type"`
	} `json:"parent"`
	StartDate civil.Date `json:"start_date"`
	EndDate   civil.Date `json:"end_date"`
	CreatedBy struct {
		Gid          string `json:"gid"`
		Name         string `json:"name"`
		ResourceType string `json:"resource_type"`
	} `json:"created_by"`
	Effort struct {
		Type  string `json:"type"`
		Value int    `json:"value"`
	} `json:"effort"`
}

type GetAllocationsConfig struct {
	Parent string
}

// GetAllocations returns all allocations
func (service *Service) GetAllocations(config *GetAllocationsConfig) ([]Allocation, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be nil")
	}
	allocations := []Allocation{}

	params := url.Values{}
	params.Set("parent", config.Parent)
	params.Set("limit", fmt.Sprintf("%v", limitDefault)) // pagination only if workspace is specified

	for {
		_allocations := []Allocation{}

		requestConfig := go_http.RequestConfig{
			Url:           service.url(fmt.Sprintf("allocations?%s", params.Encode())),
			ResponseModel: &_allocations,
		}
		_, _, nextPage, e := service.getData(&requestConfig)
		if e != nil {
			return nil, e
		}

		allocations = append(allocations, _allocations...)

		if nextPage == nil {
			break
		}
		if nextPage.Offset == "" {
			break
		}

		params.Set("offset", nextPage.Offset)
	}

	return allocations, nil
}

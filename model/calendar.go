package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
)

type CalendarEntry struct {
	URI              string             `json:"uri"`
	DateChanges      []DateChange       `json:"date_changes"`
	Description      ReleaseDescription `json:"description"`
	PublicationState PublicationState   `json:"publication_state"`
}

type ReleaseType struct {
	DataAttributes []coreModel.DataAttribute `json:"data_attributes"`
	ID             string                    `json:"id"`
	Label          coreModel.Localisation    `json:"label"`
	Language       string                    `json:"language"`
	Name           string                    `json:"name"`
	Value          string                    `json:"value"`
	IsChecked      bool                      `json:"is_checked"`
	IsDisabled     bool                      `json:"is_disabled"`
	IsRequired     bool                      `json:"is_required"`
	Count          int                       `json:"count"`
	SubTypes       map[string]ReleaseType    `json:"sub_types"`
}

type SortOption struct {
	LocaleKey string `json:"locale_key"`
	Plural    int    `json:"plural"`
	Value     string `json:"value"`
	Disabled  bool   `json:"disabled"`
}

type Sort struct {
	Mode    string       `json:"mode"`
	Options []SortOption `json:"options"`
}

type DateFieldset struct {
	Input            coreModel.InputDate
	HasValidationErr bool
	ValidationErr    coreModel.Error
}

type Entries struct {
	Count int             `json:"count"`
	Items []CalendarEntry `json:"items"`
}

type Calendar struct {
	coreModel.Page

	ReleaseTypes        map[string]ReleaseType  `json:"release_types"`
	Sort                Sort                    `json:"sort"`
	Keywords            string                  `json:"keywords"`
	BeforeDate          DateFieldset            `json:"before_date"`
	AfterDate           DateFieldset            `json:"after_date"`
	Entries             Entries                 `json:"entries"`
	KeywordSearch       coreModel.CompactSearch `json:"keyword_search"`
	TotalSearchPosition int                     `json:"total_search_position,omitempty"`
}

func (calendar Calendar) FuncIsFilterSearchPresent() bool {
	return calendar.KeywordSearch.SearchTerm != ""
}

func (calendar Calendar) FuncIsFilterCensusPresent() bool {
	for i := range calendar.ReleaseTypes {
		if calendar.ReleaseTypes[i].Name == "census" {
			return calendar.ReleaseTypes[i].IsChecked
		}
	}
	return false
}

func (calendar Calendar) FuncIsFilterDatePresent() bool {
	isBeforeDatePresent := func() bool {
		return calendar.BeforeDate.Input.InputValueDay != "" ||
			calendar.BeforeDate.Input.InputValueMonth != "" ||
			calendar.BeforeDate.Input.InputValueYear != ""
	}

	isAfterDatePresent := func() bool {
		return calendar.AfterDate.Input.InputValueDay != "" ||
			calendar.AfterDate.Input.InputValueMonth != "" ||
			calendar.AfterDate.Input.InputValueYear != ""
	}

	return isBeforeDatePresent() || isAfterDatePresent()
}

package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type CalendarEntry struct {
	URI         string             `json:"uri"`
	DateChanges []DateChange       `json:"date_changes"`
	Description ReleaseDescription `json:"description"`
}

type ReleaseType struct {
	Id        string                 `json:"id"`
	Label     string                 `json:"label"`
	LocaleKey string                 `json:"locale_key"`
	Plural    int                    `json:"plural"`
	Language  string                 `json:"langugage"`
	Name      string                 `json:"name"`
	Value     string                 `json:"value"`
	Checked   bool                   `json:"checked"`
	Count     int                    `json:"count"`
	SubTypes  map[string]ReleaseType `json:"sub_types"`
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

type Calendar struct {
	coreModel.Page

	ReleaseTypes  map[string]ReleaseType  `json:"release_types"`
	Sort          Sort                    `json:"sort"`
	Keywords      string                  `json:"keywords"`
	BeforeDate    coreModel.InputDate     `json:"before_date"`
	AfterDate     coreModel.InputDate     `json:"after_date"`
	Entries       []CalendarEntry         `json:"entries"`
	KeywordSearch coreModel.CompactSearch `json:"keyword_search"`
}

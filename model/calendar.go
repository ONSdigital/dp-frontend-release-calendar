package model

import (
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type CalendarEntry struct {
	URI         string             `json:"uri"`
	DateChanges []DateChange       `json:"date_changes"`
	Description ReleaseDescription `json:"description"`
}

type ReleaseType struct {
	Label    string                 `json:"label"`
	Checked  bool                   `json:"value"`
	Count    int                    `json:"count"`
	SubTypes map[string]ReleaseType `json:"sub_types"`
}

type SortOption = queryparams.SortOption

type Sort struct {
	Mode    string       `json:"mode"`
	Options []SortOption `json:"options"`
}

type Date struct {
	Day   string `json:"day"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

type Calendar struct {
	coreModel.Page
	ReleaseTypes  map[string]ReleaseType  `json:"release_types"`
	Sort          Sort                    `json:"sort"`
	Keywords      string                  `json:"keywords"`
	BeforeDate    Date                    `json:"before_date"`
	AfterDate     Date                    `json:"after_date"`
	Entries       []CalendarEntry         `json:"entries"`
	KeywordSearch coreModel.CompactSearch `json:"keyword_search"`
}

package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type CalendarPagination struct {
	CurrentPage  int            `json:"current_page"`
	CalendarItem []CalendarItem `json:"calendar_item"`
	TotalPages   int            `json:"total_pages"`
	Limit        int            `json:"limit"`
}

type CalendarItem struct {
	URI         string             `json:"uri"`
	Description ReleaseDescription `json:"description"`
}

type ReleaseType struct {
	Label    string                 `json:"label"`
	Checked  bool                   `json:"value"`
	Count    int                    `json:"count"`
	SubTypes map[string]ReleaseType `json:"sub_types"`
}

type SortOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
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
	ReleaseTypes       map[string]ReleaseType  `json:"release_types"`
	Sort               Sort                    `json:"sort"`
	Keywords           string                  `json:"keywords"`
	BeforeDate         Date                    `json:"before_date"`
	AfterDate          Date                    `json:"after_date"`
	CalendarPagination CalendarPagination      `json:"calendar_pagination"`
	KeywordSearch      coreModel.CompactSearch `json:"keyword_search"`
}

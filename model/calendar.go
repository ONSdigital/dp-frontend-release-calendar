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

type Filter struct {
	Name  string
	Value string
}

type SortOption struct {
	Label string
	Value string
}
type Sort struct {
	Mode    string
	Options []SortOption
}

type Calendar struct {
	coreModel.Page
	Filters            []Filter           `json:"filters"`
	Sort               Sort               `json:"sort"`
	Keywords           string             `json:"keywords"`
	CalendarPagination CalendarPagination `json:"calendar_pagination"`
}

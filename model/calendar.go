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

type Calendar struct {
	coreModel.Page
	CalendarPagination CalendarPagination `json:"calendar_pagination"`
}

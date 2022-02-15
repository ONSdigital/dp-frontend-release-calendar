package model

import (
	"time"

	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type CalendarPagination struct {
	CurrentPage  int            `json:"current_page"`
	CalendarItem []CalendarItem `json:"calendar_item"`
	TotalPages   int            `json:"total_pages"`
	Limit        int            `json:"limit"`
}

type CalendarItem struct {
	URI            string         `json:"uri"`
	Title          string         `json:"title"`
	Summary        string         `json:"summary"`
	ReleaseDate    time.Time      `json:"releaseDate"`
	Published      bool           `json:"published"`
	Cancelled      bool           `json:"cancelled"`
	ContactDetails ContactDetails `json:"contact_details"`
	NextRelease    string         `json:"nextRelease"`
}

type Calendar struct {
	coreModel.Page
	CalendarPagination CalendarPagination `json:"calendar_pagination"`
}

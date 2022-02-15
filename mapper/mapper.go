package mapper

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

func CreateRelease(ctx context.Context, basePage coreModel.Page, cfg config.Config) model.Release {
	release := model.Release{
		Page: basePage,
	}

	release.BetaBannerEnabled = true
	release.Metadata.Title = "Test Release Page"
	release.ContactDetails.Name = "Test contact name"

	return release
}

func CreateCalendar(ctx context.Context, basePage coreModel.Page, cfg config.Config) model.Calendar {
	calendar := model.Calendar{
		Page: basePage,
	}
	calendar.BetaBannerEnabled = true
	calendar.Metadata.Title = "Test Release Calendar"

	item1 := model.CalendarItem{
		URI:            "/releases/title1",
		Title:          "Title 1",
		Summary:        "A summary for Title 1",
		ReleaseDate:    time.Now().AddDate(0, 0, -10),
		Published:      true,
		Cancelled:      false,
		ContactDetails: model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		NextRelease:    "To be announced",
	}

	item2 := model.CalendarItem{
		URI:            "/releases/title2",
		Title:          "Title 2",
		Summary:        "A summary for Title 2",
		ReleaseDate:    time.Now().AddDate(0, 0, -15),
		Published:      false,
		Cancelled:      true,
		ContactDetails: model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
	}

	item3 := model.CalendarItem{
		URI:            "/releases/title3",
		Title:          "Title 3",
		Summary:        "A summary for Title 3",
		ReleaseDate:    time.Now().AddDate(0, 0, 5),
		Published:      false,
		Cancelled:      false,
		ContactDetails: model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
	}

	calendar.CalendarPagination.CurrentPage = 1
	calendar.CalendarPagination.TotalPages = 100
	calendar.CalendarPagination.Limit = 10
	calendar.CalendarPagination.CalendarItem = make([]model.CalendarItem, 3)
	calendar.CalendarPagination.CalendarItem[0] = item1
	calendar.CalendarPagination.CalendarItem[1] = item2
	calendar.CalendarPagination.CalendarItem[2] = item3

	return calendar
}

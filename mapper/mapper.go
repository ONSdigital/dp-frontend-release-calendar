package mapper

import (
	"context"

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
		URL:         "localhost",
		Title:       "title1",
		Description: "description1",
	}

	item2 := model.CalendarItem{
		URL:         "localhost",
		Title:       "title2",
		Description: "description2",
	}

	calendar.CalendarPagination.CurrentPage = 1
	calendar.CalendarPagination.TotalPages = 100
	calendar.CalendarPagination.Limit = 10
	calendar.CalendarPagination.CalendarItem = make([]model.CalendarItem, 2)
	calendar.CalendarPagination.CalendarItem[0] = item1
	calendar.CalendarPagination.CalendarItem[1] = item2

	return calendar
}

package mapper

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

const SixteensVersion = "77f1d9b"

func CreateRelease(ctx context.Context, basePage coreModel.Page, cfg config.Config) model.Release {
	release := model.Release{
		Page:     basePage,
		Markdown: []string{"markdown 1", "markdown 2"},
		RelatedDocuments: []model.Link{
			{
				Title:   "Document 1",
				Summary: "This is document 1",
				URI:     "/doc/1",
			},
		},
		RelatedDatasets: []model.Link{
			{
				Title:   "Dataset 1",
				Summary: "This is dataset 1",
				URI:     "/dataset/1",
			},
		},
		RelatedMethodology: []model.Link{
			{
				Title:   "Methodology",
				Summary: "This is methodology 1",
				URI:     "/methodology/1",
			},
		},
		RelatedMethodologyArticle: []model.Link{
			{
				Title:   "Methodology Article",
				Summary: "This is methodology article 1",
				URI:     "/methodology/article/1",
			},
		},
		Links: []model.Link{
			{
				Title:   "Link 1",
				Summary: "This is link 1",
				URI:     "/link/1",
			},
		},
		DateChanges: []model.DateChange{
			{
				Date:         "2022-02-15T11:12:05.592Z",
				ChangeNotice: "This release has changed",
			},
		},
		Description: model.ReleaseDescription{
			Title:   "Release title",
			Summary: "Release summary",
			Contact: model.ContactDetails{
				Email:     "contact@ons.gov.uk",
				Name:      "Contact name",
				Telephone: "029",
			},
			NationalStatistic:  true,
			ReleaseDate:        "2020-07-08T23:00:00.000Z",
			NextRelease:        "January 2021",
			Published:          true,
			Finalised:          true,
			Cancelled:          true,
			CancellationNotice: []string{"cancelled for a reason"},
			ProvisionalDate:    "July 2020",
		},
	}

	release.FeatureFlags.SixteensVersion = SixteensVersion
	release.BetaBannerEnabled = true
	release.Metadata.Title = "Test Release Page"

	return release
}

func CreateCalendar(ctx context.Context, basePage coreModel.Page, cfg config.Config) model.Calendar {
	calendar := model.Calendar{
		Page: basePage,
	}
	calendar.FeatureFlags.SixteensVersion = SixteensVersion
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

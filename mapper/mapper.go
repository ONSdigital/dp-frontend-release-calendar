package mapper

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
)

const SixteensVersion = "77f1d9b"

func CreateRelease(basePage coreModel.Page, release releasecalendar.Release) model.Release {
	result := model.Release{
		Page:     basePage,
		Markdown: release.Markdown,
		Description: model.ReleaseDescription{
			Title:   release.Description.Title,
			Summary: release.Description.Summary,
			Contact: model.ContactDetails{
				Email:     release.Description.Contact.Email,
				Name:      release.Description.Contact.Name,
				Telephone: release.Description.Contact.Telephone,
			},
			NationalStatistic:  release.Description.NationalStatistic,
			ReleaseDate:        release.Description.ReleaseDate,
			NextRelease:        release.Description.NextRelease,
			Published:          release.Description.Published,
			Finalised:          release.Description.Finalised,
			Cancelled:          release.Description.Cancelled,
			CancellationNotice: release.Description.CancellationNotice,
			ProvisionalDate:    release.Description.ProvisionalDate,
		},
	}

	result.RelatedDatasets = mapLink(release.RelatedDatasets)
	result.RelatedDocuments = mapLink(release.RelatedDocuments)
	result.RelatedMethodology = mapLink(release.RelatedMethodology)
	result.RelatedMethodologyArticle = mapLink(release.RelatedMethodologyArticle)
	result.Links = mapLink(release.Links)

	result.DateChanges = []model.DateChange{}
	for _, dc := range release.DateChanges {
		result.DateChanges = append(result.DateChanges, model.DateChange{
			Date:         dc.Date,
			ChangeNotice: dc.ChangeNotice,
		})
	}

	result.FeatureFlags.SixteensVersion = SixteensVersion
	result.BetaBannerEnabled = true
	result.Metadata.Title = release.Description.Title
	result.URI = release.URI
	return result
}

func mapLink(links []releasecalendar.Link) []model.Link {
	res := []model.Link{}
	for _, l := range links {
		res = append(res, model.Link{
			Title:   l.Title,
			Summary: l.Summary,
			URI:     l.URI,
		})
	}
	return res
}

func CreateCalendar(_ context.Context, basePage coreModel.Page, _ config.Config) model.Calendar {
	calendar := model.Calendar{
		Page: basePage,
	}
	calendar.FeatureFlags.SixteensVersion = SixteensVersion
	calendar.BetaBannerEnabled = true
	calendar.Metadata.Title = "Test Release Calendar"

	item1 := model.CalendarItem{
		URI: "/releases/title1",
		Description: model.ReleaseDescription{
			Title:       "Title 1",
			Summary:     "A summary for Title 1",
			ReleaseDate: time.Now().AddDate(0, 0, -10).UTC().Format(time.RFC3339),
			Published:   true,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
			NextRelease: "To be announced",
		},
	}

	item2 := model.CalendarItem{
		URI: "/releases/title2",
		Description: model.ReleaseDescription{
			Title:       "Title 2",
			Summary:     "A summary for Title 2",
			ReleaseDate: time.Now().AddDate(0, 0, -15).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   true,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item3 := model.CalendarItem{
		URI: "/releases/title3",
		Description: model.ReleaseDescription{
			Title:       "Title 3",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
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

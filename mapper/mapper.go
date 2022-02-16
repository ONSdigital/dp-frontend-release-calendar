package mapper

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
)

func CreatePreviousReleases(_ context.Context, basePage coreModel.Page, _ config.Config) model.PreviousReleases {
	previousReleases := model.PreviousReleases{
		Page: basePage,
	}

	previousReleases.BetaBannerEnabled = true
	previousReleases.Metadata.Title = "Personal well-being in the UK, quarterly: July 2021 to September 2021"
	previousReleases.ContactDetails.Name = "Test contact name"

	return previousReleases
}

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
	calendar.BetaBannerEnabled = true
	calendar.Metadata.Title = "Release Calendar"

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
		// URL:          "localhost",
		// Title:        "Public Sector Employment, UK: September 2021",
		// Description:  "description1",
		// ReleaseDate:  "10 December 2021 9:30am",
		// ReleaseState: "Published",
	}

	// item2 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Labour market in the regions of the UK: December 2021",
	// 	Description:  "description2",
	// 	ReleaseDate:  "10 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item3 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Personal well-being in the UK, quarterly: July 2021 to September 2021",
	// 	Description:  "description3",
	// 	ReleaseDate:  "10 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item4 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Labour market statistics time series: December 2021",
	// 	Description:  "description4",
	// 	ReleaseDate:  "10 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item5 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "UK labour market: December 2021",
	// 	Description:  "description5",
	// 	ReleaseDate:  "10 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item6 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Earnings and employment from Pay As You Earn Real Time Information, UK: December 2021",
	// 	Description:  "description6",
	// 	ReleaseDate:  "10 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item7 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Civil partnerships in England and Wales: 2020",
	// 	Description:  "description7",
	// 	ReleaseDate:  "9 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item8 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Understanding towns: industry analysis",
	// 	Description:  "description8",
	// 	ReleaseDate:  "9 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item9 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Disaggregating annual subnational gross value added (GVA) to lower levels of geography",
	// 	Description:  "description9",
	// 	ReleaseDate:  "9 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	// item10 := model.CalendarItem{
	// 	URL:          "localhost",
	// 	Title:        "Coronavirus (COVID-19) Infection Survey, UK: 8 December 2021",
	// 	Description:  "description10",
	// 	ReleaseDate:  "9 December 2021 9:30am",
	// 	ReleaseState: "Published",
	// }

	calendar.CalendarPagination.CurrentPage = 1
	calendar.CalendarPagination.TotalPages = 100
	calendar.CalendarPagination.Limit = 10
	calendar.CalendarPagination.CalendarItem = make([]model.CalendarItem, 10)
	calendar.CalendarPagination.CalendarItem[0] = item1
	calendar.CalendarPagination.CalendarItem[1] = item2
	calendar.CalendarPagination.CalendarItem[2] = item3
	// calendar.CalendarPagination.CalendarItem[3] = item4
	// calendar.CalendarPagination.CalendarItem[4] = item5
	// calendar.CalendarPagination.CalendarItem[5] = item6
	// calendar.CalendarPagination.CalendarItem[6] = item7
	// calendar.CalendarPagination.CalendarItem[7] = item8
	// calendar.CalendarPagination.CalendarItem[8] = item9
	// calendar.CalendarPagination.CalendarItem[9] = item10

	return calendar
}

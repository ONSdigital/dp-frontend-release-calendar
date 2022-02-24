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

	previousReleases.Description.Title = "Personal well-being in the UK, quarterly: July 2021 to September 2021"
	previousReleases.Description.Summary = "Quarterly estimates of life satisfaction, feeling that the things done in life are worthwhile, happiness and anxiety at the UK level, created using the Annual Population Survey (APS)."

	previousReleases.ReleaseHistory = make([]model.Link, 10)
	previousReleases.ReleaseHistory[0] = model.Link{
		Title:   "8 December 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[1] = model.Link{
		Title:   "26 July 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[2] = model.Link{
		Title:   "22 April 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[3] = model.Link{
		Title:   "21 January 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[4] = model.Link{
		Title:   "17 October 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[5] = model.Link{
		Title:   "21 July 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[6] = model.Link{
		Title:   "23 April 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[7] = model.Link{
		Title:   "28 January 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[8] = model.Link{
		Title:   "18 October 2019",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[9] = model.Link{
		Title:   "17 July 2019",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}

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
			Title:       "Public Sector Employment, UK: September 2021",
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
			Title:       "Labour market in the regions of the UK: December 2021",
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
			Title:       "Personal well-being in the UK, quarterly: July 2021 to September 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item4 := model.CalendarItem{
		URI: "/releases/title4",
		Description: model.ReleaseDescription{
			Title:       "Labour market statistics time series: December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item5 := model.CalendarItem{
		URI: "/releases/title5",
		Description: model.ReleaseDescription{
			Title:       "UK labour market: December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item6 := model.CalendarItem{
		URI: "/releases/title6",
		Description: model.ReleaseDescription{
			Title:       "Earnings and employment from Pay As You Earn Real Time Information, UK: December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item7 := model.CalendarItem{
		URI: "/releases/title7",
		Description: model.ReleaseDescription{
			Title:       "Civil partnerships in England and Wales: 2020",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item8 := model.CalendarItem{
		URI: "/releases/title8",
		Description: model.ReleaseDescription{
			Title:       "Understanding towns: industry analysis",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item9 := model.CalendarItem{
		URI: "/releases/title9",
		Description: model.ReleaseDescription{
			Title:       "Disaggregating annual subnational gross value added (GVA) to lower levels of geography",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item10 := model.CalendarItem{
		URI: "/releases/title10",
		Description: model.ReleaseDescription{
			Title:       "Coronavirus (COVID-19) Infection Survey, UK: 8 December 2021",
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
	calendar.CalendarPagination.CalendarItem = make([]model.CalendarItem, 10)
	calendar.CalendarPagination.CalendarItem[0] = item1
	calendar.CalendarPagination.CalendarItem[1] = item2
	calendar.CalendarPagination.CalendarItem[2] = item3
	calendar.CalendarPagination.CalendarItem[3] = item4
	calendar.CalendarPagination.CalendarItem[4] = item5
	calendar.CalendarPagination.CalendarItem[5] = item6
	calendar.CalendarPagination.CalendarItem[6] = item7
	calendar.CalendarPagination.CalendarItem[7] = item8
	calendar.CalendarPagination.CalendarItem[8] = item9
	calendar.CalendarPagination.CalendarItem[9] = item10

	return calendar
}

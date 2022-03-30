package mapper

import (
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	sitesearch "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	Convey("Given a release and a base page", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		release := releasecalendar.Release{
			URI:      "/releases/test",
			Markdown: []string{"markdown1", "markdown 2"},
			RelatedDocuments: []releasecalendar.Link{
				{
					Title:   "Document 1",
					Summary: "This is document 1",
					URI:     "/doc/1",
				},
				{
					Title:   "Document 2",
					Summary: "This is document 2",
					URI:     "/doc/2",
				},
			},
			RelatedDatasets: []releasecalendar.Link{
				{
					Title:   "Dataset 1",
					Summary: "This is dataset 1",
					URI:     "/dataset/1",
				},
				{
					Title:   "Dataset 2",
					Summary: "This is dataset 2",
					URI:     "/dataset/2",
				},
			},
			RelatedMethodology: []releasecalendar.Link{
				{
					Title:   "Methodology",
					Summary: "This is methodology 1",
					URI:     "/methodology/1",
				},
				{
					Title:   "Methodology 2",
					Summary: "This is methodology 2",
					URI:     "/methodology/2",
				},
			},
			RelatedMethodologyArticle: []releasecalendar.Link{
				{
					Title:   "Methodology Article",
					Summary: "This is methodology article 1",
					URI:     "/methodology/article/1",
				},
				{
					Title:   "Methodology Article 2",
					Summary: "This is methodology article 2",
					URI:     "/methodology/article/2",
				},
			},
			Links: []releasecalendar.Link{
				{
					Title:   "Link 1",
					Summary: "This is link 1",
					URI:     "/link/1",
				},
				{
					Title:   "Link 2",
					Summary: "This is link 2",
					URI:     "/link/2",
				},
			},
			DateChanges: []releasecalendar.ReleaseDateChange{
				{
					Date:         "2022-02-15T11:12:05.592Z",
					ChangeNotice: "This release has changed",
				},
				{
					Date:         "2022-02-22T22:02:22.202Z",
					ChangeNotice: "Yet another change",
				},
			},
			Description: releasecalendar.ReleaseDescription{
				Title:   "Release title",
				Summary: "Release summary",
				Contact: releasecalendar.Contact{
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
				CancellationNotice: []string{"cancelled for a reason", "another reason"},
				ProvisionalDate:    "July 2020",
			},
		}

		Convey("CreateRelease maps correctly to a model object", func() {
			model := CreateRelease(basePage, release)

			So(model.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
			So(model.SiteDomain, ShouldEqual, basePage.SiteDomain)
			So(model.BetaBannerEnabled, ShouldBeTrue)
			So(model.Metadata.Title, ShouldEqual, release.Description.Title)
			So(model.URI, ShouldEqual, release.URI)
			So(model.Markdown, ShouldResemble, release.Markdown)
			assertLinks(release.RelatedDatasets, model.RelatedDatasets)
			assertLinks(release.RelatedDocuments, model.RelatedDocuments)
			assertLinks(release.RelatedMethodology, model.RelatedMethodology)
			assertLinks(release.RelatedMethodologyArticle, model.RelatedMethodologyArticle)
			assertLinks(release.Links, model.Links)
			assertDateChanges(release.DateChanges, model.DateChanges)
			So(model.Description.Title, ShouldEqual, release.Description.Title)
			So(model.Description.Summary, ShouldEqual, release.Description.Summary)
			So(model.Description.Contact.Name, ShouldEqual, release.Description.Contact.Name)
			So(model.Description.Contact.Email, ShouldEqual, release.Description.Contact.Email)
			So(model.Description.Contact.Telephone, ShouldEqual, release.Description.Contact.Telephone)
			So(model.Description.NationalStatistic, ShouldEqual, release.Description.NationalStatistic)
			So(model.Description.ReleaseDate, ShouldEqual, release.Description.ReleaseDate)
			So(model.Description.Published, ShouldEqual, release.Description.Published)
			So(model.Description.Finalised, ShouldEqual, release.Description.Finalised)
			So(model.Description.Cancelled, ShouldEqual, release.Description.Cancelled)
			So(model.Description.CancellationNotice, ShouldResemble, release.Description.CancellationNotice)
			So(model.Description.ProvisionalDate, ShouldEqual, release.Description.ProvisionalDate)
		})
	})
}

func TestReleaseCalendarMapper(t *testing.T) {
	Convey("Given a Release Calendar and a base page", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		releaseResponse := sitesearch.ReleaseResponse{
			Took: 100,
			Breakdown: sitesearch.Breakdown{
				Total: 11,
			},
			Releases: []sitesearch.Release{
				{
					URI: "/releases/title1",
					DateChanges: []sitesearch.ReleaseDateChange{
						{
							Date:         "2015-09-22T12:30:23.221Z",
							ChangeNotice: "Something happened to change the date",
						},
					},
					Description: sitesearch.ReleaseDescription{
						Title:       "Public Sector Employment, UK: September 2021",
						Summary:     "A summary for Title 1",
						ReleaseDate: time.Now().AddDate(0, 0, -10).UTC().Format(time.RFC3339),
						Published:   true,
						Finalised:   true,
						Postponed:   true,
						Contact:     &sitesearch.Contact{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
						NextRelease: "To be announced",
					},
				},
				{
					URI: "/releases/title2",
					Description: sitesearch.ReleaseDescription{
						Title:       "Labour market in the regions of the UK: December 2021",
						Summary:     "A summary for Title 2",
						ReleaseDate: time.Now().AddDate(0, 0, -15).UTC().Format(time.RFC3339),
						Published:   false,
						Cancelled:   true,
						Contact:     &sitesearch.Contact{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
					},
				},
				{
					URI: "/releases/title3",
					Description: sitesearch.ReleaseDescription{
						Title:       "Personal well-being in the UK, quarterly: July 2021 to September 2021",
						Summary:     "A summary for Title 3",
						ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
						Published:   false,
						Cancelled:   false,
					},
				},
				{
					URI: "/releases/title4",
					Description: sitesearch.ReleaseDescription{
						Title:       "Labour market statistics time series: December 2021",
						Summary:     "A summary for Title 3",
						ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
						Published:   false,
						Cancelled:   false,
						Contact:     &sitesearch.Contact{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
					},
				},
				{
					URI: "/releases/title5",
					Description: sitesearch.ReleaseDescription{
						Title:       "UK labour market: December 2021",
						Summary:     "A summary for Title 3",
						ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
						Published:   false,
						Cancelled:   false,
						Contact:     &sitesearch.Contact{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
					},
				},
			},
		}

		params := queryparams.ValidatedParams{
			Limit:      5,
			Offset:     0,
			AfterDate:  queryparams.Date{},
			BeforeDate: queryparams.Date{},
			Keywords:   "everything",
			Sort:       queryparams.RelDateAsc,
			Upcoming:   true,
		}

		Convey("CreateReleaseCalendar maps correctly to a model Calendar object", func() {
			calendar := CreateReleaseCalendar(basePage, params, releaseResponse)

			So(calendar.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
			So(calendar.SiteDomain, ShouldEqual, basePage.SiteDomain)
			So(calendar.BetaBannerEnabled, ShouldBeTrue)
			So(calendar.Metadata.Title, ShouldEqual, "Release Calendar")
			So(calendar.Keywords, ShouldEqual, params.Keywords)
			So(calendar.Sort, ShouldResemble, model.Sort{Mode: params.Sort.String(), Options: queryparams.SortOptions})
			So(calendar.BeforeDate, ShouldResemble, model.Date{Day: params.BeforeDate.DayString(), Month: params.BeforeDate.MonthString(), Year: params.BeforeDate.YearString()})
			So(calendar.AfterDate, ShouldResemble, model.Date{Day: params.AfterDate.DayString(), Month: params.AfterDate.MonthString(), Year: params.AfterDate.YearString()})
			So(calendar.ReleaseTypes, ShouldResemble, mapReleases(params, releaseResponse))
			So(calendar.Pagination.TotalPages, ShouldEqual, 3)
			So(calendar.Pagination.CurrentPage, ShouldEqual, 1)
			So(calendar.Pagination.Limit, ShouldEqual, 5)
			for i, r := range calendar.Entries {
				So(r.URI, ShouldEqual, releaseResponse.Releases[i].URI)
				assertSiteSearchDateChanges(releaseResponse.Releases[i].DateChanges, r.DateChanges)
				So(r.Description.Title, ShouldEqual, releaseResponse.Releases[i].Description.Title)
				So(r.Description.Summary, ShouldEqual, releaseResponse.Releases[i].Description.Summary)
				So(r.Description.NationalStatistic, ShouldEqual, releaseResponse.Releases[i].Description.NationalStatistic)
				So(r.Description.ReleaseDate, ShouldEqual, releaseResponse.Releases[i].Description.ReleaseDate)
				So(r.Description.Published, ShouldEqual, releaseResponse.Releases[i].Description.Published)
				So(r.Description.Finalised, ShouldEqual, releaseResponse.Releases[i].Description.Finalised)
				So(r.Description.Cancelled, ShouldEqual, releaseResponse.Releases[i].Description.Cancelled)
				So(r.Description.CancellationNotice, ShouldResemble, releaseResponse.Releases[i].Description.CancellationNotice)
				So(r.Description.ProvisionalDate, ShouldEqual, releaseResponse.Releases[i].Description.ProvisionalDate)
				if releaseResponse.Releases[i].Description.Contact != nil {
					So(r.Description.Contact.Name, ShouldEqual, releaseResponse.Releases[i].Description.Contact.Name)
					So(r.Description.Contact.Email, ShouldEqual, releaseResponse.Releases[i].Description.Contact.Email)
					So(r.Description.Contact.Telephone, ShouldEqual, releaseResponse.Releases[i].Description.Contact.Telephone)
				} else {
					So(r.Description.Contact, ShouldBeZeroValue)
				}
			}
		})
	})
}

// assertLinks checks that the actual model Link content is equal to the expected release Link
func assertLinks(expected []releasecalendar.Link, actual []model.Link) {
	So(len(actual), ShouldEqual, len(expected))
	for i := range expected {
		So(actual[i].URI, ShouldEqual, expected[i].URI)
		So(actual[i].Title, ShouldEqual, expected[i].Title)
		So(actual[i].Summary, ShouldEqual, expected[i].Summary)
	}
}

// assertDateChanges checks that the actual model DateChanges content is equal to the expected release ReleaseDateChanges
func assertDateChanges(expected []releasecalendar.ReleaseDateChange, actual []model.DateChange) {
	So(len(actual), ShouldEqual, len(expected))
	for i := range expected {
		So(actual[i].Date, ShouldEqual, expected[i].Date)
		So(actual[i].ChangeNotice, ShouldEqual, expected[i].ChangeNotice)
	}
}

func assertSiteSearchDateChanges(expected []sitesearch.ReleaseDateChange, actual []model.DateChange) {
	So(len(actual), ShouldEqual, len(expected))
	for i := range expected {
		So(actual[i].Date, ShouldEqual, expected[i].Date)
		So(actual[i].ChangeNotice, ShouldEqual, expected[i].ChangeNotice)
	}
}

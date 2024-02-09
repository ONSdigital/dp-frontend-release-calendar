package mapper

import (
	"errors"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	sitesearch "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mocks"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a release and a base page", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		releaseResponse := releasecalendar.Release{
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
				Survey:             "census",
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
			lang := "cy"
			crumbLabelHome := "Hafan"
			crumbLabelReleaseCalendar := "Calendr datganiadau"
			crumbLabelCancelled := "Canslwyd"
			serviceMessage := "Service Message"
			emergencyBannerTitle := "Emergency Title"
			emergencyBannerType := "notable-death"
			emergencyBannerDescription := "Emergency Description"
			emergencyBannerURI := "https://example.com/emergency"
			emergencyBannerLinkText := "Attention, this is an emergency. There's an emergency going on."
			bannerData := zebedee.EmergencyBanner{
				Title:       emergencyBannerTitle,
				Type:        emergencyBannerType,
				Description: emergencyBannerDescription,
				URI:         emergencyBannerURI,
				LinkText:    emergencyBannerLinkText,
			}
			release := CreateRelease(basePage, releaseResponse, lang, "/prefix/releasecalendar", serviceMessage, bannerData)

			So(release.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
			So(release.SiteDomain, ShouldEqual, basePage.SiteDomain)
			So(release.BetaBannerEnabled, ShouldBeTrue)
			So(release.ServiceMessage, ShouldEqual, serviceMessage)
			So(release.EmergencyBanner.Title, ShouldEqual, emergencyBannerTitle)
			So(release.EmergencyBanner.Type, ShouldEqual, emergencyBannerType)
			So(release.EmergencyBanner.Description, ShouldEqual, emergencyBannerDescription)
			So(release.EmergencyBanner.URI, ShouldEqual, emergencyBannerURI)
			So(release.EmergencyBanner.LinkText, ShouldEqual, emergencyBannerLinkText)
			So(release.Metadata.Title, ShouldEqual, releaseResponse.Description.Title)
			So(release.URI, ShouldEqual, releaseResponse.URI)
			So(release.Markdown, ShouldResemble, releaseResponse.Markdown)
			assertLinks(releaseResponse.RelatedDatasets, release.RelatedDatasets)
			assertLinks(releaseResponse.RelatedDocuments, release.RelatedDocuments)
			assertLinks(releaseResponse.RelatedMethodology, release.RelatedMethodology)
			assertLinks(releaseResponse.RelatedMethodologyArticle, release.RelatedMethodologyArticle)
			assertLinks(releaseResponse.RelatedAPIDatasets, release.RelatedAPIDatasets)
			assertLinks(releaseResponse.Links, release.Links)
			assertDateChanges(releaseResponse.DateChanges, release.DateChanges)
			So(release.Description.Title, ShouldEqual, releaseResponse.Description.Title)
			So(release.Description.Summary, ShouldEqual, releaseResponse.Description.Summary)
			So(release.Description.Contact.Name, ShouldEqual, releaseResponse.Description.Contact.Name)
			So(release.Description.Contact.Email, ShouldEqual, releaseResponse.Description.Contact.Email)
			So(release.Description.Contact.Telephone, ShouldEqual, releaseResponse.Description.Contact.Telephone)
			So(release.Description.NationalStatistic, ShouldEqual, releaseResponse.Description.NationalStatistic)
			So(release.Description.WelshStatistic, ShouldEqual, releaseResponse.Description.WelshStatistic)
			So(release.Description.Census2021, ShouldEqual, releaseResponse.Census())
			So(release.Description.ReleaseDate, ShouldEqual, releaseResponse.Description.ReleaseDate)
			So(release.Description.Published, ShouldEqual, releaseResponse.Description.Published)
			So(release.Description.Finalised, ShouldEqual, releaseResponse.Description.Finalised)
			So(release.Description.Cancelled, ShouldEqual, releaseResponse.Description.Cancelled)
			So(release.Description.CancellationNotice, ShouldResemble, releaseResponse.Description.CancellationNotice)
			So(release.Description.ProvisionalDate, ShouldEqual, releaseResponse.Description.ProvisionalDate)
			So(release.Breadcrumb, ShouldResemble, []coreModel.TaxonomyNode{
				{
					Title: crumbLabelHome,
					URI:   "/",
				},
				{
					Title: crumbLabelReleaseCalendar,
					URI:   "/prefix/releasecalendar",
				},
				{
					Title: crumbLabelCancelled,
					URI:   "/prefix/releasecalendar?release-type=type-cancelled",
				},
			})
			So(release.PublicationState, ShouldResemble, model.PublicationState{
				Type: "cancelled",
			})
		})
	})
}

func TestReleaseCalendarMapper(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a Release Calendar and a base page", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		releaseResponse := sitesearch.ReleaseResponse{
			Took: 100,
			Breakdown: sitesearch.Breakdown{
				Total:       11,
				Provisional: 1,
				Confirmed:   4,
				Published:   5,
				Cancelled:   1,
				Census:      3,
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
						Finalised:   true,
						Census:      true,
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
					},
				},
			},
		}

		params := queryparams.ValidatedParams{
			Limit:       5,
			Offset:      0,
			Page:        1,
			AfterDate:   queryparams.Date{},
			BeforeDate:  queryparams.Date{},
			Keywords:    "everything",
			Sort:        queryparams.RelDateAsc,
			ReleaseType: queryparams.Upcoming,
		}

		cfg := config.Config{DefaultMaximumSearchResults: 1000}

		Convey("CreateReleaseCalendar maps correctly to a model Calendar object", func() {
			lang := "cy"
			metaTitle := "Calendr datganiadau"
			serviceMessage := "Service Message"
			emergencyBannerTitle := "Emergency Title"
			emergencyBannerType := "notable-death"
			emergencyBannerDescription := "Emergency Description"
			emergencyBannerURI := "https://example.com/emergency"
			emergencyBannerLinkText := "Attention, this is an emergency. There's an emergency going on."
			bannerData := zebedee.EmergencyBanner{
				Title:       emergencyBannerTitle,
				Type:        emergencyBannerType,
				Description: emergencyBannerDescription,
				URI:         emergencyBannerURI,
				LinkText:    emergencyBannerLinkText,
			}

			calendar := CreateReleaseCalendar(basePage, params, releaseResponse, cfg, lang, serviceMessage, bannerData, nil)

			So(calendar.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
			So(calendar.SiteDomain, ShouldEqual, basePage.SiteDomain)
			So(calendar.BetaBannerEnabled, ShouldBeTrue)
			So(calendar.ServiceMessage, ShouldEqual, serviceMessage)
			So(calendar.EmergencyBanner.Title, ShouldEqual, emergencyBannerTitle)
			So(calendar.EmergencyBanner.Type, ShouldEqual, emergencyBannerType)
			So(calendar.EmergencyBanner.Description, ShouldEqual, emergencyBannerDescription)
			So(calendar.EmergencyBanner.URI, ShouldEqual, emergencyBannerURI)
			So(calendar.EmergencyBanner.LinkText, ShouldEqual, emergencyBannerLinkText)
			So(calendar.Metadata.Title, ShouldEqual, metaTitle)
			So(calendar.KeywordSearch.SearchTerm, ShouldEqual, params.Keywords)
			So(calendar.Sort, ShouldResemble, model.Sort{Mode: params.Sort.String(), Options: mapSortOptions(params)})
			So(calendar.BeforeDate, ShouldResemble, model.DateFieldset{
				Input: coreModel.InputDate{
					Language:        lang,
					Id:              "before-date",
					InputNameDay:    "before-day",
					InputNameMonth:  "before-month",
					InputNameYear:   "before-year",
					InputValueDay:   params.BeforeDate.DayString(),
					InputValueMonth: params.BeforeDate.MonthString(),
					InputValueYear:  params.BeforeDate.YearString(),
					Title: coreModel.Localisation{
						LocaleKey: "ReleasedBefore",
						Plural:    1,
					},
					Description: coreModel.Localisation{
						LocaleKey: "DateFilterDescription",
						Plural:    1,
					},
				}})
			So(calendar.AfterDate, ShouldResemble, model.DateFieldset{
				Input: coreModel.InputDate{
					Language:        lang,
					Id:              "after-date",
					InputNameDay:    "after-day",
					InputNameMonth:  "after-month",
					InputNameYear:   "after-year",
					InputValueDay:   params.AfterDate.DayString(),
					InputValueMonth: params.AfterDate.MonthString(),
					InputValueYear:  params.AfterDate.YearString(),
					Title: coreModel.Localisation{
						LocaleKey: "ReleasedAfter",
						Plural:    1,
					},
					Description: coreModel.Localisation{
						LocaleKey: "DateFilterDescription",
						Plural:    1,
					},
				}})
			So(calendar.ReleaseTypes, ShouldResemble, mapReleases(params, releaseResponse, lang))
			So(calendar.Pagination.TotalPages, ShouldEqual, 3)
			So(calendar.Pagination.CurrentPage, ShouldEqual, 1)
			So(calendar.Pagination.Limit, ShouldEqual, 5)
			So(calendar.TotalSearchPosition, ShouldEqual, 0)
			So(calendar.Entries.Count, ShouldEqual, 11)
			for i, r := range calendar.Entries.Items {
				So(r.URI, ShouldEqual, releaseResponse.Releases[i].URI)
				assertSiteSearchDateChanges(releaseResponse.Releases[i].DateChanges, r.DateChanges)
				So(r.Description.Title, ShouldEqual, releaseResponse.Releases[i].Description.Title)
				So(r.Description.Summary, ShouldEqual, releaseResponse.Releases[i].Description.Summary)
				So(r.Description.ReleaseDate, ShouldEqual, releaseResponse.Releases[i].Description.ReleaseDate)
				So(r.Description.Published, ShouldEqual, releaseResponse.Releases[i].Description.Published)
				So(r.Description.Finalised, ShouldEqual, releaseResponse.Releases[i].Description.Finalised)
				So(r.Description.Cancelled, ShouldEqual, releaseResponse.Releases[i].Description.Cancelled)
				So(r.Description.ProvisionalDate, ShouldEqual, releaseResponse.Releases[i].Description.ProvisionalDate)
				So(r.Description.Contact, ShouldBeZeroValue)
			}

			expectedStates := []model.PublicationState{
				{
					Type: "published",
				},
				{
					Type: "cancelled",
				},
				{
					Type:    "upcoming",
					SubType: "confirmed",
				},
				{
					Type:    "upcoming",
					SubType: "provisional",
				},
				{
					Type:    "upcoming",
					SubType: "provisional",
				},
			}
			for i, r := range calendar.Entries.Items {
				So(r.PublicationState, ShouldResemble, expectedStates[i])
			}
		})
	})
}

func TestReleaseCalendarErrorMapper(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a Release Calendar and a base page", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		Convey("CreateReleaseCalendarError maps correctly to a model Calendar object", func() {
			lang := "cy"
			metaTitle := "Calendr datganiadau"
			err := errors.New("test error message")
			errTitleKey := "ReleaseCalendarErrorTitleValidation"

			calendar := CreateReleaseCalendarError(basePage, lang, errTitleKey, err)

			So(calendar.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
			So(calendar.SiteDomain, ShouldEqual, basePage.SiteDomain)
			So(calendar.BetaBannerEnabled, ShouldBeTrue)
			So(calendar.Metadata.Title, ShouldEqual, metaTitle)
			So(calendar.GlobalError.Title.LocaleKey, ShouldEqual, errTitleKey)
			So(calendar.GlobalError.Message, ShouldEqual, err.Error())
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

func TestGetStartEndPage(t *testing.T) {
	Convey("Given a set of parameters expressing: the 'current page number', out of a 'total number of pages', and the 'window size'", t, func() {
		testcases := []struct{ current, total, window, exStart, exEnd int }{
			{current: 1, total: 1, window: 1, exStart: 1, exEnd: 1},

			{current: 1, total: 2, window: 1, exStart: 2, exEnd: 2},
			{current: 2, total: 2, window: 1, exStart: 1, exEnd: 1},

			{current: 1, total: 3, window: 2, exStart: 1, exEnd: 2},
			{current: 2, total: 3, window: 2, exStart: 2, exEnd: 3},
			{current: 3, total: 3, window: 2, exStart: 2, exEnd: 3},

			{current: 1, total: 3, window: 3, exStart: 1, exEnd: 3},
			{current: 2, total: 3, window: 3, exStart: 1, exEnd: 3},
			{current: 3, total: 3, window: 3, exStart: 1, exEnd: 3},

			{current: 3, total: 4, window: 3, exStart: 2, exEnd: 4},
			{current: 3, total: 4, window: 5, exStart: 1, exEnd: 4},

			{current: 28, total: 32, window: 5, exStart: 26, exEnd: 30},
			{current: 31, total: 32, window: 5, exStart: 28, exEnd: 32},
		}
		Convey("check the generated start and end page numbers are correct", func() {
			for _, tc := range testcases {
				sp, ep := getWindowStartEndPage(tc.current, tc.total, tc.window)
				So(sp, ShouldEqual, tc.exStart)
				So(ep, ShouldEqual, tc.exEnd)
			}
		})
	})
}

func TestTotalSearchPositions(t *testing.T) {
	Convey("Given a set of parameters expressing: the 'current page number', items per page ", t, func() {
		testcases := []struct{ currentPageNumber, itemsPerPage, exTotalSearchPosition int }{
			{currentPageNumber: 1, itemsPerPage: 10, exTotalSearchPosition: 0},
			{currentPageNumber: 2, itemsPerPage: 10, exTotalSearchPosition: 10},
			{currentPageNumber: 3, itemsPerPage: 10, exTotalSearchPosition: 20},
			{currentPageNumber: 4, itemsPerPage: 10, exTotalSearchPosition: 30},
			{currentPageNumber: 5, itemsPerPage: 10, exTotalSearchPosition: 40},
			{currentPageNumber: 6, itemsPerPage: 10, exTotalSearchPosition: 50},
			{currentPageNumber: 7, itemsPerPage: 10, exTotalSearchPosition: 60},
		}
		Convey("check the generated start and end page numbers are correct", func() {
			for _, tc := range testcases {
				sp := getTotalSearchPosition(tc.currentPageNumber, tc.itemsPerPage)
				So(sp, ShouldEqual, tc.exTotalSearchPosition)
			}
		})
	})
}

func TestGetPageURL(t *testing.T) {
	Convey("Given a set of Validated parameters, and a route path", t, func() {
		testcases := []struct {
			params   queryparams.ValidatedParams
			path     string
			expected string
		}{
			{
				params: queryparams.ValidatedParams{
					Limit:       10,
					Page:        2,
					AfterDate:   queryparams.MustParseDate("2021-11-30"),
					Keywords:    "test",
					Sort:        queryparams.TitleAZ,
					ReleaseType: queryparams.Published,
					Highlight:   true,
				},
				path:     "/test-prefix/releasecalendar",
				expected: "/test-prefix/releasecalendar?after-day=30&after-month=11&after-year=2021&highlight=true&keywords=test&limit=10&page=2&release-type=type-published&sort=alphabetical-az",
			},
			{
				params: queryparams.ValidatedParams{
					Limit:       25,
					Page:        5,
					BeforeDate:  queryparams.MustParseDate("2022-04-01"),
					Sort:        queryparams.RelDateDesc,
					ReleaseType: queryparams.Upcoming,
					Provisional: true,
					Postponed:   true,
					Census:      true,
				},
				path:     "/releasecalendar",
				expected: "/releasecalendar?before-day=1&before-month=4&before-year=2022&census=true&limit=25&page=5&release-type=type-upcoming&sort=date-newest&subtype-postponed=true&subtype-provisional=true",
			},
		}

		Convey("check the generated page url is correct", func() {
			for _, tc := range testcases {
				So(getPageURL(tc.params.Page, tc.params, tc.path), ShouldEqual, tc.expected)
			}
		})
	})
}

package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	sitesearch "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mocks"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	"github.com/ONSdigital/dp-renderer/helper"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	lang         = "en"
	accessToken  = "token"
	collectionID = "collection"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

func TestUnitHandlers(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := gomock.Any()

	Convey("test setStatusCode", t, func() {

		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:27700", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles internal server error", func() {
			req := httptest.NewRequest("GET", "http://localhost:27700", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test API", t, func() {
		mockRenderClient := NewMockRenderClient(mockCtrl)
		mockConfig, _ := config.Get()

		w := httptest.NewRecorder()
		router := mux.NewRouter()

		Convey("test Release", func() {
			mockApiClient := NewMockReleaseCalendarAPI(mockCtrl)
			url := "/releases/test"
			router.HandleFunc(url, Release(*mockConfig, mockRenderClient, mockApiClient))

			r := releasecalendar.Release{
				URI: url,
				Description: releasecalendar.ReleaseDescription{
					Title: "Test release",
				},
			}

			Convey("it returns 200 when rendered successfully", func() {
				mockApiClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, url).Return(&r, nil)
				mockRenderClient.EXPECT().NewBasePageModel()
				mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "release")

				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("it returns 200 when rendered successfully without headers or cookies", func() {
				mockApiClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, url).Return(&r, nil)
				mockRenderClient.EXPECT().NewBasePageModel()
				mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "release")

				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("it returns 500 when there is an error getting the release from the api", func() {
				mockApiClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, url).Return(nil, errors.New("error reading data"))
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("test ReleaseCalendar", func() {
			mockSearchClient := NewMockSearchAPI(mockCtrl)
			url := "/releasecalendar/test"
			router.HandleFunc(url, ReleaseCalendar(*mockConfig, mockRenderClient, mockSearchClient))
			r := sitesearch.ReleaseResponse{
				Releases: []sitesearch.Release{
					{URI: url,
						Description: sitesearch.ReleaseDescription{Title: "Release Calendar Entry Test"},
					},
				},
			}

			Convey("it returns 200 when rendered successfully", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, nil)
				mockRenderClient.EXPECT().NewBasePageModel()
				mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "calendar")

				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("it returns 200 when rendered successfully without headers or cookies", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, "", "", lang, defaultParams()).Return(r, nil)
				mockRenderClient.EXPECT().NewBasePageModel()
				mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "calendar")

				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("it returns 400 when there is an error in one of the parameters", func() {
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s?limit=-1", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})

			Convey("it returns 500 when there is an error getting the releases from the search api", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, errors.New("error reading data"))
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("test Calendar (calendar/releasecalendar)", func() {
			mockSearchClient := NewMockSearchAPI(mockCtrl)
			url := "/calendar/releasecalendar"
			router.HandleFunc(url, ReleaseCalendarICSEntries(*mockConfig, mockSearchClient))

			Convey("it returns 200 when an ICS file is generated successfully with a single calendar entry", func() {
				single := sitesearch.ReleaseResponse{
					Releases: []sitesearch.Release{
						{URI: url,
							Description: sitesearch.ReleaseDescription{
								Title:       "Release Calendar Entry Test 1",
								ReleaseDate: "2022-03-15T07:30:00Z",
							},
						},
					},
				}
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(single, nil)
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				Convey("and the ICS file payload is as expected", func() {
					payload := w.Body.Bytes()
					So(bytes.HasPrefix(payload, []byte(`BEGIN:VCALENDAR`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`Release Calendar Entry Test`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(url)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`20220315T073000`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`END:VCALENDAR`)), ShouldBeTrue)
				})

			})

			Convey("it returns 200 when an ICS file is generated successfully with multiple calendar entries", func() {
				multiple := sitesearch.ReleaseResponse{
					Releases: []sitesearch.Release{
						{URI: url + "test1",
							Description: sitesearch.ReleaseDescription{
								Title:       "Release Calendar Entry Test 1",
								ReleaseDate: "2022-03-15T07:30:00Z",
							},
						},
						{URI: url + "test2",
							Description: sitesearch.ReleaseDescription{
								Title:       "Release Calendar Entry Test 2",
								ReleaseDate: "2022-03-16T08:00:00Z",
							},
						},
					},
				}
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(multiple, nil)
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				Convey("and the ICS file payload is as expected", func() {
					payload := w.Body.Bytes()
					So(bytes.HasPrefix(payload, []byte(`BEGIN:VCALENDAR`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`Release Calendar Entry Test 1`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(url+"test1")), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`20220315T073000`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`Release Calendar Entry Test 2`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(url+"test2")), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`20220316T080000`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`END:VCALENDAR`)), ShouldBeTrue)
				})
			})

			Convey("it returns a well formed but empty ICS file when there are no upcoming releases", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(sitesearch.ReleaseResponse{}, nil)
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				payload := w.Body.Bytes()
				So(bytes.HasPrefix(payload, []byte(`BEGIN:VCALENDAR`)), ShouldBeTrue)
				So(bytes.Contains(payload, []byte(`END:VCALENDAR`)), ShouldBeTrue)
				So(len(payload), ShouldBeBetween, 100, 250)
			})

			Convey("it returns 500 when there is an error getting the releases from the search api", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(sitesearch.ReleaseResponse{}, errors.New("error reading data"))
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
				setRequestHeaders(req)

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

	})
}

func setRequestHeaders(req *http.Request) {
	headers.SetAuthToken(req, accessToken)
	headers.SetCollectionID(req, collectionID)
}

func defaultParams() url.Values {
	values := url.Values{}
	values.Set("limit", "10")
	values.Set("page", "1")
	values.Set("offset", "0")
	values.Set("fromDate", "")
	values.Set("toDate", "")
	values.Set("sort", queryparams.RelDateDesc.BackendString())
	values.Set("keywords", "")
	values.Set("query", "")
	values.Set("release-type", queryparams.Published.Name())
	values.Set("highlight", "true")

	return values
}

func defaultICSParams() url.Values {
	values := url.Values{}
	values.Set("limit", "1000")
	values.Set("toDate", time.Now().AddDate(0, 3, 0).Format(queryparams.DateFormat))
	values.Set("sort", queryparams.RelDateAsc.BackendString())
	values.Set("release-type", queryparams.Upcoming.String())

	return values
}

func TestICalDate(t *testing.T) {
	ds := []struct{ date, expected string }{
		{date: "1st Jan 2020", expected: ""},
		{date: "21-03-2021", expected: ""},
		{date: "2021-03-04T12:10:00", expected: ""},
		{date: "2021-03-04T12:10:00Z", expected: "20210304T121000Z"},
		{date: "2021-03-04T12:10:00.000Z", expected: "20210304T121000Z"},
		{date: "2021-03-04T12:10:00+05:00", expected: "20210304T071000Z"},
	}
	for _, tc := range ds {
		Convey("given a date string "+tc.date, t, func() {
			Convey("then the iCalDate returns the date formatted according to the iCal standard", func() {
				icd := iCalDate(context.Background(), tc.date)
				So(icd, ShouldEqual, tc.expected)
			})
		})
	}
}

type printer func(b []byte) (int, error)

func (p printer) Write(b []byte) (int, error) {
	return p(b)
}

func TestToICSFile(t *testing.T) {
	Convey("given a list of resources", t, func() {
		resources := []sitesearch.Release{{URI: "/release/stuff", Description: sitesearch.ReleaseDescription{Title: "A Release Title", ReleaseDate: "2021-03-04T12:10:00Z"}}}
		Convey("and a bad printer that fails", func() {
			printerError := errors.New("this is a bad-printer error")
			badPrinter := printer(func(b []byte) (int, error) { return 0, printerError })
			Convey("verify that the toICSFile function returns the error generated by the bad printer", func() {
				err := toICSFile(context.Background(), resources, badPrinter)
				So(err, ShouldEqual, printerError)
			})
		})

		Convey("and a good printer that does not fail", func() {
			goodPrinter := new(bytes.Buffer)
			Convey("verify that the toICSFile function correctly prints the ICS file for the given releases", func() {
				err := toICSFile(context.Background(), resources, goodPrinter)
				So(err, ShouldBeNil)
				So(goodPrinter.Bytes(), ShouldNotBeNil)
			})
		})
	})
}

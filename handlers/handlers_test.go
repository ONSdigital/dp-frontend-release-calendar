package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	sitesearch "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mocks"
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

	Convey("test API (Release and ReleaseCalendar)", t, func() {
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
	values.Set("sort", "release_date_desc")
	values.Set("keywords", "")
	values.Set("query", "")
	values.Set("release-type", "type-published")
	values.Set("highlight", "true")

	return values
}

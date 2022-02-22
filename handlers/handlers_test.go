package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	gomock "github.com/golang/mock/gomock"
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

	Convey("test Release", t, func() {
		url := "/releases/test"
		r := releasecalendar.Release{
			URI: url,
			Description: releasecalendar.ReleaseDescription{
				Title: "Test release",
			},
		}

		mockApiClient := NewMockReleaseCalendarAPI(mockCtrl)
		mockRenderClient := NewMockRenderClient(mockCtrl)
		mockConfig := config.Config{}

		router := mux.NewRouter()
		router.HandleFunc(url, Release(mockConfig, mockRenderClient, mockApiClient))

		w := httptest.NewRecorder()

		Convey("it returns 200 when rendered succesfully", func() {
			mockApiClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, url).Return(&r, nil)
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "release")

			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 200 when rendered succesfully without headers or cookies", func() {
			mockApiClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, url).Return(&r, nil)
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "release")

			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 500 when there is an error getting the release from the api", func() {
			mockApiClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, url).Return(nil, errors.New(("error reading data")))
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

}

func setRequestHeaders(req *http.Request) {
	headers.SetAuthToken(req, accessToken)
	headers.SetCollectionID(req, collectionID)
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	sunsetDate         = "2025-08-29"
	sunsetLink         = "https://www.ons.gov.uk"
	deprecationDate    = "2025-08-29T10:00:00Z"
	deprecationMessage = "The releases data endpoint is deprecated"
)

func TestEndpointDeprecation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := gomock.Any()

	Convey("Test deprecation", t, func() {
		mockConfig, _ := config.Get()
		mockAPIClient := NewMockReleaseCalendarAPI(mockCtrl)
		root := "/releases"
		w := httptest.NewRecorder()
		router := mux.NewRouter()
		r := releasecalendar.Release{
			Description: releasecalendar.ReleaseDescription{
				Title: "Test release",
			},
		}

		titleSegment := strings.ReplaceAll(strings.ToLower(r.Description.Title), " ", "")
		r.URI = fmt.Sprintf("%s/%s", root, titleSegment)

		Convey("Test '/releases/{release-title}/data' endpoint deprecation", func() {
			Convey("When deprecation is enabled", func() {
				mockConfig.Deprecation.DeprecateEndpoint = true
				mockConfig.Deprecation.Link = sunsetLink
				mockConfig.Deprecation.Deprecation = deprecationDate
				mockConfig.Deprecation.DeprecationMessage = deprecationMessage

				Convey("And the sunset date has passed", func() {
					mockConfig.Deprecation.Sunset = sunsetDate

					Convey("And the release is retrieved successfully", func() {
						req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s/data", root, titleSegment), http.NoBody)

						router.HandleFunc(root+"/{release-title}/data", ReleaseData(*mockConfig, mockAPIClient))

						router.ServeHTTP(w, req)

						parsedSunset, _ := parseTime(sunsetDate)
						parsedDeprecation, _ := parseTime(deprecationDate)

						Convey("Then it returns 404", func() {
							So(w.Code, ShouldEqual, http.StatusNotFound)
							So(w.Header().Get("content-type"), ShouldEqual, "application/json")
							So(w.Header().Get("deprecation"), ShouldEqual, fmt.Sprintf("@%d", parsedDeprecation.Unix()))
							So(w.Header().Get("sunset"), ShouldEqual, parsedSunset.Format(time.RFC1123))
							So(w.Header().Get("link"), ShouldEqual, fmt.Sprintf("<%s>; rel=\"sunset\"", sunsetLink))
						})
					})
				})

				Convey("And the sunset date has not passed", func() {
					mockConfig.Deprecation.Sunset = "2026-08-29"

					Convey("And the release is retrieved successfully", func() {
						req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s/data", root, titleSegment), http.NoBody)

						router.HandleFunc(root+"/{release-title}/data", ReleaseData(*mockConfig, mockAPIClient))

						router.ServeHTTP(w, req)

						Convey("Then it returns 200", func() {
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})
				})
			})

			Convey("When deprecation is disabled", func() {
				mockConfig.Deprecation.DeprecateEndpoint = false

				js, _ := json.Marshal(r)
				Convey("And the release is retrieved successfully", func() {
					mockAPIClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, r.URI).Return(&r, nil)

					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s/data", root, titleSegment), http.NoBody)
					if err := setRequestHeaders(req); err != nil {
						t.Fatalf("unable to set request headers, error: %v", err)
					}

					router.HandleFunc(root+"/{release-title}/data", ReleaseData(*mockConfig, mockAPIClient))

					router.ServeHTTP(w, req)

					Convey("Then it returns 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
						So(w.Header().Get("content-type"), ShouldEqual, "application/json")
						So(w.Body.Bytes(), ShouldResemble, js)
					})
				})
			})
		})
	})
}

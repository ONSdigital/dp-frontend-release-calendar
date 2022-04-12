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
		mockBabbageAPI := NewMockBabbageAPI(mockCtrl)
		mockConfig, _ := config.Get()

		w := httptest.NewRecorder()
		router := mux.NewRouter()

		Convey("test Release", func() {
			mockApiClient := NewMockReleaseCalendarAPI(mockCtrl)
			url := "/releases/test"
			maxAge := 670
			router.HandleFunc(url, Release(*mockConfig, mockRenderClient, mockApiClient, mockBabbageAPI))

			r := releasecalendar.Release{
				URI: url,
				Description: releasecalendar.ReleaseDescription{
					Title: "Test release",
				},
			}

			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)

			Convey("When there is an error getting the release from the release calendar API", func() {
				apiError := errors.New("error reading data")
				Convey("And the request uses headers", func() {
					setRequestHeaders(req)
					mockApiClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, url).Return(nil, apiError)

					Convey("Then it returns 500", func() {
						router.ServeHTTP(w, req)

						So(w.Code, ShouldEqual, http.StatusInternalServerError)
					})
				})

				Convey("And the request does not use headers", func() {
					mockApiClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, url).Return(&r, nil).Return(nil, apiError)

					Convey("Then it returns 500", func() {
						router.ServeHTTP(w, req)

						So(w.Code, ShouldEqual, http.StatusInternalServerError)
					})
				})
			})

			Convey("When there is no problem with the request", func() {
				mockRenderClient.EXPECT().NewBasePageModel()
				mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "release")

				Convey("And the request uses headers", func() {
					setRequestHeaders(req)
					mockApiClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, url).Return(&r, nil)

					Convey("And Babbage calculates the cache max age successfully", func() {
						mockBabbageAPI.EXPECT().GetMaxAge(ctx, url, mockConfig.MaxAgeKey).Return(maxAge, nil)
						expectedCacheControlHeader := fmt.Sprintf("public, max-age=%d", maxAge)

						Convey("Then it returns 200 and the right cache header", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusOK)
							So(w.Header().Get("Cache-Control"), ShouldEqual, expectedCacheControlHeader)
						})
					})

					Convey("And there is an error calling Babbage", func() {
						mockBabbageAPI.EXPECT().GetMaxAge(ctx, url, mockConfig.MaxAgeKey).Return(maxAge, errors.New("Error on Babbage"))

						Convey("Then it returns 200 and the default cache header", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusOK)
							So(w.Header().Get("Cache-Control"), ShouldEqual, "public, max-age=0")
						})
					})

				})

				Convey("And the request does not use headers", func() {
					mockApiClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, url).Return(&r, nil)

					Convey("And Babbage calculates the cache max age successfully", func() {
						mockBabbageAPI.EXPECT().GetMaxAge(ctx, url, mockConfig.MaxAgeKey).Return(maxAge, nil)
						expectedCacheControlHeader := fmt.Sprintf("public, max-age=%d", maxAge)

						Convey("Then it returns 200 and the right cache header", func() {
							router.ServeHTTP(w, req)
							So(w.Code, ShouldEqual, http.StatusOK)
							So(w.Header().Get("Cache-Control"), ShouldEqual, expectedCacheControlHeader)
						})
					})

					Convey("And there is an error calling Babbage", func() {
						mockBabbageAPI.EXPECT().GetMaxAge(ctx, url, mockConfig.MaxAgeKey).Return(maxAge, errors.New("Error on Babbage"))

						Convey("Then it returns 200 and the default cache header", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusOK)
							So(w.Header().Get("Cache-Control"), ShouldEqual, "public, max-age=0")
						})
					})
				})
			})
		})

		Convey("test ReleaseCalendar", func() {
			mockSearchClient := NewMockSearchAPI(mockCtrl)
			url := "/releasecalendar"
			maxAge := 422
			router.HandleFunc(url, ReleaseCalendar(*mockConfig, mockRenderClient, mockSearchClient, mockBabbageAPI))

			r := sitesearch.ReleaseResponse{
				Releases: []sitesearch.Release{
					{
						URI:         url,
						Description: sitesearch.ReleaseDescription{Title: "Release Calendar Entry Test"},
					},
				},
			}

			Convey("Given a request without parameters", func() {
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", url), nil)

				Convey("When there is an error getting the releases from the search API", func() {
					apiError := errors.New("error reading data")
					Convey("And the request uses headers", func() {
						setRequestHeaders(req)
						mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, apiError)

						Convey("Then it returns 500", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusInternalServerError)
						})
					})

					Convey("And the request does not use headers", func() {
						mockSearchClient.EXPECT().GetReleases(ctx, "", "", lang, defaultParams()).Return(r, apiError)

						Convey("Then it returns 500", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusInternalServerError)
						})
					})
				})

				Convey("When there is no problem with the request", func() {
					mockRenderClient.EXPECT().NewBasePageModel()
					mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "calendar")

					Convey("And the request uses headers", func() {
						setRequestHeaders(req)
						mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, nil)

						Convey("And Babbage calculates the cache max age successfully", func() {
							mockBabbageAPI.EXPECT().GetMaxAge(ctx, "/releasecalendar", mockConfig.MaxAgeKey).Return(maxAge, nil)
							expectedCacheControlHeader := fmt.Sprintf("public, max-age=%d", maxAge)

							Convey("Then it returns 200 and the right cache header", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusOK)
								So(w.Header().Get("Cache-Control"), ShouldEqual, expectedCacheControlHeader)
							})
						})

						Convey("And there is an error calling Babbage", func() {
							mockBabbageAPI.EXPECT().GetMaxAge(ctx, "/releasecalendar", mockConfig.MaxAgeKey).Return(maxAge, errors.New("Error on Babbage"))

							Convey("Then it returns 200 and the default cache header", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusOK)
								So(w.Header().Get("Cache-Control"), ShouldEqual, "public, max-age=0")
							})
						})
					})

					Convey("And the request does not use headers", func() {
						mockSearchClient.EXPECT().GetReleases(ctx, "", "", lang, defaultParams()).Return(r, nil)

						Convey("And Babbage calculates the cache max age successfully", func() {
							mockBabbageAPI.EXPECT().GetMaxAge(ctx, "/releasecalendar", mockConfig.MaxAgeKey).Return(maxAge, nil)
							expectedCacheControlHeader := fmt.Sprintf("public, max-age=%d", maxAge)

							Convey("Then it returns 200 and the right cache header", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusOK)
								So(w.Header().Get("Cache-Control"), ShouldEqual, expectedCacheControlHeader)
							})
						})

						Convey("And there is an error calling Babbage", func() {
							mockBabbageAPI.EXPECT().GetMaxAge(ctx, "/releasecalendar", mockConfig.MaxAgeKey).Return(maxAge, errors.New("Error on Babbage"))

							Convey("Then it returns 200 and the default cache header", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusOK)
								So(w.Header().Get("Cache-Control"), ShouldEqual, "public, max-age=0")
							})
						})
					})
				})
			})

			Convey("Given a request with parameters", func() {

				Convey("When the limit parameter is negative", func() {
					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s?limit=-1", url), nil)

					Convey("Then it returns 400", func() {
						router.ServeHTTP(w, req)

						So(w.Code, ShouldEqual, http.StatusBadRequest)
					})
				})

				//TODO: Add test cases for parameter validation
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
	values.Set("release-type", "type-cancelled")
	values.Set("highlight", "true")

	return values
}

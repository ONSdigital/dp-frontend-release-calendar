package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	sitesearch "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mocks"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	"github.com/ONSdigital/dp-renderer/v2/helper"

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

// TODO TestUnitHandlers needs refactoring - it is way too complex and therefore hard to maintain.
// This needs to be split into different test functions.
// Overlapping with other unit tests from over packages should be removed
//
//nolint:gocognit // needs refactoring see comment above
func TestUnitHandlers(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := gomock.Any()

	Convey("test setStatusCode", t, func() {
		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:27700", http.NoBody)
			w := httptest.NewRecorder()
			err := &testCliError{}

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles internal server error", func() {
			req := httptest.NewRequest("GET", "http://localhost:27700", http.NoBody)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test handler", t, func() {
		mockRenderClient := NewMockRenderClient(mockCtrl)
		mockConfig, _ := config.Get()

		w := httptest.NewRecorder()
		router := mux.NewRouter()

		Convey("test Release endpoints", func() {
			mockZebedeeClient := NewMockZebedeeClient(mockCtrl)
			mockAPIClient := NewMockReleaseCalendarAPI(mockCtrl)
			root := "/releases"
			r := releasecalendar.Release{
				Description: releasecalendar.ReleaseDescription{
					Title: "Test release",
				},
			}
			titleSegment := strings.ReplaceAll(strings.ToLower(r.Description.Title), " ", "")
			r.URI = fmt.Sprintf("%s/%s", root, titleSegment)

			Convey("test '/releases/{release-title}'", func() {
				router.HandleFunc(root+"/{release-title}", Release(*mockConfig, mockRenderClient, mockAPIClient, mockZebedeeClient))

				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s", root, titleSegment), http.NoBody)
				Convey("When there is an error getting the release from the release calendar API", func() {
					apiError := errors.New("error reading data")
					Convey("And the request uses headers", func() {
						if err := setRequestHeaders(req); err != nil {
							t.Fatalf("unable to set request headers, error: %v", err)
						}
						mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
						mockAPIClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, r.URI).Return(nil, apiError)

						Convey("Then it returns 500", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusInternalServerError)
						})
					})

					Convey("And the request does not use headers", func() {
						mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
						mockAPIClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, r.URI).Return(&r, nil).Return(nil, apiError)

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
						if err := setRequestHeaders(req); err != nil {
							t.Fatalf("unable to set request headers, error: %v", err)
						}
						mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
						mockAPIClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, r.URI).Return(&r, nil)

						Convey("Then it returns 200", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})

					Convey("And the request does not use headers", func() {
						mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
						mockAPIClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, r.URI).Return(&r, nil)

						Convey("Then it returns 200", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})

					Convey("And the request uses eTags", func() {
						mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
						mockAPIClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, r.URI).Return(&r, nil)

						Convey("Then the eTag header is set", func() {
							router.ServeHTTP(w, req)
							eTag := w.Header().Get("ETag")

							So(eTag, ShouldNotBeEmpty)
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})
				})

				Convey("When the response includes a migrationLink", func() {
					redirect := "/redirect1"
					releaseWithMigrationLink := releasecalendar.Release{
						Description: releasecalendar.ReleaseDescription{
							Title:         "Test release",
							MigrationLink: redirect,
						},
						URI: "/releases/myrelease",
					}

					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", releaseWithMigrationLink.URI), http.NoBody)

					mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
					mockAPIClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, releaseWithMigrationLink.URI).Return(&releaseWithMigrationLink, nil)

					Convey("Then it returns 308", func() {
						router.ServeHTTP(w, req)
						location := w.Result().Header.Get("Location")
						So(w.Code, ShouldEqual, http.StatusPermanentRedirect)
						So(location, ShouldEqual, redirect)
					})
				})
			})

			Convey("test '/releases/{release-title}/data' endpoint", func() {
				dataSegment := "data"
				router.HandleFunc(root+"/{release-title}/"+dataSegment, ReleaseData(*mockConfig, mockAPIClient))

				js, _ := json.Marshal(r)
				Convey("when the release is retrieved successfully", func() {
					mockAPIClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, r.URI).Return(&r, nil)

					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s/%s", root, titleSegment, dataSegment), http.NoBody)
					if err := setRequestHeaders(req); err != nil {
						t.Fatalf("unable to set request headers, error: %v", err)
					}

					router.ServeHTTP(w, req)

					Convey("it returns 200 with the expected json payload ", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
						So(w.Body.Bytes(), ShouldResemble, js)
					})
					Convey("and the content type is 'application/json' ", func() {
						So(w.Header().Get("content-type"), ShouldEqual, "application/json")
					})
				})

				Convey("when the release is retrieved successfully without headers or cookies", func() {
					mockAPIClient.EXPECT().GetLegacyRelease(ctx, "", "", lang, r.URI).Return(&r, nil)

					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s/%s", root, titleSegment, dataSegment), http.NoBody)

					router.ServeHTTP(w, req)

					Convey("it returns 200 with the expected json payload ", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
						So(w.Body.Bytes(), ShouldResemble, js)
					})
					Convey("and the content type is 'application/json' ", func() {
						So(w.Header().Get("content-type"), ShouldEqual, "application/json")
					})
				})

				Convey("it returns 500 when there is an error getting the release from the api", func() {
					mockAPIClient.EXPECT().GetLegacyRelease(ctx, accessToken, collectionID, lang, r.URI).Return(nil, errors.New("error reading data"))
					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s/%s/%s", root, titleSegment, dataSegment), http.NoBody)
					if err := setRequestHeaders(req); err != nil {
						t.Fatalf("unable to set request headers, error: %v", err)
					}
					router.ServeHTTP(w, req)

					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("test ReleaseCalendar endpoints", func() {
			mockSearchClient := NewMockSearchAPI(mockCtrl)
			mockZebedeeClient := NewMockZebedeeClient(mockCtrl)

			Convey("test '/releasecalendar' endpoint", func() {
				endpoint := "/releasecalendar"
				router.HandleFunc(endpoint, ReleaseCalendar(*mockConfig, mockRenderClient, mockSearchClient, mockZebedeeClient))
				r := sitesearch.ReleaseResponse{
					Releases: []sitesearch.Release{
						{
							URI:         "/releases/releasecalendarentrytest",
							Description: sitesearch.ReleaseDescription{Title: "Release Calendar Entry Test"},
						},
					},
				}

				Convey("Given a request without parameters", func() {
					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)

					Convey("When there is an error getting the releases from the search API", func() {
						apiError := errors.New("error reading data")

						Convey("And the request uses headers", func() {
							if err := setRequestHeaders(req); err != nil {
								t.Fatalf("unable to set request headers, error: %v", err)
							}
							mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
							mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, apiError)

							Convey("Then it returns 500", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusInternalServerError)
							})
						})

						Convey("And the request does not use headers", func() {
							mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
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
							if err := setRequestHeaders(req); err != nil {
								t.Fatalf("unable to set request headers, error: %v", err)
							}
							mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
							mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, nil)

							Convey("Then it returns 200", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusOK)
							})
						})

						Convey("And the request does not use headers", func() {
							mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
							mockSearchClient.EXPECT().GetReleases(ctx, "", "", lang, defaultParams()).Return(r, nil)

							Convey("Then it returns 200", func() {
								router.ServeHTTP(w, req)

								So(w.Code, ShouldEqual, http.StatusOK)
							})
						})
					})
				})

				Convey("Given a request with parameters", func() {
					Convey("When parameters are invalid", func() {
						mockRenderClient.EXPECT().NewBasePageModel()
						mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "calendar")
						mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")

						req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s?limit=-1&release-type=type-sf&after-year=dad&before-day=44&after-month-99&sort=date-blah&page=dogs", endpoint), http.NoBody)

						Convey("Then it returns 200", func() {
							router.ServeHTTP(w, req)

							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})
				})
			})

			Convey("test '/releasecalendar/data'", func() {
				endpoint := "/releasecalendar/data"
				router.HandleFunc(endpoint, ReleaseCalendarData(*mockConfig, mockSearchClient))
				r := sitesearch.ReleaseResponse{
					Releases: []sitesearch.Release{
						{
							URI:         "/releases/releasecalendarentrytest",
							Description: sitesearch.ReleaseDescription{Title: "Release Calendar Entry Test"},
						},
					},
				}

				js, _ := json.Marshal(r)
				Convey("when the release calendar entries are retrieved successfully", func() {
					mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, nil)

					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)
					if err := setRequestHeaders(req); err != nil {
						t.Fatalf("unable to set request headers, error: %v", err)
					}

					router.ServeHTTP(w, req)

					Convey("it returns 200 with the expected json payload ", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
						So(w.Body.Bytes(), ShouldResemble, js)
					})
					Convey("and the content type is 'application/json' ", func() {
						So(w.Header().Get("content-type"), ShouldEqual, "application/json")
					})
				})

				Convey("when the release calendar entries are retrieved successfully without headers or cookies", func() {
					mockSearchClient.EXPECT().GetReleases(ctx, "", "", lang, defaultParams()).Return(r, nil)

					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)

					router.ServeHTTP(w, req)

					Convey("it returns 200 with the expected json payload ", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
						So(w.Body.Bytes(), ShouldResemble, js)
					})
					Convey("and the content type is 'application/json' ", func() {
						So(w.Header().Get("content-type"), ShouldEqual, "application/json")
					})
				})

				Convey("it returns 400 when there is an error in one of the parameters", func() {
					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s?limit=-1", endpoint), http.NoBody)
					if err := setRequestHeaders(req); err != nil {
						t.Fatalf("unable to set request headers, error: %v", err)
					}

					router.ServeHTTP(w, req)

					So(w.Code, ShouldEqual, http.StatusBadRequest)
				})

				Convey("it returns 500 when there is an error getting the releases from the search api", func() {
					mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultParams()).Return(r, errors.New("error reading data"))
					req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)
					if err := setRequestHeaders(req); err != nil {
						t.Fatalf("unable to set request headers, error: %v", err)
					}

					router.ServeHTTP(w, req)

					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("test calendar/releasecalendar endpoint", func() {
			mockSearchClient := NewMockSearchAPI(mockCtrl)
			endpoint := "/calendar/releasecalendar"
			router.HandleFunc(endpoint, ReleaseCalendarICSEntries(*mockConfig, mockSearchClient))

			Convey("it returns 200 when an ICS file is generated successfully with a single calendar entry", func() {
				single := sitesearch.ReleaseResponse{
					Releases: []sitesearch.Release{
						{URI: "/releases/releasecalendarentrytest1",
							Description: sitesearch.ReleaseDescription{
								Title:       "Release Calendar Entry Test 1",
								ReleaseDate: "2022-03-15T07:30:00Z",
							},
						},
					},
				}
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(single, nil)
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)
				if err := setRequestHeaders(req); err != nil {
					t.Fatalf("unable to set request headers, error: %v", err)
				}

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				Convey("and the ICS file payload is as expected", func() {
					payload := w.Body.Bytes()
					So(bytes.HasPrefix(payload, []byte(`BEGIN:VCALENDAR`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`Release Calendar Entry Test 1`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`/releases/releasecalendarentrytest1`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`20220315T073000`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`END:VCALENDAR`)), ShouldBeTrue)
				})
			})

			Convey("it returns 200 when an ICS file is generated successfully with multiple calendar entries", func() {
				multiple := sitesearch.ReleaseResponse{
					Releases: []sitesearch.Release{
						{URI: "/releases/releasecalendarentrytest1",
							Description: sitesearch.ReleaseDescription{
								Title:       "Release Calendar Entry Test 1",
								ReleaseDate: "2022-03-15T07:30:00Z",
							},
						},
						{URI: "/releases/releasecalendarentrytest2",
							Description: sitesearch.ReleaseDescription{
								Title:       "Release Calendar Entry Test 2",
								ReleaseDate: "2022-03-16T08:00:00Z",
							},
						},
					},
				}
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(multiple, nil)
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)
				if err := setRequestHeaders(req); err != nil {
					t.Fatalf("unable to set request headers, error: %v", err)
				}

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				Convey("and the ICS file payload is as expected", func() {
					payload := w.Body.Bytes()
					So(bytes.HasPrefix(payload, []byte(`BEGIN:VCALENDAR`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`Release Calendar Entry Test 1`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`/releases/releasecalendarentrytest1`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`20220315T073000`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`Release Calendar Entry Test 2`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`/releases/releasecalendarentrytest2`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`20220316T080000`)), ShouldBeTrue)
					So(bytes.Contains(payload, []byte(`END:VCALENDAR`)), ShouldBeTrue)
				})
			})

			Convey("it returns a well formed but empty ICS file when there are no upcoming releases", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(sitesearch.ReleaseResponse{}, nil)
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)
				if err := setRequestHeaders(req); err != nil {
					t.Fatalf("unable to set request headers, error: %v", err)
				}

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				payload := w.Body.Bytes()
				So(bytes.HasPrefix(payload, []byte(`BEGIN:VCALENDAR`)), ShouldBeTrue)
				So(bytes.Contains(payload, []byte(`END:VCALENDAR`)), ShouldBeTrue)
				So(len(payload), ShouldBeBetween, 100, 250)
			})

			Convey("it returns 500 when there is an error getting the releases from the search api", func() {
				mockSearchClient.EXPECT().GetReleases(ctx, accessToken, collectionID, lang, defaultICSParams()).Return(sitesearch.ReleaseResponse{}, errors.New("error reading data"))
				req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:27700%s", endpoint), http.NoBody)
				if err := setRequestHeaders(req); err != nil {
					t.Fatalf("unable to set request headers, error: %v", err)
				}

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func setRequestHeaders(req *http.Request) error {
	if err := headers.SetAuthToken(req, accessToken); err != nil {
		return err
	}

	if err := headers.SetCollectionID(req, collectionID); err != nil {
		return err
	}

	return nil
}

func defaultParams() url.Values {
	values := url.Values{}
	values.Set("limit", "10")
	values.Set("page", "1")
	values.Set("offset", "0")
	values.Set("sort", queryparams.RelDateDesc.BackendString())
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

func TestCreateRSSFeed(t *testing.T) {
	Convey("test createRSSFeed", t, func() {
		lang := "en"
		collectionID := "collection"
		accessToken := "token"
		validatedParams := queryparams.ValidatedParams{}

		mockCtrl := gomock.NewController(t)
		mockSearchClient := NewMockSearchAPI(mockCtrl)

		Convey("when GetReleases returns success", func() {
			mockSearchClient.EXPECT().GetReleases(
				gomock.Any(),
				accessToken,
				collectionID,
				lang,
				gomock.Any(),
			).Return(sitesearch.ReleaseResponse{}, nil)

			req := httptest.NewRequest("GET", "http://localhost:27700", http.NoBody)
			w := httptest.NewRecorder()

			err := createRSSFeed(context.Background(), w, req, lang, collectionID, accessToken, mockSearchClient, validatedParams)

			Convey("it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("it should set the Content-Type header to 'application/rss+xml'", func() {
				contentType := w.Header().Get("Content-Type")
				So(contentType, ShouldEqual, "application/rss+xml")
			})
		})

		Convey("when GetReleases returns an error", func() {
			mockSearchClient.EXPECT().GetReleases(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sitesearch.ReleaseResponse{}, errors.New("mocked error"))

			req := httptest.NewRequest("GET", "http://localhost:27700", http.NoBody)
			w := httptest.NewRecorder()

			err := createRSSFeed(context.Background(), w, req, lang, collectionID, accessToken, mockSearchClient, validatedParams)

			Convey("it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when parsing release date fails", func() {
			mockSearchClient.EXPECT().GetReleases(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sitesearch.ReleaseResponse{
				Releases: []sitesearch.Release{{Description: sitesearch.ReleaseDescription{ReleaseDate: "invalid date"}}},
			}, nil)

			req := httptest.NewRequest("GET", "http://localhost:27700", http.NoBody)
			w := httptest.NewRecorder()

			err := createRSSFeed(context.Background(), w, req, lang, collectionID, accessToken, mockSearchClient, validatedParams)

			Convey("it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

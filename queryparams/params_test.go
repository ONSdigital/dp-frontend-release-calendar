package queryparams

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-renderer/v2/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetIntValidator(t *testing.T) {
	Convey("given an intValidator parameterised with a maximum and minimum value", t, func() {
		validator := getIntValidator(0, 1000)

		Convey("and a set of int values as strings", func() {
			limits := []struct {
				value   string
				exValue int
				exError error
			}{
				{value: "XXX", exValue: 0, exError: errors.New("value contains non numeric characters")},
				{value: "-1", exValue: 0, exError: errors.New("value is below the minimum value (0)")},
				{value: "1001", exValue: 0, exError: fmt.Errorf("value is above the maximum value (1000)")},
				{value: "0", exValue: 0, exError: nil},
				{value: "123", exValue: 123, exError: nil},
				{value: "1000", exValue: 1000, exError: nil},
			}

			Convey("check that the validator correctly validates the limit, giving the expected results", func() {
				for _, ls := range limits {
					v, e := validator(ls.value)

					So(v, ShouldEqual, ls.exValue)
					So(e, ShouldResemble, ls.exError)
				}
			})
		})
	})
}

func TestGetLimit(t *testing.T) {
	Convey("Given a list of params", t, func() {
		ctx := context.Background()
		params := make(url.Values)
		defaultValue := 8
		maxValue := 55
		Convey("And it does not include a limit param", func() {
			Convey("When we call GetLimit", func() {
				res, err := GetLimit(ctx, params, defaultValue, maxValue)
				Convey("Then the default value is returned", func() {
					So(err, ShouldBeNil)
					So(res, ShouldEqual, defaultValue)
				})
			})
		})
		Convey("And it includes a limit param", func() {
			Convey("And it is empty", func() {
				params.Set("limit", "")
				Convey("When we call GetLimit", func() {
					res, err := GetLimit(ctx, params, defaultValue, maxValue)
					Convey("Then the default value is returned", func() {
						So(err, ShouldBeNil)
						So(res, ShouldEqual, defaultValue)
					})
				})
			})
			Convey("And it is valid", func() {
				limit := 0
				params.Set("limit", strconv.Itoa(limit))
				Convey("When we call GetLimit", func() {
					res, err := GetLimit(ctx, params, defaultValue, maxValue)
					Convey("Then the value is returned", func() {
						So(err, ShouldBeNil)
						So(res, ShouldEqual, limit)
					})
				})
			})
			Convey("And it is lower than 0", func() {
				params.Set("limit", "-1")
				Convey("When we call GetLimit", func() {
					_, err := GetLimit(ctx, params, defaultValue, maxValue)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid limit parameter: value is below the minimum value (0)")
					})
				})
			})
			Convey("And it is higher than the maximum", func() {
				limit := maxValue + 1
				params.Set("limit", strconv.Itoa(limit))
				Convey("When we call GetLimit", func() {
					_, err := GetLimit(ctx, params, defaultValue, maxValue)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid limit parameter: value is above the maximum value (55)")
					})
				})
			})
			Convey("And it is not a number", func() {
				params.Set("limit", "seven")
				Convey("When we call GetLimit", func() {
					_, err := GetLimit(ctx, params, defaultValue, maxValue)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid limit parameter: value contains non numeric characters")
					})
				})
			})
		})
	})
}

func TestGetPage(t *testing.T) {
	Convey("Given a list of params", t, func() {
		ctx := context.Background()
		params := make(url.Values)
		maxPage := 10
		Convey("And it does not include a page param", func() {
			Convey("When we call GetPage", func() {
				res, err := GetPage(ctx, params, maxPage)
				Convey("Then the default value is returned", func() {
					So(err, ShouldBeNil)
					So(res, ShouldEqual, 1)
				})
			})
		})
		Convey("And it includes a page param", func() {
			Convey("And it is empty", func() {
				params.Set("page", "")
				Convey("When we call GetPage", func() {
					res, err := GetPage(ctx, params, maxPage)
					Convey("Then the default value is returned", func() {
						So(err, ShouldBeNil)
						So(res, ShouldEqual, 1)
					})
				})
			})
			Convey("And it is valid", func() {
				limit := 1
				params.Set("page", strconv.Itoa(limit))
				Convey("When we call GetPage", func() {
					res, err := GetPage(ctx, params, maxPage)
					Convey("Then the value is returned", func() {
						So(err, ShouldBeNil)
						So(res, ShouldEqual, limit)
					})
				})
			})
			Convey("And it is lower than 1", func() {
				params.Set("page", "0")
				Convey("When we call GetPage", func() {
					_, err := GetPage(ctx, params, maxPage)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid page parameter: value is below the minimum value (1)")
					})
				})
			})
			Convey("And it is higher than the maximum", func() {
				page := maxPage + 1
				params.Set("page", strconv.Itoa(page))
				Convey("When we call GetPage", func() {
					_, err := GetPage(ctx, params, maxPage)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid page parameter: value is above the maximum value (10)")
					})
				})
			})
			Convey("And it is not a number", func() {
				params.Set("page", "three")
				Convey("When we call GetPage", func() {
					_, err := GetPage(ctx, params, maxPage)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid page parameter: value contains non numeric characters")
					})
				})
			})
		})
	})
}

func TestCalculateOffset(t *testing.T) {
	Convey("given a range of page numbers and sizes", t, func() {
		testcases := []struct{ pageNumber, pageSize, expectedOffset int }{
			{pageNumber: 0, pageSize: 10, expectedOffset: 0},
			{pageNumber: 1, pageSize: 0, expectedOffset: 0},
			{pageNumber: 1, pageSize: 10, expectedOffset: 0},
			{pageNumber: 2, pageSize: 10, expectedOffset: 10},
			{pageNumber: 3, pageSize: 10, expectedOffset: 20},
		}

		Convey("check that the calculated offset is correct", func() {
			for _, tc := range testcases {
				So(CalculateOffset(tc.pageNumber, tc.pageSize), ShouldEqual, tc.expectedOffset)
			}
		})
	})
}

func TestCalculatePageNumber(t *testing.T) {
	Convey("given a range of item offsets (0 based) and page sizes (1 based)", t, func() {
		testcases := []struct{ offset, pageSize, expectedPage int }{
			{offset: 0, pageSize: 0, expectedPage: 1},
			{offset: 1, pageSize: 0, expectedPage: 1},
			{offset: 0, pageSize: 1, expectedPage: 1},
			{offset: 1, pageSize: 1, expectedPage: 2},
			{offset: 2, pageSize: 1, expectedPage: 3},
			{offset: 9, pageSize: 10, expectedPage: 1},
			{offset: 10, pageSize: 10, expectedPage: 2},
			{offset: 10, pageSize: 5, expectedPage: 3},
		}

		Convey("check that the calculated page number is correct", func() {
			for _, tc := range testcases {
				So(CalculatePageNumber(tc.offset, tc.pageSize), ShouldEqual, tc.expectedPage)
			}
		})
	})
}

func TestGetSortOrder(t *testing.T) {
	Convey("given a set of erroneous sort string options", t, func() {
		badSortOptions := []string{"dont sort", "sort-by-date", "date-ascending", "score"}

		Convey("When we call GetSortOrder(), it returns an error and the default sort option", func() {
			for _, bso := range badSortOptions {
				v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{bso}}, RelDateDesc.String())

				So(v, ShouldEqual, RelDateDesc)
				So(e, ShouldNotBeNil)
				So(e.Error(), ShouldEqual, "invalid sort parameter: invalid sort option string")
			}
		})
	})

	Convey("given a set of good sort string options", t, func() {
		goodSortOptions := []struct {
			given   string
			exValue Sort
		}{
			{given: "date-oldest", exValue: RelDateAsc},
			{given: "date-newest", exValue: RelDateDesc},
			{given: "alphabetical-az", exValue: TitleAZ},
			{given: "alphabetical-za", exValue: TitleZA},
		}

		Convey("When we call GetSortOrder(), it returns the right value", func() {
			for _, gso := range goodSortOptions {
				v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{gso.given}}, RelDateDesc.String())
				So(v, ShouldEqual, gso.exValue)
				So(e, ShouldBeNil)
			}
		})
	})

	Convey("given the 'relevance' sort option", t, func() {
		sortOption := []string{"relevance"}

		Convey("When keywords have been set", func() {
			params := url.Values{SortName: sortOption, Keywords: []string{"keywords set"}}
			Convey("Then GetSortOrder() returns the relevance value", func() {
				v, e := GetSortOrder(context.Background(), params, RelDateDesc.String())
				So(v, ShouldEqual, Relevance)
				So(e, ShouldBeNil)
			})
		})

		Convey("When keywords have not been set", func() {
			params := url.Values{SortName: sortOption}

			Convey("And a valid default sort option has been configured", func() {
				defaultSort := RelDateAsc
				Convey("Then GetSortOrder() returns the default value", func() {
					v, e := GetSortOrder(context.Background(), params, defaultSort.String())
					So(v, ShouldEqual, RelDateAsc)
					So(e, ShouldBeNil)
				})
			})

			Convey("And a wrong default sort option has been configured", func() {
				defaultSort := "random"
				Convey("Then GetSortOrder() returns the RelDateDesc value", func() {
					v, e := GetSortOrder(context.Background(), params, defaultSort)
					So(v, ShouldEqual, RelDateDesc)
					So(e, ShouldBeNil)
				})
			})
		})
	})

	Convey("given a wrong default sort option", t, func() {
		defaultSort := "random"

		Convey("And no sort parameter has been provided", func() {
			params := url.Values{}
			Convey("Then GetSortOrder() returns the RelDateDesc value", func() {
				v, e := GetSortOrder(context.Background(), params, defaultSort)
				So(v, ShouldEqual, RelDateDesc)
				So(e, ShouldBeNil)
			})
		})
	})
}

func TestGetKeywords(t *testing.T) {
	Convey("given a keyword string", t, func() {
		var keywords string
		Convey("if the string is empty, the default value is passed back as being verified", func() {
			def := "default keywords"
			v, e := GetKeywords(context.Background(), url.Values{}, def)

			So(v, ShouldEqual, def)
			So(e, ShouldBeNil)
		})
		Convey("if the string is not empty, the unaltered string is passed back as being verified", func() {
			keywords = "a b cd"
			v, e := GetKeywords(context.Background(), url.Values{Keywords: []string{keywords}}, "default")

			So(v, ShouldEqual, keywords)
			So(e, ShouldBeNil)
		})
	})
}

func TestGetBoolean(t *testing.T) {
	Convey("Given a set of strings to be parsed as a boolean value", t, func() {
		defaultValue := true
		paramName := "bool-attr"
		Convey("And they are not valid representations of a boolean value", func() {
			bvs := []string{"not right", "correct", "right", "wrong", "maybe"}
			for _, bb := range bvs {
				Convey(fmt.Sprintf("When we call the GetBoolean function for a parameter with value %s", bb), func() {
					values := url.Values{paramName: []string{bb}}
					v, err := GetBoolean(context.Background(), values, paramName, defaultValue)
					Convey("Then an error is returned", func() {
						So(v, ShouldEqual, defaultValue)
						So(err, ShouldResemble, fmt.Errorf(`invalid boolean value for parameter "%s"`, paramName))
					})
				})
			}
		})
		Convey("And they are valid representations of a boolean value", func() {
			bvs := map[string]bool{"false": false, "T": true, "TRUE": true, "0": false, "1": true}
			for bb, expected := range bvs {
				Convey(fmt.Sprintf("When we call the GetBoolean function for a parameter with value %s", bb), func() {
					values := url.Values{paramName: []string{bb}}
					v, err := GetBoolean(context.Background(), values, paramName, defaultValue)
					Convey("Then the right boolean value is returned", func() {
						So(v, ShouldEqual, expected)
						So(err, ShouldBeNil)
					})
				})
			}
		})
		Convey("And they are empty strings", func() {
			Convey("When we call the GetBoolean function", func() {
				values := url.Values{paramName: {""}}
				v, err := GetBoolean(context.Background(), values, paramName, defaultValue)
				So(v, ShouldEqual, defaultValue)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetReleaseType(t *testing.T) {
	Convey("given a set of erroneous release-type option strings", t, func() {
		badReleaseTypes := []string{"coming-up", "finished", "done"}
		defaultType := Published
		Convey("When we call GetReleaseType(), it returns an error and the default sort option", func() {
			for _, rt := range badReleaseTypes {
				v, e := GetReleaseType(context.Background(), url.Values{Type: []string{rt}}, defaultType)

				So(v, ShouldEqual, defaultType)
				So(e, ShouldNotBeNil)
				So(e.Error(), ShouldEqual, "invalid release-type parameter: invalid release type string")
			}
		})
	})

	Convey("given a set of good release-type option", t, func() {
		goodReleaseTypes := []struct {
			given   string
			exValue ReleaseType
		}{
			{given: "type-upcoming", exValue: Upcoming},
			{given: "type-published", exValue: Published},
			{given: "type-cancelled", exValue: Cancelled},
		}
		defaultType := Published
		Convey("When we call GetReleaseType(), it returns the right value", func() {
			for _, grt := range goodReleaseTypes {
				v, e := GetReleaseType(context.Background(), url.Values{Type: []string{grt.given}}, defaultType)

				So(v, ShouldEqual, grt.exValue)
				So(e, ShouldBeNil)
			}
		})
	})
}

func TestGetDates(t *testing.T) {
	Convey("given a set of date parameters as strings", t, func() {
		testcases := []struct {
			testDescription                    string
			afterDay, afterMonth, afterYear    string
			beforeDay, beforeMonth, beforeYear string
			exFromDate, exToDate               string
			exFromError, exToError             []model.ErrorItem
		}{
			{
				testDescription: "for missing year value",
				afterDay:        "", afterMonth: "", afterYear: "",
				beforeDay: "", beforeMonth: "", beforeYear: "",
				exFromDate: "", exToDate: "",
				exFromError: nil, // year is only required when day and/or month provided
				exToError:   nil, // year is only required when day and/or month provided
			},
			{
				testDescription: "for valid day and missing year value",
				afterDay:        "1", afterMonth: "", afterYear: "",
				beforeDay: "", beforeMonth: "", beforeYear: "",
				exFromDate: "", exToDate: "",
				exFromError: []model.ErrorItem{
					{
						Description: model.Localisation{
							Text: "Enter a released AFTER year",
						},
						ID:  "fromDate-error",
						URL: "#fromDate-error",
					},
				},
			},
			{
				testDescription: "for valid month and year value assumed day set",
				afterDay:        "", afterMonth: "11", afterYear: "2021",
				beforeDay: "", beforeMonth: "5", beforeYear: "2024",
				exFromDate: "2021-11-01", exToDate: "2024-05-01",
				exFromError: nil,
				exToError:   nil,
			},
			{
				testDescription: "for valid day and year value assumed month set",
				afterDay:        "5", afterMonth: "", afterYear: "2023",
				beforeDay: "31", beforeMonth: "", beforeYear: "2024",
				exFromDate: "2023-01-05", exToDate: "2024-01-31",
				exFromError: nil,
				exToError:   nil,
			},
			{
				testDescription: "for invalid day and valid year value",
				afterDay:        "35", afterMonth: "", afterYear: "2023",
				beforeDay: "", beforeMonth: "", beforeYear: "",
				exFromDate: "", exToDate: "",
				exFromError: []model.ErrorItem{
					{
						Description: model.Localisation{
							Text: "invalid 35 parameter: value is above the maximum value (31)",
						},
						ID:  "fromDate-error",
						URL: "#fromDate-error",
					},
					{
						Description: model.Localisation{
							Text: "invalid day (35) of month (1) in year (2023)",
						},
						ID:  "fromDate-error",
						URL: "#fromDate-error",
					},
				},
				exToError: nil,
			},
			{
				testDescription: "for invalid day of month value",
				afterDay:        "32", afterMonth: "2", afterYear: "2021",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2021",
				exFromDate: "", exToDate: "2021-12-31",
				exFromError: []model.ErrorItem{
					{
						Description: model.Localisation{
							Text: "invalid 32 parameter: value is above the maximum value (31)",
						},
						ID:  "fromDate-error",
						URL: "#fromDate-error",
					},
					{
						Description: model.Localisation{
							Text: "invalid day (32) of month (2) in year (2021)",
						},
						ID:  "fromDate-error",
						URL: "#fromDate-error",
					},
				},
			},
			{
				testDescription: "for invalid 29th of February outside of leap year",
				afterDay:        "29", afterMonth: "2", afterYear: "2021",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2021",
				exFromDate: "", exToDate: "2021-12-31",
				exFromError: []model.ErrorItem{
					{
						Description: model.Localisation{
							Text: "invalid day (29) of month (2) in year (2021)",
						},
						ID:  "fromDate-error",
						URL: "#fromDate-error",
					},
				},
			},
			{
				testDescription: "for valid 29th February on leap year",
				afterDay:        "29", afterMonth: "2", afterYear: "2020",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2020",
				exFromDate: "2020-02-29", exToDate: "2020-12-31",
				exFromError: nil,
				exToError:   nil,
			},
			{
				testDescription: "for valid date range: start date before end date",
				afterDay:        "28", afterMonth: "2", afterYear: "2021",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2021",
				exFromDate: "2021-02-28", exToDate: "2021-12-31",
				exFromError: nil,
			},
			{
				testDescription: "for invalid date range: start date after end date",
				afterDay:        "28", afterMonth: "2", afterYear: "2021",
				beforeDay: "1", beforeMonth: "02", beforeYear: "2021",
				exFromDate: "2021-02-28", exToDate: "2021-02-01",
				exFromError: nil,
			},
		}

		Convey("check that the validator correctly validates the dates, giving the expected results", func() {
			for _, tc := range testcases {
				Convey(tc.testDescription, func() {
					params := make(url.Values)
					params.Set("after-year", tc.afterYear)
					params.Set("after-month", tc.afterMonth)
					params.Set("after-day", tc.afterDay)
					params.Set("before-year", tc.beforeYear)
					params.Set("before-month", tc.beforeMonth)
					params.Set("before-day", tc.beforeDay)

					from, fErr := GetStartDate(context.Background(), params)
					to, tErr := GetEndDate(context.Background(), params)

					So(fErr, ShouldResemble, tc.exFromError)
					So(tErr, ShouldResemble, tc.exToError)
					So(from.String(), ShouldEqual, tc.exFromDate)
					So(to.String(), ShouldEqual, tc.exToDate)
				})
			}
		})
	})
}

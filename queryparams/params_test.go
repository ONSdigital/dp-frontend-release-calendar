package queryparams

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"testing"

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
				{value: "XXX", exValue: 0, exError: errors.New("Value contains non numeric characters")},
				{value: "-1", exValue: 0, exError: errors.New("Value is below the minimum value (0)")},
				{value: "1001", exValue: 0, exError: fmt.Errorf("Value is above the maximum value (1000)")},
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
						So(err.Error(), ShouldEqual, "Value is below the minimum value (0)")
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
						So(err.Error(), ShouldEqual, "Value is above the maximum value (55)")
					})
				})
			})
			Convey("And it is not a number", func() {
				params.Set("limit", "seven")
				Convey("When we call GetLimit", func() {
					_, err := GetLimit(ctx, params, defaultValue, maxValue)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "Value contains non numeric characters")
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
						So(err.Error(), ShouldEqual, "Value is below the minimum value (1)")
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
						So(err.Error(), ShouldEqual, "Value is above the maximum value (10)")
					})
				})
			})
			Convey("And it is not a number", func() {
				params.Set("page", "three")
				Convey("When we call GetPage", func() {
					_, err := GetPage(ctx, params, maxPage)
					Convey("Then an error is returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "Value contains non numeric characters")
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

func TestSort(t *testing.T) {
	Convey("given a set of erroneous sort string options", t, func() {
		badSortOptions := []string{"dont sort", "sort-by-date", "date-ascending", "score"}

		Convey("parsing produces an error and returns the Invalid sort option", func() {
			for _, bso := range badSortOptions {
				v, e := ParseSort(bso)

				So(v, ShouldEqual, Invalid)
				So(e, ShouldNotBeNil)
			}
		})

		Convey("and calling GetSortOrder() returns the same error and sort option", func() {
			for _, bso := range badSortOptions {
				v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{bso}}, RelDateDesc)

				So(v, ShouldEqual, Invalid)
				So(e, ShouldNotBeNil)
			}
		})

		Convey("but a good sort option string is parsed without error, and the appropriate Sort option returned", func() {
			goodSortOptions := []struct {
				given   string
				exValue Sort
			}{
				{given: "date-oldest", exValue: RelDateAsc},
				{given: "date-newest", exValue: RelDateDesc},
				{given: "alphabetical-az", exValue: TitleAZ},
				{given: "alphabetical-za", exValue: TitleZA},
			}

			for _, gso := range goodSortOptions {
				v, e := ParseSort(gso.given)

				So(v, ShouldEqual, gso.exValue)
				So(e, ShouldBeNil)

				v, e = GetSortOrder(context.Background(), url.Values{SortName: []string{gso.given}}, RelDateDesc)
				So(v, ShouldEqual, gso.exValue)
				So(e, ShouldBeNil)

			}
		})
		Convey("except for the 'relevance' sort option - this parses as normal", func() {
			v, e := ParseSort("relevance")

			So(v, ShouldEqual, Relevance)
			So(e, ShouldBeNil)

			Convey("but can only be set if a keyword has also been set", func() {
				v, e = GetSortOrder(context.Background(), url.Values{SortName: []string{"relevance"}, Keywords: []string{"keywords set"}}, RelDateDesc)
				So(v, ShouldEqual, Relevance)
				So(e, ShouldBeNil)

				v, e = GetSortOrder(context.Background(), url.Values{SortName: []string{"relevance"}}, RelDateDesc)
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

func TestReleaseType(t *testing.T) {
	Convey("given a set of erroneous release-type option strings", t, func() {
		badReleaseTypes := []string{"coming-up", "finished", "done"}

		Convey("parsing produces an error and returns the InvalidReleaseType ReleaseType", func() {
			for _, rt := range badReleaseTypes {
				v, e := ParseReleaseType(rt)

				So(v, ShouldEqual, InvalidReleaseType)
				So(e, ShouldNotBeNil)
			}
		})

		Convey("but a good release-type option string is parsed without error, and the appropriate ReleaseType returned", func() {
			goodReleaseTypes := []struct {
				given   string
				exValue ReleaseType
			}{
				{given: "type-upcoming", exValue: Upcoming},
				{given: "type-published", exValue: Published},
				{given: "type-cancelled", exValue: Cancelled},
			}

			for _, grt := range goodReleaseTypes {
				v, e := ParseReleaseType(grt.given)

				So(v, ShouldEqual, grt.exValue)
				So(e, ShouldBeNil)
			}
		})
	})
}

func TestDatesFromParams(t *testing.T) {
	Convey("given a set of day month and year numbers as strings", t, func() {
		testcases := []struct {
			afterDay, afterMonth, afterYear    string
			beforeDay, beforeMonth, beforeYear string
			exFromDate, exToDate               string
			exError                            error
		}{
			{
				afterDay: "32", afterMonth: "2", afterYear: "2021",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2021",
				exFromDate: "", exToDate: "",
				exError: errors.New("Value is above the maximum value (31)"),
			},
			{
				afterDay: "29", afterMonth: "2", afterYear: "2021",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2021",
				exFromDate: "", exToDate: "",
				exError: errors.New("invalid day (29) of month (2) in year (2021)"),
			},
			{
				afterDay: "28", afterMonth: "2", afterYear: "2021",
				beforeDay: "31", beforeMonth: "12", beforeYear: "2021",
				exFromDate: "2021-02-28", exToDate: "2021-12-31",
				exError: nil,
			},
			{
				afterDay: "28", afterMonth: "2", afterYear: "2021",
				beforeDay: "1", beforeMonth: "02", beforeYear: "2021",
				exFromDate: "", exToDate: "",
				exError: errors.New("invalid dates - 'after' after 'before'"),
			},
		}

		Convey("check that the validator correctly validates the dates, giving the expected results", func() {
			for _, tc := range testcases {
				params := make(url.Values)
				params.Set("after-year", tc.afterYear)
				params.Set("after-month", tc.afterMonth)
				params.Set("after-day", tc.afterDay)
				params.Set("before-year", tc.beforeYear)
				params.Set("before-month", tc.beforeMonth)
				params.Set("before-day", tc.beforeDay)

				from, to, err := DatesFromParams(context.Background(), params)

				So(err, ShouldResemble, tc.exError)
				So(from.String(), ShouldEqual, tc.exFromDate)
				So(to.String(), ShouldEqual, tc.exToDate)
			}
		})
	})
}

func TestParamsAsFrontendQuery(t *testing.T) {
	Convey("Given a set of validated parameters as a ValidatedParam struct", t, func() {
		vp := ValidatedParams{
			Limit:       10,
			Page:        2,
			Offset:      10,
			AfterDate:   MustParseDate("2020-01-01"),
			Keywords:    "some keywords",
			Sort:        TitleAZ,
			Provisional: true,
			Confirmed:   true,
			Postponed:   true,
			Census:      true,
		}

		Convey("And the release type is upcoming", func() {
			vp.ReleaseType = Upcoming

			Convey("When we call AsFrontendQuery", func() {
				uv := vp.AsFrontendQuery()
				Convey("Then the validated parameters are correctly returned in an url.Values mapping", func() {
					So(uv.Get(Limit), ShouldEqual, "10")
					So(uv.Get(Page), ShouldEqual, "2")
					So(uv.Get(YearAfter), ShouldEqual, "2020")
					So(uv.Get(MonthAfter), ShouldEqual, "1")
					So(uv.Get(DayAfter), ShouldEqual, "1")
					So(uv.Get(YearBefore), ShouldEqual, "")
					So(uv.Get(MonthBefore), ShouldEqual, "")
					So(uv.Get(DayBefore), ShouldEqual, "")
					So(uv.Get(Keywords), ShouldEqual, "some keywords")
					So(uv.Get(SortName), ShouldEqual, vp.Sort.String())
					So(uv.Get(Type), ShouldEqual, vp.ReleaseType.String())
					So(uv.Get(Provisional.String()), ShouldEqual, "true")
					So(uv.Get(Confirmed.String()), ShouldEqual, "true")
					So(uv.Get(Postponed.String()), ShouldEqual, "true")
					So(uv.Get(Census), ShouldEqual, "true")
					So(uv.Get(Highlight), ShouldEqual, "")

					Convey("And any validated parameters not needed are absent from the url.Values mapping", func() {
						So(uv.Get(Offset), ShouldEqual, "")
						So(uv.Get(Query), ShouldEqual, "")
					})
				})
			})
		})

		Convey("And the release type is not upcoming", func() {
			vp.ReleaseType = Cancelled

			Convey("When we call AsFrontendQuery", func() {
				uv := vp.AsFrontendQuery()
				Convey("Then the validated parameters are correctly returned in an url.Values mapping", func() {
					So(uv.Get(Limit), ShouldEqual, "10")
					So(uv.Get(Page), ShouldEqual, "2")
					So(uv.Get(YearAfter), ShouldEqual, "2020")
					So(uv.Get(MonthAfter), ShouldEqual, "1")
					So(uv.Get(DayAfter), ShouldEqual, "1")
					So(uv.Get(YearBefore), ShouldEqual, "")
					So(uv.Get(MonthBefore), ShouldEqual, "")
					So(uv.Get(DayBefore), ShouldEqual, "")
					So(uv.Get(Keywords), ShouldEqual, "some keywords")
					So(uv.Get(SortName), ShouldEqual, vp.Sort.String())
					So(uv.Get(Type), ShouldEqual, vp.ReleaseType.String())
					So(uv.Get(Census), ShouldEqual, "true")
					So(uv.Get(Highlight), ShouldEqual, "")

					Convey("And any validated parameters not needed are absent from the url.Values mapping", func() {
						So(uv.Get(Offset), ShouldEqual, "")
						So(uv.Get(Query), ShouldEqual, "")
						So(uv.Get(Provisional.String()), ShouldEqual, "")
						So(uv.Get(Confirmed.String()), ShouldEqual, "")
						So(uv.Get(Postponed.String()), ShouldEqual, "")
					})
				})
			})
		})
	})
}

func TestParamsAsBackendQuery(t *testing.T) {
	Convey("Given a set of validated parameters as a ValidatedParam struct", t, func() {
		vp := ValidatedParams{
			Limit:       10,
			Page:        2,
			Offset:      10,
			AfterDate:   MustParseDate("2020-01-01"),
			BeforeDate:  MustParseDate("2022-09-19"),
			Keywords:    "some keywords",
			Provisional: true,
			Confirmed:   true,
			Postponed:   true,
			Highlight:   true,
			Census:      false,
		}

		Convey("And the release type is upcoming", func() {
			vp.ReleaseType = Upcoming
			Convey("And we are sorting by date in ascending order", func() {
				vp.Sort = RelDateAsc
				Convey("When we call AsBackendQuery", func() {
					uv := vp.AsBackendQuery()
					Convey("Then the validated parameters are correctly returned in an url.Values mapping", func() {
						So(uv.Get(Limit), ShouldEqual, "10")
						So(uv.Get(Page), ShouldEqual, "2")
						So(uv.Get(Offset), ShouldEqual, "10")
						So(uv.Get(YearAfter), ShouldEqual, "2020")
						So(uv.Get(MonthAfter), ShouldEqual, "1")
						So(uv.Get(DayAfter), ShouldEqual, "1")
						So(uv.Get(YearBefore), ShouldEqual, "2022")
						So(uv.Get(MonthBefore), ShouldEqual, "9")
						So(uv.Get(DayBefore), ShouldEqual, "19")
						So(uv.Get(Query), ShouldEqual, "some keywords")
						So(uv.Get(Type), ShouldEqual, vp.ReleaseType.String())
						So(uv.Get(Provisional.String()), ShouldEqual, "true")
						So(uv.Get(Confirmed.String()), ShouldEqual, "true")
						So(uv.Get(Postponed.String()), ShouldEqual, "true")
						So(uv.Get(Census), ShouldEqual, "")
						So(uv.Get(Highlight), ShouldEqual, "true")

						Convey("And the date sort order is inverted", func() {
							So(uv.Get(SortName), ShouldEqual, RelDateDesc.BackendString())
						})
						Convey("And any validated parameters not needed are absent from the url.Values mapping", func() {
							So(uv.Get(Keywords), ShouldEqual, "")
						})
					})
				})
			})
			Convey("And we are sorting by date in descending order", func() {
				vp.Sort = RelDateDesc
				Convey("When we call AsBackendQuery", func() {
					uv := vp.AsBackendQuery()
					Convey("Then the validated parameters are correctly returned in an url.Values mapping", func() {
						So(uv.Get(Limit), ShouldEqual, "10")
						So(uv.Get(Page), ShouldEqual, "2")
						So(uv.Get(Offset), ShouldEqual, "10")
						So(uv.Get(YearAfter), ShouldEqual, "2020")
						So(uv.Get(MonthAfter), ShouldEqual, "1")
						So(uv.Get(DayAfter), ShouldEqual, "1")
						So(uv.Get(YearBefore), ShouldEqual, "2022")
						So(uv.Get(MonthBefore), ShouldEqual, "9")
						So(uv.Get(DayBefore), ShouldEqual, "19")
						So(uv.Get(Query), ShouldEqual, "some keywords")
						So(uv.Get(Type), ShouldEqual, vp.ReleaseType.String())
						So(uv.Get(Provisional.String()), ShouldEqual, "true")
						So(uv.Get(Confirmed.String()), ShouldEqual, "true")
						So(uv.Get(Postponed.String()), ShouldEqual, "true")
						So(uv.Get(Census), ShouldEqual, "")
						So(uv.Get(Highlight), ShouldEqual, "true")

						Convey("And the date sort order is inverted", func() {
							So(uv.Get(SortName), ShouldEqual, RelDateAsc.BackendString())
						})
						Convey("And any validated parameters not needed are absent from the url.Values mapping", func() {
							So(uv.Get(Keywords), ShouldEqual, "")
						})
					})
				})
			})
		})

		Convey("And the release type is not upcoming", func() {
			vp.ReleaseType = Published

			Convey("When we call AsBackendQuery", func() {
				uv := vp.AsBackendQuery()
				Convey("Then the validated parameters are correctly returned in an url.Values mapping", func() {
					So(uv.Get(Limit), ShouldEqual, "10")
					So(uv.Get(Page), ShouldEqual, "2")
					So(uv.Get(Offset), ShouldEqual, "10")
					So(uv.Get(YearAfter), ShouldEqual, "2020")
					So(uv.Get(MonthAfter), ShouldEqual, "1")
					So(uv.Get(DayAfter), ShouldEqual, "1")
					So(uv.Get(YearBefore), ShouldEqual, "2022")
					So(uv.Get(MonthBefore), ShouldEqual, "9")
					So(uv.Get(DayBefore), ShouldEqual, "19")
					So(uv.Get(Query), ShouldEqual, "some keywords")
					So(uv.Get(SortName), ShouldEqual, vp.Sort.BackendString())
					So(uv.Get(Type), ShouldEqual, vp.ReleaseType.String())
					So(uv.Get(Census), ShouldEqual, "")
					So(uv.Get(Highlight), ShouldEqual, "true")

					Convey("And any validated parameters not needed are absent from the url.Values mapping", func() {
						So(uv.Get(Keywords), ShouldEqual, "")
						So(uv.Get(Provisional.String()), ShouldEqual, "")
						So(uv.Get(Confirmed.String()), ShouldEqual, "")
						So(uv.Get(Postponed.String()), ShouldEqual, "")
					})
				})
			})
		})
	})
}

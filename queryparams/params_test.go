package queryparams

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntValidator(t *testing.T) {
	Convey("given an IntValidator parameterised with a maximum and minimum value", t, func() {
		validator := GetIntValidator(0, 1000)

		Convey("and a set of int values as strings representing a given parameter name", func() {

			limits := []struct {
				name    string
				value   string
				exValue int
				exError error
			}{
				{name: "limit", value: "XXX", exValue: 0, exError: errors.New("limit search parameter provided with non numeric characters")},
				{name: "limit", value: "-1", exValue: 0, exError: errors.New("limit search parameter provided with a value that is below the minimum value")},
				{name: "limit", value: "1001", exValue: 0, exError: fmt.Errorf("limit search parameter provided with a value that is above the maximum value")},
				{name: "limit", value: "0", exValue: 0, exError: nil},
				{name: "limit", value: "123", exValue: 123, exError: nil},
				{name: "limit", value: "1000", exValue: 1000, exError: nil},
			}

			Convey("check that the validator correctly validates the limit, giving the expected results", func() {
				for _, ls := range limits {
					v, e := validator(ls.name, ls.value)

					So(v, ShouldEqual, ls.exValue)
					So(e, ShouldResemble, ls.exError)
				}
			})
		})
	})
}

func TestGetLimit(t *testing.T) {
	Convey("given an IntValidator for a limit, and a set of limits as strings", t, func() {
		validator := GetIntValidator(0, 1000)
		limits := []struct {
			given   string
			exValue int
			exError error
		}{
			{given: "XXX", exValue: 0, exError: errors.New("limit search parameter provided with non numeric characters")},
			{given: "-1", exValue: 0, exError: errors.New("limit search parameter provided with a value that is below the minimum value")},
			{given: "1001", exValue: 0, exError: fmt.Errorf("limit search parameter provided with a value that is above the maximum value")},
			{given: "0", exValue: 0, exError: nil},
			{given: "1000", exValue: 1000, exError: nil},
		}

		Convey("check that the validator correctly validates the limit", func() {
			for _, ls := range limits {
				v, e := validator(Limit, ls.given)

				So(v, ShouldEqual, ls.exValue)
				So(e, ShouldResemble, ls.exError)
			}
		})
	})
}

func TestPageValidator(t *testing.T) {
	Convey("given a page validator, and a set of page numbers as strings", t, func() {
		validator := GetIntValidator(1, 100)
		offsets := []struct {
			given   string
			exValue int
			exError error
		}{
			{given: "XXX", exValue: 0, exError: errors.New("page search parameter provided with non numeric characters")},
			{given: "0", exValue: 0, exError: errors.New("page search parameter provided with a value that is below the minimum value")},
			{given: "-1", exValue: 0, exError: errors.New("page search parameter provided with a value that is below the minimum value")},
			{given: "101", exValue: 0, exError: errors.New("page search parameter provided with a value that is above the maximum value")},
			{given: "1", exValue: 1, exError: nil},
			{given: "100", exValue: 100, exError: nil},
		}

		Convey("check that the validator correctly validates the page number", func() {
			for _, ps := range offsets {
				v, e := validator(Page, ps.given)

				So(v, ShouldEqual, ps.exValue)
				So(e, ShouldResemble, ps.exError)
			}
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

func TestParseSortFromFrontend(t *testing.T) {
	Convey("given a release type of Published", t, func() {
		releaseType := Published

		Convey("when the sort order is not given", func() {
			sort, err := ParseSortFromFrontend("", releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[PublishedNewest])
		})

		Convey("when the sort order is 'newest'", func() {
			sort, err := ParseSortFromFrontend(Newest, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[PublishedNewest])
		})

		Convey("when the sort order is 'oldest'", func() {
			sort, err := ParseSortFromFrontend(Oldest, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[PublishedOldest])
		})

		Convey("when the sort order is not a date option (such as 'relevance')", func() {
			sort, err := ParseSortFromFrontend(RelevanceLabel, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[Relevance])
		})
	})

	Convey("given a release type of Upcoming", t, func() {
		releaseType := Upcoming

		Convey("when the sort order is not given", func() {
			sort, err := ParseSortFromFrontend("", releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[UpcomingNewest])
		})

		Convey("when the sort order is 'newest", func() {
			sort, err := ParseSortFromFrontend(Newest, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[UpcomingNewest])
		})

		Convey("when the sort order is 'oldest'", func() {
			sort, err := ParseSortFromFrontend(Oldest, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[UpcomingOldest])
		})

		Convey("when the sort order is not a date option (such as 'alphabetical A-Z')", func() {
			sort, err := ParseSortFromFrontend(AlphaUser, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[Alpha])
		})
	})

	Convey("given a release type of Cancelled", t, func() {
		releaseType := Cancelled

		Convey("when the sort order is not given", func() {
			sort, err := ParseSortFromFrontend("", releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[PublishedNewest])
		})

		Convey("when the sort order is 'newest'", func() {
			sort, err := ParseSortFromFrontend(Newest, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[PublishedNewest])
		})

		Convey("when the sort order is 'oldest'", func() {
			sort, err := ParseSortFromFrontend(Oldest, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[PublishedOldest])
		})

		Convey("when the sort order is not a date option (such as 'alphabetical Z-A')", func() {
			sort, err := ParseSortFromFrontend(ReverseAlphaUser, releaseType)
			So(err, ShouldBeNil)
			So(sort, ShouldResemble, sortValues[ReverseAlpha])
		})
	})
}

func TestSortOrder(t *testing.T) {
	Convey("given a valid release type", t, func() {
		releaseType := Published

		Convey("when GetSortOrder() is called with an invalid sort option", func() {
			v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{"dont sort"}}, releaseType, MustParseSort(RelevanceLabel))

			So(e, ShouldNotBeNil)
			So(v, ShouldResemble, sortValues[InvalidSort])
		})

		Convey("when GetSortOrder() is called with a valid sort option such as 'alphabetical AZ'", func() {
			v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{AlphaUser}}, releaseType, MustParseSort(RelevanceLabel))

			So(e, ShouldBeNil)
			So(v, ShouldResemble, sortValues[Alpha])
		})

		Convey("when GetSortOrder() is called with the 'relevance' sort option and a keyword has been set", func() {
			v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{RelevanceLabel}, Keywords: []string{"keywords set"}}, releaseType, MustParseSort(AlphaUser))

			So(e, ShouldBeNil)
			So(v, ShouldResemble, sortValues[Relevance])
		})

		Convey("when GetSortOrder() is called with the 'relevance' sort option but a keyword has NOT been set", func() {
			v, e := GetSortOrder(context.Background(), url.Values{SortName: []string{RelevanceLabel}}, releaseType, MustParseSort(AlphaUser))

			So(e, ShouldBeNil)
			So(v, ShouldResemble, sortValues[Alpha])
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
	Convey("given a set of strings to be parsed as a boolean value", t, func() {
		var bvs []string
		Convey("if the strings are not valid representations of a boolean value", func() {
			bvs = []string{"not right", "correct", "right", "wrong", "maybe"}
			Convey("the correct error is returned", func() {
				for _, bb := range bvs {
					v, s, err := GetBoolean(context.Background(), url.Values{"bool-attr": []string{bb}}, "bool-attr", true)
					So(v, ShouldBeFalse)
					So(s, ShouldBeFalse)
					So(err, ShouldResemble, errors.New(`invalid boolean value for "bool-attr"`))
				}
			})
			Convey("if the strings are valid representations of a boolean value", func() {
				bvs = []string{"false", "T", "TRUE", "0", "1"}
				Convey("the correct boolean value is returned without error", func() {
					for _, gb := range bvs {
						v, s, err := GetBoolean(context.Background(), url.Values{"bool-attr": []string{gb}}, "bool-attr", true)
						So(v, ShouldBeIn, true, false)
						So(s, ShouldBeTrue)
						So(err, ShouldBeNil)
					}
				})
			})
			Convey("if the strings are empty strings", func() {
				Convey("the correct boolean value corresponding to the default value is returned without error", func() {
					v, s, err := GetBoolean(context.Background(), url.Values{}, "bool-attr", true)
					So(v, ShouldBeTrue)
					So(s, ShouldBeFalse)
					So(err, ShouldBeNil)
				})
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
				exError: errors.New("after-day search parameter provided with a value that is above the maximum value"),
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

func TestParamsAsQuery(t *testing.T) {
	Convey("given a set of validated parameters as a ValidatedParam struct", t, func() {
		vp := ValidatedParams{Limit: 10, Page: 2, Offset: 10, AfterDate: MustParseDate("2020-01-01"), Keywords: "some keywords", Sort: sortValues[Alpha], ReleaseType: Upcoming, Provisional: true, Census: true}

		Convey("verify that the validated parameters are correctly returned in an url.Values mapping", func() {
			uv := vp.AsQuery()
			So(uv.Get(Limit), ShouldEqual, "10")
			So(uv.Get(Page), ShouldEqual, "2")
			So(uv.Get(YearAfter), ShouldEqual, "2020")
			So(uv.Get(MonthAfter), ShouldEqual, "1")
			So(uv.Get(DayAfter), ShouldEqual, "1")
			So(uv.Get(YearBefore), ShouldEqual, "")
			So(uv.Get(MonthBefore), ShouldEqual, "")
			So(uv.Get(DayBefore), ShouldEqual, "")
			So(uv.Get(Keywords), ShouldEqual, "some keywords")
			So(uv.Get(SortName), ShouldEqual, sortValues[Alpha].feValue)
			So(uv.Get(Type), ShouldEqual, Upcoming.String())
			So(uv.Get(Provisional.String()), ShouldEqual, "true")
			So(uv.Get(Confirmed.String()), ShouldEqual, "false")
			So(uv.Get(Postponed.String()), ShouldEqual, "false")
			So(uv.Get(Census), ShouldEqual, "true")
			So(uv.Get(Highlight), ShouldEqual, "false")

			Convey("and any validated parameters not needed are absent from the url.Values mapping", func() {
				So(uv.Get(Offset), ShouldEqual, "")
			})
		})
	})
}

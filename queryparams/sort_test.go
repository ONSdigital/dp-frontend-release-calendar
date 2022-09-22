package queryparams

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSortString(t *testing.T) {
	Convey("Given a set of sort values", t, func() {
		testcases := []struct {
			value          Sort
			expectedString string
		}{
			{
				value:          RelDateAsc,
				expectedString: "date-oldest",
			},
			{
				value:          RelDateDesc,
				expectedString: "date-newest",
			},
			{
				value:          TitleAZ,
				expectedString: "alphabetical-az",
			},
			{
				value:          TitleZA,
				expectedString: "alphabetical-za",
			},
			{
				value:          Relevance,
				expectedString: "relevance",
			},
		}

		Convey("When we call String(), we get the appropriate front end value ", func() {
			for _, tc := range testcases {
				So(tc.value.String(), ShouldEqual, tc.expectedString)
			}
		})
	})
}

func TestSortBackendString(t *testing.T) {
	Convey("Given a set of sort values", t, func() {
		testcases := []struct {
			value          Sort
			expectedString string
		}{
			{
				value:          RelDateAsc,
				expectedString: "release_date_asc",
			},
			{
				value:          RelDateDesc,
				expectedString: "release_date_desc",
			},
			{
				value:          TitleAZ,
				expectedString: "title_asc",
			},
			{
				value:          TitleZA,
				expectedString: "title_desc",
			},
			{
				value:          Relevance,
				expectedString: "relevance",
			},
		}

		Convey("When we call BackendString(), we get the appropriate front end value ", func() {
			for _, tc := range testcases {
				So(tc.value.BackendString(), ShouldEqual, tc.expectedString)
			}
		})
	})
}

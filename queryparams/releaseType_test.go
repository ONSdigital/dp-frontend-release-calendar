package queryparams

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReleaseTypeName(t *testing.T) {
	Convey("Given a set of sort values", t, func() {
		testcases := []struct {
			value          ReleaseType
			expectedString string
		}{
			{
				value:          Upcoming,
				expectedString: "type-upcoming",
			},
			{
				value:          Published,
				expectedString: "type-published",
			},
			{
				value:          Cancelled,
				expectedString: "type-cancelled",
			},
			{
				value:          Provisional,
				expectedString: "subtype-provisional",
			},
			{
				value:          Confirmed,
				expectedString: "subtype-confirmed",
			},
			{
				value:          Postponed,
				expectedString: "subtype-postponed",
			},
		}

		Convey("When we call Name(), we get the appropriate value ", func() {
			for _, tc := range testcases {
				So(tc.value.Name(), ShouldEqual, tc.expectedString)
			}
		})
	})
}

func TestReleaseTypeString(t *testing.T) {
	Convey("Given a set of sort values", t, func() {
		testcases := []struct {
			value          ReleaseType
			expectedString string
		}{
			{
				value:          Upcoming,
				expectedString: "type-upcoming",
			},
			{
				value:          Published,
				expectedString: "type-published",
			},
			{
				value:          Cancelled,
				expectedString: "type-cancelled",
			},
			{
				value:          Provisional,
				expectedString: "subtype-provisional",
			},
			{
				value:          Confirmed,
				expectedString: "subtype-confirmed",
			},
			{
				value:          Postponed,
				expectedString: "subtype-postponed",
			},
		}

		Convey("When we call String(), we get the appropriate value ", func() {
			for _, tc := range testcases {
				So(tc.value.String(), ShouldEqual, tc.expectedString)
			}
		})
	})
}

func TestReleaseTypeLabel(t *testing.T) {
	Convey("Given a set of sort values", t, func() {
		testcases := []struct {
			value          ReleaseType
			expectedString string
		}{
			{
				value:          Upcoming,
				expectedString: "Upcoming",
			},
			{
				value:          Published,
				expectedString: "Published",
			},
			{
				value:          Cancelled,
				expectedString: "Cancelled",
			},
			{
				value:          Provisional,
				expectedString: "Provisional",
			},
			{
				value:          Confirmed,
				expectedString: "Confirmed",
			},
			{
				value:          Postponed,
				expectedString: "Postponed",
			},
		}

		Convey("When we call Label(), we get the appropriate value ", func() {
			for _, tc := range testcases {
				So(tc.value.Label(), ShouldEqual, tc.expectedString)
			}
		})
	})
}

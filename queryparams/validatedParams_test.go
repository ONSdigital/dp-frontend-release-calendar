package queryparams

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAsFrontendQuery(t *testing.T) {
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

func TestAsBackendQuery(t *testing.T) {
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
						So(uv.Get(DateFrom), ShouldEqual, "2020-01-01")
						So(uv.Get(DateTo), ShouldEqual, "2022-09-19")
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
						So(uv.Get(DateFrom), ShouldEqual, "2020-01-01")
						So(uv.Get(DateTo), ShouldEqual, "2022-09-19")
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
					So(uv.Get(DateFrom), ShouldEqual, "2020-01-01")
					So(uv.Get(DateTo), ShouldEqual, "2022-09-19")
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

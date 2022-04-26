package model_test

import (
	"testing"

	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	Convey("IsFilterSearchPresent should detect the presence or absence of a search term", t, func() {
		Convey("When a search term is absent", func() {
			calendar := model.Calendar{}
			So(calendar.IsFilterSearchPresent(), ShouldBeFalse)
		})

		Convey("When a search term is present", func() {
			calendar := model.Calendar{}
			calendar.KeywordSearch.SearchTerm = "populated"
			So(calendar.IsFilterSearchPresent(), ShouldBeTrue)
		})
	})

	Convey("IsFilterDatePresent should detect the presence or absence of a date", t, func() {
		Convey("When both dates are absent", func() {
			calendar := model.Calendar{}
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)
		})

		Convey("When an AfterDate is present", func() {
			calendar := model.Calendar{}
			calendar.AfterDate.InputValueDay = "01"
			calendar.AfterDate.InputValueMonth = "01"
			calendar.AfterDate.InputValueYear = "2000"
			So(calendar.IsFilterDatePresent(), ShouldBeTrue)

			calendar.AfterDate.InputValueMonth = ""
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)

			calendar.AfterDate.InputValueMonth = "01"
			calendar.AfterDate.InputValueDay = ""
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)

			calendar.AfterDate.InputValueDay = "01"
			calendar.AfterDate.InputValueYear = ""
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)
		})

		Convey("When a BeforeDate is present", func() {
			calendar := model.Calendar{}
			calendar.BeforeDate.InputValueDay = "01"
			calendar.BeforeDate.InputValueMonth = "01"
			calendar.BeforeDate.InputValueYear = "2000"
			So(calendar.IsFilterDatePresent(), ShouldBeTrue)

			calendar.BeforeDate.InputValueMonth = ""
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)

			calendar.BeforeDate.InputValueMonth = "01"
			calendar.BeforeDate.InputValueDay = ""
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)

			calendar.BeforeDate.InputValueDay = "01"
			calendar.BeforeDate.InputValueYear = ""
			So(calendar.IsFilterDatePresent(), ShouldBeFalse)
		})

		Convey("When both dates are present", func() {
			calendar := model.Calendar{}
			calendar.AfterDate.InputValueDay = "01"
			calendar.AfterDate.InputValueMonth = "01"
			calendar.AfterDate.InputValueYear = "2000"
			calendar.BeforeDate.InputValueDay = "01"
			calendar.BeforeDate.InputValueMonth = "01"
			calendar.BeforeDate.InputValueYear = "2000"
			So(calendar.IsFilterDatePresent(), ShouldBeTrue)
		})
	})
}

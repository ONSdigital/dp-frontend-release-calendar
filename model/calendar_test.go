package model_test

import (
	"testing"

	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	Convey("FuncIsFilterSearchPresent should detect the presence or absence of a search term", t, func() {
		Convey("When a search term is absent", func() {
			calendar := model.Calendar{}
			So(calendar.FuncIsFilterSearchPresent(), ShouldBeFalse)
		})

		Convey("When a search term is present", func() {
			calendar := model.Calendar{}
			calendar.KeywordSearch.SearchTerm = "populated"
			So(calendar.FuncIsFilterSearchPresent(), ShouldBeTrue)
		})
	})

	Convey("FuncIsFilterDatePresent should detect the presence or absence of a date", t, func() {
		Convey("When both dates are absent", func() {
			calendar := model.Calendar{}
			So(calendar.FuncIsFilterDatePresent(), ShouldBeFalse)
		})

		Convey("When an AfterDate is present", func() {
			calendar := model.Calendar{}
			calendar.AfterDate.Input.InputValueDay = "01"
			calendar.AfterDate.Input.InputValueMonth = "01"
			calendar.AfterDate.Input.InputValueYear = "2000"
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)

			calendar.AfterDate.Input.InputValueMonth = "01"
			calendar.AfterDate.Input.InputValueDay = ""
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)

			calendar.AfterDate.Input.InputValueDay = "01"
			calendar.AfterDate.Input.InputValueYear = ""
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)
		})

		Convey("When a BeforeDate is present", func() {
			calendar := model.Calendar{}
			calendar.BeforeDate.Input.InputValueDay = "01"
			calendar.BeforeDate.Input.InputValueMonth = "01"
			calendar.BeforeDate.Input.InputValueYear = "2000"
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)

			calendar.BeforeDate.Input.InputValueMonth = "01"
			calendar.BeforeDate.Input.InputValueDay = ""
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)

			calendar.BeforeDate.Input.InputValueDay = "01"
			calendar.BeforeDate.Input.InputValueYear = ""
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)
		})

		Convey("When both dates are present", func() {
			calendar := model.Calendar{}
			calendar.AfterDate.Input.InputValueDay = "01"
			calendar.AfterDate.Input.InputValueMonth = "01"
			calendar.AfterDate.Input.InputValueYear = "2000"
			calendar.BeforeDate.Input.InputValueDay = "01"
			calendar.BeforeDate.Input.InputValueMonth = "01"
			calendar.BeforeDate.Input.InputValueYear = "2000"
			So(calendar.FuncIsFilterDatePresent(), ShouldBeTrue)
		})
	})

	Convey("FuncIsFilterCensusPresent should detect if census is checked", t, func() {
		Convey("When census is not checked", func() {
			calendar := model.Calendar{}
			So(calendar.FuncIsFilterCensusPresent(), ShouldBeFalse)
		})

		Convey("When census is checked", func() {
			calendar := model.Calendar{}
			calendar.ReleaseTypes = map[string]model.ReleaseType{
				"type-census": {
					Name:      "census",
					Value:     "type-census",
					ID:        "release-type-census",
					Language:  "en",
					IsChecked: true,
					Count:     2,
				},
			}
			So(calendar.FuncIsFilterCensusPresent(), ShouldBeTrue)
		})
	})
}

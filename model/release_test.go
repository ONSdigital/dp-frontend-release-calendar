package model_test

import (
	"testing"

	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitReleaseModel(t *testing.T) {
	Convey("FuncGetPostponementReason", t, func() {
		Convey("When the release is not postponed", func() {
			release := model.Release{}
			So(release.FuncGetPostponementReason(), ShouldEqual, "")
		})

		Convey("When the release is postponed but has no reason", func() {
			release := model.Release{}
			release.DateChanges = []model.DateChange{
				{
					Date:         "2020-01-01T00:00:00.000Z",
					ChangeNotice: "",
				},
			}
			So(release.FuncGetPostponementReason(), ShouldEqual, "")
		})

		Convey("When the release is postponed with a reason", func() {
			release := model.Release{}
			reason := "Postponed due to a shortage of marmalade sandwiches"
			release.DateChanges = []model.DateChange{
				{
					Date:         "2020-01-01T00:00:00.000Z",
					ChangeNotice: "",
				},
				{
					Date:         "2020-02-01T00:00:00.000Z",
					ChangeNotice: reason,
				},
			}
			So(release.FuncGetPostponementReason(), ShouldEqual, reason)
		})
	})
}

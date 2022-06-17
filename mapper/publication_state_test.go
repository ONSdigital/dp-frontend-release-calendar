package mapper

import (
	"testing"

	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitPublicationState(t *testing.T) {
	Convey("GetPublicationState", t, func() {
		Convey("Should mark am upcoming, provisional, but cancelled release as cancelled", func() {
			description := model.ReleaseDescription{
				Cancelled: true,
				Finalised: false,
				Published: false,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type: "cancelled",
			})
		})

		Convey("Should mark an upcoming, confirmed, but cancelled release as cancelled", func() {
			description := model.ReleaseDescription{
				Cancelled: true,
				Finalised: true,
				Published: false,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type: "cancelled",
			})
		})

		Convey("Should mark a published but cancelled release as cancelled", func() {
			description := model.ReleaseDescription{
				Cancelled: true,
				Finalised: true,
				Published: true,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type: "cancelled",
			})
		})

		Convey("Should mark a published, confirmed release as published", func() {
			description := model.ReleaseDescription{
				Cancelled: false,
				Finalised: true,
				Published: true,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type: "published",
			})
		})

		Convey("Should mark a published, provisional release as published", func() {
			// Florence should not allow a release to get into this state, but what
			// Florence should do and what Florence does do not necessarily align.
			description := model.ReleaseDescription{
				Cancelled: false,
				Finalised: false,
				Published: true,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type: "published",
			})
		})

		Convey("Should mark an unpublished, provisional release as <upcoming, provisional>", func() {
			description := model.ReleaseDescription{
				Cancelled: false,
				Finalised: false,
				Published: false,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type:    "upcoming",
				SubType: "provisional",
			})
		})

		Convey("Should mark an unpublished, confirmed release as <upcoming, confirmed>", func() {
			description := model.ReleaseDescription{
				Cancelled: false,
				Finalised: true,
				Published: false,
			}
			publicationState := GetPublicationState(description, []model.DateChange{})
			So(publicationState, ShouldResemble, model.PublicationState{
				Type:    "upcoming",
				SubType: "confirmed",
			})
		})

		Convey("Should mark an unpublished, confirmed release with date changes", func() {
			Convey("As <upcoming, postponed> when it is pushed back", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-16T08:30:00.000Z",
				}
				publicationState := GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Pushed back",
						Date:         "2022-01-15T08:30:00.000Z",
					},
				})
				So(publicationState, ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "postponed",
				})
			})

			Convey("As <upcoming, postponed> when it is pushed back repeatedly", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-17T08:30:00.000Z",
				}
				publicationState := GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Pushed back once",
						Date:         "2022-01-15T08:30:00.000Z",
					},
					{
						ChangeNotice: "Pushed back again",
						Date:         "2022-01-16T08:30:00.000Z",
					},
				})
				So(publicationState, ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "postponed",
				})
			})

			Convey("As <upcoming, postponed> when it is brought forward and then pushed back", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-16T08:30:00.000Z",
				}
				publicationState := GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Brought forward",
						Date:         "2022-01-15T08:30:00.000Z",
					},
					{
						ChangeNotice: "Pushed back again",
						Date:         "2022-01-14T08:30:00.000Z",
					},
				})
				So(publicationState, ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "postponed",
				})
			})

			Convey("As <upcoming, confirmed> when it is brought forward", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-14T08:30:00.000Z",
				}
				publicationState := GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Brought forward",
						Date:         "2022-01-15T08:30:00.000Z",
					},
				})
				So(publicationState, ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "confirmed",
				})
			})

			Convey("As <upcoming, confirmed> when it is brought forward repeatedly", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-13T08:30:00.000Z",
				}
				publicationState := GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Brought forward once",
						Date:         "2022-01-15T08:30:00.000Z",
					},
					{
						ChangeNotice: "Brought forward again",
						Date:         "2022-01-14T08:30:00.000Z",
					},
				})
				So(publicationState, ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "confirmed",
				})
			})

			Convey("As <upcoming, confirmed> when it is pushed back and brought forward", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-14T08:30:00.000Z",
				}
				publicationState := GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Pushed back",
						Date:         "2022-01-15T08:30:00.000Z",
					},
					{
						ChangeNotice: "Brought forward",
						Date:         "2022-01-16T08:30:00.000Z",
					},
				})
				So(publicationState, ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "confirmed",
				})
			})

			Convey("As <upcoming, confirmed> when the ReleaseData is an invalid timestamp", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "junk",
				}

				So(GetPublicationState(description, []model.DateChange{}), ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "confirmed",
				})

				So(GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Pushed back",
						Date:         "2022-01-15T08:30:00.000Z",
					},
				}), ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "confirmed",
				})
			})

			Convey("As <upcoming, confirmed> when the DateChanges contain an invalid timestamp", func() {
				description := model.ReleaseDescription{
					Cancelled:   false,
					Finalised:   true,
					Published:   false,
					ReleaseDate: "2022-01-15T08:30:00.000Z",
				}

				So(GetPublicationState(description, []model.DateChange{
					{
						ChangeNotice: "Pushed back",
						Date:         "junk",
					},
				}), ShouldResemble, model.PublicationState{
					Type:    "upcoming",
					SubType: "confirmed",
				})
			})
		})
	})
}

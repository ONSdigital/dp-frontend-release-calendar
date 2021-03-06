package mapper

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	"github.com/ONSdigital/log.go/v2/log"
)

// GetPublicationState returns the publication state
func GetPublicationState(description model.ReleaseDescription, dateChanges []model.DateChange) model.PublicationState {
	var state model.PublicationState

	switch {
	case description.Cancelled:
		state = model.PublicationState{
			Type: "cancelled",
		}
	case description.Published:
		state = model.PublicationState{
			Type: "published",
		}
	default:
		state = model.PublicationState{
			Type: "upcoming",
		}

		switch {
		case description.Finalised:
			state.SubType = "confirmed"
			if isPostponed(description.ReleaseDate, dateChanges) {
				state.SubType = "postponed"
			}
		default:
			state.SubType = "provisional"
		}
	}

	return state
}

func isPostponed(releaseDate string, dateChanges []model.DateChange) bool {
	parseTimestamp := func(timestamp string) (time.Time, error) {
		t, err := time.Parse(time.RFC3339, timestamp)

		if err != nil {
			log.Error(context.Background(), "failed to parse timestamp", err)
			return time.Time{}, err
		}

		return t, nil
	}

	totalDateChanges := len(dateChanges)

	if totalDateChanges > 0 {
		tReleaseDate, errReleaseDate := parseTimestamp(releaseDate)
		tLatestDateChange, errLatestDateChange := parseTimestamp(dateChanges[totalDateChanges-1].Date)
		return errReleaseDate == nil && errLatestDateChange == nil && tReleaseDate.After(tLatestDateChange)
	}

	return false
}

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// IsEndpointDeprecated takes a deprecation config and checks if an endpoint deprecation
// flag is on. If on and deprecation date has passed; it sets the headers and returns a 404
// with a message. It returns false if the deprecation date has not passed or there is an error.
func IsEndpointDeprecated(w http.ResponseWriter, r *http.Request, deprecationConfig config.Deprecation) bool {
	ctx := r.Context()

	if deprecationConfig.DeprecateEndpoint {
		now := time.Now().UTC()

		parsedDeprecation, err := parseTime(deprecationConfig.Deprecation)
		if err != nil {
			log.Error(ctx, "unable to parse deprecation date", err)
			return false
		}
		deprecationUnix := fmt.Sprintf("@%d", parsedDeprecation.Unix())

		parsedSunset, err := parseTime(deprecationConfig.Sunset)
		if err != nil {
			log.Error(ctx, "unable to parse sunset date", err)
			return false
		}

		if parsedSunset.Before(now) {
			w.Header().Set("content-type", "application/json")

			if deprecationConfig.Deprecation != "" {
				w.Header().Set("Deprecation", deprecationUnix)
			}

			if deprecationConfig.Sunset != "" {
				w.Header().Set("Sunset", parsedSunset.Format(time.RFC1123))
			}

			if deprecationConfig.Link != "" {
				w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"sunset\"", deprecationConfig.Link))
			}

			w.WriteHeader(http.StatusNotFound)

			if _, err = w.Write([]byte(deprecationConfig.DeprecationMessage)); err != nil {
				log.Error(ctx, "unable to write deprecation message", err)
				return false
			}
			return true
		}

		return true
	}

	return false
}

// parseTime takes a time string and parses it to the time it represents.
func parseTime(timeStr string) (time.Time, error) {
	for _, timeFmt := range []string{time.RFC3339, time.DateOnly, time.DateTime} {
		if parsedTime, err := time.Parse(timeFmt, timeStr); err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, errors.New("invalid time format")
}

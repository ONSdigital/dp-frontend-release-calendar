package queryparams

import (
	"errors"
	"strings"
)

type ReleaseType int

const (
	invalidReleaseType ReleaseType = iota
	Upcoming
	Published
	Cancelled
	Provisional
	Confirmed
	Postponed
)

var relTypeValues = map[ReleaseType]struct{ name, label string }{
	Upcoming:    {name: "type-upcoming", label: "Upcoming"},
	Published:   {name: "type-published", label: "Published"},
	Cancelled:   {name: "type-cancelled", label: "Cancelled"},
	Provisional: {name: "subtype-provisional", label: "Provisional"},
	Confirmed:   {name: "subtype-confirmed", label: "Confirmed"},
	Postponed:   {name: "subtype-postponed", label: "Postponed"},
}

func parseReleaseType(s string) (ReleaseType, error) {
	for rt, rtv := range relTypeValues {
		if strings.EqualFold(s, rtv.name) {
			return rt, nil
		}
	}

	return invalidReleaseType, errors.New("invalid release type string")
}

func (rt ReleaseType) Name() string {
	return relTypeValues[rt].name
}

func (rt ReleaseType) Label() string {
	return relTypeValues[rt].label
}

func (rt ReleaseType) String() string {
	return rt.Name()
}

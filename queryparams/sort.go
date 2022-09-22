package queryparams

import (
	"errors"
	"strings"
)

type Sort int

const (
	invalidSort Sort = iota
	RelDateAsc
	RelDateDesc
	TitleAZ
	TitleZA
	Relevance
)

var sortValues = map[Sort]struct{ feValue, beValue string }{
	RelDateAsc:  {feValue: "date-oldest", beValue: "release_date_asc"},
	RelDateDesc: {feValue: "date-newest", beValue: "release_date_desc"},
	TitleAZ:     {feValue: "alphabetical-az", beValue: "title_asc"},
	TitleZA:     {feValue: "alphabetical-za", beValue: "title_desc"},
	Relevance:   {feValue: "relevance", beValue: "relevance"},
}

func parseSort(sort string) (Sort, error) {
	for s, sv := range sortValues {
		if strings.EqualFold(sort, sv.feValue) {
			return s, nil
		}
	}

	return invalidSort, errors.New("invalid sort option string")
}

func (s Sort) String() string {
	return sortValues[s].feValue
}

func (s Sort) BackendString() string {
	return sortValues[s].beValue
}

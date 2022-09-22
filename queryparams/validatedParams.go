package queryparams

import (
	"net/url"
	"strconv"
)

type ValidatedParams struct {
	Limit       int
	Page        int
	Offset      int
	AfterDate   Date
	BeforeDate  Date
	Keywords    string
	Sort        Sort
	ReleaseType ReleaseType
	Provisional bool
	Confirmed   bool
	Postponed   bool
	Census      bool
	Highlight   bool
}

// AsBackendQuery converts to a url.Values object with parameters as expected by the api
func (vp ValidatedParams) AsBackendQuery() url.Values {
	return vp.asQuery(true)
}

// AsFrontendQuery converts to a url.Values object with parameters
func (vp ValidatedParams) AsFrontendQuery() url.Values {
	return vp.asQuery(false)
}

func (vp ValidatedParams) asQuery(isBackend bool) url.Values {
	var query = make(url.Values)
	setValue(query, Limit, strconv.Itoa(vp.Limit))
	setValue(query, Page, strconv.Itoa(vp.Page))

	if isBackend {
		setValue(query, Offset, strconv.Itoa(vp.Offset))
		setValue(query, Query, vp.Keywords)
		setValue(query, SortName, vp.getSortBackendString())
		setValue(query, DateFrom, vp.AfterDate.String())
		setValue(query, DateTo, vp.BeforeDate.String())
	} else {
		setValue(query, Keywords, vp.Keywords)
		setValue(query, SortName, vp.Sort.String())
		setValue(query, YearBefore, vp.BeforeDate.YearString())
		setValue(query, MonthBefore, vp.BeforeDate.MonthString())
		setValue(query, DayBefore, vp.BeforeDate.DayString())
		setValue(query, YearAfter, vp.AfterDate.YearString())
		setValue(query, MonthAfter, vp.AfterDate.MonthString())
		setValue(query, DayAfter, vp.AfterDate.DayString())
	}

	setValue(query, Type, vp.ReleaseType.String())
	if vp.ReleaseType == Upcoming {
		setBoolValue(query, Provisional.String(), vp.Provisional)
		setBoolValue(query, Confirmed.String(), vp.Confirmed)
		setBoolValue(query, Postponed.String(), vp.Postponed)
	}
	setBoolValue(query, Census, vp.Census)
	setBoolValue(query, Highlight, vp.Highlight)

	return query
}

func (vp ValidatedParams) getSortBackendString() string {
	// Newest is now defined as 'the closest date to today' and so its
	// meaning in terms of algorithmic definition (ascending/descending)
	// is reversed depending on the release-type that is being viewed
	if vp.ReleaseType == Upcoming && vp.Sort == RelDateDesc {
		return RelDateAsc.BackendString()
	} else if vp.ReleaseType == Upcoming && vp.Sort == RelDateAsc {
		return RelDateDesc.BackendString()
	}
	return vp.Sort.BackendString()
}

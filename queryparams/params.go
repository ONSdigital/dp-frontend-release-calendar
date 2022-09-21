package queryparams

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
)

const (
	Limit       = "limit"
	Page        = "page"
	Offset      = "offset"
	SortName    = "sort"
	DayBefore   = "before-day"
	DayAfter    = "after-day"
	MonthBefore = "before-month"
	MonthAfter  = "after-month"
	YearBefore  = "before-year"
	YearAfter   = "after-year"
	Keywords    = "keywords"
	Query       = "query"
	DateFrom    = "fromDate"
	DateTo      = "toDate"
	Type        = "release-type"
	Census      = "census"
	Highlight   = "highlight"
)

type intValidator func(valueAsString string) (int, error)

// getIntValidator returns an IntValidator object using the min and max values provided
func getIntValidator(minValue, maxValue int) intValidator {
	return func(valueAsString string) (int, error) {
		value, err := strconv.Atoi(valueAsString)
		if err != nil {
			return 0, fmt.Errorf("Value contains non numeric characters")
		}
		if value < minValue {
			return 0, fmt.Errorf("Value is below the minimum value (%d)", minValue)
		}
		if value > maxValue {
			return 0, fmt.Errorf("Value is above the maximum value (%d)", maxValue)
		}

		return value, nil
	}
}

var (
	dayValidator   = getIntValidator(1, 31)
	monthValidator = getIntValidator(1, 12)
	yearValidator  = getIntValidator(1900, 2150)
)

// GetLimit validates and returns the "limit" parameter
func GetLimit(ctx context.Context, params url.Values, defaultValue int, maxValue int) (int, error) {
	validator := getIntValidator(0, maxValue)
	return validateAndGetIntParam(ctx, params, Limit, defaultValue, validator)
}

// GetPage validates and returns the "page" parameter
func GetPage(ctx context.Context, params url.Values, maxPage int) (int, error) {
	defaultPage := 1
	validator := getIntValidator(1, maxPage)
	return validateAndGetIntParam(ctx, params, Page, defaultPage, validator)
}

func validateAndGetIntParam(ctx context.Context, params url.Values, paramName string, defaultValue int, validator intValidator) (int, error) {
	var (
		limit = defaultValue
		err   error
	)
	asString := params.Get(paramName)
	if asString != "" {
		limit, err = validator(asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": paramName, "value": asString})
			return 0, err
		}
	}

	return limit, nil
}

func GetSortOrder(ctx context.Context, params url.Values, defaultValue Sort) (Sort, error) {
	var (
		sort = defaultValue
		err  error
	)
	asString := params.Get(SortName)
	if asString != "" {
		sort, err = ParseSort(asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": Page, "value": asString})
			return Invalid, err
		}
	}

	// When keywords are empty in this case, force the sort order back to the default.
	if params.Get(Keywords) == "" && sort == Relevance {
		return defaultValue, nil
	}

	return sort, nil
}

func GetKeywords(_ context.Context, params url.Values, defaultValue string) (string, error) {
	keywords := defaultValue

	value := params.Get(Keywords)
	if value != "" {
		// Define any validation rules here. At present there are none, so we pass the given value directly
		keywords = value
	}

	return keywords, nil
}

func GetReleaseType(ctx context.Context, params url.Values, defaultValue ReleaseType) (ReleaseType, error) {
	var (
		relType = defaultValue
		err     error
	)
	asString := params.Get(Type)
	if asString != "" {
		relType, err = ParseReleaseType(asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": Type, "value": asString})
			return InvalidReleaseType, err
		}
	}

	return relType, nil
}

// GetBoolean finds a boolean parameter and returns a default value if not present
// It returns the default value together with an error if the value can't be parsed to a boolean
func GetBoolean(ctx context.Context, params url.Values, name string, defaultValue bool) (bool, error) {
	asString := params.Get(name)
	if asString == "" {
		return defaultValue, nil
	}

	upcoming, err := strconv.ParseBool(asString)
	if err != nil {
		log.Warn(ctx, "invalid boolean value", log.Data{"param": name, "value": asString})
		return defaultValue, fmt.Errorf("invalid boolean value for parameter %q", name)
	}

	return upcoming, nil
}

func DatesFromParams(ctx context.Context, params url.Values) (Date, Date, error) {
	var (
		from, to         time.Time
		fromDate, toDate Date
	)

	yearString, monthString, dayString := params.Get(YearAfter), params.Get(MonthAfter), params.Get(DayAfter)
	if yearString != "" && monthString != "" && dayString != "" {
		year, err := yearValidator(yearString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": YearAfter, "value": yearString})
			return Date{}, Date{}, err
		}
		month, err := monthValidator(monthString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": MonthAfter, "value": monthString})
			return Date{}, Date{}, err
		}
		day, err := dayValidator(dayString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": DayAfter, "value": dayString})
			return Date{}, Date{}, err
		}
		from = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		_, m, _ := from.Date()
		if m != time.Month(month) {
			log.Warn(ctx, "invalid day of month", log.Data{DayAfter: dayString, MonthAfter: monthString, YearAfter: yearString})
			return Date{}, Date{}, fmt.Errorf("invalid day (%s) of month (%s) in year (%s)", dayString, monthString, yearString)
		}
		fromDate = DateFromTime(from)
	}

	yearString, monthString, dayString = params.Get(YearBefore), params.Get(MonthBefore), params.Get(DayBefore)
	if yearString != "" && monthString != "" && dayString != "" {
		year, err := yearValidator(yearString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": YearBefore, "value": yearString})
			return Date{}, Date{}, err
		}
		month, err := monthValidator(monthString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": MonthBefore, "value": monthString})
			return Date{}, Date{}, err
		}
		day, err := dayValidator(dayString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": DayBefore, "value": dayString})
			return Date{}, Date{}, err
		}
		to = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		_, m, _ := to.Date()
		if m != time.Month(month) {
			log.Warn(ctx, "invalid day of month", log.Data{DayBefore: dayString, MonthBefore: monthString, YearBefore: yearString})
			return Date{}, Date{}, fmt.Errorf("invalid day (%s) of month (%s) in year (%s)", dayString, monthString, yearString)
		}
		toDate = DateFromTime(to)
	}

	if !from.IsZero() && !to.IsZero() {
		if from.After(to) {
			log.Warn(ctx, "invalid dates - from after to", log.Data{DateFrom: fromDate, DateTo: toDate})
			return Date{}, Date{}, errors.New("invalid dates - 'after' after 'before'")
		}
	}

	return fromDate, toDate, nil
}

// CalculateOffset returns the offset (0 based) into a list, given a page number (1 based) and the size of a page.
// A pageNumber <= 0 or a pageSize <= 0 will give an offset of 0
func CalculateOffset(pageNumber, pageSize int) int {
	if pageNumber <= 0 || pageSize <= 0 {
		return 0
	}

	return (pageNumber * pageSize) - pageSize
}

// CalculatePageNumber returns the page number (1 based) containing the offset(th) (0 based) element in a list, given a page size of pageSize.
// An offset <= 0 or pageSize <= 0 will give a page number of 1, i.e. the first page
func CalculatePageNumber(offset, pageSize int) int {
	if offset <= 0 || pageSize <= 0 {
		return 1
	}

	if (offset+1)%pageSize == 0 {
		return (offset + 1) / pageSize
	}

	return ((offset + 1) / pageSize) + 1
}

type Sort int

const (
	Invalid Sort = iota
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
	Invalid:     {feValue: "invalid", beValue: "invalid"},
}

func ParseSort(sort string) (Sort, error) {
	for s, sv := range sortValues {
		if strings.EqualFold(sort, sv.feValue) {
			return s, nil
		}
	}

	return Invalid, errors.New("invalid sort option string")
}

func MustParseSort(sort string) Sort {
	s, err := ParseSort(sort)
	if err != nil {
		panic("invalid sort string: " + sort)
	}

	return s
}

func (s Sort) String() string {
	return sortValues[s].feValue
}

func (s Sort) BackendString() string {
	return sortValues[s].beValue
}

type Date struct {
	date       time.Time
	y, m, d    int
	ys, ms, ds string
}

const DateFormat = "2006-01-02"

func MustParseDate(dateAsString string) Date {
	d, err := ParseDate(dateAsString)
	if err != nil {
		panic("invalid date string: " + dateAsString)
	}

	return d
}

func ParseDate(dateAsString string) (Date, error) {
	if dateAsString == "" {
		return Date{}, nil
	}
	t, err := time.Parse(DateFormat, dateAsString)
	if err != nil {
		return Date{}, err
	}

	return DateFromTime(t), nil
}

func DateFromTime(t time.Time) Date {
	if t.IsZero() {
		return Date{}
	}
	date := Date{date: t}
	y, m, d := t.Date()
	date.y, date.m, date.d = y, int(m), d
	date.ys, date.ms, date.ds = strconv.Itoa(y), strconv.Itoa(int(m)), strconv.Itoa(d)

	return date
}

func (d Date) String() string {
	if d.date.IsZero() {
		return ""
	}

	return d.date.UTC().Format(DateFormat)
}

func (d Date) YearString() string {
	return d.ys
}

func (d Date) MonthString() string {
	return d.ms
}

func (d Date) DayString() string {
	return d.ds
}

type ReleaseType int

const (
	InvalidReleaseType ReleaseType = iota
	Upcoming
	Published
	Cancelled
	Provisional
	Confirmed
	Postponed
)

var relTypeValues = map[ReleaseType]struct{ name, label string }{
	Upcoming:           {name: "type-upcoming", label: "Upcoming"},
	Published:          {name: "type-published", label: "Published"},
	Cancelled:          {name: "type-cancelled", label: "Cancelled"},
	Provisional:        {name: "subtype-provisional", label: "Provisional"},
	Confirmed:          {name: "subtype-confirmed", label: "Confirmed"},
	Postponed:          {name: "subtype-postponed", label: "Postponed"},
	InvalidReleaseType: {name: "Invalid", label: "Invalid"},
}

func ParseReleaseType(s string) (ReleaseType, error) {
	for rt, rtv := range relTypeValues {
		if strings.EqualFold(s, rtv.name) {
			return rt, nil
		}
	}

	return InvalidReleaseType, errors.New("invalid release type string")
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

	setValue(query, YearBefore, vp.BeforeDate.YearString())
	setValue(query, MonthBefore, vp.BeforeDate.MonthString())
	setValue(query, DayBefore, vp.BeforeDate.DayString())

	setValue(query, YearAfter, vp.AfterDate.YearString())
	setValue(query, MonthAfter, vp.AfterDate.MonthString())
	setValue(query, DayAfter, vp.AfterDate.DayString())

	if isBackend {
		setValue(query, Offset, strconv.Itoa(vp.Offset))
		setValue(query, Query, vp.Keywords)
		setValue(query, SortName, vp.getSortBackendString())
	} else {
		setValue(query, Keywords, vp.Keywords)
		setValue(query, SortName, vp.Sort.String())
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

func setValue(query url.Values, key string, value string) {
	if value != "" {
		query.Set(key, value)
	}
}

func setBoolValue(query url.Values, key string, value bool) {
	if value {
		query.Set(key, strconv.FormatBool(value))
	}
}

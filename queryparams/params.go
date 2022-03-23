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
	Published   = "type-published"
	Cancelled   = "type-cancelled"
	Upcoming    = "type-upcoming"
	Provisional = "subtype-provisional"
	Confirmed   = "subtype-confirmed"
	Postponed   = "subtype-postponed"
	Census      = "census"
)

func ParamGet(params url.Values, key string, defaultValue string) string {
	valueAsString := params.Get(key)
	if valueAsString == "" {
		return defaultValue
	}

	return valueAsString
}

type IntValidator func(name string, valueAsString string) (int, error)

func GetIntValidator(minValue, maxValue int) IntValidator {
	return func(name string, valueAsString string) (int, error) {
		value, err := strconv.Atoi(valueAsString)
		if err != nil {
			return 0, fmt.Errorf("%s search parameter provided with non numeric characters", name)
		}
		if value < minValue {
			return 0, fmt.Errorf("%s search parameter provided with a value that is below the minimum value", name)
		}
		if value > maxValue {
			return 0, fmt.Errorf("%s search parameter provided with a value that is above the maximum value", name)
		}

		return value, nil
	}
}

var (
	dayValidator   = GetIntValidator(1, 31)
	monthValidator = GetIntValidator(1, 12)
	yearValidator  = GetIntValidator(1900, 2150)
)

func GetLimit(ctx context.Context, params url.Values, defaultValue int, validator IntValidator) (int, error) {
	var (
		limit = defaultValue
		err   error
	)
	asString := params.Get(Limit)
	if asString != "" {
		limit, err = validator(Limit, asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": Limit, "value": asString})
			return 0, err
		}
	}

	return limit, nil
}

func GetPage(ctx context.Context, params url.Values, defaultValue int, validator IntValidator) (int, error) {
	var (
		limit = defaultValue
		err   error
	)
	asString := params.Get(Page)
	if asString != "" {
		limit, err = validator(Page, asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": Page, "value": asString})
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

func GetBoolean(ctx context.Context, params url.Values, name string, defaultValue bool) (bool, bool, error) {
	asString := params.Get(name)
	if asString == "" {
		return defaultValue, false, nil
	}

	upcoming, err := strconv.ParseBool(asString)
	if err != nil {
		log.Warn(ctx, fmt.Sprintf("invalid boolean value for %q", name), log.Data{"param": name, "value": asString})
		return false, false, fmt.Errorf("invalid boolean value for %q", name)
	}

	return upcoming, true, nil
}

func DatesFromParams(ctx context.Context, params url.Values) (Date, Date, error) {
	var (
		from, to         time.Time
		fromDate, toDate Date
	)

	yearString, monthString, dayString := params.Get(YearAfter), params.Get(MonthAfter), params.Get(DayAfter)
	if yearString != "" && monthString != "" && dayString != "" {
		year, err := yearValidator(YearAfter, yearString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": YearAfter, "value": yearString})
			return Date{}, Date{}, err
		}
		month, err := monthValidator(MonthAfter, monthString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": MonthAfter, "value": monthString})
			return Date{}, Date{}, err
		}
		day, err := dayValidator(DayAfter, dayString)
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
		year, err := yearValidator(YearBefore, yearString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": YearBefore, "value": yearString})
			return Date{}, Date{}, err
		}
		month, err := monthValidator(MonthBefore, monthString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": MonthBefore, "value": monthString})
			return Date{}, Date{}, err
		}
		day, err := dayValidator(DayBefore, dayString)
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
)

var sortNames = map[Sort]string{RelDateAsc: "release_date_asc", RelDateDesc: "release_date_desc", TitleAZ: "title_asc", TitleZA: "title_desc", Invalid: "invalid"}

func ParseSort(sort string) (Sort, error) {
	for s, sn := range sortNames {
		if strings.EqualFold(sort, sn) {
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
	return sortNames[s]
}

type SortOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

var SortOptions = []SortOption{
	{Label: "Newest", Value: "release_date_desc"},
	{Label: "Oldest", Value: "release_date_asc"},
	{Label: "TitleAZ", Value: "title_desc"},
	{Label: "TitleZA", Value: "title_asc"},
}

type Date struct {
	date       time.Time
	y, m, d    int
	ys, ms, ds string
}

const DateFormat = "2006-01-02"

func DateFromString(dateAsString string) (Date, error) {
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

type ValidatedParams struct {
	Limit       int
	Offset      int
	AfterDate   Date
	BeforeDate  Date
	Keywords    string
	Sort        Sort
	Published   bool
	Cancelled   bool
	Upcoming    bool
	Provisional bool
	Confirmed   bool
	Postponed   bool
	Census      bool
}
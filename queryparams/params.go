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
	Provisional = "subtype-provisional"
	Confirmed   = "subtype-confirmed"
	Postponed   = "subtype-postponed"
	Census      = "census"
	Highlight   = "highlight"
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

func GetBackwardsCompatibleReleaseType(ctx context.Context, params url.Values, defaultValue ReleaseType) (ReleaseType, error) {
	if params.Get("release-type") == "" {
		switch {
		case params.Get("type-upcoming") != "":
			return Upcoming, nil
		case params.Get("type-published") != "":
			return Published, nil
		default:
			return Cancelled, nil
		}
	}

	return GetReleaseType(ctx, params, defaultValue)
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
	Relevance
)

var feSortNames = map[Sort]string{RelDateAsc: "date-oldest", RelDateDesc: "date-newest", TitleAZ: "alphabetical-az", TitleZA: "alphabetical-za", Relevance: "relevance", Invalid: "invalid"}
var beSortNames = map[Sort]string{RelDateAsc: "release_date_asc", RelDateDesc: "release_date_desc", TitleAZ: "title_asc", TitleZA: "title_desc", Relevance: "relevance", Invalid: "invalid"}

func ParseSort(sort string) (Sort, error) {
	for s, sn := range feSortNames {
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
	return feSortNames[s]
}

func (s Sort) BackendString() string {
	return beSortNames[s]
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

func (d Date) Date() (int, int, int) {
	return d.y, d.m, d.d
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
)

var relTypeNames = map[ReleaseType]string{Upcoming: "type-upcoming", Published: "type-published", Cancelled: "type-cancelled", InvalidReleaseType: "Invalid"}

func ParseReleaseType(s string) (ReleaseType, error) {
	for rt, rtn := range relTypeNames {
		if strings.EqualFold(s, rtn) {
			return rt, nil
		}
	}

	return InvalidReleaseType, errors.New("invalid release type string")
}

func MustParseReleaseType(s string) ReleaseType {
	rt, err := ParseReleaseType(s)
	if err != nil {
		panic("invalid release type string: " + s)
	}

	return rt
}

func (rt ReleaseType) String() string {
	return relTypeNames[rt]
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

func (vp ValidatedParams) AsQuery() url.Values {
	var query = make(url.Values)
	query.Set(Limit, strconv.Itoa(vp.Limit))
	query.Set(Page, strconv.Itoa(vp.Page))

	query.Set(YearBefore, vp.BeforeDate.YearString())
	query.Set(MonthBefore, vp.BeforeDate.MonthString())
	query.Set(DayBefore, vp.BeforeDate.DayString())

	query.Set(YearAfter, vp.AfterDate.YearString())
	query.Set(MonthAfter, vp.AfterDate.MonthString())
	query.Set(DayAfter, vp.AfterDate.DayString())

	query.Set(Keywords, vp.Keywords)
	query.Set(SortName, vp.Sort.String())
	query.Set(Type, vp.ReleaseType.String())
	if vp.ReleaseType == Upcoming {
		query.Set(Provisional, strconv.FormatBool(vp.Provisional))
		query.Set(Confirmed, strconv.FormatBool(vp.Confirmed))
		query.Set(Postponed, strconv.FormatBool(vp.Postponed))
	}
	query.Set(Census, strconv.FormatBool(vp.Census))
	query.Set(Highlight, strconv.FormatBool(vp.Highlight))

	return query
}

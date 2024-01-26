package queryparams

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
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
			return 0, fmt.Errorf("value contains non numeric characters")
		}
		if value < minValue {
			return 0, fmt.Errorf("value is below the minimum value (%d)", minValue)
		}
		if value > maxValue {
			return 0, fmt.Errorf("value is above the maximum value (%d)", maxValue)
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
func GetLimit(ctx context.Context, params url.Values, defaultValue, maxValue int) (int, error) {
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
			return 0, fmt.Errorf("invalid %s parameter: %s", paramName, err.Error())
		}
	}

	return limit, nil
}

// GetSortOrder validates and returns the "sort" parameter
func GetSortOrder(ctx context.Context, params url.Values, defaultValue string) (Sort, error) {
	defaultSort, err := parseSort(defaultValue)
	if err != nil {
		log.Warn(ctx, fmt.Sprintf("Invalid config value for default sort. Using %s as default", RelDateDesc.String()), log.Data{"value": defaultValue})
		defaultSort = RelDateDesc
	}

	sort := defaultSort
	asString := params.Get(SortName)
	if asString != "" {
		sort, err = parseSort(asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": SortName, "value": asString})
			return defaultSort, fmt.Errorf("invalid %s parameter: %s", SortName, err.Error())
		}
	}

	// When keywords are empty in this case, force the sort order back to the default.
	if params.Get(Keywords) == "" && sort == Relevance {
		sort = defaultSort
	}

	return sort, nil
}

// GetKeywords validates and returns the "keywords" parameter
func GetKeywords(_ context.Context, params url.Values, defaultValue string) (string, error) {
	keywords := defaultValue

	value := params.Get(Keywords)
	if value != "" {
		// Define any validation rules here. At present there are none, so we pass the given value directly
		keywords = value
	}

	return keywords, nil
}

// GetReleaseType validates and returns the "release-type" parameter
func GetReleaseType(ctx context.Context, params url.Values, defaultValue ReleaseType) (ReleaseType, error) {
	var (
		relType = defaultValue
		err     error
	)
	asString := params.Get(Type)
	if asString != "" {
		relType, err = parseReleaseType(asString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": Type, "value": asString})
			return defaultValue, fmt.Errorf("invalid %s parameter: %s", Type, err.Error())
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

// ErrInvalidDateInput is return when input date is invalid e.g. 31 Feb
type ErrInvalidDateInput struct {
	msg string
}

func (e ErrInvalidDateInput) Error() string { return e.msg }

// Assumed indicates whether values are assumed
type Assumed struct {
	isDay, isMonth bool
}

// GetDates finds the date from and date to parameters
func GetDates(ctx context.Context, params url.Values) (startDate, endDate Date, err error) {
	var (
		startTime, endTime               time.Time
		assumedStartTime, assumedEndTime Assumed
	)

	yearAfterString, monthAfterString, dayAfterString := params.Get(YearAfter), params.Get(MonthAfter), params.Get(DayAfter)
	yearBeforeString, monthBeforeString, dayBeforeString := params.Get(YearBefore), params.Get(MonthBefore), params.Get(DayBefore)
	logData := log.Data{
		"year_after": yearAfterString, "month_after": monthAfterString, "day_after": dayAfterString,
		"year_before": yearBeforeString, "month_before": monthBeforeString, "day_before": DayBefore,
	}

	startTime, assumedStartTime, err = getValidTimestamp(yearAfterString, monthAfterString, dayAfterString)
	if err != nil {
		log.Warn(ctx, "invalid date, startDate", log.FormatErrors([]error{err}), logData)
		return Date{}, Date{}, err
	}

	startDate = DateFromTime(startTime)
	startDate.assumedDay = assumedStartTime.isDay
	startDate.assumedMonth = assumedStartTime.isMonth

	endTime, assumedEndTime, err = getValidTimestamp(yearBeforeString, monthBeforeString, dayBeforeString)
	if err != nil {
		log.Warn(ctx, "invalid date, endDate", log.FormatErrors([]error{err}), logData)
		return Date{}, Date{}, err
	}

	endDate = DateFromTime(endTime)
	endDate.assumedDay = assumedEndTime.isDay
	endDate.assumedMonth = assumedEndTime.isMonth

	if !startTime.IsZero() && !endTime.IsZero() && startTime.After(endTime) {
		log.Warn(ctx, "invalid date range: start date after end date", log.Data{DateFrom: startDate, DateTo: endDate})
		return Date{}, Date{}, errors.New("invalid dates: start date after end date")
	}

	return startDate, endDate, nil
}

func getValidTimestamp(year, month, day string) (time.Time, Assumed, error) {
	if (month != "" || day != "") && year == "" {
		return time.Time{}, Assumed{}, ErrInvalidDateInput{msg: "Enter a year"}
	}

	if year == "" {
		return time.Time{}, Assumed{}, nil
	}

	var aTimes Assumed
	if year != "" && month == "" {
		month = "1"
		aTimes.isMonth = true
	}

	if year != "" && day == "" {
		day = "1"
		aTimes.isDay = true
	}

	y, err := yearValidator(year)
	if err != nil {
		return time.Time{}, Assumed{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", year, err.Error())}
	}

	m, err := monthValidator(month)
	if err != nil {
		return time.Time{}, Assumed{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", month, err.Error())}
	}

	d, err := dayValidator(day)
	if err != nil {
		return time.Time{}, Assumed{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", day, err.Error())}
	}

	timestamp := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)

	// Check the day is valid for the month in the year, e.g. day 30 cannot be in month 2 (February)
	_, mo, _ := timestamp.Date()
	if mo != time.Month(m) {
		return time.Time{}, Assumed{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid day (%s) of month (%s) in year (%s)", day, month, year)}
	}

	return timestamp, aTimes, nil
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

func setValue(query url.Values, key, value string) {
	if value != "" {
		query.Set(key, value)
	}
}

func setBoolValue(query url.Values, key string, value bool) {
	if value {
		query.Set(key, strconv.FormatBool(value))
	}
}

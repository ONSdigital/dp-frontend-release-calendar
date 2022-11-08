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

// GetDates finds the date from and date to parameters
func GetDates(ctx context.Context, params url.Values) (Date, Date, error) {
	var (
		from, to         time.Time
		fromDate, toDate Date
	)

	yearString, monthString, dayString := params.Get(YearAfter), params.Get(MonthAfter), params.Get(DayAfter)
	if yearString != "" && monthString != "" && dayString != "" {
		year, err := yearValidator(yearString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": YearAfter, "value": yearString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", YearAfter, err.Error())}
		}
		month, err := monthValidator(monthString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": MonthAfter, "value": monthString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", MonthAfter, err.Error())}
		}
		day, err := dayValidator(dayString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": DayAfter, "value": dayString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", DayAfter, err.Error())}
		}
		from = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		_, m, _ := from.Date()
		if m != time.Month(month) {
			log.Warn(ctx, "invalid day of month", log.Data{DayAfter: dayString, MonthAfter: monthString, YearAfter: yearString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid day (%s) of month (%s) in year (%s)", dayString, monthString, yearString)}
		}
		fromDate = DateFromTime(from)
	}

	yearString, monthString, dayString = params.Get(YearBefore), params.Get(MonthBefore), params.Get(DayBefore)
	if yearString != "" && monthString != "" && dayString != "" {
		year, err := yearValidator(yearString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": YearBefore, "value": yearString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", YearBefore, err.Error())}
		}
		month, err := monthValidator(monthString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": MonthBefore, "value": monthString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", MonthBefore, err.Error())}
		}
		day, err := dayValidator(dayString)
		if err != nil {
			log.Warn(ctx, err.Error(), log.Data{"param": DayBefore, "value": dayString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid %s parameter: %s", DayBefore, err.Error())}
		}
		to = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		_, m, _ := to.Date()
		if m != time.Month(month) {
			log.Warn(ctx, "invalid day of month", log.Data{DayBefore: dayString, MonthBefore: monthString, YearBefore: yearString})
			return Date{}, Date{}, ErrInvalidDateInput{msg: fmt.Sprintf("invalid day (%s) of month (%s) in year (%s)", dayString, monthString, yearString)}
		}
		toDate = DateFromTime(to)
	}

	if !from.IsZero() && !to.IsZero() && from.After(to) {
		log.Warn(ctx, "invalid dates: from after to", log.Data{DateFrom: fromDate, DateTo: toDate})
		return Date{}, Date{}, errors.New("invalid dates: 'after' after 'before'")
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

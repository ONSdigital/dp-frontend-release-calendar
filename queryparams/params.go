package queryparams

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	core "github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	Limit       = "limit"
	Page        = "page"
	Offset      = "offset"
	SortName    = "sort"
	DayBefore   = "before-day"
	DayAfter    = "after-day"
	Before      = "before"
	MonthBefore = Before + "-month"
	After       = "after"
	MonthAfter  = After + "-month"
	YearBefore  = "before-year"
	YearAfter   = "after-year"
	Keywords    = "keywords"
	Query       = "query"
	DateFrom    = "fromDate"
	DateFromErr = DateFrom + "-error"
	DateTo      = "toDate"
	DateToErr   = DateTo + "-error"
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
			return 0, fmt.Errorf("enter a number")
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

// GetStartDate returns the validated date from parameters
func GetStartDate(params url.Values) (startDate Date, validationErrs []core.ErrorItem) {
	var startTime time.Time

	yearAfterString, monthAfterString, dayAfterString := params.Get(YearAfter), params.Get(MonthAfter), params.Get(DayAfter)
	startDate.ds = dayAfterString
	startDate.ms = monthAfterString
	startDate.ys = yearAfterString

	if (monthAfterString != "" || dayAfterString != "") && yearAfterString == "" {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: "Enter a released after year",
			},
			ID:  DateFromErr,
			URL: fmt.Sprintf("#%s", DateFromErr),
		})
		startDate.hasValidationErr = true
		return startDate, validationErrs
	}

	var assumedDay, assumedMonth bool
	if yearAfterString != "" && monthAfterString == "" {
		monthAfterString = "1"
		assumedMonth = true
	}

	if yearAfterString != "" && dayAfterString == "" {
		dayAfterString = "1"
		assumedDay = true
	}

	startTime, validationErrs = getValidTimestamp(yearAfterString, monthAfterString, dayAfterString, DateFromErr, After)
	if len(validationErrs) > 0 {
		startDate.hasValidationErr = true
		return startDate, validationErrs
	}

	startDate = DateFromTime(startTime)
	startDate.assumedDay = assumedDay
	startDate.assumedMonth = assumedMonth
	startDate.hasValidationErr = false

	return startDate, nil
}

// GetDates returns the validated date to parameters
func GetEndDate(params url.Values) (endDate Date, validationErrs []core.ErrorItem) {
	var endTime time.Time

	yearBeforeString, monthBeforeString, dayBeforeString := params.Get(YearBefore), params.Get(MonthBefore), params.Get(DayBefore)
	endDate.ds = dayBeforeString
	endDate.ms = monthBeforeString
	endDate.ys = yearBeforeString

	if (monthBeforeString != "" || dayBeforeString != "") && yearBeforeString == "" {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: "Enter a released before year",
			},
			ID:  DateToErr,
			URL: fmt.Sprintf("#%s", DateToErr),
		})
		endDate.hasValidationErr = true
		return endDate, validationErrs
	}

	var assumedDay, assumedMonth bool
	if yearBeforeString != "" && monthBeforeString == "" {
		monthBeforeString = "1"
		assumedMonth = true
	}

	if yearBeforeString != "" && dayBeforeString == "" {
		dayBeforeString = "1"
		assumedDay = true
	}

	endTime, validationErrs = getValidTimestamp(yearBeforeString, monthBeforeString, dayBeforeString, DateToErr, Before)
	if len(validationErrs) > 0 {
		endDate.hasValidationErr = true
		return endDate, validationErrs
	}

	endDate = DateFromTime(endTime)
	endDate.assumedDay = assumedDay
	endDate.assumedMonth = assumedMonth
	endDate.hasValidationErr = false

	return endDate, nil
}

// getValidTimestamp returns a valid timestamp or an error
func getValidTimestamp(year, month, day, fieldsetID, fieldsetStr string) (time.Time, []core.ErrorItem) {
	if year == "" || month == "" || day == "" {
		return time.Time{}, []core.ErrorItem{}
	}

	var validationErrs []core.ErrorItem

	y, err := yearValidator(year)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: fmt.Sprintf("%s for released %s year", CapitalizeFirstLetter(err.Error()), fieldsetStr),
			},
			ID:  fieldsetID,
			URL: fmt.Sprintf("#%s", fieldsetID),
		})
	}

	m, err := monthValidator(month)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: fmt.Sprintf("%s for released %s month", CapitalizeFirstLetter(err.Error()), fieldsetStr),
			},
			ID:  fieldsetID,
			URL: fmt.Sprintf("#%s", fieldsetID),
		})
	}

	d, err := dayValidator(day)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: fmt.Sprintf("%s for released %s day", CapitalizeFirstLetter(err.Error()), fieldsetStr),
			},
			ID:  fieldsetID,
			URL: fmt.Sprintf("#%s", fieldsetID),
		})
	}

	// Throw errors back to user before further validation
	if len(validationErrs) > 0 {
		return time.Time{}, validationErrs
	}

	timestamp := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)

	// Check the day is valid for the month in the year, e.g. day 30 cannot be in month 2 (February)
	_, mo, _ := timestamp.Date()
	if mo != time.Month(m) {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: "Enter a valid date",
			},
			ID:  fieldsetID,
			URL: fmt.Sprintf("#%s", fieldsetID),
		})
	}

	return timestamp, validationErrs
}

// CapitalizeFirstLetter is a helper function that transforms the first letter of a string to uppercase
func CapitalizeFirstLetter(input string) string {
	switch {
	case len(input) <= 0:
		return input
	case len(input) == 1:
		return strings.ToUpper(input)
	default:
		return strings.ToUpper(input[:1]) + input[1:]
	}
}

// ValidateDateRange returns an error and 'to' date if the 'from' date is after than the 'to' date
func ValidateDateRange(from, to Date) (end Date, err error) {
	startDate, err := ParseDate(from.String())
	if err != nil {
		return Date{}, err
	}
	endDate, err := ParseDate(to.String())
	if err != nil {
		return Date{}, err
	}

	startTime, _ := getValidTimestamp(startDate.YearString(), startDate.MonthString(), startDate.DayString(), "", "")
	endTime, _ := getValidTimestamp(endDate.YearString(), endDate.MonthString(), endDate.DayString(), "", "")
	if startTime.After(endTime) {
		end = to
		end.hasValidationErr = true
		return end, fmt.Errorf("enter a released before year that is later than %s", startDate.YearString())
	}
	return Date{}, nil
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

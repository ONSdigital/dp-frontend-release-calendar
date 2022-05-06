package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	search "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mapper"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"

	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
)

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// Release will load a release page
func Release(cfg config.Config, rc RenderClient, api ReleaseCalendarAPI) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		release(w, r, accessToken, collectionID, lang, rc, api, cfg)
	})
}

func release(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, api ReleaseCalendarAPI, cfg config.Config) {
	ctx := req.Context()

	uri := strings.TrimPrefix(req.URL.EscapedPath(), cfg.PrivateRoutingPrefix)
	release, err := api.GetLegacyRelease(ctx, userAccessToken, collectionID, lang, uri)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateRelease(basePage, *release)

	rc.BuildPage(w, m, "release")
}

func ReleaseCalendar(cfg config.Config, rc RenderClient, api SearchAPI) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		releaseCalendar(w, r, accessToken, collectionID, lang, rc, api, cfg)
	})
}

func releaseCalendar(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, api SearchAPI, cfg config.Config) {
	ctx := req.Context()
	params := req.URL.Query()
	validatedParams := queryparams.ValidatedParams{}

	pageSize, err := queryparams.GetLimit(ctx, params, cfg.DefaultLimit, queryparams.GetIntValidator(0, cfg.DefaultMaximumLimit))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid %s parameter", queryparams.Limit), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.Limit, strconv.Itoa(pageSize))
	validatedParams.Limit = pageSize

	pageNumber, err := queryparams.GetPage(ctx, params, 1, queryparams.GetIntValidator(1, cfg.DefaultMaximumSearchResults/cfg.DefaultLimit))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid %s parameter", queryparams.Page), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.Page, strconv.Itoa(pageNumber))
	validatedParams.Page = pageNumber
	offset := queryparams.CalculateOffset(pageNumber, pageSize)
	params.Set(queryparams.Offset, strconv.Itoa(offset))
	validatedParams.Offset = offset

	fromDate, toDate, err := queryparams.DatesFromParams(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.DateFrom, fromDate.String())
	validatedParams.AfterDate = fromDate
	params.Set(queryparams.DateTo, toDate.String())
	validatedParams.BeforeDate = toDate

	sort, err := queryparams.GetSortOrder(ctx, params, queryparams.MustParseSort(cfg.DefaultSort))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.SortName, sort.BackendString())
	validatedParams.Sort = sort

	keywords, err := queryparams.GetKeywords(ctx, params, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.Keywords, keywords)
	validatedParams.Keywords = keywords
	params.Set(queryparams.Query, keywords)

	releaseType, err := queryparams.GetReleaseType(ctx, params, queryparams.Published)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.Type, releaseType.String())
	validatedParams.ReleaseType = releaseType

	provisional, set, err := queryparams.GetBoolean(ctx, params, queryparams.Provisional.String(), false)
	validatedParams.Provisional = provisional
	if provisional || set {
		params.Set(queryparams.Provisional.String(), strconv.FormatBool(provisional))
	}
	confirmed, set, err := queryparams.GetBoolean(ctx, params, queryparams.Confirmed.String(), false)
	validatedParams.Confirmed = confirmed
	if confirmed || set {
		params.Set(queryparams.Confirmed.String(), strconv.FormatBool(confirmed))
	}
	postponed, set, err := queryparams.GetBoolean(ctx, params, queryparams.Postponed.String(), false)
	validatedParams.Postponed = postponed
	if postponed || set {
		params.Set(queryparams.Postponed.String(), strconv.FormatBool(postponed))
	}

	census, set, err := queryparams.GetBoolean(ctx, params, queryparams.Census, false)
	validatedParams.Census = census
	if census || set {
		params.Set(queryparams.Census, strconv.FormatBool(census))
	}

	highlight, set, err := queryparams.GetBoolean(ctx, params, queryparams.Highlight, true)
	validatedParams.Highlight = highlight
	if highlight || set {
		params.Set(queryparams.Highlight, strconv.FormatBool(highlight))
	}

	releases, err := api.GetReleases(ctx, userAccessToken, collectionID, lang, params)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	calendar := mapper.CreateReleaseCalendar(basePage, validatedParams, releases, cfg)

	rc.BuildPage(w, calendar, "calendar")
}

func ReleaseCalendarICSEntries(cfg config.Config, api SearchAPI) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		releaseCalendarICSEntries(w, r, accessToken, collectionID, lang, api, cfg)
	})
}

func releaseCalendarICSEntries(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, api SearchAPI, cfg config.Config) {
	ctx := req.Context()
	params := req.URL.Query()

	params.Set(queryparams.Limit, strconv.Itoa(cfg.DefaultMaximumSearchResults))
	params.Set(queryparams.SortName, queryparams.RelDateAsc.BackendString())
	params.Set(queryparams.DateTo, time.Now().AddDate(0, 3, 0).Format(queryparams.DateFormat))
	params.Set(queryparams.Type, queryparams.Upcoming.String())

	releases, err := api.GetReleases(ctx, userAccessToken, collectionID, lang, params)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	fileWriter := new(bytes.Buffer)
	err = toICSFile(ctx, releases.Releases, fileWriter)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Character-Encoding", "UTF8")
	w.Header().Set("Content-Disposition", "attachment; filename=releases.ics")
	if _, err = w.Write(fileWriter.Bytes()); err != nil {
		setStatusCode(req, w, err)
		return
	}

}

func toICSFile(ctx context.Context, releases []search.Release, w io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	printLine := func(s string) {
		if _, e := fmt.Fprintln(w, s); e != nil {
			panic(e)
		}
	}
	printLine("BEGIN:VCALENDAR")
	printLine("PRODID:-//Office for National Statistics//NONSGML//EN")
	printLine("VERSION:2.0")
	printLine("CALSCALE:GREGORIAN")
	for _, r := range releases {
		printLine("BEGIN:VEVENT")
		printLine("DTSTAMP:" + time.Now().UTC().Format(iCalDateFormat))
		releaseDate := iCalDate(ctx, r.Description.ReleaseDate)
		printLine("DTSTART:" + releaseDate)
		printLine("DTEND:" + releaseDate)
		printLine("SUMMARY:" + r.Description.Title)
		printLine("UID:" + r.URI)
		printLine("STATUS:" + releaseStatus(r))
		printLine("DESCRIPTION:" + r.Description.Summary)
		printLine("END:VEVENT")
	}
	printLine("END:VCALENDAR")

	return nil
}

const iCalDateFormat = "20060102T150405Z"

func iCalDate(ctx context.Context, dateRFC3339 string) string {
	dateiCal, err := time.Parse(time.RFC3339, dateRFC3339)
	if err != nil {
		log.Warn(ctx, "iCalDate::unrecognised date format", log.Data{"date": dateRFC3339, "error": err})
		return ""
	}

	return dateiCal.UTC().Format(iCalDateFormat)
}

func releaseStatus(r search.Release) string {
	switch {
	case r.Description.Cancelled:
		return queryparams.Cancelled.Label()
	case r.Description.Published:
		return queryparams.Published.Label()
	case r.Description.Finalised:
		if r.DateChanges != nil {
			return queryparams.Postponed.Label()
		}
		return queryparams.Confirmed.Label()
	default:
		return queryparams.Provisional.Label()
	}
}

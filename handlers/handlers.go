package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

func release(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, api ReleaseCalendarAPI, _ config.Config) {
	ctx := req.Context()

	release, err := api.GetLegacyRelease(ctx, userAccessToken, collectionID, lang, req.URL.EscapedPath())
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateRelease(basePage, *release)

	rc.BuildPage(w, m, "release")
}

func PreviousReleasesSample(cfg config.Config, rc RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		previousReleasesSample(w, req, rc, cfg)
	}
}

func previousReleasesSample(w http.ResponseWriter, req *http.Request, rc RenderClient, cfg config.Config) {
	ctx := req.Context()
	basePage := rc.NewBasePageModel()
	m := mapper.CreatePreviousReleases(ctx, basePage, cfg)

	rc.BuildPage(w, m, "previousreleases")
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
	params.Set(queryparams.SortName, sort.String())
	validatedParams.Sort = sort

	keywords, err := queryparams.GetKeywords(ctx, params, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.Set(queryparams.Keywords, keywords)
	validatedParams.Keywords = keywords
	params.Set(queryparams.Query, keywords)

	// TODO Upcoming is the only Release Type to be parsed as present until the extended calendar query is added
	upcoming, set, err := queryparams.GetBoolean(ctx, params, queryparams.Upcoming, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validatedParams.Upcoming = upcoming
	if upcoming || set {
		params.Set(queryparams.Upcoming, strconv.FormatBool(upcoming))
	}

	releases, err := api.GetReleases(ctx, userAccessToken, collectionID, lang, params)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	calendar := mapper.CreateReleaseCalendar(basePage, validatedParams, releases)

	rc.BuildPage(w, calendar, "calendar")
}

func CalendarSample(cfg config.Config, rc RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		calendarSample(w, req, rc, cfg)
	}
}

func calendarSample(w http.ResponseWriter, req *http.Request, rc RenderClient, cfg config.Config) {
	ctx := req.Context()
	basePage := rc.NewBasePageModel()
	m := mapper.CreateCalendar(ctx, basePage, cfg)

	rc.BuildPage(w, m, "calendar")
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	search "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mapper"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	core "github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/gorilla/feeds"

	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	iCalDateFormat = "20060102T150405Z"
	defaultMaxAge  = 5 // 5 seconds
	homepagePath   = "/"
)

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		status = err.Code()
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

func setCacheHeader(ctx context.Context, w http.ResponseWriter, babbage BabbageAPI, uri, key string) {
	maxAge, err := babbage.GetMaxAge(ctx, uri, key)
	if err != nil {
		// Do not cache
		maxAge = defaultMaxAge
		log.Warn(ctx,
			fmt.Sprintf("Couldn't find max age from Babbage, using default %d sec", maxAge),
			log.Data{"uri": uri, "err": err.Error()})
	}
	w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
}

// Release will load a release page
func Release(cfg config.Config, rc RenderClient, api ReleaseCalendarAPI, babbage BabbageAPI, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		ctx := r.Context()
		releaseURI := strings.TrimPrefix(r.URL.EscapedPath(), cfg.RoutingPrefix)

		homepageContent, err := zc.GetHomepageContent(ctx, accessToken, collectionID, lang, homepagePath)
		if err != nil {
			log.Warn(ctx, "unable to get homepage content", log.FormatErrors([]error{err}), log.Data{"homepage_content": err})
		}

		release, err := api.GetLegacyRelease(ctx, accessToken, collectionID, lang, releaseURI)
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		basePage := rc.NewBasePageModel()
		m := mapper.CreateRelease(basePage, *release, lang, cfg.CalendarPath(), homepageContent.ServiceMessage, homepageContent.EmergencyBanner)

		setCacheHeader(ctx, w, babbage, releaseURI, cfg.BabbageMaxAgeKey)

		rc.BuildPage(w, m, "release")
	})
}

func ReleaseData(cfg config.Config, api ReleaseCalendarAPI) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		release, err := api.GetLegacyRelease(r.Context(), accessToken, collectionID, lang, strings.TrimSuffix(strings.TrimPrefix(r.URL.EscapedPath(), cfg.RoutingPrefix), "/data"))
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		data, err := json.Marshal(release)
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		w.Header().Set("content-type", "application/json")
		if _, err = w.Write(data); err != nil {
			setStatusCode(r, w, err)
			return
		}
	})
}

func ReleaseCalendar(cfg config.Config, rc RenderClient, api SearchAPI, babbage BabbageAPI, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		var err error
		ctx := r.Context()
		params := r.URL.Query()

		homepageContent, err := zc.GetHomepageContent(ctx, accessToken, collectionID, lang, homepagePath)
		if err != nil {
			log.Warn(ctx, "unable to get homepage content", log.FormatErrors([]error{err}), log.Data{"homepage_content": err})
		}

		validatedParams, validationErrs := validateParamsAsFrontend(ctx, params, cfg)
		if len(validationErrs) > 0 {
			calendar := mapper.CreateReleaseCalendar(rc.NewBasePageModel(), validatedParams, search.ReleaseResponse{}, cfg, lang, homepageContent.ServiceMessage, homepageContent.EmergencyBanner, validationErrs)
			setCacheHeader(ctx, w, babbage, "/releasecalendar", cfg.BabbageMaxAgeKey)
			rc.BuildPage(w, calendar, "calendar")
			return
		}

		if _, rssParam := params["rss"]; rssParam {
			r.Header.Set("Accept", "application/rss+xml")
			if err = createRSSFeed(ctx, w, r, lang, collectionID, accessToken, api, validatedParams); err != nil {
				setStatusCode(r, w, err)
				return
			}
			return
		}

		releases, err := api.GetReleases(ctx, accessToken, collectionID, lang, validatedParams.AsBackendQuery())
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		calendar := mapper.CreateReleaseCalendar(rc.NewBasePageModel(), validatedParams, releases, cfg, lang, homepageContent.ServiceMessage, homepageContent.EmergencyBanner, nil)
		setCacheHeader(ctx, w, babbage, "/releasecalendar", cfg.BabbageMaxAgeKey)
		rc.BuildPage(w, calendar, "calendar")
	})
}

func ReleaseCalendarData(cfg config.Config, api SearchAPI) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		ctx := r.Context()
		params := r.URL.Query()

		validatedParams, err := validateParams(ctx, params, cfg)
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		releases, err := api.GetReleases(ctx, accessToken, collectionID, lang, validatedParams.AsBackendQuery())
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		data, err := json.Marshal(releases)
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		w.Header().Set("content-type", "application/json")
		if _, err = w.Write(data); err != nil {
			setStatusCode(r, w, err)
			return
		}
	})
}

func validateParamsAsFrontend(ctx context.Context, params url.Values, cfg config.Config) (vp queryparams.ValidatedParams, validationErrs []core.ErrorItem) {
	validatedParams := queryparams.ValidatedParams{}

	limit, err := queryparams.GetLimit(ctx, params, cfg.DefaultLimit, cfg.DefaultMaximumLimit)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: err.Error(),
			},
		})
	}
	validatedParams.Limit = limit

	pageNumber, err := queryparams.GetPage(ctx, params, cfg.DefaultMaximumSearchResults/cfg.DefaultLimit)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: err.Error(),
			},
		})
	}
	validatedParams.Page = pageNumber

	validatedParams.Offset = queryparams.CalculateOffset(pageNumber, limit)

	fromDate, vErrs := queryparams.GetStartDate(params)
	if len(vErrs) > 0 {
		validationErrs = append(validationErrs, vErrs...)
	}
	validatedParams.AfterDate = fromDate

	toDate, vErrs := queryparams.GetEndDate(params)
	if len(vErrs) > 0 {
		validationErrs = append(validationErrs, vErrs...)
	}
	if fromDate.String() != "" && toDate.String() != "" {
		toDate, err = queryparams.ValidateDateRange(fromDate, toDate)
		if err != nil {
			validationErrs = append(validationErrs, core.ErrorItem{
				Description: core.Localisation{
					Text: queryparams.CapitalizeFirstLetter(err.Error()),
				},
				ID:  queryparams.DateToErr,
				URL: fmt.Sprintf("#%s", queryparams.DateToErr),
			})
		}
	}
	validatedParams.BeforeDate = toDate

	sort, err := queryparams.GetSortOrder(ctx, params, cfg.DefaultSort)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: err.Error(),
			},
		})
	}
	validatedParams.Sort = sort

	keywords, err := queryparams.GetKeywords(ctx, params, "")
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: err.Error(),
			},
		})
	}
	validatedParams.Keywords = keywords

	releaseType, err := queryparams.GetReleaseType(ctx, params, queryparams.Published)
	if err != nil {
		validationErrs = append(validationErrs, core.ErrorItem{
			Description: core.Localisation{
				Text: err.Error(),
			},
		})
	}
	validatedParams.ReleaseType = releaseType

	validatedParams.Provisional, _ = queryparams.GetBoolean(ctx, params, queryparams.Provisional.String(), false)
	validatedParams.Confirmed, _ = queryparams.GetBoolean(ctx, params, queryparams.Confirmed.String(), false)
	validatedParams.Postponed, _ = queryparams.GetBoolean(ctx, params, queryparams.Postponed.String(), false)
	validatedParams.Census, _ = queryparams.GetBoolean(ctx, params, queryparams.Census, false)
	validatedParams.Highlight, _ = queryparams.GetBoolean(ctx, params, queryparams.Highlight, true)

	return validatedParams, validationErrs
}

func validateParams(ctx context.Context, params url.Values, cfg config.Config) (queryparams.ValidatedParams, error) {
	validatedParams := queryparams.ValidatedParams{}

	limit, err := queryparams.GetLimit(ctx, params, cfg.DefaultLimit, cfg.DefaultMaximumLimit)
	if err != nil {
		return validatedParams, &clientErr{err}
	}
	validatedParams.Limit = limit

	pageNumber, err := queryparams.GetPage(ctx, params, cfg.DefaultMaximumSearchResults/cfg.DefaultLimit)
	if err != nil {
		return validatedParams, &clientErr{err}
	}
	validatedParams.Page = pageNumber

	validatedParams.Offset = queryparams.CalculateOffset(pageNumber, limit)

	fromDate, vErrs := queryparams.GetStartDate(params)
	if len(vErrs) > 0 {
		for _, err := range vErrs {
			log.Error(ctx, "invalid date", fmt.Errorf("startdate field error: %s", err.Description.Text))
		}
		return validatedParams, &clientErr{fmt.Errorf("invalid startDate")}
	}
	validatedParams.AfterDate = fromDate

	toDate, vErrs := queryparams.GetEndDate(params)
	if len(vErrs) > 0 {
		for _, err := range vErrs {
			log.Error(ctx, "invalid date", fmt.Errorf("endDate field error: %s", err.Description.Text))
		}
		return validatedParams, &clientErr{fmt.Errorf("invalid endDate")}
	}
	validatedParams.BeforeDate = toDate

	if fromDate.String() != "" && toDate.String() != "" {
		_, err = queryparams.ValidateDateRange(fromDate, toDate)
		if err != nil {
			return validatedParams, &clientErr{err}
		}
	}

	sort, err := queryparams.GetSortOrder(ctx, params, cfg.DefaultSort)
	if err != nil {
		return validatedParams, &clientErr{err}
	}
	validatedParams.Sort = sort

	keywords, err := queryparams.GetKeywords(ctx, params, "")
	if err != nil {
		return validatedParams, &clientErr{err}
	}
	validatedParams.Keywords = keywords

	releaseType, err := queryparams.GetReleaseType(ctx, params, queryparams.Published)
	if err != nil {
		return validatedParams, &clientErr{err}
	}
	validatedParams.ReleaseType = releaseType

	validatedParams.Provisional, _ = queryparams.GetBoolean(ctx, params, queryparams.Provisional.String(), false)
	validatedParams.Confirmed, _ = queryparams.GetBoolean(ctx, params, queryparams.Confirmed.String(), false)
	validatedParams.Postponed, _ = queryparams.GetBoolean(ctx, params, queryparams.Postponed.String(), false)
	validatedParams.Census, _ = queryparams.GetBoolean(ctx, params, queryparams.Census, false)
	validatedParams.Highlight, _ = queryparams.GetBoolean(ctx, params, queryparams.Highlight, true)

	return validatedParams, nil
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
	if err = toICSFile(ctx, releases.Releases, fileWriter); err != nil {
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
	for i := range releases {
		printLine("BEGIN:VEVENT")
		printLine("DTSTAMP:" + time.Now().UTC().Format(iCalDateFormat))
		releaseDate := iCalDate(ctx, releases[i].Description.ReleaseDate)
		printLine("DTSTART:" + releaseDate)
		printLine("DTEND:" + releaseDate)
		printLine("SUMMARY:" + releases[i].Description.Title)
		printLine("UID:" + releases[i].URI)
		printLine("STATUS:" + releaseStatus(releases[i]))
		printLine("DESCRIPTION:" + releases[i].Description.Summary)
		printLine("END:VEVENT")
	}
	printLine("END:VCALENDAR")

	return nil
}

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

func createRSSFeed(ctx context.Context, w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string, api SearchAPI, validatedParams queryparams.ValidatedParams) error {
	var err error
	uriPrefix := "https://www.ons.gov.uk"
	releases, err := api.GetReleases(ctx, accessToken, collectionID, lang, validatedParams.AsBackendQuery())
	if err != nil {
		setStatusCode(r, w, err)
		return err
	}

	w.Header().Set("Content-Type", "application/rss+xml")

	feed := &feeds.Feed{
		Title:       "ONS Release Calendar RSS Feed.",
		Link:        &feeds.Link{Href: "https://www.ons.gov.uk/releasecalendar"},
		Description: "Latest ONS releases",
	}

	feed.Items = []*feeds.Item{}
	for i := range releases.Releases {
		release := &releases.Releases[i]
		date, parseErr := time.Parse("2006-01-02T15:04:05.000Z", release.Description.ReleaseDate)
		if parseErr != nil {
			return fmt.Errorf("error parsing time: %s", parseErr)
		}
		item := &feeds.Item{
			Title:       release.Description.Title,
			Link:        &feeds.Link{Href: uriPrefix + release.URI},
			Description: release.Description.Summary,
			Id:          uriPrefix + release.URI,
			Created:     date,
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return fmt.Errorf("error converting to rss: %s", err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(rss))
	if err != nil {
		return fmt.Errorf("error writing rss to response: %s", err)
	}

	return nil
}

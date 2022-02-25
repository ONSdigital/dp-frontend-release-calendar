package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mapper"
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

	release, err := api.GetLegacyRelease(ctx, userAccessToken, collectionID, lang, req.URL.EscapedPath())
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateRelease(basePage, *release)

	rc.BuildPage(w, m, "release")
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

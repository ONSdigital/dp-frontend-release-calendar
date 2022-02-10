package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/mapper"
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

func ReleaseSample(cfg config.Config, rc RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		releaseSample(w, req, rc, cfg)
	}
}

func releaseSample(w http.ResponseWriter, req *http.Request, rc RenderClient, cfg config.Config) {
	ctx := req.Context()
	basePage := rc.NewBasePageModel()
	m := mapper.CreateRelease(ctx, basePage, cfg)

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

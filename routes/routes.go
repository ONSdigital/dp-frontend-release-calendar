package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/handlers"

	render "github.com/ONSdigital/dp-renderer"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
	Render             *render.Render
	ReleaseCalendarAPI *releasecalendar.Client
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)

	r.StrictSlash(true).Path("/releases/{uri:.*}").Methods("GET").HandlerFunc(handlers.Release(*cfg, c.Render, c.ReleaseCalendarAPI))
	r.StrictSlash(true).Path("/previousreleasessample").Methods("GET").HandlerFunc(handlers.PreviousReleasesSample(*cfg, c.Render))
	r.StrictSlash(true).Path("/calendarsample").Methods("GET").HandlerFunc(handlers.CalendarSample(*cfg, c.Render))
}

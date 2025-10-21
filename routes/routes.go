package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	search "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/handlers"

	render "github.com/ONSdigital/dis-design-system-go"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
	Render             *render.Render
	ReleaseCalendarAPI *releasecalendar.Client
	SearchAPI          *search.Client
	ZebedeeClient      *zebedee.Client
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)

	r.StrictSlash(true).Path(cfg.RoutingPrefix + "/releases/{uri}").Methods("GET").HandlerFunc(handlers.Release(*cfg, c.Render, c.ReleaseCalendarAPI, c.ZebedeeClient))
	r.StrictSlash(true).Path(cfg.RoutingPrefix + "/releases/{uri}/data").Methods("GET").HandlerFunc(handlers.ReleaseData(*cfg, c.ReleaseCalendarAPI))
	r.StrictSlash(true).Path(cfg.RoutingPrefix + "/releasecalendar").Methods("GET").HandlerFunc(handlers.ReleaseCalendar(*cfg, c.Render, c.SearchAPI, c.ZebedeeClient))
	r.StrictSlash(true).Path(cfg.RoutingPrefix + "/releasecalendar/data").Methods("GET").HandlerFunc(handlers.ReleaseCalendarData(*cfg, c.SearchAPI))
	r.StrictSlash(true).Path(cfg.RoutingPrefix + "/calendar/releasecalendar").Methods("GET").HandlerFunc(handlers.ReleaseCalendarICSEntries(*cfg, c.SearchAPI))
}

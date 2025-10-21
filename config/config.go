package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
)

type Config struct {
	APIRouterURL                string `envconfig:"API_ROUTER_URL"`
	BindAddr                    string `envconfig:"BIND_ADDR"`
	Debug                       bool   `envconfig:"DEBUG"`
	DefaultLimit                int    `envconfig:"DEFAULT_LIMIT"`
	DefaultMaximumLimit         int    `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	DefaultMaximumSearchResults int    `envconfig:"DEFAULT_MAXIMUM_SEARCH_RESULTS"`
	DefaultSort                 string `envconfig:"DEFAULT_SORT"`
	Deprecation                 Deprecation
	FeedbackAPIURL              string        `envconfig:"FEEDBACK_API_URL"`
	GracefulShutdownTimeout     time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckCriticalTimeout  time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthCheckInterval         time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	IsPublishing                bool          `envconfig:"IS_PUBLISHING"`
	PatternLibraryAssetsPath    string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	RoutingPrefix               string        `envconfig:"ROUTING_PREFIX"`
	SiteDomain                  string        `envconfig:"SITE_DOMAIN"`
	SupportedLanguages          []string      `envconfig:"SUPPORTED_LANGUAGES"`
}

type Deprecation struct {
	DeprecateEndpoint  bool   `envconfig:"DEPRECATE_ENDPOINT"`
	Deprecation        string `envconfig:"DEPRECATION"`
	DeprecationMessage string `envconfig:"DEPRECATION_MESSAGE"`
	Link               string `envconfig:"LINK"`
	Sunset             string `envconfig:"SUNSET"`
}

var cfg *Config

var RendererVersion = "v0.2.0"

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	var err error

	cfg, err = get()
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.PatternLibraryAssetsPath = "http://localhost:9002/dist/assets"
	} else {
		cfg.PatternLibraryAssetsPath = fmt.Sprintf("//cdn.ons.gov.uk/dis-design-system-go/%s", RendererVersion)
	}

	cfg.RoutingPrefix = validateRoutingPrefix(cfg.RoutingPrefix)

	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		APIRouterURL:                "http://localhost:23200/v1",
		BindAddr:                    ":27700",
		Debug:                       false,
		DefaultLimit:                10,
		DefaultMaximumLimit:         100,
		DefaultMaximumSearchResults: 1000,
		DefaultSort:                 queryparams.RelDateDesc.String(),
		Deprecation: Deprecation{
			DeprecateEndpoint:  false,
			Deprecation:        "", // could be of format "2025-08-29T10:00:00Z, 2025-08-29 15:04:05 and 2025-08-29"
			DeprecationMessage: "",
			Link:               "",
			Sunset:             "", // could be of format "2025-08-29T10:00:00Z, 2025-08-29 15:04:05 and 2025-08-29"
		},
		FeedbackAPIURL:             "http://localhost:23200/v1/feedback",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		IsPublishing:               false,
		RoutingPrefix:              "",
		SiteDomain:                 "localhost",
		SupportedLanguages:         []string{"en", "cy"},
	}

	return cfg, envconfig.Process("", cfg)
}

func validateRoutingPrefix(prefix string) string {
	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		return "/" + prefix
	}

	return prefix
}

func (cfg *Config) CalendarPath() string {
	return cfg.RoutingPrefix + "/releasecalendar"
}

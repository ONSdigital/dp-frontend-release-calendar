package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
)

type Config struct {
	APIRouterURL                  string        `envconfig:"API_ROUTER_URL"`
	BabbageMaxAgeKey              string        `envconfig:"BABBAGE_MAXAGE_KEY"`
	BabbageURL                    string        `envconfig:"BABBAGE_URL"`
	BindAddr                      string        `envconfig:"BIND_ADDR"`
	Debug                         bool          `envconfig:"DEBUG"`
	DefaultLimit                  int           `envconfig:"DEFAULT_LIMIT"`
	DefaultMaximumLimit           int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	DefaultMaximumSearchResults   int           `envconfig:"DEFAULT_MAXIMUM_SEARCH_RESULTS"`
	DefaultSort                   string        `envconfig:"DEFAULT_SORT"`
	EnableBabbageCalculatedMaxAge bool          `envconfig:"ENABLE_BABBAGE_CALCULATED_MAX_AGE"`
	GracefulShutdownTimeout       time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckCriticalTimeout    time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthCheckInterval           time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	IsPublishing                  bool          `envconfig:"IS_PUBLISHING"`
	PatternLibraryAssetsPath      string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	RoutingPrefix                 string        `envconfig:"ROUTING_PREFIX"`
	SiteDomain                    string        `envconfig:"SITE_DOMAIN"`
	SupportedLanguages            []string      `envconfig:"SUPPORTED_LANGUAGES"`
	EnableOfficialStatisticsLogo  bool          `envconfig:"ENABLE_OFFICIAL_STATISTICS_LOGO"`
}

var cfg *Config

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
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/b3361a2"
	}

	cfg.RoutingPrefix = validateRoutingPrefix(cfg.RoutingPrefix)

	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		APIRouterURL:                  "http://localhost:23200/v1",
		BabbageMaxAgeKey:              "",
		BabbageURL:                    "http://localhost:8080",
		BindAddr:                      ":27700",
		Debug:                         false,
		DefaultLimit:                  10,
		DefaultMaximumLimit:           100,
		DefaultMaximumSearchResults:   1000,
		DefaultSort:                   queryparams.RelDateDesc.String(),
		EnableBabbageCalculatedMaxAge: true,
		GracefulShutdownTimeout:       5 * time.Second,
		HealthCheckCriticalTimeout:    90 * time.Second,
		HealthCheckInterval:           30 * time.Second,
		IsPublishing:                  false,
		RoutingPrefix:                 "",
		SiteDomain:                    "localhost",
		SupportedLanguages:            []string{"en", "cy"},
		EnableOfficialStatisticsLogo:  false,
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

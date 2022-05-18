package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
)

type Config struct {
	APIRouterURL                string        `envconfig:"API_ROUTER_URL"`
	BabbageURL                  string        `envconfig:"BABBAGE_URL"`
	MaxAgeKey                   string        `envconfig:"BABBAGE_MAXAGE_KEY"`
	BindAddr                    string        `envconfig:"BIND_ADDR"`
	Debug                       bool          `envconfig:"DEBUG"`
	DefaultLimit                int           `envconfig:"DEFAULT_LIMIT"`
	DefaultMaximumLimit         int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	DefaultMaximumSearchResults int           `envconfig:"DEFAULT_MAXIMUM_SEARCH_RESULTS"`
	DefaultSort                 string        `envconfig:"DEFAULT_SORT"`
	GracefulShutdownTimeout     time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval         time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout  time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	PatternLibraryAssetsPath    string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SupportedLanguages          []string      `envconfig:"SUPPORTED_LANGUAGES"`
	SiteDomain                  string        `envconfig:"SITE_DOMAIN"`
	PrivateRoutingPrefix        string        `envconfig:"PRIVATE_ROUTING_PREFIX"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {

	cfg, err := get()

	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.PatternLibraryAssetsPath = "http://localhost:9002/dist/assets"
	} else {
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/e9dc639"
	}

	cfg.PrivateRoutingPrefix = validatePrivatePrefix(cfg.PrivateRoutingPrefix)

	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		APIRouterURL:                "http://localhost:23200/v1",
		BabbageURL:                  "http://localhost:8080",
		MaxAgeKey:                   "",
		BindAddr:                    ":27700",
		Debug:                       false,
		DefaultLimit:                10,
		DefaultMaximumLimit:         100,
		DefaultSort:                 queryparams.RelDateDesc.String(),
		DefaultMaximumSearchResults: 1000,
		GracefulShutdownTimeout:     5 * time.Second,
		HealthCheckInterval:         30 * time.Second,
		HealthCheckCriticalTimeout:  90 * time.Second,
		SupportedLanguages:          []string{"en", "cy"},
		SiteDomain:                  "localhost",
		PrivateRoutingPrefix:        "",
	}

	return cfg, envconfig.Process("", cfg)
}

func validatePrivatePrefix(prefix string) string {
	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		return "/" + prefix
	}

	return prefix
}

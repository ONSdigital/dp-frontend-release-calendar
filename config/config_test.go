package config

import (
	"os"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the values should be set to the expected defaults", func() {
				So(cfg.Debug, ShouldBeFalse)
				So(cfg.BindAddr, ShouldEqual, ":27700")
				So(cfg.PatternLibraryAssetsPath, ShouldEqual, "//cdn.ons.gov.uk/dp-design-system/119ecec")
				So(cfg.SupportedLanguages, ShouldResemble, []string{"en", "cy"})
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.APIRouterURL, ShouldEqual, "http://localhost:23200/v1")
				So(cfg.BabbageURL, ShouldEqual, "http://localhost:8080")
				So(cfg.MaxAgeKey, ShouldEqual, "")
				So(cfg.DefaultLimit, ShouldEqual, 10)
				So(cfg.DefaultMaximumLimit, ShouldEqual, 100)
				So(cfg.DefaultSort, ShouldEqual, "date-newest")
				So(cfg.DefaultMaximumSearchResults, ShouldEqual, 1000)
				So(cfg.SiteDomain, ShouldEqual, "localhost")
				So(cfg.RoutingPrefix, ShouldEqual, "")
			})

			Convey("Then a second call to config should return the same config", func() {
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}

func TestValidateRoutingPrefix(t *testing.T) {
	Convey("when a routing prefix is not set", t, func() {
		So(validateRoutingPrefix(""), ShouldEqual, "")
	})
	Convey("when a routing prefix is set without an initial '/'", t, func() {
		So(validateRoutingPrefix("a-prefix"), ShouldEqual, "/a-prefix")
	})
	Convey("when a valid routing prefix is set with an initial '/'", t, func() {
		So(validateRoutingPrefix("/a-prefix"), ShouldEqual, "/a-prefix")
	})
}

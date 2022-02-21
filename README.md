# dp-frontend-release-calendar
Release Calendar frontend controller

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                 | Description
| ---------------------------- | ----------------------- | -----------
| BIND_ADDR                    | :27700                  | The host and port to bind to
| DEBUG                        | false                   | Enable debug mode
| SITE_DOMAIN                  | localhost               |
| PATTERN_LIBRARY_ASSETS_PATH  | ""                      | Pattern library location
| SUPPORTED_LANGUAGES          | [2]string{"en", "cy"}   | Supported languages
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                      | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s                     | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                     | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.


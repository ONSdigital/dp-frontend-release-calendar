# dp-frontend-release-calendar

## Release Calendar frontend controller

Provides server side rendering of the Release Calendar and Release pages.

A Release Calender page is constructed from metadata drawn from the Search API.

A Release page is constructed from the data drawn from the Release API.

### Getting started

* Run the Digital Publishing Design System (`dp-design-system`) in a
  separate shell with `./run.sh`
* Run `make debug`
* In your browser, visit one of:
  - `http://localhost:27700/calendarsample`
  - `http://localhost:27700/releases/{topic}` where `{topic}` exists in `zebedee/master/releases/`
  - `http://localhost:27700/previousreleasessample`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable           | Default                   | Description|
|--------------------------------|---------------------------| -----------|
| BIND_ADDR                      | :27700                    | The host and port to bind to|
| DEBUG                          | false                     | Enable debug mode|
| API_ROUTER_URL                 | http://localhost:23200/v1 | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)|
| BABBAGE_URL                    | http://localhost:8080     | The URL of [babbage](https://github.com/ONSdigital/babbage)|
| SITE_DOMAIN                    | localhost                 ||
| PATTERN_LIBRARY_ASSETS_PATH    | ""                        | Pattern library location|
| SUPPORTED_LANGUAGES            | [2]string{"en", "cy"}     | Supported languages|
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                        | The graceful shutdown timeout in seconds (`time.Duration` format)|
| HEALTHCHECK_INTERVAL           | 30s                       | Time between self-healthchecks (`time.Duration` format)|
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)|
| DEFAULT_LIMIT                  | 10                        | The default size of (number of search results on) a page|
| DEFAULT_MAXIMUM_LIMIT          | 100                       | The default maximum size of (number of search results on) a page|
| DEFAULT_MAXIMUM_SEARCH_RESULTS | 1000                      | The default maximum number of search results that will be paged|
| DEFAULT_SORT                   | "release_date_desc"       | The default sort order of search results |
| BABBAGE_MAXAGE_KEY             | ""                        | The key required to get the max age value from babbage |
| ROUTING_PREFIX                 | ""                        | Any routing prefix for the service |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

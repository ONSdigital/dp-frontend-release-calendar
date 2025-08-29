# dp-frontend-release-calendar

## Release Calendar frontend controller

Provides server side rendering of the Release Calendar and Release pages.

A Release Calender page is constructed from metadata drawn from the Search API.
See [search service architecture docs here](https://github.com/ONSdigital/dp-search-api/tree/develop/architecture#search-service-architecture)

A Release page is constructed from the data drawn from the Release API.

### Getting started

Run `make help` to see full list of make targets.

* Run the Digital Publishing Design System (`dp-design-system`) in a
  separate shell with `./run.sh`
* Run `make debug`
* In your browser, visit one of:
  * `http://localhost:27700/releasecalendar`
  * `http://localhost:27700/releases/{topic}` where `{topic}` exists in `zebedee/master/releases/`
* For document data underlying each page, visit one of:
  * `http://localhost:27700/releasecalendar/data`
  * `http://localhost:27700/releases/{topic}/data`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable           | Default                     | Description                                                                                                        |
|--------------------------------|-----------------------------|--------------------------------------------------------------------------------------------------------------------|
| API_ROUTER_URL                 | <http://localhost:23200/v1> | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)                                        |
| BIND_ADDR                      | :27700                      | The host and port to bind to                                                                                       |
| DEBUG                          | false                       | Enable debug mode                                                                                                  |
| DEFAULT_LIMIT                  | 10                          | The default size of (number of search results on) a page                                                           |
| DEFAULT_MAXIMUM_LIMIT          | 100                         | The default maximum size of (number of search results on) a page                                                   |
| DEFAULT_MAXIMUM_SEARCH_RESULTS | 1000                        | The default maximum number of search results that will be paged                                                    |
| DEFAULT_SORT                   | "release_date_desc"         | The default sort order of search results                                                                           |
| DEPRECATION_DATE               | ""                          | The date in which the decision was made to deprecate an endpoint e.g. "Wed, 11 Jul 2025 23:59:59 GMT"              |
| DEPRECATION_MESSAGE            | ""                          | Message to be given on API response to deprecated endpoint                                                         |
| ENDPOINT_DEPRECATION           | false                       | Enable endpoint deprecation                                                                                        |
| FEEDBACK_API_URL               | [http://localhost:23200/v1/feedback](http://localhost:23200/v1/feedback) | The public `dp-api-router` address for feedback, not the internal one |
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                          | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                         | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| HEALTHCHECK_INTERVAL           | 30s                         | Time between self-healthchecks (`time.Duration` format)                                                            |
| IS_PUBLISHING                  | false                       | Mode in which the service is running                                                                               |
| PATTERN_LIBRARY_ASSETS_PATH    | ""                          | Pattern library location                                                                                           |
| ROUTING_PREFIX                 | ""                          | Any routing prefix for the service                                                                                 |
| SITE_DOMAIN                    | localhost                   |                                                                                                                    |
| SUNSET_DATE                    | ""                          | The date when this service will cease to return data on its deprecated endpoints and instead return blanket 404 status codes, e.g. "Fri, 11 Aug 2025 23:59:59 GMT"               |
| SUNSET_LINK                    | ""                          | A url to further information of the deprecation of the service or endpoints                                        |
| SUPPORTED_LANGUAGES            | []string{"en", "cy"}        | Supported languages                                                                                                |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2025, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

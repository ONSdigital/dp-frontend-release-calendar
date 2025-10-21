# dp-frontend-release-calendar

## Release Calendar frontend controller

Provides server side rendering of the Release Calendar and Release pages.

A Release Calender page is constructed from metadata drawn from the Search API.
See [search service architecture docs here](https://github.com/ONSdigital/dp-search-api/tree/develop/architecture#search-service-architecture)

A Release page is constructed from the data drawn from the Release API.

### Getting started

Run `make help` to see full list of make targets.

* Run [`dis-design-system-go`](https://github.com/ONSdigital/dis-design-system-go) in a separate shell to generate static assets (css/js)
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
| FEEDBACK_API_URL               | [http://localhost:23200/v1/feedback](http://localhost:23200/v1/feedback) | The public `dp-api-router` address for feedback, not the internal one |
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                          | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                         | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| HEALTHCHECK_INTERVAL           | 30s                         | Time between self-healthchecks (`time.Duration` format)                                                            |
| IS_PUBLISHING                  | false                       | Mode in which the service is running                                                                               |
| PATTERN_LIBRARY_ASSETS_PATH    | ""                          | Pattern library location                                                                                           |
| ROUTING_PREFIX                 | ""                          | Any routing prefix for the service                                                                                 |
| SITE_DOMAIN                    | localhost                   |                                                                                                                    |
| SUPPORTED_LANGUAGES            | []string{"en", "cy"}        | Supported languages                                                                                                |

### Deprecation Configuration

The following environment variables are for deprecating an endpoint in this service i.e. `/releases/data`.

| Environment variable           | Default                     | Description                                                                                                        |
|--------------------------------|-----------------------------|--------------------------------------------------------------------------------------------------------------------|
| DEPRECATE_ENDPOINT             | false                       | Enable endpoint deprecation                                                                                        |
| DEPRECATION                    | ""                          | The date in which the decision was made to deprecate an endpoint                                                   |
| DEPRECATION_MESSAGE            | ""                          | Message to be given on API response to deprecated endpoint                                                         |
| LINK                           | ""                          | A url to further information of the deprecation of the service or endpoints                                        |
| SUNSET                         | ""                          | The date when this service will cease to return data on a deprecated endpoint and instead return a 404 status code with a message                   |

To deprecate an endpoint all values should be set on the environmental variables as follows;

```json
[
  {
    "DEPRECATE_ENDPOINT": true,
    "DEPRECATION": "2025-08-29",
    "DEPRECATION_MESSAGE": "This endpoint is now deprecated.",
    "LINK": "https://www.ons.gov.uk",
    "SUNSET": "2025-08-29"
  }
]
```

All dates are parsed according to either RFC3339 ("2006-01-02T15:04:05Z07:00"), DateOnly ("2006-01-02") or DateTime ("2006-01-02 15:04:05") as defined in [Go's `time` package](https://pkg.go.dev/time#pkg-constants).

Run with endpoint deprecation on

```make
make debug DEPRECATE_ENDPOINT=true SUNSET="2025-08-29" LINK="https://www.ons.gov.uk" DEPRECATION="2025-08-29T10:00:00Z" DEPRECATION_MESSAGE="The release data endpoint is now deprecated"
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2025, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

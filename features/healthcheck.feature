Feature: Healthcheck endpoint should inform the health of service

    Scenario: Returning a OK (200) status when health endpoint called  
        Given the release calendar is running
        And the downstream service is healthy
        And I wait 2 seconds for the healthcheck to be available
        When I GET "/health"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status":"OK",
                "version":{
                    "git_commit":"3t7e5s1t4272646ef477f8ed755",
                    "language":"go",
                    "language_version":"go1.17.8",
                    "version":"v1.2.3"
                },
                "checks":[
                    {
                        "name":"Release Calendar API",
                        "status":"OK",
                        "status_code":200,
                        "message":"release-calendar-api is ok"
                    }
                ]
            }
        """

    Scenario: Returning a WARNING (429) status when one downstream service is warning
        Given the release calendar is running  
        And the downstream service is warning
        And I wait 2 seconds for the healthcheck to be available
        When I GET "/health"
        Then the HTTP status code should be "429"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status": "WARNING",
                "version": {
                    "git_commit": "3t7e5s1t4272646ef477f8ed755",
                    "language": "go",
                    "language_version": "go1.17.8",
                    "version": "v1.2.3"
                },
                "checks": [
                    {
                        "name": "Release Calendar API",
                        "status": "WARNING",
                        "status_code": 429,
                        "message": "release-calendar-api is degraded, but at least partially functioning"
                    }
                ]
            }
        """

    Scenario: Returning a WARNING (429) status when one downstream service is critical and critical timeout has not expired  
        Given the release calendar is running
        And the downstream service is failing
        And I wait 2 seconds for the healthcheck to be available
        When I GET "/health"
        Then the HTTP status code should be "429"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status": "WARNING",
                "version": {
                    "git_commit": "3t7e5s1t4272646ef477f8ed755",
                    "language": "go",
                    "language_version": "go1.17.8",
                    "version": "v1.2.3"
                },
                "checks": [
                    {
                        "name": "Release Calendar API",
                        "status": "CRITICAL",
                        "status_code": 500,
                        "message": "release-calendar-api functionality is unavailable or non-functioning"
                    }
                ]
            }
        """

    Scenario: Returning a CRITICAL (500) status when health endpoint called
        Given the release calendar is running
        And the downstream service is failing
        And I wait 2 seconds for the healthcheck to be available
        When I GET "/health"
        And I wait 4 seconds to pass the critical timeout
        And I GET "/health"
        Then the HTTP status code should be "500"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status": "CRITICAL",
                "version": {
                    "git_commit": "3t7e5s1t4272646ef477f8ed755",
                    "language": "go",
                    "language_version": "go1.17.8",
                    "version": "v1.2.3"
                },
                "checks": [
                    {
                        "name": "Release Calendar API",
                        "status": "CRITICAL",
                        "status_code": 500,
                        "message": "release-calendar-api functionality is unavailable or non-functioning"
                    }
                ]
            }
        """

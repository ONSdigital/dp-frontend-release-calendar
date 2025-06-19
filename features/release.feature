@ReleaseLegacy
Feature: Release Legacy
  Scenario: Get a single Release calendar entry
    Given there is a Release Calendar API that gives a successful response for "/releases/myrelease"
    And the release calendar is running
    When I navigate to "/releases/myrelease"
    And the page should have the following content
      """
          {
              "#main h1": "My test release"
          }
      """
  Scenario: Get a single Release calendar entry with a migration Link
    Given there is a Release Calendar API that gives a successful response for "/releases/myrelease" with a migration link
    And the release calendar is running
    When I GET "/releases/myrelease"
    Then the HTTP status code should be "308"
    And the response header "Location" should be "/redirect1"

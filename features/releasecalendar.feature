Feature: Release Calendar

  Scenario: GET /releasecalendar and checking for zero results
    Given there is a Search API that gives a successful response and returns 0 results
    And the release calendar is running
    When I navigate to "/releasecalendar"
    And the page should have the following content
      """
          {
              "#main h1": "Release calendar",
              "#results": "0 results"
          }
      """

  Scenario: GET /releasecalendar and checking for one result
    Given there is a Search API that gives a successful response and returns 1 results
    And the release calendar is running
    When I navigate to "/releasecalendar"
    And the page should have the following content
      """
          {
              "#main h1": "Release calendar",
              "#results": "1 result",
              ".ons-pagination__position": "Page 1 of 1"
          }
      """

  Scenario: GET /releasecalendar and checking for 10 results
    Given there is a Search API that gives a successful response and returns 10 results
    And the release calendar is running
    When I navigate to "/releasecalendar"
    And the page should have the following content
      """
          {
              "#main h1": "Release calendar",
              "#results": "10 results",
              ".ons-pagination__position": "Page 1 of 1"
          }
      """

  Scenario: GET /releasecalendar and checking for 11 results
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar"
    And the page should have the following content
      """
          {
              "#main h1": "Release calendar",
              "#results": "11 results",
              ".ons-pagination__position": "Page 1 of 2"
          }
      """

  Scenario: GET /releasecalendar with a invalid page number (exceeding total pages)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?page=3"
    And the page should have the following content
      """
          {
              ".ons-list__link": "invalid page parameter: value is above total pages (2)"
          }
      """

  Scenario: GET /releasecalendar with a invalid page number (exceeding max pages)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?page=101"
    And the page should have the following content
      """
          {
              ".ons-list__link": "invalid page parameter: value is above the maximum value (100)"
          }
      """

  Scenario: GET /releasecalendar with a invalid page number (not a number)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?page=dog"
    And the page should have the following content
      """
          {
              ".ons-list__link": "invalid page parameter: enter a number"
          }
      """

  Scenario: GET /releasecalendar with a invalid page number (below min pages)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?page=-1"
    And the page should have the following content
      """
          {
              ".ons-list__link": "invalid page parameter: value is below the minimum value (1)"
          }
      """

  Scenario: GET /releasecalendar with an invalid after date (invalid day)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=32&after-month=1&after-year=2000&before-day=1&before-month=1&before-year=2020"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a real date"
          }
      """

  Scenario: GET /releasecalendar with an invalid after date (invalid month)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=1&after-month=13&after-year=2000&before-day=1&before-month=1&before-year=2020"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a real date"
          }
      """

  Scenario: GET /releasecalendar with an invalid after date (invalid year)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=1&after-month=1&after-year=abcd&before-day=1&before-month=1&before-year=2020"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a number for released after year"
          }
      """

  Scenario: GET /releasecalendar with an invalid before date (invalid day)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=1&after-month=1&after-year=2000&before-day=32&before-month=1&before-year=2020"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a real date"
          }
      """

  Scenario: GET /releasecalendar with an invalid before date (invalid month)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=1&after-month=1&after-year=2000&before-day=1&before-month=13&before-year=2020"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a real date"
          }
      """

  Scenario: GET /releasecalendar with an invalid before date (invalid year)
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=1&after-month=1&after-year=2000&before-day=1&before-month=1&before-year=abcd"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a number for released before year"
          }
      """

  Scenario: GET /releasecalendar with before date earlier than after date
    Given there is a Search API that gives a successful response and returns 11 results
    And the release calendar is running
    When I navigate to "/releasecalendar?after-day=1&after-month=1&after-year=2020&before-day=1&before-month=1&before-year=2000"
    And the page should have the following content
      """
          {
              ".ons-list__link": "Enter a released before year that is later than 2020"
          }
      """
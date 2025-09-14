package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	"github.com/ONSdigital/dp-frontend-release-calendar/service"
	"github.com/ONSdigital/dp-frontend-release-calendar/service/mocks"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/cucumber/godog"
	"github.com/maxcnunes/httpfake"
	"github.com/stretchr/testify/assert"
)

// HealthCheckTest represents a test healthcheck struct that mimics the real healthcheck struct
type HealthCheckTest struct {
	Status    string                  `json:"status"`
	Version   healthcheck.VersionInfo `json:"version"`
	Uptime    time.Duration           `json:"uptime"`
	StartTime time.Time               `json:"start_time"`
	Checks    []*Check                `json:"checks"`
}

// Check represents a health status of a registered app that mimics the real check struct
// As the component test needs to access fields that are not exported in the real struct
type Check struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StatusCode  int        `json:"status_code"`
	Message     string     `json:"message"`
	LastChecked *time.Time `json:"last_checked"`
	LastSuccess *time.Time `json:"last_success"`
	LastFailure *time.Time `json:"last_failure"`
}

// RegisterSteps registers the specific steps needed to do component tests for the release calendar
func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I should receive the following health JSON response:$`, c.iShouldReceiveTheFollowingHealthJSONResponse)
	ctx.Step(`^the downstream service is (healthy|warning|failing)$`, c.theDownstreamServiceStatus)
	ctx.Step(`^the release calendar is running$`, c.theReleaseCalendarIsRunning)
	ctx.Step(`^there is a Search API that gives a successful response and returns ([1-9]\d*|0) results`, c.thereIsASearchAPIThatGivesASuccessfulResponseAndReturnsResults)
	ctx.Step(`^there is a Release Calendar API that gives a successful response for "([^"]*)"$`, c.thereIsAReleaseAPIThatGivesASuccessfulResponseFor)
	ctx.Step(`^there is a Release Calendar API that gives a successful response for "([^"]*)" with a migration link`, c.thereIsAReleaseAPIThatGivesASuccessfulResponseForWithMigrationLink)
}

func (c *Component) theReleaseCalendarIsRunning() error {
	ctx := context.Background()

	initFunctions := &mocks.InitialiserMock{
		DoGetHTTPServerFunc:   c.getHTTPServer,
		DoGetHealthCheckFunc:  getHealthCheckOK,
		DoGetHealthClientFunc: c.getHealthClient,
	}

	serviceList := service.NewServiceList(initFunctions)

	c.svc = service.New()
	if err := c.svc.Init(ctx, c.Config, serviceList); err != nil {
		log.Error(ctx, "failed to init service", err)
		return err
	}

	svcErrors := make(chan error, 1)

	c.StartTime = time.Now()

	c.svc.Run(ctx, svcErrors)
	c.ServiceRunning = true

	return nil
}

func (c *Component) theDownstreamServiceStatus(status string) error {
	var statusCode int
	switch status {
	case "healthy":
		statusCode = 200
	case "warning":
		statusCode = 429
	case "failing":
		statusCode = 500
	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	return c.setDownstreamServiceStatus(statusCode)
}

func (c *Component) setDownstreamServiceStatus(statusCode int) error {
	c.FakeAPIRouter.healthRequest.Lock()
	defer c.FakeAPIRouter.healthRequest.Unlock()

	c.FakeAPIRouter.healthRequest.CustomHandle = healthCheckStatusHandle(statusCode)

	return nil
}

func healthCheckStatusHandle(status int) httpfake.Responder {
	return func(w http.ResponseWriter, _ *http.Request, rh *httpfake.Request) {
		rh.Lock()
		defer rh.Unlock()
		w.WriteHeader(status)
	}
}

func (c *Component) iShouldReceiveTheFollowingHealthJSONResponse(expectedResponse *godog.DocString) error {
	var healthResponse, expectedHealth HealthCheckTest

	responseBody, err := io.ReadAll(c.APIFeature.HTTPResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response of release calendar component - error: %v", err)
	}

	err = json.Unmarshal(responseBody, &healthResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response of release calendar component - error: %v", err)
	}

	err = json.Unmarshal([]byte(expectedResponse.Content), &expectedHealth)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected health response - error: %v", err)
	}

	c.validateHealthCheckResponse(healthResponse, expectedHealth)

	return c.ErrorFeature.StepError()
}

func (c *Component) validateHealthCheckResponse(healthResponse, expectedResponse HealthCheckTest) {
	maxExpectedStartTime := c.StartTime.Add((c.Config.HealthCheckInterval + 1) * time.Second)

	assert.Equal(&c.ErrorFeature, expectedResponse.Status, healthResponse.Status)
	assert.True(&c.ErrorFeature, healthResponse.StartTime.After(c.StartTime))
	assert.True(&c.ErrorFeature, healthResponse.StartTime.Before(maxExpectedStartTime))
	assert.Greater(&c.ErrorFeature, healthResponse.Uptime.Seconds(), float64(0))

	c.validateHealthVersion(healthResponse.Version, expectedResponse.Version, maxExpectedStartTime)

	for i, checkResponse := range healthResponse.Checks {
		c.validateHealthCheck(checkResponse, expectedResponse.Checks[i])
	}
}

func (c *Component) validateHealthVersion(versionResponse, expectedVersion healthcheck.VersionInfo, maxExpectedStartTime time.Time) {
	assert.True(&c.ErrorFeature, versionResponse.BuildTime.Before(maxExpectedStartTime))
	assert.Equal(&c.ErrorFeature, expectedVersion.GitCommit, versionResponse.GitCommit)
	assert.Equal(&c.ErrorFeature, expectedVersion.Language, versionResponse.Language)
	assert.NotEmpty(&c.ErrorFeature, versionResponse.LanguageVersion)
	assert.Equal(&c.ErrorFeature, expectedVersion.Version, versionResponse.Version)
}

func (c *Component) validateHealthCheck(checkResponse, expectedCheck *Check) {
	maxExpectedHealthCheckTime := c.StartTime.Add((c.Config.HealthCheckInterval + c.Config.HealthCheckCriticalTimeout + 1) * time.Second)

	assert.Equal(&c.ErrorFeature, expectedCheck.Name, checkResponse.Name)
	assert.Equal(&c.ErrorFeature, expectedCheck.Status, checkResponse.Status)
	assert.Equal(&c.ErrorFeature, expectedCheck.StatusCode, checkResponse.StatusCode)
	assert.Equal(&c.ErrorFeature, expectedCheck.Message, checkResponse.Message)
	assert.True(&c.ErrorFeature, checkResponse.LastChecked.Before(maxExpectedHealthCheckTime))
	assert.True(&c.ErrorFeature, checkResponse.LastChecked.After(c.StartTime))

	if expectedCheck.StatusCode == 200 {
		assert.True(&c.ErrorFeature, checkResponse.LastSuccess.Before(maxExpectedHealthCheckTime))
		assert.True(&c.ErrorFeature, checkResponse.LastSuccess.After(c.StartTime))
	} else {
		assert.True(&c.ErrorFeature, checkResponse.LastFailure.Before(maxExpectedHealthCheckTime))
		assert.True(&c.ErrorFeature, checkResponse.LastFailure.After(c.StartTime))
	}
}

func (c *Component) thereIsASearchAPIThatGivesASuccessfulResponseAndReturnsResults(count int) error {
	c.FakeAPIRouter.searchReleasesRequest.Lock()
	defer c.FakeAPIRouter.searchReleasesRequest.Unlock()

	c.FakeAPIRouter.searchReleasesRequest.Response = generateReleasesResponse(count)

	return nil
}

func (c *Component) thereIsAReleaseAPIThatGivesASuccessfulResponseFor() error {
	c.FakeAPIRouter.releaseRequest.Lock()
	defer c.FakeAPIRouter.releaseRequest.Unlock()

	mockedResult := releasecalendar.Release{
		Description: releasecalendar.ReleaseDescription{
			Title: "My test release",
		},
	}

	c.FakeAPIRouter.releaseRequest = c.FakeAPIRouter.fakeHTTP.NewHandler().Get("/releases/legacy")
	c.FakeAPIRouter.releaseRequest.Response = generateReleaseEntryResponse(mockedResult)

	return nil
}

func (c *Component) thereIsAReleaseAPIThatGivesASuccessfulResponseForWithMigrationLink() error {
	c.FakeAPIRouter.releaseRequest.Lock()
	defer c.FakeAPIRouter.releaseRequest.Unlock()

	mockedResult := releasecalendar.Release{
		Description: releasecalendar.ReleaseDescription{
			Title:         "My test release with migration link",
			MigrationLink: "/redirect1",
		},
	}

	c.FakeAPIRouter.releaseRequest = c.FakeAPIRouter.fakeHTTP.NewHandler().Get("/releases/legacy")
	c.FakeAPIRouter.releaseRequest.Response = generateReleaseEntryResponse(mockedResult)

	return nil
}

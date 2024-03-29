package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
)

const babbageServiceName = "Babbage"

// BabbageClient is a client to call Babbage
type BabbageClient struct {
	cli dphttp.Clienter
	url string
}

// NewBabbageClient creates a new instance of BabbageClient with a given babbage url
func NewBabbageClient(babbageURL string) *BabbageClient {
	hcClient := health.NewClient(babbageServiceName, babbageURL)

	return &BabbageClient{
		cli: hcClient.Client,
		url: babbageURL,
	}
}

// Checker calls babbage health endpoint and returns a check object to the caller.
func (c *BabbageClient) Checker(ctx context.Context, check *healthcheck.CheckState) error {
	hcClient := health.Client{
		Client: c.cli,
		URL:    c.url,
		Name:   babbageServiceName,
	}

	return hcClient.Checker(ctx, check)
}

func (c *BabbageClient) get(ctx context.Context, uri string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, err
	}
	return c.cli.Do(ctx, req)
}

// GetMaxAge calls the relevant Babbage endpoint to find out the max age for the requested content uri
func (c *BabbageClient) GetMaxAge(ctx context.Context, contentURI, key string) (int, error) {
	var babbageURI string
	if contentURI == "/releasecalendar" {
		// There is a specific endpoint for the release calendar max age
		babbageURI = fmt.Sprintf("%s/releasecalendarmaxage?key=%s", c.url, key)
	} else {
		babbageURI = fmt.Sprintf("%s/maxage?uri=%s&key=%s", c.url, contentURI, key)
	}
	resp, err := c.get(ctx, babbageURI)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid response from babbage. Status %d", resp.StatusCode)
	}

	defer closeResponseBody(ctx, resp)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	maxAge, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}

	return maxAge, nil
}

// closeResponseBody closes the response body and logs an error containing the context if unsuccessful
func closeResponseBody(ctx context.Context, resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		log.Error(ctx, "error closing http response body", err)
	}
}

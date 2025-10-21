package handlers

import (
	"context"
	"io"
	"net/url"

	"github.com/ONSdigital/dis-design-system-go/model"
	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	search "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
)

// To mock interfaces in this file
//go:generate mockgen -source=clients.go -destination=mock_clients.go -package=handlers github.com/ONSdigital/dp-frontend-release-calendar/handlers

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

// RenderClient is an interface with methods for rendering a template
type RenderClient interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() model.Page
}

// ReleaseCalendarAPI is an interface for the Release Calendar API
type ReleaseCalendarAPI interface {
	GetLegacyRelease(ctx context.Context, userAccessToken, collectionID, lang, uri string) (*releasecalendar.Release, error)
}

// SearchAPI is an interface to the Search API for searching the Release (Calendars)
type SearchAPI interface {
	GetReleases(ctx context.Context, userAccessToken, collectionID, lang string, query url.Values) (search.ReleaseResponse, error)
}

// BabbageClient is an interface to Babbage
type BabbageAPI interface {
	GetMaxAge(ctx context.Context, contentURI, key string) (int, error)
}

// ZebedeeClient is an interface for zebedee client
type ZebedeeClient interface {
	GetHomepageContent(ctx context.Context, userAccessToken, collectionID, lang, path string) (m zebedee.HomepageContent, err error)
}

package steps

import (
	"github.com/maxcnunes/httpfake"
)

// FakeAPI contains all the information for a fake component API
type FakeAPI struct {
	fakeHTTP              *httpfake.HTTPFake
	healthRequest         *httpfake.Request
	searchRequest         *httpfake.Request
	searchReleasesRequest *httpfake.Request
	navigationRequest     *httpfake.Request
	releaseRequest        *httpfake.Request
}

// NewFakeAPI creates a new fake component API
func NewFakeAPI() *FakeAPI {
	return &FakeAPI{
		fakeHTTP: httpfake.New(),
	}
}

// Close closes the fake API
func (f *FakeAPI) Close() {
	f.fakeHTTP.Close()
}

package mapper

import (
	"context"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	search "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

func createTableOfContents(
	description model.ReleaseDescription,
	relatedDocuments []model.Link,
	relatedDatasets []model.Link,
	dateChanges []model.DateChange,
	releaseHistory []model.Link,
	codeOfPractice bool,
) coreModel.TableOfContents {
	toc := coreModel.TableOfContents{
		AriaLabelLocaliseKey: "TableOfContents",
		TitleLocaliseKey:     "Contents",
	}

	sections := make(map[string]coreModel.ContentSection)
	displayOrder := make([]string, 0)

	if description.Summary != "" {
		sections["summary"] = coreModel.ContentSection{
			Current: false,
			Title:   "Summary",
		}
		displayOrder = append(displayOrder, "summary")
	}

	if len(relatedDocuments) > 0 {
		sections["publications"] = coreModel.ContentSection{
			Current: false,
			Title:   "Publications",
		}
		displayOrder = append(displayOrder, "publications")
	}

	if len(relatedDatasets) > 0 {
		sections["data"] = coreModel.ContentSection{
			Current: false,
			Title:   "Data",
		}
		displayOrder = append(displayOrder, "data")
	}

	if (model.ContactDetails{} != description.Contact) {
		sections["contactdetails"] = coreModel.ContentSection{
			Current: false,
			Title:   "Contact details",
		}
		displayOrder = append(displayOrder, "contactdetails")
	}

	if len(dateChanges) > 0 {
		sections["changestothisreleasedate"] = coreModel.ContentSection{
			Current: false,
			Title:   "Changes to this release date",
		}
		displayOrder = append(displayOrder, "changestothisreleasedate")
	}

	if len(releaseHistory) > 0 {
		sections["releasehistory"] = coreModel.ContentSection{
			Current: false,
			Title:   "Release history",
		}
		displayOrder = append(displayOrder, "releasehistory")
	}

	if codeOfPractice {
		sections["codeofpractice"] = coreModel.ContentSection{
			Current: false,
			Title:   "Code of Practice",
		}
		displayOrder = append(displayOrder, "codeofpractice")
	}

	toc.Sections = sections
	toc.DisplayOrder = displayOrder

	return toc
}

func CreatePreviousReleases(_ context.Context, basePage coreModel.Page, _ config.Config) model.PreviousReleases {
	previousReleases := model.PreviousReleases{
		Page: basePage,
	}

	previousReleases.BetaBannerEnabled = true
	previousReleases.Metadata.Title = "Personal well-being in the UK, quarterly: July 2021 to September 2021"

	previousReleases.Description.Title = "Personal well-being in the UK, quarterly: July 2021 to September 2021"
	previousReleases.Description.Summary = "Quarterly estimates of life satisfaction, feeling that the things done in life are worthwhile, happiness and anxiety at the UK level, created using the Annual Population Survey (APS)."

	previousReleases.ReleaseHistory = make([]model.Link, 10)
	previousReleases.ReleaseHistory[0] = model.Link{
		Title:   "8 December 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[1] = model.Link{
		Title:   "26 July 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[2] = model.Link{
		Title:   "22 April 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[3] = model.Link{
		Title:   "21 January 2021",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[4] = model.Link{
		Title:   "17 October 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[5] = model.Link{
		Title:   "21 July 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[6] = model.Link{
		Title:   "23 April 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[7] = model.Link{
		Title:   "28 January 2020",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[8] = model.Link{
		Title:   "18 October 2019",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}
	previousReleases.ReleaseHistory[9] = model.Link{
		Title:   "17 July 2019",
		URI:     "localhost",
		Summary: "Updated with latest data",
	}

	previousReleases.TableOfContents = createTableOfContents(
		previousReleases.Description,
		nil,
		nil,
		nil,
		previousReleases.ReleaseHistory,
		false,
	)

	previousReleases.Pagination.CurrentPage = 6
	previousReleases.Pagination.TotalPages = 10
	previousReleases.Pagination.Limit = 10
	previousReleases.Pagination.PagesToDisplay = []coreModel.PageToDisplay{
		{PageNumber: 5, URL: "previousreleasessample/5"},
		{PageNumber: 6, URL: "previousreleasessample/6"},
		{PageNumber: 7, URL: "previousreleasessample/7"},
	}
	previousReleases.Pagination.FirstAndLastPages = []coreModel.PageToDisplay{
		{PageNumber: 1, URL: "previousreleasessample/1"},
		{PageNumber: 100, URL: "previousreleasessample/100"},
	}
	previousReleases.Pagination.LimitOptions = []int{10, 25}

	return previousReleases
}

func CreateRelease(basePage coreModel.Page, release releasecalendar.Release) model.Release {
	result := model.Release{
		Page:     basePage,
		Markdown: release.Markdown,
		Description: model.ReleaseDescription{
			Title:   release.Description.Title,
			Summary: release.Description.Summary,
			Contact: model.ContactDetails{
				Email:     release.Description.Contact.Email,
				Name:      release.Description.Contact.Name,
				Telephone: release.Description.Contact.Telephone,
			},
			NationalStatistic:  release.Description.NationalStatistic,
			ReleaseDate:        release.Description.ReleaseDate,
			NextRelease:        release.Description.NextRelease,
			Published:          release.Description.Published,
			Finalised:          release.Description.Finalised,
			Cancelled:          release.Description.Cancelled,
			CancellationNotice: release.Description.CancellationNotice,
			ProvisionalDate:    release.Description.ProvisionalDate,
		},
	}

	result.RelatedDatasets = mapLink(release.RelatedDatasets)
	result.RelatedDocuments = mapLink(release.RelatedDocuments)
	result.RelatedMethodology = mapLink(release.RelatedMethodology)
	result.RelatedMethodologyArticle = mapLink(release.RelatedMethodologyArticle)
	result.Links = mapLink(release.Links)

	result.DateChanges = []model.DateChange{}
	for _, dc := range release.DateChanges {
		result.DateChanges = append(result.DateChanges, model.DateChange{
			Date:         dc.Date,
			ChangeNotice: dc.ChangeNotice,
		})
	}

	result.BetaBannerEnabled = true
	result.Metadata.Title = release.Description.Title
	result.URI = release.URI
	result.CodeOfPractice = true

	result.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Release calendar",
			URI:   "/calendar",
		},
		{
			Title: "Published", // TODO Set this from data
			URI:   "/calendar", // TODO Integrate with Search API
		},
		{
			Title: release.Description.Title,
		},
	}

	result.TableOfContents = createTableOfContents(
		result.Description,
		result.RelatedDocuments,
		result.RelatedDatasets,
		result.DateChanges,
		result.ReleaseHistory,
		result.CodeOfPractice,
	)

	return result
}

func mapLink(links []releasecalendar.Link) []model.Link {
	res := []model.Link{}
	for _, l := range links {
		res = append(res, model.Link{
			Title:   l.Title,
			Summary: l.Summary,
			URI:     l.URI,
		})
	}
	return res
}

func CreateReleaseCalendar(basePage coreModel.Page, params queryparams.ValidatedParams, response search.ReleaseResponse, cfg config.Config) model.Calendar {
	calendar := model.Calendar{
		Page: basePage,
	}
	calendar.BetaBannerEnabled = true
	calendar.Metadata.Title = helper.Localise("ReleaseCalendarPageTitle", calendar.Language, 1)
	calendar.KeywordSearch = coreModel.CompactSearch{
		ElementId:        "keyword-search",
		InputName:        "keywords",
		Language:         calendar.Language,
		LabelLocaliseKey: "ReleaseCalendarPageSearchKeywords",
		SearchTerm:       params.Keywords,
	}

	calendar.Sort = model.Sort{
		Mode:    params.Sort.String(),
		Options: mapSortOptions(params),
	}

	calendar.AfterDate = coreModel.InputDate{
		Language:        calendar.Language,
		Id:              "after-date",
		InputNameDay:    "after-day",
		InputNameMonth:  "after-month",
		InputNameYear:   "after-year",
		InputValueDay:   params.AfterDate.DayString(),
		InputValueMonth: params.AfterDate.MonthString(),
		InputValueYear:  params.AfterDate.YearString(),
		Title:           helper.Localise("ReleasedAfter", calendar.Language, 1),
		Description:     helper.Localise("DateFilterDescription", calendar.Language, 1),
	}

	calendar.BeforeDate = coreModel.InputDate{
		Language:        calendar.Language,
		Id:              "before-date",
		InputNameDay:    "before-day",
		InputNameMonth:  "before-month",
		InputNameYear:   "before-year",
		InputValueDay:   params.BeforeDate.DayString(),
		InputValueMonth: params.BeforeDate.MonthString(),
		InputValueYear:  params.BeforeDate.YearString(),
		Title:           helper.Localise("ReleasedBefore", calendar.Language, 1),
		Description:     helper.Localise("DateFilterDescription", calendar.Language, 1),
	}

	calendar.ReleaseTypes = mapReleases(params, response, calendar.Language)

	totalResults := cfg.DefaultMaximumSearchResults
	if totalResults > response.Breakdown.Total {
		totalResults = response.Breakdown.Total
	}
	calendar.Pagination.TotalPages = queryparams.CalculatePageNumber(totalResults-1, params.Limit)
	calendar.Pagination.CurrentPage = queryparams.CalculatePageNumber(params.Offset, params.Limit)
	calendar.Pagination.Limit = params.Limit
	calendar.Pagination.PagesToDisplay = getPagesToDisplay(params, calendar.Pagination.TotalPages, defaultWindowSize)
	calendar.Pagination.FirstAndLastPages = getFirstAndLastPages(params, calendar.Pagination.TotalPages)
	calendar.Pagination.LimitOptions = []int{10, 25}

	for _, release := range response.Releases {
		calendar.Entries = append(calendar.Entries, calendarEntryFromRelease(release))
	}

	return calendar
}

const defaultWindowSize = 5

func getPagesToDisplay(params queryparams.ValidatedParams, totalPages, windowSize int) []coreModel.PageToDisplay {
	start, end := getWindowStartEndPage(params.Page, totalPages, windowSize)

	var pagesToDisplay []coreModel.PageToDisplay
	for i := start; i <= end; i++ {
		pagesToDisplay = append(pagesToDisplay, coreModel.PageToDisplay{
			PageNumber: i,
			URL:        getPageURL(i, params),
		})
	}

	return pagesToDisplay
}

func getFirstAndLastPages(params queryparams.ValidatedParams, totalPages int) []coreModel.PageToDisplay {
	return []coreModel.PageToDisplay{
		{
			PageNumber: 1,
			URL:        getPageURL(1, params),
		},
		{
			PageNumber: totalPages,
			URL:        getPageURL(totalPages, params),
		},
	}
}

// getWindowStartEndPage calculates the start and end page of the moving window of size windowSize, over the set of pages
// whose current page is currentPage, and whose size is totalPages
// It is an error to pass a parameter whose value is < 1, or a currentPage > totalPages, and the function will panic in this case
func getWindowStartEndPage(currentPage, totalPages, windowSize int) (int, int) {
	if currentPage < 1 || totalPages < 1 || windowSize < 1 || currentPage > totalPages {
		panic("invalid parameters for getWindowStartEndPage - see documentation")
	}
	switch {
	case windowSize == 1:
		se := (currentPage % totalPages) + 1
		return se, se
	case windowSize >= totalPages:
		return 1, totalPages
	}

	windowOffset := getWindowOffset(windowSize)
	start := currentPage - windowOffset
	switch {
	case start <= 0:
		start = 1
	case start > totalPages-windowSize+1:
		start = totalPages - windowSize + 1
	}

	end := start + windowSize - 1
	if end > totalPages {
		end = totalPages
	}

	return start, end
}

func getPageURL(page int, params queryparams.ValidatedParams) (pageURL string) {
	query := params.AsQuery()
	query.Set("page", strconv.Itoa(page))

	return "/releasecalendar?" + query.Encode()
}

func getWindowOffset(windowSize int) int {
	if windowSize%2 == 0 {
		return (windowSize / 2) - 1
	}

	return windowSize / 2
}

func calendarEntryFromRelease(release search.Release) model.CalendarEntry {
	result := model.CalendarEntry{
		URI:         release.URI,
		DateChanges: dateChanges(release.DateChanges),
		Description: model.ReleaseDescription{
			Title:           release.Description.Title,
			Summary:         release.Description.Summary,
			ReleaseDate:     release.Description.ReleaseDate,
			ProvisionalDate: release.Description.ProvisionalDate,
			Published:       release.Description.Published,
			Cancelled:       release.Description.Cancelled,
			Finalised:       release.Description.Finalised,
		},
	}

	if highlight := release.Highlight; highlight != nil {
		switch {
		case highlight.Title != "":
			result.Description.Title = highlight.Title
		case highlight.Summary != "":
			result.Description.Summary = highlight.Summary
		}
	}

	return result
}

func CreateCalendar(_ context.Context, basePage coreModel.Page, _ config.Config) model.Calendar {
	calendar := model.Calendar{
		Page: basePage,
	}
	calendar.BetaBannerEnabled = true
	calendar.Metadata.Title = helper.Localise("ReleaseCalendarPageTitle", calendar.Language, 1)

	item1 := model.CalendarEntry{
		URI: "/releases/title1",
		Description: model.ReleaseDescription{
			Title:       "Public Sector Employment, UK: September 2021",
			Summary:     "A summary for Title 1",
			ReleaseDate: time.Now().AddDate(0, 0, -10).UTC().Format(time.RFC3339),
			Published:   true,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
			NextRelease: "To be announced",
		},
	}

	item2 := model.CalendarEntry{
		URI: "/releases/title2",
		Description: model.ReleaseDescription{
			Title:       "Labour market in the regions of the UK: December 2021",
			Summary:     "A summary for Title 2",
			ReleaseDate: time.Now().AddDate(0, 0, -15).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   true,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item3 := model.CalendarEntry{
		URI: "/releases/title3",
		Description: model.ReleaseDescription{
			Title:       "Personal well-being in the UK, quarterly: July 2021 to September 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item4 := model.CalendarEntry{
		URI: "/releases/title4",
		Description: model.ReleaseDescription{
			Title:       "Labour market statistics time series: December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item5 := model.CalendarEntry{
		URI: "/releases/title5",
		Description: model.ReleaseDescription{
			Title:       "UK labour market: December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item6 := model.CalendarEntry{
		URI: "/releases/title6",
		Description: model.ReleaseDescription{
			Title:       "Earnings and employment from Pay As You Earn Real Time Information, UK: December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item7 := model.CalendarEntry{
		URI: "/releases/title7",
		Description: model.ReleaseDescription{
			Title:       "Civil partnerships in England and Wales: 2020",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item8 := model.CalendarEntry{
		URI: "/releases/title8",
		Description: model.ReleaseDescription{
			Title:       "Understanding towns: industry analysis",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item9 := model.CalendarEntry{
		URI: "/releases/title9",
		Description: model.ReleaseDescription{
			Title:       "Disaggregating annual subnational gross value added (GVA) to lower levels of geography",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	item10 := model.CalendarEntry{
		URI: "/releases/title10",
		Description: model.ReleaseDescription{
			Title:       "Coronavirus (COVID-19) Infection Survey, UK: 8 December 2021",
			Summary:     "A summary for Title 3",
			ReleaseDate: time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339),
			Published:   false,
			Cancelled:   false,
			Contact:     model.ContactDetails{Name: "test publisher", Email: "testpublisher@ons.gov.uk"},
		},
	}

	calendar.Pagination.CurrentPage = 6
	calendar.Pagination.TotalPages = 10
	calendar.Pagination.Limit = 10
	calendar.Pagination.PagesToDisplay = []coreModel.PageToDisplay{
		{PageNumber: 5, URL: "calendarsample/5"},
		{PageNumber: 6, URL: "calendarsample/6"},
		{PageNumber: 7, URL: "calendarsample/7"},
	}
	calendar.Pagination.FirstAndLastPages = []coreModel.PageToDisplay{
		{PageNumber: 1, URL: "calendarsample/1"},
		{PageNumber: 100, URL: "calendarsample/100"},
	}
	calendar.Pagination.LimitOptions = []int{10, 25}

	calendar.Entries = make([]model.CalendarEntry, 10)
	calendar.Entries[0] = item1
	calendar.Entries[1] = item2
	calendar.Entries[2] = item3
	calendar.Entries[3] = item4
	calendar.Entries[4] = item5
	calendar.Entries[5] = item6
	calendar.Entries[6] = item7
	calendar.Entries[7] = item8
	calendar.Entries[8] = item9
	calendar.Entries[9] = item10

	calendar.ReleaseTypes = map[string]model.ReleaseType{
		"type-published": {
			Name:      "release-type",
			Value:     "type-published",
			Id:        "release-type-published",
			LocaleKey: "FilterReleaseTypePublished",
			Plural:    1,
			Language:  calendar.Language,
			Checked:   true,
			Count:     450,
		},
		"type-upcoming": {
			Name:      "release-type",
			Value:     "type-upcoming",
			Id:        "release-type-upcoming",
			LocaleKey: "FilterReleaseTypeUpcoming",
			Plural:    1,
			Language:  calendar.Language,
			Checked:   false,
			Count:     234,
		},
		"type-cancelled": {
			Name:      "release-type",
			Value:     "type-cancelled",
			Id:        "release-type-cancelled",
			LocaleKey: "FilterReleaseTypeCancelled",
			Plural:    1,
			Language:  calendar.Language,
			Checked:   false,
			Count:     0,
		},
	}

	calendar.Sort = model.Sort{
		Mode: "alphabetical-az",
		Options: []model.SortOption{
			{
				LocaleKey: "ReleaseCalendarSortOptionDateNewest",
				Plural:    1,
				Value:     "date-newest",
			},
			{
				LocaleKey: "ReleaseCalendarSortOptionDateOldest",
				Plural:    1,
				Value:     "date-oldest",
			},
			{
				LocaleKey: "ReleaseCalendarSortOptionAlphabeticalAZ",
				Plural:    1,
				Value:     "alphabetical-az",
			},
			{
				LocaleKey: "ReleaseCalendarSortOptionAlphabeticalZA",
				Plural:    1,
				Value:     "alphabetical-za",
			},
			{
				LocaleKey: "ReleaseCalendarSortOptionRelevance",
				Plural:    1,
				Value:     "relevance",
			},
		},
	}

	calendar.BeforeDate = coreModel.InputDate{
		Language:        calendar.Language,
		Id:              "before-date",
		InputNameDay:    "before-day",
		InputNameMonth:  "before-month",
		InputNameYear:   "before-year",
		InputValueDay:   "1",
		InputValueMonth: "2",
		InputValueYear:  "2050",
		Title:           helper.Localise("ReleasedBefore", calendar.Language, 1),
		Description:     helper.Localise("DateFilterDescription", calendar.Language, 1),
	}

	calendar.AfterDate = coreModel.InputDate{
		Language:        calendar.Language,
		Id:              "after-date",
		InputNameDay:    "after-day",
		InputNameMonth:  "after-month",
		InputNameYear:   "after-year",
		InputValueDay:   "5",
		InputValueMonth: "6",
		InputValueYear:  "1950",
		Title:           helper.Localise("ReleasedAfter", calendar.Language, 1),
		Description:     helper.Localise("DateFilterDescription", calendar.Language, 1),
	}

	calendar.KeywordSearch = coreModel.CompactSearch{
		ElementId:        "keyword-search",
		InputName:        "keywords",
		Language:         calendar.Language,
		LabelLocaliseKey: "ReleaseCalendarPageSearchKeywords",
		SearchTerm:       "zip zap zoo",
	}

	return calendar
}

func dateChanges(changes []search.ReleaseDateChange) []model.DateChange {
	var modelChanges = make([]model.DateChange, len(changes))
	for i, c := range changes {
		modelChanges[i].Date = c.Date
		modelChanges[i].ChangeNotice = c.ChangeNotice
	}

	return modelChanges
}

func mapReleases(params queryparams.ValidatedParams, response search.ReleaseResponse, language string) map[string]model.ReleaseType {
	checkType := func(given, want queryparams.ReleaseType) bool { return given == want }
	return map[string]model.ReleaseType{
		"type-published": {
			Name:      "release-type",
			Value:     "type-published",
			Id:        "release-type-published",
			LocaleKey: "FilterReleaseTypePublished",
			Plural:    1,
			Language:  language,
			Checked:   checkType(params.ReleaseType, queryparams.Published),
			Count:     response.Breakdown.Published,
		},
		"type-upcoming": {
			Name:      "release-type",
			Value:     "type-upcoming",
			Id:        "release-type-upcoming",
			LocaleKey: "FilterReleaseTypeUpcoming",
			Plural:    1,
			Language:  language,
			Checked:   checkType(params.ReleaseType, queryparams.Upcoming),
			Count:     response.Breakdown.Provisional + response.Breakdown.Confirmed + response.Breakdown.Postponed,
		},
		"type-cancelled": {
			Name:      "release-type",
			Value:     "type-cancelled",
			Id:        "release-type-cancelled",
			LocaleKey: "FilterReleaseTypeCancelled",
			Plural:    1,
			Language:  language,
			Checked:   checkType(params.ReleaseType, queryparams.Cancelled),
			Count:     response.Breakdown.Cancelled,
		},
	}
}

func mapSortOptions(params queryparams.ValidatedParams) []model.SortOption {
	return []model.SortOption{
		{
			LocaleKey: "ReleaseCalendarSortOptionDateNewest",
			Plural:    1,
			Value:     queryparams.RelDateDesc.String(),
			Disabled:  false,
		},
		{
			LocaleKey: "ReleaseCalendarSortOptionDateOldest",
			Plural:    1,
			Value:     queryparams.RelDateAsc.String(),
			Disabled:  false,
		},
		{
			LocaleKey: "ReleaseCalendarSortOptionAlphabeticalAZ",
			Plural:    1,
			Value:     queryparams.TitleAZ.String(),
			Disabled:  false,
		},
		{
			LocaleKey: "ReleaseCalendarSortOptionAlphabeticalZA",
			Plural:    1,
			Value:     queryparams.TitleZA.String(),
			Disabled:  false,
		},
		{
			LocaleKey: "ReleaseCalendarSortOptionRelevance",
			Plural:    1,
			Value:     queryparams.Relevance.String(),
			Disabled: func(keywords string) bool {
				if keywords == "" {
					return true
				}
				return false
			}(params.Keywords),
		},
	}
}

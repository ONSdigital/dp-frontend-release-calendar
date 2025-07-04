package mapper

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/releasecalendar"
	search "github.com/ONSdigital/dp-api-clients-go/v2/site-search"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-release-calendar/config"
	"github.com/ONSdigital/dp-frontend-release-calendar/model"
	"github.com/ONSdigital/dp-frontend-release-calendar/queryparams"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/russross/blackfriday/v2"
)

func createTableOfContents(
	description model.ReleaseDescription,
	relatedDocuments []model.Link,
	relatedDatasets []model.Link,
	relatedAPIDatasets []model.Link,
	dateChanges []model.DateChange,
	aboutTheData bool,
	_ []model.Link,
	relatedMethodology []model.Link,
	relatedMethodologyArticle []model.Link,
	links []model.Link,
	markdown []string,
) coreModel.TableOfContents {
	toc := coreModel.TableOfContents{
		AriaLabel: coreModel.Localisation{
			LocaleKey: "TableOfContents",
			Plural:    1,
		},
		Title: coreModel.Localisation{
			LocaleKey: "Contents",
			Plural:    1,
		},
	}

	sections := make(map[string]coreModel.ContentSection)
	displayOrder := make([]string, 0)

	if description.Summary != "" {
		sections["summary"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionSummary",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "summary")
	}

	if len(relatedDocuments) > 0 {
		sections["publications"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionPublications",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "publications")
	}

	if len(relatedDatasets) > 0 || len(relatedAPIDatasets) > 0 {
		sections["data"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionData",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "data")
	}

	if len(relatedMethodology) > 0 || len(relatedMethodologyArticle) > 0 {
		sections["methodology"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionMethodology",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "methodology")
	}

	if (model.ContactDetails{} != description.Contact) {
		sections["contactdetails"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionContactDetails",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "contactdetails")
	}

	if len(dateChanges) > 0 {
		sections["changestothisreleasedate"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionDateChanges",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "changestothisreleasedate")
	}

	if aboutTheData {
		sections["aboutthedata"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionAboutTheData",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "aboutthedata")
	}

	if len(markdown) > 0 {
		sections["prereleaseaccesslist"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "ReleaseSectionPreReleaseAccessList",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "prereleaseaccesslist")
	}

	if len(links) > 0 {
		sections["links"] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				LocaleKey: "YouMightAlsoBeInterestedIn",
				Plural:    1,
			},
		}
		displayOrder = append(displayOrder, "links")
	}

	toc.Sections = sections
	toc.DisplayOrder = displayOrder

	return toc
}

func mapEmergencyBanner(bannerData zebedee.EmergencyBanner) coreModel.EmergencyBanner {
	var mappedEmergencyBanner coreModel.EmergencyBanner
	emptyBannerObj := zebedee.EmergencyBanner{}
	if bannerData != emptyBannerObj {
		mappedEmergencyBanner.Title = bannerData.Title
		mappedEmergencyBanner.Type = strings.ReplaceAll(bannerData.Type, "_", "-")
		mappedEmergencyBanner.Description = bannerData.Description
		mappedEmergencyBanner.URI = bannerData.URI
		mappedEmergencyBanner.LinkText = bannerData.LinkText
	}
	return mappedEmergencyBanner
}

func createPreGTMJavaScript(title string, description model.ReleaseDescription) []template.JS {
	var releaseStatus string
	var censusTag string

	releaseDate := helper.DateFormatYYYYMMDD(description.ReleaseDate)
	releaseTime := helper.TimeFormat24h(description.ReleaseDate)

	switch {
	case description.Cancelled:
		releaseStatus = "cancelled"
	case description.Published:
		releaseStatus = "published"
	default:
		releaseStatus = "upcoming"
	}

	if description.Census2021 {
		censusTag = "census"
		return []template.JS{
			template.JS(`dataLayer.push({` + //nolint:gosec // input is controlled by app
				`"analyticsOptOut": getUsageCookieValue(),
				"gtm.whitelist": ["google","hjtc","lcl"],
				"gtm.blacklist": ["customScripts","sp","adm","awct","k","d","j"],
				"contentTitle": "` + title + `",
				"release-status": "` + releaseStatus + `",
				"release-date": "` + releaseDate + `",
				"release-time": "` + releaseTime + `",
				"release-date-status": "` + description.ProvisionalDate + `",
				"next-release-date": "` + description.NextRelease + `",
				"contact-name": "` + description.Contact.Name + `",
				"tag": "` + censusTag + `"
			});`),
		}
	}
	return []template.JS{
		template.JS(`dataLayer.push({` + //nolint:gosec // input is controlled by app
			`"analyticsOptOut": getUsageCookieValue(),
			"gtm.whitelist": ["google","hjtc","lcl"],
			"gtm.blacklist": ["customScripts","sp","adm","awct","k","d","j"],
			"contentTitle": "` + title + `",
			"release-status": "` + releaseStatus + `",
			"release-date": "` + releaseDate + `",
			"release-time": "` + releaseTime + `",
			"release-date-status": "` + description.ProvisionalDate + `",
			"next-release-date": "` + description.NextRelease + `",
			"contact-name": "` + description.Contact.Name + `",
		});`),
	}
}

func CreateRelease(cfg config.Config, basePage coreModel.Page, release releasecalendar.Release, lang, path, serviceMessage string, emergencyBannerContent zebedee.EmergencyBanner) model.Release {
	result := model.Release{
		Page:     basePage,
		Markdown: convertMarkdownToHTML(release.Markdown),
		Description: model.ReleaseDescription{
			Title:   release.Description.Title,
			Summary: release.Description.Summary,
			Contact: model.ContactDetails{
				Email:     release.Description.Contact.Email,
				Name:      release.Description.Contact.Name,
				Telephone: release.Description.Contact.Telephone,
			},
			NationalStatistic:  release.Description.NationalStatistic,
			WelshStatistic:     release.Description.WelshStatistic,
			Census2021:         release.Census(),
			ReleaseDate:        release.Description.ReleaseDate,
			NextRelease:        release.Description.NextRelease,
			Published:          release.Description.Published,
			Finalised:          release.Description.Finalised,
			Cancelled:          release.Description.Cancelled,
			CancellationNotice: release.Description.CancellationNotice,
			ProvisionalDate:    release.Description.ProvisionalDate,
		},
	}
	result.FeatureFlags.IsPublishing = cfg.IsPublishing
	result.Language = lang
	result.Type = "releaseCalendar"
	result.ServiceMessage = serviceMessage
	result.EmergencyBanner = mapEmergencyBanner(emergencyBannerContent)
	result.RelatedDatasets = mapLink(release.RelatedDatasets)
	result.RelatedAPIDatasets = mapLink(release.RelatedAPIDatasets)
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

	result.PublicationState = GetPublicationState(result.Description, result.DateChanges)

	result.BetaBannerEnabled = true
	result.Metadata.Title = release.Description.Title
	result.URI = release.URI
	result.AboutTheData = result.Description.NationalStatistic || result.Description.WelshStatistic || result.Description.Census2021

	result.Breadcrumb = mapBreadcrumbTrail(result.Description, result.Language, path)

	result.TableOfContents = createTableOfContents(
		result.Description,
		result.RelatedDocuments,
		result.RelatedDatasets,
		result.RelatedAPIDatasets,
		result.DateChanges,
		result.AboutTheData,
		result.RelatedAPIDatasets,
		result.RelatedMethodology,
		result.RelatedMethodologyArticle,
		result.Links,
		result.Markdown,
	)
	result.PreGTMJavaScript = createPreGTMJavaScript(result.Metadata.Title, result.Description)

	if !result.Description.Finalised && result.Description.ProvisionalDate == "" {
		result.Description.ProvisionalDate = helper.DateTimeOnsDatePatternFormat(result.Description.ReleaseDate, result.Language)
	}
	return result
}

func mapBreadcrumbTrail(description model.ReleaseDescription, language, path string) []coreModel.TaxonomyNode {
	selectState := func(description model.ReleaseDescription) (string, queryparams.ReleaseType) {
		if description.Cancelled {
			return "BreadcrumbCancelled", queryparams.Cancelled
		}

		if description.Published {
			return "BreadcrumbPublished", queryparams.Published
		}

		return "BreadcrumbUpcoming", queryparams.Upcoming
	}

	localeKey, releaseType := selectState(description)

	return []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("BreadcrumbHome", language, 1),
			URI:   "/",
		},
		{
			Title: helper.Localise("BreadcrumbReleaseCalendar", language, 1),
			URI:   path,
		},
		{
			Title: helper.Localise(localeKey, language, 1),
			URI:   fmt.Sprintf("%s?release-type=%s", path, releaseType.String()),
		},
	}
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

func CreateReleaseCalendar(basePage coreModel.Page, params queryparams.ValidatedParams, response search.ReleaseResponse, cfg config.Config, lang, serviceMessage string, emergencyBannerContent zebedee.EmergencyBanner, validationErrs []coreModel.ErrorItem) model.Calendar {
	calendar := model.Calendar{
		Page: basePage,
	}
	calendar.FeatureFlags.IsPublishing = cfg.IsPublishing
	calendar.Language = lang
	calendar.Type = "releaseCalendar"
	calendar.URI = "/releasecalendar"
	calendar.BetaBannerEnabled = true
	calendar.ServiceMessage = serviceMessage
	calendar.EmergencyBanner = mapEmergencyBanner(emergencyBannerContent)
	calendar.Metadata.Title = helper.Localise("ReleaseCalendarPageTitle", calendar.Language, 1)
	calendar.FeatureFlags.FeedbackAPIURL = cfg.FeedbackAPIURL
	calendar.KeywordSearch = coreModel.CompactSearch{
		ElementId: "keyword-search",
		InputName: "keywords",
		Language:  calendar.Language,
		Label: coreModel.Localisation{
			LocaleKey: "ReleaseCalendarPageSearchKeywords",
			Plural:    1,
		},
		SearchTerm: params.Keywords,
	}

	calendar.Sort = model.Sort{
		Mode:    params.Sort.String(),
		Options: mapSortOptions(params),
	}

	itemsPerPage := params.Limit

	totalResults := cfg.DefaultMaximumSearchResults
	if totalResults > response.Breakdown.Total {
		totalResults = response.Breakdown.Total
	}
	currentPage := queryparams.CalculatePageNumber(params.Offset, itemsPerPage)

	calendar.Pagination.TotalPages = queryparams.CalculatePageNumber(totalResults-1, itemsPerPage)
	calendar.Pagination.CurrentPage = currentPage
	calendar.Pagination.Limit = itemsPerPage
	calendar.Pagination.PagesToDisplay = getPagesToDisplay(params, cfg.CalendarPath(), calendar.Pagination.TotalPages, defaultWindowSize)
	calendar.Pagination.FirstAndLastPages = getFirstAndLastPages(params, cfg.CalendarPath(), calendar.Pagination.TotalPages)
	calendar.Pagination.LimitOptions = []int{10, 25}
	calendar.TotalSearchPosition = getTotalSearchPosition(currentPage, itemsPerPage)
	calendar.RSSLink = fmt.Sprintf("releasecalendar?rss&%s", params.AsFrontendQuery().Encode())

	if currentPage > calendar.Pagination.TotalPages {
		validationErrs = append(validationErrs, coreModel.ErrorItem{
			Description: coreModel.Localisation{
				Text: fmt.Sprintf("invalid page parameter: value is above total pages (%d)", calendar.Pagination.TotalPages),
			},
		})
		response = search.ReleaseResponse{}
	}

	for i := range response.Releases {
		calendar.Entries.Items = append(calendar.Entries.Items, calendarEntryFromRelease(response.Releases[i], cfg.RoutingPrefix))
	}

	calendar.Entries.Count = response.Breakdown.Total
	calendar.ReleaseTypes = mapReleases(params, response, calendar.Language)

	var fdErrDescription, tdErrDescription []coreModel.Localisation
	if len(validationErrs) > 0 {
		calendar.Error = coreModel.Error{
			Title:      calendar.Metadata.Title,
			ErrorItems: validationErrs,
			Language:   lang,
		}

		for _, err := range validationErrs {
			switch err.ID {
			case params.AfterDate.GetFieldsetErrID():
				fdErrDescription = append(fdErrDescription, err.Description)
			case params.BeforeDate.GetFieldsetErrID():
				tdErrDescription = append(tdErrDescription, err.Description)
			}
		}
	}

	calendar.AfterDate = coreModel.DateFieldset{
		Language:                 lang,
		ValidationErrDescription: fdErrDescription,
		ErrorID:                  params.AfterDate.GetFieldsetErrID(),
		Input: coreModel.InputDate{
			Language:              lang,
			Id:                    "after-date",
			InputNameDay:          "after-day",
			InputNameMonth:        "after-month",
			InputNameYear:         "after-year",
			InputValueDay:         params.AfterDate.DayString(),
			InputValueMonth:       params.AfterDate.MonthString(),
			InputValueYear:        params.AfterDate.YearString(),
			HasDayValidationErr:   params.AfterDate.HasDayValidationErr(),
			HasMonthValidationErr: params.AfterDate.HasMonthValidationErr(),
			HasYearValidationErr:  params.AfterDate.HasYearValidationErr(),
			DataAttributes: []coreModel.DataAttribute{
				{
					Key: "invalid-date",
					Value: coreModel.Localisation{
						LocaleKey: "ValidationInvalidDate",
						Plural:    1,
					},
				},
			},
			DayDataAttributes: []coreModel.DataAttribute{
				{
					Key: "pattern-mismatch",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationPatternMismatch", lang, 1, "after", "day"),
					},
				},
			},
			MonthDataAttributes: []coreModel.DataAttribute{
				{
					Key: "pattern-mismatch",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationPatternMismatch", lang, 1, "after", "month"),
					},
				},
			},
			YearDataAttributes: []coreModel.DataAttribute{
				{
					Key: "value-missing",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationYearMissing", lang, 1, "after"),
					},
				},
				{
					Key: "pattern-mismatch",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationPatternMismatch", lang, 1, "after", "year"),
					},
				},
			},
			Title: coreModel.Localisation{
				LocaleKey: "ReleasedAfter",
				Plural:    1,
			},
			Description: coreModel.Localisation{
				LocaleKey: "DateFilterAfterDescription",
				Plural:    1,
			},
		},
	}

	calendar.BeforeDate = coreModel.DateFieldset{
		Language:                 lang,
		ValidationErrDescription: tdErrDescription,
		ErrorID:                  params.BeforeDate.GetFieldsetErrID(),
		Input: coreModel.InputDate{
			Language:              lang,
			Id:                    "before-date",
			InputNameDay:          "before-day",
			InputNameMonth:        "before-month",
			InputNameYear:         "before-year",
			InputValueDay:         params.BeforeDate.DayString(),
			InputValueMonth:       params.BeforeDate.MonthString(),
			InputValueYear:        params.BeforeDate.YearString(),
			HasDayValidationErr:   params.BeforeDate.HasDayValidationErr(),
			HasMonthValidationErr: params.BeforeDate.HasMonthValidationErr(),
			HasYearValidationErr:  params.BeforeDate.HasYearValidationErr(),
			Title: coreModel.Localisation{
				LocaleKey: "ReleasedBefore",
				Plural:    1,
			},
			Description: coreModel.Localisation{
				LocaleKey: "DateFilterBeforeDescription",
				Plural:    1,
			},
			DataAttributes: []coreModel.DataAttribute{
				{
					Key: "invalid-range",
					Value: coreModel.Localisation{
						LocaleKey: "ValidationInvalidDateRange",
						Plural:    1,
					},
				},
				{
					Key: "invalid-date",
					Value: coreModel.Localisation{
						LocaleKey: "ValidationInvalidDate",
						Plural:    1,
					},
				},
			},
			DayDataAttributes: []coreModel.DataAttribute{
				{
					Key: "pattern-mismatch",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationPatternMismatch", lang, 1, "before", "day"),
					},
				},
			},
			MonthDataAttributes: []coreModel.DataAttribute{
				{
					Key: "pattern-mismatch",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationPatternMismatch", lang, 1, "before", "month"),
					},
				},
			},
			YearDataAttributes: []coreModel.DataAttribute{
				{
					Key: "value-missing",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationYearMissing", lang, 1, "before"),
					},
				},
				{
					Key: "pattern-mismatch",
					Value: coreModel.Localisation{
						Text: helper.Localise("ValidationPatternMismatch", lang, 1, "before", "year"),
					},
				},
			},
		},
	}

	return calendar
}

// The current page number sits within a window, and the window size determines the
// number of pages around the current page. For example a window size of 3 with the
// current page shown in () would give:
// - (1) 2 3 at the start of the page range
// - 8 9 (10) at the end of the page range
// - 1 ... 5 (6) 7 ... 10 in the middle of the page range
const defaultWindowSize = 5

func getPagesToDisplay(params queryparams.ValidatedParams, path string, totalPages, windowSize int) []coreModel.PageToDisplay {
	start, end := getWindowStartEndPage(params.Page, totalPages, windowSize)

	var pagesToDisplay []coreModel.PageToDisplay
	for i := start; i <= end; i++ {
		pagesToDisplay = append(pagesToDisplay, coreModel.PageToDisplay{
			PageNumber: i,
			URL:        getPageURL(i, params, path),
		})
	}

	return pagesToDisplay
}

func getTotalSearchPosition(currentPage, itemsPerPage int) int {
	totalSearchPosition := (currentPage - 1) * itemsPerPage
	return totalSearchPosition
}

func getFirstAndLastPages(params queryparams.ValidatedParams, path string, totalPages int) []coreModel.PageToDisplay {
	return []coreModel.PageToDisplay{
		{
			PageNumber: 1,
			URL:        getPageURL(1, params, path),
		},
		{
			PageNumber: totalPages,
			URL:        getPageURL(totalPages, params, path),
		},
	}
}

// getWindowStartEndPage calculates the start and end page of the moving window of size windowSize, over the set of pages
// whose current page is currentPage, and whose size is totalPages
// It is an error to pass a parameter whose value is < 1, or a currentPage > totalPages, and the function will panic in this case
func getWindowStartEndPage(currentPage, totalPages, windowSize int) (start, end int) {
	if totalPages < 1 || windowSize < 1 {
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
	start = currentPage - windowOffset
	switch {
	case start <= 0:
		start = 1
	case start > totalPages-windowSize+1:
		start = totalPages - windowSize + 1
	}

	end = start + windowSize - 1
	if end > totalPages {
		end = totalPages
	}

	return start, end
}

func getPageURL(page int, params queryparams.ValidatedParams, path string) (pageURL string) {
	query := params.AsFrontendQuery()
	query.Set("page", strconv.Itoa(page))

	return path + "?" + query.Encode()
}

func getWindowOffset(windowSize int) int {
	if windowSize%2 == 0 {
		return (windowSize / 2) - 1
	}

	return windowSize / 2
}

func calendarEntryFromRelease(release search.Release, uriPrivatePrefix string) model.CalendarEntry {
	result := model.CalendarEntry{
		URI:         uriPrivatePrefix + release.URI,
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

	result.PublicationState = GetPublicationState(result.Description, result.DateChanges)

	return result
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
	generateLabel := func(localeKey, language string, plural, count int) string {
		if count > 0 {
			return fmt.Sprintf("%s (%d)", helper.Localise(localeKey, language, plural), count)
		}
		return helper.Localise(localeKey, language, plural)
	}
	return map[string]model.ReleaseType{
		"type-published": {
			Name:  "release-type",
			Value: "type-published",
			ID:    "release-type-published",
			Label: coreModel.Localisation{
				Text: generateLabel("FilterReleaseTypePublished", language, 1, response.Breakdown.Published),
			},
			Language:  language,
			IsChecked: checkType(params.ReleaseType, queryparams.Published),
			Count:     response.Breakdown.Published,
		},
		"type-upcoming": {
			Name:  "release-type",
			Value: "type-upcoming",
			ID:    "release-type-upcoming",
			Label: coreModel.Localisation{
				Text: generateLabel("FilterReleaseTypeUpcoming", language, 1, response.Breakdown.Provisional+response.Breakdown.Confirmed+response.Breakdown.Postponed),
			},
			Language:  language,
			IsChecked: checkType(params.ReleaseType, queryparams.Upcoming),
			Count:     response.Breakdown.Provisional + response.Breakdown.Confirmed + response.Breakdown.Postponed,
		},
		"type-cancelled": {
			Name:  "release-type",
			Value: "type-cancelled",
			ID:    "release-type-cancelled",
			Label: coreModel.Localisation{
				Text: generateLabel("FilterReleaseTypeCancelled", language, 1, response.Breakdown.Cancelled),
			},
			Language:  language,
			IsChecked: checkType(params.ReleaseType, queryparams.Cancelled),
			Count:     response.Breakdown.Cancelled,
		},
		"type-census": {
			Name:  "census",
			Value: "true",
			ID:    "release-type-census",
			Label: coreModel.Localisation{
				Text: generateLabel("FilterReleaseTypeCensus", language, 1, response.Breakdown.Census),
			},
			Language:  language,
			IsChecked: params.Census,
			Count:     response.Breakdown.Census,
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
				return keywords == ""
			}(params.Keywords),
		},
	}
}

func convertMarkdownToHTML(markdowns []string) []string {
	markdownHTML := []string{}
	for _, markdown := range markdowns {
		html := blackfriday.Run([]byte(markdown))
		markdownHTML = append(markdownHTML, string(html))
	}
	return markdownHTML
}

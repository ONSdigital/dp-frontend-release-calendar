package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type Link struct {
	Title   string `json:"title"`
	URI     string `json:"uri"`
	Summary string `json:"summary"`
	// Description string `json:"description"`
}

type ContactDetails struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

// type DateChange struct {
// 	PreviousDate    string `json:"date"`
// 	ReasonForChange string `json:"reason_for_change"`
// }

// type StatusFlags struct {
// 	NotYetPublished           bool `json:"not_yet_published"`
// 	InsufficientDataCancelled bool `json:"insufficient_data_cancelled"`
// 	InsufficientDataPostponed bool `json:"insufficient_data_postponed"`
// 	MergeCancelled            Link `json:"merge_cancelled"`
// }

// type StatusDates struct {
// 	Provisional bool   `json:"provisional"`
// 	Cancelled   bool   `json:"cancelled"`
// 	ReleaseDate string `json:"release_date"`
// 	NextDate    string `json:"next_date"`
// }

type Release struct {
	coreModel.Page
	Markdown                  []string           `json:"markdown"`
	RelatedDocuments          []Link             `json:"related_documents"`
	RelatedDatasets           []Link             `json:"related_datasets"`
	RelatedMethodology        []Link             `json:"related_methodology"`
	RelatedMethodologyArticle []Link             `json:"related_methodology_article"`
	Links                     []Link             `json:"links"`
	DateChanges               []DateChange       `json:"date_changes"`
	Description               ReleaseDescription `json:"description"`
	// ContactDetails            ContactDetails `json:"contact_details"`
	// ReleaseHistory            []Link         `json:"release_history"`
	// StatusFlags               StatusFlags    `json:"status_flags"`
	// StatusDates               StatusDates    `json:"status_dates"`
}

type DateChange struct {
	ChangeNotice string `json:"change_notice"`
	Date         string `json:"previous_date"`
}

type ReleaseDescription struct {
	CancellationNotice []string       `json:"cancellation_notice"`
	Cancelled          bool           `json:"cancelled"`
	Contact            ContactDetails `json:"contact"`
	Finalised          bool           `json:"finalised"`
	NationalStatistic  bool           `json:"national_statistic"`
	NextRelease        string         `json:"next_release"`
	ProvisionalDate    string         `json:"provisional_date"`
	Published          bool           `json:"published"`
	ReleaseDate        string         `json:"release_date"`
	Summary            string         `json:"summary"`
	Title              string         `json:"title"`
}

type PreviousReleases struct {
	coreModel.Page
	Markdown       string `json:"markdown"`
	ReleaseHistory []Link `json:"release_history"`
}

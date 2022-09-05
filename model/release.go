package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type Link struct {
	Title   string `json:"title"`
	URI     string `json:"uri"`
	Summary string `json:"summary"`
}

type ContactDetails struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

// PublicationState represents Types 'cancelled', 'published', or 'upcoming'
// SubTypes of 'upcoming' are 'confirmed', 'postponed', or 'provisional'
type PublicationState struct {
	Type    string `json:"type"`
	SubType string `json:"sub_type"`
}

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
	ReleaseHistory            []Link             `json:"release_history"`
	AboutTheData              bool               `json:"about_the_data"`
	PublicationState          PublicationState   `json:"publication_state"`
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
	WelshStatistic     bool           `json:"welsh_statistic"`
	Census2021         bool           `json:"census_2021"`
	NextRelease        string         `json:"next_release"`
	ProvisionalDate    string         `json:"provisional_date"`
	Published          bool           `json:"published"`
	ReleaseDate        string         `json:"release_date"`
	Summary            string         `json:"summary"`
	Title              string         `json:"title"`
}

type PreviousReleases struct {
	coreModel.Page
	Description    ReleaseDescription `json:"description"`
	ReleaseHistory []Link             `json:"release_history"`
}

// FuncGetPostponementReason Gets the most recent postponement reason, if available
func (release Release) FuncGetPostponementReason() string {
	reason := ""
	totalDateChanges := len(release.DateChanges)

	if totalDateChanges > 0 {
		reason = release.DateChanges[totalDateChanges-1].ChangeNotice
	}

	return reason
}

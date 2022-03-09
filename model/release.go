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
	// TODO Provisional entry for modelling history
	// ReleaseHistory            []Link         `json:"release_history"`
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

// TODO Provisional model for previous releases
type PreviousReleases struct {
	coreModel.Page
	Description    ReleaseDescription `json:"description"`
	ReleaseHistory []Link             `json:"release_history"`
}

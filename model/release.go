package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type Link struct {
	Title string `json:"title"`
	URI   string `json:"uri"`
	Index int    `json:"index"`
}

type ContactDetails struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type Release struct {
	coreModel.Page
	Markdown                  string         `json:"markdown"`
	RelatedDocuments          []Link         `json:"related_documents"`
	RelatedDatasets           []Link         `json:"related_datasets"`
	RelatedMethodology        []Link         `json:"related_methodology"`
	RelatedMethodologyArticle []Link         `json:"related_methodology_article"`
	Links                     []Link         `json:"links"`
	ContactDetails            ContactDetails `json:"contact_details"`
}

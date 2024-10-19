package core

import (
	"github.com/mdw-tools/hugoinho/contracts"
)

type MetadataValidationHandler struct {
	slugs map[string]struct{}
}

func NewMetadataValidationHandler() *MetadataValidationHandler {
	return &MetadataValidationHandler{slugs: make(map[string]struct{})}
}

func (this *MetadataValidationHandler) Handle(article *contracts.Article) {
	if article.Metadata.Title == "" {
		article.Error = StackTraceError(errBlankMetadataTitle)
		return
	}

	if article.Metadata.Slug == "" {
		article.Error = StackTraceError(errBlankMetadataSlug)
		return
	}

	if article.Metadata.Date.IsZero() {
		article.Error = StackTraceError(errBlankMetadataDate)
		return
	}

	_, found := this.slugs[article.Metadata.Slug]
	if found {
		article.Error = StackTraceError(errRepeatedMetadataSlug)
		return
	}

	this.slugs[article.Metadata.Slug] = struct{}{}
}

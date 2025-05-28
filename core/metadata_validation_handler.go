package core

import (
	"fmt"

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
		article.Error = fmt.Errorf("[%s] %w", article.Source.Path, StackTraceError(errBlankMetadataTitle))
		return
	}

	if article.Metadata.Slug == "" {
		article.Error = fmt.Errorf("[%s] %w", article.Source.Path, StackTraceError(errBlankMetadataSlug))
		return
	}

	if article.Metadata.Date.IsZero() {
		article.Error = fmt.Errorf("[%s] %w", article.Source.Path, StackTraceError(errBlankMetadataDate))
		return
	}

	_, found := this.slugs[article.Metadata.Slug]
	if found {
		article.Error = fmt.Errorf("[%s] %w", article.Source.Path, StackTraceError(errRepeatedMetadataSlug))
		return
	}

	this.slugs[article.Metadata.Slug] = struct{}{}
}

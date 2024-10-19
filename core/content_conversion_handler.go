package core

import (
	"github.com/mdw-tools/hugoinho/contracts"
)

type ContentConversionHandler struct {
	inner converter
}

type converter interface {
	Convert(content string) (string, error)
}

func NewContentConversionHandler(inner converter) *ContentConversionHandler {
	return &ContentConversionHandler{inner: inner}
}

func (this *ContentConversionHandler) Handle(article *contracts.Article) {
	_, original := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	converted, err := this.inner.Convert(original)
	if err != nil {
		article.Error = StackTraceError(err)
		return
	}

	article.Content.Original = original
	article.Content.Converted = converted
}

package core

import (
	"fmt"
	"strings"

	"github.com/mdw-tools/hugoinho/contracts"
)

type MetadataParsingHandler struct{}

func NewMetadataParsingHandler() *MetadataParsingHandler {
	return &MetadataParsingHandler{}
}

func (this *MetadataParsingHandler) Handle(article *contracts.Article) {
	if strings.TrimSpace(article.Source.Data) == "" {
		article.Error = StackTraceError(errMissingMetadata)
		return
	}

	metadata, _ := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	if len(metadata) == 0 {
		article.Error = StackTraceError(errMissingMetadataDivider)
		return
	}

	parser := NewMetadataParser(strings.Split(metadata, "\n"))
	err := parser.Parse()
	if err != nil {
		article.Error = fmt.Errorf("[%s] %w", article.Source.Path, err)
		return
	}

	article.Metadata = parser.Parsed()
}

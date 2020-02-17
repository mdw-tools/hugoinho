package core

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
)

type MetadataParsingHandler struct{}

func NewMetadataParsingHandler() *MetadataParsingHandler {
	return &MetadataParsingHandler{}
}

func (this *MetadataParsingHandler) Handle(article *contracts.Article) error {
	if strings.TrimSpace(article.Source.Data) == "" {
		return NewStackTraceError(errMissingMetadata)
	}

	metadata, _ := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	if len(metadata) == 0 {
		return NewStackTraceError(errMissingMetadataDivider)
	}

	parser := NewMetadataParser(strings.Split(metadata, "\n"))
	err := parser.Parse()
	if err != nil {
		return err
	}

	article.Metadata = parser.Parsed()
	return nil
}

type MetadataParser struct {
	lines  []string
	parsed contracts.ArticleMetadata

	parsedTitle bool
	parsedIntro bool
	parsedSlug  bool
	parsedDraft bool
	parsedDate  bool
	parsedTags  bool
}

func NewMetadataParser(lines []string) *MetadataParser {
	return &MetadataParser{lines: lines}
}

func (this *MetadataParser) Parse() error {
	for _, line := range this.lines {
		key, value := divide(line, ":")

		switch key {
		case "title":
			err := this.parseTitle(value)
			if err != nil {
				return err
			}
		case "intro":
			err := this.parseIntro(value)
			if err != nil {
				return err
			}
		case "slug":
			err := this.parseSlug(value)
			if err != nil {
				return err
			}
		case "draft":
			err := this.parseDraft(value)
			if err != nil {
				return err
			}
		case "date":
			err := this.parseDate(value)
			if err != nil {
				return err
			}
		case "tags":
			err := this.parseTags(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *MetadataParser) parseTitle(value string) error {
	if this.parsedTitle {
		return NewStackTraceError(errDuplicateMetadataTitle)
	}
	if value == "" {
		return errBlankMetadataTitle
	}
	this.parsed.Title = value
	this.parsedTitle = true
	return nil
}
func (this *MetadataParser) parseIntro(value string) error {
	if this.parsedIntro {
		return NewStackTraceError(errDuplicateMetadataIntro)
	}
	if value == "" {
		return NewStackTraceError(errBlankMetadataIntro)
	}
	this.parsed.Intro = value
	this.parsedIntro = true
	return nil
}
func (this *MetadataParser) parseSlug(value string) error {
	if this.parsedSlug {
		return NewStackTraceError(errDuplicateMetadataSlug)
	}
	if value == "" {
		return NewStackTraceError(errBlankMetadataSlug)
	}
	if strings.ToLower(value) != value {
		return NewStackTraceError(errInvalidMetadataSlug)
	}
	parsed, _ := url.Parse(value)
	if parsed.Path != parsed.EscapedPath() {
		return NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataSlug, value))
	}
	this.parsed.Slug = value
	this.parsedSlug = true
	return nil
}
func (this *MetadataParser) parseDraft(value string) error {
	if this.parsedDraft {
		return NewStackTraceError(errDuplicateMetadataDraft)
	}

	switch value {
	case "true":
		this.parsed.Draft = true
		this.parsedDraft = true
	case "false":
		this.parsed.Draft = false
		this.parsedDraft = true
	case "":
		return NewStackTraceError(errBlankMetadataDraft)
	default:
		return NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataDraft, value))
	}
	return nil
}

func (this *MetadataParser) parseDate(value string) error {
	if this.parsedDate {
		return NewStackTraceError(errDuplicateMetadataDate)
	}
	if value == "" {
		return NewStackTraceError(errBlankMetadataDate)
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return NewStackTraceError(fmt.Errorf("%w with value: [%s] err: %v", errInvalidMetadataDate, value, err))
	}
	this.parsed.Date = parsed
	this.parsedDate = true
	return nil
}

func (this *MetadataParser) parseTags(value string) error {
	if this.parsedTags {
		return NewStackTraceError(errDuplicateMetadataTags)
	}
	if value == "" {
		return NewStackTraceError(errBlankMetadataTags)
	}
	tags := strings.Fields(value)
	for _, tag := range tags {
		if !isValidTag(tag) {
			return NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataTags, value))
		}
	}
	this.parsed.Tags = tags
	this.parsedTags = true
	return nil
}

func isValidTag(tag string) bool {
	for _, c := range tag {
		if !(isSpace(c) || isDash(c) || isNumber(c) || isLowerAlpha(c)) {
			return false
		}
	}
	return true
}

func (this *MetadataParser) Parsed() contracts.ArticleMetadata {
	return this.parsed
}

var (
	errMissingMetadata        = errors.New("article lacks metadata")
	errMissingMetadataDivider = errors.New("article lacks metadata divider")

	errDuplicateMetadataTitle = errors.New("duplicate metadata title")
	errDuplicateMetadataIntro = errors.New("duplicate metadata intro")
	errDuplicateMetadataSlug  = errors.New("duplicate metadata slug")
	errDuplicateMetadataDraft = errors.New("duplicate metadata draft")
	errDuplicateMetadataDate  = errors.New("duplicate metadata date")
	errDuplicateMetadataTags  = errors.New("duplicate metadata tags")

	errInvalidMetadataSlug  = errors.New("invalid metadata slug")
	errInvalidMetadataDraft = errors.New("invalid metadata draft")
	errInvalidMetadataDate  = errors.New("invalid metadata date")
	errInvalidMetadataTags  = errors.New("invalid metadata tags")

	errBlankMetadataSlug  = errors.New("blank metadata slug")
	errBlankMetadataDraft = errors.New("blank metadata draft")
	errBlankMetadataTitle = errors.New("blank metadata title")
	errBlankMetadataIntro = errors.New("blank metadata intro")
	errBlankMetadataDate  = errors.New("blank metadata date")
	errBlankMetadataTags  = errors.New("blank metadata tags")
)
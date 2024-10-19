package core

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mdw-tools/hugoinho/contracts"
)

type MetadataParser struct {
	lines  []string
	parsed contracts.ArticleMetadata

	parsedTitle  bool
	parsedIntro  bool
	parsedSlug   bool
	parsedDraft  bool
	parsedDate   bool
	parsedTopics bool
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
		case "topics":
			err := this.parseTopics(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *MetadataParser) parseTitle(value string) error {
	if this.parsedTitle {
		return StackTraceError(errDuplicateMetadataTitle)
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
		return StackTraceError(errDuplicateMetadataIntro)
	}
	this.parsed.Intro = value
	this.parsedIntro = true
	return nil
}
func (this *MetadataParser) parseSlug(value string) error {
	if this.parsedSlug {
		return StackTraceError(errDuplicateMetadataSlug)
	}
	if value == "" {
		return StackTraceError(errBlankMetadataSlug)
	}
	if strings.ToLower(value) != value {
		return StackTraceError(errInvalidMetadataSlug)
	}
	parsed, _ := url.Parse(value)
	if parsed.Path != parsed.EscapedPath() {
		return StackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataSlug, value))
	}
	this.parsed.Slug = value
	this.parsedSlug = true
	return nil
}
func (this *MetadataParser) parseDraft(value string) error {
	if this.parsedDraft {
		return StackTraceError(errDuplicateMetadataDraft)
	}

	switch value {
	case "true":
		this.parsed.Draft = true
		this.parsedDraft = true
	case "false":
		this.parsed.Draft = false
		this.parsedDraft = true
	case "":
		return StackTraceError(errBlankMetadataDraft)
	default:
		return StackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataDraft, value))
	}
	return nil
}
func (this *MetadataParser) parseDate(value string) error {
	if this.parsedDate {
		return StackTraceError(errDuplicateMetadataDate)
	}
	if value == "" {
		return StackTraceError(errBlankMetadataDate)
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return StackTraceError(fmt.Errorf("%w with value: [%s] err: %v", errInvalidMetadataDate, value, err))
	}
	this.parsed.Date = parsed
	this.parsedDate = true
	return nil
}
func (this *MetadataParser) parseTopics(value string) error {
	if this.parsedTopics {
		return StackTraceError(errDuplicateMetadataTopics)
	}
	unique := make(map[string]struct{})
	topics := strings.Fields(value)
	for _, topic := range topics {
		if !isValidTopic(topic) {
			return StackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataTopics, value))
		}
		unique[topic] = struct{}{}
	}
	if len(unique) != len(topics) {
		return StackTraceError(fmt.Errorf("%w: [%s] (repeated values)", errInvalidMetadataTopics, value))
	}
	this.parsed.Topics = topics
	this.parsedTopics = true
	return nil
}

func isValidTopic(topic string) bool {
	for _, c := range topic {
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

	errDuplicateMetadataTitle  = errors.New("duplicate metadata title")
	errDuplicateMetadataIntro  = errors.New("duplicate metadata intro")
	errDuplicateMetadataSlug   = errors.New("duplicate metadata slug")
	errDuplicateMetadataDraft  = errors.New("duplicate metadata draft")
	errDuplicateMetadataDate   = errors.New("duplicate metadata date")
	errDuplicateMetadataTopics = errors.New("duplicate metadata topics")

	errInvalidMetadataSlug   = errors.New("invalid metadata slug")
	errInvalidMetadataDraft  = errors.New("invalid metadata draft")
	errInvalidMetadataDate   = errors.New("invalid metadata date")
	errInvalidMetadataTopics = errors.New("invalid metadata topics")

	errRepeatedMetadataSlug = errors.New("repeated metadata slug")

	errBlankMetadataSlug  = errors.New("blank metadata slug")
	errBlankMetadataDraft = errors.New("blank metadata draft")
	errBlankMetadataTitle = errors.New("blank metadata title")
	errBlankMetadataDate  = errors.New("blank metadata date")
)

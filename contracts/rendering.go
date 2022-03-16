package contracts

import (
	"errors"
	"time"
)

type Renderer interface {
	Render(any) (string, error)
}

const (
	HomePageTemplateName = "home.tmpl"
	ArticleTemplateName  = "article.tmpl"
	TopicsTemplateName   = "topics.tmpl"
)

var (
	ErrUnsupportedRenderingType = errors.New("unsupported rendering type")
	ErrRenderingFailure         = errors.New("failed to render template")
)

type (
	RenderedHomePage struct {
		Pages []RenderedArticleSummary
	}

	RenderedArticle struct {
		Slug    string
		Title   string
		Intro   string
		Date    time.Time
		Topics  []string
		Content string
	}

	RenderedArticleSummary struct {
		Slug   string
		Title  string
		Intro  string
		Date   time.Time
		Topics []string
		Draft  bool
	}

	RenderedTopicsListing struct {
		Topics []RenderedTopicListing
	}

	RenderedTopicListing struct {
		Topic    string
		Articles []RenderedArticleSummary
	}
)

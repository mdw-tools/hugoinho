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
	ArchivesTemplateName = "archives.tmpl"
	ArticleTemplateName  = "article.tmpl"
	TopicsTemplateName   = "topics.tmpl"

	// FeedTemplateName is optional; the feed is only rendered when a site defines it.
	FeedTemplateName = "feed.tmpl"
)

var (
	ErrUnsupportedRenderingType = errors.New("unsupported rendering type")
	ErrRenderingFailure         = errors.New("failed to render template")
	ErrMissingTemplate          = errors.New("missing template")
)

type (
	RenderedHomePage struct {
		ProminentTopics []string
		Pages           []RenderedArticleSummary
	}

	RenderedArchivesPage struct {
		Pages []RenderedArticleSummary
	}

	RenderedArticle struct {
		Slug      string
		Title     string
		Intro     string
		Date      time.Time
		Topics    []string
		Canonical string
		Content   string
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

	// RenderedFeed holds every article (newest first) for a syndication feed (RSS, Atom, etc.).
	// Unlike RenderedArticle (as rendered by the article template), the Content is not marked
	// as safe HTML, so templates emit it escaped, which is what XML feeds require.
	RenderedFeed struct {
		Articles []RenderedArticle
	}
)

package core

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mdw-tools/hugoinho/contracts"
)

// FeedRenderingHandler renders all articles, newest first, into feed.xml at the
// root of the output folder. Rendering the feed is optional: when the site does
// not define a feed template, nothing is written.
type FeedRenderingHandler struct {
	articles []contracts.RenderedArticle
	renderer contracts.Renderer
	disk     RenderingFileSystem
	output   string
}

func NewFeedRenderingHandler(renderer contracts.Renderer, disk RenderingFileSystem, output string) *FeedRenderingHandler {
	return &FeedRenderingHandler{renderer: renderer, disk: disk, output: output}
}

func (this *FeedRenderingHandler) Handle(article *contracts.Article) {
	this.articles = append(this.articles, contracts.RenderedArticle{
		Slug:      article.Metadata.Slug,
		Title:     article.Metadata.Title,
		Intro:     article.Metadata.Intro,
		Date:      article.Metadata.Date,
		Topics:    article.Metadata.Topics,
		Canonical: article.Metadata.Canonical,
		Content:   article.Content.Converted,
	})
}

func (this *FeedRenderingHandler) Finalize() error {
	if len(this.articles) == 0 {
		return nil
	}
	slices.SortStableFunc(this.articles, func(a, b contracts.RenderedArticle) int {
		if byDate := -a.Date.Compare(b.Date); byDate != 0 {
			return byDate
		}
		return strings.Compare(a.Slug, b.Slug)
	})

	rendered, err := this.renderer.Render(contracts.RenderedFeed{Articles: this.articles})
	if errors.Is(err, contracts.ErrMissingTemplate) {
		return nil
	}
	if err != nil {
		return err
	}

	err = this.disk.MkdirAll(this.output, 0755)
	if err != nil {
		return err
	}
	return this.disk.WriteFile(filepath.Join(this.output, "feed.xml"), []byte(rendered), 0644)
}

package core

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestFeedRenderingHandlerSuite(t *testing.T) {
	suite.Run(&FeedRenderingHandlerSuite{T: suite.New(t)}, suite.Options.UnitTests())
}

type FeedRenderingHandlerSuite struct {
	*suite.T

	handler  *FeedRenderingHandler
	renderer *FakeRenderer
	disk     *InMemoryFileSystem
}

func (this *FeedRenderingHandlerSuite) Setup() {
	this.renderer = NewFakeRenderer()
	this.disk = NewInMemoryFileSystem()
	this.handler = NewFeedRenderingHandler(this.renderer, this.disk, "output/folder")
}
func feedArticle(slug string, day int) *contracts.Article {
	return &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:      slug,
			Title:     "Title " + slug,
			Intro:     "Intro " + slug,
			Topics:    []string{"topic"},
			Date:      Date(2024, 1, day),
			Canonical: "https://example.com" + slug,
		},
		Content: contracts.ArticleContent{Converted: fmt.Sprintf("<p>%s</p>", slug)},
	}
}
func renderedFeedArticle(slug string, day int) contracts.RenderedArticle {
	return contracts.RenderedArticle{
		Slug:      slug,
		Title:     "Title " + slug,
		Intro:     "Intro " + slug,
		Topics:    []string{"topic"},
		Date:      Date(2024, 1, day),
		Canonical: "https://example.com" + slug,
		Content:   fmt.Sprintf("<p>%s</p>", slug),
	}
}
func (this *FeedRenderingHandlerSuite) handleAndFinalize() error {
	this.handler.Handle(feedArticle("/b", 1))
	this.handler.Handle(feedArticle("/c", 3))
	this.handler.Handle(feedArticle("/a", 1))
	return this.handler.Finalize()
}

func (this *FeedRenderingHandlerSuite) TestNoArticles_NothingRendered() {
	err := this.handler.Finalize()
	this.So(err, should.BeNil)
	this.So(this.renderer.all, should.BeEmpty)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *FeedRenderingHandlerSuite) TestArticlesRenderedNewestFirst_WrittenToFeedXML() {
	this.renderer.result = "RENDERED"

	err := this.handleAndFinalize()

	this.So(err, should.BeNil)
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedFeed{
		Articles: []contracts.RenderedArticle{
			renderedFeedArticle("/c", 3),
			renderedFeedArticle("/a", 1), // same date: ordered by slug
			renderedFeedArticle("/b", 1),
		},
	})
	this.So(this.disk.Files, should.Contain, "output/folder/feed.xml")
	this.So(this.disk.Files["output/folder/feed.xml"].Content(), should.Equal, "RENDERED")
}
func (this *FeedRenderingHandlerSuite) TestMissingFeedTemplate_NothingWritten_NoError() {
	this.renderer.err = fmt.Errorf("%w: feed.tmpl", contracts.ErrMissingTemplate)

	err := this.handleAndFinalize()

	this.So(err, should.BeNil)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *FeedRenderingHandlerSuite) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, renderErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *FeedRenderingHandlerSuite) TestMkdirAllErrorReturned() {
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder"] = mkdirErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, mkdirErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *FeedRenderingHandlerSuite) TestWriteFileErrorReturned() {
	writeErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/feed.xml"] = writeErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, writeErr)
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/feed.xml")
}

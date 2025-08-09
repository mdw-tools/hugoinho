package core

import (
	"errors"
	"testing"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestTopicPageRenderingHandlerFixture(t *testing.T) {
	suite.Run(&TopicPageRenderingHandlerFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type TopicPageRenderingHandlerFixture struct {
	*suite.T

	handler  *TopicPageRenderingHandler
	disk     *InMemoryFileSystem
	renderer *FakeRenderer
}

func (this *TopicPageRenderingHandlerFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.renderer = NewFakeRenderer()
	this.handler = NewTopicPageRenderingHandler(this.disk, this.renderer, "output/folder")
	this.handleArticles()
}

func (this *TopicPageRenderingHandlerFixture) handleArticles() {
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug1",
			Title:  "title1",
			Intro:  "intro1",
			Date:   Date(2020, 1, 1),
			Topics: []string{"a", "b"},
		},
	})
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug2",
			Title:  "title2",
			Intro:  "intro2",
			Date:   Date(2020, 2, 2),
			Topics: []string{"b", "c"},
		},
	})
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug3",
			Title:  "title3",
			Intro:  "intro3",
			Date:   Date(2020, 3, 3),
			Topics: []string{"c"},
		},
	})
}

func (this *TopicPageRenderingHandlerFixture) assertHandledArticlesRendered() {
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedTopicsListing{
		Topics: []contracts.RenderedTopicListing{
			{
				Topic: "b",
				Articles: []contracts.RenderedArticleSummary{
					{
						Slug:  "/slug2",
						Title: "title2",
						Intro: "intro2",
						Date:  Date(2020, 2, 2),
					},
					{
						Slug:  "/slug1",
						Title: "title1",
						Intro: "intro1",
						Date:  Date(2020, 1, 1),
					},
				},
			},
			{
				Topic: "c",
				Articles: []contracts.RenderedArticleSummary{
					{
						Slug:  "/slug3",
						Title: "title3",
						Intro: "intro3",
						Date:  Date(2020, 3, 3),
					},
					{
						Slug:  "/slug2",
						Title: "title2",
						Intro: "intro2",
						Date:  Date(2020, 2, 2),
					},
				},
			},
			//{Topic: "a"}, // This topic only has one article and will be omitted from the listing.
		},
	})
}

func (this *TopicPageRenderingHandlerFixture) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handler.Finalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.Contain, "output/folder")
	if this.So(this.disk.Files, should.Contain, "output/folder/topics/index.html") {
		file := this.disk.Files["output/folder/topics/index.html"]
		this.So(file.Content(), should.Equal, "RENDERED")
	}
}

func (this *TopicPageRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handler.Finalize()

	this.So(err, should.WrapError, renderErr)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *TopicPageRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder/topics"] = mkdirErr

	err := this.handler.Finalize()

	this.So(err, should.WrapError, mkdirErr)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *TopicPageRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/topics/index.html"] = writeFileErr

	err := this.handler.Finalize()

	this.So(err, should.WrapError, writeFileErr)
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/topics/index.html")
}

package core

import (
	"errors"
	"testing"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestArticleRenderingHandlerFixture(t *testing.T) {
	should.Run(&ArticleRenderingHandlerFixture{T: should.New(t)}, should.Options.UnitTests())
}

type ArticleRenderingHandlerFixture struct {
	*should.T
	handler  *ArticleRenderingHandler
	renderer *FakeRenderer
	disk     *InMemoryFileSystem
	article  *contracts.Article
}

func (this *ArticleRenderingHandlerFixture) Setup() {
	this.renderer = NewFakeRenderer()
	this.disk = NewInMemoryFileSystem()
	this.handler = NewArticleRenderingHandler(this.disk, this.renderer, "output/folder")

	this.article = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft:  false,
			Slug:   "/slug",
			Title:  "Title",
			Intro:  "Intro",
			Topics: []string{"a", "b"},
			Date:   Date(2020, 2, 8),
		},
		Content: contracts.ArticleContent{
			Converted: "CONTENT",
		},
	}
}

func (this *ArticleRenderingHandlerFixture) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	this.handler.Handle(this.article)

	this.So(this.article.Error, should.BeNil)
	this.assertArticleDataRendered()
	this.So(this.disk.Files, should.Contain, "output/folder/slug")
	if !this.So(this.disk.Files, should.Contain, "output/folder/slug/index.html") {
		file := this.disk.Files["output/folder/slug/index.html"]
		this.So(file.Content(), should.Equal, "RENDERED")
	}
}

func (this *ArticleRenderingHandlerFixture) assertArticleDataRendered() {
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedArticle{
		Slug:    this.article.Metadata.Slug,
		Title:   this.article.Metadata.Title,
		Intro:   this.article.Metadata.Intro,
		Date:    this.article.Metadata.Date,
		Topics:  this.article.Metadata.Topics,
		Content: this.article.Content.Converted,
	})
}

func (this *ArticleRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	this.handler.Handle(this.article)

	this.So(this.article.Error, should.WrapError, renderErr)
	this.assertArticleDataRendered()
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *ArticleRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder/slug"] = mkdirErr

	this.handler.Handle(this.article)

	this.So(this.article.Error, should.WrapError, mkdirErr)
	this.assertArticleDataRendered()
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *ArticleRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/slug/index.html"] = writeFileErr

	this.handler.Handle(this.article)

	this.So(this.article.Error, should.WrapError, writeFileErr)
	this.assertArticleDataRendered()
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/slug/index.html")
}

///////////////////////////////////////////////////////////////////

type FakeRenderer struct {
	all      []any
	rendered any
	result   string
	err      error
}

func NewFakeRenderer() *FakeRenderer {
	return &FakeRenderer{}
}

func (this *FakeRenderer) Render(rendered any) (string, error) {
	this.all = append(this.all, rendered)
	this.rendered = rendered
	return this.result, this.err
}

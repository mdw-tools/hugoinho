package core

import (
	"errors"
	"strings"
	"testing"

	"github.com/mdw-go/testing/v2/better"
	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestListRenderingHandlerSuite(t *testing.T) {
	suite.Run(&ListRenderingHandlerSuite{T: suite.New(t)}, suite.Options.UnitTests())
}

type ListRenderingHandlerSuite struct {
	*suite.T

	handler  *ListRenderingHandler
	renderer *FakeRenderer
	disk     *InMemoryFileSystem
}

var (
	articleA = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft:  false,
			Slug:   "/a",
			Title:  "A",
			Intro:  "aa",
			Topics: []string{"topic-a"},
			Date:   Date(2023, 7, 7),
		},
	}
	articleB = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft:  true,
			Slug:   "/b",
			Title:  "B",
			Intro:  "bb",
			Topics: []string{"topic-b"},
			Date:   Date(2023, 7, 8),
		},
	}
	articleC = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft:  false,
			Slug:   "/c",
			Title:  "C",
			Intro:  "cc",
			Topics: []string{"topic-c"},
			Date:   Date(2023, 7, 9),
		},
	}
)

func (this *ListRenderingHandlerSuite) filter(article *contracts.Article) bool {
	return article.Metadata.Title < "C"
}
func (this *ListRenderingHandlerSuite) sorter(i, j contracts.RenderedArticleSummary) int {
	return strings.Compare(i.Title, j.Title)
}
func (this *ListRenderingHandlerSuite) assertHandledArticlesRendered() {
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedListPage{
		Title: "TITLE",
		Pages: []contracts.RenderedArticleSummary{
			{
				Slug:   "/a",
				Title:  "A",
				Intro:  "aa",
				Date:   Date(2023, 7, 7),
				Topics: []string{"topic-a"},
				Draft:  false,
			},
			{
				Slug:   "/b",
				Title:  "B",
				Intro:  "bb",
				Date:   Date(2023, 7, 8),
				Topics: []string{"topic-b"},
				Draft:  true,
			},
		},
	})
}
func (this *ListRenderingHandlerSuite) Setup() {
	this.renderer = NewFakeRenderer()
	this.disk = NewInMemoryFileSystem()
	this.handler = NewListRenderingHandler(this.filter, this.sorter, this.renderer, this.disk, "output/folder", "TITLE")
}
func (this *ListRenderingHandlerSuite) handleAndFinalize() error {
	this.handler.Handle(articleC)
	this.handler.Handle(articleB)
	this.handler.Handle(articleA)
	return this.handler.Finalize()
}
func (this *ListRenderingHandlerSuite) TestNoArticles_NothingToRender() {
	this.handler.Handle(articleC) // will be filtered out
	err := this.handler.Finalize()
	this.So(err, should.BeNil)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *ListRenderingHandlerSuite) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handleAndFinalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.Contain, "output/folder")
	this.So(this.disk.Files, better.Contain, "output/folder/index.html")
	file := this.disk.Files["output/folder/index.html"]
	this.So(file.Content(), should.Equal, "RENDERED")
}
func (this *ListRenderingHandlerSuite) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, renderErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *ListRenderingHandlerSuite) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder"] = mkdirErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, mkdirErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *ListRenderingHandlerSuite) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/index.html"] = writeFileErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, writeFileErr)
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/index.html")
}

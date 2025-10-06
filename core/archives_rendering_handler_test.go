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

func TestArchivesRenderingHandlerSuite(t *testing.T) {
	suite.Run(&ArchivesRenderingHandlerSuite{T: suite.New(t)}, suite.Options.UnitTests())
}

type ArchivesRenderingHandlerSuite struct {
	*suite.T

	handler  *ArchivesRenderingHandler
	renderer *FakeRenderer
	disk     *InMemoryFileSystem
}

func (this *ArchivesRenderingHandlerSuite) filter(article *contracts.Article) bool {
	return article.Metadata.Title < "C"
}
func (this *ArchivesRenderingHandlerSuite) sorter(i, j contracts.RenderedArticleSummary) int {
	return strings.Compare(i.Title, j.Title)
}
func (this *ArchivesRenderingHandlerSuite) assertHandledArticlesRendered() {
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedArchivesPage{
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
			{
				Slug:   "/b/2",
				Title:  "B2",
				Intro:  "bb",
				Date:   Date(2023, 7, 8),
				Topics: []string{"topic-b"},
				Draft:  true,
			},
		},
	})
}
func (this *ArchivesRenderingHandlerSuite) Setup() {
	this.renderer = NewFakeRenderer()
	this.disk = NewInMemoryFileSystem()
	this.handler = NewArchivesRenderingHandler(this.filter, this.sorter, this.renderer, this.disk, "output/folder")
}
func (this *ArchivesRenderingHandlerSuite) handleAndFinalize() error {
	this.handler.Handle(articleA)
	this.handler.Handle(articleB)
	this.handler.Handle(articleB2)
	this.handler.Handle(articleC)
	return this.handler.Finalize()
}
func (this *ArchivesRenderingHandlerSuite) TestNoArticles_NothingToRender() {
	this.handler.Handle(articleC) // will be filtered out
	err := this.handler.Finalize()
	this.So(err, should.BeNil)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *ArchivesRenderingHandlerSuite) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handleAndFinalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.Contain, "output/folder")
	this.So(this.disk.Files, better.Contain, "output/folder/archives/index.html")
	file := this.disk.Files["output/folder/archives/index.html"]
	this.So(file.Content(), should.Equal, "RENDERED")
}
func (this *ArchivesRenderingHandlerSuite) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, renderErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *ArchivesRenderingHandlerSuite) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder/archives"] = mkdirErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, mkdirErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *ArchivesRenderingHandlerSuite) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/archives/index.html"] = writeFileErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, writeFileErr)
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/archives.html")
}

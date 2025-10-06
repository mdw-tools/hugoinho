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
	suite.Run(&HomepageRenderingHandlerSuite{T: suite.New(t)}, suite.Options.UnitTests())
}

type HomepageRenderingHandlerSuite struct {
	*suite.T

	handler  *HomepageRenderingHandler
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
	articleB2 = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft:  true,
			Slug:   "/b/2",
			Title:  "B2",
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

func (this *HomepageRenderingHandlerSuite) filter(article *contracts.Article) bool {
	return article.Metadata.Title < "C"
}
func (this *HomepageRenderingHandlerSuite) sorter(i, j contracts.RenderedArticleSummary) int {
	return strings.Compare(i.Title, j.Title)
}
func (this *HomepageRenderingHandlerSuite) assertHandledArticlesRendered() {
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedHomePage{
		ProminentTopics: []string{"topic-b", "topic-a"},
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
func (this *HomepageRenderingHandlerSuite) Setup() {
	this.renderer = NewFakeRenderer()
	this.disk = NewInMemoryFileSystem()
	this.handler = NewHomepageRenderingHandler(this.filter, this.sorter, this.renderer, this.disk, "output/folder")
}
func (this *HomepageRenderingHandlerSuite) handleAndFinalize() error {
	this.handler.Handle(articleA)
	this.handler.Handle(articleB)
	this.handler.Handle(articleB2)
	this.handler.Handle(articleC)
	return this.handler.Finalize()
}
func (this *HomepageRenderingHandlerSuite) TestNoArticles_NothingToRender() {
	this.handler.Handle(articleC) // will be filtered out
	err := this.handler.Finalize()
	this.So(err, should.BeNil)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *HomepageRenderingHandlerSuite) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handleAndFinalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.Contain, "output/folder")
	this.So(this.disk.Files, better.Contain, "output/folder/index.html")
	file := this.disk.Files["output/folder/index.html"]
	this.So(file.Content(), should.Equal, "RENDERED")
}
func (this *HomepageRenderingHandlerSuite) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, renderErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *HomepageRenderingHandlerSuite) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder"] = mkdirErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, mkdirErr)
	this.So(this.disk.Files, should.BeEmpty)
}
func (this *HomepageRenderingHandlerSuite) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/index.html"] = writeFileErr

	err := this.handleAndFinalize()

	this.So(err, should.WrapError, writeFileErr)
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/index.html")
}

func TestLeaderboardSuite(t *testing.T) {
	suite.Run(&LeaderboardSuiteSuite{T: suite.New(t)}, suite.Options.UnitTests())
}

type LeaderboardSuiteSuite struct {
	*suite.T
}

func (this *LeaderboardSuiteSuite) Test() {
	l := leaderboard[string]{
		"a": 42,
		"b": 43,
		"c": 44,
		"d": 45,
	}
	this.So(l.TopN(3), should.Equal, []string{"d", "c", "b"})
}

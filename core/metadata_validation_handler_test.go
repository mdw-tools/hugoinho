package core

import (
	"testing"
	"time"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestMetadataValidationHandlerFixture(t *testing.T) {
	suite.Run(&MetadataValidationHandlerFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type MetadataValidationHandlerFixture struct {
	*suite.T

	handler *MetadataValidationHandler
	article *contracts.Article
}

func (this *MetadataValidationHandlerFixture) Setup() {
	this.handler = NewMetadataValidationHandler()
	this.article = &contracts.Article{
		Source: contracts.ArticleSource{
			Path: "/the/article/path",
		},
		Metadata: contracts.ArticleMetadata{
			Draft:  false,
			Slug:   "/slug1",
			Title:  "Title",
			Intro:  "Introduction",
			Topics: []string{"a", "b", "c"},
			Date:   Date(2020, 2, 2),
		},
	}
}

func (this *MetadataValidationHandlerFixture) TestAllPresentAndAccountedFor() {
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.BeNil)
}
func (this *MetadataValidationHandlerFixture) TestMissingTitle_Err() {
	this.article.Metadata.Title = ""
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.WrapError, errBlankMetadataTitle)
	this.So(this.article.Error.Error(), should.Contain, this.article.Source.Path)
}
func (this *MetadataValidationHandlerFixture) TestMissingSlug_Err() {
	this.article.Metadata.Slug = ""
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.WrapError, errBlankMetadataSlug)
	this.So(this.article.Error.Error(), should.Contain, this.article.Source.Path)
}
func (this *MetadataValidationHandlerFixture) TestMissingDate_Err() {
	this.article.Metadata.Date = time.Time{}
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.WrapError, errBlankMetadataDate)
	this.So(this.article.Error.Error(), should.Contain, this.article.Source.Path)
}
func (this *MetadataValidationHandlerFixture) TestUniqueSlugs_OK() {
	this.assertHandleWithSlugOK("a")
	this.assertHandleWithSlugOK("b")
	this.assertHandleWithSlugOK("c")
}
func (this *MetadataValidationHandlerFixture) TestRepeatedSlugs_Err() {
	this.assertHandleWithSlugOK("A")
	this.assertHandleWithSlugOK("b")
	this.assertHandleWithSlugOK("c")
	this.assertHandleWithSlugFAIL("A")
}
func (this *MetadataValidationHandlerFixture) assertHandleWithSlugOK(slug string) {
	this.article.Metadata.Slug = slug
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.BeNil)
}
func (this *MetadataValidationHandlerFixture) assertHandleWithSlugFAIL(slug string) {
	this.article.Metadata.Slug = slug
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.WrapError, errRepeatedMetadataSlug)
}

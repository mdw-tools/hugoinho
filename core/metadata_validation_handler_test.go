package core

import (
	"testing"
	"time"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestMetadataValidationHandlerFixture(t *testing.T) {
	should.Run(&MetadataValidationHandlerFixture{T: should.New(t)}, should.Options.UnitTests())
}

type MetadataValidationHandlerFixture struct {
	*should.T

	handler *MetadataValidationHandler
	article *contracts.Article
}

func (this *MetadataValidationHandlerFixture) Setup() {
	this.handler = NewMetadataValidationHandler()
	this.article = &contracts.Article{
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
}
func (this *MetadataValidationHandlerFixture) TestMissingSlug_Err() {
	this.article.Metadata.Slug = ""
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.WrapError, errBlankMetadataSlug)
}
func (this *MetadataValidationHandlerFixture) TestMissingDate_Err() {
	this.article.Metadata.Date = time.Time{}
	this.handler.Handle(this.article)
	this.So(this.article.Error, should.WrapError, errBlankMetadataDate)
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

package core

import (
	"errors"
	"testing"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/testing/should"
)

func TestContentConversionHandlerFixture(t *testing.T) {
	should.Run(&ContentConversionHandlerFixture{T: should.New(t)}, should.Options.UnitTests())
}

type ContentConversionHandlerFixture struct {
	*should.T

	converter *ContentConversionHandler
	inner     *FakeConverter
	input     chan contracts.Article
	output    chan contracts.Article
}

func (this *ContentConversionHandlerFixture) Setup() {
	this.inner = NewFakeConverter()
	this.converter = NewContentConversionHandler(this.inner)
}

func (this *ContentConversionHandlerFixture) formatSourceData(original string) string {
	return "\n" + contracts.METADATA_CONTENT_DIVIDER + "\n" + original
}

func (this *ContentConversionHandlerFixture) TestValidContentParsedAndConverted() {
	article := &contracts.Article{Source: contracts.ArticleSource{Data: this.formatSourceData("content1")}}

	this.converter.Handle(article)

	this.So(article, should.Equal, &contracts.Article{
		Error:   nil,
		Source:  contracts.ArticleSource{Data: this.formatSourceData("content1")},
		Content: contracts.ArticleContent{Original: "content1", Converted: "content1 (CONVERTED)"},
	})
}

func (this *ContentConversionHandlerFixture) TestInvalidContentElicitsError() {
	article := &contracts.Article{Source: contracts.ArticleSource{Data: this.formatSourceData("content1")}}
	conversionError := errors.New("conversion error")
	this.inner.err = conversionError

	this.converter.Handle(article)

	this.So(article.Error, should.WrapError, conversionError)
	this.So(article.Content, should.Equal, contracts.ArticleContent{})
}

////////////////////////////////////////////////////

type FakeConverter struct {
	err error
}

func NewFakeConverter() *FakeConverter {
	return &FakeConverter{}
}

func (this *FakeConverter) Convert(content string) (string, error) {
	return content + " (CONVERTED)", this.err
}

package core

import (
	"fmt"
	"testing"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestListenerFixture(t *testing.T) {
	suite.Run(&ListenerFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type ListenerFixture struct {
	*suite.T
	input   chan contracts.Article
	output  chan contracts.Article
	handler *FakeHandler
}

func (this *ListenerFixture) Setup() {
	this.input = make(chan contracts.Article, 10)
	this.output = make(chan contracts.Article, 10)
	this.handler = NewFakeHandler()
}

func (this *ListenerFixture) TestEachArticleHandledIfNotErrantAndPassedOn() {
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "A"}}
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "B"}, Error: contracts.ErrDroppedArticle}
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "C"}}
	close(this.input)

	Listen(this.input, this.output, this.handler)

	this.So(gather(this.output), should.Equal, []contracts.Article{
		{Content: contracts.ArticleContent{Original: "A", Converted: "A1"}},
		{Content: contracts.ArticleContent{Original: "B", Converted: ""}, Error: contracts.ErrDroppedArticle},
		{Content: contracts.ArticleContent{Original: "C", Converted: "C2"}},
	})
}

func (this *ListenerFixture) TestFinalizeCalledIfDefinedOnHandler() {
	close(this.input)

	handler := NewFakeFinalizingHandler()

	Listen(this.input, this.output, handler)

	this.So(handler.called, should.Equal, 1)
}

func (this *ListenerFixture) TestFinalizeErrPassedOnIfNonNil() {
	close(this.input)

	handler := NewFakeFinalizingHandler()
	handler.err = contracts.ErrDroppedArticle

	Listen(this.input, this.output, handler)

	this.So(handler.called, should.Equal, 1)
	this.So(len(this.output), should.Equal, 1)
	this.So(<-this.output, should.Equal, contracts.Article{Error: contracts.ErrDroppedArticle})
}

///////////////////////////////////////////////////////////////

type FakeHandler struct {
	calls int
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{}
}

func (this *FakeHandler) Handle(article *contracts.Article) {
	this.calls++
	article.Content.Converted = article.Content.Original + fmt.Sprint(this.calls)
}

//////////////////////////////////////////////////////////////

type FakeFinalizingHandler struct {
	called int
	err    error
}

func (this *FakeFinalizingHandler) Handle(*contracts.Article) {
	panic("NOT NEEDED")
}

func NewFakeFinalizingHandler() *FakeFinalizingHandler {
	return &FakeFinalizingHandler{}
}

func (this *FakeFinalizingHandler) Finalize() error {
	this.called++
	return this.err
}

//////////////////////////////////////////////////////////////

func gather(output chan contracts.Article) (pages []contracts.Article) {
	for page := range output {
		pages = append(pages, page)
	}
	return pages
}

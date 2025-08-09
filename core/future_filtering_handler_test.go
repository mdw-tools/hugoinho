package core

import (
	"testing"
	"time"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestFutureFilteringHandlerFixture(t *testing.T) {
	suite.Run(&FutureFilteringHandlerFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type FutureFilteringHandlerFixture struct {
	*suite.T

	present time.Time
	past    time.Time
	future  time.Time
}

func (this *FutureFilteringHandlerFixture) Setup() {
	this.present = time.Now()
	this.past = this.present.Add(-time.Second)
	this.future = this.present.Add(time.Second)
}

func (this *FutureFilteringHandlerFixture) article(date time.Time) *contracts.Article {
	return &contracts.Article{Metadata: contracts.ArticleMetadata{Date: date}}
}

func (this *FutureFilteringHandlerFixture) TestDisabled_LetEverythingThrough() {
	disabled := NewFutureFilteringHandler(this.present, false)

	past := this.article(this.past)
	disabled.Handle(past)
	this.So(past.Error, should.BeNil)

	present := this.article(this.present)
	disabled.Handle(present)
	this.So(present.Error, should.BeNil)

	future := this.article(this.future)
	disabled.Handle(future)
	this.So(future.Error, should.BeNil)
}

func (this *FutureFilteringHandlerFixture) TestEnabled_AnythingAfterNowDropped() {
	enabled := NewFutureFilteringHandler(this.present, true)

	past := this.article(this.past)
	enabled.Handle(past)
	this.So(past.Error, should.BeNil)

	present := this.article(this.present)
	enabled.Handle(present)
	this.So(present.Error, should.BeNil)

	future := this.article(this.future)
	enabled.Handle(future)
	this.So(future.Error, should.WrapError, contracts.ErrDroppedArticle)
}

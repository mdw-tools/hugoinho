package core

import (
	"testing"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestDraftFilteringHandlerFixture(t *testing.T) {
	should.Run(&DraftFilteringHandlerFixture{T: should.New(t)}, should.Options.UnitTests())
}

type DraftFilteringHandlerFixture struct {
	*should.T
}

func (this *DraftFilteringHandlerFixture) article(draft bool) *contracts.Article {
	return &contracts.Article{Metadata: contracts.ArticleMetadata{Draft: draft}}
}

func (this *DraftFilteringHandlerFixture) TestDisabled_LetEverythingThrough() {
	handler := NewDraftFilteringHandler(false)

	draft := this.article(true)
	handler.Handle(draft)
	this.So(draft.Error, should.BeNil)

	nonDraft := this.article(false)
	handler.Handle(nonDraft)
	this.So(nonDraft.Error, should.BeNil)
}

func (this *DraftFilteringHandlerFixture) TestEnabled_AnyDraftsDropped() {
	handler := NewDraftFilteringHandler(true)

	nonDraft := this.article(false)
	handler.Handle(nonDraft)
	this.So(nonDraft.Error, should.BeNil)

	draft := this.article(true)
	handler.Handle(draft)
	this.So(draft.Error, should.WrapError, contracts.ErrDroppedArticle)
}

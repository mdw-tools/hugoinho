package core

import (
	"fmt"
	"time"

	"github.com/mdw-tools/hugoinho/contracts"
)

type FutureFilteringHandler struct {
	now     time.Time
	enabled bool
}

func NewFutureFilteringHandler(now time.Time, enabled bool) *FutureFilteringHandler {
	return &FutureFilteringHandler{now: now, enabled: enabled}
}

func (this *FutureFilteringHandler) Handle(article *contracts.Article) {
	if !this.enabled {
		return
	}
	if !article.Metadata.Date.After(this.now) {
		return
	}
	article.Error = fmt.Errorf(
		"%w: %s (can be published on %s)",
		contracts.ErrDroppedArticle,
		article.Metadata.Slug,
		article.Metadata.Date.Format("January 2, 2006"),
	)
}

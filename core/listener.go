package core

import (
	"github.com/mdw-tools/hugoinho/contracts"
)

func Listen(in, out chan contracts.Article, handler contracts.Handler) {
	defer close(out)
	defer finalize(handler, out)

	for article := range in {
		if article.Error == nil {
			handler.Handle(&article)
		}
		out <- article
	}
}

func finalize(handler contracts.Handler, out chan contracts.Article) {
	finalizer, ok := handler.(contracts.Finalizer)
	if !ok {
		return
	}

	err := finalizer.Finalize()
	if err == nil {
		return
	}

	out <- contracts.Article{Error: err}
}

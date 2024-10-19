package core

import (
	"strings"

	"github.com/mdw-tools/hugoinho/contracts"
)

type PathLoader struct {
	files  contracts.Walk
	root   string
	output chan contracts.Article
	err    error
}

func NewPathLoader(
	files contracts.Walk,
	root string,
	output chan contracts.Article,
) *PathLoader {
	return &PathLoader{
		files:  files,
		root:   root,
		output: output,
	}
}

func (this *PathLoader) Start() {
	defer close(this.output)
	files := this.files.Walk(this.root)
	for file := range files {
		if file.Error != nil {
			this.err = StackTraceError(file.Error)
			return
		}
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		this.output <- contracts.Article{
			Source: contracts.ArticleSource{Path: file.Path},
		}
	}
}

func (this *PathLoader) Finalize() error {
	return this.err
}

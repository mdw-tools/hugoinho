package core

import "github.com/mdw-tools/hugoinho/contracts"

type FileReadingHandler struct {
	disk contracts.ReadFile
}

func NewFileReadingHandler(disk contracts.ReadFile) *FileReadingHandler {
	return &FileReadingHandler{disk: disk}
}

func (this *FileReadingHandler) Handle(article *contracts.Article) {
	raw, err := this.disk.ReadFile(article.Source.Path)
	if err != nil {
		article.Error = StackTraceError(err)
	} else {
		article.Source.Data = string(raw)
	}
}

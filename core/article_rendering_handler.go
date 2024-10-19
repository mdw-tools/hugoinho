package core

import (
	"path/filepath"

	"github.com/mdw-tools/hugoinho/contracts"
)

type RenderingFileSystem interface {
	contracts.MkdirAll
	contracts.WriteFile
}

type ArticleRenderingHandler struct {
	disk     RenderingFileSystem
	renderer contracts.Renderer
	output   string
}

func NewArticleRenderingHandler(
	disk RenderingFileSystem,
	renderer contracts.Renderer,
	output string,
) *ArticleRenderingHandler {
	return &ArticleRenderingHandler{
		disk:     disk,
		renderer: renderer,
		output:   output,
	}
}

func (this *ArticleRenderingHandler) Handle(article *contracts.Article) {
	data := contracts.RenderedArticle{
		Slug:    article.Metadata.Slug,
		Title:   article.Metadata.Title,
		Intro:   article.Metadata.Intro,
		Date:    article.Metadata.Date,
		Topics:  article.Metadata.Topics,
		Content: article.Content.Converted,
	}

	rendered, err := this.renderer.Render(data)
	if err != nil {
		article.Error = StackTraceError(err)
		return
	}

	folder := filepath.Join(this.output, article.Metadata.Slug)
	err = this.disk.MkdirAll(folder, 0755)
	if err != nil {
		article.Error = StackTraceError(err)
		return
	}

	err = this.disk.WriteFile(filepath.Join(folder, "index.html"), []byte(rendered), 0644)
	if err != nil {
		article.Error = StackTraceError(err)
		return
	}
}

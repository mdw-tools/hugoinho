package core

import (
	"path/filepath"
	"slices"

	"github.com/mdw-tools/hugoinho/contracts"
)

type ArchivesRenderingHandler struct {
	pages    []contracts.RenderedArticleSummary
	filter   contracts.Filter
	sorter   contracts.Sorter
	renderer contracts.Renderer
	disk     RenderingFileSystem
	output   string
}

func NewArchivesRenderingHandler(
	filter contracts.Filter,
	sorter contracts.Sorter,
	renderer contracts.Renderer,
	disk RenderingFileSystem,
	output string,
) *ArchivesRenderingHandler {
	return &ArchivesRenderingHandler{
		filter:   filter,
		sorter:   sorter,
		renderer: renderer,
		disk:     disk,
		output:   output,
	}
}
func (this *ArchivesRenderingHandler) Handle(article *contracts.Article) {
	if !this.filter(article) {
		return
	}
	this.pages = append(this.pages, contracts.RenderedArticleSummary{
		Slug:   article.Metadata.Slug,
		Title:  article.Metadata.Title,
		Intro:  article.Metadata.Intro,
		Date:   article.Metadata.Date,
		Topics: article.Metadata.Topics,
		Draft:  article.Metadata.Draft,
	})
}
func (this *ArchivesRenderingHandler) Finalize() error {
	if len(this.pages) == 0 {
		return nil
	}
	rendered, err := this.renderer.Render(contracts.RenderedArchivesPage{
		Pages: slices.SortedStableFunc(slices.Values(this.pages), this.sorter),
	})
	if err != nil {
		return StackTraceError(err)
	}

	folder := filepath.Join(this.output, "archives")
	err = this.disk.MkdirAll(folder, 0755)
	if err != nil {
		return StackTraceError(err)
	}

	err = this.disk.WriteFile(filepath.Join(folder, "index.html"), []byte(rendered), 0644)
	if err != nil {
		return StackTraceError(err)
	}

	return nil
}

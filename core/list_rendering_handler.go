package core

import (
	"path/filepath"
	"sort"

	"github.com/mdw-tools/hugoinho/contracts"
)

type ListRenderingHandler struct {
	listing  []contracts.RenderedArticleSummary
	filter   contracts.Filter
	sorter   contracts.Sorter
	renderer contracts.Renderer
	disk     RenderingFileSystem
	output   string
	title    string
}

func NewListRenderingHandler(
	filter contracts.Filter,
	sorter contracts.Sorter,
	renderer contracts.Renderer,
	disk RenderingFileSystem,
	output, title string,
) *ListRenderingHandler {
	return &ListRenderingHandler{
		filter:   filter,
		sorter:   sorter,
		renderer: renderer,
		disk:     disk,
		output:   output,
		title:    title,
	}
}
func (this *ListRenderingHandler) Handle(article *contracts.Article) {
	if !this.filter(article) {
		return
	}
	this.listing = append(this.listing, contracts.RenderedArticleSummary{
		Slug:   article.Metadata.Slug,
		Title:  article.Metadata.Title,
		Intro:  article.Metadata.Intro,
		Date:   article.Metadata.Date,
		Topics: article.Metadata.Topics,
		Draft:  article.Metadata.Draft,
	})
}
func (this *ListRenderingHandler) Finalize() error {
	if len(this.listing) == 0 {
		return nil
	}

	sort.SliceStable(this.listing, func(i, j int) bool {
		return this.sorter(this.listing[i], this.listing[j]) < 0
	})

	rendered, err := this.renderer.Render(contracts.RenderedListPage{
		Title: this.title,
		Pages: this.listing,
	})
	if err != nil {
		return StackTraceError(err)
	}

	err = this.disk.MkdirAll(this.output, 0755)
	if err != nil {
		return StackTraceError(err)
	}

	err = this.disk.WriteFile(filepath.Join(this.output, "index.html"), []byte(rendered), 0644)
	if err != nil {
		return StackTraceError(err)
	}

	return nil
}

package core

import (
	"cmp"
	"maps"
	"path/filepath"
	"slices"

	"github.com/mdw-tools/hugoinho/contracts"
)

type HomepageRenderingHandler struct {
	pages    []contracts.RenderedArticleSummary
	topics   leaderboard[string]
	filter   contracts.Filter
	sorter   contracts.Sorter
	renderer contracts.Renderer
	disk     RenderingFileSystem
	output   string
}

func NewHomepageRenderingHandler(
	filter contracts.Filter,
	sorter contracts.Sorter,
	renderer contracts.Renderer,
	disk RenderingFileSystem,
	output string,
) *HomepageRenderingHandler {
	return &HomepageRenderingHandler{
		filter:   filter,
		sorter:   sorter,
		renderer: renderer,
		disk:     disk,
		output:   output,
		topics:   make(leaderboard[string]),
	}
}
func (this *HomepageRenderingHandler) Handle(article *contracts.Article) {
	if !this.filter(article) {
		return
	}
	for topic := range slices.Values(article.Metadata.Topics) {
		this.topics[topic]++
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
func (this *HomepageRenderingHandler) Finalize() error {
	if len(this.pages) == 0 {
		return nil
	}
	rendered, err := this.renderer.Render(contracts.RenderedHomePage{
		ProminentTopics: this.topics.TopN(30),
		Pages:           slices.SortedStableFunc(slices.Values(this.pages), this.sorter)[:min(len(this.pages), 10)],
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

type leaderboard[T cmp.Ordered] map[T]int

func (this leaderboard[T]) compare(i, j T) int {
	rank := -cmp.Compare(this[i], this[j])
	if rank == 0 {
		return cmp.Compare(i, j)
	}
	return rank
}
func (this leaderboard[T]) TopN(n int) (result []T) {
	return slices.SortedStableFunc(maps.Keys(this), this.compare)[:min(n, len(this))]
}

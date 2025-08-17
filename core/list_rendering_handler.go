package core

import (
	"cmp"
	"maps"
	"path/filepath"
	"slices"

	"github.com/mdw-tools/hugoinho/contracts"
)

type ListRenderingHandler struct {
	listing  []contracts.RenderedArticleSummary
	topics   leaderboard[string]
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
		topics:   make(map[string]int),
	}
}
func (this *ListRenderingHandler) Handle(article *contracts.Article) {
	if !this.filter(article) {
		return
	}
	for topic := range slices.Values(article.Metadata.Topics) {
		this.topics[topic]++
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
	this.listing = slices.SortedFunc(slices.Values(this.listing), this.sorter)

	rendered, err := this.renderer.Render(contracts.RenderedListPage{
		Title:           this.title,
		LatestArticle:   this.listing[0],
		ProminentTopics: this.topics.TopN(30),
		Pages:           this.listing,
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

type leaderboard[T comparable] map[T]int

func (this leaderboard[T]) TopN(n int) (result []T) {
	return take(n, slices.SortedStableFunc(maps.Keys(this), func(i, j T) int {
		return -cmp.Compare(this[i], this[j])
	}))
}
func take[T any](i int, s []T) []T {
	if len(s) > i {
		return s[:i]
	}
	return s
}

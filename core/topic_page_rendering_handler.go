package core

import (
	"path/filepath"
	"slices"
	"sort"

	"github.com/mdw-tools/hugoinho/contracts"
)

type TopicPageRenderingHandler struct {
	disk     RenderingFileSystem
	renderer contracts.Renderer
	output   string
	topics   map[string][]contracts.RenderedArticleSummary
}

func NewTopicPageRenderingHandler(
	disk RenderingFileSystem,
	renderer contracts.Renderer,
	output string,
) *TopicPageRenderingHandler {
	return &TopicPageRenderingHandler{
		disk:     disk,
		renderer: renderer,
		output:   output,
		topics:   make(map[string][]contracts.RenderedArticleSummary),
	}
}

func (this *TopicPageRenderingHandler) Handle(article *contracts.Article) {
	for _, topic := range article.Metadata.Topics {
		this.topics[topic] = append(this.topics[topic], contracts.RenderedArticleSummary{
			Slug:  article.Metadata.Slug,
			Title: article.Metadata.Title,
			Intro: article.Metadata.Intro,
			Date:  article.Metadata.Date,
		})
	}
}

func (this *TopicPageRenderingHandler) Finalize() error {
	rendered, err := this.renderer.Render(this.prepareRendering())
	if err != nil {
		return StackTraceError(err)
	}

	folder := filepath.Join(this.output, "topics")
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

func (this *TopicPageRenderingHandler) prepareRendering() (full contracts.RenderedTopicsListing) {
	for _, topic := range this.sortTopics() {
		articles := this.topics[topic]
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].Date.After(articles[j].Date)
		})
		full.Topics = append(full.Topics, contracts.RenderedTopicListing{
			Topic:    topic,
			Articles: articles,
		})
	}
	return full
}

func (this *TopicPageRenderingHandler) sortTopics() (topics []string) {
	for topic := range this.topics {
		if len(this.topics[topic]) == 1 {
			continue
		}
		topics = append(topics, topic)
	}
	slices.Sort(topics)
	return topics
}

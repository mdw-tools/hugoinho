package core

import "github.com/mdw-tools/hugoinho/contracts"

type Pipeline struct {
	clock    contracts.Clock
	config   contracts.Config
	disk     contracts.FileSystem
	renderer contracts.Renderer
}

func NewPipeline(
	clock contracts.Clock,
	config contracts.Config,
	disk contracts.FileSystem,
	renderer contracts.Renderer,
) *Pipeline {
	return &Pipeline{
		clock:    clock,
		config:   config,
		disk:     disk,
		renderer: renderer,
	}
}
func (this *Pipeline) Run() (out chan contracts.Article) {
	out = this.goLoad()
	out = this.goListen(out, NewFileReadingHandler(this.disk))
	out = this.goListen(out, NewMetadataParsingHandler())
	out = this.goListen(out, NewMetadataValidationHandler())
	out = this.goListen(out, NewDraftFilteringHandler(!this.config.BuildDrafts))
	out = this.goListen(out, NewFutureFilteringHandler(this.clock(), !this.config.BuildFuture))
	out = this.goListen(out, NewContentConversionHandler(NewGoldmarkMarkdownConverter()))
	out = this.goListen(out, NewArticleRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, NewTopicPageRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, NewArchivesRenderingHandler(
		filterAll,
		sortByDateDescending,
		this.renderer,
		this.disk,
		this.config.TargetRoot,
	))
	out = this.goListen(out, NewHomepageRenderingHandler(
		filterAll,
		sortByDateDescending,
		this.renderer,
		this.disk,
		this.config.TargetRoot,
	))
	return out
}
func (this *Pipeline) goLoad() (out chan contracts.Article) {
	out = make(chan contracts.Article)
	go NewPathLoader(this.disk, this.config.ContentRoot, out).Start()
	return out
}
func (this *Pipeline) goListen(in chan contracts.Article, handler contracts.Handler) (out chan contracts.Article) {
	out = make(chan contracts.Article)
	go Listen(in, out, handler)
	return out
}
func filterAll(*contracts.Article) bool { return true }
func sortByDateDescending(i, j contracts.RenderedArticleSummary) int {
	return -i.Date.Compare(j.Date)
}

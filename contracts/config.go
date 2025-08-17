package contracts

type Config struct {
	TemplateDir string
	ContentRoot string
	TargetRoot  string
	BasePath    string
	BuildDrafts bool
	BuildFuture bool
}

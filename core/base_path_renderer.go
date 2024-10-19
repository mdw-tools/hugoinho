package core

import (
	"fmt"
	"regexp"

	"github.com/mdw-tools/hugoinho/contracts"
)

type BasePathRenderer struct {
	inner    contracts.Renderer
	pattern  *regexp.Regexp
	basePath string
}

func NewBasePathRenderer(inner contracts.Renderer, basePath string) *BasePathRenderer {
	return &BasePathRenderer{
		inner:    inner,
		pattern:  regexp.MustCompile(`href="/([^/])`),
		basePath: fmt.Sprintf(`href="%s/$1`, basePath),
	}
}
func (this *BasePathRenderer) Render(v any) (string, error) {
	output, err := this.inner.Render(v)
	return this.pattern.ReplaceAllString(output, this.basePath), err
}

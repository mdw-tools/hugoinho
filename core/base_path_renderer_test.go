package core

import (
	"strings"
	"testing"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestBasePathRendererSuite(t *testing.T) {
	suite.Run(&BasePathRendererSuite{T: suite.New(t)}, suite.Options.UnitTests())
}

type BasePathRendererSuite struct {
	*suite.T
	inner    *FakeRenderer
	renderer contracts.Renderer
}

func (this *BasePathRendererSuite) Setup() {
	this.inner = NewFakeRenderer()
	this.renderer = NewBasePathRenderer(this.inner, "/base-path")
}

func (this *BasePathRendererSuite) TestRenderer() {
	this.inner.result = strings.Join([]string{
		`<a href="/yes">link</a>`,
		`<a href="//no">link</a>`,
		`<a href="https://no">link</a>`,
	}, "\n")

	output, err := this.renderer.Render(42)

	this.So(this.inner.all, should.Equal, []any{42})
	this.So(err, should.BeNil)
	this.So(output, should.Equal, strings.Join([]string{
		`<a href="/base-path/yes">link</a>`,
		`<a href="//no">link</a>`,
		`<a href="https://no">link</a>`,
	}, "\n"))
}

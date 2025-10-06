package core

import (
	"testing"
	"text/template"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestTemplateRendererFixture(t *testing.T) {
	suite.Run(&TemplateRendererFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type TemplateRendererFixture struct {
	*suite.T

	templates *template.Template
	renderer  *TemplateRenderer
}

func (this *TemplateRendererFixture) Setup() {
	this.parseTemplate(contracts.HomePageTemplateName)
	this.parseTemplate(contracts.ArchivesTemplateName)
	this.parseTemplate(contracts.ArticleTemplateName)
	this.parseTemplate(contracts.TopicsTemplateName)
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.BeNil)
}
func (this *TemplateRendererFixture) parseTemplate(name string) {
	var err error
	if this.templates == nil {
		this.templates = template.New(name)
	} else {
		this.templates = this.templates.New(name)
	}
	this.templates, err = this.templates.Parse(name)
	this.So(err, should.BeNil)
}

func (this *TemplateRendererFixture) TestMissingHomePageTemplate_ValidateErr() {
	this.templates = nil
	this.parseTemplate(contracts.ArchivesTemplateName)
	this.parseTemplate(contracts.ArticleTemplateName)
	this.parseTemplate(contracts.TopicsTemplateName)
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestMissingTopicsTemplate_ValidateErr() {
	this.templates = nil
	this.parseTemplate(contracts.ArchivesTemplateName)
	this.parseTemplate(contracts.ArticleTemplateName)
	this.parseTemplate(contracts.HomePageTemplateName)
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestMissingArticleTemplate_ValidateErr() {
	this.templates = nil
	this.parseTemplate(contracts.HomePageTemplateName)
	this.parseTemplate(contracts.ArchivesTemplateName)
	this.parseTemplate(contracts.TopicsTemplateName)
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestMissingArchivesTemplate_ValidateErr() {
	this.templates = nil
	this.parseTemplate(contracts.HomePageTemplateName)
	this.parseTemplate(contracts.ArticleTemplateName)
	this.parseTemplate(contracts.TopicsTemplateName)
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestCanRenderTypesCorrespondingToTemplates() {
	home, homeErr := this.renderer.Render(contracts.RenderedHomePage{})
	this.So(homeErr, should.BeNil)
	this.So(home, should.Equal, contracts.HomePageTemplateName)

	archives, archivesErr := this.renderer.Render(contracts.RenderedArchivesPage{})
	this.So(archivesErr, should.BeNil)
	this.So(archives, should.Equal, contracts.ArchivesTemplateName)

	article, articleErr := this.renderer.Render(contracts.RenderedArticle{})
	this.So(articleErr, should.BeNil)
	this.So(article, should.Equal, contracts.ArticleTemplateName)

	topics, topicsErr := this.renderer.Render(contracts.RenderedTopicsListing{})
	this.So(topicsErr, should.BeNil)
	this.So(topics, should.Equal, contracts.TopicsTemplateName)
}

func (this *TemplateRendererFixture) TestCannotRenderUnknownTypes() {
	home, homeErr := this.renderer.Render(42)
	this.So(homeErr, should.WrapError, contracts.ErrUnsupportedRenderingType)
	this.So(home, should.BeEmpty)
}

func (this *TemplateRendererFixture) TestRenderError() {
	this.prepareRendererWithBadTemplate()

	rendered, err := this.renderer.Render(contracts.RenderedTopicsListing{})
	this.So(err, should.WrapError, contracts.ErrRenderingFailure)
	this.So(rendered, should.BeEmpty)
}

func (this *TemplateRendererFixture) prepareRendererWithBadTemplate() {
	var err error
	t := template.New(contracts.TopicsTemplateName)
	t, err = t.Parse("{{ .UnknownField }}")
	this.So(err, should.BeNil)

	this.renderer = NewTemplateRenderer(t)
}

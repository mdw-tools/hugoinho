package core

import (
	"errors"
	"testing"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestTemplateLoaderFixture(t *testing.T) {
	should.Run(&TemplateLoaderFixture{T: should.New(t)}, should.Options.UnitTests())
}

type TemplateLoaderFixture struct {
	*should.T

	disk   *InMemoryFileSystem
	loader *TemplateLoader
}

func (this *TemplateLoaderFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.loader = NewTemplateLoader(this.disk, "templates")
	_ = this.disk.MkdirAll("templates", 0755)
	_ = this.disk.WriteFile("templates/supplemental-template.tmpl", []byte(""), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.HomePageTemplateName, []byte(ValidHomePageTemplate), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.TopicsTemplateName, []byte(ValidTopicsPageTemplate), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.ArticleTemplateName, []byte(ValidArticlePageTemplate), 0644)
}

func (this *TemplateLoaderFixture) TestLoadsEachTemplate() {
	templates, err := this.loader.Load()

	this.So(err, should.BeNil)
	this.So(templates.Lookup("supplemental-template.tmpl"), should.NOT.BeNil)
	this.So(templates.Lookup(contracts.HomePageTemplateName), should.NOT.BeNil)
	this.So(templates.Lookup(contracts.TopicsTemplateName), should.NOT.BeNil)
	this.So(templates.Lookup(contracts.ArticleTemplateName), should.NOT.BeNil)
}

func (this *TemplateLoaderFixture) TestNonTemplateFilesIgnored() {
	_ = this.disk.WriteFile("templates/not-a-template", []byte(ValidHomePageTemplate), 0644)

	templates, err := this.loader.Load()

	this.So(err, should.BeNil)
	this.So(templates.Lookup("not-a-template"), should.BeNil)
	this.So(templates.Lookup(contracts.HomePageTemplateName), should.NOT.BeNil)
	this.So(templates.Lookup(contracts.TopicsTemplateName), should.NOT.BeNil)
	this.So(templates.Lookup(contracts.ArticleTemplateName), should.NOT.BeNil)
}

func (this *TemplateLoaderFixture) TestInvalidTemplateFiles_Error() {
	_ = this.disk.WriteFile("templates/invalid-template.tmpl", []byte("{{ .invalid {{{}{{{})"), 0644)

	templates, err := this.loader.Load()

	this.So(err, should.NOT.BeNil)
	this.So(templates, should.BeNil)
}

func (this *TemplateLoaderFixture) TestWalkFileErr() {
	gophers := errors.New("GOPHERS")
	this.disk.ErrWalkFunc["templates/"+contracts.HomePageTemplateName] = gophers

	templates, err := this.loader.Load()

	this.So(err, should.WrapError, gophers)
	this.So(templates, should.BeNil)
}

func (this *TemplateLoaderFixture) TestReadFileErr() {
	gophers := errors.New("GOPHERS")
	this.disk.ErrReadFile["templates/"+contracts.HomePageTemplateName] = gophers

	templates, err := this.loader.Load()

	this.So(err, should.WrapError, gophers)
	this.So(templates, should.BeNil)
}

const (
	ValidHomePageTemplate    = ``
	ValidTopicsPageTemplate  = ``
	ValidArticlePageTemplate = ``
)

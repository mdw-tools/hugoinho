package core

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"text/template"

	"github.com/mdwhatcott/huguinho/contracts"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(templates *template.Template) *TemplateRenderer {
	return &TemplateRenderer{templates: templates}
}

func (this *TemplateRenderer) Validate() error {
	rendered, err := this.Render(contracts.RenderedHomePage{})
	if err != nil {
		return err
	}
	fmt.Println("HOME:", rendered)
	if rendered == "" {
		return errors.New("missing rendered content (template must not have been provided)")
	}

	rendered, err = this.Render(contracts.RenderedTopicsListing{})
	if err != nil {
		return err
	}
	fmt.Println("TOPICS:", rendered)
	if rendered == "" {
		return errors.New("missing rendered content (template must not have been provided)")
	}

	rendered, err = this.Render(contracts.RenderedArticle{})
	if err != nil {
		return err
	}
	fmt.Println("ARTICLE:", rendered)
	if rendered == "" {
		return errors.New("missing rendered content (template must not have been provided)")
	}

	return nil
}

func (this *TemplateRenderer) Render(v interface{}) (string, error) {
	switch v.(type) {

	case contracts.RenderedArticle:
		return this.render(contracts.ArticleTemplateName, v)

	case contracts.RenderedTopicsListing:
		return this.render(contracts.TopicsTemplateName, v)

	case contracts.RenderedHomePage:
		return this.render(contracts.HomePageTemplateName, v)

	default:
		return "", fmt.Errorf(
			"%w [%v]: %v",
			contracts.ErrUnsupportedRenderingType,
			reflect.TypeOf(v).Name(), v,
		)
	}
}

func (this *TemplateRenderer) render(name string, data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	err := this.templates.ExecuteTemplate(buffer, name, data)
	if err != nil {
		return "", fmt.Errorf(
			"%w failed to render template [%s] (err: %v) with data of type [%v]: %+v",
			contracts.ErrRenderingFailure,
			name, err, reflect.TypeOf(data), data,
		)
	}
	return buffer.String(), nil
}
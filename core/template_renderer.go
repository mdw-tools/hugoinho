package core

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"reflect"

	"github.com/mdw-tools/hugoinho/contracts"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(templates *template.Template) *TemplateRenderer {
	return &TemplateRenderer{templates: templates}
}

func (this *TemplateRenderer) Validate() (result error) {
	pages := []any{
		contracts.RenderedHomePage{},
		contracts.RenderedArchivesPage{},
		contracts.RenderedTopicsListing{},
		contracts.RenderedArticle{},
	}
	for _, page := range pages {
		if _, err := this.Render(page); err != nil {
			result = errors.Join(result, err)
		}
	}
	return result
}

func (this *TemplateRenderer) Render(v any) (string, error) {
	switch instance := v.(type) {

	case contracts.RenderedArticle:
		return this.render(contracts.ArticleTemplateName, renderedArticle{
			RenderedArticle: instance,
			Content:         template.HTML(instance.Content),
		})

	case contracts.RenderedArchivesPage:
		return this.render(contracts.ArchivesTemplateName, v)

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

func (this *TemplateRenderer) render(name string, data any) (string, error) {
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

// renderedArticle is an internal representation of the article with the Content
// (emitted from Markdown converter) marked as safe HTML.
type renderedArticle struct {
	contracts.RenderedArticle
	Content template.HTML
}

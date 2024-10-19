package core

import (
	"strings"
	"text/template"

	"github.com/mdw-tools/hugoinho/contracts"
)

type TemplateLoaderFileSystem interface {
	contracts.ReadFile
	contracts.Walk
}

type TemplateLoader struct {
	disk   TemplateLoaderFileSystem
	folder string
}

func NewTemplateLoader(disk TemplateLoaderFileSystem, folder string) *TemplateLoader {
	return &TemplateLoader{disk: disk, folder: folder}
}

func (this *TemplateLoader) Load() (templates *template.Template, err error) {
	for entry := range this.disk.Walk(this.folder) {
		if entry.Error != nil {
			return nil, StackTraceError(entry.Error)
		}
		if templates == nil {
			templates = template.New(entry.Name())
		} else {
			templates = templates.New(entry.Name())
		}
		if !strings.HasSuffix(entry.Name(), ".tmpl") {
			continue
		}
		all, err := this.disk.ReadFile(entry.Path)
		if err != nil {
			return nil, StackTraceError(err)
		}
		templates, err = templates.Parse(string(all))
		if err != nil {
			return nil, StackTraceError(err)
		}
	}
	return templates, nil
}

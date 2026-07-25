package core

import (
	"path/filepath"
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
			return nil, entry.Error
		}
		if !strings.HasSuffix(entry.Name(), ".tmpl") {
			continue
		}
		templateName := this.templateName(entry)
		if templates == nil {
			templates = template.New(templateName)
		} else {
			templates = templates.New(templateName)
		}
		all, err := this.disk.ReadFile(entry.Path)
		if err != nil {
			return nil, err
		}
		templates, err = templates.Parse(string(all))
		if err != nil {
			return nil, err
		}
	}
	return templates, nil
}

// templateName computes the template name from an entry's path.
// Top-level files use their filename (e.g., "home.tmpl").
// Subdirectory files use their relative path (e.g., "subdir/home.tmpl").
func (this *TemplateLoader) templateName(entry contracts.FileSystemEntry) string {
	// Error is ignored because Walk only yields entries under the folder,
	// so filepath.Rel is guaranteed to succeed.
	rel, _ := filepath.Rel(this.folder, entry.Path)
	return rel
}

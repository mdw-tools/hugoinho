package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mdw-go/testing/v2/better"
	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
)

func TestPipelineRunnerFixture(t *testing.T) {
	suite.Run(&PipelineRunnerFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type PipelineRunnerFixture struct {
	*suite.T

	log      *bytes.Buffer
	started  time.Time
	finished time.Time
	args     []string
	disk     *InMemoryFileSystem
	runner   *PipelineRunner
}

func (this *PipelineRunnerFixture) Setup() {
	this.log = new(bytes.Buffer)
	this.disk = NewInMemoryFileSystem()

	this.file("content/a.md", ContentA)
	this.file("content/b.md", ContentB)
	this.file("content/c.md", ContentC)

	this.file("templates/listing.tmpl", TemplateHome)
	this.file("templates/topics.tmpl", TemplateTopics)
	this.file("templates/article.tmpl", TemplateArticle)
}

func (this *PipelineRunnerFixture) buildRunner() *PipelineRunner {
	this.started = Date(2022, 1, 1)
	this.finished = this.started.Add(time.Millisecond)
	this.runner = NewPipelineRunner("version", this.args, this.disk, this.Now, log.New(this.log, "", 0))
	return this.runner
}
func (this *PipelineRunnerFixture) Now() time.Time {
	defer func() {
		this.started = this.finished
	}()
	return this.started
}

func (this *PipelineRunnerFixture) arg(values ...string) {
	this.args = append(this.args, values...)
}
func (this *PipelineRunnerFixture) file(path, content string) {
	_ = this.disk.WriteFile(path, []byte(content), 0644)
}
func (this *PipelineRunnerFixture) ls(root string) {
	files := this.disk.Walk(root)
	for file := range files {
		this.Println("Path:", filepath.Join(root, file.Path))
	}
}
func (this *PipelineRunnerFixture) assertFolder(path string) {
	dir := this.disk.Files[path]
	this.So(dir, better.NOT.BeNil)
	this.So(dir.IsDir(), should.BeTrue)
}
func (this *PipelineRunnerFixture) assertFile(path, expectedContent string) {
	this.Println("Path:", path)
	file := this.disk.Files[path]
	this.So(file, better.NOT.BeNil)
	actual := strings.ReplaceAll(strings.TrimSpace(file.Content()), "\n", `\n`)
	expected := strings.ReplaceAll(strings.TrimSpace(expectedContent), "\n", `\n`)
	this.So(actual, should.Equal, expected)
}

func (this *PipelineRunnerFixture) TestBadConfigPreventsProcessing_Error() {
	this.arg("-invalid", "=l2k3j")
	errs := this.buildRunner().Run()
	this.So(errs, should.Equal, 1)
	this.assertOriginalDiskState()
}

func (this *PipelineRunnerFixture) TestTemplateLoadFailurePreventsProcessing_Error() {
	this.disk.ErrReadFile["templates/listing.tmpl"] = errors.New("gophers")
	errs := this.buildRunner().Run()
	this.So(errs, should.Equal, 1)
	this.assertOriginalDiskState()
}

func (this *PipelineRunnerFixture) TestTemplateValidationFailurePreventsProcessing_Error() {
	this.file("templates/listing.tmpl", `{{ .INVALID }}`)
	errs := this.buildRunner().Run()
	this.So(errs, should.Equal, 1)
	this.assertOriginalDiskState()
}

func (this *PipelineRunnerFixture) TestValidConfigAndTemplates_PipelineRuns() {
	this.arg("-base-path", "/base-path")
	this.assertOriginalDiskState()

	errs := this.buildRunner().Run()

	this.So(errs, should.Equal, 0)
	this.assertRenderedDiskState()
}

func (this *PipelineRunnerFixture) assertRenderedDiskState() {
	this.So(len(this.disk.Files), should.Equal, 14)
	files, _ := json.MarshalIndent(this.disk.Files, "", "  ")
	this.Println("FILES:", string(files))

	this.assertFolder("rendered")
	this.assertFolder("rendered/topics")
	this.assertFolder("rendered/article-a")
	this.assertFolder("rendered/article-b")

	this.assertFile("rendered/index.html", RenderedListDescending)
	this.assertFile("rendered/topics/index.html", RenderedTopics)
	this.assertFile("rendered/article-a/index.html", RenderedArticleA)
}

func (this *PipelineRunnerFixture) assertOriginalDiskState() {
	this.So(len(this.disk.Files), should.Equal, 6) // 3 articles, 3 templates
}

const (
	ContentA = `
slug:   /article-a/
title:  Article A
intro:  The introduction for Article A.
topics: important misc
date:   2021-02-08

+++

This is the first article.
`

	ContentB = `
slug:   /article-b/
title:  Article B
intro:  The introduction for Article B.
topics: important
date:   2021-02-09

+++

This is the second article.
`

	ContentC = `
slug:   /article-c/
title:  Article C
intro:  The introduction for Article C.
topics: important
date:   2021-02-10
draft:  true

+++

This is the third article.
`

	TemplateHome = `
{{ range .Pages }}
Slug:   {{ .Slug }}
Title:  {{ .Title }}
Date:   {{ .Date.Format "2006-01-02" }}
Intro:  {{ .Intro }}
Topics: {{ range .Topics }}{{ . }} {{ end }}
------------------------------------------------------------------
{{ end }}
`

	TemplateArticle = `
Title:  {{ .Title }}
Intro:  {{ .Intro }}
Date:   {{ .Date.Format "2006-01-02" }}
Topics: {{ range .Topics }}{{ . }} {{ end }}
Content:

{{ .Content }}
`

	TemplateTopics = `
{{ range .Topics }}{{ .Topic }}
{{ range .Articles }}
	Slug:  {{ .Slug }}
	Title: {{ .Title }}
	Intro: {{ .Intro }}
	Date:  {{ .Date.Format "2006-01-02" }}
{{ end }}
{{ end }}
`

	RenderedArticleA = `
Title:  Article A
Intro:  The introduction for Article A.
Date:   2021-02-08
Topics: important misc 
Content:

<p>This is the first article.</p>
`

	RenderedTopics = `
important

	Slug:  /article-b/
	Title: Article B
	Intro: The introduction for Article B.
	Date:  2021-02-09

	Slug:  /article-a/
	Title: Article A
	Intro: The introduction for Article A.
	Date:  2021-02-08

`

	RenderedListAscending = `
Slug:   /article-a/
Title:  Article A
Date:   2021-02-08
Intro:  The introduction for Article A.
Topics: important misc 
------------------------------------------------------------------

Slug:   /article-b/
Title:  Article B
Date:   2021-02-09
Intro:  The introduction for Article B.
Topics: important 
------------------------------------------------------------------
`

	RenderedListDescending = `
Slug:   /article-b/
Title:  Article B
Date:   2021-02-09
Intro:  The introduction for Article B.
Topics: important 
------------------------------------------------------------------

Slug:   /article-a/
Title:  Article A
Date:   2021-02-08
Intro:  The introduction for Article A.
Topics: important misc 
------------------------------------------------------------------
`
)

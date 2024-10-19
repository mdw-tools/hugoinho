package core

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestReporterFixture(t *testing.T) {
	should.Run(&ReporterFixture{T: should.New(t)}, should.Options.UnitTests())
}

type ReporterFixture struct {
	*should.T
}

func (this *ReporterFixture) Test() {
	started := time.Now()
	stopped := started.Add(time.Millisecond * 42).Add(time.Microsecond * 1)

	stream := make(chan contracts.Article)
	go this.load(stream)

	logger := new(bytes.Buffer)
	reporter := NewReporter(started, log.New(logger, "", 0))
	reporter.ProcessStream(stream)
	reporter.RenderFinalReport(stopped)

	this.So(reporter.Errors(), should.Equal, 1)
	this.So(logger.String(), should.Equal, strings.Join([]string{
		"[INFO] published article: /a",
		"[INFO] dropped article",
		"[INFO] published article: /c",
		"[WARN] error: GOPHERS",
		"[INFO] published article: /e",
		"[INFO] errors encountered:  1",
		"[INFO] dropped articles:    1",
		"[INFO] published articles:  3",
		"[INFO] processing duration: 42ms",
		"",
	}, "\n"))
}

func (this *ReporterFixture) load(stream chan contracts.Article) {
	defer close(stream)
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/a"}}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/b"}, Error: contracts.ErrDroppedArticle}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/c"}}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/d"}, Error: errors.New("GOPHERS")}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/e"}}
}

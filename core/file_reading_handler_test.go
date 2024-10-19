package core

import (
	"errors"
	"testing"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestFileReaderFixture(t *testing.T) {
	should.Run(&FileReaderFixture{T: should.New(t)}, should.Options.UnitTests())
}

type FileReaderFixture struct {
	*should.T
	reader *FileReadingHandler
	files  *InMemoryFileSystem
}

func (this *FileReaderFixture) Setup() {
	this.files = NewInMemoryFileSystem()
	this.reader = NewFileReadingHandler(this.files)

	_ = this.files.WriteFile("/file1", []byte("FILE1"), 0644)
}

func (this *FileReaderFixture) TestRead() {
	article := &contracts.Article{Source: contracts.ArticleSource{Path: "/file1"}}
	this.reader.Handle(article)
	this.So(article, should.Equal, &contracts.Article{
		Source: contracts.ArticleSource{Path: "/file1", Data: "FILE1"},
	})
}

func (this *FileReaderFixture) TestReadError() {
	article := &contracts.Article{Source: contracts.ArticleSource{Path: "/file1"}}
	this.files.ErrReadFile["/file1"] = readError

	this.reader.Handle(article)

	this.So(article.Error, should.WrapError, readError)
	this.So(article, should.Equal, &contracts.Article{
		Error:  article.Error,
		Source: contracts.ArticleSource{Path: "/file1"},
	})
}

var readError = errors.New("read error")

package core

import (
	"testing"

	"github.com/mdw-go/testing/should"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestMetadataParserFixture(t *testing.T) {
	should.Run(&MetadataParserFixture{T: should.New(t)}, should.Options.UnitTests())
}

type MetadataParserFixture struct {
	*should.T
	parser  *MetadataParsingHandler
	article *contracts.Article
}

func (this *MetadataParserFixture) Setup() {
	this.parser = NewMetadataParsingHandler()
	this.article = &contracts.Article{}
}

func (this *MetadataParserFixture) appendSourceLine(line string) {
	this.article.Source.Data += line + "\n"
}
func (this *MetadataParserFixture) appendDividerAndContent() {
	this.appendSourceLine("\n+++\n")
	this.appendSourceLine("article content")
}
func (this *MetadataParserFixture) appendMetadataWithContent(lines ...string) {
	for _, line := range lines {
		this.appendSourceLine(line)
	}
	this.appendDividerAndContent()
}

func (this *MetadataParserFixture) TestBlankFile_Err() {
	this.appendSourceLine(" \t \n  ")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errMissingMetadata)
}
func (this *MetadataParserFixture) TestMissingMetadataDivider_Err() {
	this.appendSourceLine("title: This is the title")
	this.appendSourceLine("")
	this.appendSourceLine("This is the article content")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errMissingMetadataDivider)
}
func (this *MetadataParserFixture) TestDuplicateTitle_Err() {
	this.appendMetadataWithContent(
		"title: This is the title",
		"title: This is another title",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errDuplicateMetadataTitle)
}
func (this *MetadataParserFixture) TestBlankTitle_Err() {
	this.appendMetadataWithContent("title: ")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errBlankMetadataTitle)
}
func (this *MetadataParserFixture) TestDuplicateIntro_Err() {
	this.appendMetadataWithContent(
		"intro: This is the intro",
		"intro: This is another intro",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errDuplicateMetadataIntro)
}
func (this *MetadataParserFixture) TestDuplicateSlug_Err() {
	this.appendMetadataWithContent(
		"slug: /this/is/the/slug",
		"slug: /this/is/another/slug",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errDuplicateMetadataSlug)
}
func (this *MetadataParserFixture) TestInvalidSlug_Err() {
	this.appendMetadataWithContent("slug: /this/slug/is invalid")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataSlug)
}
func (this *MetadataParserFixture) TestInvalidCasingSlug_Err() {
	this.appendMetadataWithContent("slug: /THIS/SLUG/IS/INVALID")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataSlug)
}
func (this *MetadataParserFixture) TestBlankSlug_Err() {
	this.appendMetadataWithContent("slug: ")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errBlankMetadataSlug)
}
func (this *MetadataParserFixture) TestBlankDraft_Err() {
	this.appendMetadataWithContent("draft: ")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errBlankMetadataDraft)
}
func (this *MetadataParserFixture) TestInvalidDraft_Err() {
	this.appendMetadataWithContent("draft: not-true-or-false")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataDraft)
}
func (this *MetadataParserFixture) TestMetadataDraftFalse_Valid() {
	this.appendMetadataWithContent("draft: false")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.BeNil)
	this.So(this.article.Error, should.BeNil)
}
func (this *MetadataParserFixture) TestDuplicateDraft_Err() {
	this.appendMetadataWithContent(
		"draft: false",
		"draft: true",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errDuplicateMetadataDraft)
}
func (this *MetadataParserFixture) TestBlankDate_Err() {
	this.appendMetadataWithContent("date: ")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errBlankMetadataDate)
}
func (this *MetadataParserFixture) TestInvalidDate_Err() {
	this.appendMetadataWithContent("date: not-a-date")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataDate)
}
func (this *MetadataParserFixture) TestDuplicateDate_Err() {
	this.appendMetadataWithContent(
		"date: 2020-02-01",
		"date: 2020-02-02",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errDuplicateMetadataDate)
}
func (this *MetadataParserFixture) TestInvalidTopics_Err() {
	this.appendMetadataWithContent("topics: invalid?!")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataTopics)
}
func (this *MetadataParserFixture) TestInvalidCasingTopics_Err() {
	this.appendMetadataWithContent("topics: INVALID")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataTopics)
}
func (this *MetadataParserFixture) TestInvalidRepeatedTopics_Err() {
	this.appendMetadataWithContent("topics: repeated something-else repeated")

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errInvalidMetadataTopics)
}
func (this *MetadataParserFixture) TestDuplicateTopics_Err() {
	this.appendMetadataWithContent(
		"topics: a b c",
		"topics: x y z",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.WrapError, errDuplicateMetadataTopics)
}

func (this *MetadataParserFixture) TestAllValidAttributes() {
	this.appendMetadataWithContent(
		"title:  This is the title ",
		"intro:  This is the intro ",
		"slug:   /this/is/the/slug ",
		"draft:  true              ",
		"date:   2020-02-16        ",
		"topics: a-a b c           ",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.BeNil)
	this.So(this.article.Metadata.Title, should.Equal, "This is the title")
	this.So(this.article.Metadata.Intro, should.Equal, "This is the intro")
	this.So(this.article.Metadata.Slug, should.Equal, "/this/is/the/slug")
	this.So(this.article.Metadata.Draft, should.BeTrue)
	this.So(this.article.Metadata.Date, should.Equal, Date(2020, 2, 16))
	this.So(this.article.Metadata.Topics, should.Equal, []string{"a-a", "b", "c"})
}

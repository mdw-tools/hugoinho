package core

import (
	"bytes"
	"testing"

	"github.com/mdw-go/testing/v2/should"
	"github.com/mdw-go/testing/v2/suite"
	"github.com/mdw-tools/hugoinho/contracts"
)

func TestCLIParserFixture(t *testing.T) {
	suite.Run(&CLIParserFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type CLIParserFixture struct {
	*suite.T

	output *bytes.Buffer
	args   []string
}

func (this *CLIParserFixture) Setup() {
	this.output = new(bytes.Buffer)
}

func (this *CLIParserFixture) Parse() (contracts.Config, error) {
	parser := NewCLIParser("version", this.args)
	parser.flags.SetOutput(this.output)
	return parser.Parse()
}

func (this *CLIParserFixture) TestDefaults() {
	this.args = []string{}
	config, err := this.Parse()
	this.So(err, should.BeNil)
	this.So(config, should.Equal, contracts.Config{
		TemplateDir: "templates",
		ContentRoot: "content",
		TargetRoot:  "rendered",
		BasePath:    "",
		BuildDrafts: false,
		BuildFuture: false,
	})
}

func (this *CLIParserFixture) TestCustomValues() {
	this.args = []string{
		"-templates", "other-templates",
		"-content", "other-content",
		"-target", "other-rendered",
		"-base-path", "/path",
		"-with-drafts",
		"-with-future",
	}
	config, err := this.Parse()
	this.So(err, should.BeNil)
	this.So(config, should.Equal, contracts.Config{
		TemplateDir: "other-templates",
		ContentRoot: "other-content",
		TargetRoot:  "other-rendered",
		BasePath:    "/path",
		BuildDrafts: true,
		BuildFuture: true,
	})
}

func (this *CLIParserFixture) TestMissingTemplatesFolder() {
	this.args = []string{"-templates", ""}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
}

func (this *CLIParserFixture) TestMissingContentFolder() {
	this.args = []string{"-content", ""}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
}

func (this *CLIParserFixture) TestMissingTargetFolder() {
	this.args = []string{"-target", ""}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
}

func (this *CLIParserFixture) TestBogusValue() {
	this.args = []string{"-bogus"}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
}

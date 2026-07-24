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

func (this *CLIParserFixture) TestPathTraversalInTarget() {
	this.args = []string{"-target", "../../etc"}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
	this.So(err.Error(), should.Contain, "path traversal")
}

func (this *CLIParserFixture) TestPathTraversalInContent() {
	this.args = []string{"-content", "../../secret"}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
	this.So(err.Error(), should.Contain, "path traversal")
}

func (this *CLIParserFixture) TestPathTraversalInTemplates() {
	this.args = []string{"-templates", "../../etc/templates"}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
	this.So(err.Error(), should.Contain, "path traversal")
}

func (this *CLIParserFixture) TestValidNestedPath() {
	this.args = []string{"-content", "content/posts", "-target", "rendered/html"}
	config, err := this.Parse()
	this.So(err, should.BeNil)
	this.So(config.ContentRoot, should.Equal, "content/posts")
	this.So(config.TargetRoot, should.Equal, "rendered/html")
}

func (this *CLIParserFixture) TestPathTraversalWindowsStyle() {
	this.args = []string{"-target", "..\\etc"}
	config, err := this.Parse()
	this.So(err, should.WrapError, ErrInvalidConfig)
	this.So(config, should.Equal, contracts.Config{})
	this.So(err.Error(), should.Contain, "path traversal")
}

func (this *CLIParserFixture) TestHasPathTraversalClean() {
	this.So(hasPathTraversal("content/posts"), should.BeFalse)
	this.So(hasPathTraversal("templates"), should.BeFalse)
	this.So(hasPathTraversal("rendered/html"), should.BeFalse)
}

func (this *CLIParserFixture) TestHasPathTraversalTraversal() {
	this.So(hasPathTraversal("../../etc"), should.BeTrue)
	this.So(hasPathTraversal("..\\secret"), should.BeTrue)
	this.So(hasPathTraversal("a/../b"), should.BeTrue)
}

func (this *CLIParserFixture) TestSanitizeForError() {
	this.So(sanitizeForError("../../etc"), should.Equal, "<traversal>/<traversal>/etc")
	this.So(sanitizeForError("normal/path"), should.Equal, "normal/path")
}

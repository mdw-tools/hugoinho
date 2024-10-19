package core

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/mdw-tools/hugoinho/contracts"
)

type CLIParser struct {
	args   []string
	flags  *flag.FlagSet
	buffer *bytes.Buffer
}

func NewCLIParser(version string, args []string) *CLIParser {
	flags := flag.NewFlagSet(fmt.Sprintf("hugoinho @ %s", version), flag.ContinueOnError)
	buffer := new(bytes.Buffer)
	flags.SetOutput(buffer)

	return &CLIParser{
		args:   args,
		flags:  flags,
		buffer: buffer,
	}
}

func (this *CLIParser) Parse() (config contracts.Config, err error) {
	this.stringFlag("author   ", "Blog author name.                  ", "author   ", &config.Author)
	this.stringFlag("templates", "Directory with html templates.     ", "templates", &config.TemplateDir)
	this.stringFlag("content  ", "Directory with markdown content.   ", "content  ", &config.ContentRoot)
	this.stringFlag("target   ", "Directory for rendered html.       ", "rendered ", &config.TargetRoot)
	this.stringFlag("base-path", "Initial path of rendered html.     ", "         ", &config.BasePath)
	this.boolFlag("with-drafts", "When set, include drafts.             ", false, &config.BuildDrafts)
	this.boolFlag("with-future", "When set, include future articles.    ", false, &config.BuildFuture)

	err = this.flags.Parse(this.args)
	if err != nil {
		return contracts.Config{}, this.composeError(err)
	}

	err = validateConfig(config)
	if err != nil {
		return contracts.Config{}, this.composeError(err)
	}

	return config, nil
}

func (this *CLIParser) composeError(err error) error {
	return fmt.Errorf("%w: %v\n%s", ErrInvalidConfig, err, this.buffer.String())
}

func (this *CLIParser) stringFlag(name, description, value string, s *string) {
	this.flags.StringVar(s,
		strings.TrimSpace(name),
		strings.TrimSpace(value),
		strings.TrimSpace(description),
	)
}

func (this *CLIParser) boolFlag(name, description string, value bool, b *bool) {
	this.flags.BoolVar(b,
		strings.TrimSpace(name),
		value,
		strings.TrimSpace(description),
	)
}

func validateConfig(config contracts.Config) error {
	if config.TemplateDir == "" {
		return errors.New("template directory is required")
	}
	if config.ContentRoot == "" {
		return errors.New("content directory is required")
	}
	if config.TargetRoot == "" {
		return errors.New("target directory is required")
	}
	return nil
}

var ErrInvalidConfig = errors.New("invalid config")

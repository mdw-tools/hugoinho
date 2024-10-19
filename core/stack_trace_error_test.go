package core

import (
	"errors"
	"testing"

	"github.com/mdw-go/testing/should"
)

func TestTraceErrorFixture(t *testing.T) {
	should.Run(&StackTraceErrorFixture{T: should.New(t)}, should.Options.UnitTests())
}

type StackTraceErrorFixture struct {
	*should.T
}

func (this *StackTraceErrorFixture) Test() {
	gopherErr := errors.New("gophers")
	err := StackTraceError(gopherErr)
	if this.So(err, should.WrapError, gopherErr) {
		this.So(err.Error(), should.Contain, "gophers")
		this.So(err.Error(), should.Contain, "stack:")
	}
}

func (this *StackTraceErrorFixture) TestNil() {
	var err error
	err = StackTraceError(err)
	this.So(err, should.BeNil)
}

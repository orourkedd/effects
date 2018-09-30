package effects

import (
	"context"
	"time"
)

// Cmd -
type Cmd interface {
	GetValue() interface{}
	SetValue(interface{}) error
}

// Context -
type Context interface {
	Child() Context
	Do(cmd Cmd) error
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
	Abort(...interface{}) bool
}

// RealContext is an effects-as-data context
type RealContext struct {
	Context     context.Context
	Interpreter func(Cmd, Context) error
}

// Child -
func (ctx RealContext) Child() Context {
	return RealContext{
		Interpreter: ctx.Interpreter,
		Context:     ctx,
	}
}

// Abort -
func (ctx RealContext) Abort(args ...interface{}) bool {
	return false
}

// Do processes a command
func (ctx RealContext) Do(cmd Cmd) error {
	return ctx.Interpreter(cmd, ctx)
}

// Deadline -
func (ctx RealContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.Context.Deadline()
}

// Done -
func (ctx RealContext) Done() <-chan struct{} {
	return ctx.Context.Done()
}

// Err -
func (ctx RealContext) Err() error {
	return ctx.Context.Err()
}

// Value -
func (ctx RealContext) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

// NewContext -
func NewContext(interpreter func(Cmd, Context) error) Context {
	return RealContext{
		Interpreter: interpreter,
		Context:     context.Background(),
	}
}

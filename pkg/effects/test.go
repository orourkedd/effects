package effects

import (
	"context"
	"reflect"
	"time"
)

// TestContext is an effects-as-data context
type TestContext struct {
	Context     context.Context
	Interpreter func(interface{}, Context) error
	Args        []interface{}
	Expected    [][]interface{}
	ShouldAbort bool
	CmdQueue    []func(interface{})
	CmdIndex    int
}

// Child -
func (ctx *TestContext) Child() Context {
	return &TestContext{
		Context:     ctx,
		ShouldAbort: true,
	}
}

// Abort -
func (ctx *TestContext) Abort(args ...interface{}) bool {
	if ctx.ShouldAbort {
		return true
	}
	// record stuff here
	return false
}

// Do processes a command
func (ctx *TestContext) Do(cmd interface{}) error {
	ctx.CmdQueue[ctx.CmdIndex](cmd)
	ctx.CmdIndex++
	return nil
}

// Deadline -
func (ctx *TestContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.Context.Deadline()
}

// Done -
func (ctx *TestContext) Done() <-chan struct{} {
	return ctx.Context.Done()
}

// Err -
func (ctx *TestContext) Err() error {
	return ctx.Context.Err()
}

// Value -
func (ctx *TestContext) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

// Cmd -
func (ctx *TestContext) Cmd(fn interface{}) {
	f := func(cmd interface{}) {
		value := reflect.ValueOf(fn)

		if value.Kind() != reflect.Func {
			panic("Check must receive a function.")
		}

		if value.Type().NumIn() != 1 {
			panic("Function can only take 1 argument")
		}

		value.Call([]reflect.Value{reflect.ValueOf(cmd)})
	}
	ctx.CmdQueue = append(ctx.CmdQueue, f)
}

// NewTestContext -
func NewTestContext() *TestContext {
	return &TestContext{
		Context: context.Background(),
	}
}

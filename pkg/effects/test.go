package effects

import (
	"context"
	"reflect"
	"time"
)

// TestContext is an effects-as-data context
type TestContext struct {
	Context      context.Context
	Interpreter  func(interface{}, Context) error
	Args         []interface{}
	Expected     [][]interface{}
	CmdChan      chan interface{}
	ContinueChan chan interface{}
	ShouldAbort  bool
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
	ctx.CmdChan <- cmd
	ctx.ContinueChan <- struct{}{}
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

// Next -
func (ctx *TestContext) Next() interface{} {
	return <-ctx.CmdChan
}

// Continue -
func (ctx *TestContext) Continue() interface{} {
	return <-ctx.ContinueChan
}

// Wait -
func (ctx *TestContext) Wait() {
	ctx.ContinueChan <- struct{}{}
}

// Start -
func (ctx *TestContext) Start(fn func()) {
	go func() {
		fn()
		ctx.Wait()
	}()
}

// Cmd -
func (ctx *TestContext) Cmd(fn interface{}) {
	value := reflect.ValueOf(fn)

	if value.Kind() != reflect.Func {
		panic("Check must receive a function.")
	}

	if value.Type().NumIn() != 1 {
		panic("Function can only take 1 argument")
	}

	cmd := <-ctx.CmdChan
	value.Call([]reflect.Value{reflect.ValueOf(cmd)})
	<-ctx.ContinueChan
}

// End -
func (ctx *TestContext) End(fn ...func()) {
	<-ctx.ContinueChan
	if len(fn) > 0 {
		fn[0]()
	}
}

// NewTestContext -
func NewTestContext() *TestContext {
	return &TestContext{
		Context:      context.Background(),
		CmdChan:      make(chan interface{}),
		ContinueChan: make(chan interface{}),
	}
}

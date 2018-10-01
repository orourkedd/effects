package effects

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"gotest.tools/assert"
)

// TestContext is an effects-as-data context
type TestContext struct {
	Context     context.Context
	Interpreter func(Cmd, Context) error
	Args        []interface{}
	Expected    [][]interface{}
	CallIndex   int
	Testing     *testing.T
	ShouldAbort bool
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
	if !ctx.ShouldAbort {
		return false
	}
	ctx.Args = args
	ctx.CallIndex++
	// ctx.Context.(*TestContext).CallLog = append(ctx.Context.(*TestContext).CallLog, FunctionCall(args...))
	return true
}

// Do processes a command
func (ctx *TestContext) Do(cmd Cmd) error {
	before := ctx.Expected[ctx.CallIndex][0]
	after := ctx.Expected[ctx.CallIndex][1]

	// Do assertion
	expected := reflect.ValueOf(before).Interface()
	actual := reflect.ValueOf(cmd).Elem().Interface()
	fmt.Println(expected, actual)
	assert.Equal(ctx.Testing, expected, actual)

	// Set pointer to new value
	ptr := reflect.ValueOf(cmd).Elem()
	if ptr.CanSet() {
		ptr.Set(reflect.ValueOf(after))
	}

	ctx.CallIndex++

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

// NewTestContext -
func NewTestContext(t *testing.T, expected [][]interface{}) *TestContext {
	return &TestContext{
		Context:  context.Background(),
		Expected: expected,
		Testing:  t,
	}
}

// FunctionCall -
func FunctionCall(args ...interface{}) []interface{} {
	return args
}

// Afterer -
type Afterer struct {
	CallLog []interface{}
}

// After -
func (a *Afterer) After(after interface{}) []interface{} {
	a.CallLog = append(a.CallLog, after)
	return a.CallLog
}

// Before -
func Before(before interface{}) *Afterer {
	a := &Afterer{
		CallLog: []interface{}{before},
	}
	return a
}

package effects

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// TestContext is an effects-as-data context
type TestContext struct {
	Context     context.Context
	Interpreter func(interface{}, Context) error
	Args        []interface{}
	Expected    [][]interface{}
	CallIndex   int
}

// Child -
func (ctx *TestContext) Child() Context {
	return &TestContext{
		Context: ctx,
	}
}

// Abort -
func (ctx *TestContext) Abort(args ...interface{}) bool {
	ctx.Args = args
	ctx.CallIndex++
	// ctx.Context.(*TestContext).CallLog = append(ctx.Context.(*TestContext).CallLog, FunctionCall(args...))
	return true
}

// Do processes a command
func (ctx *TestContext) Do(cmd interface{}) error {
	// cmdType := reflect.TypeOf(cmd)
	// cmdValue := reflect.ValueOf(cmd)
	// fmt.Println("value:", cmdValue)
	// fmt.Println("Type:", cmdType)
	// copiedCmd := reflect.New(cmdType)
	// copier.Copy(&copiedCmd, cmdValue)
	// fmt.Println("Do", cmd, copiedCmd)
	// ctx.CallLog = append(ctx.CallLog, copiedCmd)
	fmt.Println("Expected:", ctx.Expected[ctx.CallIndex][0])
	fmt.Println("Actual:", cmd)
	fmt.Println(reflect.DeepEqual(ctx.Expected[ctx.CallIndex][0], cmd))
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
func NewTestContext(expected [][]interface{}) *TestContext {
	return &TestContext{
		Context:  context.Background(),
		Expected: expected,
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

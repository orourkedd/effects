package effects

import (
	"context"
	"fmt"
	"github.com/sanity-io/litter"
	"reflect"
	"testing"
	"time"
)

type TestContext struct {
	Context     context.Context
	Parent      *TestContext
	Interpreter func(interface{}, Context) error
	Args        []interface{}
	Expected    [][]interface{}
	ShouldAbort bool
	CmdQueue    []func(interface{}) error
	CmdIndex    int
	FnArgs      []interface{}
	FnErr       error
	T           *testing.T
}

func (ctx *TestContext) Do(cmd interface{}) error {
	if ctx.CmdIndex >= len(ctx.CmdQueue) {
		ctx.T.Fatalf("attempting to process a command (number %d in your function) not specified in your test.  You'll need to add another ctx.Cmd(...) to your test to account for this command:\n%s", ctx.CmdIndex+1, litter.Sdump(cmd))
	}
	err := ctx.CmdQueue[ctx.CmdIndex](cmd)
	ctx.CmdIndex++
	return err
}

func (ctx *TestContext) DoSeries(cmds interface{}) error {
	if ctx.CmdIndex >= len(ctx.CmdQueue) {
		panic(fmt.Sprintf("attempting to process a command (%d) not specified in test", ctx.CmdIndex+1))
	}
	err := ctx.CmdQueue[ctx.CmdIndex](cmds)
	ctx.CmdIndex++
	return err
}

func (ctx *TestContext) DoConcurrent(cmds interface{}) error {
	if ctx.CmdIndex >= len(ctx.CmdQueue) {
		panic(fmt.Sprintf("attempting to process a command (%d) not specified in test", ctx.CmdIndex+1))
	}
	err := ctx.CmdQueue[ctx.CmdIndex](cmds)
	ctx.CmdIndex++
	return err
}

func (ctx *TestContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.Context.Deadline()
}

func (ctx *TestContext) Done() <-chan struct{} {
	return ctx.Context.Done()
}

func (ctx *TestContext) Err() error {
	return ctx.FnErr
}

func (ctx *TestContext) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

func (ctx *TestContext) Cmd(fn interface{}) {
	f := func(cmd interface{}) error {
		value := reflect.ValueOf(fn)

		// Verify that a function is passed in
		if value.Kind() != reflect.Func {
			panic(fmt.Sprintf("ctx.Cmd(...) must receive a function.  In your test, you're passing in a value of type `%v`", value.Type()))
		}

		// Verify that the function takes on 1 argument
		if value.Type().NumIn() != 1 {
			panic(fmt.Sprintf("ctx.Cmd(...) must receive a function that takes only 1 argument.  In your test, you're passing in a function that takes %d arguments", value.Type().NumIn()))
		}

		// Verify that the function's argument is a pointer or a slice of pointers
		expectedType := value.Type().In(0)
		actualType := reflect.TypeOf(cmd)

		if expectedType.Kind() == reflect.Slice {
			if expectedType.Elem().Kind() != reflect.Ptr {
				panic("ctx.Cmd(...) must receive a function that takes a single argument of kind ptr (pointer) or a slice of pointers")
			}
		} else {
			if expectedType.Kind() != reflect.Ptr {
				panic("ctx.Cmd(...) must receive a function that takes a single argument of kind ptr (pointer) or a slice of pointers")
			}
		}

		// Verify that the function's argument type is the same as the type that comes from the non-test code
		if expectedType != actualType {
			panic(fmt.Sprintf("Your test expected a command of type %v, but the actual command was of type %v", expectedType, actualType))
		}

		results := value.Call([]reflect.Value{reflect.ValueOf(cmd)})

		// If the function returns nothing, return nil
		if len(results) == 0 {
			return nil
		}

		// Verify that the values returned from the function is an error
		err, ok := results[0].Interface().(error)

		if !ok {
			panic(fmt.Sprintf("functions passed to ctx.Cmd(...) must return an error or return nothing.  In your test, the function is returning a value of type `%v`", results[0].Type()))
		}

		return err
	}
	ctx.CmdQueue = append(ctx.CmdQueue, f)
}

func (ctx *TestContext) Finished(t *testing.T) {
	if ctx.CmdIndex != len(ctx.CmdQueue) {
		t.Fatalf("expected %d cmds to be processed but processed %d", len(ctx.CmdQueue), ctx.CmdIndex)
	}
}

func NewTestContext(t *testing.T) *TestContext {
	return &TestContext{
		Context: context.Background(),
		T:       t,
	}
}

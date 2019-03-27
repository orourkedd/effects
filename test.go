package effects

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type TestContext struct {
	Context            context.Context
	Parent             *TestContext
	Interpreter        func(interface{}, Context) error
	Args               []interface{}
	Expected           [][]interface{}
	ShouldAbort        bool
	CmdQueue           []func(interface{}) error
	CmdQueueWithResult []func(...interface{}) (interface{}, error)
	CmdIndex           int
	FnArgs             []interface{}
	FnErr              error
}

func (ctx *TestContext) Do(cmd interface{}) error {
	if ctx.CmdIndex >= len(ctx.CmdQueue) {
		panic(fmt.Sprintf("attempting to process a command (%d) not specified in test", ctx.CmdIndex+1))
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

		if value.Kind() != reflect.Func {
			panic("Check must receive a function.")
		}

		if value.Type().NumIn() != 1 {
			panic("Function can only take 1 argument")
		}

		results := value.Call([]reflect.Value{reflect.ValueOf(cmd)})

		if len(results) == 0 {
			return nil
		}

		err := results[0].Interface().(error)

		return err
	}
	ctx.CmdQueue = append(ctx.CmdQueue, f)

	fWithResult := func(args ...interface{}) (interface{}, error) {
		value := reflect.ValueOf(fn)

		if value.Kind() != reflect.Func {
			panic("Check must receive a function.")
		}

		argValues := []reflect.Value{}
		for _, a := range args {
			argValues = append(argValues, reflect.ValueOf(a))
		}

		results := value.Call(argValues)

		if len(results) != 2 {
			panic("function must return 2 values")
		}

		err, ok := results[1].Interface().(error)

		if ok {
			return results[0].Interface(), err
		} else {
			return results[0].Interface(), nil
		}
	}
	ctx.CmdQueueWithResult = append(ctx.CmdQueueWithResult, fWithResult)
}

func (ctx *TestContext) Finished() error {
	if ctx.CmdIndex != len(ctx.CmdQueue) {
		return fmt.Errorf("expected %d cmds to be processed but processed %d", len(ctx.CmdQueue), ctx.CmdIndex)
	}
	return nil
}

func NewTestContext() *TestContext {
	return &TestContext{
		Context: context.Background(),
	}
}

package effects

import (
	"context"
	"reflect"
	"time"
)

// TestContext is an effects-as-data context
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

// Child -
func (ctx *TestContext) Child() Context {
	return &TestContext{
		Context:     ctx,
		ShouldAbort: true,
		Parent:      ctx,
	}
}

// Return -
func (ctx *TestContext) Return() interface{} {
	result, err := ctx.Parent.CmdQueueWithResult[ctx.Parent.CmdIndex](ctx.FnArgs...)
	ctx.Parent.CmdIndex++
	ctx.FnErr = err
	return result
}

// Abort -
func (ctx *TestContext) Abort(args ...interface{}) bool {
	if ctx.ShouldAbort {
		ctx.FnArgs = args
		return true
	}
	// record stuff here
	return false
}

// Do processes a command
func (ctx *TestContext) Do(cmd interface{}) error {
	err := ctx.CmdQueue[ctx.CmdIndex](cmd)
	ctx.CmdIndex++
	return err
}

// DoSeries processes a command
func (ctx *TestContext) DoSeries(cmds interface{}) error {
	ctx.CmdQueue[ctx.CmdIndex](cmds)
	ctx.CmdIndex++
	return nil
}

// DoConcurrent processes a command
func (ctx *TestContext) DoConcurrent(cmds interface{}) error {
	ctx.CmdQueue[ctx.CmdIndex](cmds)
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
	return ctx.FnErr
}

// Value -
func (ctx *TestContext) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

// Cmd -
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

// NewTestContext -
func NewTestContext() *TestContext {
	return &TestContext{
		Context: context.Background(),
	}
}

package effects

import (
	"context"
	"reflect"
	"sync"
	"time"
)

// Context -
type Context interface {
	Child() Context
	Do(interface{}) error
	DoSeries(interface{}) error
	DoConcurrent(interface{}) error
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
	Abort(...interface{}) bool
	Return() interface{}
}

// RealContext is an effects-as-data context
type RealContext struct {
	Context     context.Context
	Interpreter func(interface{}, Context) error
}

// Child -
func (ctx RealContext) Child() Context {
	return RealContext{
		Interpreter: ctx.Interpreter,
		Context:     ctx,
	}
}

// Return -
func (ctx RealContext) Return() interface{} {
	return nil
}

// Abort -
func (ctx RealContext) Abort(args ...interface{}) bool {
	return false
}

// Do processes a command
func (ctx RealContext) Do(cmd interface{}) error {
	return ctx.Interpreter(cmd, ctx)
}

// DoSeries processes a command
func (ctx RealContext) DoSeries(cmds interface{}) error {
	s := reflect.ValueOf(cmds)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	list := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		list[i] = s.Index(i).Interface()
	}

	for _, cmd := range list {
		err := ctx.Interpreter(cmd, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// DoConcurrent processes a command
func (ctx RealContext) DoConcurrent(cmds interface{}) error {
	s := reflect.ValueOf(cmds)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	list := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		list[i] = s.Index(i).Interface()
	}

	wg := sync.WaitGroup{}
	wg.Add(len(list))

	var err error

	for _, cmd := range list {
		go func(c interface{}) {
			defer wg.Done()

			cmdErr := ctx.Interpreter(c, ctx)
			if cmdErr != nil {
				err = cmdErr
			}
		}(cmd)
	}
	wg.Wait()

	return err
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
func NewContext(interpreter func(interface{}, Context) error) Context {
	return RealContext{
		Interpreter: interpreter,
		Context:     context.Background(),
	}
}

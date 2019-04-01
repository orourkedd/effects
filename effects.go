package effects

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Context interface {
	Do(interface{}) error
	DoSeries(interface{}) error
	DoConcurrent(interface{}) error
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type RealContext struct {
	Context     context.Context
	Interpreter func(Context, interface{}) error
}

func InterpretSafely(ctx RealContext, cmd interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.New(r.(string))
		}
	}()
	return ctx.Interpreter(ctx, cmd)
}

func (ctx RealContext) Do(cmd interface{}) error {
	value := reflect.ValueOf(cmd)
	if value.Kind() != reflect.Ptr {
		return errors.New("ctx.Do(...) must receive a ptr")
	}

	if reflect.ValueOf(cmd).IsNil() {
		return errors.New("ctx.Do(...) cannot receive a nil ptr")
	}
	return InterpretSafely(ctx, cmd)
}

func (ctx RealContext) DoSeries(cmds interface{}) error {
	s := reflect.ValueOf(cmds)

	if s.Kind() != reflect.Slice {
		return fmt.Errorf("a slice of cmd pointers must be passed to `DoSeries` but a `%v` was passed instead", s.Kind())
	}

	list := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Kind() != reflect.Ptr {
			return fmt.Errorf("a slice of ptrs must be passed to `DoSeries` but the slice contains a `%v` at index %d", s.Index(i).Kind(), i)
		}
		list[i] = s.Index(i).Interface()
	}

	for _, cmd := range list {
		err := ctx.Do(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx RealContext) DoConcurrent(cmds interface{}) error {
	s := reflect.ValueOf(cmds)

	if s.Kind() != reflect.Slice {
		return fmt.Errorf("a slice of cmd pointers must be passed to `DoConcurrent` but a `%v` was passed instead", s.Kind())
	}

	list := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Kind() != reflect.Ptr {
			return fmt.Errorf("a slice of ptrs must be passed to `DoConcurrent` but the slice contains a `%v` at index %d", s.Index(i).Kind(), i)
		}
		list[i] = s.Index(i).Interface()
	}

	wg := sync.WaitGroup{}
	wg.Add(len(list))

	var err error

	for _, cmd := range list {
		go func(c interface{}) {
			defer wg.Done()

			cmdErr := ctx.Do(c)
			if cmdErr != nil {
				err = cmdErr
			}
		}(cmd)
	}
	wg.Wait()

	return err
}

func (ctx RealContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.Context.Deadline()
}

func (ctx RealContext) Done() <-chan struct{} {
	return ctx.Context.Done()
}

func (ctx RealContext) Err() error {
	return ctx.Context.Err()
}

func (ctx RealContext) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

func NewContext(ctx context.Context, interpreter func(Context, interface{}) error) Context {
	return RealContext{
		Interpreter: interpreter,
		Context:     ctx,
	}
}

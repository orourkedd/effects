package effects_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/orourkedd/effects"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Now struct {
	Time time.Time
}

type Panic struct{}

type ErrorOut struct{}

type NeverReturn struct {
	ContextDone bool
}

var now = time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)

func interpreter(command interface{}, ctx effects.Context) error {
	switch cmd := command.(type) {
	case *Now:
		cmd.Time = now

	case *Panic:
		panic("oops")

	case *ErrorOut:
		return errors.New("oops")

	case *NeverReturn:
		select {
		case <-ctx.Done():
			cmd.ContextDone = true
			return ctx.Err()
		case <-time.After(1000 * time.Hour):
			return nil
		}

	default:
		panic(fmt.Sprintf("Unknown command type: %T", cmd))
	}

	return nil
}

func TestEffectsBasic(t *testing.T) {
	fn := func(ctx effects.Context) (time.Time, error) {
		n := Now{}
		err := ctx.Do(&n)
		if err != nil {
			return time.Time{}, err
		}
		return n.Time, nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	result, err := fn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, now, result)
}

func TestEffectsNonPtrCmd(t *testing.T) {
	fn := func(ctx effects.Context) (time.Time, error) {
		n := Now{}
		err := ctx.Do(n)
		if err != nil {
			return time.Time{}, err
		}
		return n.Time, nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	_, err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "ctx.Do(...) must receive a ptr", err.Error())
}

func TestEffectsHandleInterpreterPanic(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := Panic{}
		err := ctx.Do(&n)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "oops", err.Error())
}

func TestEffectsHandleInterpreterError(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := ErrorOut{}
		err := ctx.Do(&n)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "oops", err.Error())
}

func TestEffectsSeries(t *testing.T) {
	fn := func(ctx effects.Context) ([]*Now, error) {
		n := []*Now{{}, {}}
		err := ctx.DoSeries(n)
		if err != nil {
			return nil, err
		}
		return n, nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	result, err := fn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, []*Now{{Time: now}, {Time: now}}, result)
}

func TestEffectsPassPointerToSliceToDoSeries(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := []*Now{{}, {}}
		err := ctx.DoSeries(&n)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "a slice of cmd pointers must be passed to `DoSeries` but a `ptr` was passed instead", err.Error())
}

func TestEffectsPassSliceOfNonPtrsToDoSeries(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := []Now{{}, {}}
		err := ctx.DoSeries(n)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "a slice of ptrs must be passed to `DoSeries` but the slice contains a `struct` at index 0", err.Error())
}

func TestEffectsPassNonSliceToDoSeries(t *testing.T) {
	fn := func(ctx effects.Context) error {
		err := ctx.DoSeries(true)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "a slice of cmd pointers must be passed to `DoSeries` but a `bool` was passed instead", err.Error())
}

func TestEffectsConcurrent(t *testing.T) {
	fn := func(ctx effects.Context) ([]*Now, error) {
		n := []*Now{{}, {}}
		err := ctx.DoConcurrent(n)
		if err != nil {
			return nil, err
		}
		return n, nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	result, err := fn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, []*Now{{Time: now}, {Time: now}}, result)
}

func TestEffectsPassPointerToSliceToDoConcurrent(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := []*Now{{}, {}}
		err := ctx.DoConcurrent(&n)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "a slice of cmd pointers must be passed to `DoConcurrent` but a `ptr` was passed instead", err.Error())
}

func TestEffectsPassSliceOfNonPtrsToDoConcurrent(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := []Now{{}, {}}
		err := ctx.DoConcurrent(n)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "a slice of ptrs must be passed to `DoConcurrent` but the slice contains a `struct` at index 0", err.Error())
}

func TestEffectsPassNonSliceToDoConcurrent(t *testing.T) {
	fn := func(ctx effects.Context) error {
		err := ctx.DoConcurrent(true)
		if err != nil {
			return err
		}
		return nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "a slice of cmd pointers must be passed to `DoConcurrent` but a `bool` was passed instead", err.Error())
}

func TestEffectsPassNilPtr(t *testing.T) {
	fn := func(ctx effects.Context) (time.Time, error) {
		var n *Now
		err := ctx.Do(n)
		if err != nil {
			return time.Time{}, err
		}
		return n.Time, nil
	}

	ctx := effects.NewContext(context.Background(), interpreter)

	_, err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "ctx.Do(...) cannot receive a nil ptr", err.Error())
}

func TestEffectsUseParentContext(t *testing.T) {
	fn := func(ctx effects.Context) error {
		n := NeverReturn{}
		err := ctx.Do(&n)
		// validate that the context timeout code path was exercised
		assert.True(t, n.ContextDone)
		if err != nil {
			return err
		}
		return nil
	}

	done := make(chan struct{})

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	ctx := effects.NewContext(timeoutCtx, interpreter)

	go func() {
		err := fn(ctx)
		assert.NotNil(t, err)
		// validate that context deadline exceeded error was returned
		assert.Equal(t, "context deadline exceeded", err.Error())

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		// noop
	case <-time.After(time.Second * 1):
		assert.Fail(t, "effects context did not timeout")
	}

	<-done
}

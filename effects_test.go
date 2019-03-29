package effects_test

import (
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

var now = time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)

func interpreter(command interface{}, ctx effects.Context) error {
	switch cmd := command.(type) {
	case *Now:
		cmd.Time = now

	case *Panic:
		panic("oops")

	case *ErrorOut:
		return errors.New("oops")

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

	ctx := effects.NewContext(interpreter)

	result, err := fn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, now, result)
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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

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

	ctx := effects.NewContext(interpreter)

	_, err := fn(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "ctx.Do(...) cannot receive a nil ptr", err.Error())
}

package effects_test

import (
	"errors"
	"github.com/orourkedd/effects"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Get struct {
	URL  string
	Body string
}

func testRunnerFn(ctx effects.Context) (string, error) {
	// Get current time
	n := Now{}
	err := ctx.Do(&n)
	if err != nil {
		return "", err
	}

	// HTTP request
	g := Get{
		URL: "https://www.swapi.co/api/people/1",
	}
	err = ctx.Do(&g)
	if err != nil {
		return "", err
	}

	// Do a list of commands in a series
	times := []*Now{{}, {}, {}}
	err = ctx.DoSeries(times)
	if err != nil {
		return "", err
	}

	// Do a list of commands in parallel
	timesConcurrent := []*Now{{}, {}, {}}
	err = ctx.DoConcurrent(timesConcurrent)
	if err != nil {
		return "", err
	}

	return g.Body, nil
}

func TestEffectsTestRunner(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		assert.Equal(t, cmd, &Now{})
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	// Series
	ctx.Cmd(func(cmds []*Now) {
		assert.Equal(t, 3, len(cmds))
		for _, n := range cmds {
			n.Time = time.Now()
		}
	})

	// Concurrent
	ctx.Cmd(func(cmds []*Now) {
		assert.Equal(t, 3, len(cmds))
		for _, n := range cmds {
			n.Time = time.Now()
		}
	})

	body, err := testRunnerFn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "{...}", body)
	assert.Nil(t, ctx.Finished())
}

func TestEffectsTestRunnerErrorSingle(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) error {
		assert.Equal(t, cmd, &Now{})
		return errors.New("oops")
	})

	body, err := testRunnerFn(ctx)
	assert.Equal(t, err.Error(), "oops")
	assert.Equal(t, "", body)
	assert.Nil(t, ctx.Finished())
}

func TestEffectsTestRunnerErrorSingleTwoDeep(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		assert.Equal(t, cmd, &Now{})
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) error {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = ""
		return errors.New("oops")
	})

	body, err := testRunnerFn(ctx)
	assert.Equal(t, err.Error(), "oops")
	assert.Equal(t, "", body)
	assert.Nil(t, ctx.Finished())
}

func TestEffectsTestRunnerErrorInSeries(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		assert.Equal(t, cmd, &Now{})
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	// Series
	ctx.Cmd(func(cmds []*Now) error {
		return errors.New("oops")
	})

	body, err := testRunnerFn(ctx)
	assert.Equal(t, err.Error(), "oops")
	assert.Equal(t, "", body)
	assert.Nil(t, ctx.Finished())
}

func TestEffectsTestRunnerErrorInConcurrent(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		assert.Equal(t, cmd, &Now{})
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	// Series
	ctx.Cmd(func(cmds []*Now) {
		assert.Equal(t, 3, len(cmds))
		for _, n := range cmds {
			n.Time = time.Now()
		}
	})

	// Concurrent
	ctx.Cmd(func(cmds []*Now) error {
		return errors.New("oops")
	})

	body, err := testRunnerFn(ctx)
	assert.Equal(t, err.Error(), "oops")
	assert.Equal(t, "", body)
	assert.Nil(t, ctx.Finished())
}

func TestEffectsTestRunnerTooManyStepsInTest(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) error {
		return errors.New("oops")
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	body, err := testRunnerFn(ctx)
	assert.Equal(t, err.Error(), "oops")
	assert.Equal(t, "", body)
	finishError := ctx.Finished()
	assert.NotNil(t, ctx.Finished())
	assert.Equal(t, finishError.Error(), "expected 2 cmds to be processed but processed 1")
}

func TestEffectsTestRunnerNoStepsSingle(t *testing.T) {
	ctx := effects.NewTestContext()

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "attempting to process a command (1) not specified in test", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTooFewStepsSingle(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		assert.Equal(t, cmd, &Now{})
		cmd.Time = time.Now()
	})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "attempting to process a command (2) not specified in test", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTooFewStepsSeries(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "attempting to process a command (3) not specified in test", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTooFewStepsConcurrent(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) {
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	ctx.Cmd(func(cmds []*Now) {
		assert.Equal(t, 3, len(cmds))
		for _, n := range cmds {
			n.Time = time.Now()
		}
	})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "attempting to process a command (4) not specified in test", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestExpectsWrongType(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "Your test expected a command of type *effects_test.Get, but the actual command was of type *effects_test.Now", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestExpectsWrongTypeWithResult(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Get) error {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
		return nil
	})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "Your test expected a command of type *effects_test.Get, but the actual command was of type *effects_test.Now", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestCmdFunctionReturnsNonError(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) string { return "" })

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "functions passed to ctx.Cmd(...) must return an error or return nothing.  In your test, the function is returning a value of type `string`", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestCmdShouldOnlyTakeOneArgument(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(c1 *Now, c2 *Now) {})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "ctx.Cmd(...) must receive a function that takes only 1 argument.  In your test, you're passing in a function that takes 2 arguments", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestPassesNonFunctionToCmd(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd("NOT A FUNCTION")

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "ctx.Cmd(...) must receive a function.  In your test, you're passing in a value of type `string`", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestCmdShouldTakeAFunctionWithAPtrArgument(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(c1 Now) {})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "ctx.Cmd(...) must receive a function that takes a single argument of kind ptr (pointer) or a slice of pointers", r)
		}
	}()

	testRunnerFn(ctx)
}

func TestEffectsTestRunnerTestCmdShouldTakeAFunctionWithASliceOfPtrArgument(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(c1 []Now) {})

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "ctx.Cmd(...) must receive a function that takes a single argument of kind ptr (pointer) or a slice of pointers", r)
		}
	}()

	testRunnerFn(ctx)
}

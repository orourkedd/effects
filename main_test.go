package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"orourkedd.com/effects/pkg/effects"
)

func TestFn(t *testing.T) {
	expected := [][]interface{}{
		effects.Before(Now{}).After(Now{Time: time.Now()}),
		effects.Before(Get{
			URL: "https://www.swapi.co/api/people/1",
		}).After(Get{
			URL:  "https://www.swapi.co/api/people/1",
			Body: "{...}",
		}),
		effects.Before(effects.FunctionCall(foo, "foo")).After(time.Now()),
	}

	ctx := effects.NewTestContext(t, expected)
	err := fn(ctx)
	assert.Nil(t, err)
}

func TestFoo(t *testing.T) {
	n := time.Now()
	expected := [][]interface{}{
		effects.Before(Now{}).After(Now{Time: n}),
		effects.Before(effects.FunctionCall(bar)).After(nil),
	}

	ctx := effects.NewTestContext(t, expected)
	now, err := foo(ctx, "foo")
	assert.Equal(t, n, now)
	assert.Nil(t, err)
}

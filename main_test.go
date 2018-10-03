package main

import (
	"errors"
	"testing"
	"time"

	"github.com/orourkedd/effects/pkg/effects"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
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

	// Concurent
	ctx.Cmd(func(cmds []*Now) {
		assert.Equal(t, 3, len(cmds))
		for _, n := range cmds {
			n.Time = time.Now()
		}
	})

	body, err := fn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "{...}", body)
}

func TestErrorHandling(t *testing.T) {
	ctx := effects.NewTestContext()

	ctx.Cmd(func(cmd *Now) error {
		return errors.New("oops")
	})

	body, err := fn(ctx)
	assert.Equal(t, body, "")
	assert.Equal(t, "oops", err.Error())
}

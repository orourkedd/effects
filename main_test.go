package main

import (
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

	body, err := fn(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "{...}", body)
}

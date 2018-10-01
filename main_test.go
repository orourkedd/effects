package main

import (
	"testing"
	"time"

	"github.com/orourkedd/effects/pkg/effects"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	ctx := effects.NewTestContext()

	var body string
	var err error

	ctx.Start(func() {
		body, err = fn(ctx)
	})

	ctx.Cmd(func(cmd *Now) {
		assert.Equal(t, cmd, &Now{})
		cmd.Time = time.Now()
	})

	ctx.Cmd(func(cmd *Get) {
		assert.Equal(t, cmd, &Get{URL: "https://www.swapi.co/api/people/1"})
		cmd.Body = "{...}"
	})

	ctx.End(func() {
		assert.Nil(t, err)
		assert.Equal(t, "{...}", body)
	})
}

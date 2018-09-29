package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"orourkedd.com/effects/pkg/effects"
)

func TestMain(t *testing.T) {
	n := time.Now()
	expected := [][]interface{}{
		effects.Before(Now{}).After(Now{Time: n}),
		effects.Before(effects.FunctionCall(foo, "foo")).After(n),
	}
	ctx := effects.NewTestContext(expected)
	err := fn(ctx)
	assert.Nil(t, err)
	// for i, value := range expected {
	// actual := ctx.CallLog[i]
	// expected := value[0]
	// assert.Equal(t, expected, actual)
	// fmt.Println("actual:", actual)
	// fmt.Println("expected:", expected)
	// fmt.Println(reflect.DeepEqual(actual, expected))
	// assert.True(t, reflect.DeepEqual(expected, actual))
	// }
}

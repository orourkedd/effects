package main

import (
	"fmt"
	"log"
	"time"

	"github.com/orourkedd/effects/pkg/effects"
)

func main() {
	ctx := effects.NewContext(interpreter)
	result, err := fn(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(result)
}

func fn(ctx effects.Context) (string, error) {
	if ctx.Abort() { // pass in the args for asserting
		return "", nil
	}

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
		log.Fatal(err)
	}

	// pass a child context to the next function.  This is how the test framework will know
	// where to create seams in your code.
	// n2, err := foo(ctx.Child(), "foo")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(n2)

	return g.Body, nil
}

func foo(ctx effects.Context, value string) (time.Time, error) {
	if ctx.Abort(foo, value) { // pass in the args for asserting
		return time.Time{}, nil
	}

	// Get current time
	n := Now{}
	err := ctx.Do(&n)
	if err != nil {
		return time.Time{}, nil
	}

	bar(ctx.Child())

	return n.Time, nil
}

func bar(ctx effects.Context) {
	if ctx.Abort(bar) {
		return
	}

	// Get current time
	n := Now{}
	ctx.Do(&n)
	fmt.Println(n.Time)
}

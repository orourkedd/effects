package main

import (
	"fmt"
	"log"
	"time"

	"orourkedd.com/effects/pkg/effects"
)

func main() {
	// effects.NewContext() returns a struct that implements context.Context
	ctx := effects.NewContext(interpreter)
	err := fn(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func fn(ctx effects.Context) error {
	// Get current time
	n := Now{}
	err := ctx.Do(&n)
	if err != nil {
		return err
	}
	// fmt.Println(n.Time)

	// HTTP request
	g := Get{
		URL: "https://www.swapi.co/api/people/1",
	}
	err = ctx.Do(&g)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(g.Body)

	// pass a child context to the next function.  This is how the test framework will know
	// where to create seams in your code.
	n2, err := foo(ctx.Child(), "foo")
	if err != nil {
		return err
	}
	fmt.Println(n2)

	return nil
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

func bar(ctx effects.Context) interface{} {
	if ctx.Abort(bar) {
		fmt.Println("abort bar")
		return nil
	}

	// Get current time
	n := Now{}
	ctx.Do(&n)
	fmt.Println(n.Time)
	return nil
}

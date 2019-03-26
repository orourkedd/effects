package main

import (
	"fmt"
	"github.com/orourkedd/effects/pkg/effects"
	"log"
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

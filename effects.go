package main

import (
	"log"
	"time"

	"github.com/imroc/req"
	"orourkedd.com/effects/pkg/effects"
)

// Now -
type Now struct {
	Time time.Time
}

// Get -
type Get struct {
	URL  string
	Body string
}

// Define an interpreter function for each application
func interpreter(command interface{}, ctx effects.Context) error {
	switch cmd := command.(type) {
	case *Now:
		cmd.Time = time.Now()
		return nil

	case *Get:
		resp, err := req.Get(cmd.URL)
		if err != nil {
			return err
		}
		cmd.Body = string(resp.Bytes())

	default:
		log.Fatalf("Unknown command type: %T", cmd)
	}

	return nil
}

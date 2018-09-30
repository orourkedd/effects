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

// SetValue -
func (cmd Now) SetValue(value interface{}) error {
	cmd.Time = value.(time.Time)
	return nil
}

// GetValue -
func (cmd Now) GetValue() interface{} {
	return cmd.Time
}

// Get -
type Get struct {
	URL  string
	Body string
}

// SetValue -
func (cmd Get) SetValue(value interface{}) error {
	cmd.Body = value.(string)
	return nil
}

// GetValue -
func (cmd Get) GetValue() interface{} {
	return cmd.Body
}

// Define an interpreter function for each application
func interpreter(command effects.Cmd, ctx effects.Context) error {
	switch cmd := command.(type) {
	case *Now:
		cmd.SetValue(time.Now())
		return nil

	case *Get:
		resp, err := req.Get(cmd.URL)
		if err != nil {
			return err
		}
		cmd.SetValue(string(resp.Bytes()))

	default:
		log.Fatalf("Unknown command type: %T", cmd)
	}

	return nil
}

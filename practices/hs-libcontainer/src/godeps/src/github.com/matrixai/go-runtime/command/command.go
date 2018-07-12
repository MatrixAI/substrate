package command

import (
	"fmt"
	"github.com/matrixai/go-runtime/context"
)

// Command interface specifies the necessary operations
// needed for each API call.
type Command interface {
	Execute(*context.Context) (interface{}, error)
}

// Type is an enum representing command type
type Type int

const (
	// CREATE Creates the container environment
	CREATE Type = iota
	// START Start the processes within a container
	START
	// KILL Send a signal to a running container
	KILL
	// DELETE Deletes a container
	DELETE
	// STATE the state of the container
	STATE
)

// Factory creates a Command based on ctype.
func Factory(ctype Type) (cmd Command, err error) {
	switch ctype {
	case CREATE:
		fmt.Println()
		cmd = &CreateCommand{}
	// TODO: more commands to add
	default:
		err = fmt.Errorf("unknown command type")
	}
	return
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExecuteCmd validates that executing the root command itself does not return an error.
// e.g. an invalid flag name causes an error.
func TestExecuteCmd(t *testing.T) {
	err := executeCmd()
	assert.NoError(t, err)
}

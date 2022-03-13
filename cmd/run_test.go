package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBoolOrFalse(t *testing.T) {
	assert.Equal(t, parseBoolOrFalse("True"), true)
	assert.Equal(t, parseBoolOrFalse("TRUE"), true)
	assert.Equal(t, parseBoolOrFalse("true"), true)
	assert.Equal(t, parseBoolOrFalse(""), false)
	assert.Equal(t, parseBoolOrFalse("False"), false)
	assert.Equal(t, parseBoolOrFalse("FALSE"), false)
	assert.Equal(t, parseBoolOrFalse("false"), false)
}

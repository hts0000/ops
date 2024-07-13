package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrZero(t *testing.T) {
	var (
		intPtr   *int
		strPtr   *string
		floatPtr *float64
		boolPtr  *bool

		intVal   = 42
		strVal   = "hello"
		floatVal = 3.14
		boolVal  = true
	)

	assert.Equal(t, 0, GetOrZero(intPtr))
	assert.Equal(t, "", GetOrZero(strPtr))
	assert.Equal(t, 0.0, GetOrZero(floatPtr))
	assert.Equal(t, false, GetOrZero(boolPtr))

	assert.Equal(t, 42, GetOrZero(&intVal))
	assert.Equal(t, "hello", GetOrZero(&strVal))
	assert.Equal(t, 3.14, GetOrZero(&floatVal))
	assert.Equal(t, true, GetOrZero(&boolVal))
}

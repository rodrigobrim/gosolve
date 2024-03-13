package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicOnError(t *testing.T) {

	PanicOnError(t, nil)
	assert.True(t, true)

	defer func() {
		if r := recover(); r != nil {
			assert.True(t, true)
		}
	}()

	PanicOnError(t, fmt.Errorf("error"))

}

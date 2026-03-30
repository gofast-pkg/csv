package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWarning(t *testing.T) {
	t.Run("should return a new warning", func(t *testing.T) {
		w := NewWarning()

		assert.NotNil(t, w)
		assert.Empty(t, w)
	})
}

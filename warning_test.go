package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWarning(t *testing.T) {
	t.Run("should return a new warning", func(t *testing.T) {
		w := NewWarning()
		assert.NotNil(t, w)
	})
}

func TestWarning_Wrap(t *testing.T) {
	t.Run("should add values to the warning", func(t *testing.T) {
		w := NewWarning()
		warn := NewWarning()
		warn["key"] = []string{"value"}

		w.Wrap(warn)
		assert.EqualValues(t, w["key"], warn["key"])
	})
	t.Run("should append values to the warning", func(t *testing.T) {
		expected := []string{"value", "value2"}
		w := NewWarning()
		w["key"] = []string{"value"}
		warn := NewWarning()
		warn["key"] = []string{"value2"}

		w.Wrap(warn)
		assert.EqualValues(t, w["key"], expected)
	})
	t.Run("should not add values to the warning", func(t *testing.T) {
		expected := []string{"value"}
		w := NewWarning()
		w["key"] = []string{"value"}
		warn := NewWarning()

		w.Wrap(warn)
		assert.EqualValues(t, w["key"], expected)
	})
}

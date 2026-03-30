package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultContextKey   = "key"
	defaultContextValue = "value"
)

func TestNewDecoder(t *testing.T) {
	t.Run("should return an error because newInstanceFunc is nil", func(t *testing.T) {
		decoder, err := NewDecoder(ConfigDecoder{
			SaveInstanceFunc:    func(_ Decoder, _ any) error { return nil },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return nil },
		})
		require.ErrorIs(t, err, ErrInvalidConfigDecoder)
		assert.Nil(t, decoder)
	})
	t.Run("should return an error because saveInstanceFunc is nil", func(t *testing.T) {
		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:     func(_ Decoder) (any, error) { return &struct{}{}, nil },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return nil },
		})
		require.ErrorIs(t, err, ErrInvalidConfigDecoder)
		assert.Nil(t, decoder)
	})
	t.Run("should return an error because warningInstanceFunc is nil", func(t *testing.T) {
		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:  func(_ Decoder) (any, error) { return &struct{}{}, nil },
			SaveInstanceFunc: func(_ Decoder, _ any) error { return nil },
		})
		require.ErrorIs(t, err, ErrInvalidConfigDecoder)
		assert.Nil(t, decoder)
	})
	t.Run("should return a Decoder", func(t *testing.T) {
		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:     func(_ Decoder) (any, error) { return &struct{}{}, nil },
			SaveInstanceFunc:    func(_ Decoder, _ any) error { return nil },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return nil },
		})
		require.NoError(t, err)
		assert.NotNil(t, decoder)
	})
}

func TestDecoder_ContextSet(t *testing.T) {
	t.Run("should set the context", func(t *testing.T) {
		d := &decoder{context: make(map[string]any)}

		d.ContextSet(defaultContextKey, defaultContextValue)
		assert.Equal(t, defaultContextValue, d.context[defaultContextKey])
	})
	t.Run("should override the context", func(t *testing.T) {
		d := &decoder{context: make(map[string]any)}

		d.ContextSet(defaultContextKey, defaultContextValue)
		assert.Equal(t, defaultContextValue, d.context[defaultContextKey])

		value2 := "value2"
		d.ContextSet(defaultContextKey, value2)
		assert.Equal(t, value2, d.context[defaultContextKey])
	})
}

func TestDecoder_ContextGet(t *testing.T) {
	t.Run("when trying to get an unexistant value", func(t *testing.T) {
		d := &decoder{context: make(map[string]any)}

		v, ok := d.ContextGet(defaultContextKey)
		assert.False(t, ok)
		assert.Nil(t, v)
	})
	t.Run("when get an existent value", func(t *testing.T) {
		d := &decoder{context: make(map[string]any)}

		d.ContextSet(defaultContextKey, defaultContextValue)
		v, ok := d.ContextGet(defaultContextKey)

		assert.True(t, ok)
		assert.Equal(t, defaultContextValue, v)
	})
}

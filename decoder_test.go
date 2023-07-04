package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	cfg := ConfigDecoder{
		NewInstanceFunc:  func() any { return nil },
		SaveInstanceFunc: func(any, Decoder) error { return nil },
	}

	t.Run("should return a new decoder", func(t *testing.T) {
		d, err := NewDecoder(cfg)
		if assert.NoError(t, err) {
			assert.NotNil(t, d)
		}
	})
	t.Run("should return an error because the reader is nil", func(t *testing.T) {
		d, err := NewDecoder(ConfigDecoder{})
		if assert.Error(t, err) {
			assert.Nil(t, d)
			assert.EqualError(t, err, ErrConfDecoder.Error())
		}
	})
}

func TestDecoder_ContextSet(t *testing.T) {
	t.Run("should set the context", func(t *testing.T) {
		key := "key"
		value := "value"
		d := &decoder{context: make(map[string]string)}

		d.ContextSet(key, value)
		assert.Equal(t, d.context[key], value)
	})
	t.Run("should override the context", func(t *testing.T) {
		key := "key"
		value := "value"
		d := &decoder{context: make(map[string]string)}

		d.ContextSet(key, value)
		assert.Equal(t, d.context[key], value)

		value = "value2"
		d.ContextSet(key, value)
		assert.Equal(t, d.context[key], value)
	})
}

func TestDecoder_ContextGet(t *testing.T) {
	t.Run("should not find the value associated to key", func(t *testing.T) {
		d := &decoder{context: make(map[string]string)}

		v, ok := d.ContextGet("key")
		if assert.False(t, ok) {
			assert.Equal(t, v, "")
		}
	})
	t.Run("should find the value associated to key", func(t *testing.T) {
		key := "key"
		value := "value"
		d := &decoder{context: make(map[string]string)}

		d.ContextSet(key, value)
		v, ok := d.ContextGet(key)
		if assert.True(t, ok) {
			assert.Equal(t, v, value)
		}
	})
}

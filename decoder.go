package csv

import "errors"

// Errors supported by the decoder.
var (
	ErrInvalidConfigDecoder = errors.New("ConfigDecoder is invalid")
)

// ConfigDecoder is the configuration to create a new decoder.
type ConfigDecoder struct {
	NewInstanceFunc     func(Decoder) (any, error)
	SaveInstanceFunc    func(Decoder, any) error
	WarningInstanceFunc func(Decoder, Warning) error
}

// Decoder is the interface that wraps the basic methods to decode a csv file
// with a specific process for each line of the file.
//
//mockery:generate: true
type Decoder interface {
	ContextSet(key string, value any)
	ContextGet(key string) (value any, found bool)
	newInstance() (any, error)
	saveInstance(any) error
	warningUnusedField(warn Warning) error
}

type decoder struct {
	context             map[string]any
	newInstanceFunc     func(Decoder) (any, error)
	saveInstanceFunc    func(Decoder, any) error
	warningInstanceFunc func(Decoder, Warning) error
}

func (c ConfigDecoder) isValid() bool {
	if c.NewInstanceFunc == nil {
		return false
	}
	if c.SaveInstanceFunc == nil {
		return false
	}
	if c.WarningInstanceFunc == nil {
		return false
	}

	return true
}

// NewDecoder returns a new decoder with the specific configuration.
// If the configuration is not valid, the function returns an error of type
// ErrInvalidConfigDecoder.
func NewDecoder(conf ConfigDecoder) (Decoder, error) {
	if !conf.isValid() {
		return nil, ErrInvalidConfigDecoder
	}

	return &decoder{
		newInstanceFunc:     conf.NewInstanceFunc,
		saveInstanceFunc:    conf.SaveInstanceFunc,
		warningInstanceFunc: conf.WarningInstanceFunc,
		context:             make(map[string]any),
	}, nil
}

// ContextGet returns the value associated to the key passed as parameter.
func (d *decoder) ContextGet(key string) (any, bool) {
	if v, ok := d.context[key]; ok {
		return v, true
	}

	return nil, false
}

// ContextSet sets a value in the context of the decoder.
// If the key already exists, the value is overridden.
// The context is used to share data between the different functions that
// compose the process to decode a csv file.
func (d *decoder) ContextSet(key string, value any) {
	d.context[key] = value
}

func (d *decoder) newInstance() (any, error) {
	return d.newInstanceFunc(d)
}

func (d *decoder) saveInstance(obj any) error {
	return d.saveInstanceFunc(d, obj)
}

func (d *decoder) warningUnusedField(warn Warning) error {
	return d.warningInstanceFunc(d, warn)
}

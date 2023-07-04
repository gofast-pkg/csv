package csv

// Decoder is the interface that wraps the basic methods to decode a csv file
// with a specific process for each line of the file.
type Decoder interface {
	ContextSet(key, value string)
	ContextGet(key string) (value string, found bool)
	newInstance() any
	saveInstance(any) error
}

type decoder struct {
	newInstanceFunc  func() any
	saveInstanceFunc func(any, Decoder) error
	context          map[string]string
}

// ConfigDecoder is the struct that contains the configuration to create a new
// decoder.
type ConfigDecoder struct {
	NewInstanceFunc  func() any
	SaveInstanceFunc func(any, Decoder) error
}

func (c ConfigDecoder) isValid() bool {
	if c.NewInstanceFunc != nil && c.SaveInstanceFunc != nil {
		return true
	}

	return false
}

// NewDecoder returns a new decoder with the configuration passed as parameter.
// If the configuration is not valid, the function returns an error of type
// ErrConfDecoder.
func NewDecoder(conf ConfigDecoder) (Decoder, error) {
	if !conf.isValid() {
		return nil, ErrConfDecoder
	}

	return &decoder{
		newInstanceFunc:  conf.NewInstanceFunc,
		saveInstanceFunc: conf.SaveInstanceFunc,
		context:          make(map[string]string),
	}, nil
}

// ContextSet sets a value in the context of the decoder.
// If the key already exists, the value is overridden.
// The context is used to share data between the different functions that
// compose the process to decode a csv file.
func (d *decoder) ContextSet(key, value string) {
	d.context[key] = value
}

// ContextGet returns the value associated to the key passed as parameter.
func (d *decoder) ContextGet(key string) (string, bool) {
	if v, ok := d.context[key]; ok {
		return v, true
	}

	return "", false
}

func (d *decoder) newInstance() any {
	return d.newInstanceFunc()
}

func (d *decoder) saveInstance(obj any) error {
	return d.saveInstanceFunc(obj, d)
}

package csv

import "github.com/jszwec/csvutil"

// Warning collect relevants informations about the process on the csv file.
type Warning map[string]string

// NewWarning create a warning instance.
// For read a warning it's possible to iterate over it like a map.
func NewWarning() Warning {
	return make(map[string]string)
}

func (w *Warning) unusedFields(decoder *csvutil.Decoder) {
	header := decoder.Header()
	for _, i := range decoder.Unused() {
		if header[i] == "" {
			continue
		}

		(*w)[header[i]] = decoder.Record()[i]
	}
}

package csv

import "github.com/jszwec/csvutil"

// Warning collect relevants informations about the process on the csv file
type Warning map[string][]string

// NewWarning return a new warning
// For read a warning it's possible to iterate over it like a map
func NewWarning() Warning {
	return make(map[string][]string)
}

// Wrap relevants informations from the warn to the warning
func (w *Warning) Wrap(warn Warning) {
	for key, values := range warn {
		for _, value := range values {
			w.addValues(key, value)
		}
	}
}

func (w *Warning) unusedFields(decoder *csvutil.Decoder) {
	header := decoder.Header()
	for _, i := range decoder.Unused() {
		if header[i] == "" {
			continue
		}
		w.addValues(header[i], decoder.Record()[i])
	}
}

func (w *Warning) addValues(key, value string) {
	if _, ok := (*w)[key]; !ok {
		(*w)[key] = []string{value}

		return
	}
	(*w)[key] = append((*w)[key], value)
}

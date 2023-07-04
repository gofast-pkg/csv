package csv

import (
	"encoding/csv"
	"io"

	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
)

// ExtensionFile for the csv file type
const (
	ExtensionFile = ".csv"
)

// CSV interface for new csv reader
type CSV interface {
	DecodeWithDecoder(d Decoder) (Warning, error)
	Decode(obj any) (Warning, error)
}

type csvReader struct {
	decoder *csvutil.Decoder
}

// New create a new CSV reader from io.Reader.
// Separator is the separator used in the CSV file
func New(r io.Reader, separator rune) (CSV, error) {
	var err error
	reader := csvReader{}

	if r == nil {
		return nil, ErrNilReader
	}
	csvReader := csv.NewReader(r)
	csvReader.Comma = separator
	if reader.decoder, err = csvutil.NewDecoder(csvReader); err != nil {
		return nil, errors.Wrap(ErrNewDecoder, err.Error())
	}
	reader.decoder.DisallowMissingColumns = true

	return &reader, nil
}

// Decode decode the CSV file into the obj
// Each call to Decode read and process one line of the CSV file
// Return a Warning object that contains all unused fields.
func (r *csvReader) Decode(obj any) (Warning, error) {
	var err error

	warn := NewWarning()
	if err = r.decoder.Decode(obj); err != nil {
		if err == io.EOF {
			return nil, err
		}

		return nil, errors.Wrap(ErrOBJDecode, err.Error())
	}
	// aggregate unused fields in a Warning
	warn.unusedFields(r.decoder)

	return warn, nil
}

// DecodeWithDecoder decode the CSV file at the end with the preconfigured decoder
// Decoder must be preconfigured with the NewDecoder function
// Return a Warning object that contains all unused fields.
func (r *csvReader) DecodeWithDecoder(d Decoder) (Warning, error) {
	if d == nil {
		return nil, ErrDecoder
	}
	warn := NewWarning()
	for {
		obj := d.newInstance()
		entryWarn, err := r.Decode(obj)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		warn.Wrap(entryWarn)

		if err = d.saveInstance(obj); err != nil {
			return nil, errors.Wrap(ErrRecorder, err.Error())
		}
	}

	return warn, nil
}

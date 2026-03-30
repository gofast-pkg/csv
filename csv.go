// Package csv allow to parse csv file line by line and transform each row in a data model.
// Expose a DecoderConfig to define custom behavior for each row.
package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/jszwec/csvutil"
)

// Errors supported by the csv operations.
var (
	ErrNilOBJ              = errors.New("obj is nil")
	ErrNilReader           = errors.New("io.reader is nil")
	ErrNilDecoder          = errors.New("decoder is nil")
	ErrOBJDecode           = errors.New("fails to decode into specific object")
	ErrToCreateNewInstance = errors.New("fails to create new obj instance")
	ErrToSaveNewInstance   = errors.New("fails to save new obj instance")
	ErrToCheckWarning      = errors.New("fails to check the unused fields")
)

// ExtensionFile for the csv file type.
const (
	ExtensionFile = ".csv"
)

// CSV interface creates a new interface ready to parse a csv file.
//
//mockery:generate: true
type CSV interface {
	// DecodeWithDecoder parse the CSV file with a custom decoder.
	// Decoder allows you to apply specific behavior to each line of the csv.
	// (see: Decoder type and NewDecoder func)
	DecodeWithDecoder(d Decoder) error
	// Decode the csv file loaded into the obj instance.
	// Each call to Decode read and process one line of the CSV file.
	// Decode return a Warning object iterable like a map to check fields
	// which are not used.
	Decode(obj any) (Warning, error)
}

type parser struct {
	decoder *csvutil.Decoder
}

// New create a new CSV reader from an io.Reader.
// Separator is the separator used in the CSV file.
func New(r io.Reader, separator rune) (CSV, error) {
	var err error

	if r == nil {
		return nil, ErrNilReader
	}

	csvReader := csv.NewReader(r)
	csvReader.Comma = separator

	p := parser{}
	if p.decoder, err = csvutil.NewDecoder(csvReader); err != nil {
		return nil, fmt.Errorf("fails to create new csv decoder: %w", err)
	}
	p.decoder.DisallowMissingColumns = true

	return &p, nil
}

func (p *parser) Decode(obj any) (Warning, error) {
	var err error

	if obj == nil {
		return nil, ErrNilOBJ
	}

	if err = p.decoder.Decode(obj); err != nil {
		if err == io.EOF {
			return nil, err
		}

		return nil, errors.Join(err, ErrOBJDecode)
	}

	warn := NewWarning()
	warn.unusedFields(p.decoder)

	return warn, nil
}

func (p *parser) DecodeWithDecoder(d Decoder) error {
	if d == nil {
		return ErrNilDecoder
	}

	for {
		err := p.decodeWithDecoder(d)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}
	}
}

func (p *parser) decodeWithDecoder(d Decoder) error {
	obj, err := d.newInstance()
	if err != nil {
		return errors.Join(err, ErrToCreateNewInstance)
	}

	warning, err := p.Decode(obj)
	if err != nil {
		return errors.Join(err, ErrOBJDecode)
	}

	if err = d.warningUnusedField(warning); err != nil {
		return errors.Join(err, ErrToCheckWarning)
	}

	if err = d.saveInstance(obj); err != nil {
		return errors.Join(err, ErrToSaveNewInstance)
	}

	return nil
}

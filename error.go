package csv

import "github.com/pkg/errors"

// error for the csv parser
var (
	ErrConfDecoder            = errors.New("error to create a csv.NewDecoder")
	ErrNewCSVReader           = errors.New("error to create a csv.NewCSV, reader is invalid")
	ErrOBJDecode              = errors.New("error to decode into specific object")
	ErrNewDecoder             = errors.New("error to create a csvutil.NewDecoder")
	ErrDecoder                = errors.New("decoder is nil")
	ErrBadConfiguration       = errors.New("configuration are not full setup to decode csv")
	ErrNoCreateObjectToDecode = errors.New("no create object function to decode")
	ErrRecorder               = errors.New("function save instance was return an error")
	ErrNilReader              = errors.New("reader is nil")
)

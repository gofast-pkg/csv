package csv

import (
	"io"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

/*
** Test data stuff
 */

// testfiles
const (
	testFilePath      = "testdata/testfile.csv"
	testEmptyFilePath = "testdata/testfile_empty.csv"
	testWrongFilePath = "testdata/testfile_wrong_csv_format.csv"
)

type DataTestLong struct {
	Name      string `csv:"name"`
	Type      string `csv:"type"`
	MainColor string `csv:"main_color"`
	Size      string `csv:"size"`
}

type DataTestShort struct {
	Name string `csv:"name"`
	Type string `csv:"type"`
}

/*
** Unit tests
 */

func TestNew(t *testing.T) {
	t.Run("should return a new csv", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if assert.NoError(t, err) {
			assert.NotNil(t, csvReader)
		}
	})
	t.Run("should return an error because the reader is nil", func(t *testing.T) {
		csvReader, err := New(nil, ';')
		if assert.Error(t, err) {
			assert.Nil(t, csvReader)
			assert.EqualError(t, err, ErrNilReader.Error())
		}
	})
	t.Run("should return an error because the reader is invalid", func(t *testing.T) {
		reader, err := os.Open(testEmptyFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if assert.Error(t, err) {
			assert.Nil(t, csvReader)
			assert.ErrorIs(t, err, ErrNewDecoder)
		}
	})
}

func TestCSV_Decode(t *testing.T) {
	t.Run("should return an error because the config is nil", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		warn, err := csvReader.Decode(nil)
		if assert.Error(t, err) {
			assert.Nil(t, warn)
			assert.ErrorIs(t, err, ErrOBJDecode)
		}
	})
	t.Run("should read file until EOF", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		w := NewWarning()
		for {
			var warn Warning
			warn, err = csvReader.Decode(&DataTestLong{})
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Error(err)
			}
			w.Wrap(warn)
		}

		if assert.Error(t, err) {
			assert.ErrorIs(t, err, io.EOF)
			if assert.NotNil(t, w) {
				assert.Empty(t, w)
			}
		}
	})
}

func TestCSV_DecodeWithDecoder(t *testing.T) {
	t.Run("should return an error because the config is nil", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		warn, err := csvReader.DecodeWithDecoder(nil)
		if assert.Error(t, err) {
			assert.Nil(t, warn)
			assert.ErrorIs(t, err, ErrDecoder)
		}
	})
	t.Run("should decode csv file until the end", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		var list []DataTestLong
		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc: func() any { return &DataTestLong{} },
			SaveInstanceFunc: func(obj any, d Decoder) error {
				list = append(list, *(obj.(*DataTestLong)))

				return nil
			},
		})
		if err != nil {
			t.Error(err)
		}

		warn, err := csvReader.DecodeWithDecoder(decoder)

		if assert.NoError(t, err) {
			assert.Empty(t, warn)
		}
	})
	t.Run("should return an error because the decoder is invalid", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc: func() any { return &DataTestLong{} },
			SaveInstanceFunc: func(obj any, d Decoder) error {
				return errors.New("bad decoder")
			},
		})
		if err != nil {
			t.Error(err)
		}

		warn, err := csvReader.DecodeWithDecoder(decoder)

		if assert.Error(t, err) {
			assert.Nil(t, warn)
			assert.ErrorIs(t, err, ErrRecorder)
		}
	})
	t.Run("should return an error because the format doesn't match", func(t *testing.T) {
		reader, err := os.Open(testWrongFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc: func() any { return &DataTestShort{} },
			SaveInstanceFunc: func(obj any, d Decoder) error {
				return nil
			},
		})
		if err != nil {
			t.Error(err)
		}

		warn, err := csvReader.DecodeWithDecoder(decoder)

		if assert.Error(t, err) {
			assert.Nil(t, warn)
			assert.ErrorIs(t, err, ErrOBJDecode)
		}
	})
	t.Run("should decode csv file until the end with warning", func(t *testing.T) {
		reader, err := os.Open(testFilePath)
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()

		csvReader, err := New(reader, ';')
		if err != nil {
			t.Error(err)
		}

		var list []DataTestShort
		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc: func() any { return &DataTestShort{} },
			SaveInstanceFunc: func(obj any, d Decoder) error {
				list = append(list, *(obj.(*DataTestShort)))

				return nil
			},
		})
		if err != nil {
			t.Error(err)
		}

		expectedWarn := map[string][]string{
			"main_color": {"black", "blue"},
			"size":       {"big", "small"},
		}
		warn, err := csvReader.DecodeWithDecoder(decoder)

		if assert.NoError(t, err) {
			assert.NotEmpty(t, warn)
			assert.EqualValues(t, warn, expectedWarn)
		}
	})
}

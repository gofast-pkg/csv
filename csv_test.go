package csv

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// files used for the unit tests.
const (
	withValidCSVFile  = "testdata/valid.csv"
	withEmptyCSVFile  = "testdata/empty.csv"
	withNoDataCSVFile = "testdata/nodata.csv"
)

// DataTestLong contract to the valid csv test file.
type DataTestLong struct {
	Name      string `csv:"name"`
	Type      string `csv:"type"`
	MainColor string `csv:"main_color"`
	Size      string `csv:"size"`
}

// DataTestShort partial contract to the valid csv test file.
type DataTestShort struct {
	Name string `csv:"name"`
	Type string `csv:"type"`
}

// DataTestInvalid is an invalid contract to the valid csv test file.
type DataTestInvalid struct {
	Name      int      `csv:"name"`
	Type      string   `csv:"type"`
	MainColor []string `csv:"main_color"`
	Size      string   `csv:"size"`
}

func TestNew(t *testing.T) {
	t.Run("should return an error because the reader is nil", func(t *testing.T) {
		csvReader, err := New(nil, ';')
		require.ErrorIs(t, err, ErrNilReader)
		assert.Nil(t, csvReader)
	})
	t.Run("should return an error because the reader is empty", func(t *testing.T) {
		reader, err := os.Open(withEmptyCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.ErrorContains(t, err, "fails to create new csv decoder")
		assert.Nil(t, csvReader)
	})
	t.Run("should return a new csv", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		if assert.NoError(t, err) {
			assert.NotNil(t, csvReader)
		}
	})
}

func TestCSV_Decode(t *testing.T) {
	t.Run("should return an error because the obj instance is nil", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		warn, err := csvReader.Decode(nil)
		require.ErrorIs(t, err, ErrNilOBJ)
		assert.Empty(t, warn)
	})
	t.Run("should return an io.EOF error", func(t *testing.T) {
		reader, err := os.Open(withNoDataCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		warn, err := csvReader.Decode(&DataTestLong{})
		require.ErrorIs(t, err, io.EOF)
		assert.Nil(t, warn)
	})
	t.Run(
		"should return an error because the obj does not match with the csv file",
		func(t *testing.T) {
			reader, err := os.Open(withValidCSVFile)
			require.NoError(t, err)
			defer func() { require.NoError(t, reader.Close()) }()

			csvReader, err := New(reader, ';')
			require.NoError(t, err)

			warn, err := csvReader.Decode(&DataTestInvalid{})
			require.ErrorIs(t, err, ErrOBJDecode)
			assert.Nil(t, warn)
		})
	t.Run("should read valid csv file", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		cases := []DataTestLong{
			{
				Name:      "Root",
				Type:      "Dog",
				MainColor: "black",
				Size:      "big",
			},
			{
				Name:      "Toto",
				Type:      "Human",
				MainColor: "blue",
				Size:      "small",
			},
		}

		for _, expected := range cases {
			data := DataTestLong{}
			warn, err := csvReader.Decode(&data)

			require.NoError(t, err)
			assert.Empty(t, warn)
			assert.Equal(t, expected, data)
		}
	})
	t.Run("should read valid csv file with warning", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		cases := []DataTestShort{
			{
				Name: "Root",
				Type: "Dog",
			},
			{
				Name: "Toto",
				Type: "Human",
			},
		}
		warnings := []map[string]string{
			{
				"main_color": "black",
				"size":       "big",
			},
			{
				"main_color": "blue",
				"size":       "small",
			},
		}

		for i, expected := range cases {
			data := DataTestShort{}
			warn, err := csvReader.Decode(&data)

			require.NoError(t, err)
			assert.EqualValues(t, warnings[i], warn)
			assert.Equal(t, expected, data)
		}
	})
}

func TestCSV_DecodeWithDecoder(t *testing.T) {
	t.Run("should return an error because the decoder is nil", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		err = csvReader.DecodeWithDecoder(nil)
		require.ErrorIs(t, err, ErrNilDecoder)
	})
	t.Run("should return an error when calling new instance func", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:     func(_ Decoder) (any, error) { return nil, assert.AnError },
			SaveInstanceFunc:    func(_ Decoder, _ any) error { return nil },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return nil },
		})
		require.NoError(t, err)

		err = csvReader.DecodeWithDecoder(decoder)
		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return an error when calling warning unused field func", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:     func(_ Decoder) (any, error) { return &DataTestLong{}, nil },
			SaveInstanceFunc:    func(_ Decoder, _ any) error { return nil },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return assert.AnError },
		})
		require.NoError(t, err)

		err = csvReader.DecodeWithDecoder(decoder)
		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return an error when calling save instance func", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:     func(_ Decoder) (any, error) { return &DataTestLong{}, nil },
			SaveInstanceFunc:    func(_ Decoder, _ any) error { return assert.AnError },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return nil },
		})
		require.NoError(t, err)

		err = csvReader.DecodeWithDecoder(decoder)
		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return an error to decode obj", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc:     func(_ Decoder) (any, error) { return &DataTestInvalid{}, nil },
			SaveInstanceFunc:    func(_ Decoder, _ any) error { return nil },
			WarningInstanceFunc: func(_ Decoder, _ Warning) error { return nil },
		})
		require.NoError(t, err)

		err = csvReader.DecodeWithDecoder(decoder)
		require.ErrorIs(t, err, ErrOBJDecode)
	})
	t.Run("should decode the csv file without warning", func(t *testing.T) {
		reader, err := os.Open(withValidCSVFile)
		require.NoError(t, err)
		defer func() { require.NoError(t, reader.Close()) }()

		csvReader, err := New(reader, ';')
		require.NoError(t, err)

		gotten := []DataTestLong{}
		expected := []DataTestLong{
			{
				Name:      "Root",
				Type:      "Dog",
				MainColor: "black",
				Size:      "big",
			},
			{
				Name:      "Toto",
				Type:      "Human",
				MainColor: "blue",
				Size:      "small",
			},
		}
		decoder, err := NewDecoder(ConfigDecoder{
			NewInstanceFunc: func(_ Decoder) (any, error) { return &DataTestLong{}, nil },
			SaveInstanceFunc: func(_ Decoder, obj any) error {
				if v, ok := obj.(*DataTestLong); ok {
					gotten = append(gotten, *v)

					return nil
				}

				return assert.AnError
			},
			WarningInstanceFunc: func(_ Decoder, warn Warning) error {
				assert.Empty(t, warn)

				return nil
			},
		})
		require.NoError(t, err)

		err = csvReader.DecodeWithDecoder(decoder)
		require.NoError(t, err)
		assert.Equal(t, expected, gotten)
	})
}

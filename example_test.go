package csv_test

import (
	"fmt"
	"os"

	"github.com/gofast-pkg/csv"
	"github.com/stretchr/testify/assert"
)

const withValidCSVFile = "testdata/valid.csv"

func ExampleNew() {
	type Model struct {
		Name      string `csv:"name"`
		Type      string `csv:"type"`
		MainColor string `csv:"main_color"`
		Size      string `csv:"size"`
	}

	reader, err := os.Open(withValidCSVFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = reader.Close(); err != nil {
			panic(err)
		}
	}()

	csvReader, err := csv.New(reader, ';')
	if err != nil {
		panic(err)
	}

	exampleContextKey := "example-key"

	list := []Model{}
	cfg := csv.ConfigDecoder{
		NewInstanceFunc: func(_ csv.Decoder) (any, error) { return &Model{}, nil },
		SaveInstanceFunc: func(dec csv.Decoder, obj any) error {
			// manipulate data in the decoder context
			v, ok := dec.ContextGet(exampleContextKey)
			if !ok {
				return assert.AnError
			}
			if count, ok := v.(*int); ok {
				(*count)++
				dec.ContextSet(exampleContextKey, count)
			}

			if v, ok := obj.(*Model); ok {
				list = append(list, *v)

				return nil
			}

			return assert.AnError
		},
		WarningInstanceFunc: func(_ csv.Decoder, warn csv.Warning) error {
			for k, v := range warn {
				fmt.Printf("warning with key %s value %s\n", k, v)
			}

			return nil
		},
	}
	decoder, err := csv.NewDecoder(cfg)
	if err != nil {
		panic(err)
	}

	// records some datas in the decoder context
	count := int(0)
	decoder.ContextSet(exampleContextKey, &count)

	err = csvReader.DecodeWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	for _, v := range list {
		fmt.Println(v)
	}

	fmt.Println(count)

	// Output:
	// {Root Dog black big}
	// {Toto Human blue small}
	// 2
}

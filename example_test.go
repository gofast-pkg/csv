package csv_test

import (
	"fmt"
	"os"

	"github.com/gofast-pkg/csv"
)

const testFilePath = "testdata/testfile.csv"

func ExampleNewWarning() {
	w := csv.NewWarning()
	w["key"] = []string{"value1", "value2"}
	for key, values := range w {
		for _, value := range values {
			fmt.Println(key, value)
		}
	}
	// Output:
	// key value1
	// key value2
}

func ExampleNew() {
	type Model struct {
		Name      string `csv:"name"`
		Type      string `csv:"type"`
		MainColor string `csv:"main_color"`
		Size      string `csv:"size"`
	}

	reader, err := os.Open(testFilePath)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	csvReader, err := csv.New(reader, ';')
	if err != nil {
		panic(err)
	}

	list := []Model{}
	cfg := csv.ConfigDecoder{
		NewInstanceFunc: func() any { return &Model{} },
		SaveInstanceFunc: func(obj any, d csv.Decoder) error {
			if v, ok := d.ContextGet("type"); ok {
				obj.(*Model).Type = v
			}
			list = append(list, *(obj.(*Model)))

			return nil
		},
	}
	decoder, err := csv.NewDecoder(cfg)
	if err != nil {
		panic(err)
	}

	warn, err := csvReader.DecodeWithDecoder(decoder)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(warn))
	for _, v := range list {
		fmt.Println(v)
	}
	// Output:
	// 0
	// {Root Dog black big}
	// {Toto Human blue small}
}

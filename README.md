# CSV

[![Static Badge](https://img.shields.io/badge/project%20use%20codesystem-green?link=https%3A%2F%2Fgithub.com%2Fgofast-pkg%2Fcodesystem)](https://github.com/gofast-pkg/codesystem)
![Build](https://github.com/gofast-pkg/csv/actions/workflows/ci.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/gofast-pkg/csv.svg)](https://pkg.go.dev/github.com/gofast-pkg/csv)
[![codecov](https://codecov.io/gh/gofast-pkg/csv/branch/main/graph/badge.svg?token=7TCE3QB21E)](https://codecov.io/gh/gofast-pkg/csv)
[![Release](https://img.shields.io/github/release/gofast-pkg/csv?style=flat-square)](https://github.com/gofast-pkg/csv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofast-pkg/csv)](https://goreportcard.com/report/github.com/gofast-pkg/csv)
[![codebeat badge](https://codebeat.co/badges/1771f3ed-bead-4953-bd72-da2c8819962c)](https://codebeat.co/projects/github-com-gofast-pkg-csv-main)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgofast-pkg%2Fcsv.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgofast-pkg%2Fcsv?ref=badge_shield)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/gofast-pkg/csv/blob/main/LICENSE)

Package CSV is a package to parse a full csv reader inside a slice structure.
Allow to configure a decoder to transform directly csv entries into csv structure.

I had this specific need because i would to parse many csv files for differents type structure.

This package allow to configure dynamicly the structure builder to parse each entries files.

If you need an easy parser to decode csv file basicly, prefer use the [csvutil package](https://pkg.go.dev/github.com/jszwec/csvutil)

## Install

``` bash
$> go get github.com/gofast-pkg/csv@latest
```

## Usage

``` Golang
import github.com/gofast-pkg/csv

func main() {
  reader, err := os.Open(filePath)
  if err != nil {
    panic(err)
  }
  defer reader.Close()

  csvReader, err := csv.New(reader, ';')
  if err != nil {
    panic(err)
  }
}
```

Examples are provided on the [go doc reference](https://pkg.go.dev/github.com/gofast-pkg/csv)

## Contributing

&nbsp;:grey_exclamation:&nbsp; Use issues for everything

Read more informations with the [CONTRIBUTING_GUIDE](./.github/CONTRIBUTING.md)

For all changes, please update the CHANGELOG.txt file by replacing the existant content.

Thank you &nbsp;:pray:&nbsp;&nbsp;:+1:&nbsp;

<a href="https://github.com/gofast-pkg/csv/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=gofast-pkg/csv" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## Licence

[MIT](https://github.com/gofast-pkg/csv/blob/main/LICENSE)

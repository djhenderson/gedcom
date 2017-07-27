# gedcom

Go package to parse GEDCOM files into Go structures.
Also write these to GEDCOM files.

[![GoDoc](https://godoc.org/github.com/djhenderson/gedcom?status.svg)](https://godoc.org/github.com/djhenderson/gedcom)

## Note

This is a fork of [github.com/iand/gedcom](https://github.com/iand/gedcom) with major revisions and additions.

## Usage

The package provides a Decoder with a single Decode method that returns a Gedcom struct. Use the NewDecoder method to create a new decoder.

This example shows how to parse a GEDCOM file and list all the individuals. In this example the entire input file is read into memory, but the decoder is streaming so it should be able to deal with very large files: just pass an appropriate Reader.


	package main

	import (
		"bytes"
		"github.com/djhenderson/gedcom"
		"io/ioutil"
	)

    func main() {
		data, _ := ioutil.ReadFile("testdata/kennedy.ged") // read all data

		d := gedcom.NewDecoder(bytes.NewReader(data)) // make the decoder

		g, _ := d.Decode() // decode all data

		for _, rec := range g.Individual { // for all individuals
			if len(rec.Name) > 0 { // check that individual has a name
				println(rec.Name[0].Name) // print the name
			}
		}
	}

The structures produced by the Decoder are in [types.go](types.go) and correspond roughly 1:1 to the structures in the [GEDCOM specification](http://homepages.rootsweb.ancestry.com/~pmcbride/gedcom/55gctoc.htm).

This package does not implement the entire GEDCOM specification, I'm still working on it. It's been tested with all available of GEDCOM files. It has been extensively tested with non-ASCII character sets and with pathalogical cases such as the [GEDCOM 5.5 Torture Test Files](http://www.geditcom.com/gedcom.html).

## Installation

Simply run

	go get github.com/djhenderson/gedcom

Documentation is at [https://godoc.org/github.com/djhenderson/gedcom](https://godoc.org/github.com/djhenderson/gedcom)

## Authors

* [Ian Davis](https://github.com/iand) - <http://iandavis.com/>
* [Doug Henderson](https://github.com/djhenderson)

## Contributing

* Do submit your changes as a pull request
* Do your best to adhere to the existing coding conventions and idioms.
* Do run `go fmt` on the code before committing
* Do feel free to add yourself to the [`CREDITS`](CREDITS) file and the
  corresponding Contributors list in the the [`README.md`](README.md).
  Alphabetical order applies.
* Don't touch the [`AUTHORS`](AUTHORS) file. An existing author will add you if
  your contributions are significant enough.
* Do note that in order for any non-trivial changes to be merged (as a rule
  of thumb, additions larger than about 15 lines of code), an explicit
  Public Domain Dedication needs to be on record from you. Please include
  a copy of the statement found in the [`WAIVER`](WAIVER) file with your pull request

## License

This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.

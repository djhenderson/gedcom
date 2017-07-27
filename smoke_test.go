package gedcom

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// TestSmoke test the example code for [README.md](README.md).
// In this version, check all errors.
func TestSmoke(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/kennedy.ged") // read all data
	if err != nil {
		t.Fatal("ioutl.ReadFile failed:", err)
	}

	d := NewDecoder(bytes.NewReader(data)) // make the decoder

	g, err := d.Decode() // decode all data
	if err != nil {
		t.Fatal("Decode failed (see log.txt):", err)
	}

	for _, rec := range g.Individual { // for all individuals
		if len(rec.Name) > 0 { // check that individual has a name
			println(rec.Name[0].Name) // print the name
		}
	}
}

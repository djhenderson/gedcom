/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

import (
	"io"
	"testing"
)

type example struct {
	input []byte
	level int
	tag   string
	value string
	xref  string
}

var examples = []example{
	{[]byte("1 SEX F\n"), 1, `SEX`, `F`, ""},
	{[]byte(" 1 SEX F\n"), 1, `SEX`, `F`, ""},
	{[]byte("  \r\n\t 1 SEX F\n"), 1, `SEX`, `F`, ""},
	{[]byte("  \r\n\t 1     SEX      F\n"), 1, `SEX`, `F`, ""},
	{[]byte("1 SEX F\r"), 1, `SEX`, `F`, ""},
	{[]byte("1 SEX F \r"), 1, `SEX`, `F `, ""},
	{[]byte("0 HEAD\r"), 0, `HEAD`, ``, ""},
	{[]byte("0 @OTHER@ SUBM\n"), 0, `SUBM`, ``, "OTHER"},
	// leading BOM
	//{[]byte("\xef\xbb\xbf0 HEAD\r"), 0, `HEAD`, ``, ""},
	// email address with escaped @
	//{[]byte("1 EMAIL toaddress@@example.com"), 1, `EMAIL`, `toaddress@@example.com`, ""},
	// email address with naked @
	//{[]byte("1 EMAIL toaddress@example.com"), 1, `EMAIL`, `toaddress@example.com`, ""},
	// value includes text and xref
	//{[]byte("    2 SOUR Book, Simple @S10@"), 2, `SOUR`, `Book, Simple`, "S10"},
	// value includes text and xref
	//{[]byte("  1 ROLE CHIL @I1@"), 1, `ROLE`, `CHIL`, "I1"},
	// value includes non-ASCII character
	//{[]byte("1 NAME Æthelreda //"), 1, `NAME`, `Æthelreda`, ""},
}

func TestNextTagFound(t *testing.T) {
	s := &scanner{}
	for _, ex := range examples {
		s.reset()
		offset, err := s.nextTag(ex.input)

		if err != nil {
			t.Fatalf(`nextTag for "%s" returned error "%v", expected no error`, ex.input, err)
		}

		if offset == 0 {
			t.Fatalf(`nextTag for "%s" did not find tag, expected it to find`, ex.input)
		}

		if s.level != ex.level {
			t.Errorf(`nextTag for "%s" returned level %d, expected %d`, ex.input, s.level, ex.level)
		}

		if string(s.tag) != ex.tag {
			t.Errorf(`nextTag for "%s" returned tag "%s", expected "%s"`, ex.input, s.tag, ex.tag)
		}

		if string(s.value) != ex.value {
			t.Errorf(`nextTag for "%s" returned value "%s", expected "%s"`, ex.input, s.value, ex.value)
		}

		if string(s.xref) != ex.xref {
			t.Errorf(`nextTag for "%s" returned xref "%s", expected "%s"`, ex.input, s.xref, ex.xref)
		}

	}

}

var examplesNot = [][]byte{
	[]byte("1 SEX F"),
	[]byte(" 1 SEX F "),
}

func TestNextTagNotFound(t *testing.T) {
	s := &scanner{}
	for _, ex := range examplesNot {
		s.reset()
		_, err := s.nextTag(ex)

		if err != io.EOF {
			t.Fatalf(`nextTag for "%s" returned unexpected error "%v", expected io.EOF`, ex, err)
		}

	}

}

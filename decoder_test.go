/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
)

var (
	data []byte
)

func init() {
	const logFileName = "./log.txt"

	log.SetFlags(log.Flags() | log.Lshortfile)
	logFile, err := os.Create(logFileName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.SetOutput(logFile)

	data, err = ioutil.ReadFile("testdata/allged.ged")
	if err != nil {
		panic(err)
	}
}

func TestStructuresAreInitialized(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()

	if err != nil {
		t.Fatalf("Result of decoding gedcom gave error, expected no error")
	}

	if g == nil {
		t.Fatalf("Result of decoding gedcom was nil, expected valid object")
	}
	if g.Individual == nil {
		t.Fatalf("Individual list was nil, expected valid slice")
	}

	if g.Family == nil {
		t.Fatalf("Family list was nil, expected valid slice")
	}

	if g.Media == nil {
		t.Fatalf("Media list was nil, expected valid slice")
	}

	if g.Repository == nil {
		t.Fatalf("Repository list was nil, expected valid slice")
	}

	if g.Source == nil {
		t.Fatalf("Source list was nil, expected valid slice")
	}

	if g.Submitter == nil {
		t.Fatalf("Submitter list was nil, expected valid slice")
	}

}

func TestIndividual(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()

	if len(g.Individual) != 8 {
		t.Fatalf("Individual list length was %d, expected 8", len(g.Individual))
	}

	i1 := g.Individual[0]

	if i1.Xref != "@PERSON1@" {
		t.Logf("i1.Xref (%T) = %q\n", i1.Xref, i1.Xref)
		t.Errorf(`Individual 0 xref was %q, expected "@PERSON1@"`, i1.Xref)
	}
}

func TestIndiSex(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	if i1.Sex != "M" {
		t.Errorf(`Individual 0 sex "%q" names, expected "M"`, i1.Sex)
	}
}

func TestIndiName1(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	if len(i1.Name) != 2 {
		t.Fatalf(`Individual 0 had %d names, expected 2`, len(i1.Name))
	}
}

func TestIndiName2(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	// 1 NAME given name /surname/jr.
	// 2 SOUR @SOURCE1@
	// 3 PAGE 42
	// 3 DATA
	// 4 DATE BEF 1 JAN 1900
	// 4 TEXT a sample text
	// 5 CONT Sample text continued here. The word TE
	// 5 CONC ST should not be broken!
	// 3 QUAY 0
	// 3 NOTE A note
	// 4 CONT Note continued here. The word TE
	// 4 CONC ST should not be broken!
	// 2 NOTE Personal Name note
	// 3 CONT Note continued here. The word TE
	// 3 CONC ST should not be broken!

	name1 := &NameRecord{
		Level: 1,
		Name:  "given name /surname/jr.",
		Citation: CitationRecords{
			&CitationRecord{
				Level:   2,
				Value:   "@SOURCE1@",
				Page:    "42",
				Quality: "0",
				Data: DataRecords{
					&DataRecord{
						Level: 3,
						Date:  "BEF 1 JAN 1900",
						Text:  "a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: NoteRecords{
					&NoteRecord{
						Level: 3,
						Note:  "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: NoteRecords{
			&NoteRecord{
				Level: 2,
				Note:  "Personal Name note\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	diff := deep.Equal(i1.Name[0], name1)
	if diff != nil {
		fmt.Println(strings.Repeat("=", 80))
		for n, d := range diff {
			fmt.Println(n, strings.Repeat("-", 80))
			fmt.Println(d)
		}
		fmt.Println(strings.Repeat("=", 80))
		t.Error(diff)
	}
}

func TestIndiEvents1(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	if len(i1.Event) != 31 {
		t.Fatalf(`Individual 0 had %d events, expected 31`, len(i1.Event))
	}
}

func TestIndiEvents2(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	// 1 BIRT
	// 2 DATE 31 DEC 1997
	// 2 PLAC The place
	// 2 SOUR @SOURCE1@
	// 3 PAGE 42
	// 3 DATA
	// 4 DATE 31 DEC 1900
	// 4 TEXT a sample text
	// 5 CONT Sample text continued here. The word TE
	// 5 CONC ST should not be broken!
	// 3 QUAY 2
	// 3 NOTE A note
	// 4 CONT Note continued here. The word TE
	// 4 CONC ST should not be broken!
	// 2 NOTE BIRTH event note (the event of entering into life)
	// 3 CONT Note continued here. The word TE
	// 3 CONC ST should not be broken!
	// 2 FAMC @PARENTS@

	event1 := &EventRecord{
		Level: 1,
		Tag:   "BIRT",
		Date: &DateRecord{
			Level: 2,
			Tag:   "DATE",
			Date:  "31 DEC 1997",
		},
		Place: &PlaceRecord{
			Level: 2,
			Tag:   "PLAC",
			Name:  "The place",
		},
		Parents: FamilyLinks{
			&FamilyLink{
				Level: 2,
				Tag:   "FAMC",
				Value: "@PARENTS@",
			},
		},
		Citation: []*CitationRecord{
			&CitationRecord{
				Level:   2,
				Value:   "@SOURCE1@",
				Page:    "42",
				Quality: "2",
				Data: []*DataRecord{
					&DataRecord{
						Level: 3,
						Date:  "31 DEC 1900",
						Text:  "a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: []*NoteRecord{
					&NoteRecord{
						Level: 3,
						Note:  "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: []*NoteRecord{
			&NoteRecord{
				Level: 2,
				Note:  "BIRTH event note (the event of entering into life)\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	diff := deep.Equal(i1.Event[0], event1)
	if diff != nil {
		//		t.Errorf("Individual 0 event 0 was: %q\nExpected: %q\n",
		//			spew.Sdump(i1.Event[0]), spew.Sdump(event1))
		fmt.Println(strings.Repeat("=", 80))
		for n, d := range diff {
			fmt.Println(n, strings.Repeat("-", 80))
			fmt.Println(d)
		}
		fmt.Println(strings.Repeat("=", 80))
		t.Error(diff)
	}
}

func TestIndiAttribute1(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	if len(i1.Attribute) != 7 {
		t.Fatalf(`Individual 0 had %d attributes, expected 7`, len(i1.Attribute))
	}
}

func TestIndiAttribute2(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	att1 := &AttributeRecord{
		Level: 1,
		Tag:   "CAST",
		Value: "Cast name",
		Date: &DateRecord{
			Level: 2,
			Tag:   "DATE", // fix
			Date:  "31 DEC 1997",
		},
		Place: &PlaceRecord{
			Level: 2,
			Tag:   "PLAC", // fix
			Name:  "The place",
		},
		Citation: []*CitationRecord{
			&CitationRecord{
				Level:   2,
				Value:   "@SOURCE1@",
				Page:    "42",
				Quality: "3",
				Data: []*DataRecord{
					&DataRecord{
						Level: 3,
						Date:  "31 DEC 1900",
						Text:  "a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: []*NoteRecord{
					&NoteRecord{
						Level: 3,
						Note:  "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: []*NoteRecord{
			&NoteRecord{
				Level: 2,
				Note:  "CASTE event note (the name of an individual's rank or status in society, based   on racial or religious differences, or differences in wealth, inherited   rank, profession, occupation, etc)\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	if len(i1.Attribute) > 0 {

		if !reflect.DeepEqual(i1.Attribute[0], att1) {
			t.Errorf("Individual 0 attribute 0 was: \n%q\nExpected: \n%q\n",
				spew.Sdump(i1.Attribute[0]), spew.Sdump(att1))
		}
	}
}

func TestIndiParents1(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	if len(i1.Parents) != 2 {
		t.Fatalf(`Individual 0 had %d parent families, expected 2`, len(i1.Parents))
	}

}

func TestIndiParents2(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()
	i1 := g.Individual[0]

	//0 @PARENTS@ FAM
	//1 HUSB @PERSON5@
	//1 CHIL @PERSON1@

	fam1 := &FamilyLink{
		Level: 1,
		Tag:   "FAMC",
		Value: "@PARENTS@",
		Note: []*NoteRecord{
			&NoteRecord{
				Level: 2,
				Note:  "Note about the link to parents\nNote continued here. The word TEST should not be broken!",
			},
			&NoteRecord{
				Level: 2,
				Note:  "Another note about the link to parents\nNote continued here. The word TEST should not be broken!",
			},
		},
		Pedigree: &PedigreeRecord{
			Level:    2,
			Pedigree: "birth",
		},
	}

	if !reflect.DeepEqual(i1.Parents[0], fam1) {
		t.Errorf("Family 0 parents 0 was: \n%q\nExpected: \n%q\n",
			spew.Sdump(i1.Parents[0]), spew.Sdump(fam1))
	}

}

func TestSubmitter(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()

	if len(g.Submitter) != 1 {
		t.Fatalf("Submitter list length was %d, expected 1", len(g.Submitter))
	}

}

func TestFamily(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()

	if len(g.Family) != 4 {
		t.Fatalf("Family list length was %d, expected 4", len(g.Family))
	}

}

func TestSource1(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()

	if len(g.Source) != 1 {
		t.Fatalf("Source list length was %d, expected 1", len(g.Source))
	}

}

func TestSource2(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, _ := d.Decode()

	//0 @SOURCE1@ SOUR
	sour1 := &SourceRecord{
		Xref:  "@SOURCE1@",
		Title: "Title of source\nTitle continued here. The word TEST should not be broken!",
		Author: &AuthorRecord{
			Level:  1,
			Author: "Author of source\nAuthor continued here. The word TEST should not be broken!",
		},
		Publication:  "Source publication facts\nPublication facts continued here. The word TEST should not be broken!",
		Abbreviation: "Short title",
		Note: NoteRecords{
			&NoteRecord{
				Level: 1,
				Note:  "A note about the family\nNote continued here. The word TEST should not be broken!",
			},
		},
		Text: "Citation from source\nCitation continued here. The word TEST should not be broken!",
		Data: &DataRecord{Level: 1,
			Event: EventRecords{
				&EventRecord{
					Level: 2,
					Tag:   "EVEN",
					Value: "BIRT, CHR",
					Date: &DateRecord{
						Level: 3,
						Tag:   "DATE", // fix
						Date:  "FROM 1 JAN 1980 TO 1 FEB 1982",
					},
					Place: &PlaceRecord{
						Level: 3,
						Tag:   "PLAC", // fix
						Name:  "Place",
					},
				},
				&EventRecord{
					Level: 2,
					Tag:   "EVEN",
					Value: "DEAT",
					Date: &DateRecord{
						Level: 3,
						Tag:   "DATE", // fix
						Date:  "FROM 1 JAN 1980 TO 1 FEB 1982",
					},
					Place: &PlaceRecord{
						Level: 3,
						Tag:   "PLAC", // fix
						Name:  "Another place",
					},
				},
			},
			Agency: "Resposible agency",
			Note: NoteRecords{
				&NoteRecord{
					Level: 2,
					Note:  "A note about whatever\nNote continued here. The word TEST should not be broken!",
				},
			},
		},

		Media: MediaLinks{
			&MediaLink{
				Media: &MediaRecord{
					Level:    1,
					Format:   "bmp",
					Title:    "A bmp picture",
					FileName: `\\network\drive\path\file name.bmp`,
					Note: NoteRecords{
						&NoteRecord{
							Level: 2,
							Note:  "A note\nNote continued here. The word TEST should not be broken!",
						},
					},
				},
			},
		},
		Change: &ChangeRecord{
			Level: 1,
			Date: &DateRecord{
				Level: 2,
				Tag:   "DATE", // fix
				Date:  "1 APR 1998",
				Time:  "12:34:56.789",
			},
			Note: NoteRecords{
				&NoteRecord{
					Level: 2,
					Note:  "A note\nNote continued here. The word TEST should not be broken!",
				},
			},
		},
		//1 _MYOWNTAG This is a non-standard tag. Not recommended but allowed
	}

	if !reflect.DeepEqual(g.Source[0], sour1) {
		a1 := g.Source[0].String()
		a2 := sour1.String()
		t.Log(a1)
		t.Log(a2)
		t.Errorf("Family 0 parents 0 was: \n%q\nExpected: \n%q\n",
			spew.Sdump(g.Source[0]), spew.Sdump(sour1))
	}

}

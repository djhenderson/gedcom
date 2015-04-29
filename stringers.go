/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

import (
	"fmt"
	"log"
	"strings"
)

func LongString(level int, xref string, tag string, longString string) []string {
	var s string
	var ss []string

	parts := strings.Split(longString, "\n")
	for i, part := range parts {
		sXref := ""
		if i == 0 && xref != "" {
			sXref = fmt.Sprintf("@%s@ ", xref)
		}
		s = fmt.Sprintf("%d %s%s %s", level, sXref, tag, part)
		ss = append(ss, s)
		if i == 0 {
			tag = "CONT"
			level += 1
		}
	}
	return ss
}

func (r *AddressRecord) String() string {
	var ss []string
	var s string

	sas := LongString(r.Level, "", "ADDR", r.Full)
	ss = append(ss, sas...)

	if r.Line1 != "" {
		s = fmt.Sprintf("%d ADR1 %s", r.Level+1, r.Line1)
		ss = append(ss, s)
	}

	if r.Line2 != "" {
		s = fmt.Sprintf("%d ADR2 %s", r.Level+1, r.Line2)
		ss = append(ss, s)
	}

	if r.City != "" {
		s = fmt.Sprintf("%d CITY %s", r.Level+1, r.City)
		ss = append(ss, s)
	}

	if r.State != "" {
		s = fmt.Sprintf("%d STAE %s", r.Level+1, r.State)
		ss = append(ss, s)
	}

	if r.PostalCode != "" {
		s = fmt.Sprintf("%d POST %s", r.Level+1, r.PostalCode)
		ss = append(ss, s)
	}

	if r.Country != "" {
		s = fmt.Sprintf("%d CTRY %s", r.Level+1, r.Country)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r *BusinessRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d CORP %s", r.Level, r.BusinessName)
	ss = append(ss, s)

	if r.Address != nil {
		s = fmt.Sprintf("%s", r.Address)
		ss = append(ss, s)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			s = fmt.Sprintf("%d PHON %s", r.Level+1, phone)
			ss = append(ss, s)
		}
	}

	return strings.Join(ss, "\n")
}

func (r *ChangeRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d CHAN %s", r.Level, r.Date)
	ss = append(ss, s)

	if r.Date != nil {
		s = fmt.Sprintf("%s", r.Date)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r *CharacterSetRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d CHAR %s", r.Level, r.CharacterSet)
	ss = append(ss, s)

	if r.Version != "" {
		s = fmt.Sprintf("%d VERS %s", r.Level+1, r.Version)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r *ChildStatusRecord) String() string {
	log.Fatal("ChildStatusRecord stringer not implemented")
	return ""
}

func (r ChildStatusRecords) String() string {
	var ss []string
	var s string
	log.Printf("ChildStatusRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *CitationRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d SOUR @%s@", r.Level, r.Source.Xref)
	ss = append(ss, s)

	if r.Page != "" {
		s = fmt.Sprintf("%d PAGE %s", r.Level+1, r.Page)
		ss = append(ss, s)
	}

	if r.Data != nil {
		for _, data := range r.Data {
			s = fmt.Sprintf("%s", data)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Quality != "" {
		s = fmt.Sprintf("%d PAGE %s", r.Level+1, r.Quality)
		ss = append(ss, s)
	}

	if r.Media != nil {
		for _, media := range r.Media {
			s = fmt.Sprintf("%d OBJE @%s@", r.Level+1, media.Xref)
			ss = append(ss, s)
			s = fmt.Sprintf("%s", media)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r CitationRecords) String() string {
	var ss []string
	var s string

	for _, note := range r {
		s = fmt.Sprintf("%s", *note)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *DataRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d DATA %s", r.Level, r.Data)
	ss = append(ss, s)

	if r.Date != "" {
		s = fmt.Sprintf("%d DATE %s", r.Level+1, r.Date)
		ss = append(ss, s)
	}

	if r.Copyright != "" {
		s = fmt.Sprintf("%d COPR %s", r.Level+1, r.Copyright)
		ss = append(ss, s)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			s = fmt.Sprintf("%d TEXT %s", r.Level+1, text)
			ss = append(ss, s)
		}
	}

	return strings.Join(ss, "\n")
}

func (r DataRecords) String() string {
	var ss []string
	var s string
	log.Printf("DataRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *DateRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d DATE %s", r.Level+1, r.Date)
	ss = append(ss, s)

	if r.Time != "" {
		s = fmt.Sprintf("%d TIME %s", r.Level+1, r.Time)
		ss = append(ss, s)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			s = fmt.Sprintf("%d TEXT %s", r.Level+1, text)
			ss = append(ss, s)
		}
	}

	return strings.Join(ss, "\n")
}

func (r EventDefinitionRecord) String() string {
	log.Fatal("EventDefinition Stringer not implemented")
	return ""
}

func (r EventDefinitionRecords) String() string {
	var ss []string
	var s string

	log.Printf("EventDefinitionRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *EventRecord) String() string {
	var ss []string
	var s string

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}
	s = fmt.Sprintf("%d %s%s %s", r.Level, id, r.Tag, r.Value)
	ss = append(ss, s)

	if r.Type != "" {
		s = fmt.Sprintf("%d TYPE %s", r.Level+1, r.Type)
		ss = append(ss, s)
	}

	if r.Date != nil {
		s = fmt.Sprintf("%d DATE %s", r.Level+1, r.Date)
		ss = append(ss, s)
	}

	if r.Place != nil {
		s = fmt.Sprintf("%s", r.Place)
		if s != "" {
			ss = append(ss, s)
		}
	}

	if r.Address != nil {
		s = fmt.Sprintf("%s", r.Address)
		if s != "" {
			ss = append(ss, s)
		}
	}

	if r.Parents != nil {

	}
	if r.Age != "" {
		s = fmt.Sprintf("%d AGE %s", r.Level+1, r.Age)
		ss = append(ss, s)
	}

	if r.Agency != "" {
		s = fmt.Sprintf("%d AGNC %s", r.Level+1, r.Agency)
		ss = append(ss, s)
	}

	if r.Cause != "" {
		s = fmt.Sprintf("%d CAUS %s", r.Level+1, r.Cause)
		ss = append(ss, s)
	}

	if r.Temple != "" {
		s = fmt.Sprintf("%d TEMP %s", r.Level+1, r.Temple)
		ss = append(ss, s)
	}

	if r.Media != nil {
		for _, media := range r.Media {
			s = fmt.Sprintf("%s", media)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Citation != nil {
		s = fmt.Sprintf("%s", r.Citation)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = fmt.Sprintf("%s", r.Change)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r EventRecords) String() string {
	var ss []string
	var s string

	log.Printf("EventRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}
	if ss == nil {
		return ""
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *FamilyLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d %s @%s@", r.Level, r.Tag, r.Family.Xref)
	ss = append(ss, s)

	if r.Pedigree != "" {
		s = fmt.Sprintf("%d PEDI %s", r.Level+1, r.Pedigree)
		ss = append(ss, s)
	}

	if r.Adopted != "" {
		s = fmt.Sprintf("%d ADOP %s", r.Level+1, r.Adopted)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r FamilyLinks) String() string {
	var ss []string
	var s string

	log.Printf("FamilyLinks type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}
	if ss == nil {
		return ""
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *FamilyRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d @%s@ FAM", r.Level, r.Xref)
	ss = append(ss, s)

	if r.Husband != nil {
		s = fmt.Sprintf("%s", r.Husband)
		if s != "" {
			ss = append(ss, s)
		}
	}

	if r.Wife != nil {
		s = fmt.Sprintf("%s", r.Wife)
		if s != "" {
			ss = append(ss, s)
		}
	}

	if r.NumChildren > 0 {
		s = fmt.Sprint("%d NCHI %d", r.Level+1, r.NumChildren)
		ss = append(ss, s)
	}

	if r.Child != nil {
		for _, child := range r.Child {
			s = fmt.Sprintf("%s", child)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Event != nil {
		for _, event := range r.Event {
			s = fmt.Sprintf("%s", event)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Citation != nil {
		s = fmt.Sprintf("%s", r.Citation)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = fmt.Sprintf("%s", r.Change)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r FamilyRecords) String() string {
	var ss []string
	var s string

	log.Printf("FamilyRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *GedcomRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d GEDC", r.Level)
	ss = append(ss, s)

	s = fmt.Sprintf("%d VERS %s", r.Level+1, r.Version)
	ss = append(ss, s)

	s = fmt.Sprintf("%d FORM %s", r.Level+1, r.Form)
	ss = append(ss, s)

	return strings.Join(ss, "\n")

}

func (r *HeaderRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d HEAD", r.Level)
	ss = append(ss, s)

	if r.CharacterSet != nil {
		s = fmt.Sprintf("%s", r.CharacterSet)
		ss = append(ss, s)
	}

	if r.SourceSystem != nil {
		s = fmt.Sprintf("%s", r.SourceSystem)
		ss = append(ss, s)
	}

	if r.Destination != "" {
		s = fmt.Sprintf("%d DEST %s", r.Level+1, r.Destination)
		ss = append(ss, s)
	}

	if r.Date != nil {
		s = fmt.Sprintf("%d DATE %s", r.Level+1, r.Date)
		ss = append(ss, s)
	}

	if r.Time != "" {
		s = fmt.Sprintf("%d TIME %s", r.Level+1, r.Time)
		ss = append(ss, s)
	}

	if r.FileName != "" {
		s = fmt.Sprintf("%d FILE %s", r.Level+1, r.FileName)
		ss = append(ss, s)
	}

	if r.Gedcom != nil {
		s = fmt.Sprintf("%s", r.Gedcom)
		ss = append(ss, s)
	}

	if r.Language != "" {
		s = fmt.Sprintf("%d LANG %s", r.Level+1, r.Language)
		ss = append(ss, s)
	}

	if r.Copyright != "" {
		s = fmt.Sprintf("%d COPR %s", r.Level+1, r.Copyright)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", &r.Note)
		ss = append(ss, s)
	}

	if r.Submitter != nil {
		s = fmt.Sprintf("%d SUBM @%s@", r.Level+1, r.Xref)
		ss = append(ss, s)
	}

	if r.Submission != nil {
		s = fmt.Sprintf("%d SUBN @%s@", r.Level+1, r.Xref)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")

}

func (r *IndividualLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d %s @%s@", r.Level+1, r.Tag, r.Xref)
	ss = append(ss, s)

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r IndividualLinks) String() string {
	var ss []string
	var s string

	log.Printf("IndividualLinks type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *IndividualRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d @%s@ INDI", r.Level, r.Xref)
	ss = append(ss, s)

	for _, name := range r.Name {
		s = fmt.Sprintf("%s", name)
		ss = append(ss, s)
	}

	if r.Sex != "" {
		s = fmt.Sprintf("%d SEX %s", r.Level+1, r.Sex)
		ss = append(ss, s)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			s = fmt.Sprintf("%s", event)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	//Attribute []*EventRecord      // ?

	if r.Parents != nil {
		for _, famc := range r.Parents {
			s = fmt.Sprintf("%s", famc)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Family != nil {
		for _, fams := range r.Family {
			s = fmt.Sprintf("%s", fams)
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	//Address   []*AddressRecord    // INDI.ADDR
	//Health    string              // INDI.HEAL
	//History   []*HistoryRecord    // INDI.HIST

	if r.Media != nil {
		s = fmt.Sprintf("%s", r.Media)
		ss = append(ss, s)
	}

	if r.Citation != nil {
		s = fmt.Sprintf("%s", r.Citation)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = fmt.Sprintf("%s", r.Change)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r IndividualRecords) String() string {
	var ss []string
	var s string

	log.Printf("IndividualRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *MediaLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d OBJE @%s@", r.Level, r.Xref)
	ss = append(ss, s)

	if r.Format != "" {
		s = fmt.Sprintf("%d FORM %s", r.Level+1, r.Format)
		ss = append(ss, s)
	}

	if r.Title != "" {
		s = fmt.Sprintf("%d TITL %s", r.Level+1, r.Title)
		ss = append(ss, s)
	}

	if r.FileName != "" {
		s = fmt.Sprintf("%d FILE %s", r.Level+1, r.FileName)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", &r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r MediaLinks) String() string {
	var ss []string
	var s string

	log.Printf("MediaLinks type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *MediaRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d OBJE", r.Level)
	ss = append(ss, s)

	if r.Format != "" {
		s = fmt.Sprintf("%d FORM %s", r.Level+1, r.Format)
		ss = append(ss, s)
	}

	if r.Title != "" {
		s = fmt.Sprintf("%d TITL %s", r.Level+1, r.Title)
		ss = append(ss, s)
	}

	if r.FileName != "" {
		s = fmt.Sprintf("%d FILE %s", r.Level+1, r.FileName)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", &r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r MediaRecords) String() string {
	var ss []string
	var s string

	log.Printf("MediaRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *NameRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d NAME %s", r.Level+1, r.Name)
	ss = append(ss, s)

	if r.Citation != nil {
		s = fmt.Sprintf("%s", r.Citation)
		ss = append(ss, s)
	}

	if r.Note != nil {
		for _, note := range r.Note {
			s = fmt.Sprintf("%s", note.Note)
			ss = append(ss, s)
		}
	}

	return strings.Join(ss, "\n")
}

func (r NameRecords) String() string {
	var ss []string
	var s string

	log.Printf("NameRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *NoteRecord) String() string {
	var ss []string
	var s string

	//s = fmt.Sprintf("%d NOTE %s", r.Level, r.Note)
	//ss = append(ss, s)
	sas := LongString(r.Level, r.Xref, "NOTE", r.Note)
	ss = append(ss, sas...)

	if r.Citation != nil {
		s = fmt.Sprintf("%s", r.Citation)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r NoteRecords) String() string {
	var ss []string
	var s string

	log.Printf("NoteRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *PlacePartRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d %s %s", r.Level, r.Tag, r.Part)
	ss = append(ss, s)

	if r.Jurisdiction != "" {
		s = fmt.Sprintf("%d JURI %s", r.Level+1, r.Jurisdiction)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r *PlaceRecord) String() string {
	var ss []string
	var s string

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}

	s = fmt.Sprintf("%d %sPLAC %s", r.Level, id, r.Name)
	ss = append(ss, s)

	if r.Name != "" {
		s = fmt.Sprintf("%d PLAS %s", r.Level+1, r.Name)
		ss = append(ss, s)
	}

	if r.Modifier != "" {
		s = fmt.Sprintf("%d PLAM %s", r.Level+1, r.Modifier)
		ss = append(ss, s)
	}

	if r.Parts != nil {
		for _, part := range r.Parts {
			s = fmt.Sprintf("%s", part)
			ss = append(ss, s)
		}
	}

	if r.Citation != nil {
		s = fmt.Sprintf("%s", r.Citation)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = fmt.Sprintf("%s", r.Change)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r PlaceRecords) String() string {
	var ss []string
	var s string

	log.Printf("PlaceRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}
	if ss == nil {
		return ""
	}

	if ss == nil {
		return ""
	}
	return strings.Join(ss, "\n")
}

func (r *RootRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s", r.Header)
	ss = append(ss, s)

	if r.Submitter != nil {
		s = fmt.Sprintf("%s", r.Submitter)
		ss = append(ss, s)
	}

	if r.Submission != nil {
		s = fmt.Sprintf("%s", r.Submission)
		ss = append(ss, s)
	}

	if r.Place != nil {
		for _, x := range r.Place {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	if r.Event != nil {
		for _, x := range r.Event {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	for _, x := range r.Individual {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	for _, x := range r.Family {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}

	if r.Repository != nil {
		for _, x := range r.Repository {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	if r.Source != nil {
		for _, x := range r.Source {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	if r.Note != nil {
		for _, x := range r.Note {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	if r.Media != nil {
		for _, x := range r.Media {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	if r.ChildStatus != nil {
		for _, x := range r.ChildStatus {
			s = fmt.Sprintf("%s", x)
			ss = append(ss, s)
		}
	}

	s = fmt.Sprintf("%s", r.Trailer)
	ss = append(ss, s)

	return strings.Join(ss, "\n")
}

func (r *SchemaRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d SCHEMA", r.Level)
	ss = append(ss, s)

	return strings.Join(ss, "\n")
}

func (r *SourceRecord) String() string {
	var ss []string
	var s string

	if r.Name != "" {
		s = fmt.Sprintf("%d NAME %s", r.Level+1, r.Name)
		ss = append(ss, s)
	}

	if r.Title != "" {
		s = fmt.Sprintf("%d TITL %s", r.Level+1, r.Title)
		ss = append(ss, s)

	}

	if r.Text != nil {
		for _, text := range r.Text {
			s = fmt.Sprintf("%d TEXT %s", r.Level+1, text)
			ss = append(ss, s)
		}
	}

	if r.Data != nil {
		s = fmt.Sprintf("%s", r.Data)
		ss = append(ss, s)
	}

	if r.Media != nil {
		s = fmt.Sprintf("%s", r.Media)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = fmt.Sprintf("%s", r.Note)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r SourceRecords) String() string {
	var ss []string
	var s string
	log.Printf("SourceRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x)
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

func (r *SubmissionLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d SUBN @%s@", r.Level, r.Xref)
	ss = append(ss, s)

	// TODO
	log.Println("Submission Link Stringer TODO")

	return strings.Join(ss, "\n")
}

func (r *SubmissionRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d @%s@ SUBN", r.Level, r.Xref)
	ss = append(ss, s)

	// TODO
	log.Println("SubmissionRecord Stringer TODO")

	return strings.Join(ss, "\n")
}

func (r *SubmitterLink) String() string {

	var ss []string
	var s string

	s = fmt.Sprintf("%d SUBN @%s@", r.Level, r.Xref)
	ss = append(ss, s)

	// TODO
	log.Println("Submitter Link Stringer TODO")

	return strings.Join(ss, "\n")
}

func (r *SubmitterRecord) String() string {
	var ss []string
	var s, sXref0, sXrefn string

	if r.Level == 0 {
		sXrefn, sXref0 = "", fmt.Sprintf("@%s@ ", r.Xref)
	} else {
		sXref0, sXrefn = "", fmt.Sprintf(" @%s@", r.Xref)
	}
	s = fmt.Sprintf("%d %sSUBM%s", r.Level, sXref0, sXrefn)
	ss = append(ss, s)

	// TODO
	log.Println("SubmitterRecord Stringer TODO")

	return strings.Join(ss, "\n")
}

func (r *SystemRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%d SOUR %s", r.Level, r.Id)
	ss = append(ss, s)

	if r.ProductName != "" {
		s = fmt.Sprintf("%d NAME %s", r.Level+1, r.ProductName)
		ss = append(ss, s)
	}
	if r.Version != "" {
		s = fmt.Sprintf("%d VERS %s", r.Level+1, r.Version)
		ss = append(ss, s)
	}

	if r.Business != nil {
		s = fmt.Sprintf("%s", r.Business)
		ss = append(ss, s)
	}

	if r.SourceData != nil {
		s = fmt.Sprintf("%s", r.SourceData)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (r *TrailerRecord) String() string {
	return fmt.Sprintf("%d TRLR", 0)
}

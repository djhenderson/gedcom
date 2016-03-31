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

// indent emits spaces based on the level number
func indent(i int) string {
	const spaces = "                    "
	return spaces[:i*2]
}

// LongString formats a long string using CONT and CONC lines
func LongString(level int, xref string, tag string, longString string) []string {
	var s string
	var ss []string

	sXref0, sXrefN := "", ""
	if xref != "" {
		sXref0 = fmt.Sprintf(" @%s@", xref)
		if level != 0 {
			sXrefN, sXref0 = sXref0, sXrefN
		}
	}

	if longString == "" {
		s = fmt.Sprintf("%s%d%s %s%s", indent(level), level, sXref0, tag, sXrefN)
		ss = append(ss, s)
		return ss
	}

	parts := strings.Split(longString, "\n")
	for i, part := range parts {
		s = fmt.Sprintf("%s%d%s %s%s %s", indent(level), level, sXref0, tag, sXrefN, part)
		ss = append(ss, s)
		if i == 0 {
			tag = "CONT"
			level++
			sXref0, sXrefN = "", ""
		}
	}
	return ss
}

// String stringifies a GEDCOM address record
func (r *AddressRecord) String() string {
	var ss []string
	var s string

	sas := LongString(r.Level, "", "ADDR", r.Full)
	ss = append(ss, sas...)

	if r.Line1 != "" {
		s = fmt.Sprintf("%s%d ADR1 %s", indent(r.Level+1), r.Level+1, r.Line1)
		ss = append(ss, s)
	}

	if r.Line2 != "" {
		s = fmt.Sprintf("%s%d ADR2 %s", indent(r.Level+1), r.Level+1, r.Line2)
		ss = append(ss, s)
	}

	if r.Line3 != "" {
		s = fmt.Sprintf("%s%d ADR3 %s", indent(r.Level+1), r.Level+1, r.Line3)
		ss = append(ss, s)
	}

	if r.City != "" {
		s = fmt.Sprintf("%s%d CITY %s", indent(r.Level+1), r.Level+1, r.City)
		ss = append(ss, s)
	}

	if r.State != "" {
		s = fmt.Sprintf("%s%d STAE %s", indent(r.Level+1), r.Level+1, r.State)
		ss = append(ss, s)
	}

	if r.PostalCode != "" {
		s = fmt.Sprintf("%s%d POST %s", indent(r.Level+1), r.Level+1, r.PostalCode)
		ss = append(ss, s)
	}

	if r.Country != "" {
		s = fmt.Sprintf("%s%d CTRY %s", indent(r.Level+1), r.Level+1, r.Country)
		ss = append(ss, s)
	}

	if r.Phone != "" {
		s = fmt.Sprintf("%s%d PHON %s", indent(r.Level+1), r.Level+1, r.Phone)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of address records
func (r AddressRecords) String() string {
	var ss []string
	var s string

	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM bibliography record
func (r *BibliographyRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d BIBL %s", indent(r.Level), r.Level, r.Value)
	ss = append(ss, s)

	for _, comp := range r.Component {
		sas := LongString(r.Level+1, "", "COMP", comp)
		ss = append(ss, sas...)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM blob record
func (r *BlobRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d BLOB", indent(r.Level), r.Level)
	ss = append(ss, s)

	parts := strings.Split(r.Data, "\n")
	for _, data := range parts {
		s := fmt.Sprintf("%s%d CONT %s", indent(r.Level+1), r.Level+1, data)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM business record
func (r *BusinessRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d CORP %s", indent(r.Level), r.Level, r.BusinessName)
	ss = append(ss, s)

	if r.Address != nil {
		s = r.Address.String()
		ss = append(ss, s)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			s = fmt.Sprintf("%s%d PHON %s", indent(r.Level+1), r.Level+1, phone)
			ss = append(ss, s)
		}
	}

	if r.WebSite != "" {
		s = fmt.Sprintf("%s%d WWW %s", indent(r.Level), r.Level, r.WebSite)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM call number record
func (r *CallNumberRecord) String() string {
	var ss []string
	var s string
	s = fmt.Sprintf("%s%d CALN %s", indent(r.Level), r.Level, r.CallNumber)
	ss = append(ss, s)

	if r.Media != "" {
		s = fmt.Sprintf("%s%d MEDI %s", indent(r.Level+1), r.Level+1, r.Media)
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM change record
func (r *ChangeRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d CHAN", indent(r.Level), r.Level)
	ss = append(ss, s)

	if r.Date != nil {
		s = r.Date.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM character set record
func (r *CharacterSetRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d CHAR %s", indent(r.Level), r.Level, r.CharacterSet)
	ss = append(ss, s)

	if r.Version != "" {
		s = fmt.Sprintf("%s%d VERS %s", indent(r.Level+1), r.Level+1, r.Version)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM child status record
func (r *ChildStatusRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d @%s@ CSTA", indent(r.Level), r.Level, r.Xref)
	ss = append(ss, s)

	s = fmt.Sprintf("%s%d NAME %s", indent(r.Level+1), r.Level+1, r.Name)
	ss = append(ss, s)

	return strings.Join(ss, "\n")
}

// String stringifies a slice of child status records
func (r ChildStatusRecords) String() string {
	var ss []string
	var s string

	// log.Printf("ChildStatusRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM citatio record
func (r *CitationRecord) String() string {
	var ss []string
	var s string

	xref := ""
	if r.Source != nil {
		xref = r.Source.Xref
	}
	sas := LongString(r.Level, xref, "SOUR", r.Value)
	ss = append(ss, sas...)

	if r.Page != "" {
		s = fmt.Sprintf("%s%d PAGE %s", indent(r.Level+1), r.Level+1, r.Page)
		ss = append(ss, s)
	}

	if r.Reference != "" {
		s = fmt.Sprintf("%s%d REF %s", indent(r.Level+1), r.Level+1, r.Reference)
		ss = append(ss, s)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			s = event.String()
			ss = append(ss, s)
		}
	}

	if r.Data != nil {
		for _, data := range r.Data {
			s = data.String()
			ss = append(ss, s)
		}
	}

	if r.Text != nil {
		for _, text := range r.Text {
			sas := LongString(r.Level+1, "", "TEXT", text)
			ss = append(ss, sas...)
		}
	}

	if r.Quality != "" {
		s = fmt.Sprintf("%s%d QUAY %s", indent(r.Level+1), r.Level+1, r.Quality)
		ss = append(ss, s)
	}

	if r.Media != nil {
		for _, media := range r.Media {
			s = media.String()
			ss = append(ss, s)
		}
	}

	if r.CONS != "" {
		s = fmt.Sprintf("%s%d CONS %s", indent(r.Level+1), r.Level+1, r.CONS)
		ss = append(ss, s)
	}

	if r.Direct != "" {
		s = fmt.Sprintf("%s%d DIRE %s", indent(r.Level+1), r.Level+1, r.Direct)
		ss = append(ss, s)
	}

	if r.SourceQuality != "" {
		s = fmt.Sprintf("%s%d SOQU %s", indent(r.Level+1), r.Level+1, r.SourceQuality)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of citation records
func (r CitationRecords) String() string {
	var ss []string
	var s string

	//log.Printf("CitationRecords type(r): %T\n", r)
	for _, note := range r {
		//log.Printf("CitationRecords type(note): %T\n", note)
		//log.Printf("CitationRecords type(*note): %T\n", *note)
		s = note.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM data record
func (r *DataRecord) String() string {
	var ss []string
	var s string

	//s = fmt.Sprintf("%s%d DATA %s", indent(r.Level), r.Level, r.Data)
	//ss = append(ss, s)
	sas := LongString(r.Level, "", "DATA", r.Data)
	ss = append(ss, sas...)

	if r.Date != "" {
		s = fmt.Sprintf("%s%d DATE %s", indent(r.Level+1), r.Level+1, r.Date)
		ss = append(ss, s)
	}

	if r.Copyright != "" {
		s = fmt.Sprintf("%s%d COPR %s", indent(r.Level+1), r.Level+1, r.Copyright)
		ss = append(ss, s)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			//s = fmt.Sprintf("%s%d TEXT %s", indent(r.Level+1), r.Level+1, text)
			//ss = append(ss, s)
			sas := LongString(r.Level+1, "", "TEXT", text)
			ss = append(ss, sas...)
		}
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of data records
func (r DataRecords) String() string {
	var ss []string
	var s string

	//log.Printf("DataRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM date record
func (r *DateRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d DATE %s", indent(r.Level), r.Level, r.Date)
	ss = append(ss, s)

	if r.Day != "" {
		s = fmt.Sprintf("%s%d DATD %s", indent(r.Level+1), r.Level+1, r.Day)
		ss = append(ss, s)
	}

	if r.Month != "" {
		s = fmt.Sprintf("%s%d DATM %s", indent(r.Level+1), r.Level+1, r.Month)
		ss = append(ss, s)
	}

	if r.Year != "" {
		s = fmt.Sprintf("%s%d DATY %s", indent(r.Level+1), r.Level+1, r.Year)
		ss = append(ss, s)
	}

	if r.Full != "" {
		s = fmt.Sprintf("%s%d DATF %s", indent(r.Level+1), r.Level+1, r.Full)
		ss = append(ss, s)
	}

	if r.Short != "" {
		s = fmt.Sprintf("%s%d DATS %s", indent(r.Level+1), r.Level+1, r.Short)
		ss = append(ss, s)
	}

	if r.Time != "" {
		s = fmt.Sprintf("%s%d TIME %s", indent(r.Level+1), r.Level+1, r.Time)
		ss = append(ss, s)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			//s = fmt.Sprintf("%s%d TEXT %s", indent(r.Level+1), r.Level+1, text)
			//ss = append(ss, s)
			sas := LongString(r.Level+1, "", "TEXT", text)
			ss = append(ss, sas...)
		}
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM event definition record
func (r EventDefinitionRecord) String() string {
	log.Fatal("EventDefinition Stringer not implemented")
	return ""
}

// String stringifies a slice of event definition records
func (r EventDefinitionRecords) String() string {
	var ss []string
	var s string

	//log.Printf("EventDefinitionRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM event records
func (r *EventRecord) String() string {
	var ss []string
	var s string

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}
	s = fmt.Sprintf("%s%d %s%s %s", indent(r.Level), r.Level, id, r.Tag, r.Value)
	ss = append(ss, s)

	if r.Type != "" {
		s = fmt.Sprintf("%s%d TYPE %s", indent(r.Level+1), r.Level+1, r.Type)
		ss = append(ss, s)
	}

	if r.Name != "" {
		s = fmt.Sprintf("%s%d NAME %s", indent(r.Level+1), r.Level+1, r.Name)
		ss = append(ss, s)
	}

	if r.Primary_ != "" {
		s = fmt.Sprintf("%s%d _PRIM %s", indent(r.Level+1), r.Level+1, r.Primary_)
		ss = append(ss, s)
	}

	if r.Date != nil {
		s = r.Date.String()
		ss = append(ss, s)
	}

	if r.Place != nil {
		s = r.Place.String()
		ss = append(ss, s)
	}

	if r.Role != nil {
		s = r.Role.String()
		ss = append(ss, s)
	}
	if r.Address != nil {
		s = r.Address.String()
		ss = append(ss, s)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			s = fmt.Sprintf("%s%d PHON %s", indent(r.Level+1), r.Level+1, phone)
			ss = append(ss, s)
		}
	}

	if r.Parents != nil {
		s = r.Parents.String()
		ss = append(ss, s)
	}

	if r.Husband != nil {
		s = r.Husband.String()
		ss = append(ss, s)
	}

	if r.Wife != nil {
		s = r.Wife.String()
		ss = append(ss, s)
	}

	if r.Spouse != nil {
		s = r.Spouse.String()
		ss = append(ss, s)
	}

	if r.Age != "" {
		s = fmt.Sprintf("%s%d AGE %s", indent(r.Level+1), r.Level+1, r.Age)
		ss = append(ss, s)
	}

	if r.Agency != "" {
		s = fmt.Sprintf("%s%d AGNC %s", indent(r.Level+1), r.Level+1, r.Agency)
		ss = append(ss, s)
	}

	if r.Cause != "" {
		s = fmt.Sprintf("%s%d CAUS %s", indent(r.Level+1), r.Level+1, r.Cause)
		ss = append(ss, s)
	}

	if r.Temple != "" {
		s = fmt.Sprintf("%s%d TEMP %s", indent(r.Level+1), r.Level+1, r.Temple)
		ss = append(ss, s)
	}

	if r.Status != "" {
		s = fmt.Sprintf("%s%d STAT %s", indent(r.Level+1), r.Level+1, r.Status)
		ss = append(ss, s)
	}

	if r.Media != nil {
		for _, media := range r.Media {
			s = media.String()
			if s != "" {
				ss = append(ss, s)
			}
		}
	}

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	if r.UID_ != nil {
		for _, uid := range r.UID_ {
			s = fmt.Sprintf("%s%d _UID %s", indent(r.Level+1), r.Level+1, uid)
			ss = append(ss, s)
		}
	}

	if r.UpdateTime_ != "" {
		s = fmt.Sprintf("%s%d _UPD %s", indent(r.Level+1), r.Level+1, r.UpdateTime_)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of event records
func (r EventRecords) String() string {
	var ss []string
	var s string

	//log.Printf("EventRecords type(r): %T\n", r)
	for _, x := range r {
		s = fmt.Sprintf("%s", x.String())
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM link to a family record
func (r *FamilyLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d %s @%s@", indent(r.Level), r.Level, r.Tag, r.Family.Xref)
	ss = append(ss, s)

	if r.Pedigree != "" {
		s = fmt.Sprintf("%s%d PEDI %s", indent(r.Level+1), r.Level+1, r.Pedigree)
		ss = append(ss, s)
	}

	if r.Adopted != "" {
		s = fmt.Sprintf("%s%d ADOP %s", indent(r.Level+1), r.Level+1, r.Adopted)
		ss = append(ss, s)
	}

	if r.Primary_ != "" {
		s = fmt.Sprintf("%s%d _PRIMARY %s", indent(r.Level+1), r.Level+1, r.Primary_)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of links to family records
func (r FamilyLinks) String() string {
	var ss []string
	var s string

	//log.Printf("FamilyLinks type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM family record
func (r *FamilyRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d @%s@ FAM", indent(r.Level), r.Level, r.Xref)
	ss = append(ss, s)

	if r.Husband != nil {
		s = r.Husband.String()
		ss = append(ss, s)
	}

	if r.Wife != nil {
		s = r.Wife.String()
		ss = append(ss, s)
	}

	if r.NumChildren > 0 {
		s = fmt.Sprintf("%s%d NCHI %d", indent(r.Level+1), r.Level+1, r.NumChildren)
		ss = append(ss, s)
	}

	if r.Child != nil {
		for _, child := range r.Child {
			s = child.String()
			ss = append(ss, s)
		}
	}

	if r.Event != nil {
		for _, event := range r.Event {
			s = event.String()
			ss = append(ss, s)
		}
	}

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			s = subm.String()
			ss = append(ss, s)
		}
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	if r.UID_ != nil {
		for _, uid := range r.UID_ {
			s = fmt.Sprintf("%s%d _UID %s", indent(r.Level+1), r.Level+1, uid)
			ss = append(ss, s)
		}
	}

	if r.UpdateTime_ != "" {
		s = fmt.Sprintf("%s%d _UPD %s", indent(r.Level+1), r.Level+1, r.UpdateTime_)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of family records
func (r FamilyRecords) String() string {
	var ss []string
	var s string

	//log.Printf("FamilyRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM footnote record
func (r *FootnoteRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d FOOT %s", indent(r.Level), r.Level, r.Value)
	ss = append(ss, s)

	for _, comp := range r.Component {
		sas := LongString(r.Level+1, "", "COMP", comp)
		ss = append(ss, sas...)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM gedcom record
func (r *GedcomRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d GEDC", indent(r.Level), r.Level)
	ss = append(ss, s)

	s = fmt.Sprintf("%s%d VERS %s", indent(r.Level+1), r.Level+1, r.Version)
	ss = append(ss, s)

	s = fmt.Sprintf("%s%d FORM %s", indent(r.Level+1), r.Level+1, r.Form)
	ss = append(ss, s)

	return strings.Join(ss, "\n")

}

// String stringifies a GEDCOM header record
func (r *HeaderRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d HEAD", indent(r.Level), r.Level)
	ss = append(ss, s)

	if r.CharacterSet != nil {
		s = r.CharacterSet.String()
		ss = append(ss, s)
	}

	if r.SourceSystem != nil {
		s = r.SourceSystem.String()
		ss = append(ss, s)
	}

	if r.Destination != "" {
		s = fmt.Sprintf("%s%d DEST %s", indent(r.Level+1), r.Level+1, r.Destination)
		ss = append(ss, s)
	}

	if r.Date != nil {
		s = r.Date.String()
		ss = append(ss, s)
	}

	if r.Time != "" {
		s = fmt.Sprintf("%s%d TIME %s", indent(r.Level+1), r.Level+1, r.Time)
		ss = append(ss, s)
	}

	if r.FileName != "" {
		s = fmt.Sprintf("%s%d FILE %s", indent(r.Level+1), r.Level+1, r.FileName)
		ss = append(ss, s)
	}

	if r.Gedcom != nil {
		s = r.Gedcom.String()
		ss = append(ss, s)
	}

	if r.Language != "" {
		s = fmt.Sprintf("%s%d LANG %s", indent(r.Level+1), r.Level+1, r.Language)
		ss = append(ss, s)
	}

	if r.Copyright != "" {
		s = fmt.Sprintf("%s%d COPR %s", indent(r.Level+1), r.Level+1, r.Copyright)
		ss = append(ss, s)
	}

	if r.Place != nil {
		s = r.Place.String()
		ss = append(ss, s)
	}

	if r.RootPerson_ != nil {
		s = r.RootPerson_.String()
		ss = append(ss, s)
	}

	if r.HomePerson_ != nil {
		s = r.HomePerson_.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			s = subm.String()
			ss = append(ss, s)
		}
	}

	if r.Submission != nil {
		for _, subn := range r.Submission {
			s = subn.String()
			ss = append(ss, s)
		}
	}

	if r.Schema != nil {
		s = r.Schema.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM history record
func (r *HistoryRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d HIST %s", indent(r.Level), r.Level, r.History)
	ss = append(ss, s)

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of history records
func (r HistoryRecords) String() string {
	var ss []string
	var s string

	for _, x := range r {
		s = fmt.Sprintf("%s", x.String())
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM individual links
func (r *IndividualLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d %s @%s@", indent(r.Level), r.Level, r.Tag, r.Individual.Xref)
	ss = append(ss, s)

	if r.Relationship != "" {
		s = fmt.Sprintf("%s%d RELA %s", indent(r.Level+1), r.Level+1, r.Relationship)
		ss = append(ss, s)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			s = event.String()
			ss = append(ss, s)
		}
	}

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of links to individual records
func (r IndividualLinks) String() string {
	var ss []string
	var s string

	//log.Printf("IndividualLinks type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM individual record
func (r *IndividualRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d @%s@ INDI", indent(r.Level), r.Level, r.Xref)
	ss = append(ss, s)

	for _, name := range r.Name {
		s = name.String()
		ss = append(ss, s)
	}

	if r.Restriction != "" {
		s = fmt.Sprintf("%s%d RESN %s", indent(r.Level+1), r.Level+1, r.Restriction)
		ss = append(ss, s)
	}

	if r.Sex != "" {
		s = fmt.Sprintf("%s%d SEX %s", indent(r.Level+1), r.Level+1, r.Sex)
		ss = append(ss, s)
	}

	if r.ProfilePicture_ != nil {
		s = r.ProfilePicture_.String()
		ss = append(ss, s)
	}

	if r.CONL != "" {
		s = fmt.Sprintf("%s%d CONL %s", indent(r.Level+1), r.Level+1, r.CONL)
		ss = append(ss, s)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			s = event.String()
			ss = append(ss, s)
		}
	}

	if r.Attribute != "" {
		s = fmt.Sprintf("%s%d ATTR %s", indent(r.Level+1), r.Level+1, r.Attribute)
		ss = append(ss, s)
	}

	if r.Parents != nil {
		for _, famc := range r.Parents {
			s = famc.String()
			ss = append(ss, s)
		}
	}

	if r.Family != nil {
		for _, fams := range r.Family {
			s = fams.String()
			ss = append(ss, s)
		}
	}

	if r.Address != nil {
		s = r.Address.String()
		ss = append(ss, s)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			s = fmt.Sprintf("%s%d PHON %s", indent(r.Level+1), r.Level+1, phone)
			ss = append(ss, s)
		}
	}

	if r.Health != "" {
		s = fmt.Sprintf("%s%d HEAL %s", indent(r.Level+1), r.Level+1, r.Health)
		ss = append(ss, s)
	}

	if r.History != nil {
		s = r.History.String()
		ss = append(ss, s)
	}

	if r.Quality != "" {
		s = fmt.Sprintf("%s%d QUAY %s", indent(r.Level+1), r.Level+1, r.Quality)
		ss = append(ss, s)
	}

	if r.Living != "" {
		s = fmt.Sprintf("%s%d LVG %s", indent(r.Level+1), r.Level+1, r.Living)
		ss = append(ss, s)
	}

	if r.AncestralFileNumber != nil {
		for _, afn := range r.AncestralFileNumber {
			s = fmt.Sprintf("%s%d AFN %s", indent(r.Level+1), r.Level+1, afn)
			ss = append(ss, s)
		}
	}

	if r.RecordFileNumber != "" {
		s = fmt.Sprintf("%s%d RFN %s", indent(r.Level+1), r.Level+1, r.RecordFileNumber)
		ss = append(ss, s)
	}

	if r.UserReferenceNumber != nil {
		s = r.UserReferenceNumber.String()
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	if r.UID_ != nil {
		for _, uid := range r.UID_ {
			s = fmt.Sprintf("%s%d _UID %s", indent(r.Level+1), r.Level+1, uid)
			ss = append(ss, s)
		}
	}

	if r.Email != "" {
		s = fmt.Sprintf("%s%d EMAIL %s", indent(r.Level+1), r.Level+1, r.Email)
		ss = append(ss, s)
	}

	if r.WebSite != "" {
		s = fmt.Sprintf("%s%d WWW %s", indent(r.Level+1), r.Level+1, r.WebSite)
		ss = append(ss, s)
	}

	if r.Media != nil {
		s = r.Media.String()
		ss = append(ss, s)
	}

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Associated != nil {
		s = r.Associated.String()
		ss = append(ss, s)
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			s = subm.String()
			ss = append(ss, s)
		}
	}

	if r.ANCI != nil {
		for _, subm := range r.ANCI {
			s = subm.String()
			ss = append(ss, s)
		}
	}

	if r.DESI != nil {
		for _, subm := range r.DESI {
			s = subm.String()
			ss = append(ss, s)
		}
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	if r.Alias != "" {
		s = fmt.Sprintf("%s%d ALIA %s", indent(r.Level+1), r.Level+1, r.Alias)
		ss = append(ss, s)
	}

	if r.Father != nil {
		s = r.Father.String()
		ss = append(ss, s)
	}

	if r.Mother != nil {
		s = r.Mother.String()
		ss = append(ss, s)
	}

	if r.Miscellaneous != nil {
		for _, misc := range r.Miscellaneous {
			s = fmt.Sprintf("%s%d MISC %s", indent(r.Level+1), r.Level+1, misc)
			ss = append(ss, s)
		}
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of individual records
func (r IndividualRecords) String() string {
	var ss []string
	var s string

	//log.Printf("IndividualRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM media links
func (r *MediaLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d %s @%s@", indent(r.Level), r.Level, r.Tag, r.Media.Xref)
	ss = append(ss, s)

	return strings.Join(ss, "\n")
}

// String stringifies a slice of links to media records
func (r MediaLinks) String() string {
	var ss []string
	var s string

	//log.Printf("MediaLinks type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM media record
func (r *MediaRecord) String() string {
	var ss []string
	var s string

	id0, idN := "", ""
	if r.Xref != "" {
		id0 = fmt.Sprintf(" @%s@ ", r.Xref)

		if r.Level != 0 {
			idN, id0 = id0, idN
		}
	}

	s = fmt.Sprintf("%s%d%s %s%s", indent(r.Level), r.Level, id0, "OBJE", idN)
	ss = append(ss, s)

	if r.Format != "" {
		s = fmt.Sprintf("%s%d FORM %s", indent(r.Level+1), r.Level+1, r.Format)
		ss = append(ss, s)
	}

	if r.URL_ != "" {
		s = fmt.Sprintf("%s%d _URL %s", indent(r.Level+1), r.Level+1, r.URL_)
		ss = append(ss, s)
	}

	if r.FileName != "" {
		s = fmt.Sprintf("%s%d FILE %s", indent(r.Level+1), r.Level+1, r.FileName)
		ss = append(ss, s)
	}

	if r.Title != "" {
		s = fmt.Sprintf("%s%d TITL %s", indent(r.Level+1), r.Level+1, r.Title)
		ss = append(ss, s)
	}

	if r.Date != "" {
		s = fmt.Sprintf("%s%d DATE %s", indent(r.Level+1), r.Level+1, r.Date)
		ss = append(ss, s)
	}

	if r.Author != "" {
		s = fmt.Sprintf("%s%d AUTH %s", indent(r.Level+1), r.Level+1, r.Author)
		ss = append(ss, s)
	}

	if r.Text != "" {
		sas := LongString(r.Level, "", "TEXT", r.Text)
		ss = append(ss, sas...)
	}

	if r.Date_ != "" {
		s = fmt.Sprintf("%s%d _DATE %s", indent(r.Level+1), r.Level+1, r.Date_)
		ss = append(ss, s)
	}

	if r.AstId_ != "" {
		s = fmt.Sprintf("%s%d _ASTID %s", indent(r.Level+1), r.Level+1, r.AstId_)
		ss = append(ss, s)
	}

	if r.AstType_ != "" {
		s = fmt.Sprintf("%s%d _ASTTYP %s", indent(r.Level+1), r.Level+1, r.AstType_)
		ss = append(ss, s)
	}

	if r.AstDesc_ != "" {
		s = fmt.Sprintf("%s%d _ASTDESC %s", indent(r.Level+1), r.Level+1, r.AstDesc_)
		ss = append(ss, s)
	}

	if r.AstLoc_ != "" {
		s = fmt.Sprintf("%s%d _ASTLOC %s", indent(r.Level+1), r.Level+1, r.AstLoc_)
		ss = append(ss, s)
	}

	if r.AstPerm_ != "" {
		s = fmt.Sprintf("%s%d _ASTPERM %s", indent(r.Level+1), r.Level+1, r.AstPerm_)
		ss = append(ss, s)
	}

	if r.AstUpPid_ != "" {
		s = fmt.Sprintf("%s%d _ASTUPPID %s", indent(r.Level+1), r.Level+1, r.AstUpPid_)
		ss = append(ss, s)
	}

	if r.BinaryLargeObject != nil {
		s = r.BinaryLargeObject.String()
		ss = append(ss, s)
	}

	if r.UserReferenceNumber != nil {
		s = r.UserReferenceNumber.String()
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of media records
func (r MediaRecords) String() string {
	var ss []string
	var s string

	//log.Printf("MediaRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM name record
func (r *NameRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d NAME %s", indent(r.Level), r.Level, r.Name)
	ss = append(ss, s)

	if r.Prefix != "" {
		s = fmt.Sprintf("%s%d NPFX %s", indent(r.Level+1), r.Level+1, r.Prefix)
		ss = append(ss, s)
	}

	if r.GivenName != "" {
		s = fmt.Sprintf("%s%d GIVN %s", indent(r.Level+1), r.Level+1, r.GivenName)
		ss = append(ss, s)
	}

	if r.MiddleName_ != "" {
		s = fmt.Sprintf("%s%d _MIDN %s", indent(r.Level+1), r.Level+1, r.MiddleName_)
		ss = append(ss, s)
	}

	if r.SurnamePrefix != "" {
		s = fmt.Sprintf("%s%d SPFX %s", indent(r.Level+1), r.Level+1, r.SurnamePrefix)
		ss = append(ss, s)
	}

	if r.Surname != "" {
		s = fmt.Sprintf("%s%d SURN %s", indent(r.Level+1), r.Level+1, r.Surname)
		ss = append(ss, s)
	}

	if r.Suffix != "" {
		s = fmt.Sprintf("%s%d NSFX %s", indent(r.Level+1), r.Level+1, r.Suffix)
		ss = append(ss, s)
	}

	if r.PreferedGivenName_ != "" {
		s = fmt.Sprintf("%s%d _PGVN %s", indent(r.Level+1), r.Level+1, r.PreferedGivenName_)
		ss = append(ss, s)
	}

	if r.NameType != "" {
		s = fmt.Sprintf("%s%d _PRIM %s", indent(r.Level+1), r.Level+1, r.Primary_)
		ss = append(ss, s)
	}

	if r.Primary_ != "" {
		s = fmt.Sprintf("%s%d _PRIM %s", indent(r.Level+1), r.Level+1, r.Primary_)
		ss = append(ss, s)
	}

	if r.AKA_ != nil {
		for _, aka := range r.AKA_ {
			s = fmt.Sprintf("%s%d _AKA %s", indent(r.Level+1), r.Level+1, aka)
			ss = append(ss, s)
		}
	}

	if r.Nickname != nil {
		for _, nick := range r.Nickname {
			s = fmt.Sprintf("%s%d NICK %s", indent(r.Level+1), r.Level+1, nick)
			ss = append(ss, s)
		}
	}

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		for _, note := range r.Note {
			s = note.String()
			ss = append(ss, s)
		}
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of name records
func (r NameRecords) String() string {
	var ss []string
	var s string

	//log.Printf("NameRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM note records
func (r *NoteRecord) String() string {
	var ss []string
	var s string

	sas := LongString(r.Level, r.Xref, "NOTE", r.Note)
	ss = append(ss, sas...)

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.UserReferenceNumber != nil {
		s = r.UserReferenceNumber.String()
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of note records
func (r NoteRecords) String() string {
	var ss []string
	var s string

	//log.Printf("NoteRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM place part record
func (r *PlacePartRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d %s %s", indent(r.Level), r.Level, r.Tag, r.Part)
	ss = append(ss, s)

	if r.Jurisdiction != "" {
		s = fmt.Sprintf("%s%d JURI %s", indent(r.Level+1), r.Level+1, r.Jurisdiction)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM place record
func (r *PlaceRecord) String() string {
	var ss []string
	var s string

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}

	s = fmt.Sprintf("%s%d %sPLAC %s", indent(r.Level), r.Level, id, r.Name)
	ss = append(ss, s)

	if r.Form != "" {
		s = fmt.Sprintf("%s%d PLAS %s", indent(r.Level+1), r.Level+1, r.Form)
		ss = append(ss, s)
	}

	if r.ShortName != "" {
		s = fmt.Sprintf("%s%d PLAS %s", indent(r.Level+1), r.Level+1, r.ShortName)
		ss = append(ss, s)
	}

	if r.Modifier != "" {
		s = fmt.Sprintf("%s%d PLAM %s", indent(r.Level+1), r.Level+1, r.Modifier)
		ss = append(ss, s)
	}

	if r.Parts != nil {
		for _, part := range r.Parts {
			s = part.String()
			ss = append(ss, s)
		}
	}

	if r.Citation != nil {
		s = r.Citation.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of place records
func (r PlaceRecords) String() string {
	var ss []string
	var s string

	//log.Printf("PlaceRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM link to a repository record
func (r *RepositoryLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d REPO @%s@", indent(r.Level), r.Level, r.Repository.Xref)
	ss = append(ss, s)

	if r.CallNumber != nil {
		s = r.CallNumber.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM repository record
func (r *RepositoryRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d @%s@ REPO", indent(r.Level), r.Level, r.Xref)
	ss = append(ss, s)

	if r.Name != "" {
		s = fmt.Sprintf("%s%d NAME %s", indent(r.Level+1), r.Level+1, r.Name)
		ss = append(ss, s)
	}

	if r.Address != nil {
		s = r.Address.String()
		ss = append(ss, s)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			s = fmt.Sprintf("%s%d PHON %s", indent(r.Level+1), r.Level+1, phone)
			ss = append(ss, s)
		}
	}

	if r.WebSite != "" {
		s = fmt.Sprintf("%s%d WWW %s", indent(r.Level+1), r.Level+1, r.WebSite)
		ss = append(ss, s)
	}

	if r.UserReferenceNumber != nil {
		s = r.UserReferenceNumber.String()
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM role record
func (r *RoleRecord) String() string {
	var ss []string
	var s string

	spacer := " "
	if r.Role == "" {
		spacer = ""
	}
	s = fmt.Sprintf("%s%d ROLE %s%s@%s@", indent(r.Level), r.Level, r.Role, spacer, r.Individual.Xref)
	ss = append(ss, s)

	if r.Principal != "" {
		s = fmt.Sprintf("%s%d PRIN %s", indent(r.Level+1), r.Level+1, r.Principal)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of role record
func (r RoleRecords) String() string {
	var ss []string
	var s string

	//log.Printf("RoleRecords type(r): %T\n", r)
	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM root record, i.e. the whole file
func (r *RootRecord) String() string {
	var ss []string
	var s string

	if r.Header != nil {
		s = r.Header.String()
		ss = append(ss, s)
	}

	if r.Submission != nil {
		for _, subn := range r.Submission {
			s = subn.String()
			ss = append(ss, s)
		}
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			s = subm.String()
			ss = append(ss, s)
		}
	}

	if r.Place != nil {
		for _, x := range r.Place {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Event != nil {
		for _, x := range r.Event {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Individual != nil {
		for _, x := range r.Individual {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Family != nil {
		for _, x := range r.Family {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Repository != nil {
		for _, x := range r.Repository {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Source != nil {
		for _, x := range r.Source {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Note != nil {
		for _, x := range r.Note {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Media != nil {
		for _, x := range r.Media {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.ChildStatus != nil {
		for _, x := range r.ChildStatus {
			s = x.String()
			ss = append(ss, s)
		}
	}

	if r.Trailer != nil {
		s = r.Trailer.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM schema record
func (r *SchemaRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d SCHEMA", indent(r.Level), r.Level)
	ss = append(ss, s)

	for _, data := range r.Data {
		level := int(data[0])
		s = indent(level) + data
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM short title record
func (r *ShortTitleRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d SHTI %s", indent(r.Level), r.Level, r.ShortTitle)
	ss = append(ss, s)

	if r.Indexed != "" {
		s = fmt.Sprintf("%s%d INDX %s", indent(r.Level+1), r.Level+1, r.Indexed)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM level 0 source record
func (r *SourceRecord) String() string {
	var ss []string
	var s string

	sas := LongString(r.Level, r.Xref, "SOUR", r.Value)
	ss = append(ss, sas...)

	if r.Name != "" {
		s = fmt.Sprintf("%s%d NAME %s", indent(r.Level+1), r.Level+1, r.Name)
		ss = append(ss, s)
	}

	if r.Title != "" {
		//s = fmt.Sprintf("%s%d TITL %s", indent(r.Level+1), r.Level+1, r.Title)
		//ss = append(ss, s)
		sas := LongString(r.Level+1, "", "TITL", r.Title)
		ss = append(ss, sas...)
	}

	if r.Author != "" {
		s = fmt.Sprintf("%s%d AUTH %s", indent(r.Level+1), r.Level+1, r.Author)
		ss = append(ss, s)
	}

	if r.Abbreviation != "" {
		s = fmt.Sprintf("%s%d ABBR %s", indent(r.Level+1), r.Level+1, r.Abbreviation)
		ss = append(ss, s)
	}

	if r.Publication != "" {
		s = fmt.Sprintf("%s%d PUBL %s", indent(r.Level+1), r.Level+1, r.Publication)
		ss = append(ss, s)
	}

	if r.Parenthesized_ != "" {
		s = fmt.Sprintf("%s%d _PAREN %s", indent(r.Level+1), r.Level+1, r.Parenthesized_)
		ss = append(ss, s)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			//s = fmt.Sprintf("%s%d TEXT %s", indent(r.Level+1), r.Level+1, text)
			//ss = append(ss, s)
			sas := LongString(r.Level+1, "", "TEXT", text)
			ss = append(ss, sas...)
		}
	}

	if r.Data != nil {
		s = r.Data.String()
		ss = append(ss, s)
	}

	if r.Footnote != nil {
		s = r.Footnote.String()
		ss = append(ss, s)
	}

	if r.Bibliography != nil {
		s = r.Bibliography.String()
		ss = append(ss, s)
	}

	if r.Repository != nil {
		s = r.Repository.String()
		ss = append(ss, s)
	}

	if r.UserReferenceNumber != nil {
		s = r.UserReferenceNumber.String()
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	if r.ShortAuthor != "" {
		s = fmt.Sprintf("%s%d SHAU %s", indent(r.Level+1), r.Level+1, r.ShortAuthor)
		ss = append(ss, s)
	}

	if r.ShortTitle != nil {
		s = r.ShortTitle.String()
		ss = append(ss, s)
	}

	if r.Media != nil {
		s = r.Media.String()
		ss = append(ss, s)
	}

	if r.Note != nil {
		s = r.Note.String()
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
	}
	return strings.Join(ss, "\n")
}

// String stringifies a slice of source records
func (r SourceRecords) String() string {
	var ss []string
	var s string

	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM link to a submission record
func (r *SubmissionLink) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d SUBN @%s@", indent(r.Level), r.Level, r.Submission.Xref)
	ss = append(ss, s)

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM level 0 submission record
func (r *SubmissionRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d @%s@ SUBN", indent(r.Level), r.Level, r.Xref)
	ss = append(ss, s)

	if r.FamilyFileName != "" {
		s = fmt.Sprintf("%s%d FAMF %s", indent(r.Level+1), r.Level+1, r.FamilyFileName)
		ss = append(ss, s)
	}

	if r.Temple != "" {
		s = fmt.Sprintf("%s%d TEMP %s", indent(r.Level+1), r.Level+1, r.Temple)
		ss = append(ss, s)
	}

	if r.Ancestors != "" {
		s = fmt.Sprintf("%s%d ANCE %s", indent(r.Level+1), r.Level+1, r.Ancestors)
		ss = append(ss, s)
	}

	if r.Descendents != "" {
		s = fmt.Sprintf("%s%d DESC %s", indent(r.Level+1), r.Level+1, r.Descendents)
		ss = append(ss, s)
	}

	if r.Ordinance != "" {
		s = fmt.Sprintf("%s%d ORDI %s", indent(r.Level+1), r.Level+1, r.Ordinance)
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM link to a submitter record
func (r *SubmitterLink) String() string {

	var ss []string
	var s string

	s = fmt.Sprintf("%s%d %s @%s@", indent(r.Level), r.Level, r.Tag, r.Submitter.Xref)
	ss = append(ss, s)

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM level 0 submitter record
func (r *SubmitterRecord) String() string {
	var ss []string
	var s, sXref0, sXrefn string

	if r.Level == 0 {
		sXrefn, sXref0 = "", fmt.Sprintf("@%s@ ", r.Xref)
	} else {
		sXref0, sXrefn = "", fmt.Sprintf(" @%s@", r.Xref)
	}
	s = fmt.Sprintf("%s%d %sSUBM%s", indent(r.Level), r.Level, sXref0, sXrefn)
	ss = append(ss, s)

	if r.Name != "" {
		s = fmt.Sprintf("%s%d NAME %s", indent(r.Level+1), r.Level+1, r.Name)
		ss = append(ss, s)
	}

	if r.Address != nil {
		s = r.Address.String()
		ss = append(ss, s)
	}

	if r.Country != "" {
		s = fmt.Sprintf("%s%d CTRY %s", indent(r.Level+1), r.Level+1, r.Country)
		ss = append(ss, s)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			s = fmt.Sprintf("%s%d PHON %s", indent(r.Level+1), r.Level+1, phone)
			ss = append(ss, s)
		}
	}

	if r.Email != "" {
		s = fmt.Sprintf("%s%d EMAIL %s", indent(r.Level+1), r.Level+1, r.Email)
		ss = append(ss, s)
	}

	if r.WebSite != "" {
		s = fmt.Sprintf("%s%d WWW %s", indent(r.Level+1), r.Level+1, r.WebSite)
		ss = append(ss, s)
	}

	if r.Language != "" {
		s = fmt.Sprintf("%s%d LANG %s", indent(r.Level+1), r.Level+1, r.Language)
		ss = append(ss, s)
	}

	if r.Media != nil {
		s = r.Media.String()
		ss = append(ss, s)
	}

	if r.RecordFileNumber != "" {
		s = fmt.Sprintf("%s%d RFN %s", indent(r.Level+1), r.Level+1, r.RecordFileNumber)
		ss = append(ss, s)
	}

	if r.STAL != "" {
		s = fmt.Sprintf("%s%d STAL %s", indent(r.Level+1), r.Level+1, r.STAL)
		ss = append(ss, s)
	}

	if r.NUMB != "" {
		s = fmt.Sprintf("%s%d STAL %s", indent(r.Level+1), r.Level+1, r.NUMB)
		ss = append(ss, s)
	}

	if r.RIN != "" {
		s = fmt.Sprintf("%s%d RIN %s", indent(r.Level+1), r.Level+1, r.RIN)
		ss = append(ss, s)
	}

	if r.Change != nil {
		s = r.Change.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a slice of submitter records
func (r SubmitterRecords) String() string {
	var ss []string
	var s string

	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM system record
func (r *SystemRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d SOUR %s", indent(r.Level), r.Level, r.SystemName)
	ss = append(ss, s)

	if r.ProductName != "" {
		s = fmt.Sprintf("%s%d NAME %s", indent(r.Level+1), r.Level+1, r.ProductName)
		ss = append(ss, s)
	}

	if r.Version != "" {
		s = fmt.Sprintf("%s%d VERS %s", indent(r.Level+1), r.Level+1, r.Version)
		ss = append(ss, s)
	}

	if r.Business != nil {
		s = r.Business.String()
		ss = append(ss, s)
	}

	if r.SourceData != nil {
		s = r.SourceData.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

// String stringifies a GEDCOM level 0 trailer record
func (r *TrailerRecord) String() string {
	return fmt.Sprintf("%s%d TRLR\n", indent(0), 0)
}

// String stringifies a GEDCOM user reference number record
func (r *UserReferenceNumberRecord) String() string {
	var ss []string
	var s string

	s = fmt.Sprintf("%s%d REFN %s", indent(r.Level), r.Level, r.UserReferenceNumber)
	ss = append(ss, s)

	if r.Type != "" {
		s = fmt.Sprintf("%s%d TYPE %s", indent(r.Level+1), r.Level+1, r.Type)
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

// String stringifies a slice of user reference number records
func (r UserReferenceNumberRecords) String() string {
	var ss []string
	var s string

	for _, x := range r {
		s = x.String()
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

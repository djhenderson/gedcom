/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// CheckError logs a fatal error when passe a non-nil error code.
func CheckError(err error) {
	if err != nil {
		log.Fatalln("writers: i/o error: ", err.Error())
	}
}

// LongWrite formats a long string using CONT and CONC lines
func LongWrite(w io.Writer, level int, xref string, tag string, longString string) (int, error) {

	sXref0, sXrefN := "", ""
	if xref != "" {
		sXref0 = fmt.Sprintf(" @%s@", xref)
		if level != 0 {
			sXrefN, sXref0 = sXref0, sXrefN
		}
	}

	if longString == "" {
		_, err := fmt.Fprintf(w, "%s%d%s %s%s\n", indent(level), level, sXref0, tag, sXrefN)
		CheckError(err)

		return 0, nil
	}

	parts := strings.Split(longString, "\n")
	for i, part := range parts {
		_, err := fmt.Fprintf(w, "%s%d%s %s%s %s\n", indent(level), level, sXref0, tag, sXrefN, part)
		CheckError(err)

		if i == 0 {
			tag = "CONT"
			level++
			sXref0, sXrefN = "", ""
		}
	}
	return 0, nil
}

// Write formats and writes a GEDCOM address record
func (r *AddressRecord) Write(w io.Writer) (int, error) {

	_, err := LongWrite(w, r.Level, "", "ADDR", r.Full)
	CheckError(err)

	if r.Line1 != "" {
		_, err := fmt.Fprintf(w, "%s%d ADR1 %s\n", indent(r.Level+1), r.Level+1, r.Line1)
		CheckError(err)
	}

	if r.Line2 != "" {
		_, err := fmt.Fprintf(w, "%s%d ADR2 %s\n", indent(r.Level+1), r.Level+1, r.Line2)
		CheckError(err)
	}

	if r.Line3 != "" {
		_, err := fmt.Fprintf(w, "%s%d ADR3 %s\n", indent(r.Level+1), r.Level+1, r.Line3)
		CheckError(err)
	}

	if r.City != "" {
		_, err := fmt.Fprintf(w, "%s%d CITY %s\n", indent(r.Level+1), r.Level+1, r.City)
		CheckError(err)
	}

	if r.State != "" {
		_, err := fmt.Fprintf(w, "%s%d STAE %s\n", indent(r.Level+1), r.Level+1, r.State)
		CheckError(err)
	}

	if r.PostalCode != "" {
		_, err := fmt.Fprintf(w, "%s%d POST %s\n", indent(r.Level+1), r.Level+1, r.PostalCode)
		CheckError(err)
	}

	if r.Country != "" {
		_, err := fmt.Fprintf(w, "%s%d CTRY %s\n", indent(r.Level+1), r.Level+1, r.Country)
		CheckError(err)
	}

	if r.Phone != "" {
		_, err := fmt.Fprintf(w, "%s%d PHON %s\n", indent(r.Level+1), r.Level+1, r.Phone)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of address records
func (r AddressRecords) Write(w io.Writer) (int, error) {

	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM bibliography record
func (r *BibliographyRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d BIBL %s\n", indent(r.Level), r.Level, r.Value)
	CheckError(err)

	for _, comp := range r.Component {
		_, err := LongWrite(w, r.Level+1, "", "COMP", comp)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM blob record
func (r *BlobRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d BLOB\n", indent(r.Level), r.Level)
	CheckError(err)

	parts := strings.Split(r.Data, "\n")
	for _, data := range parts {
		_, err := fmt.Fprintf(w, "%s%d CONT %s\n", indent(r.Level+1), r.Level+1, data)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM business record
func (r *BusinessRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d CORP %s\n", indent(r.Level), r.Level, r.BusinessName)
	CheckError(err)

	if r.Address != nil {
		_, err := r.Address.Write(w)
		CheckError(err)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			_, err := fmt.Fprintf(w, "%s%d PHON %s\n", indent(r.Level+1), r.Level+1, phone)
			CheckError(err)
		}
	}

	if r.WebSite != "" {
		_, err := fmt.Fprintf(w, "%s%d WWW %s\n", indent(r.Level), r.Level, r.WebSite)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM call number record
func (r *CallNumberRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d CALN %s\n", indent(r.Level), r.Level, r.CallNumber)
	CheckError(err)

	// In this unique case, the value is a string, not a media record or link
	if r.Media != "" {
		_, err := fmt.Fprintf(w, "%s%d MEDI %s\n", indent(r.Level+1), r.Level+1, r.Media)
		CheckError(err)
	}
	return 0, nil
}

// Write formats and writes a GEDCOM change record
func (r *ChangeRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d CHAN\n", indent(r.Level), r.Level)
	CheckError(err)

	if r.Date != nil {
		_, err := r.Date.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM character set record
func (r *CharacterSetRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d CHAR %s\n", indent(r.Level), r.Level, r.CharacterSet)
	CheckError(err)

	if r.Version != "" {
		_, err := fmt.Fprintf(w, "%s%d VERS %s\n", indent(r.Level+1), r.Level+1, r.Version)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM child status record
func (r *ChildStatusRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d @%s@ CSTA\n", indent(r.Level), r.Level, r.Xref)
	CheckError(err)

	_, err = fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level+1), r.Level+1, r.Name)
	CheckError(err)

	return 0, nil
}

// Write formats and writes a slice of child status records
func (r ChildStatusRecords) Write(w io.Writer) (int, error) {

	// log.Printf("ChildStatusRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM citatio record
func (r *CitationRecord) Write(w io.Writer) (int, error) {

	xref := ""
	if r.Source != nil {
		xref = r.Source.Xref
	}
	_, err := LongWrite(w, r.Level, xref, "SOUR", r.Value)
	CheckError(err)

	if r.Page != "" {
		_, err := fmt.Fprintf(w, "%s%d PAGE %s\n", indent(r.Level+1), r.Level+1, r.Page)
		CheckError(err)
	}

	if r.Reference != "" {
		_, err := fmt.Fprintf(w, "%s%d REF %s\n", indent(r.Level+1), r.Level+1, r.Reference)
		CheckError(err)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			_, err := event.Write(w)
			CheckError(err)
		}
	}

	if r.Data != nil {
		for _, data := range r.Data {
			_, err := data.Write(w)
			CheckError(err)
		}
	}

	if r.Text != nil {
		for _, text := range r.Text {
			_, err := LongWrite(w, r.Level+1, "", "TEXT", text)
			CheckError(err)
		}
	}

	if r.Quality != "" {
		_, err := fmt.Fprintf(w, "%s%d QUAY %s\n", indent(r.Level+1), r.Level+1, r.Quality)
		CheckError(err)
	}

	if r.Media != nil {
		_, err := r.Media.Write(w)
		CheckError(err)
	}

	if r.CONS != "" {
		_, err := fmt.Fprintf(w, "%s%d CONS %s\n", indent(r.Level+1), r.Level+1, r.CONS)
		CheckError(err)
	}

	if r.Direct != "" {
		_, err := fmt.Fprintf(w, "%s%d DIRE %s\n", indent(r.Level+1), r.Level+1, r.Direct)
		CheckError(err)
	}

	if r.SourceQuality != "" {
		_, err := fmt.Fprintf(w, "%s%d SOQU %s\n", indent(r.Level+1), r.Level+1, r.SourceQuality)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of citation records
func (r CitationRecords) Write(w io.Writer) (int, error) {

	//log.Printf("CitationRecords type(r): %T\n", r)
	for _, note := range r {
		//log.Printf("CitationRecords type(note): %T\n", note)
		//log.Printf("CitationRecords type(*note): %T\n", *note)
		_, err := note.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM data record
func (r *DataRecord) Write(w io.Writer) (int, error) {

	_, err := LongWrite(w, r.Level, "", "DATA", r.Data)
	CheckError(err)

	if r.Date != "" {
		_, err := fmt.Fprintf(w, "%s%d DATE %s\n", indent(r.Level+1), r.Level+1, r.Date)
		CheckError(err)
	}

	if r.Copyright != "" {
		_, err := fmt.Fprintf(w, "%s%d COPR %s\n", indent(r.Level+1), r.Level+1, r.Copyright)
		CheckError(err)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			_, err := LongWrite(w, r.Level+1, "", "TEXT", text)
			CheckError(err)
		}
	}

	return 0, nil
}

// Write formats and writes a slice of data records
func (r DataRecords) Write(w io.Writer) (int, error) {

	//log.Printf("DataRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM date record
func (r *DateRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d DATE %s\n", indent(r.Level), r.Level, r.Date)
	CheckError(err)

	if r.Day != "" {
		_, err := fmt.Fprintf(w, "%s%d DATD %s\n", indent(r.Level+1), r.Level+1, r.Day)
		CheckError(err)
	}

	if r.Month != "" {
		_, err := fmt.Fprintf(w, "%s%d DATM %s\n", indent(r.Level+1), r.Level+1, r.Month)
		CheckError(err)
	}

	if r.Year != "" {
		_, err := fmt.Fprintf(w, "%s%d DATY %s\n", indent(r.Level+1), r.Level+1, r.Year)
		CheckError(err)
	}

	if r.Full != "" {
		_, err := fmt.Fprintf(w, "%s%d DATF %s\n", indent(r.Level+1), r.Level+1, r.Full)
		CheckError(err)

	}

	if r.Short != "" {
		_, err := fmt.Fprintf(w, "%s%d DATS %s\n", indent(r.Level+1), r.Level+1, r.Short)
		CheckError(err)
	}

	if r.Time != "" {
		_, err := fmt.Fprintf(w, "%s%d TIME %s\n", indent(r.Level+1), r.Level+1, r.Time)
		CheckError(err)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			_, err := LongWrite(w, r.Level+1, "", "TEXT", text)
			CheckError(err)
		}
	}

	return 0, nil
}

// Write formats and writes a GEDCOM event definition record
func (r EventDefinitionRecord) Write(w io.Writer) (int, error) {
	log.Fatal("EventDefinition.Write() not implemented")
	return 0, nil
}

// Write formats and writes a slice of event definition records
func (r EventDefinitionRecords) Write(w io.Writer) (int, error) {

	//log.Printf("EventDefinitionRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM event records
func (r *EventRecord) Write(w io.Writer) (int, error) {

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}
	_, err := fmt.Fprintf(w, "%s%d %s%s %s\n", indent(r.Level), r.Level, id, r.Tag, r.Value)
	CheckError(err)

	if r.Type != "" {
		_, err := fmt.Fprintf(w, "%s%d TYPE %s\n", indent(r.Level+1), r.Level+1, r.Type)
		CheckError(err)
	}

	if r.Name != "" {
		_, err := fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level+1), r.Level+1, r.Name)
		CheckError(err)
	}

	if r.Primary_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _PRIM %s\n", indent(r.Level+1), r.Level+1, r.Primary_)
		CheckError(err)
	}

	if r.Date != nil {
		_, err := r.Date.Write(w)
		CheckError(err)
	}

	if r.Place != nil {
		_, err := r.Place.Write(w)
		CheckError(err)
	}

	if r.Role != nil {
		_, err := r.Role.Write(w)
		CheckError(err)
	}
	if r.Address != nil {
		_, err := r.Address.Write(w)
		CheckError(err)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			_, err := fmt.Fprintf(w, "%s%d PHON %s\n", indent(r.Level+1), r.Level+1, phone)
			CheckError(err)
		}
	}

	if r.Parents != nil {
		_, err := r.Parents.Write(w)
		CheckError(err)
	}

	if r.Husband != nil {
		_, err := r.Husband.Write(w)
		CheckError(err)
	}

	if r.Wife != nil {
		_, err := r.Wife.Write(w)
		CheckError(err)
	}

	if r.Spouse != nil {
		_, err := r.Spouse.Write(w)
		CheckError(err)
	}

	if r.Age != "" {
		_, err := fmt.Fprintf(w, "%s%d AGE %s\n", indent(r.Level+1), r.Level+1, r.Age)
		CheckError(err)
	}

	if r.Agency != "" {
		_, err := fmt.Fprintf(w, "%s%d AGNC %s\n", indent(r.Level+1), r.Level+1, r.Agency)
		CheckError(err)
	}

	if r.Cause != "" {
		_, err := fmt.Fprintf(w, "%s%d CAUS %s\n", indent(r.Level+1), r.Level+1, r.Cause)
		CheckError(err)
	}

	if r.Temple != "" {
		_, err := fmt.Fprintf(w, "%s%d TEMP %s\n", indent(r.Level+1), r.Level+1, r.Temple)
		CheckError(err)
	}

	if r.Status != "" {
		_, err := fmt.Fprintf(w, "%s%d STAT %s\n", indent(r.Level+1), r.Level+1, r.Status)
		CheckError(err)
	}

	if r.Media != nil {
		_, err := r.Media.Write(w)
		CheckError(err)
	}

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	if r.UID_ != nil {
		for _, uid := range r.UID_ {
			_, err := fmt.Fprintf(w, "%s%d _UID %s\n", indent(r.Level+1), r.Level+1, uid)
			CheckError(err)
		}
	}

	if r.UpdateTime_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _UPD %s\n", indent(r.Level+1), r.Level+1, r.UpdateTime_)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of event records
func (r EventRecords) Write(w io.Writer) (int, error) {

	//log.Printf("EventRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := fmt.Fprintf(w, "%s", x.String( /**/ ))
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM link to a family record
func (r *FamilyLink) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d %s @%s@\n", indent(r.Level), r.Level, r.Tag, r.Family.Xref)
	CheckError(err)

	if r.Pedigree != "" {
		_, err := fmt.Fprintf(w, "%s%d PEDI %s\n", indent(r.Level+1), r.Level+1, r.Pedigree)
		CheckError(err)
	}

	if r.Adopted != "" {
		_, err := fmt.Fprintf(w, "%s%d ADOP %s\n", indent(r.Level+1), r.Level+1, r.Adopted)
		CheckError(err)
	}

	if r.Primary_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _PRIMARY %s\n", indent(r.Level+1), r.Level+1, r.Primary_)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of links to family records
func (r FamilyLinks) Write(w io.Writer) (int, error) {

	//log.Printf("FamilyLinks type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM family record
func (r *FamilyRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d @%s@ FAM\n", indent(r.Level), r.Level, r.Xref)
	CheckError(err)

	if r.Husband != nil {
		_, err := r.Husband.Write(w)
		CheckError(err)
	}

	if r.Wife != nil {
		_, err := r.Wife.Write(w)
		CheckError(err)
	}

	if r.NumChildren > 0 {
		_, err := fmt.Fprintf(w, "%s%d NCHI %d\n", indent(r.Level+1), r.Level+1, r.NumChildren)
		CheckError(err)
	}

	if r.Child != nil {
		for _, child := range r.Child {
			_, err := child.Write(w)
			CheckError(err)
		}
	}

	if r.Event != nil {
		for _, event := range r.Event {
			_, err := event.Write(w)
			CheckError(err)
		}
	}

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			_, err := subm.Write(w)
			CheckError(err)
		}
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	if r.UID_ != nil {
		for _, uid := range r.UID_ {
			_, err := fmt.Fprintf(w, "%s%d _UID %s\n", indent(r.Level+1), r.Level+1, uid)
			CheckError(err)
		}
	}

	if r.UpdateTime_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _UPD %s\n", indent(r.Level+1), r.Level+1, r.UpdateTime_)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of family records
func (r FamilyRecords) Write(w io.Writer) (int, error) {

	//log.Printf("FamilyRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM footnote record
func (r *FootnoteRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d FOOT %s\n", indent(r.Level), r.Level, r.Value)
	CheckError(err)

	for _, comp := range r.Component {
		_, err := LongWrite(w, r.Level+1, "", "COMP", comp)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM gedcom record
func (r *GedcomRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d GEDC\n", indent(r.Level), r.Level)
	CheckError(err)

	_, err = fmt.Fprintf(w, "%s%d VERS %s\n", indent(r.Level+1), r.Level+1, r.Version)
	CheckError(err)

	_, err = fmt.Fprintf(w, "%s%d FORM %s\n", indent(r.Level+1), r.Level+1, r.Form)
	CheckError(err)

	return 0, nil

}

// Write formats and writes a GEDCOM header record
func (r *HeaderRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d HEAD\n", indent(r.Level), r.Level)
	CheckError(err)

	if r.CharacterSet != nil {
		_, err := r.CharacterSet.Write(w)
		CheckError(err)
	}

	if r.SourceSystem != nil {
		_, err := r.SourceSystem.Write(w)
		CheckError(err)
	}

	if r.Destination != "" {
		_, err := fmt.Fprintf(w, "%s%d DEST %s\n", indent(r.Level+1), r.Level+1, r.Destination)
		CheckError(err)
	}

	if r.Date != nil {
		_, err := r.Date.Write(w)
		CheckError(err)
	}

	if r.Time != "" {
		_, err := fmt.Fprintf(w, "%s%d TIME %s\n", indent(r.Level+1), r.Level+1, r.Time)
		CheckError(err)
	}

	if r.FileName != "" {
		_, err := fmt.Fprintf(w, "%s%d FILE %s\n", indent(r.Level+1), r.Level+1, r.FileName)
		CheckError(err)
	}

	if r.Gedcom != nil {
		_, err := r.Gedcom.Write(w)
		CheckError(err)
	}

	if r.Language != "" {
		_, err := fmt.Fprintf(w, "%s%d LANG %s\n", indent(r.Level+1), r.Level+1, r.Language)
		CheckError(err)
	}

	if r.Copyright != "" {
		_, err := fmt.Fprintf(w, "%s%d COPR %s\n", indent(r.Level+1), r.Level+1, r.Copyright)
		CheckError(err)
	}

	if r.Place != nil {
		_, err := r.Place.Write(w)
		CheckError(err)
	}

	if r.RootPerson_ != nil {
		_, err := r.RootPerson_.Write(w)
		CheckError(err)
	}

	if r.HomePerson_ != nil {
		_, err := r.HomePerson_.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			_, err := subm.Write(w)
			CheckError(err)
		}
	}

	if r.Submission != nil {
		for _, subn := range r.Submission {
			_, err := subn.Write(w)
			CheckError(err)
		}
	}

	if r.Schema != nil {
		_, err := r.Schema.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM history record
func (r *HistoryRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d HIST %s\n", indent(r.Level), r.Level, r.History)
	CheckError(err)

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of history records
func (r HistoryRecords) Write(w io.Writer) (int, error) {

	for _, x := range r {
		_, err := fmt.Fprintf(w, "%s\n", x.String( /**/ ))
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM individual links
func (r *IndividualLink) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d %s @%s@\n", indent(r.Level), r.Level, r.Tag, r.Individual.Xref)
	CheckError(err)

	if r.Relationship != "" {
		_, err := fmt.Fprintf(w, "%s%d RELA %s\n", indent(r.Level+1), r.Level+1, r.Relationship)
		CheckError(err)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			_, err := event.Write(w)
			CheckError(err)
		}
	}

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of links to individual records
func (r IndividualLinks) Write(w io.Writer) (int, error) {

	//log.Printf("IndividualLinks type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM individual record
func (r *IndividualRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d @%s@ INDI\n", indent(r.Level), r.Level, r.Xref)
	CheckError(err)

	for _, name := range r.Name {
		_, err := name.Write(w)
		CheckError(err)
	}

	if r.Restriction != "" {
		_, err := fmt.Fprintf(w, "%s%d RESN %s\n", indent(r.Level+1), r.Level+1, r.Restriction)
		CheckError(err)
	}

	if r.Sex != "" {
		_, err := fmt.Fprintf(w, "%s%d SEX %s\n", indent(r.Level+1), r.Level+1, r.Sex)
		CheckError(err)
	}

	if r.ProfilePicture_ != nil {
		_, err := r.ProfilePicture_.Write(w)
		CheckError(err)
	}

	if r.CONL != "" {
		_, err := fmt.Fprintf(w, "%s%d CONL %s\n", indent(r.Level+1), r.Level+1, r.CONL)
		CheckError(err)
	}

	if r.Event != nil {
		for _, event := range r.Event {
			_, err := event.Write(w)
			CheckError(err)
		}
	}

	if r.Attribute != "" {
		_, err := fmt.Fprintf(w, "%s%d ATTR %s\n", indent(r.Level+1), r.Level+1, r.Attribute)
		CheckError(err)
	}

	if r.Parents != nil {
		for _, famc := range r.Parents {
			_, err := famc.Write(w)
			CheckError(err)
		}
	}

	if r.Family != nil {
		for _, fams := range r.Family {
			_, err := fams.Write(w)
			CheckError(err)
		}
	}

	if r.Address != nil {
		_, err := r.Address.Write(w)
		CheckError(err)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			_, err := fmt.Fprintf(w, "%s%d PHON %s\n", indent(r.Level+1), r.Level+1, phone)
			CheckError(err)
		}
	}

	if r.Health != "" {
		_, err := fmt.Fprintf(w, "%s%d HEAL %s\n", indent(r.Level+1), r.Level+1, r.Health)
		CheckError(err)
	}

	if r.History != nil {
		_, err := r.History.Write(w)
		CheckError(err)
	}

	if r.Quality != "" {
		_, err := fmt.Fprintf(w, "%s%d QUAY %s\n", indent(r.Level+1), r.Level+1, r.Quality)
		CheckError(err)
	}

	if r.Living != "" {
		_, err := fmt.Fprintf(w, "%s%d LVG %s\n", indent(r.Level+1), r.Level+1, r.Living)
		CheckError(err)
	}

	if r.AncestralFileNumber != nil {
		for _, afn := range r.AncestralFileNumber {
			_, err := fmt.Fprintf(w, "%s%d AFN %s\n", indent(r.Level+1), r.Level+1, afn)
			CheckError(err)
		}
	}

	if r.RecordFileNumber != "" {
		_, err := fmt.Fprintf(w, "%s%d RFN %s\n", indent(r.Level+1), r.Level+1, r.RecordFileNumber)
		CheckError(err)
	}

	if r.UserReferenceNumber != nil {
		_, err := r.UserReferenceNumber.Write(w)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	if r.UID_ != nil {
		for _, uid := range r.UID_ {
			_, err := fmt.Fprintf(w, "%s%d _UID %s\n", indent(r.Level+1), r.Level+1, uid)
			CheckError(err)
		}
	}

	if r.Email != "" {
		_, err := fmt.Fprintf(w, "%s%d EMAIL %s\n", indent(r.Level+1), r.Level+1, r.Email)
		CheckError(err)
	}

	if r.WebSite != "" {
		_, err := fmt.Fprintf(w, "%s%d WWW %s\n", indent(r.Level+1), r.Level+1, r.WebSite)
		CheckError(err)
	}

	if r.Media != nil {
		_, err := r.Media.Write(w)
		CheckError(err)
	}

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Associated != nil {
		_, err := r.Associated.Write(w)
		CheckError(err)
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			_, err := subm.Write(w)
			CheckError(err)
		}
	}

	if r.ANCI != nil {
		for _, subm := range r.ANCI {
			_, err := subm.Write(w)
			CheckError(err)
		}
	}

	if r.DESI != nil {
		for _, subm := range r.DESI {
			_, err := subm.Write(w)
			CheckError(err)
		}
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	if r.Alias != "" {
		_, err := fmt.Fprintf(w, "%s%d ALIA %s\n", indent(r.Level+1), r.Level+1, r.Alias)
		CheckError(err)
	}

	if r.Father != nil {
		_, err := r.Father.Write(w)
		CheckError(err)
	}

	if r.Mother != nil {
		_, err := r.Mother.Write(w)
		CheckError(err)
	}

	if r.Miscellaneous != nil {
		for _, misc := range r.Miscellaneous {
			_, err := fmt.Fprintf(w, "%s%d MISC %s\n", indent(r.Level+1), r.Level+1, misc)
			CheckError(err)
		}
	}

	return 0, nil
}

// Write formats and writes a slice of individual records
func (r IndividualRecords) Write(w io.Writer) (int, error) {

	//log.Printf("IndividualRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM media links
func (r *MediaLink) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d %s @%s@\n", indent(r.Level), r.Level, r.Tag, r.Media.Xref)
	CheckError(err)

	return 0, nil
}

// Write formats and writes a slice of links to media records
func (r MediaLinks) Write(w io.Writer) (int, error) {

	//log.Printf("MediaLinks type(r): %T\n", r)
	for _, x := range r {
		var err error
		if x.Value != "" {
			_, err = x.Write(w)
		} else {
			_, err = x.Media.Write(w)
		}
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM media record
func (r *MediaRecord) Write(w io.Writer) (int, error) {

	id0, idN := "", ""
	if r.Xref != "" {
		id0 = fmt.Sprintf(" @%s@ ", r.Xref)

		if r.Level != 0 {
			idN, id0 = id0, idN
		}
	}

	_, err := fmt.Fprintf(w, "%s%d%s %s%s\n", indent(r.Level), r.Level, id0, "OBJE", idN)
	CheckError(err)

	if r.Format != "" {
		_, err := fmt.Fprintf(w, "%s%d FORM %s\n", indent(r.Level+1), r.Level+1, r.Format)
		CheckError(err)
	}

	if r.URL_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _URL %s\n", indent(r.Level+1), r.Level+1, r.URL_)
		CheckError(err)
	}

	if r.FileName != "" {
		_, err := fmt.Fprintf(w, "%s%d FILE %s\n", indent(r.Level+1), r.Level+1, r.FileName)
		CheckError(err)

	}

	if r.Title != "" {
		_, err := fmt.Fprintf(w, "%s%d TITL %s\n", indent(r.Level+1), r.Level+1, r.Title)
		CheckError(err)
	}

	if r.Date != "" {
		_, err := fmt.Fprintf(w, "%s%d DATE %s\n", indent(r.Level+1), r.Level+1, r.Date)
		CheckError(err)
	}

	if r.Author != "" {
		_, err := fmt.Fprintf(w, "%s%d AUTH %s\n", indent(r.Level+1), r.Level+1, r.Author)
		CheckError(err)
	}

	if r.Text != "" {
		_, err := LongWrite(w, r.Level, "", "TEXT", r.Text)
		CheckError(err)
	}

	if r.Date_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _DATE %s\n", indent(r.Level+1), r.Level+1, r.Date_)
		CheckError(err)
	}

	if r.AstId_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _ASTID %s\n", indent(r.Level+1), r.Level+1, r.AstId_)
		CheckError(err)
	}

	if r.AstType_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _ASTTYP %s\n", indent(r.Level+1), r.Level+1, r.AstType_)
		CheckError(err)
	}

	if r.AstDesc_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _ASTDESC %s\n", indent(r.Level+1), r.Level+1, r.AstDesc_)
		CheckError(err)
	}

	if r.AstLoc_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _ASTLOC %s\n", indent(r.Level+1), r.Level+1, r.AstLoc_)
		CheckError(err)
	}

	if r.AstPerm_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _ASTPERM %s\n", indent(r.Level+1), r.Level+1, r.AstPerm_)
		CheckError(err)
	}

	if r.AstUpPid_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _ASTUPPID %s\n", indent(r.Level+1), r.Level+1, r.AstUpPid_)
		CheckError(err)
	}

	if r.BinaryLargeObject != nil {
		_, err := r.BinaryLargeObject.Write(w)
		CheckError(err)
	}

	if r.UserReferenceNumber != nil {
		_, err := r.UserReferenceNumber.Write(w)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of media records
func (r MediaRecords) Write(w io.Writer) (int, error) {

	//log.Printf("MediaRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM name record
func (r *NameRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level), r.Level, r.Name)
	CheckError(err)

	if r.Prefix != "" {
		_, err := fmt.Fprintf(w, "%s%d NPFX %s\n", indent(r.Level+1), r.Level+1, r.Prefix)
		CheckError(err)
	}

	if r.GivenName != "" {
		_, err := fmt.Fprintf(w, "%s%d GIVN %s\n", indent(r.Level+1), r.Level+1, r.GivenName)
		CheckError(err)
	}

	if r.MiddleName_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _MIDN %s\n", indent(r.Level+1), r.Level+1, r.MiddleName_)
		CheckError(err)
	}

	if r.SurnamePrefix != "" {
		_, err := fmt.Fprintf(w, "%s%d SPFX %s\n", indent(r.Level+1), r.Level+1, r.SurnamePrefix)
		CheckError(err)
	}

	if r.Surname != "" {
		_, err := fmt.Fprintf(w, "%s%d SURN %s\n", indent(r.Level+1), r.Level+1, r.Surname)
		CheckError(err)
	}

	if r.Suffix != "" {
		_, err := fmt.Fprintf(w, "%s%d NSFX %s\n", indent(r.Level+1), r.Level+1, r.Suffix)
		CheckError(err)
	}

	if r.PreferedGivenName_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _PGVN %s\n", indent(r.Level+1), r.Level+1, r.PreferedGivenName_)
		CheckError(err)
	}

	if r.NameType != "" {
		_, err := fmt.Fprintf(w, "%s%d _PRIM %s\n", indent(r.Level+1), r.Level+1, r.Primary_)
		CheckError(err)
	}

	if r.Primary_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _PRIM %s\n", indent(r.Level+1), r.Level+1, r.Primary_)
		CheckError(err)
	}

	if r.AKA_ != nil {
		for _, aka := range r.AKA_ {
			_, err := fmt.Fprintf(w, "%s%d _AKA %s\n", indent(r.Level+1), r.Level+1, aka)
			CheckError(err)
		}
	}

	if r.Nickname != nil {
		for _, nick := range r.Nickname {
			_, err := fmt.Fprintf(w, "%s%d NICK %s\n", indent(r.Level+1), r.Level+1, nick)
			CheckError(err)
		}
	}

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		for _, note := range r.Note {
			_, err := note.Write(w)
			CheckError(err)
		}
	}

	return 0, nil
}

// Write formats and writes a slice of name records
func (r NameRecords) Write(w io.Writer) (int, error) {

	//log.Printf("NameRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM note records
func (r *NoteRecord) Write(w io.Writer) (int, error) {

	_, err := LongWrite(w, r.Level, r.Xref, "NOTE", r.Note)
	CheckError(err)

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.UserReferenceNumber != nil {
		_, err := r.UserReferenceNumber.Write(w)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of note records
func (r NoteRecords) Write(w io.Writer) (int, error) {

	//log.Printf("NoteRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM place part record
func (r *PlacePartRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d %s %s\n", indent(r.Level), r.Level, r.Tag, r.Part)
	CheckError(err)

	if r.Jurisdiction != "" {
		_, err := fmt.Fprintf(w, "%s%d JURI %s\n", indent(r.Level+1), r.Level+1, r.Jurisdiction)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM place record
func (r *PlaceRecord) Write(w io.Writer) (int, error) {

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}

	_, err := fmt.Fprintf(w, "%s%d %sPLAC %s\n", indent(r.Level), r.Level, id, r.Name)
	CheckError(err)

	if r.Form != "" {
		_, err := fmt.Fprintf(w, "%s%d PLAS %s\n", indent(r.Level+1), r.Level+1, r.Form)
		CheckError(err)
	}

	if r.ShortName != "" {
		_, err := fmt.Fprintf(w, "%s%d PLAS %s\n", indent(r.Level+1), r.Level+1, r.ShortName)
		CheckError(err)
	}

	if r.Modifier != "" {
		_, err := fmt.Fprintf(w, "%s%d PLAM %s\n", indent(r.Level+1), r.Level+1, r.Modifier)
		CheckError(err)
	}

	if r.Parts != nil {
		for _, part := range r.Parts {
			_, err := part.Write(w)
			CheckError(err)
		}
	}

	if r.Citation != nil {
		_, err := r.Citation.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of place records
func (r PlaceRecords) Write(w io.Writer) (int, error) {

	//log.Printf("PlaceRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM link to a repository record
func (r *RepositoryLink) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d REPO @%s@\n", indent(r.Level), r.Level, r.Repository.Xref)
	CheckError(err)

	if r.CallNumber != nil {
		_, err := r.CallNumber.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of repository links
func (r RepositoryLinks) Write(w io.Writer) (int, error) {

	//log.Printf("RepositoryLinks type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM repository record
func (r *RepositoryRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d @%s@ REPO\n", indent(r.Level), r.Level, r.Xref)
	CheckError(err)

	if r.Name != "" {
		_, err := fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level+1), r.Level+1, r.Name)
		CheckError(err)
	}

	if r.Address != nil {
		_, err := r.Address.Write(w)
		CheckError(err)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			_, err := fmt.Fprintf(w, "%s%d PHON %s\n", indent(r.Level+1), r.Level+1, phone)
			CheckError(err)
		}
	}

	if r.WebSite != "" {
		_, err := fmt.Fprintf(w, "%s%d WWW %s\n", indent(r.Level+1), r.Level+1, r.WebSite)
		CheckError(err)
	}

	if r.UserReferenceNumber != nil {
		_, err := r.UserReferenceNumber.Write(w)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of repository records
func (r RepositoryRecords) Write(w io.Writer) (int, error) {

	//log.Printf("RepositoryRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM role record
func (r *RoleRecord) Write(w io.Writer) (int, error) {

	spacer := " "
	if r.Role == "" {
		spacer = ""
	}
	_, err := fmt.Fprintf(w, "%s%d ROLE %s%s@%s@\n", indent(r.Level), r.Level, r.Role, spacer, r.Individual.Xref)
	CheckError(err)

	if r.Principal != "" {
		_, err := fmt.Fprintf(w, "%s%d PRIN %s\n", indent(r.Level+1), r.Level+1, r.Principal)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of role record
func (r RoleRecords) Write(w io.Writer) (int, error) {

	//log.Printf("RoleRecords type(r): %T\n", r)
	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM root record, i.e. the whole file
func (r *RootRecord) Write(w io.Writer) (int, error) {

	if r.Header != nil {
		_, err := r.Header.Write(w)
		CheckError(err)
	}

	if r.Submission != nil {
		for _, subn := range r.Submission {
			_, err := subn.Write(w)
			CheckError(err)
		}
	}

	if r.Submitter != nil {
		for _, subm := range r.Submitter {
			_, err := subm.Write(w)
			CheckError(err)
		}
	}

	if r.Place != nil {
		for _, x := range r.Place {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Event != nil {
		for _, x := range r.Event {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Individual != nil {
		for _, x := range r.Individual {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Family != nil {
		for _, x := range r.Family {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Repository != nil {
		for _, x := range r.Repository {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Source != nil {
		for _, x := range r.Source {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Note != nil {
		for _, x := range r.Note {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Media != nil {
		_, err := r.Media.Write(w)
		CheckError(err)
	}

	if r.ChildStatus != nil {
		for _, x := range r.ChildStatus {
			_, err := x.Write(w)
			CheckError(err)
		}
	}

	if r.Trailer != nil {
		_, err := r.Trailer.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM schema record
func (r *SchemaRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d SCHEMA\n", indent(r.Level), r.Level)
	CheckError(err)

	for _, data := range r.Data {
		level := int(data[0])
		_, err := fmt.Fprintf(w, "%s\n", indent(level)+data)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM short title record
func (r *ShortTitleRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d SHTI %s\n", indent(r.Level), r.Level, r.ShortTitle)
	CheckError(err)

	if r.Indexed != "" {
		_, err := fmt.Fprintf(w, "%s%d INDX %s\n", indent(r.Level+1), r.Level+1, r.Indexed)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM level 0 source record
func (r *SourceRecord) Write(w io.Writer) (int, error) {

	_, err := LongWrite(w, r.Level, r.Xref, "SOUR", r.Value)
	CheckError(err)

	if r.Name != "" {
		_, err := fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level+1), r.Level+1, r.Name)
		CheckError(err)
	}

	if r.Title != "" {
		_, err := LongWrite(w, r.Level+1, "", "TITL", r.Title)
		CheckError(err)
	}

	if r.Author != "" {
		_, err := fmt.Fprintf(w, "%s%d AUTH %s\n", indent(r.Level+1), r.Level+1, r.Author)
		CheckError(err)
	}

	if r.Abbreviation != "" {
		_, err := fmt.Fprintf(w, "%s%d ABBR %s\n", indent(r.Level+1), r.Level+1, r.Abbreviation)
		CheckError(err)
	}

	if r.Publication != "" {
		_, err := fmt.Fprintf(w, "%s%d PUBL %s\n", indent(r.Level+1), r.Level+1, r.Publication)
		CheckError(err)
	}

	if r.Parenthesized_ != "" {
		_, err := fmt.Fprintf(w, "%s%d _PAREN %s\n", indent(r.Level+1), r.Level+1, r.Parenthesized_)
		CheckError(err)
	}

	if r.Text != nil {
		for _, text := range r.Text {
			_, err := LongWrite(w, r.Level+1, "", "TEXT", text)
			CheckError(err)
		}
	}

	if r.Data != nil {
		_, err := r.Data.Write(w)
		CheckError(err)
	}

	if r.Footnote != nil {
		_, err := r.Footnote.Write(w)
		CheckError(err)
	}

	if r.Bibliography != nil {
		_, err := r.Bibliography.Write(w)
		CheckError(err)
	}

	if r.Repository != nil {
		_, err := r.Repository.Write(w)
		CheckError(err)
	}

	if r.UserReferenceNumber != nil {
		_, err := r.UserReferenceNumber.Write(w)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	if r.ShortAuthor != "" {
		_, err := fmt.Fprintf(w, "%s%d SHAU %s\n", indent(r.Level+1), r.Level+1, r.ShortAuthor)
		CheckError(err)
	}

	if r.ShortTitle != nil {
		_, err := r.ShortTitle.Write(w)
		CheckError(err)
	}

	if r.Media != nil {
		_, err := r.Media.Write(w)
		CheckError(err)
	}

	if r.Note != nil {
		_, err := r.Note.Write(w)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}
	return 0, nil
}

// Write formats and writes a slice of source records
func (r SourceRecords) Write(w io.Writer) (int, error) {

	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}
	return 0, nil
}

// Write formats and writes a GEDCOM link to a submission record
func (r *SubmissionLink) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d SUBN @%s@\n", indent(r.Level), r.Level, r.Submission.Xref)
	CheckError(err)

	return 0, nil
}

// Write formats and writes a GEDCOM level 0 submission record
func (r *SubmissionRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d @%s@ SUBN\n", indent(r.Level), r.Level, r.Xref)
	CheckError(err)

	if r.FamilyFileName != "" {
		_, err := fmt.Fprintf(w, "%s%d FAMF %s\n", indent(r.Level+1), r.Level+1, r.FamilyFileName)
		CheckError(err)
	}

	if r.Temple != "" {
		_, err := fmt.Fprintf(w, "%s%d TEMP %s\n", indent(r.Level+1), r.Level+1, r.Temple)
		CheckError(err)
	}

	if r.Ancestors != "" {
		_, err := fmt.Fprintf(w, "%s%d ANCE %s\n", indent(r.Level+1), r.Level+1, r.Ancestors)
		CheckError(err)
	}

	if r.Descendents != "" {
		_, err := fmt.Fprintf(w, "%s%d DESC %s\n", indent(r.Level+1), r.Level+1, r.Descendents)
		CheckError(err)
	}

	if r.Ordinance != "" {
		_, err := fmt.Fprintf(w, "%s%d ORDI %s\n", indent(r.Level+1), r.Level+1, r.Ordinance)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM link to a submitter record
func (r *SubmitterLink) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d %s @%s@\n", indent(r.Level), r.Level, r.Tag, r.Submitter.Xref)
	CheckError(err)

	return 0, nil
}

// Write formats and writes a GEDCOM level 0 submitter record
func (r *SubmitterRecord) Write(w io.Writer) (int, error) {

	var sXref0, sXrefn string

	if r.Level == 0 {
		sXrefn, sXref0 = "", fmt.Sprintf("@%s@ ", r.Xref)
	} else {
		sXref0, sXrefn = "", fmt.Sprintf(" @%s@", r.Xref)
	}
	_, err := fmt.Fprintf(w, "%s%d %sSUBM%s\n", indent(r.Level), r.Level, sXref0, sXrefn)
	CheckError(err)

	if r.Name != "" {
		_, err := fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level+1), r.Level+1, r.Name)
		CheckError(err)
	}

	if r.Address != nil {
		_, err := r.Address.Write(w)
		CheckError(err)
	}

	if r.Country != "" {
		_, err := fmt.Fprintf(w, "%s%d CTRY %s\n", indent(r.Level+1), r.Level+1, r.Country)
		CheckError(err)
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			_, err := fmt.Fprintf(w, "%s%d PHON %s\n", indent(r.Level+1), r.Level+1, phone)
			CheckError(err)
		}
	}

	if r.Email != "" {
		_, err := fmt.Fprintf(w, "%s%d EMAIL %s\n", indent(r.Level+1), r.Level+1, r.Email)
		CheckError(err)
	}

	if r.WebSite != "" {
		_, err := fmt.Fprintf(w, "%s%d WWW %s\n", indent(r.Level+1), r.Level+1, r.WebSite)
		CheckError(err)
	}

	if r.Language != "" {
		_, err := fmt.Fprintf(w, "%s%d LANG %s\n", indent(r.Level+1), r.Level+1, r.Language)
		CheckError(err)
	}

	if r.Media != nil {
		_, err := r.Media.Write(w)
		CheckError(err)
	}

	if r.RecordFileNumber != "" {
		_, err := fmt.Fprintf(w, "%s%d RFN %s\n", indent(r.Level+1), r.Level+1, r.RecordFileNumber)
		CheckError(err)
	}

	if r.STAL != "" {
		_, err := fmt.Fprintf(w, "%s%d STAL %s\n", indent(r.Level+1), r.Level+1, r.STAL)
		CheckError(err)
	}

	if r.NUMB != "" {
		_, err := fmt.Fprintf(w, "%s%d STAL %s\n", indent(r.Level+1), r.Level+1, r.NUMB)
		CheckError(err)
	}

	if r.RIN != "" {
		_, err := fmt.Fprintf(w, "%s%d RIN %s\n", indent(r.Level+1), r.Level+1, r.RIN)
		CheckError(err)
	}

	if r.Change != nil {
		_, err := r.Change.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a slice of submitter records
func (r SubmitterRecords) Write(w io.Writer) (int, error) {

	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}
	return 0, nil
}

// Write formats and writes a GEDCOM system record
func (r *SystemRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d SOUR %s\n", indent(r.Level), r.Level, r.SystemName)
	CheckError(err)

	if r.ProductName != "" {
		_, err := fmt.Fprintf(w, "%s%d NAME %s\n", indent(r.Level+1), r.Level+1, r.ProductName)
		CheckError(err)
	}

	if r.Version != "" {
		_, err := fmt.Fprintf(w, "%s%d VERS %s\n", indent(r.Level+1), r.Level+1, r.Version)
		CheckError(err)
	}

	if r.Business != nil {
		_, err := r.Business.Write(w)
		CheckError(err)
	}

	if r.SourceData != nil {
		_, err := r.SourceData.Write(w)
		CheckError(err)
	}

	return 0, nil
}

// Write formats and writes a GEDCOM level 0 trailer record
func (r *TrailerRecord) Write(w io.Writer) (int, error) {
	_, err := fmt.Fprintf(w, "%s%d TRLR\n", indent(0), 0)
	CheckError(err)

	return 0, nil
}

// Write formats and writes a GEDCOM user reference number record
func (r *UserReferenceNumberRecord) Write(w io.Writer) (int, error) {

	_, err := fmt.Fprintf(w, "%s%d REFN %s\n", indent(r.Level), r.Level, r.UserReferenceNumber)
	CheckError(err)

	if r.Type != "" {
		_, err := fmt.Fprintf(w, "%s%d TYPE %s\n", indent(r.Level+1), r.Level+1, r.Type)
		CheckError(err)
	}
	return 0, nil
}

// Write formats and writes a slice of user reference number records
func (r UserReferenceNumberRecords) Write(w io.Writer) (int, error) {

	for _, x := range r {
		_, err := x.Write(w)
		CheckError(err)
	}

	return 0, nil
}

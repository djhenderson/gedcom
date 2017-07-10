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
func LongWrite(w io.Writer, level int, xref string, tag string, longString string) (nbytes int, err error) {

	sXref0, sXrefN := "", ""
	if xref != "" {
		sXref0 = fmt.Sprintf(" @%s@", xref)
		if level != 0 {
			sXrefN, sXref0 = sXref0, sXrefN
		}
	}

	if longString == "" {
		nbytes, err := fmt.Fprintf(w, "%s%d%s %s%s\n", indent(level), level, sXref0, tag, sXrefN)
		CheckError(err)

		return nbytes, err
	}

	parts := strings.Split(longString, "\n")
	for i, part := range parts {
		n, err := fmt.Fprintf(w, "%s%d%s %s%s %s\n", indent(level), level, sXref0, tag, sXrefN, part)
		CheckError(err)
		nbytes += n

		if i == 0 {
			tag = "CONT"
			level++
			sXref0, sXrefN = "", ""
		}
	}
	return nbytes, err
}

// WriteLine0 writes a level 0 line
func WriteLine0(w io.Writer, level int, xref string, tag string, value string) (n int, err error) {

	spacer := " "
	if value == "" {
		spacer = ""
	}
	n, err = fmt.Fprintf(w, "%s%d @%s@ %s%s%s\n", indent(level), level, xref, tag, spacer, value)
	CheckError(err)

	return n, err
}

// WriteLineLink writes a link line
func WriteLineLink(w io.Writer, level int, tag string, xref string) (n int, err error) {

	sXref := ""
	if xref != "" {
		sXref = fmt.Sprintf(" @%s@", xref)
	}
	n, err = fmt.Fprintf(w, "%s%d %s%s\n", indent(level), level, tag, sXref)
	CheckError(err)

	return n, err
}

// WriteLineN writes a level n line
func WriteLineN(w io.Writer, level int, tag string, value string) (n int, err error) {

	spacer := " "
	if value == "" {
		spacer = ""
	}
	n, err = fmt.Fprintf(w, "%s%d %s%s%s\n", indent(level), level, tag, spacer, value)
	CheckError(err)

	return n, err
}

// WriteLineNp1 writes a level n+1 line
func WriteLineNp1(w io.Writer, level int, tag string, value string) (n int, err error) {

	spacer := " "
	if value == "" {
		spacer = ""
	}
	n, err = fmt.Fprintf(w, "%s%d %s%s%s\n", indent(level+1), level+1, tag, spacer, value)
	CheckError(err)

	return n, err
}

// Write formats and writes a GEDCOM address record
func (r *AddressRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, "", "ADDR", r.Full)
	nbytes += n

	if r.Name_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_NAME", r.Name_)
		nbytes += n
	}

	if r.Line1 != "" {
		n, err = WriteLineNp1(w, r.Level, "ADR1", r.Line1)
		nbytes += n
	}

	if r.Line2 != "" {
		n, err = WriteLineNp1(w, r.Level, "ADR2", r.Line2)
		nbytes += n
	}

	if r.Line3 != "" {
		n, err = WriteLineNp1(w, r.Level, "ADR3", r.Line3)
		nbytes += n
	}

	if r.City != "" {
		n, err = WriteLineNp1(w, r.Level, "CITY", r.City)
		nbytes += n
	}

	if r.State != "" {
		n, err = WriteLineNp1(w, r.Level, "STAE", r.State)
		nbytes += n
	}

	if r.PostalCode != "" {
		n, err = WriteLineNp1(w, r.Level, "POST", r.PostalCode)
		nbytes += n
	}

	if r.Country != "" {
		n, err = WriteLineNp1(w, r.Level, "CTRY", r.Country)
		nbytes += n
	}

	if r.Phone != "" {
		n, err = WriteLineNp1(w, r.Level, "PHON", r.Phone)
		nbytes += n
	}

	if r.Note != nil { // Leg8
		nbytes += n
		n, err = r.Note.Write(w)
	}

	return nbytes, err
}

// Write formats and writes a slice of address records
func (r AddressRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM author record
func (r *AuthorRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, "", "AUTH", r.Author)
	nbytes += n

	if r.Abbreviation != "" {
		n, err = WriteLineNp1(w, r.Level, "ABBR", r.Abbreviation)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM bibliography record
func (r *BibliographyRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "BIBL", r.Value)
	nbytes += n

	for _, comp := range r.Component {
		n, err = LongWrite(w, r.Level+1, "", "COMP", comp)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM blob record
func (r *BlobRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "BLOB", "")
	nbytes += n

	parts := strings.Split(r.Data, "\n")
	for _, data := range parts {
		n, err = WriteLineNp1(w, r.Level, "CONT", data)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM business record
func (r *BusinessRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "CORP", r.BusinessName)
	nbytes += n

	if r.Address != nil {
		n, err = r.Address.Write(w)
		nbytes += n
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			n, err = WriteLineNp1(w, r.Level, "PHON", phone)
			nbytes += n
		}
	}

	if r.WebSite != "" {
		n, err = WriteLineNp1(w, r.Level, "WWW", r.WebSite)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM call number record
func (r *CallNumberRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "CALN", r.CallNumber)
	nbytes += n

	// In this unique case, the value is a string, not a media record or link
	if r.Media != "" {
		n, err = WriteLineNp1(w, r.Level, "MEDI", r.Media)
		nbytes += n
	}
	return nbytes, err
}

// Write formats and writes a GEDCOM change record
func (r *ChangeRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "CHAN", "")
	nbytes += n

	if r.Date != nil {
		n, err = r.Date.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		nbytes += n
		n, err = r.Note.Write(w)
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM character set record
func (r *CharacterSetRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "CHAR", r.CharacterSet)
	nbytes += n

	if r.Version != "" {
		n, err = WriteLineNp1(w, r.Level, "VERS", r.Version)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM child status record
func (r *ChildStatusRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLine0(w, r.Level, r.Xref, "CSTA", "")
	nbytes += n

	n, err = WriteLineNp1(w, r.Level, "NAME", r.Name)
	nbytes += n

	return nbytes, err
}

// Write formats and writes a slice of child status records
func (r ChildStatusRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	// log.Printf("ChildStatusRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM citatio record
func (r *CitationRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	xref := ""
	if r.Source != nil {
		xref = r.Source.Xref
	}
	n, err = LongWrite(w, r.Level, xref, "SOUR", r.Value)
	nbytes += n

	if r.ReferenceNumber != "" {
		n, err = WriteLineNp1(w, r.Level, "REFN", r.ReferenceNumber)
		nbytes += n
	}

	if r.Page != "" {
		n, err = WriteLineNp1(w, r.Level, "PAGE", r.Page)
		nbytes += n
	}

	if r.Reference != "" {
		n, err = WriteLineNp1(w, r.Level, "REF", r.Reference)
		nbytes += n
	}

	if r.FamilySearchFTID_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_FSFTID", r.FamilySearchFTID_)
		nbytes += n
	}

	if r.Event != nil {
		r.Event.Write(w)
		nbytes += n
	}

	if r.Data != nil {
		for _, data := range r.Data {
			n, err = data.Write(w)
			nbytes += n
		}
	}

	if r.Text != nil {
		for _, text := range r.Text {
			n, err = LongWrite(w, r.Level+1, "", "TEXT", text)
			nbytes += n
		}
	}

	if r.Quality != "" {
		n, err = WriteLineNp1(w, r.Level, "QUAY", r.Quality)
		nbytes += n
	}

	if r.Media != nil {
		n, err = r.Media.Write(w)
		nbytes += n
	}

	if r.CONS != "" {
		n, err = WriteLineNp1(w, r.Level, "CONS", r.CONS)
		nbytes += n
	}

	if r.Direct != "" {
		n, err = WriteLineNp1(w, r.Level, "DIRE", r.Direct)
		nbytes += n
	}

	if r.SourceQuality != "" {
		n, err = WriteLineNp1(w, r.Level, "SOQU", r.SourceQuality)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Date != "" { // Leg8
		n, err = WriteLineNp1(w, r.Level, "DATE", r.Date)
		nbytes += n
	}

	if r.Rin_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_RIN", r.Rin_)
		nbytes += n
	}

	if r.AppliesTo_ != "" { // AQ15
		n, err = WriteLineNp1(w, r.Level, "_APPLIES_TO", r.AppliesTo_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of citation records
func (r CitationRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("CitationRecords type(r): %T\n", r)
	for _, note := range r {
		//log.Printf("CitationRecords type(note): %T\n", note)
		//log.Printf("CitationRecords type(*note): %T\n", *note)
		n, err = note.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM data record
func (r *DataRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, "", "DATA", r.Data)
	nbytes += n

	if r.Date != "" {
		n, err = WriteLineNp1(w, r.Level, "DATE", r.Date)
		nbytes += n
	}

	if r.Copyright != "" {
		n, err = WriteLineNp1(w, r.Level, "COPR", r.Copyright)
		nbytes += n
	}

	if r.Text != nil {
		for _, text := range r.Text {
			n, err = LongWrite(w, r.Level+1, "", "TEXT", text)
			nbytes += n
		}
	}

	return nbytes, err
}

// Write formats and writes a slice of data records
func (r DataRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("DataRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM date record
func (r *DateRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, r.Tag, r.Date)
	nbytes += n

	if r.Day != "" {
		n, err = WriteLineNp1(w, r.Level, "DATD", r.Day)
		nbytes += n
	}

	if r.Month != "" {
		n, err = WriteLineNp1(w, r.Level, "DATM", r.Month)
		nbytes += n
	}

	if r.Year != "" {
		n, err = WriteLineNp1(w, r.Level, "DATY", r.Year)
		nbytes += n
	}

	if r.Full != "" {
		n, err = WriteLineNp1(w, r.Level, "DATF", r.Full)
		nbytes += n
	}

	if r.Short != "" {
		n, err = WriteLineNp1(w, r.Level, "DATS", r.Short)
		nbytes += n
	}

	if r.Time != "" {
		n, err = WriteLineNp1(w, r.Level, "TIME", r.Time)
		nbytes += n
	}

	if r.Text != nil {
		for _, text := range r.Text {
			n, err = LongWrite(w, r.Level+1, "", "TEXT", text)
			nbytes += n
		}
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM event definition record
func (r EventDefinitionRecord) Write(w io.Writer) (nbytes int, err error) {
	log.Fatal("EventDefinition.Write() not implemented")
	return nbytes, err
}

// Write formats and writes a slice of event definition records
func (r EventDefinitionRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("EventDefinitionRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM event records
func (r *EventRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}
	spacer := " "
	if r.Value == "" {
		spacer = ""
	}
	n, err = fmt.Fprintf(w, "%s%d %s%s%s%s\n", indent(r.Level), r.Level, id, r.Tag, spacer, r.Value)
	CheckError(err)
	nbytes += n

	if r.Type != "" {
		n, err = WriteLineNp1(w, r.Level, "TYPE", r.Type)
		nbytes += n
	}

	if r.Name != "" {
		n, err = WriteLineNp1(w, r.Level, "NAME", r.Name)
		nbytes += n
	}

	if r.Primary_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_PRIM", r.Primary_)
		nbytes += n
	}

	if r.AlternateBirth_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_ALT_BIRTH", r.AlternateBirth_)
		nbytes += n
	}

	if r.Confidential_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_CONFIDENTIAL", r.Confidential_)
		nbytes += n
	}

	if r.Date != nil {
		n, err = r.Date.Write(w)
		nbytes += n
	}

	if r.Date2_ != nil { // AQ14
		n, err = r.Date2_.Write(w)
		nbytes += n
	}

	if r.Place != nil {
		n, err = r.Place.Write(w)
		nbytes += n
	}

	if r.Place2_ != nil { // AQ14
		n, err = r.Place2_.Write(w)
		nbytes += n
	}

	if r.Description2_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_Description2", r.Description2_)
		nbytes += n
	}

	if r.Role != nil {
		n, err = r.Role.Write(w)
		nbytes += n
	}
	if r.Address != nil {
		n, err = r.Address.Write(w)
		nbytes += n
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			n, err = WriteLineNp1(w, r.Level, "PHON", phone)
			nbytes += n
		}
	}

	if r.Parents != nil {
		n, err = r.Parents.Write(w)
		nbytes += n
	}

	if r.Husband != nil {
		n, err = r.Husband.Write(w)
		nbytes += n
	}

	if r.Wife != nil {
		n, err = r.Wife.Write(w)
		nbytes += n
	}

	if r.Spouse != nil {
		n, err = r.Spouse.Write(w)
		nbytes += n
	}

	if r.Agency != "" {
		n, err = WriteLineNp1(w, r.Level, "AGNC", r.Agency)
		nbytes += n
	}

	if r.Cause != "" {
		n, err = WriteLineNp1(w, r.Level, "CAUS", r.Cause)
		nbytes += n
	}

	if r.Temple != "" {
		n, err = WriteLineNp1(w, r.Level, "TEMP", r.Temple)
		nbytes += n
	}

	if r.Status != "" {
		n, err = WriteLineNp1(w, r.Level, "STAT", r.Status)
		nbytes += n
	}

	if r.Media != nil {
		n, err = r.Media.Write(w)
		nbytes += n
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	if r.UniqueId_ != nil {
		for _, uid := range r.UniqueId_ {
			n, err = WriteLineNp1(w, r.Level, "_UID", uid)
			nbytes += n
		}
	}

	if r.UpdateTime_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_UPD", r.UpdateTime_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of event records
func (r EventRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("EventRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM link to a family record
func (r *FamilyLink) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineLink(w, r.Level, r.Tag, r.Family.Xref)
	nbytes += n

	if r.Adopted != "" {
		n, err = WriteLineNp1(w, r.Level, "ADOP", r.Adopted)
		nbytes += n
	}

	if r.Primary_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_PRIMARY", r.Primary_)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Pedigree != nil {
		n, err = r.Pedigree.Write(w)
		nbytes += n
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of links to family records
func (r FamilyLinks) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("FamilyLinks type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM family record
func (r *FamilyRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLine0(w, r.Level, r.Xref, "FAM", "")
	nbytes += n

	if r.Status_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_STAT", r.Status_)
		nbytes += n
	}

	if r.NoChildren_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_NONE", r.NoChildren_)
		nbytes += n
	}

	if r.Husband != nil {
		n, err = r.Husband.Write(w)
		nbytes += n
	}

	if r.Wife != nil {
		n, err = r.Wife.Write(w)
		nbytes += n
	}

	if r.NumChildren > 0 {
		value := fmt.Sprintf("%d", r.NumChildren)
		n, err = WriteLineNp1(w, r.Level, "NCHI", value)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	if r.Child != nil {
		n, err = r.Child.Write(w)
		nbytes += n
	}

	if r.Event != nil {
		n, err = r.Event.Write(w)
		nbytes += n
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Submitter != nil {
		r.Submitter.Write(w)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	if r.UniqueId_ != nil {
		for _, uid := range r.UniqueId_ {
			n, err = WriteLineNp1(w, r.Level, "_UID", uid)
			nbytes += n
		}
	}

	if r.UpdateTime_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_UPD", r.UpdateTime_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of family records
func (r FamilyRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("FamilyRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM footnote record
func (r *FootnoteRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "FOOT", r.Value)
	nbytes += n

	for _, comp := range r.Component {
		n, err = LongWrite(w, r.Level+1, "", "COMP", comp)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM gedcom record
func (r *GedcomRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "GEDC", "")
	nbytes += n

	n, err = WriteLineNp1(w, r.Level, "VERS", r.Version)
	nbytes += n

	n, err = WriteLineNp1(w, r.Level, "FORM", r.Form)
	nbytes += n

	return nbytes, err

}

// Write formats and writes a GEDCOM header record
func (r *HeaderRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "HEAD", "")
	nbytes += n

	if r.CharacterSet != nil {
		n, err = r.CharacterSet.Write(w)
		nbytes += n
	}

	if r.SourceSystem != nil {
		n, err = r.SourceSystem.Write(w)
		nbytes += n
	}

	if r.Destination != "" {
		n, err = WriteLineNp1(w, r.Level, "DEST", r.Destination)
		nbytes += n
	}

	if r.Date != nil {
		n, err = r.Date.Write(w)
		nbytes += n
	}

	if r.FileName != "" {
		n, err = WriteLineNp1(w, r.Level, "FILE", r.FileName)
		nbytes += n
	}

	if r.Gedcom != nil {
		n, err = r.Gedcom.Write(w)
		nbytes += n
	}

	if r.Language != "" {
		n, err = WriteLineNp1(w, r.Level, "LANG", r.Language)
		nbytes += n
	}

	if r.Copyright != "" {
		n, err = WriteLineNp1(w, r.Level, "COPR", r.Copyright)
		nbytes += n
	}

	if r.Place != nil {
		n, err = r.Place.Write(w)
		nbytes += n
	}

	if r.RootPerson_ != nil {
		n, err = r.RootPerson_.Write(w)
		nbytes += n
	}

	if r.HomePerson_ != nil {
		n, err = r.HomePerson_.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Submitter != nil {
		n, err = r.Submitter.Write(w)
		nbytes += n
	}

	if r.Submission != nil {
		r.Submission.Write(w)
		nbytes += n
	}

	if r.Schema != nil {
		n, err = r.Schema.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM history record
func (r *HistoryRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "HIST", r.History)
	nbytes += n

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of history records
func (r HistoryRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM individual links
func (r *IndividualLink) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineLink(w, r.Level, r.Tag, r.Individual.Xref)
	nbytes += n

	if r.Relationship != "" {
		n, err = WriteLineNp1(w, r.Level, "RELA", r.Relationship)
		nbytes += n
	}

	if r.Event != nil {
		n, err = r.Event.Write(w)
		nbytes += n
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Age != "" {
		n, err = WriteLineNp1(w, r.Level, "AGE", r.Age)
		nbytes += n
	}

	if r.Preferred_ != "" { // Leg8
		n, err = WriteLineNp1(w, r.Level, "_PREF", r.Preferred_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of links to individual records
func (r IndividualLinks) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("IndividualLinks type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM individual record
func (r *IndividualRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLine0(w, r.Level, r.Xref, "INDI", "")
	nbytes += n

	for _, name := range r.Name {
		n, err = name.Write(w)
		nbytes += n
	}

	if r.Status_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_STAT", r.Status_)
		nbytes += n
	}

	if r.Restriction != "" {
		n, err = WriteLineNp1(w, r.Level, "RESN", r.Restriction)
		nbytes += n
	}

	if r.Sex != "" {
		n, err = WriteLineNp1(w, r.Level, "SEX", r.Sex)
		nbytes += n
	}

	if r.ProfilePicture_ != nil {
		n, err = r.ProfilePicture_.Write(w)
		nbytes += n
	}

	if r.CONL != "" {
		n, err = WriteLineNp1(w, r.Level, "CONL", r.CONL)
		nbytes += n
	}

	if r.Event != nil {
		r.Event.Write(w)
		nbytes += n
	}

	if r.Attribute != "" {
		n, err = WriteLineNp1(w, r.Level, "ATTR", r.Attribute)
		nbytes += n
	}

	if r.Parents != nil {
		n, err = r.Parents.Write(w)
		nbytes += n
	}

	if r.Family != nil {
		n, err = r.Family.Write(w)
		nbytes += n
	}

	if r.Address != nil {
		n, err = r.Address.Write(w)
		nbytes += n
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			n, err = WriteLineNp1(w, r.Level, "PHON", phone)
			nbytes += n
		}
	}

	if r.Health != "" {
		n, err = WriteLineNp1(w, r.Level, "HEAL", r.Health)
		nbytes += n
	}

	if r.History != nil {
		n, err = r.History.Write(w)
		nbytes += n
	}

	if r.Quality != "" {
		n, err = WriteLineNp1(w, r.Level, "QUAY", r.Quality)
		nbytes += n
	}

	if r.Living != "" {
		n, err = WriteLineNp1(w, r.Level, "LVG", r.Living)
		nbytes += n
	}

	if r.AncestralFileNumber != nil {
		for _, afn := range r.AncestralFileNumber {
			n, err = WriteLineNp1(w, r.Level, "AFN", afn)
			nbytes += n
		}
	}

	if r.RecordFileNumber != "" {
		n, err = WriteLineNp1(w, r.Level, "RFN", r.RecordFileNumber)
		nbytes += n
	}

	if r.UserReferenceNumber != nil {
		n, err = r.UserReferenceNumber.Write(w)
		nbytes += n
	}

	if r.FamilySearchFTID_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_FSFTID", r.FamilySearchFTID_)
		nbytes += n
	}

	if r.FamilySearchLink_ != "" { // Leg8
		n, err = WriteLineNp1(w, r.Level, "_FSLINK", r.FamilySearchLink_)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	if r.UniqueId_ != nil {
		for _, uid := range r.UniqueId_ {
			n, err = WriteLineNp1(w, r.Level, "_UID", uid)
			nbytes += n
		}
	}

	if r.Email != "" {
		n, err = WriteLineNp1(w, r.Level, "EMAIL", r.Email)
		nbytes += n
	}

	if r.Email_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_EMAIL", r.Email_)
		nbytes += n
	}

	if r.URL_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_URL", r.URL_)
		nbytes += n
	}

	if r.WebSite != "" {
		n, err = WriteLineNp1(w, r.Level, "WWW", r.WebSite)
		nbytes += n
	}

	if r.Media != nil {
		n, err = r.Media.Write(w)
		nbytes += n
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Associated != nil {
		n, err = r.Associated.Write(w)
		nbytes += n
	}

	if r.Submitter != nil {
		n, err = r.Submitter.Write(w)
		nbytes += n
	}

	if r.ANCI != nil {
		n, err = r.ANCI.Write(w)
		nbytes += n
	}

	if r.DESI != nil {
		n, err = r.DESI.Write(w)
		nbytes += n
	}

	if r.Alias != "" {
		n, err = WriteLineNp1(w, r.Level, "ALIA", r.Alias)
		nbytes += n
	}

	if r.Father != nil {
		n, err = r.Father.Write(w)
		nbytes += n
	}

	if r.Mother != nil {
		n, err = r.Mother.Write(w)
		nbytes += n
	}

	if r.Miscellaneous != nil {
		for _, misc := range r.Miscellaneous {
			n, err = WriteLineNp1(w, r.Level, "MISC", misc)
			nbytes += n
		}
	}

	if r.PPExclude_ != "" { // Leg8
		n, err = WriteLineNp1(w, r.Level, "_PPEXCLUDE", r.PPExclude_)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	if r.Todo_ != "" { // AQ15
		n, err = WriteLineNp1(w, r.Level, "_TODO", r.Todo_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of individual records
func (r IndividualRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("IndividualRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM media links
func (r *MediaLink) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineLink(w, r.Level, r.Tag, r.Media.Xref)
	nbytes += n

	return nbytes, err
}

// Write formats and writes a slice of links to media records
func (r MediaLinks) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("MediaLinks type(r): %T\n", r)
	for _, x := range r {
		if x.Value != "" {
			n, err = x.Write(w)
		} else {
			n, err = x.Media.Write(w)
		}
		nbytes += n

	}

	return nbytes, err
}

// Write formats and writes a GEDCOM media record
func (r *MediaRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	id0, idN := "", ""
	if r.Xref != "" {
		id0 = fmt.Sprintf(" @%s@ ", r.Xref)

		if r.Level != 0 {
			idN, id0 = id0, idN
		}
	}

	n, err = fmt.Fprintf(w, "%s%d%s %s%s\n", indent(r.Level), r.Level, id0, "OBJE", idN)
	CheckError(err)
	nbytes += n

	if r.Format != "" {
		n, err = WriteLineNp1(w, r.Level, "FORM", r.Format)
		nbytes += n
	}

	if r.Url_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_URL", r.Url_)
		nbytes += n
	}

	if r.FileName != "" {
		n, err = WriteLineNp1(w, r.Level, "FILE", r.FileName)
		nbytes += n
	}

	if r.Title != "" {
		n, err = WriteLineNp1(w, r.Level, "TITL", r.Title)
		nbytes += n
	}

	if r.Date != "" {
		n, err = WriteLineNp1(w, r.Level, "DATE", r.Date)
		nbytes += n
	}

	if r.Author != "" {
		n, err = WriteLineNp1(w, r.Level, "AUTH", r.Author)
		nbytes += n
	}

	if r.Text != "" {
		n, err = LongWrite(w, r.Level, "", "TEXT", r.Text)
		nbytes += n
	}

	if r.Date_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_DATE", r.Date_)
		nbytes += n
	}

	if r.AstId_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ASTID", r.AstId_)
		nbytes += n
	}

	if r.AstType_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ASTTYP", r.AstType_)
		nbytes += n
	}

	if r.AstDesc_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ASTDESC", r.AstDesc_)
		nbytes += n
	}

	if r.AstLoc_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ASTLOC", r.AstLoc_)
		nbytes += n
	}

	if r.AstPerm_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ASTPERM", r.AstPerm_)
		nbytes += n
	}

	if r.AstUpPid_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ASTUPPID", r.AstUpPid_)
		nbytes += n
	}

	if r.BinaryLargeObject != nil {
		n, err = r.BinaryLargeObject.Write(w)
		nbytes += n
	}

	if r.UserReferenceNumber != nil {
		n, err = r.UserReferenceNumber.Write(w)
		nbytes += n
	}

	if r.Rin != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.Rin)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	if r.FsFtId_ != "" { // AQ15
		n, err = WriteLineNp1(w, r.Level, "_FSFTID", r.FsFtId_)
		nbytes += n
	}

	if r.Scbk_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_SCBK", r.Scbk_)
		nbytes += n
	}

	if r.Primary_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_PRIM", r.Primary_)
		nbytes += n
	}

	if r.Type_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_TYPE", r.Type_)
		nbytes += n
	}

	if r.Sshow_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_SSHOW", r.Sshow_)
		nbytes += n
	}

	if r.Stime_ != "" { // AQ15
		n, err = WriteLineNp1(w, r.Level, "_STIME", r.Stime_)
		nbytes += n
	}

	if r.mediaLinks != nil {
		n, err = r.mediaLinks.Write(w)
		nbytes += n
	}

	if r.SrcPp_ != "" { // AQ15
		n, err = WriteLineNp1(w, r.Level, "_SRCPP", r.SrcPp_)
		nbytes += n
	}

	if r.SrcFlip_ != "" { // AQ15
		n, err = WriteLineNp1(w, r.Level, "_SRCFLIP", r.SrcFlip_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of media records
func (r MediaRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("MediaRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM name record
func (r *NameRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "NAME", r.Name)
	nbytes += n

	if r.Prefix != "" {
		n, err = WriteLineNp1(w, r.Level, "NPFX", r.Prefix)
		nbytes += n
	}

	if r.GivenName != "" {
		n, err = WriteLineNp1(w, r.Level, "GIVN", r.GivenName)
		nbytes += n
	}

	if r.MiddleName_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_MIDN", r.MiddleName_)
		nbytes += n
	}

	if r.SurnamePrefix != "" {
		n, err = WriteLineNp1(w, r.Level, "SPFX", r.SurnamePrefix)
		nbytes += n
	}

	if r.Surname != "" {
		n, err = WriteLineNp1(w, r.Level, "SURN", r.Surname)
		nbytes += n
	}

	if r.Suffix != "" {
		n, err = WriteLineNp1(w, r.Level, "NSFX", r.Suffix)
		nbytes += n
	}

	if r.PreferedGivenName_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_PGVN", r.PreferedGivenName_)
		nbytes += n
	}

	if r.RomanizedName != "" {
		n, err = WriteLineNp1(w, r.Level, "ROMN", r.RomanizedName)
		nbytes += n
	}

	if r.PhoneticName != "" {
		n, err = WriteLineNp1(w, r.Level, "FONE", r.PhoneticName)
		nbytes += n
	}

	if r.MarriedName_ != "" { // AQ14
		n, err = WriteLineNp1(w, r.Level, "_MARNM", r.MarriedName_)
		nbytes += n
	}

	if r.NameType != "" {
		n, err = WriteLineNp1(w, r.Level, "TYPE", r.NameType)
		nbytes += n
	}

	if r.Primary_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_PRIM", r.Primary_)
		nbytes += n
	}

	if r.AlsoKnownAs_ != nil {
		for _, aka := range r.AlsoKnownAs_ {
			n, err = WriteLineNp1(w, r.Level, "_AKA", aka)
			nbytes += n
		}
	}

	if r.Nickname != nil {
		for _, nick := range r.Nickname {
			n, err = WriteLineNp1(w, r.Level, "NICK", nick)
			nbytes += n
		}
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		r.Note.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of name records
func (r NameRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("NameRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM note records
func (r *NoteRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, r.Xref, "NOTE", r.Note)
	nbytes += n

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.UserReferenceNumber != nil {
		n, err = r.UserReferenceNumber.Write(w)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of note records
func (r NoteRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("NoteRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM pedigree record
func (r *PedigreeRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "PEDI", r.Pedigree)
	nbytes += n

	if r.Husband_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_HUSB", r.Husband_)
		nbytes += n
	}

	if r.Wife_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_WIFE", r.Wife_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM place part record
func (r *PlacePartRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, r.Tag, r.Part)
	nbytes += n

	if r.Jurisdiction != "" {
		n, err = WriteLineNp1(w, r.Level, "JURI", r.Jurisdiction)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM place record
func (r *PlaceRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	id := ""
	if r.Xref != "" {
		id = fmt.Sprintf("@%s@ ", r.Xref)
	}

	n, err = fmt.Fprintf(w, "%s%d %s%s %s\n", indent(r.Level), r.Level, id, r.Tag, r.Name)
	CheckError(err)
	nbytes += n

	if r.Form != "" {
		n, err = WriteLineNp1(w, r.Level, "PLAS", r.Form)
		nbytes += n
	}

	if r.ShortName != "" {
		n, err = WriteLineNp1(w, r.Level, "PLAS", r.ShortName)
		nbytes += n
	}

	if r.Modifier != "" {
		n, err = WriteLineNp1(w, r.Level, "PLAM", r.Modifier)
		nbytes += n
	}

	if r.Parts != nil {
		for _, part := range r.Parts {
			n, err = part.Write(w)
			nbytes += n
		}
	}

	if r.Citation != nil {
		n, err = r.Citation.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of place records
func (r PlaceRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("PlaceRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM link to a repository record
func (r *RepositoryLink) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineLink(w, r.Level, "REPO", r.Repository.Xref)
	nbytes += n

	if r.CallNumber != nil {
		n, err = r.CallNumber.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of repository links
func (r RepositoryLinks) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("RepositoryLinks type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM repository record
func (r *RepositoryRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLine0(w, r.Level, r.Xref, "REPO", "")
	nbytes += n

	if r.Name != "" {
		n, err = WriteLineNp1(w, r.Level, "NAME", r.Name)
		nbytes += n
	}

	if r.Address != nil {
		n, err = r.Address.Write(w)
		nbytes += n
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			n, err = WriteLineNp1(w, r.Level, "PHON", phone)
			nbytes += n
		}
	}

	if r.WebSite != "" {
		n, err = WriteLineNp1(w, r.Level, "WWW", r.WebSite)
		nbytes += n
	}

	if r.UserReferenceNumber != nil {
		n, err = r.UserReferenceNumber.Write(w)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of repository records
func (r RepositoryRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("RepositoryRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM role record
func (r *RoleRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	spacer := ""
	xref := ""
	if r.Individual != nil {
		if r.Individual.Xref != "" {
			xref = fmt.Sprintf("@%s@", r.Individual.Xref)
		}
	}
	if r.Role != "" && xref != "" {
		spacer = " "
	}
	n, err = fmt.Fprintf(w, "%s%d ROLE %s%s%s\n", indent(r.Level), r.Level, r.Role, spacer, xref)
	CheckError(err)
	nbytes += n

	if r.Principal != "" {
		n, err = WriteLineNp1(w, r.Level, "PRIN", r.Principal)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of role record
func (r RoleRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("RoleRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM root record, i.e. the whole file
func (r *RootRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	if r.Header != nil {
		n, err = r.Header.Write(w)
		nbytes += n
	}

	if len(r.Submission) > 0 { // SUBM
		r.Submission.Write(w)
		nbytes += n
	}

	if len(r.Submitter) > 0 {
		r.Submitter.Write(w)
		nbytes += n
	}

	if len(r.Individual) > 0 { // INDI
		r.Individual.Write(w)
		nbytes += n
	}

	if len(r.Family) > 0 { // FAM
		r.Family.Write(w)
		nbytes += n
	}

	if len(r.Note) > 0 { // NOTE
		r.Note.Write(w)
		nbytes += n
	}

	if len(r.Place) > 0 { // PLAC
		r.Place.Write(w)
		nbytes += n
	}

	if len(r.Event) > 0 { // EVEN
		r.Event.Write(w)
		nbytes += n
	}

	if len(r.Media) > 0 { // OBJE
		n, err = r.Media.Write(w)
		nbytes += n
	}

	if len(r.ChildStatus) > 0 { // _CSTA
		n, err = r.ChildStatus.Write(w)
		nbytes += n
	}

	if len(r.Todo_) > 0 { // _TODO
		r.Todo_.Write(w)
		nbytes += n
	}

	if len(r.Source) > 0 { // SOUR
		r.Source.Write(w)
		nbytes += n
	}

	if len(r.Repository) > 0 { // REPO
		r.Repository.Write(w)
		nbytes += n
	}

	if r.Trailer != nil { // TRLR
		n, err = r.Trailer.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM schema record
func (r *SchemaRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = fmt.Fprintf(w, "%s%d SCHEMA\n", indent(r.Level), r.Level)
	CheckError(err)
	nbytes += n

	for _, data := range r.Data {
		level := int(data[0])
		n, err = fmt.Fprintf(w, "%s\n", indent(level)+data)
		CheckError(err)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM short title record
func (r *ShortTitleRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "SHTI", r.ShortTitle)
	nbytes += n

	if r.Indexed != "" {
		n, err = WriteLineNp1(w, r.Level, "INDX", r.Indexed)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM level 0 source record
func (r *SourceRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, r.Xref, "SOUR", r.Value)
	nbytes += n

	if r.Name != "" {
		n, err = WriteLineNp1(w, r.Level, "NAME", r.Name)
		nbytes += n
	}

	if r.Title != "" {
		n, err = LongWrite(w, r.Level+1, "", "TITL", r.Title)
		nbytes += n
	}

	if r.Author != nil {
		n, err = r.Author.Write(w)
		nbytes += n
	}

	if r.Abbreviation != "" {
		n, err = WriteLineNp1(w, r.Level, "ABBR", r.Abbreviation)
		nbytes += n
	}

	if r.Publication != "" {
		n, err = LongWrite(w, r.Level+1, "", "PUBL", r.Publication)
		nbytes += n
	}

	if r.MediaType != "" { // Leg8
		n, err = WriteLineNp1(w, r.Level, "MEDI", r.MediaType)
		nbytes += n
	}

	if r.Parenthesized_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_PAREN", r.Parenthesized_)
		nbytes += n
	}

	if r.Type_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_TYPE", r.Type_)
		nbytes += n
	}

	if r.Other_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_OTHER", r.Other_)
		nbytes += n
	}

	if r.Master_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_MASTER", r.Master_)
		nbytes += n
	}

	if r.Italic_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_ITALIC", r.Italic_)
		nbytes += n
	}

	if r.Text != nil {
		for _, text := range r.Text {
			n, err = LongWrite(w, r.Level+1, "", "TEXT", text)
			nbytes += n
		}
	}

	if r.Data != nil {
		n, err = r.Data.Write(w)
		nbytes += n
	}

	if r.Footnote != nil {
		n, err = r.Footnote.Write(w)
		nbytes += n
	}

	if r.Bibliography != nil {
		n, err = r.Bibliography.Write(w)
		nbytes += n
	}

	if r.Repository != nil {
		n, err = r.Repository.Write(w)
		nbytes += n
	}

	if r.UserReferenceNumber != nil {
		n, err = r.UserReferenceNumber.Write(w)
		nbytes += n
	}

	if r.Quality != "" {
		n, err = WriteLineNp1(w, r.Level, "QUAY", r.Quality)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	if r.ShortAuthor != "" {
		n, err = WriteLineNp1(w, r.Level, "SHAU", r.ShortAuthor)
		nbytes += n
	}

	if r.ShortTitle != nil {
		n, err = r.ShortTitle.Write(w)
		nbytes += n
	}

	if r.Media != nil {
		n, err = r.Media.Write(w)
		nbytes += n
	}

	if r.Note != nil {
		n, err = r.Note.Write(w)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}
	if r.WebTag_ != nil { // RM6
		n, err = r.WebTag_.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of source records
func (r SourceRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM link to a submission record
func (r *SubmissionLink) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineLink(w, r.Level, "SUBN", r.Submission.Xref)
	nbytes += n

	return nbytes, err
}

// Write formats and writes a slice of links to submission records
func (r SubmissionLinks) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM level 0 submission record
func (r *SubmissionRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLine0(w, r.Level, r.Xref, "SUBN", "")
	nbytes += n

	if r.FamilyFileName != "" {
		n, err = WriteLineNp1(w, r.Level, "FAMF", r.FamilyFileName)
		nbytes += n
	}

	if r.Temple != "" {
		n, err = WriteLineNp1(w, r.Level, "TEMP", r.Temple)
		nbytes += n
	}

	if r.Ancestors != "" {
		n, err = WriteLineNp1(w, r.Level, "ANCE", r.Ancestors)
		nbytes += n
	}

	if r.Descendents != "" {
		n, err = WriteLineNp1(w, r.Level, "DESC", r.Descendents)
		nbytes += n
	}

	if r.Ordinance != "" {
		n, err = WriteLineNp1(w, r.Level, "ORDI", r.Ordinance)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of submission record
func (r SubmissionRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	//log.Printf("SubmissionRecords type(r): %T\n", r)
	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM link to a submitter record
func (r *SubmitterLink) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineLink(w, r.Level, r.Tag, r.Submitter.Xref)
	nbytes += n

	return nbytes, err
}

// Write formats and writes a slice of links to submitter records
func (r SubmitterLinks) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}
	return nbytes, err
}

// Write formats and writes a GEDCOM level 0 submitter record
func (r *SubmitterRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	var sXref0, sXrefn string

	if r.Level == 0 {
		sXrefn, sXref0 = "", fmt.Sprintf("@%s@ ", r.Xref)
	} else {
		sXref0, sXrefn = "", fmt.Sprintf(" @%s@", r.Xref)
	}
	n, err = fmt.Fprintf(w, "%s%d %sSUBM%s\n", indent(r.Level), r.Level, sXref0, sXrefn)
	CheckError(err)
	nbytes += n

	if r.Name != "" {
		n, err = WriteLineNp1(w, r.Level, "NAME", r.Name)
		nbytes += n
	}

	if r.Address != nil {
		n, err = r.Address.Write(w)
		nbytes += n
	}

	if r.Country != "" {
		n, err = WriteLineNp1(w, r.Level, "CTRY", r.Country)
		nbytes += n
	}

	if r.Phone != nil {
		for _, phone := range r.Phone {
			n, err = WriteLineNp1(w, r.Level, "PHON", phone)
			nbytes += n
		}
	}

	if r.Email != "" {
		n, err = WriteLineNp1(w, r.Level, "EMAIL", r.Email)
		nbytes += n
	}

	if r.Email_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_EMAIL", r.Email_)
		nbytes += n
	}

	if r.WebSite != "" {
		n, err = WriteLineNp1(w, r.Level, "WWW", r.WebSite)
		nbytes += n
	}

	if r.Language != "" {
		n, err = WriteLineNp1(w, r.Level, "LANG", r.Language)
		nbytes += n
	}

	if r.Media != nil {
		n, err = r.Media.Write(w)
		nbytes += n
	}

	if r.RecordFileNumber != "" {
		n, err = WriteLineNp1(w, r.Level, "RFN", r.RecordFileNumber)
		nbytes += n
	}

	if r.STAL != "" {
		n, err = WriteLineNp1(w, r.Level, "STAL", r.STAL)
		nbytes += n
	}

	if r.NUMB != "" {
		n, err = WriteLineNp1(w, r.Level, "STAL", r.NUMB)
		nbytes += n
	}

	if r.RecordInternal != "" {
		n, err = WriteLineNp1(w, r.Level, "RecordInternal", r.RecordInternal)
		nbytes += n
	}

	if r.Change != nil {
		n, err = r.Change.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of submitter records
func (r SubmitterRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}
	return nbytes, err
}

// Write formats and writes a GEDCOM system record
func (r *SystemRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "SOUR", r.SystemName)
	nbytes += n

	if r.ProductName != "" {
		n, err = WriteLineNp1(w, r.Level, "NAME", r.ProductName)
		nbytes += n
	}

	if r.Version != "" {
		n, err = WriteLineNp1(w, r.Level, "VERS", r.Version)
		nbytes += n
	}

	if r.Business != nil {
		n, err = r.Business.Write(w)
		nbytes += n
	}

	if r.SourceData != nil {
		n, err = r.SourceData.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM level 0 todo record
func (r *TodoRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, r.Xref, "_TODO", r.Value)
	nbytes += n

	if r.Description != "" {
		n, err = WriteLineNp1(w, r.Level, "DESC", r.Description)
		nbytes += n
	}

	if r.Priority_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_PRIORITY", r.Priority_)
		nbytes += n
	}

	if r.Type != "" {
		n, err = WriteLineNp1(w, r.Level, "TYPE", r.Type)
		nbytes += n
	}

	if r.Status != "" {
		n, err = WriteLineNp1(w, r.Level, "STAT", r.Status)
		nbytes += n
	}

	if r.Date2_ != "" {
		n, err = WriteLineNp1(w, r.Level, "_DATE2", r.Date2_)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of todo records
func (r TodoRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM level 0 trailer record
func (r *TrailerRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, 0, "TRLR", "")
	nbytes += n

	return nbytes, err
}

// Write formats and writes a GEDCOM user reference number record
func (r *UserReferenceNumberRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = WriteLineN(w, r.Level, "REFN", r.UserReferenceNumber)
	nbytes += n

	if r.Type != "" {
		n, err = WriteLineNp1(w, r.Level, "TYPE", r.Type)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a slice of user reference number records
func (r UserReferenceNumberRecords) Write(w io.Writer) (nbytes int, err error) {
	var n int

	for _, x := range r {
		n, err = x.Write(w)
		nbytes += n
	}

	return nbytes, err
}

// Write formats and writes a GEDCOM web tag record (RM6)
func (r *WebTagRecord) Write(w io.Writer) (nbytes int, err error) {
	var n int

	n, err = LongWrite(w, r.Level, r.Xref, "_WEBTAG", r.Value)
	nbytes += n

	if r.Name != "" {
		n, err = WriteLineNp1(w, r.Level, "NAME", r.Name)
		nbytes += n
	}

	if r.URL != "" {
		n, err = WriteLineNp1(w, r.Level, "URL", r.URL)
		nbytes += n
	}

	return nbytes, err
}

/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

// Package gedcom supports decoding and encoding GEDCOM files.
package gedcom

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

// A Decoder reads and decodes GEDCOM objects from an input stream.
type Decoder struct {
	r            io.Reader
	parsers      []parser
	refs         map[string]interface{}
	LineNum      int
	warningCount int
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Count warnings
func (d *Decoder) CountWarnings() {
	d.warningCount += 1
	if d.warningCount >= 10 {
		panic("too many warnings.")
	}
}

// Decode reads the next GEDCOM-encoded value from its
// input and stores it in the value pointed to by v.
func (d *Decoder) Decode() (*RootRecord, error) {

	r := &RootRecord{
		Level:            -1,
		Album:            make(AlbumRecords, 0), // MH/FTB8
		ChildStatus:      make(ChildStatusRecords, 0),
		Event:            make(EventRecords, 0),
		EventDefinition_: make(EventDefinitionRecords, 0), // AQ14, RM6
		Family:           make(FamilyRecords, 0),
		Individual:       make(IndividualRecords, 0),
		Media:            make(MediaRecords, 0),
		Note:             make(NoteRecords, 0),
		Place:            make(PlaceRecords, 0),
		PlaceDefinition_: make(PlaceDefinitionRecords, 0), // Leg8
		Repository:       make(RepositoryRecords, 0),
		Source:           make(SourceRecords, 0),
		Submission:       make(SubmissionRecords, 0),
		Submitter:        make(SubmitterRecords, 0),
		Todo_:            make(TodoRecords, 0), // AQ15
	}

	d.refs = make(map[string]interface{})
	d.parsers = []parser{makeRootParser(d, r)}
	err := d.scan(r)
	if err != nil {
		log.Println(err.Error())
	}

	return r, err
}

func (d *Decoder) scan(r *RootRecord) (err error) {
	s := &scanner{}
	buf := make([]byte, 512)

	var n int
	n, err = d.r.Read(buf)
	if err != nil {
		log.Fatalf("Fatal: at %d: %s\n", d.LineNum, err.Error())
		return
	}

	for n > 0 {
		pos := 0

		for {
			s.reset()
			d.LineNum += 1
			var offset int
			offset, err = s.nextTag(buf[pos:n])
			pos += offset
			if err != nil {
				if err == io.EOF {
					err = nil
				} else {
					log.Printf("Error: at %d: %s\n", d.LineNum, err.Error())
				}
				break
			}

			var xref string
			if s.xref != nil && len(s.xref) > 0 {
				xref = "@" + string(s.xref) + "@"
			} else {
				xref = ""
			}
			d.parsers[len(d.parsers)-1](s.level, string(s.tag), string(s.value), xref)

		}

		// shift unparsed bytes to start of buffer
		rest := copy(buf, buf[pos:])

		// top up buffer
		var num int
		num, err = d.r.Read(buf[rest:len(buf)])
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				log.Println(err.Error())
			}
			break
		}

		n = rest + num - 1

	}
	return
}

type parser func(level int, tag string, value string, xref string) error

func (d *Decoder) pushParser(p parser) {
	d.parsers = append(d.parsers, p)
}

func (d *Decoder) popParser(level int, tag string, value string, xref string) error {
	n := len(d.parsers) - 1
	if n < 1 {
		log.Panic("MASSIVE ERROR") // TODO
	}
	d.parsers = d.parsers[0:n]

	return d.parsers[len(d.parsers)-1](level, tag, value, xref)
}

// FindFamily returns the FamilyRecord for an xref
func (d *Decoder) FindFamily(xref string) *FamilyRecord {
	ref, found := d.refs[xref]
	if !found || ref == nil {
		// log.Printf("FindFamily(%s) not found.\n", xref)
		return nil
	}
	f, found := ref.(*FamilyRecord)
	if !found {
		log.Fatalf("FindFamily(%s) not a Family\n", xref)
	}
	return f
}

// FindSource returns the SourceRecord for an xref
func (d *Decoder) FindSource(xref string) *SourceRecord {
	return d.refs[xref].(*SourceRecord)
}

// Level 0 record constructors

// album finds or creates a level 0 AlbumRecord
func (d *Decoder) album(xref string) *AlbumRecord {
	if xref == "" {
		return &AlbumRecord{}
	}

	ref, found := d.refs[xref].(*AlbumRecord)
	if !found {
		rec := &AlbumRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// attribute finds or creates a level 0 AttributeRecord or a link to one
func (d *Decoder) XXattribute(xref string) *AttributeRecord {
	if xref == "" {
		return &AttributeRecord{}
	}

	ref, found := d.refs[xref].(*AttributeRecord)
	if !found {
		rec := &AttributeRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// childStatus finds or creates a level 0 ChildStatusRecord or a link to one
func (d *Decoder) childStatus(xref string) *ChildStatusRecord {
	if xref == "" {
		return &ChildStatusRecord{}
	}

	ref, found := d.refs[xref].(*ChildStatusRecord)
	if !found {
		rec := &ChildStatusRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// event finds or creates a level 0 EventRecord or a link to one
func (d *Decoder) event(xref string) *EventRecord {
	if xref == "" {
		return &EventRecord{}
	}

	ref, found := d.refs[xref].(*EventRecord)
	if !found {
		rec := &EventRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// eventDefinition finds or creates a level 0 EventDefinitionRecord
func (d *Decoder) eventDefinition(xref string) *EventDefinitionRecord {
	if xref == "" {
		return &EventDefinitionRecord{}
	}

	ref, found := d.refs[xref].(*EventDefinitionRecord)
	if !found {
		rec := &EventDefinitionRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// family finds or creates a level 0 FamilyRecord
func (d *Decoder) family(xref string) *FamilyRecord {
	if xref == "" {
		return &FamilyRecord{}
	}

	ref := d.FindFamily(xref)
	if ref == nil {
		rec := &FamilyRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// header finds or creates a level 0 HeaderRecord
func (d *Decoder) header(xref string) *HeaderRecord {
	if xref == "" {
		xref = "HEAD"
	}

	ref, found := d.refs[xref].(*HeaderRecord)
	if !found {
		rec := &HeaderRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// individual finds or creates a level 0 IndividualRecord
func (d *Decoder) individual(xref string) *IndividualRecord {
	if xref == "" {
		return &IndividualRecord{}
	}

	ref, found := d.refs[xref].(*IndividualRecord)
	if !found {
		rec := &IndividualRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// media finds or creates a level 0 MediaRecord
func (d *Decoder) media(xref string) *MediaRecord {
	if xref == "" {
		return &MediaRecord{}
	}

	ref, found := d.refs[xref].(*MediaRecord)
	if !found {
		rec := &MediaRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// note finds or creates a level 0 NoteRecord
func (d *Decoder) note(xref string) *NoteRecord {
	if xref == "" {
		return &NoteRecord{}
	}

	ref, found := d.refs[xref].(*NoteRecord)
	if !found {
		rec := &NoteRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// place finds or creates a level 0 PlaceRecord
func (d *Decoder) place(xref string) *PlaceRecord {
	if xref == "" {
		return &PlaceRecord{Tag: "PLAC"}
	}

	ref, found := d.refs[xref].(*PlaceRecord)
	if !found {
		rec := &PlaceRecord{Xref: xref, Tag: "PLAC"}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// placeDefinition finds or creates a level 0 PlaceDefinitionRecord
func (d *Decoder) placeDefinition(xref string) *PlaceDefinitionRecord {
	if xref == "" {
		return &PlaceDefinitionRecord{}
	}

	ref, found := d.refs[xref].(*PlaceDefinitionRecord)
	if !found {
		rec := &PlaceDefinitionRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// publish finds or creates a level 0 PublishRecord
func (d *Decoder) publish(xref string) *PublishRecord {
	if xref == "" {
		return &PublishRecord{}
	}

	ref, found := d.refs[xref].(*PublishRecord)
	if !found {
		rec := &PublishRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// repository finds or creates a level 0 RepositoryRecord
func (d *Decoder) repository(xref string) *RepositoryRecord {
	if xref == "" {
		return &RepositoryRecord{}
	}

	ref, found := d.refs[xref].(*RepositoryRecord)
	if !found {
		rec := &RepositoryRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// source finds or creates a level 0 SourceRecord
func (d *Decoder) source(xref string) *SourceRecord {
	if xref == "" {
		return &SourceRecord{}
	}

	ref, found := d.refs[xref].(*SourceRecord)
	if !found {
		rec := &SourceRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// submission finds or creates a level 0 SubmissionRecord
func (d *Decoder) submission(xref string) *SubmissionRecord {
	if xref == "" {
		return &SubmissionRecord{}
	}

	ref, found := d.refs[xref].(*SubmissionRecord)
	if !found {
		rec := &SubmissionRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// submitter finds or creates a level 0 SubmitterRecord
func (d *Decoder) submitter(xref string) *SubmitterRecord {
	if xref == "" {
		return &SubmitterRecord{}
	}

	ref, found := d.refs[xref].(*SubmitterRecord)
	if !found {
		rec := &SubmitterRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// todo finds or creates a level 0 ToDoRecord
func (d *Decoder) todo(xref string) *TodoRecord {
	if xref == "" {
		return &TodoRecord{}
	}

	ref, found := d.refs[xref].(*TodoRecord)
	if !found {
		rec := &TodoRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// trailer finds or creates a level 0 TrailerRecord
func (d *Decoder) trailer(xref string) *TrailerRecord {
	if xref == "" {
		xref = "TRLR"
	}

	ref, found := d.refs[xref].(*TrailerRecord)
	if !found {
		rec := &TrailerRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

// Record parser factories

// makeAddressParser returns a parser for an AddressRecord
func makeAddressParser(d *Decoder, r *AddressRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Full == "" {
				r.Full = value
			} else {
				r.Full = r.Full + "\n" + value
			}

		case "CONC":
			r.Full = r.Full + value

		case "ADR1":
			r.Line1 = value

		case "ADR2":
			r.Line2 = value

		case "ADR3":
			r.Line3 = value

		case "CITY":
			r.City = value

		case "STAE":
			r.State = value

		case "POST":
			r.PostalCode = value

		case "CTRY":
			r.Country = value

		case "PHON":
			rec := &PhoneRecord{Level: level, Phone: value}
			r.Phone = append(r.Phone, rec)
			d.pushParser(makePhoneParser(d, rec, level))

		case "_NAME": // AQ14
			r.Name_ = value

		case "NOTE": // Leg8
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Address tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeAlbumParser returns a parser for an AlbumRecord (MH/FTB8)
func makeAlbumParser(d *Decoder, r *AlbumRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "RIN":
			r.Rin = append(r.Rin, value)

		case "TITL":
			r.Title = value

		case "_PHOTO":
			rec := &PhotoRecord{Level: level}
			r.Photo_ = append(r.Photo_, rec)
			d.pushParser(makePhotoParser(d, rec, level))

		default:
			log.Printf("unhandled Album tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeAttributeParser returns a parser for an AttributeRecord
func makeAttributeParser(d *Decoder, r *AttributeRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_UID": // MH/FTB8
			r.UniqueId_ = append(r.UniqueId_, value)

		case "RIN": // MH/FTB8
			r.Rin = append(r.Rin, value)

		case "TYPE":
			r.Type = value

		case "NAME":
			r.Name = value

			//		case "_PRIM":
			//			r.Primary_ = value

			//		case "_ALT_BIRTH": // AQ14
			//			r.AlternateBirth_ = value

			//		case "_CONFIDENTIAL": // AQ14
			//			r.Confidential_ = value

		case "DATE":
			rec := &DateRecord{Level: level, Tag: tag, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

			//		case "_DATE2": // AQ14
			//			rec := &DateRecord{Level: level, Tag: tag, Date: value}
			//			r.Date2_ = rec
			//			d.pushParser(makeDateParser(d, rec, level))

		case "PLAC":
			rec := &PlaceRecord{Level: level, Tag: tag, Name: value}
			r.Place = rec
			d.pushParser(makePlaceParser(d, rec, level))

			//		case "_PLAC2": // AQ14
			//			rec := &PlaceRecord{Level: level, Tag: tag, Name: value}
			//			r.Place2_ = rec
			//			d.pushParser(makePlaceParser(d, rec, level))

		case "_Description2": // AQ14
			r.Description2_ = value

			//		case "ROLE": // This is a kind of IndividualLink
			//			indi := d.individual(stripXref(value))
			//			rec := &RoleRecord{Level: level, Role: stripValue(value), Individual: indi}
			//			r.Role = append(r.Role, rec)
			//			d.pushParser(makeRoleParser(d, rec, level))

			//		case "ADDR":
			//			rec := &AddressRecord{Level: level, Full: value}
			//			r.Address = rec
			//			d.pushParser(makeAddressParser(d, rec, level))

			//		case "PHON":
			//			rec := &PhoneRecord{Level: level, Phone: value}
			//			r.Phone = append(r.Phone, rec)
			//			d.pushParser(makePhoneParser(d, rec, level))

			//		case "FAMC":
			//			rec := &FamilyLink{Level: level, Tag: tag, Value: value}
			//			r.Parents = append(r.Parents, rec)
			//			d.pushParser(makeFamilyLinkParser(d, rec, level))

			//		case "HUSB":
			//			husband := d.individual(stripXref(value))
			//			rec := &IndividualLink{Level: level, Tag: tag, Individual: husband}
			//			r.Husband = rec
			//			d.pushParser(makeIndividualLinkParser(d, rec, level))

			//		case "WIFE":
			//			wife := d.individual(stripXref(value))
			//			rec := &IndividualLink{Level: level, Tag: tag, Individual: wife}
			//			r.Wife = rec
			//			d.pushParser(makeIndividualLinkParser(d, rec, level))

			//		case "SPOU":
			//			spouse := d.individual(stripXref(value))
			//			rec := &IndividualLink{Level: level, Tag: tag, Individual: spouse}
			//			r.Spouse = rec
			//			d.pushParser(makeIndividualLinkParser(d, rec, level))

			//		case "AGNC":
			//			r.Agency = value

			//		case "CAUS":
			//			r.Cause = value

			//		case "TEMP":
			//			r.Temple = value

			//		case "STAT":
			//			r.Status = value

		case "QUAY":
			r.Quality = value

			//		case "OBJE":
			//			if value != "" {
			//				media := d.media(stripXref(value))
			//				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
			//				r.Media = append(r.Media, rec)
			//				d.pushParser(makeMediaLinkParser(d, rec, level))
			//			} else {
			//				rec := &MediaRecord{Level: level, Tag: tag}
			//				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
			//				r.Media = append(r.Media, link)
			//				d.pushParser(makeMediaParser(d, rec, level))
			//			}

			//		case "_UID":
			//			r.UniqueId_ = append(r.UniqueId_, value)

		case "RecordInternal":
			r.RecordInternal = value

			//		case "EMAIL":
			//			r.Email = value

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

			//		case "_UPD":
			//			r.UpdateTime_ = value

		default:
			log.Printf("unhandled Attribute tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeAuthorParser returns a parser for an AuthorRecord
func makeAuthorParser(d *Decoder, r *AuthorRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Author == "" {
				r.Author = value
			} else {
				r.Author = r.Author + "\n" + value
			}

		case "CONC":
			r.Author = r.Author + value

		case "ABBR":
			r.Abbreviation = value

		default:
			log.Printf("unhandled Author tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeBibliographyParser parses a BibliographyRecord
func makeBibliographyParser(d *Decoder, r *BibliographyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "COMP":
			r.Component = r.Component + value
			d.pushParser(makeTextParser(d, &r.Component, level))

		default:
			log.Printf("unhandled bibliography record tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeBlobParser parses a BlobRecord
func makeBlobParser(d *Decoder, r *BlobRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Data == "" {
				r.Data = value
			} else {
				r.Data = r.Data + "\n" + value
			}

		case "CONC":
			r.Data = r.Data + value

		default:
			log.Printf("unhandled Blob tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeBusinessParser parses a BusinessRecord
func makeBusinessParser(d *Decoder, r *BusinessRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			r.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "PHON":
			rec := &PhoneRecord{Level: level, Phone: value}
			r.Phone = append(r.Phone, rec)
			d.pushParser(makePhoneParser(d, rec, level))

		case "WWW":
			r.WebSite = value

		default:
			log.Printf("unhandled Business tagat %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeCallNumberParser parses a CallNumberRecord
func makeCallNumberParser(d *Decoder, r *CallNumberRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "MEDI":
			r.Media = value

		default:
			log.Printf("unhandled CallNumber tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeChangeParser parses a ChangeRecord
func makeChangeParser(d *Decoder, r *ChangeRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "DATE":
			rec := &DateRecord{Level: level, Tag: tag, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Change tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeCharacterSetParser returns a parser for an CharacterSetRecord
func makeCharacterSetParser(d *Decoder, r *CharacterSetRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			r.Version = value

		default:
			log.Printf("unhandled CharacterSet tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeChildStatusParser returns a parser for an ChildStatusRecord
func makeChildStatusParser(d *Decoder, r *ChildStatusRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			r.Name = value

		default:
			log.Printf("unhandled ChildStatus tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeCitationParser returns a parser for an CitationRecord
func makeCitationParser(d *Decoder, r *CitationRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Value == "" {
				r.Value = value
			} else {
				r.Value = r.Value + "\n" + value
			}

		case "CONC":
			r.Value = r.Value + value

		case "RIN": // MH/FTB8
			r.Rin = append(r.Rin, value)

		case "PAGE":
			r.Page = value

		case "REF":
			r.Reference = value

		case "REFN":
			r.ReferenceNumber = value

		case "_FSFTID": // AQ14
			r.FamilySearchFTID_ = value

		case "QUAY":
			r.Quality = value

		case "CONS":
			r.CONS = value

		case "DIRE":
			r.Direct = value

		case "SOQU":
			r.SourceQuality = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "EVEN":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				rec := &MediaRecord{Level: level, Tag: tag}
				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
				r.Media = append(r.Media, link)
				d.pushParser(makeMediaParser(d, rec, level))
			}

		case "DATA":
			rec := &DataRecord{Level: level, Data: value}
			r.Data = append(r.Data, rec)
			d.pushParser(makeDataParser(d, rec, level))

		case "TEXT":
			r.Text = r.Text + value
			d.pushParser(makeTextParser(d, &r.Text, level))

		case "DATE": // Leg8
			r.Date = value

		case "_RIN": // AQ14
			r.Rin_ = value

		case "_APPLIES_TO": // AQ15
			r.AppliesTo_ = value

		case "_SUBQ", "_BIBL", "_TMPLT", "TID", "FIELD", "NAME", "VALUE": // RM6
			// TODO

		default:
			log.Printf("unhandled Citation tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeDataParser returns a parser for an DataRecord
func makeDataParser(d *Decoder, r *DataRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "DATE":
			r.Date = value

		case "COPR":
			r.Copyright = value

		case "TEXT":
			r.Text = r.Text + value
			d.pushParser(makeTextParser(d, &r.Text, level))

		case "EVEN":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "AGNC":
			r.Agency = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Data tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeDateParser returns a parser for an DateRecord
func makeDateParser(d *Decoder, r *DateRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TIME":
			r.Time = value

		case "TEXT":
			r.Text = r.Text + value
			d.pushParser(makeTextParser(d, &r.Text, level))

		case "DATD":
			r.Day = value

		case "DATM":
			r.Month = value

		case "DATY":
			r.Year = value

		case "DATF":
			r.Full = value

		case "DATS":
			r.Short = value

		case "_TIMEZONE": // MH/FTB8
			r.TimeZone_ = value

		default:
			log.Printf("unhandled Date tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeEventDefinitionParser returns a parser for an EventDefinitionRecord
func makeEventDefinitionParser(d *Decoder, r *EventDefinitionRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE":
			r.Type = value

		case "TITL":
			rec := &TitleRecord{Level: level, Title: value}
			r.Title = append(r.Title, rec)
			d.pushParser(makeTitleParser(d, rec, level))

		case "ABBR":
			r.Abbreviation = value

		case "_SENT":
			r.Sentence_ = value

		case "_DESC_FLAG":
			r.DescriptionFlag_ = value

		case "_Assoc":
			r.Association_ = value

		case "_RIN":
			r.RecordInternal_ = value

		case "ROLE",
			"PLAC", "DATE", "DESC",
			"SENT", "CONT", "CONC",
			"_SENM", "_SENDOM", "_SENPOM", "_SENDPM",
			"_SENF", "_SENDOF", "_SENPOF", "_SENDPF",
			"_SENU", "_SENDOU", "_SENPOU", "_SENDPU",
			"_SEN1", "_SEN2", "_SEN3", "_SEN4",
			"_SEN5", "_SEN6", "_SEN7", "_SEN8",
			"_INC_NOTES", "_DEF", "_PP_EXCLUDE",
			"_DATE_TYPE", "_PLACE_TYPE", "_CONF_FLAG":
			// TODO

		default:
			log.Printf("unhandled Event Definition tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeEventParser returns a parser for an EventRecord
func makeEventParser(d *Decoder, r *EventRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_UID":
			r.UniqueId_ = append(r.UniqueId_, value)

		case "RIN": // MH/FTB8
			r.Rin = append(r.Rin, value)

		case "TYPE":
			r.Type = value

		case "NAME":
			r.Name = value

		case "_PRIM":
			r.Primary_ = value

		case "_ALT_BIRTH": // AQ14
			r.AlternateBirth_ = value

		case "_CONFIDENTIAL": // AQ14
			r.Confidential_ = value

		case "DATE":
			rec := &DateRecord{Level: level, Tag: tag, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "_DATE2": // AQ14
			rec := &DateRecord{Level: level, Tag: tag, Date: value}
			r.Date2_ = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "PLAC":
			rec := &PlaceRecord{Level: level, Tag: tag, Name: value}
			r.Place = rec
			d.pushParser(makePlaceParser(d, rec, level))

		case "_PLAC2": // AQ14
			rec := &PlaceRecord{Level: level, Tag: tag, Name: value}
			r.Place2_ = rec
			d.pushParser(makePlaceParser(d, rec, level))

		case "_Description2": // AQ14
			r.Description2_ = value

		case "AGE": // MH/FTB8
			r.Age = value

		case "ROLE": // This is a kind of IndividualLink
			indi := d.individual(stripXref(value))
			rec := &RoleRecord{Level: level, Role: stripValue(value), Individual: indi}
			r.Role = append(r.Role, rec)
			d.pushParser(makeRoleParser(d, rec, level))

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			r.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "PHON":
			rec := &PhoneRecord{Level: level, Phone: value}
			r.Phone = append(r.Phone, rec)
			d.pushParser(makePhoneParser(d, rec, level))

		case "FAMC":
			rec := &FamilyLink{Level: level, Tag: tag, Value: value}
			r.Parents = append(r.Parents, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "HUSB":
			husband := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: husband}
			r.Husband = rec
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "WIFE":
			wife := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: wife}
			r.Wife = rec
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "SPOU":
			spouse := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: spouse}
			r.Spouse = rec
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "AGNC":
			r.Agency = value

		case "CAUS":
			r.Cause = value

		case "TEMP":
			r.Temple = value

		case "STAT":
			r.Status = value

		case "QUAY":
			r.Quality = value

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				rec := &MediaRecord{Level: level, Tag: tag}
				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
				r.Media = append(r.Media, link)
				d.pushParser(makeMediaParser(d, rec, level))
			}

		case "RecordInternal":
			r.RecordInternal = value

		case "EMAIL":
			r.Email = value

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_UPD":
			r.UpdateTime_ = value

		default:
			log.Printf("unhandled Event tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeFamilyLinkParser returns a parser for an FamilyLink
func makeFamilyLinkParser(d *Decoder, r *FamilyLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}

		//		r.SetFamily(d.FindFamily(value))
		//		if r.family == nil {
		//			log.Printf("Warning: Family not found for '%s'\n", value)
		//			log.Printf("Line number: %d\nLevel: %d\nTag: %s\nValue: %s\nXref: %s\n",
		//				d.LineNum, level, tag, value, xref)
		//			d.CountWarnings()
		//		}

		switch tag {

		case "PEDI":
			rec := &PedigreeRecord{Level: level, Pedigree: value}
			r.Pedigree = rec
			d.pushParser(makePedigreeParser(d, rec, level))

		case "ADOP":
			r.Adopted = value

		case "_PRIMARY":
			r.Primary_ = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "SOUR": // AQ15
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		default:
			log.Printf("unhandled Family Link tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeFamilyParser returns a parser for an FamilyRecord
func makeFamilyParser(d *Decoder, r *FamilyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "RIN":
			r.Rin = append(r.Rin, value)

		case "_STAT": // AQ14
			r.Status_ = value

		case "_NONE": // AQ14
			r.NoChildren_ = value

		case "HUSB":
			husband := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: husband}
			r.Husband = rec
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "WIFE":
			wife := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: wife}
			r.Wife = rec
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "NCHI":
			pint, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Printf("NCHI = %s: %s", value, err.Error())
			}
			r.NumChildren = int(pint)

		case "CHIL":
			child := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: child}
			r.Child = append(r.Child, rec)
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "ANUL", "CENS", "DIV", "DIVF", // FAM 5.5.1
			"ENGA", "MARB", "MARC", // FAM 5.5.1
			"MARR",         // FAM 5.5.1
			"MARL", "MARS", // FAM 5.5.1
			"RESI", // FAM 5.5.1
			"EVEN", // FAM 5.5.1
			"SLGS": // LDS FAM 5.5.1
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "_UID":
			r.UniqueId_ = append(r.UniqueId_, value)

		case "RecordInternal":
			r.RecordInternal = value

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				rec := &MediaRecord{Level: level, Tag: tag}
				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
				r.Media = append(r.Media, link)
				d.pushParser(makeMediaParser(d, rec, level))
			}

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "SUBM":
			subm := d.submitter(stripXref(value))
			rec := &SubmitterLink{Level: level, Tag: tag, Submitter: subm}
			r.Submitter = append(r.Submitter, rec)
			d.pushParser(makeSubmitterLinkParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_UPD":
			r.UpdateTime_ = value

		default:
			log.Printf("unhandled Family tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}
		return nil
	}
}

// makeFootnoteParser returns a parser for an FootnoteRecord
func makeFootnoteParser(d *Decoder, r *FootnoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "COMP":
			r.Component = r.Component + value
			d.pushParser(makeTextParser(d, &r.Component, level))

		default:
			log.Printf("unhandled Footnote tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeGedcomParser returns a parser for an GedcomRecord
func makeGedcomParser(d *Decoder, r *GedcomRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			r.Version = value

		case "FORM":
			r.Form = value

		default:
			log.Printf("unhandled Gedcom tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}
		return nil
	}
}

// makeHeaderParser returns a parser for an HeaderRecord
func makeHeaderParser(d *Decoder, r *HeaderRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SOUR": // This is not a citation
			rec := &SystemRecord{Level: level, SystemName: value}
			r.SourceSystem = rec
			d.pushParser(makeSystemParser(d, rec, level))

		case "DEST":
			r.Destination = value

		case "DATE":
			rec := &DateRecord{Level: level, Tag: tag, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "TIME":
			r.Time = value

		case "GEDC":
			rec := &GedcomRecord{Level: level}
			r.Gedcom = rec
			d.pushParser(makeGedcomParser(d, rec, level))

		case "CHAR":
			rec := &CharacterSetRecord{Level: level, CharacterSet: value}
			r.CharacterSet = rec
			d.pushParser(makeCharacterSetParser(d, rec, level))

		case "LANG":
			r.Language = value

		case "FILE":
			r.FileName = value

		case "_RINS": // MH/FTB8
			r.Rins_ = value

		case "_UID": // MH/FTB8
			r.Uid_ = value

		case "_PROJECT_GUID": // MH/FTB8
			r.ProjectGuid_ = value

		case "_EXPORTED_FROM_SITE_ID": // MH/FTB8
			r.ExportedFromSiteId_ = value

		case "_SM_MERGES": // MH/FTB8
			r.SmMerges_ = value

		case "_DESCRIPTION_AWARE": // MH/FTB8
			r.DescriptionAware_ = value

		case "COPR":
			r.Copyright = value

		case "PLAC":
			rec := &PlaceRecord{Level: level, Tag: tag, Name: value}
			r.Place = rec
			d.pushParser(makePlaceParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "SUBM":
			subm := d.submitter(stripXref(value))
			rec := &SubmitterLink{Level: level, Tag: tag, Submitter: subm}
			r.Submitter = append(r.Submitter, rec)
			d.pushParser(makeSubmitterLinkParser(d, rec, level))

		case "SUBN":
			subn := d.submission(stripXref(value))
			rec := &SubmissionLink{Level: level, Submission: subn}
			r.Submission = append(r.Submission, rec)
			d.pushParser(makeSubmissionLinkParser(d, rec, level))

		case "SCHEMA":
			rec := &SchemaRecord{Level: level}
			r.Schema = rec
			d.pushParser(makeSchemaParser(d, rec, level))

		case "_ROOT":
			root := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: root}
			r.RootPerson_ = rec

		case "_HME":
			home := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: home}
			r.HomePerson_ = rec

		default:
			log.Printf("unhandled Header tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}
		return nil
	}
}

// makeHistoryParser returns a parser for an HistoryRecord
func makeHistoryParser(d *Decoder, r *HistoryRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.History == "" {
				r.History = value
			} else {
				r.History = r.History + "\n" + value
			}

		case "CONC":
			r.History = r.History + value

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		default:
			log.Printf("unhandled History tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeIndividualLinkParser returns a parser for an IndividualLink
func makeIndividualLinkParser(d *Decoder, r *IndividualLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "RELA":
			r.Relationship = value

		case "SLGC":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "AGE":
			r.Age = value

		case "_PREF": // Leg8
			r.Preferred_ = value

		default:
			log.Printf("unhandled Individual Link tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeIndividualParser returns a parser for an IndividualRecord
func makeIndividualParser(d *Decoder, r *IndividualRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "RIN":
			r.Rin = append(r.Rin, value)

		case "NAME":
			rec := &NameRecord{Level: level, Name: value}
			r.Name = append(r.Name, rec)
			d.pushParser(makeNameParser(d, rec, level))

		case "_STAT": // AQ14
			r.Status_ = value

		case "RESN":
			r.Restriction = value

		case "SEX":
			r.Sex = value

		case "ATTR", "CAST", "DSCR", "EDUC", "IDNO", // 5.5.1
			"NATI", "NCHI", "NMR", "OCCU", "PROP", "RELI", // 5.5.1
			"SSN", "TITL", "FACT": // 5.5.1
			rec := &AttributeRecord{Level: level, Tag: tag, Value: value}
			r.Attribute = append(r.Attribute, rec)
			d.pushParser(makeAttributeParser(d, rec, level))

		case "BIRT", "CHR", "DEAT", "BURI", "CREM", // INDI 5.5.1
			"ADOP",                         // INDI 5.5.1
			"BAPM", "BARM", "BASM", "BLES", // INDI 5.5.1
			"CHRA", "CONF", "FCOM", "ORDN", // INDI 5.5.1
			"NATU", "EMIG", "IMMI", // INDI 5.5.1
			"CENS", "PROB", "WILL", // INDI 5.5.1
			"GRAD", "RETI", // INDI 5.5.1
			"EVEN",                         // INDI 5.5.1
			"BAPL", "CONL", "ENDL", "SLGC", // LDS INDI 5.5.1
			"ELEC", "ILLN", "IMMIG", "MILI",
			"MILI_AWA", "MILI_RET", "RESD", "RESI",
			"TRAV", "WAR":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "FAMC":
			rec := &FamilyLink{Level: level, Tag: tag, Value: value}
			r.Parents = append(r.Parents, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "FAMS":
			rec := &FamilyLink{Level: level, Tag: tag, Value: value}
			r.Family = append(r.Family, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				rec := &MediaRecord{Level: level, Tag: tag}
				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
				r.Media = append(r.Media, link)
				d.pushParser(makeMediaParser(d, rec, level))
			}

		case "HEAL":
			r.Health = value

		case "QUAY":
			r.Quality = value

		case "LVG":
			r.Living = value

		case "AFN":
			r.AncestralFileNumber = append(r.AncestralFileNumber, value)

		case "RFN":
			r.RecordFileNumber = value

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "_FSFTID": // AQ14
			r.FamilySearchFTID_ = value

		case "_FSLINK": // Leg8
			r.FamilySearchLink_ = value

		case "_UID":
			r.UniqueId_ = append(r.UniqueId_, value)

		case "RecordInternal":
			r.RecordInternal = value

		case "EMAIL":
			r.Email = value

		case "_EMAIL": // AQ14
			r.Email_ = value

		case "_URL": // AQ14
			r.URL_ = value

		case "WWW":
			r.WebSite = value

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			r.Address = append(r.Address, rec)
			d.pushParser(makeAddressParser(d, rec, level))

		case "HIST":
			rec := &HistoryRecord{Level: level, History: value}
			r.History = append(r.History, rec)
			d.pushParser(makeHistoryParser(d, rec, level))

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "ASSO":
			assoc := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: assoc}
			r.Associated = append(r.Associated, rec)
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "SUBM":
			subm := d.submitter(stripXref(value))
			rec := &SubmitterLink{Level: level, Tag: tag, Submitter: subm}
			r.Submitter = append(r.Submitter, rec)
			d.pushParser(makeSubmitterLinkParser(d, rec, level))

		case "ANCI":
			subm := d.submitter(stripXref(value))
			rec := &SubmitterLink{Level: level, Tag: tag, Submitter: subm}
			r.ANCI = append(r.ANCI, rec)
			d.pushParser(makeSubmitterLinkParser(d, rec, level))

		case "DESI":
			subm := d.submitter(stripXref(value))
			rec := &SubmitterLink{Level: level, Tag: tag, Submitter: subm}
			r.DESI = append(r.DESI, rec)
			d.pushParser(makeSubmitterLinkParser(d, rec, level))

		case "_UPD":
			r.UpdateTime_ = value

		case "ALIA":
			r.Alias = value

		case "FATH":
			father := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: father}
			r.Father = rec

		case "MOTH":
			mother := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: mother}
			r.Mother = rec

		case "PHON":
			rec := &PhoneRecord{Level: level, Phone: value}
			r.Phone = append(r.Phone, rec)
			d.pushParser(makePhoneParser(d, rec, level))

		case "MISC":
			r.Miscellaneous = append(r.Miscellaneous, value)

		case "_PROF":
			media := d.media(stripXref(value))
			rec := &MediaLink{Level: level, Tag: tag, Media: media}
			r.ProfilePicture_ = rec
			d.pushParser(makeMediaLinkParser(d, rec, level))

		case "_PPEXCLUDE": // Leg8
			r.PPExclude_ = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_TODO": // AQ15
			r.Todo_ = append(r.Todo_, value)

		case "Anecdote": // (Custom - MH/FTB8)
			r.Anecdote = append(r.Anecdote, value)

		default:
			log.Printf("unhandled Individual tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}
		return nil
	}
}

// makeMediaLinkParser returns a parser for an MediaLink
func makeMediaLinkParser(d *Decoder, r *MediaLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		default:
			log.Printf("unhandled MediaLink tag at %d: %d %s %s\r", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeMediaParser makes a MediaRecord parser function
func makeMediaParser(d *Decoder, r *MediaRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "FORM":
			r.Format = value

		case "_URL":
			r.Url_ = value

		case "FILE":
			r.FileName = value

		case "TITL":
			r.Title = value

		case "DATE":
			r.Date = value

		case "AUTH":
			r.Author = value

		case "TEXT":
			r.Text = value
			d.pushParser(makeTextParser(d, &r.Text, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "_DATE": // (MH/FTB8)
			r.Date_ = value

		case "_PLACE": // (MH/FTB8)
			r.Place_ = value

		case "_ASTID":
			r.AstId_ = value

		case "_ASTTYP":
			r.AstType_ = value

		case "_ASTDESC":
			r.AstDesc_ = value

		case "_ASTLOC":
			r.AstLoc_ = value

		case "_ASTPERM":
			r.AstPerm_ = value

		case "_ASTUPPID":
			r.AstUpPid_ = value

		case "BLOB":
			rec := &BlobRecord{Level: level, Data: value}
			r.BinaryLargeObject = rec
			d.pushParser(makeBlobParser(d, rec, level))

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "RecordInternal":
			r.RecordInternal = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_FSFTID": // AQ15
			r.FsFtId_ = value

		case "_SCBK": // AQ14
			r.Scbk_ = value

		case "_PRIM": // AQ14(MH/FTB8)
			r.Primary_ = value

		case "_SCAN": // AQ14(MH/FTB8)
			r.Scan_ = value

		case "_TYPE": // AQ14
			r.Type_ = value

		case "_SSHOW": // AQ14
			rec := &SlideShowRecord{Level: level, Included: value}
			r.Sshow_ = rec
			d.pushParser(makeSlideShowParser(d, rec, level))

		case "_PRIM_CUTOUT": // MH/FTB8
			r.PrimCutout_ = value

		case "_CUTOUT": // MH/FTB8
			r.Cutout_ = value

		case "_POSITION": // MH/FTB8
			r.Position_ = value

		case "_ALBUM": // MH/FTB8
			r.Album_ = value

		case "_PHOTO_RIN": // MH/FTB8
			r.PhotoRin_ = value

		case "_FILESIZE": // MH/FTB8
			r.Filesize_ = value

		case "_PARENTRIN": // MH/FTB8
			r.ParentRin_ = value

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.mediaLinks = append(r.mediaLinks, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				log.Fatal("OBJE inside OBJE")
			}

		case "_SRCPP": // AQ15
			r.SrcPp_ = value

		case "_SRCFLIP": // AQ15
			r.SrcFlip_ = value

		default:
			log.Printf("unhandled Media tag at %d: %d %s %s\r", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeNameParser returns a parser for an NameRecord
func makeNameParser(d *Decoder, r *NameRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NPFX":
			r.Prefix = value

		case "GIVN":
			r.GivenName = value

		case "_MIDN":
			r.MiddleName_ = value

		case "SPFX":
			r.SurnamePrefix = value

		case "SURN":
			r.Surname = value

		case "NSFX":
			r.Suffix = value

		case "_PGVN":
			r.PreferedGivenName_ = value

		case "ROMN":
			r.RomanizedName = value

		case "FONE":
			r.PhoneticName = value

		case "_FORMERNAME": // MH/FTB8
			r.FormerName_ = value

		case "_MARNM": // AQ14, MH/FTB8
			r.MarriedName_ = value

		case "TYPE":
			r.NameType = value

		case "_PRIM":
			r.Primary_ = value

		case "_AKA":
			r.AlsoKnownAs_ = append(r.AlsoKnownAs_, value)

		case "NICK":
			r.Nickname = append(r.Nickname, value)

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Name tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeNoteParser returns a parser for an NoteRecord
func makeNoteParser(d *Decoder, r *NoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Note == "" {
				r.Note = value
			} else {
				r.Note = r.Note + "\n" + value
			}

		case "CONC":
			r.Note = r.Note + value

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "RecordInternal":
			r.RecordInternal = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_DESCRIPTION": // MH/FTB8
			r.Description_ = value

		default:
			log.Printf("unhandled Note tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePedigreeParser returns a parser for an PedigreeRecord
func makePedigreeParser(d *Decoder, r *PedigreeRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_HUSB":
			r.Husband_ = value

		case "_WIFE":
			r.Wife_ = value

		default:
			log.Printf("unhandled Pedigree tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePhoneParser returns a parser for an PhoneRecord
func makePhoneParser(d *Decoder, r *PhoneRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_TYPE": // MH/FTB8
			r.Type_ = value

		default:
			log.Printf("unhandled Phone tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePhotoParser returns a parser for an PhotoRecord (MH/FTB8)
func makePhotoParser(d *Decoder, r *PhotoRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_UID":
			r.Uid_ = value

		case "_PRIN":
			r.Prin_ = value

		default:
			log.Printf("unhandled Photo tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePlaceDefinitionParser returns a parser for an PlaceDefinitionRecord
func makePlaceDefinitionParser(d *Decoder, r *PlaceDefinitionRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PLAC":
			r.Place = value

		case "ABBR":
			r.Abbreviation = value

		default:
			log.Printf("unhandled Place Definition tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePlacePartParser returns a parser for an PlacePartRecord
func makePlacePartParser(d *Decoder, r *PlacePartRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "JURI":
			r.Jurisdiction = value

		default:
			log.Printf("unhandled Place Part tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePlaceParser returns a parser for an PlaceRecord
func makePlaceParser(d *Decoder, r *PlaceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "FORM":
			r.Form = value

		case "PLAS":
			r.ShortName = value

		case "PLAM":
			r.Modifier = value

		case "PLA0", "PLA1", "PLA2", "PLA3", "PLA4":
			rec := &PlacePartRecord{Level: level, Tag: tag, Part: value}
			r.Parts = append(r.Parts, rec)
			d.pushParser(makePlacePartParser(d, rec, level))

		case "SOUR":
			rec := &CitationRecord{Level: level, Value: value}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Place tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makePublishParser returns a parser for an PublishRecord
func makePublishParser(d *Decoder, r *PublishRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_SITEADDRESS":
			r.SiteAddress_ = value

		case "_SITENAME":
			r.SiteName_ = value

		case "_SITEID":
			r.SiteId_ = value

		case "_USERNAME":
			r.UserName_ = value

		case "_DISABLED":
			r.Disabled_ = value

		default:
			log.Printf("unhandled Publish_ tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeRepositoryLinkParser returns a parser for an RepositoryLink
func makeRepositoryLinkParser(d *Decoder, r *RepositoryLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CALN":
			rec := &CallNumberRecord{Level: level, CallNumber: value}
			r.CallNumber = rec
			d.pushParser(makeCallNumberParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Repository Link tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeRepositoryParser returns a parser for an RepositoryRecord
func makeRepositoryParser(d *Decoder, r *RepositoryRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			r.Name = value

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			r.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "PHON":
			rec := &PhoneRecord{Level: level, Phone: value}
			r.Phone = append(r.Phone, rec)
			d.pushParser(makePhoneParser(d, rec, level))

		case "WWW":
			r.WebSite = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "RecordInternal":
			r.RecordInternal = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Repository tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeRoleParser returns a parser for an RoleRecord
func makeRoleParser(d *Decoder, r *RoleRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PRIN":
			r.Principal = value

		default:
			log.Printf("unhandled Role tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeRootParser returns a parser for an RootRecord
func makeRootParser(d *Decoder, r *RootRecord) parser {
	return func(level int, tag string, value string, xref string) error {
		//log.Println(level, tag, value, xref)
		if level == 0 {
			// cases ordered approx. by frequency of appearance
			r.Level = level // always zero
			switch tag {

			case "INDI":
				rec := d.individual(xref)
				r.Individual = append(r.Individual, rec)
				d.pushParser(makeIndividualParser(d, rec, level))

			case "FAM":
				rec := d.family(xref)
				r.Family = append(r.Family, rec)
				d.pushParser(makeFamilyParser(d, rec, level))

			case "SOUR":
				rec := d.source(xref)
				rec.Value = value
				r.Source = append(r.Source, rec)
				d.pushParser(makeSourceParser(d, rec, level))

			case "EVEN":
				rec := d.event(xref)
				rec.Tag = tag
				rec.Value = value
				r.Event = append(r.Event, rec)
				d.pushParser(makeEventParser(d, rec, level))

			case "PLAC":
				rec := d.place(xref)
				rec.Name = value
				r.Place = append(r.Place, rec)
				d.pushParser(makePlaceParser(d, rec, level))

			case "REPO":
				rec := d.repository(xref)
				r.Repository = append(r.Repository, rec)
				d.pushParser(makeRepositoryParser(d, rec, level))

			case "_PLAC_DEFN": // Leg8
				rec := d.placeDefinition(xref)
				r.PlaceDefinition_ = append(r.PlaceDefinition_, rec)
				d.pushParser(makePlaceDefinitionParser(d, rec, level))

			case "_EVENT_DEFN": // AQ14
				rec := d.eventDefinition(xref)
				rec.Tag = tag
				rec.Name = value
				r.EventDefinition_ = append(r.EventDefinition_, rec)
				d.pushParser(makeEventDefinitionParser(d, rec, level))

			case "_TODO": // AQ15
				rec := d.todo(xref)
				r.Todo_ = append(r.Todo_, rec)
				d.pushParser(makeTodoParser(d, rec, level))

			case "_EVDEF": // RM6
				rec := d.eventDefinition(xref)
				r.EventDefinition_ = append(r.EventDefinition_, rec)
				d.pushParser(makeEventDefinitionParser(d, rec, level))

			case "NOTE":
				rec := d.note(xref)
				r.Note = append(r.Note, rec)
				d.pushParser(makeNoteParser(d, rec, level))

			case "OBJE":
				rec := d.media(xref)
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaParser(d, rec, level))

			case "CSTA":
				rec := d.childStatus(xref)
				r.ChildStatus = append(r.ChildStatus, rec)
				d.pushParser(makeChildStatusParser(d, rec, level))

			case "SUBM":
				rec := d.submitter(xref)
				r.Submitter = append(r.Submitter, rec)
				d.pushParser(makeSubmitterParser(d, rec, level))

			case "SUBN":
				rec := d.submission(xref)
				r.Submission = append(r.Submission, rec)
				d.pushParser(makeSubmissionParser(d, rec, level))

			case "ALBUM": // MN/FTB8
				rec := d.album(xref)
				r.Album = append(r.Album, rec)
				d.pushParser(makeAlbumParser(d, rec, level))

			case "_PUBLISH": // MN/FTB8
				rec := d.publish(xref)
				r.Publish_ = append(r.Publish_, rec)
				d.pushParser(makePublishParser(d, rec, level))

			case "HEAD":
				rec := d.header(xref)
				r.Header = rec
				d.pushParser(makeHeaderParser(d, rec, level))

			case "TRLR":
				rec := d.trailer(xref)
				r.Trailer = rec
				// There should be nothing to parse in trailer
				//d.pushParser(makeTrailerParser(d, obj, level))

			default:
				log.Printf("unhandled Root tag at %d: %d @%s@ %s %s\n", d.LineNum, level, xref, tag, value)
			}
		} else {
			log.Printf("Not level 0 Root tag at %d: %d @%s@ %s %s\n", d.LineNum, level, xref, tag, value)
		}
		return nil
	}
}

// makeSchemaParser returns a parser for an SchemaRecord
func makeSchemaParser(d *Decoder, r *SchemaRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "AGER", "AKA", "CLASS", "COMP", "CONT", "DETAIL1", "DETAIL2",
			"EVEN", "LANG", "NAME", "NOTE", "PERI", "POSB", "POSF", "PREB",
			"PREF", "PRIN", "ROLE", "SEX", "SOUR", "STYL":
			s := fmt.Sprintf("%d %s %s", level, tag, value)
			r.Data = append(r.Data, s)

		default:
			log.Printf("unhandled Schema tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeShortTitleParser returns a parser for an ShortTitleRecord
func makeShortTitleParser(d *Decoder, r *ShortTitleRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "INDX":
			r.Indexed = value

		default:
			log.Printf("unhandled Short Title tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSlideShowParser returns a parser for an SlideShowRecord
func makeSlideShowParser(d *Decoder, r *SlideShowRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "_STIME":
			r.ShowTime_ = value

		default:
			log.Printf("unhandled Slide Show tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSourceParser returns a parser for an SourceRecord
func makeSourceParser(d *Decoder, r *SourceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Value == "" {
				r.Value = value
			} else {
				r.Value = r.Value + "\n" + value
			}

		case "CONC":
			r.Value = r.Value + value

		case "RIN": // MH/FT8
			r.Rin = append(r.Rin, value)

		case "NAME":
			r.Name = value
			d.pushParser(makeTextParser(d, &r.Name, level))

		case "TITL":
			r.Title = value
			d.pushParser(makeTextParser(d, &r.Title, level))

		case "AUTH":
			rec := &AuthorRecord{Level: level, Author: value}
			r.Author = rec
			d.pushParser(makeAuthorParser(d, rec, level))

		case "ABBR":
			r.Abbreviation = value
			d.pushParser(makeTextParser(d, &r.Abbreviation, level))

		case "PUBL":
			r.Publication = value
			d.pushParser(makeTextParser(d, &r.Publication, level))

		case "MEDI":
			r.MediaType = value

		case "_PAREN":
			r.Parenthesized_ = value

		case "_MEDI": // MH/FTB8
			r.Medi_ = value

		case "_TYPE":
			r.Type_ = value

		case "_OTHER":
			r.Other_ = value

		case "_MASTER":
			r.Master_ = value

		case "_ITALIC":
			r.Italic_ = value

		case "TEXT":
			r.Text = r.Text + value
			d.pushParser(makeTextParser(d, &r.Text, level))

		case "DATA":
			rec := &DataRecord{Level: level, Data: value}
			r.Data = rec
			d.pushParser(makeDataParser(d, rec, level))

		case "SHAU":
			r.ShortAuthor = value

		case "SHTI":
			rec := &ShortTitleRecord{Level: level, ShortTitle: value}
			r.ShortTitle = rec
			d.pushParser(makeShortTitleParser(d, rec, level))

		case "FOOT":
			rec := &FootnoteRecord{Level: level, Value: value}
			r.Footnote = rec
			d.pushParser(makeFootnoteParser(d, rec, level))

		case "BIBL":
			rec := &BibliographyRecord{Level: level, Value: value}
			r.Bibliography = rec
			d.pushParser(makeBibliographyParser(d, rec, level))

		case "REPO":
			repo := d.repository(stripXref(value))
			rec := &RepositoryLink{Level: level, Repository: repo}
			r.Repository = rec
			d.pushParser(makeRepositoryLinkParser(d, rec, level))

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				rec := &MediaRecord{Level: level, Tag: tag, Value: value}
				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
				r.Media = append(r.Media, link)
				d.pushParser(makeMediaParser(d, rec, level))
			}

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "QUAY":
			r.Quality = value

		case "RecordInternal":
			r.RecordInternal = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_WEBTAG": // RM6
			rec := &WebTagRecord{Level: level, Value: value}
			r.WebTag_ = rec
			d.pushParser(makeWebTagParser(d, rec, level))

		case "_SUBQ", "_BIBL", "_TMPLT", "TID", "FIELD", "VALUE": // RM6
			// TODO

		default:
			log.Printf("unhandled Source tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSubmissionLinkParser returns a parser for an SubmissionLink
func makeSubmissionLinkParser(d *Decoder, r *SubmissionLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		default:
			log.Printf("unhandled Submission Link tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSubmissionParser returns a parser for an SubmissionRecord
func makeSubmissionParser(d *Decoder, r *SubmissionRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SUBM":
			rec := d.submitter(xref)
			r.Submitter = rec
			d.pushParser(makeSubmitterParser(d, rec, level))

		case "FAMF":
			r.FamilyFileName = value

		case "TEMP":
			r.Temple = value

		case "ANCE":
			r.Ancestors = value

		case "DESC":
			r.Descendents = value

		case "ORDI":
			r.Ordinance = value

		case "RecordInternal":
			r.RecordInternal = value

		default:
			log.Printf("unhandled Submission tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSubmitterLinkParser returns a parser for an SubmitterLink
func makeSubmitterLinkParser(d *Decoder, r *SubmitterLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		default:
			log.Printf("unhandled Submitter Link tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSubmitterParser returns a parser for an SubmitterRecord
func makeSubmitterParser(d *Decoder, r *SubmitterRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "RIN": // MH/FTB8
			r.Rin = append(r.Rin, value)

		case "NAME":
			r.Name = value
			//d.pushParser(makeTextParser(d, &r.Name, level))

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			r.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "CTRY":
			r.Country = value

		case "PHON":
			rec := &PhoneRecord{Level: level, Phone: value}
			r.Phone = append(r.Phone, rec)
			d.pushParser(makePhoneParser(d, rec, level))

		case "EMAIL":
			r.Email = value

		case "_EMAIL": // AQ14
			r.Email_ = value

		case "WWW":
			r.WebSite = value

		case "LANG":
			r.Language = value

		case "OBJE":
			if value != "" {
				media := d.media(stripXref(value))
				rec := &MediaLink{Level: level, Tag: tag, Value: value, Media: media}
				r.Media = append(r.Media, rec)
				d.pushParser(makeMediaLinkParser(d, rec, level))
			} else {
				rec := &MediaRecord{Level: level, Tag: tag}
				link := &MediaLink{Level: level, Tag: tag, Value: value, Media: rec}
				r.Media = append(r.Media, link)
				d.pushParser(makeMediaParser(d, rec, level))
			}

		case "RFN":
			r.RecordFileNumber = value

		case "RecordInternal":
			r.RecordInternal = value

		case "STAL":
			r.STAL = value

		case "NUMB":
			r.NUMB = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Submitter tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeSystemParser returns a parser for an SystemRecord
func makeSystemParser(d *Decoder, r *SystemRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			r.Version = value

		case "NAME":
			r.ProductName = value

		case "CORP":
			rec := &BusinessRecord{Level: level, BusinessName: value}
			r.Business = rec
			d.pushParser(makeBusinessParser(d, rec, level))

		case "DATA":
			rec := &DataRecord{Level: level, Data: value}
			r.SourceData = rec
			d.pushParser(makeDataParser(d, rec, level))

		case "_RTLSAVE": // MH/FTB8
			r.RtlSave_ = value

		default:
			log.Printf("unhandled System tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}
		return nil
	}
}

// makeTextParser returns a parser for an string
func makeTextParser(d *Decoder, s *string, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		// no Level here
		switch tag {

		case "CONT":
			if *s == "" {
				*s = value
			} else {
				*s = *s + "\n" + value
			}

		case "CONC":
			*s = *s + value

		default:
			log.Printf("unhandled Text tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeTitleParser returns a parser for an TitleRecord
func makeTitleParser(d *Decoder, r *TitleRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "ABBR":
			r.Abbreviation = value

		default:
			log.Printf("unhandled Title tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeTodoParser returns a parser for an TodoRecord
func makeTodoParser(d *Decoder, r *TodoRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "DESC":
			r.Description = value

		case "_PRIORITY":
			r.Priority_ = value

		case "_CAT":
			r.Category_ = value

		case "TYPE":
			r.Type = value

		case "STAT":
			r.Status = value

		case "DATE":
			r.Date = value

		case "_DATE2":
			r.Date2_ = value

		default:
			log.Printf("unhandled Todo tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeUserReferenceNumberParser returns a parser for an UserReferenceNumberRecord
func makeUserReferenceNumberParser(d *Decoder, r *UserReferenceNumberRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE":
			r.Type = value

		default:
			log.Printf("unhandled UserReferenceNumber tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// makeWebTagParser returns a parser for an WebTagRecord (RM6)
func makeWebTagParser(d *Decoder, r *WebTagRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			if r.Value == "" {
				r.Value = value
			} else {
				r.Value = r.Value + "\n" + value
			}

		case "CONC":
			r.Value = r.Value + value

		case "NAME":
			r.Name = value

		case "URL":
			r.URL = value

		default:
			log.Printf("unhandled WebTag tag at %d: %d %s %s\n", d.LineNum, level, tag, value)
		}

		return nil
	}
}

// stripXref removes value and surrounding @s from xref
func stripXref(value string) string {
	return value
	//	if value == "" {
	//		return ""
	//	}
	//	if value[0] == '@' {
	//		return strings.Trim(value, "@")
	//	}
	//	atIndex := strings.IndexByte(value, '@')
	//	if atIndex >= 0 {
	//		return strings.Trim(value[atIndex:], "@")
	//	}
	//	return ""
}

// stripValue removes @ bracketed xref from value
func stripValue(value string) string {
	return value
	//	if value == "" {
	//		return ""
	//	}
	//	atIndex := strings.IndexByte(value, '@')
	//	if atIndex >= 0 {
	//		return strings.Trim(value[:atIndex], " ")
	//	}
	//	return value
}

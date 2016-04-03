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
	"strings"
)

// A Decoder reads and decodes GEDCOM objects from an input stream.
type Decoder struct {
	r       io.Reader
	parsers []parser
	refs    map[string]interface{}
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the next GEDCOM-encoded value from its
// input and stores it in the value pointed to by v.
func (d *Decoder) Decode() (*RootRecord, error) {

	r := &RootRecord{
		Level:            -1,
		Place:            make(PlaceRecords, 0),
		Event:            make(EventRecords, 0),
		Individual:       make(IndividualRecords, 0),
		Family:           make(FamilyRecords, 0),
		Repository:       make(RepositoryRecords, 0),
		Source:           make(SourceRecords, 0),
		Media:            make(MediaLinks, 0),
		Note:             make(NoteRecords, 0),
		EventDefinition_: make(EventDefinitionRecords, 0),
		ChildStatus:      make(ChildStatusRecords, 0),
		Submission:       make(SubmissionRecords, 0),
		Submitter:        make(SubmitterRecords, 0),
	}

	d.refs = make(map[string]interface{})
	d.parsers = []parser{makeRootParser(d, r)}
	d.scan(r)

	return r, nil
}

func (d *Decoder) scan(r *RootRecord) {
	s := &scanner{}
	buf := make([]byte, 512)

	n, err := d.r.Read(buf)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for n > 0 {
		pos := 0

		for {
			s.reset()
			offset, err := s.nextTag(buf[pos:n])
			pos += offset
			if err != nil {
				if err != io.EOF {
					log.Println(err.Error())
					return
				}
				break
			}

			d.parsers[len(d.parsers)-1](s.level, string(s.tag), string(s.value), string(s.xref))

		}

		// shift unparsed bytes to start of buffer
		rest := copy(buf, buf[pos:])

		// top up buffer
		num, err := d.r.Read(buf[rest:len(buf)])
		if err != nil {
			break
		}

		n = rest + num - 1

	}

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

// Level 0 record constructors

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

	ref, found := d.refs[xref].(*FamilyRecord)
	if !found {
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
		return &PlaceRecord{}
	}

	ref, found := d.refs[xref].(*PlaceRecord)
	if !found {
		rec := &PlaceRecord{Xref: xref}
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

// makeAddressParser parses an AddressRecord
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
			r.Phone = value

		default:
			log.Printf("unhandled Address tag: %d %s %s\n", level, tag, value)
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
			r.Component = append(r.Component, value)
			d.pushParser(makeTextParser(d, &r.Component[len(r.Component)-1], level))

		default:
			log.Printf("unhandled bibliography record tag: %d %s %s\n", level, tag, value)
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
			log.Printf("unhandled Blob tag: %d %s %s\n", level, tag, value)
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
			r.Phone = append(r.Phone, value)

		case "WWW":
			r.WebSite = value

		default:
			log.Printf("unhandled Business tag: %d %s %s\n", level, tag, value)
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
			log.Printf("unhandled CallNumber tag: %d %s %s\n", level, tag, value)
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
			rec := &DateRecord{Level: level, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Change tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeCharacterSetParser(d *Decoder, r *CharacterSetRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			r.Version = value

		default:
			log.Printf("unhandled CharacterSet tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeChildStatusParser(d *Decoder, r *ChildStatusRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			r.Name = value

		default:
			log.Printf("unhandled ChildStatus tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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

		case "PAGE":
			r.Page = value

		case "REF":
			r.Reference = value

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
			r.Text = append(r.Text, value)
			d.pushParser(makeTextParser(d, &r.Text[len(r.Text)-1], level))

		default:
			log.Printf("unhandled Citation tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			r.Text = append(r.Text, value)
			d.pushParser(makeTextParser(d, &r.Text[len(r.Text)-1], level))

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
			log.Printf("unhandled Data tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeDateParser(d *Decoder, r *DateRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TIME":
			r.Time = value

		case "TEXT":
			r.Text = append(r.Text, value)
			d.pushParser(makeTextParser(d, &r.Text[len(r.Text)-1], level))

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

		default:
			log.Printf("unhandled Date tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeEventDefinitionParser(d *Decoder, r *EventDefinitionRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE", "TITL", "ABBR",
			"_SENM", "_SENDOM", "_SENPOM", "_SENDPM",
			"_SENF", "_SENDOF", "_SENPOF", "_SENDPF",
			"_SENU", "_SENDOU", "_SENPOU", "_SENDPU",
			"_DATE_TYPE", "_PLACE_TYPE", "_DESC_FLAG", "_CONF_FLAG":
			// TODO

		default:
			log.Printf("unhandled Event Definition tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeEventParser(d *Decoder, r *EventRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE":
			r.Type = value

		case "NAME":
			r.Name = value

		case "_PRIM":
			r.Primary_ = value

		case "DATE":
			rec := &DateRecord{Level: level, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "PLAC":
			rec := &PlaceRecord{Level: level, Name: value}
			r.Place = rec
			d.pushParser(makePlaceParser(d, rec, level))

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
			r.Phone = append(r.Phone, value)

		case "FAMC":
			family := d.family(stripXref(value))
			rec := &FamilyLink{Level: level, Tag: tag, Family: family}
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

		case "_UID":
			r.UID_ = append(r.UID_, value)

		case "RIN":
			r.RIN = value

		case "EMAIL":
			r.Email = value

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
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
			log.Printf("unhandled Event tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeFamilyLinkParser(d *Decoder, r *FamilyLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PEDI":
			r.Pedigree = value

		case "ADOP":
			r.Adopted = value

		case "_PRIMARY":
			r.Primary_ = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Family Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeFamilyParser(d *Decoder, r *FamilyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

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

		case "ANUL", "CENS", "DIV", "DIVF", "ENGA", "EVEN", "MARR", "MARB",
			"MARC", "MARL", "MARS", "SLGC", "SLGS":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "_UID":
			r.UID_ = append(r.UID_, value)

		case "RIN":
			r.RIN = value

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
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
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
			log.Printf("unhandled Family tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

func makeFootnoteParser(d *Decoder, r *FootnoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "COMP":
			r.Component = append(r.Component, value)
			d.pushParser(makeTextParser(d, &r.Component[len(r.Component)-1], level))

		default:
			log.Printf("unhandled Footnote tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			log.Printf("unhandled Gedcom tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

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
			rec := &DateRecord{Level: level, Date: value}
			r.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

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

		case "COPR":
			r.Copyright = value

		case "PLAC":
			rec := &PlaceRecord{Level: level, Name: value}
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
			log.Printf("unhandled Header tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

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
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		default:
			log.Printf("unhandled History tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "AGE":
			r.Age = value

		default:
			log.Printf("unhandled Individual Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeIndividualParser(d *Decoder, r *IndividualRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			rec := &NameRecord{Level: level, Name: value}
			r.Name = append(r.Name, rec)
			d.pushParser(makeNameParser(d, rec, level))

		case "RESN":
			r.Restriction = value

		case "SEX":
			r.Sex = value

		case "ADOP", "BAPL", "BAPM", "BARM", "BASM", "BIRT", "BLES", "BURI",
			"CAST", "CENS", "CHR", "CHRA", "CONF", "CREM", "DEAT", "DSCR",
			"EDUC", "ELEC", "EMIG", "ENDL", "ENGA", "EVEN", "FACT", "FCOM",
			"GRAD", "IDNO", "ILLN", "IMMI", "IMMIG", "MARR", "MILI",
			"MILI_AWA", "MILI_RET", "NATI", "NATU", "NCHI", "NMR", "OCCU",
			"ORDN", "PROB", "PROP", "RELI", "RESD", "RESI", "RETI", "SLGC",
			"SSN", "TITL", "TRAV", "WAR", "WILL":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			r.Event = append(r.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "ATTR":
			r.Attribute = value

		case "FAMC":
			family := d.family(stripXref(value))
			rec := &FamilyLink{Level: level, Tag: tag, Family: family}
			r.Parents = append(r.Parents, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "FAMS":
			family := d.family(stripXref(value))
			rec := &FamilyLink{Level: level, Tag: tag, Family: family}
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

		case "CONL":
			r.CONL = value

		case "AFN":
			r.AncestralFileNumber = append(r.AncestralFileNumber, value)

		case "RFN":
			r.RecordFileNumber = value

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "_UID":
			r.UID_ = append(r.UID_, value)

		case "RIN":
			r.RIN = value

		case "EMAIL":
			r.Email = value

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
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
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

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

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
			r.Phone = append(r.Phone, value)

		case "MISC":
			r.Miscellaneous = append(r.Miscellaneous, value)

		case "_PROF":
			media := d.media(stripXref(value))
			rec := &MediaLink{Level: level, Tag: tag, Media: media}
			r.ProfilePicture_ = rec
			d.pushParser(makeMediaLinkParser(d, rec, level))

		default:
			log.Printf("unhandled Individual tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

func makeMediaLinkParser(d *Decoder, r *MediaLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		default:
			log.Printf("unhandled MediaLink tag: %d %s %s\r", level, tag, value)
		}

		return nil
	}
}

func makeMediaParser(d *Decoder, r *MediaRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "FORM":
			r.Format = value

		case "_URL":
			r.URL_ = value

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

		case "_DATE":
			r.Date_ = value

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

		case "RIN":
			r.RIN = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Media tag: %d %s %s\r", level, tag, value)
		}

		return nil
	}
}

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

		case "TYPE":
			r.NameType = value

		case "_PRIM":
			r.Primary_ = value

		case "_AKA":
			r.AKA_ = append(r.AKA_, value)

		case "NICK":
			r.Nickname = append(r.Nickname, value)

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			r.Note = append(r.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Name tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			r.Citation = append(r.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "REFN":
			rec := &UserReferenceNumberRecord{Level: level, UserReferenceNumber: value}
			r.UserReferenceNumber = append(r.UserReferenceNumber, rec)
			d.pushParser(makeUserReferenceNumberParser(d, rec, level))

		case "RIN":
			r.RIN = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Note tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makePlacePartParser(d *Decoder, r *PlacePartRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "JURI":
			r.Jurisdiction = value

		default:
			log.Printf("unhandled Place Part tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
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
			log.Printf("unhandled Place tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			log.Printf("unhandled Repository Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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
			r.Phone = append(r.Phone, value)

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

		case "RIN":
			r.RIN = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Repository tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeRoleParser(d *Decoder, r *RoleRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PRIN":
			r.Principal = value

		default:
			log.Printf("unhandled Role tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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

			case "_EVENT_DEFN":
				rec := d.eventDefinition(xref)
				r.EventDefinition_ = append(r.EventDefinition_, rec)
				d.pushParser(makeEventDefinitionParser(d, rec, level))

			case "NOTE":
				rec := d.note(xref)
				r.Note = append(r.Note, rec)
				d.pushParser(makeNoteParser(d, rec, level))

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
				log.Printf("unhandled Root tag: %d @%s@ %s %s\n", level, xref, tag, value)
			}
		} else {
			log.Printf("Not level 0 Root tag: %d @%s@ %s %s\n", level, xref, tag, value)
		}
		return nil
	}
}

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
			log.Printf("unhandled Schema tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeShortTitleParser(d *Decoder, r *ShortTitleRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "INDX":
			r.Indexed = value

		default:
			log.Printf("unhandled Short Title tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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

		case "NAME":
			r.Name = value
			d.pushParser(makeTextParser(d, &r.Name, level))

		case "TITL":
			r.Title = value
			d.pushParser(makeTextParser(d, &r.Title, level))

		case "AUTH":
			r.Author = value
			d.pushParser(makeTextParser(d, &r.Title, level))

		case "ABBR":
			r.Abbreviation = value
			d.pushParser(makeTextParser(d, &r.Title, level))

		case "PUBL":
			r.Abbreviation = value
			d.pushParser(makeTextParser(d, &r.Title, level))

		case "_PAREN":
			r.Parenthesized_ = value

		case "TEXT":
			r.Text = append(r.Text, value)
			d.pushParser(makeTextParser(d, &r.Text[len(r.Text)-1], level))

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

		case "RIN":
			r.RIN = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Source tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeSubmissionLinkParser(d *Decoder, r *SubmissionLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		default:
			log.Printf("unhandled Submission Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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

		case "RIN":
			r.RIN = value

		default:
			log.Printf("unhandled Submission tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeSubmitterLinkParser(d *Decoder, r *SubmitterLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		default:
			log.Printf("unhandled Submitter Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeSubmitterParser(d *Decoder, r *SubmitterRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

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
			r.Phone = append(r.Phone, value)

		case "EMAIL":
			r.Email = value

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

		case "RIN":
			r.RIN = value

		case "STAL":
			r.STAL = value

		case "NUMB":
			r.NUMB = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			r.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Submitter tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

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

		default:
			log.Printf("unhandled System tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

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
			log.Printf("unhandled Text tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeUserReferenceNumberParser(d *Decoder, r *UserReferenceNumberRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE":
			r.Type = value

		default:
			log.Printf("unhandled UserReferenceNumber tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

// stripXref removes value and surrounding @s from xref
func stripXref(value string) string {
	if value == "" {
		return ""
	}
	if value[0] == '@' {
		return strings.Trim(value, "@")
	}
	atIndex := strings.IndexByte(value, '@')
	if atIndex >= 0 {
		return strings.Trim(value[atIndex:], "@")
	}
	return ""
}

// stripValue removes @ bracketed xref from value
func stripValue(value string) string {
	if value == "" {
		return ""
	}
	atIndex := strings.IndexByte(value, '@')
	if atIndex >= 0 {
		return strings.Trim(value[:atIndex], " ")
	}
	return value
}

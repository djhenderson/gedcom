/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

// Package gedcom supports decoding and encoding GEDCOM files.
package gedcom

import (
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

	g := &RootRecord{
		Level:            -1,
		Place:            make(PlaceRecords, 0),
		Event:            make(EventRecords, 0),
		Individual:       make(IndividualRecords, 0),
		Family:           make(FamilyRecords, 0),
		Repository:       make(RepositoryRecords, 0),
		Source:           make(SourceRecords, 0),
		Media:            make(MediaRecords, 0),
		Note:             make(NoteRecords, 0),
		EventDefinition_: make(EventDefinitionRecords, 0),
		ChildStatus:      make(ChildStatusRecords, 0),
		Submission:       make([]*SubmissionRecord, 0),
		Submitter:        make([]*SubmitterRecord, 0),
	}

	d.refs = make(map[string]interface{})
	d.parsers = []parser{makeRootParser(d, g)}
	d.scan(g)

	return g, nil
}

func (d *Decoder) scan(g *RootRecord) {
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
func makeAddressParser(d *Decoder, a *AddressRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			a.Full = a.Full + "\n" + value

		case "CONC":
			a.Full = a.Full + value

		case "ADR1":
			a.Line1 = value

		case "ADR2":
			a.Line2 = value

		case "ADR3":
			a.Line3 = value

		case "CITY":
			a.City = value

		case "STAE":
			a.State = value

		case "POST":
			a.PostalCode = value

		case "CTRY":
			a.Country = value

		case "PHON":
			a.Phone = value

		default:
			log.Printf("unhandled Address tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeBibliographyParser(d *Decoder, c *BibliographyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "COMP":
			c.Component = append(c.Component, value)
			d.pushParser(makeTextParser(d, &c.Component[len(c.Component)-1], level))

		default:
			log.Printf("unhandled Bibliography tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeBusinessParser(d *Decoder, b *BusinessRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			b.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "PHON":
			b.Phone = append(b.Phone, value)

		case "WWW":
			b.WebSite = value

		default:
			log.Printf("unhandled Business tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeCallNumberParser(d *Decoder, b *CallNumberRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "MEDI":
			b.Media = value

		default:
			log.Printf("unhandled CallNumber tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeChangeParser(d *Decoder, e *ChangeRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "DATE":
			rec := &DateRecord{Level: level, Date: value}
			e.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			e.Note = append(e.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Change tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeCharacterSetParser(d *Decoder, e *CharacterSetRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			e.Version = value

		default:
			log.Printf("unhandled CharacterSet tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeChildStatusParser(d *Decoder, e *ChildStatusRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			e.Name = value

		default:
			log.Printf("unhandled ChildStatus tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeCitationParser(d *Decoder, c *CitationRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PAGE":
			c.Page = value

		case "REF":
			c.Reference = value

		case "QUAY":
			c.Quality = value

		case "CONS":
			c.CONS = value

		case "DIRE":
			c.Direct = value

		case "SOQU":
			c.SourceQuality = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			c.Note = append(c.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "DATA":
			rec := &DataRecord{Level: level, Data: value}
			c.Data = append(c.Data, rec)
			d.pushParser(makeDataParser(d, rec, level))

		case "TEXT":
			c.Text = append(c.Text, value)
			d.pushParser(makeTextParser(d, &c.Text[len(c.Text)-1], level))

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

func makeEventDefinitionParser(d *Decoder, s *EventDefinitionRecord, minLevel int) parser {
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

func makeEventParser(d *Decoder, e *EventRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE":
			e.Type = value

		case "NAME":
			e.Name = value

		case "_PRIM":
			e.Primary_ = value

		case "DATE":
			rec := &DateRecord{Level: level, Date: value}
			e.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "PLAC":
			rec := &PlaceRecord{Level: level, Name: value}
			e.Place = rec
			d.pushParser(makePlaceParser(d, rec, level))

		case "ROLE": // This is a kind of IndividualLink
			indi := d.individual(stripXref(value))
			rec := &RoleRecord{Level: level, Role: stripValue(value), Individual: indi}
			e.Role = append(e.Role, rec)
			d.pushParser(makeRoleParser(d, rec, level))

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			e.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "FAMC":
			family := d.family(stripXref(value))
			rec := &FamilyLink{Level: level, Tag: tag, Family: family}
			e.Parents = append(e.Parents, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "HUSB":
			husband := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: husband}
			e.Husband = rec

		case "WIFE":
			wife := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: wife}
			e.Wife = rec

		case "SPOU":
			spouse := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: spouse}
			e.Spouse = rec

		case "AGE":
			e.Age = value

		case "AGNC":
			e.Agency = value

		case "CAUS":
			e.Cause = value

		case "TEMP":
			e.Temple = value

		case "QUAY":
			e.Quality = value

		case "_UID":
			e.UID_ = append(e.UID_, value)

		case "RIN":
			e.RIN = value

		case "EMAIL":
			e.Email = value

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			e.Citation = append(e.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			e.Note = append(e.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			e.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_UPD":
			e.UpdateTime_ = value

		default:
			log.Printf("unhandled Event tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeFamilyLinkParser(d *Decoder, f *FamilyLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PEDI":
			f.Pedigree = value

		case "ADOP":
			f.Adopted = value

		case "_PRIMARY":
			f.Primary_ = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			f.Note = append(f.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Family Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeFamilyParser(d *Decoder, f *FamilyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "HUSB":
			husband := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: husband}
			f.Husband = rec

		case "WIFE":
			wife := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: wife}
			f.Wife = rec

		case "NCHI":
			pint, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Printf("NCHI = %s: %s", value, err.Error())
			}
			f.NumChildren = int(pint)

		case "CHIL":
			child := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: child}
			f.Child = append(f.Child, rec)
			d.pushParser(makeIndividualLinkParser(d, rec, level))

		case "ANUL", "CENS", "DIV", "DIVF", "ENGA", "EVEN", "MARR", "MARB",
			"MARC", "MARL", "MARS", "SLGC", "SLGS":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			f.Event = append(f.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "_UID":
			f.UID_ = append(f.UID_, value)

		case "RIN":
			f.RIN = value

		case "OBJE":
			rec := &MediaLink{Level: level, Media: d.media(stripXref(value))}
			f.Media = append(f.Media, rec)
			d.pushParser(makeMediaLinkParser(d, rec, level))

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			f.Citation = append(f.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			f.Note = append(f.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{}
			f.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_UPD":
			f.UpdateTime_ = value

		default:
			log.Printf("unhandled Family tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

func makeFootnoteParser(d *Decoder, c *FootnoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "COMP":
			c.Component = append(c.Component, value)
			d.pushParser(makeTextParser(d, &c.Component[len(c.Component)-1], level))

		default:
			log.Printf("unhandled Footnote tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeGedcomParser(d *Decoder, i *GedcomRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			i.Version = value

		case "FORM":
			i.Form = value

		default:
			log.Printf("unhandled Gedcom tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

func makeHeaderParser(d *Decoder, h *HeaderRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SOUR": // This is not a citation
			rec := &SystemRecord{Level: level, Id: value}
			h.SourceSystem = rec
			d.pushParser(makeSystemParser(d, rec, level))

		case "DEST":
			h.Destination = value

		case "DATE":
			rec := &DateRecord{Level: level, Date: value}
			h.Date = rec
			d.pushParser(makeDateParser(d, rec, level))

		case "TIME":
			h.Time = value

		case "GEDC":
			rec := &GedcomRecord{Level: level}
			h.Gedcom = rec
			d.pushParser(makeGedcomParser(d, rec, level))

		case "CHAR":
			rec := &CharacterSetRecord{Level: level, CharacterSet: value}
			h.CharacterSet = rec
			d.pushParser(makeCharacterSetParser(d, rec, level))

		case "LANG":
			h.Language = value

		case "FILE":
			h.FileName = value

		case "COPR":
			h.Copyright = value

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			h.Note = append(h.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "SUBM":
			subm := d.submitter(stripXref(value))
			rec := &SubmitterLink{Level: level, Submitter: subm}
			h.Submitter = append(h.Submitter, rec)
			d.pushParser(makeSubmitterLinkParser(d, rec, level))

		case "SUBN":
			subn := d.submission(stripXref(value))
			rec := &SubmissionLink{Level: level, Submission: subn}
			h.Submission = append(h.Submission, rec)
			d.pushParser(makeSubmissionLinkParser(d, rec, level))

		case "SCHEMA":
			rec := &SchemaRecord{Level: level}
			h.Schema = rec
			d.pushParser(makeSchemaParser(d, rec, level))

		case "_ROOT":
			root := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: root}
			h.Root_ = rec

		default:
			log.Printf("unhandled Header tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

func makeHistoryParser(d *Decoder, n *HistoryRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			n.History = n.History + "\n" + value

		case "CONC":
			n.History = n.History + value

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			n.Citation = append(n.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		default:
			log.Printf("unhandled History tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeIndividualLinkParser(d *Decoder, n *IndividualLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SLGC":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			n.Event = append(n.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			n.Note = append(n.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Individual Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeIndividualParser(d *Decoder, i *IndividualRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			rec := &NameRecord{Level: level, Name: value}
			i.Name = append(i.Name, rec)
			d.pushParser(makeNameParser(d, rec, level))

		case "SEX":
			i.Sex = value

		case "BIRT", "CHR", "DEAT", "BURI", "CREM", "ADOP", "BAPM", "BARM",
			"BASM", "BLES", "CHRA", "CONF", "FCOM", "ORDN", "NATU", "EMIG",
			"IMMI", "CENS", "PROB", "WILL", "GRAD", "RETI", "BAPL", "EDUC",
			"ENDL", "NATI", "OCCU", "RELI", "RESI", "TITL", "ENGA", "MARR",
			"IMMIG", "ILLN", "TRAV", "RESD", "MILI", "WAR", "MILI_AWA",
			"MILI_RET", "ELEC", "SLGC", "EVEN":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			i.Event = append(i.Event, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "CAST", "DSCR", "IDNO", "NCHI", "NMR", "PROP", "SSN", "FACT":
			rec := &EventRecord{Level: level, Tag: tag, Value: value}
			i.Attribute = append(i.Attribute, rec)
			d.pushParser(makeEventParser(d, rec, level))

		case "FAMC":
			family := d.family(stripXref(value))
			rec := &FamilyLink{Level: level, Tag: tag, Family: family}
			i.Parents = append(i.Parents, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "FAMS":
			family := d.family(stripXref(value))
			rec := &FamilyLink{Level: level, Tag: tag, Family: family}
			i.Family = append(i.Family, rec)
			d.pushParser(makeFamilyLinkParser(d, rec, level))

		case "OBJE":
			rec := &MediaLink{Level: level, Media: d.media(stripXref(value))}
			i.Media = append(i.Media, rec)
			d.pushParser(makeMediaLinkParser(d, rec, level))

		case "HEAL":
			i.Health = value

		case "QUAY":
			i.Quality = value

		case "LVG":
			i.Living = value

		case "AFN":
			i.AFN = append(i.AFN, value)

		case "RFN":
			i.RefNumber = value

		case "REFN":
			rec := &ReferenceNumberRecord{Level: level, ReferenceNumber: value}
			i.ReferenceNumber = rec
			d.pushParser(makeReferenceNumberParser(d, rec, level))

		case "_UID":
			i.UID_ = append(i.UID_, value)

		case "RIN":
			i.RIN = value

		case "EMAIL":
			i.Email = value

		case "WWW":
			i.WebSite = value

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			i.Address = append(i.Address, rec)
			d.pushParser(makeAddressParser(d, rec, level))

		case "HIST":
			rec := &HistoryRecord{Level: level, History: value}
			i.History = append(i.History, rec)
			d.pushParser(makeHistoryParser(d, rec, level))

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			i.Citation = append(i.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			i.Note = append(i.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			i.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		case "_UPD":
			i.UpdateTime_ = value

		case "ALIA":
			i.Alias = value

		case "FATH":
			father := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: father}
			i.Father = rec

		case "MOTH":
			mother := d.individual(stripXref(value))
			rec := &IndividualLink{Level: level, Tag: tag, Individual: mother}
			i.Mother = rec

		case "PHON":
			i.Phone = append(i.Phone, value)

		case "MISC":
			i.Miscellaneous = append(i.Miscellaneous, value)

		default:
			log.Printf("unhandled Individual tag: %d %s %s\n", level, tag, value)
		}
		return nil
	}
}

func makeMediaLinkParser(d *Decoder, n *MediaLink, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "FORM":
			n.Format = value

		case "FILE":
			n.FileName = value

		case "TITL":
			n.Title = value

		case "DATE":
			n.Date = value

		case "AUTH":
			n.Author = value

		case "TEXT":
			n.Text = value
			d.pushParser(makeTextParser(d, &n.Text, level))

		case "_PLACE",
			"_PRIM_CUTOUT",
			"_POSITION",
			"_PHOTO_RIN",
			"_FILESIZE",
			"_PRIM",
			"_CUTOUT",
			"_DATE":
			// TODO

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			n.Note = append(n.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Media Link tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeMediaParser(d *Decoder, n *MediaRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "FORM":
			n.Format = value

		case "_URL":
			n.URL_ = value

		case "FILE":
			n.FileName = value

		case "TITL":
			n.Title = value

		case "DATE":
			n.Date = value

		case "AUTH":
			n.Author = value

		case "TEXT":
			n.Text = value
			d.pushParser(makeTextParser(d, &n.Text, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			n.Note = append(n.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "_DATE":
			n.Date_ = value

		case "_ASTID":
			n.AstId_ = value

		case "_ASTTYP":
			n.AstType_ = value

		case "_ASTDESC":
			n.AstDesc_ = value

		case "_ASTPERM":
			n.AstPerm_ = value

		case "_ASTUPPID":
			n.AstUpPid_ = value

		default:
			log.Printf("unhandled Media tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeNameParser(d *Decoder, n *NameRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NPFX":
			n.Prefix = value

		case "GIVN":
			n.GivenName = value

		case "_MIDN":
			n.MiddleName_ = value

		case "SURN":
			n.Surname = value

		case "NSFX":
			n.Suffix = value

		case "_PGVN":
			n.PreferedGivenName_ = value

		case "_PRIM":
			n.Primary_ = value

		case "_AKA":
			n.AKA_ = append(n.AKA_, value)

		case "NICK":
			n.Nickname = append(n.Nickname, value)

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			n.Citation = append(n.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			n.Note = append(n.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		default:
			log.Printf("unhandled Name tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeNoteParser(d *Decoder, n *NoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "CONT":
			n.Note = n.Note + "\n" + value

		case "CONC":
			n.Note = n.Note + value

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			n.Citation = append(n.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		default:
			log.Printf("unhandled Note tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makePlacePartParser(d *Decoder, p *PlacePartRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "JURI":
			p.Jurisdiction = value

		default:
			log.Printf("unhandled Place Part tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makePlaceParser(d *Decoder, p *PlaceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "PLAS":
			p.ShortName = value

		case "PLAM":
			p.Modifier = value

		case "PLA0", "PLA1", "PLA2", "PLA3", "PLA4":
			rec := &PlacePartRecord{Level: level, Tag: tag, Part: value}
			p.Parts = append(p.Parts, rec)
			d.pushParser(makePlacePartParser(d, rec, level))

		case "SOUR":
			sour := d.source(stripXref(value))
			rec := &CitationRecord{Level: level, Value: stripValue(value), Source: sour}
			p.Citation = append(p.Citation, rec)
			d.pushParser(makeCitationParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			p.Note = append(p.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			p.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Place tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeReferenceNumberParser(d *Decoder, r *ReferenceNumberRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "TYPE":
			r.Type = value

		default:
			log.Printf("unhandled ReferenceNumber tag: %d %s %s\n", level, tag, value)
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

		case "WWW":
			r.WebSite = value

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

func makeRootParser(d *Decoder, g *RootRecord) parser {
	return func(level int, tag string, value string, xref string) error {
		//log.Println(level, tag, value, xref)
		if level == 0 {
			// cases ordered approx. by frequency of appearance
			g.Level = level // always zero
			switch tag {

			case "INDI":
				rec := d.individual(xref)
				g.Individual = append(g.Individual, rec)
				d.pushParser(makeIndividualParser(d, rec, level))

			case "FAM":
				rec := d.family(xref)
				g.Family = append(g.Family, rec)
				d.pushParser(makeFamilyParser(d, rec, level))

			case "SOUR":
				rec := d.source(xref)
				g.Source = append(g.Source, rec)
				d.pushParser(makeSourceParser(d, rec, level))

			case "EVEN":
				rec := d.event(xref)
				rec.Tag = tag
				rec.Value = value
				g.Event = append(g.Event, rec)
				d.pushParser(makeEventParser(d, rec, level))

			case "PLAC":
				rec := d.place(xref)
				g.Place = append(g.Place, rec)
				d.pushParser(makePlaceParser(d, rec, level))

			case "REPO":
				rec := d.repository(xref)
				g.Repository = append(g.Repository, rec)
				d.pushParser(makeRepositoryParser(d, rec, level))

			case "_EVENT_DEFN":
				rec := d.eventDefinition(xref)
				g.EventDefinition_ = append(g.EventDefinition_, rec)
				d.pushParser(makeEventDefinitionParser(d, rec, level))

			case "NOTE":
				rec := d.note(xref)
				g.Note = append(g.Note, rec)
				d.pushParser(makeNoteParser(d, rec, level))

			case "OBJE":
				rec := d.media(xref)
				g.Media = append(g.Media, rec)
				d.pushParser(makeMediaParser(d, rec, level))

			case "CSTA":
				rec := d.childStatus(xref)
				g.ChildStatus = append(g.ChildStatus, rec)
				d.pushParser(makeChildStatusParser(d, rec, level))

			case "HEAD":
				rec := d.header(xref)
				g.Header = rec
				d.pushParser(makeHeaderParser(d, rec, level))

			case "SUBM":
				rec := d.submitter(xref)
				g.Submitter = append(g.Submitter, rec)
				d.pushParser(makeSubmitterParser(d, rec, level))

			case "SUBN":
				rec := d.submission(xref)
				g.Submission = append(g.Submission, rec)
				d.pushParser(makeSubmissionParser(d, rec, level))

			case "TRLR":
				rec := d.trailer(xref)
				g.Trailer = rec
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

func makeSchemaParser(d *Decoder, s *SchemaRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "AGER", "AKA", "CLASS", "COMP", "CONT", "DETAIL1", "DETAIL2",
			"EVEN", "LANG", "NAME", "NOTE", "PERI", "POSB", "POSF", "PREB",
			"PREF", "PRIN", "ROLE", "SEX", "SOUR", "STYL":
			// TODO

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

func makeSourceParser(d *Decoder, s *SourceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			s.Name = value
			d.pushParser(makeTextParser(d, &s.Name, level))

		case "TITL":
			s.Title = value
			d.pushParser(makeTextParser(d, &s.Title, level))

		case "AUTH":
			s.Author = value
			d.pushParser(makeTextParser(d, &s.Title, level))

		case "ABBR":
			s.Abbreviation = value
			d.pushParser(makeTextParser(d, &s.Title, level))

		case "PUBL":
			s.Abbreviation = value
			d.pushParser(makeTextParser(d, &s.Title, level))

		case "_PAREN":
			s.Parenthesized_ = value

		case "TEXT":
			s.Text = append(s.Text, value)
			d.pushParser(makeTextParser(d, &s.Text[len(s.Text)-1], level))

		case "DATA":
			rec := &DataRecord{Level: level, Data: value}
			s.Data = rec
			d.pushParser(makeDataParser(d, rec, level))

		case "SHAU":
			s.ShortAuthor = value

		case "SHTI":
			rec := &ShortTitleRecord{Level: level, ShortTitle: value}
			s.ShortTitle = rec // append(s.ShortTitle, rec)
			d.pushParser(makeShortTitleParser(d, rec, level))

		case "FOOT":
			rec := &FootnoteRecord{Level: level, Value: value}
			s.Footnote = rec // append(s.Footnote, rec)
			d.pushParser(makeFootnoteParser(d, rec, level))

		case "BIBL":
			rec := &BibliographyRecord{Level: level, Value: value}
			s.Bibliography = rec // append(s.Media, rec)
			d.pushParser(makeBibliographyParser(d, rec, level))

		case "REPO":
			repo := d.repository(stripXref(value))
			rec := &RepositoryLink{Level: level, Repository: repo}
			s.Repository = rec
			d.pushParser(makeRepositoryLinkParser(d, rec, level))

		case "OBJE":
			rec := &MediaLink{Level: level, Media: d.media(stripXref(value))}
			s.Media = append(s.Media, rec)
			d.pushParser(makeMediaLinkParser(d, rec, level))

		case "NOTE":
			rec := &NoteRecord{Level: level, Note: value}
			s.Note = append(s.Note, rec)
			d.pushParser(makeNoteParser(d, rec, level))

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			s.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Source tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeSubmissionLinkParser(d *Decoder, s *SubmissionLink, minLevel int) parser {
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

func makeSubmissionParser(d *Decoder, s *SubmissionRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SUBM":
			rec := d.submitter(xref)
			s.Submitter = rec
			d.pushParser(makeSubmitterParser(d, rec, level))

		case "FAMF":
			s.FamilyFileName = value

		case "TEMP":
			s.Temple = value

		case "ANCE":
			s.Ancestors = value

		case "DESC":
			s.Descendents = value

		case "ORDI":
			s.Ordinance = value

		default:
			log.Printf("unhandled Submission tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeSubmitterLinkParser(d *Decoder, s *SubmitterLink, minLevel int) parser {
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

func makeSubmitterParser(d *Decoder, s *SubmitterRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "NAME":
			s.Name = value
			//d.pushParser(makeTextParser(d, &s.Name, level))

		case "ADDR":
			rec := &AddressRecord{Level: level, Full: value}
			s.Address = rec
			d.pushParser(makeAddressParser(d, rec, level))

		case "CTRY":
			s.Country = value

		case "PHON":
			s.Phone = append(s.Phone, value)

		case "EMAIL":
			s.Email = value

		case "WWW":
			s.WebSite = value

		case "LANG":
			s.Language = value

		case "RIN":
			s.RIN = value

		case "STAL":
			s.STAL = value

		case "NUMB":
			s.NUMB = value

		case "CHAN":
			rec := &ChangeRecord{Level: level}
			s.Change = rec
			d.pushParser(makeChangeParser(d, rec, level))

		default:
			log.Printf("unhandled Submitter tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

func makeSystemParser(d *Decoder, s *SystemRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "VERS":
			s.Version = value

		case "NAME":
			s.ProductName = value

		case "CORP":
			rec := &BusinessRecord{Level: level, BusinessName: value}
			s.Business = rec
			d.pushParser(makeBusinessParser(d, rec, level))

		case "DATA":
			rec := &DataRecord{Level: level, Data: value}
			s.SourceData = rec
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
			*s = *s + "\n" + value

		case "CONC":
			*s = *s + value

		default:
			log.Printf("unhandled Text tag: %d %s %s\n", level, tag, value)
		}

		return nil
	}
}

// stripXref removes value and surrounding @s
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

/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

// AddressRecord represents an address record
type AddressRecord struct {
	Level      int    // ..ADDR level
	Full       string // ..ADDR value
	Line1      string // ..ADDR.ADR1
	Line2      string // ..ADDR.ADR2
	Line3      string // ..ADDR.ADR3
	City       string // ..ADDR.CITY
	State      string // ..ADDR.STAE
	PostalCode string // ..ADDR.POST
	Country    string // ..ADDR.CTRY
	Phone      string // ..ADDR.PHON
}

// AddressRecords represents a slice of address records
type AddressRecords []*AddressRecord

// BibliographyRecord represents a bibliography record
type BibliographyRecord struct {
	Level     int      // ..BIBL level
	Value     string   // ..BIBL value
	Component []string // ..BIBL.COMP
}

// BlobRecord represents a binary large object record
type BlobRecord struct {
	Level int    // ..BLOB level
	Data  string // ..BLOB.CONT
}

// BusinessRecord represents a business record
type BusinessRecord struct {
	Level        int            // ..HEAD.SOUR.CORP level
	BusinessName string         // ..HEAD.SOUR.CORP value
	Address      *AddressRecord // ..HEAD.SOUR.CORP.ADDR
	Phone        []string       // ..HEAD.SOUR.CORP.PHON
	WebSite      string         // ..HEAD.SOUR.CORP.WWW
}

// CallNumberRecord represents a call number record
type CallNumberRecord struct {
	Level      int    // ..REPO.CALN level
	CallNumber string // ..REPO.CALN value
	Media      string // ..REPO.CALN.MEDI
}

// ChangeRecord represents a change record
type ChangeRecord struct {
	Level int         // ..CHAN level
	Date  *DateRecord // ..CHAN.DATE
	Note  NoteRecords // ..CHAN.NOTE
}

// CharacterSetRecord represents a character set record
type CharacterSetRecord struct {
	Level        int    // ..CHAR level
	CharacterSet string // ..CHAR value
	Version      string // ..CHAR.VERS
}

// ChildStatusRecord represents a child status record
type ChildStatusRecord struct {
	Level int    // CSTA level; always 0
	Xref  string // xref_id of 0 level CSTA
	Name  string // CSTA.NAME
}

// ChildStatusRecords represents a slice of child status records
type ChildStatusRecords []*ChildStatusRecord

// CitationRecord represents a link to a source record
type CitationRecord struct {
	Level         int           // ..SOUR level; not 0
	Value         string        // value of ..SOUR excluding xref
	Source        *SourceRecord // linked source
	Page          string        // ..SOUR.PAGE
	Reference     string        // ..SOUR.REF
	Event         EventRecords  // .. SOUR.EVEN
	Data          DataRecords   // ..SOUR.DATA
	Text          []string      // ..SOUR.TEXT
	Quality       string        // ..SOUR.QUAY
	Media         MediaRecords  // ..SOUR.OBJE
	CONS          string        // ..SOUR.CONS
	Direct        string        // ..SOUR.DIRE
	SourceQuality string        // ..SOUR.SOQU
	Note          NoteRecords   // ..SOUR.NOTE
}

// CitationRecords represents a slice of links to a source
type CitationRecords []*CitationRecord

// DataRecord represents a data record
type DataRecord struct {
	Level     int          // ..DATA level
	Data      string       // value of ..DATA
	Date      string       // ..DATA.DATE
	Copyright string       // ..DATA.COPR
	Text      []string     // ..DATA.TEXT
	Event     EventRecords // ..DATA.EVEN
	Agency    string       // ..DATA.AGNC
	Note      NoteRecords  // ..DATA.NOTE
}

// DataRecords represents a slice of data records
type DataRecords []*DataRecord

// DateRecord represents a date
type DateRecord struct {
	Level int      // ..DATE level
	Date  string   // ..DATE value
	Time  string   // ..DATE.TIME value
	Text  []string // ..DATE.TEXT
	Day   string   // ..DATA.DATD
	Month string   // ..DATA.DATM
	Year  string   // ..DATA.DATY
	Full  string   // ..DATA.DATF
	Short string   // ..DATA.DATS
}

// EventDefinitionRecord represens a GEDCOM event definition record.
type EventDefinitionRecord struct {
	Level int    // Level 0 _EVENT_DEFN
	Xref  string // Level 0 xref
}

// EventDefinitionRecords represents a slice of event definition records.
type EventDefinitionRecords []*EventDefinitionRecord

// EventRecord represents a GEDCOM event record.
type EventRecord struct {
	Level       int             // ..EVEN level; 0 or higher
	Xref        string          // xref_id of level 0 ..EVEN
	Tag         string          // Event tag EVEN or BIRT or ...
	Value       string          // Event value
	Type        string          // ..EVEN.TYPE
	Name        string          // ..EVEN.NAME
	Primary_    string          // ..EVEN._PRIM
	Date        *DateRecord     // ..EVEN.DATE
	Place       *PlaceRecord    // ..EVEN.PLAC
	Role        RoleRecords     // ..EVEN.ROLE
	Address     *AddressRecord  // ..EVEN.ADDR
	Phone       []string        // ..EVEN.PHON
	Parents     FamilyLinks     // ..EVEN.FAMC
	Husband     *IndividualLink // ..EVEN.HUSB
	Wife        *IndividualLink // ..EVEN.WIFE
	Spouse      *IndividualLink // ..EVEN.SPOU
	Age         string          // ..EVEN.AGE
	Agency      string          // ..EVEN.AGNC
	Cause       string          // ..EVEN.CAUS
	Temple      string          // ..EVEN.TEMP
	Quality     string          // ..EVEN.QUAY
	Status      string          // ..EVEN.STAT
	UID_        []string        // ..EVEN._UID
	RIN         string          // ..EVEN.RIN
	Email       string          // ..EVEN.EMAIL
	Media       MediaRecords    // ..EVEN.OBJE
	Citation    CitationRecords // ..EVEN.SOUR
	Note        NoteRecords     // ..EVEN.NOTE
	Change      *ChangeRecord   // ..EVEN.CHAN
	UpdateTime_ string          // ..EVEN._UPD
}

// EventRecords represents a slice of event records.
type EventRecords []*EventRecord

// FamilyLink represents a GEDCOM link to a family record.
type FamilyLink struct {
	Level    int           //  level
	Tag      string        // tag from INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Family   *FamilyRecord // target of INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Pedigree string        // INDI.FAMC.PEDI or ...
	Adopted  string        // INDI.FAMC.ADOP or ...
	Primary_ string        // INDI.FAMC._PRIMARY or ...
	Note     NoteRecords   // INDI.FAMC.NOTE or ..
}

// FamilyLinks represents a slice of links to family records.
type FamilyLinks []*FamilyLink

// FamilyRecord represents a GEDCOM family record.
type FamilyRecord struct {
	Level               int                        // FAM level; only 0
	Xref                string                     // xref_id of FAM
	Husband             *IndividualLink            // FAM.HUSB
	Wife                *IndividualLink            // FAM.WIFE
	NumChildren         int                        // FAM.NCHI
	Child               IndividualLinks            // FAM.CHIL
	Event               EventRecords               // FAM.MARR, FAM.EVEN
	UID_                []string                   // FAM._UID
	RIN                 string                     // FAM.RIN
	UserReferenceNumber UserReferenceNumberRecords // FAM.REFN
	Media               MediaRecords               // FAM.OBJE
	Citation            CitationRecords            // FAM.SOUR
	Note                NoteRecords                // FAM.NOTE
	Submitter           SubmitterLinks             // FAM.SUBM
	Change              *ChangeRecord              // FAM.CHAN
	UpdateTime_         string                     // FAM._UPD
}

// FamilyRecords represents a slice of family records.
type FamilyRecords []*FamilyRecord

// FootnoteRecord represents a footnote record
type FootnoteRecord struct {
	Level     int      // ..FOOT level
	Value     string   // ..FOOT value
	Component []string // ..FOOT.COMP
}

// GedcomRecord represents a gedcom record
type GedcomRecord struct {
	Level   int    // HEAD.GEDC level
	Version string // HEAD.GEDC.VERS
	Form    string // HEAD.GEDC.FORM
}

// HeaderRecord represents a header record
type HeaderRecord struct {
	Level        int                 // HEAD level; always 0
	Xref         string              // Fake id; set to HEAD
	SourceSystem *SystemRecord       // HEAD.SOUR
	Destination  string              // HEAD.DEST
	Date         *DateRecord         // HEAD.DATE
	Time         string              // HEAD.TIME
	FileName     string              // HEAD.FILE
	Gedcom       *GedcomRecord       // HEAD.GEDC
	CharacterSet *CharacterSetRecord // HEAD.CHAR
	Language     string              // HEAD.LANG
	Copyright    string              // HEAD.COPR
	Place        *PlaceRecord        // HEAD.PLAC
	Root_        *IndividualLink     // HEAD._ROOT
	Note         NoteRecords         // HEAD.NOTE
	Submitter    SubmitterLinks      // HEAD.SUBM
	Submission   SubmissionLinks     // HEAD.SUBN
	Schema       *SchemaRecord       // HEAD.SCHEMA
}

// HistoryRecord represents a history record
type HistoryRecord struct {
	Level    int             // ..HIST level
	History  string          // ..HIST value
	Citation CitationRecords // ..HIST.SOUR
}

// HistoryRecords represents a slice of history records
type HistoryRecords []*HistoryRecord

// IndividualLink represents a link to an individual record
type IndividualLink struct {
	Level        int               // ..INDI level
	Tag          string            // tag from FAM.HUSB or FAM.WIFE or FAM.CHILD or INDI.ASSO
	Individual   *IndividualRecord // target of FAM.HUSB or FAM.WIFE or FAM.CHILD or INDI.ASSOC
	Relationship string            // INDI.ASSO.RELA
	Event        EventRecords      // FAM.CHILD.SLGC
	Citation     CitationRecords   // INDI.ASSO.SOUR
	Note         NoteRecords       // FAM.HUSB.NOTE or FAM.WIFE.NOTE or FAM.CHILD.NOTE
}

// IndividualLinks represents a slice of links to individual records
type IndividualLinks []*IndividualLink

// IndividualRecord represents an individual record
type IndividualRecord struct {
	Level               int                        // INDI level; always 0
	Xref                string                     // xref_id of INDI
	Name                NameRecords                // INDI.NAME
	Restriction         string                     // INDI.RESN
	Sex                 string                     // INDI.SEX
	Event               EventRecords               // INDI.BIRT, INDI.CHR, INDI.DEAT, INDI.BURI, INDI.EVEN
	Attribute           string                     // INDI.ATTR
	Parents             FamilyLinks                // INDI.FAMC
	Family              FamilyLinks                // INDI.FAMS
	Address             AddressRecords             // INDI.ADDR
	Phone               []string                   // INDI.PHON
	Media               MediaRecords               // INDI.OBJE
	Health              string                     // INDI.HEAL
	History             HistoryRecords             // INDI.HIST
	Quality             string                     // INDI.QUAY
	Living              string                     // INDI.LVG
	CONL                string                     // INDI.CONL
	AncestralFileNumber []string                   // INDI.AFN
	RecordFileNumber    string                     // INDI.RFN
	UserReferenceNumber UserReferenceNumberRecords // INDI.REFN
	RIN                 string                     // INDI.RIN
	UID_                []string                   // INDI._UID
	Email               string                     // INDI.EMAIL
	WebSite             string                     // INDI.WWW
	Citation            CitationRecords            // INDI.SOUR
	Note                NoteRecords                // INDI.NOTE
	Associated          IndividualLinks            // INDI.ASSO
	Submitter           SubmitterLinks             // INDI.SUBM
	ANCI                SubmitterLinks             // INDI.ANCI
	DESI                SubmitterLinks             // INDI.DESI
	Change              *ChangeRecord              // INDI.CHAN
	UpdateTime_         string                     // INDI._UPD
	Alias               string                     // INDI.ALIA
	Father              *IndividualLink            // INDI.FATH
	Mother              *IndividualLink            // INDI.MOTH
	Miscellaneous       []string                   // INDI.MISC
}

// IndividualRecords represents a slice of individual records
type IndividualRecords []*IndividualRecord

// MediaRecord represents a media record
type MediaRecord struct {
	Level               int                        // ..OBJE level
	Xref                string                     // xref_id of 0 level OBJE
	Format              string                     // OBJE.FORM
	URL_                string                     // OBJE._URL
	FileName            string                     // OBJE.FILE
	Title               string                     // OBJE.TITL
	Date                string                     // OBJE.DATE
	Author              string                     // OBJE.AUTH
	Text                string                     // OBJE.TEXT
	Note                NoteRecords                // OBJE.NOTE
	Date_               string                     // OBJE._DATE
	AstId_              string                     // OBJE._ASTID
	AstType_            string                     // OBJE._ASTTYP
	AstDesc_            string                     // OBJE._ASTDESC
	AstPerm_            string                     // OBJE._ASTPERM
	AstUpPid_           string                     // OBJE._ASTUPPID
	BinaryLargeObject   *BlobRecord                // OBJE.BLOB
	UserReferenceNumber UserReferenceNumberRecords // OBJE.REFN
	RIN                 string                     // OBJE.RIN
	Change              *ChangeRecord              // INDI.CHAN
}

// MediaRecords represents a slice of media records
type MediaRecords []*MediaRecord

// NameRecord represents a name record
type NameRecord struct {
	Level              int             // ..NAME level
	Name               string          // ..NAME value
	Prefix             string          // ..NAME.NPFX
	GivenName          string          // ..NAME.GIVN
	MiddleName_        string          // ..NAME._MIDN
	SurnamePrefix      string          // ..NAME.SPFX
	Surname            string          // ..NAME.SURN
	Suffix             string          // ..NAME.NSFX
	PreferedGivenName_ string          // ..NAME._PGVN
	Primary_           string          // ..NAME._PRIM
	AKA_               []string        // ..NAME._AKA
	Nickname           []string        // ..NAME.NICK
	Citation           CitationRecords // ..NAME.SOUR
	Note               NoteRecords     // ..NAME.NOTE
}

// NameRecords represents a slice of name records
type NameRecords []*NameRecord

// NoteRecord represents a note record
type NoteRecord struct {
	Level               int                        // ..NOTE level
	Xref                string                     // .. xref_value of 0 level NOTE
	Note                string                     // ..NOTE value
	Citation            CitationRecords            // ..NOTE.SOUR
	UserReferenceNumber UserReferenceNumberRecords // ..NOTE.REFN
	RIN                 string                     // ..NOTE.RIN
	Change              *ChangeRecord              // ..SOUR.CHAN
}

// NoteRecords represents a slice of note records
type NoteRecords []*NoteRecord

// PlacePartRecord represents a place part record
type PlacePartRecord struct {
	Level        int    // ..PLAC.PLAn level, n=0..4
	Tag          string // ..PLAC.PLAn tag
	Part         string // ..PLAC.PLAn value
	Jurisdiction string // ..PLAC.PLAn.JURI
}

// PlacePartRecords represents a slice of place part records
type PlacePartRecords []*PlacePartRecord

// PlaceRecord represents a place record
type PlaceRecord struct {
	Level     int              // ..PLAC level; 0 or higher
	Xref      string           // xref_id of 0 level PLAC
	Name      string           // ..PLAC value
	Form      string           // ..PLAC.FORM
	ShortName string           // ..PLAC.PLAS
	Modifier  string           // ..PLAC.PLAM
	Parts     PlacePartRecords // ..PLAC.PLAn n=0..4
	Citation  CitationRecords  // ..PLAC.SOUR
	Note      NoteRecords      // ..PLAC.NOTE
	Change    *ChangeRecord    // ..PLAC.CHAN
}

// PlaceRecords represents a slice of place records
type PlaceRecords []*PlaceRecord

// RepositoryLink represents a link to a repository record
type RepositoryLink struct {
	Level      int               // ..REPO level
	Xref       string            // xref_id of 0 level REPO
	Repository *RepositoryRecord // The linked repository
	CallNumber *CallNumberRecord // ..REPO.CALN
	Note       NoteRecords       // ..REPO.NOTE
}

// RepositoryLinks represents a slice of links to repository records
type RepositoryLinks []*RepositoryLink

// RepositoryRecord represents a repository record
type RepositoryRecord struct {
	Level               int                        // REPO level; always 0
	Xref                string                     // xref_id of 0 level REPO
	Name                string                     // REPO.NAME
	Address             *AddressRecord             // REPO.ADDR
	Phone               []string                   // REPO.PHON
	WebSite             string                     // REPO.WWW
	UserReferenceNumber UserReferenceNumberRecords // REPO.REFN
	RIN                 string                     // REPO.RIN
	Note                NoteRecords                // REPO.NOTE
	Change              *ChangeRecord              // REPO.CHAN
}

// RepositoryRecords represents a slice of repository records
type RepositoryRecords []*RepositoryRecord

// RoleRecord represents a role record
type RoleRecord struct {
	Level      int               // ..ROLE level
	Role       string            // ..ROLE no-ref value
	Individual *IndividualRecord // ..ROLE ref value
	Principal  string            // ..ROLE.PRIN
}

// RoleRecords represents a slice of role records
type RoleRecords []*RoleRecord

// RootRecord represents a the root record of a GEDCOM file.
type RootRecord struct {
	Level            int                    // root level
	Header           *HeaderRecord          // HEAD
	Submitter        SubmitterRecords       // SUBM
	Submission       SubmissionRecords      // SUBN
	Place            PlaceRecords           // PLAC
	Event            EventRecords           // EVEN
	Individual       IndividualRecords      // INDI
	Family           FamilyRecords          // FAM
	Repository       RepositoryRecords      // REPO
	Source           SourceRecords          // SOUR
	Media            MediaRecords           // OBJE
	Note             NoteRecords            // NOTE
	EventDefinition_ EventDefinitionRecords // _EVENT_DEFN
	ChildStatus      ChildStatusRecords     // CSTA
	Trailer          *TrailerRecord         // TRLR
}

// SchemaRecord represents a schema record
type SchemaRecord struct {
	Level int      // ..SCHEMA level
	Data  []string // schema data
}

// ShortTitleRecord represents a short title record
type ShortTitleRecord struct {
	Level      int    // ..SHTI level
	ShortTitle string // ..SHTI value
	Indexed    string // ..SHTI.INDX
}

// SourceRecord represents a source record.
// Note: A CitationRecord is a link to a source record
type SourceRecord struct {
	Level               int                        // ..SOUR level; always 0
	Xref                string                     // xref_id of 0 level SOUR
	Value               string                     // ..SOUR value
	Name                string                     // ..SOUR.NAME
	Title               string                     // ..SOUR.TITL
	Author              string                     // ..SOUR.AUTH
	Abbreviation        string                     // ..SOUR.ABBR
	Publication         string                     // ..SOUR.PUBL
	Parenthesized_      string                     // ..SOUR._PAREN (PAF5)
	Text                []string                   // ..SOUR.TEXT
	Data                *DataRecord                // ..SOUR.DATA
	Footnote            *FootnoteRecord            // ..SOUR.FOOT
	Bibliography        *BibliographyRecord        // ..SOUR.BIBL
	Repository          *RepositoryLink            // ..SOUR.REPO
	UserReferenceNumber UserReferenceNumberRecords // ..SOUR.REFN
	RIN                 string                     // ..SOUR.RIN
	ShortAuthor         string                     // ..SOUR.SHAU
	ShortTitle          *ShortTitleRecord          // ..SOUR.SHTI
	Media               MediaRecords               // ..SOUR.OBJE
	Note                NoteRecords                // ..SOUR.NOTE
	Change              *ChangeRecord              // ..SOUR.CHAN
}

// SourceRecords represents a slice of source records
type SourceRecords []*SourceRecord

// SubmissionLink represents a link to a submission record
type SubmissionLink struct {
	Level      int               // ..SUBN level
	Submission *SubmissionRecord // target of ..SUBN
}

// SubmissionLinks represents a slice of links to submission records
type SubmissionLinks []*SubmissionLink

// SubmissionRecord represents a submission record.
type SubmissionRecord struct {
	Level          int              // SUBN level; always 0
	Xref           string           // xref_id of SUBN
	Submitter      *SubmitterRecord // SUBN.SUBM
	FamilyFileName string           // SUBN.FAMF
	Temple         string           // SUBN.TEMP
	Ancestors      string           // SUBN.ANCE
	Descendents    string           // SUBN.DESC
	Ordinance      string           // SUBN.ORDI
	RIN            string           // SUBN.RIN
}

// SubmissionRecords represents a slice of submission records
type SubmissionRecords []*SubmissionRecord

// SubmitterLink represents a link to a submitter record
type SubmitterLink struct {
	Level     int              // ..SUBM level
	Tag       string           // ..SUBM link tag
	Submitter *SubmitterRecord // target of ..SUBM
}

// SubmitterLinks represents a slice of links to submitter records
type SubmitterLinks []*SubmitterLink

// SubmitterRecord represents a submitter record
type SubmitterRecord struct {
	Level            int            // SUBM level; always 0
	Xref             string         // xref_id of SUBM
	Name             string         // SUBM.NAME
	Address          *AddressRecord // SUBM.ADDR
	Country          string         // SUBM.CTRY
	Phone            []string       // SUBM.PHON
	Email            string         // SUBM.EMAIL
	WebSite          string         // SUBM.WWW
	Language         string         // SUBM.LANG
	Media            MediaRecords   // SUBM.OBJE
	RecordFileNumber string         // SUBM.RFN
	STAL             string         // SUBM.STAL
	NUMB             string         // SUBM.NUMB
	RIN              string         // SUBM.RIN
	Change           *ChangeRecord  // SUBM.CHAN
}

// SubmitterRecords represents a slice of submitter records
type SubmitterRecords []*SubmitterRecord

// SystemRecord represents a system record
type SystemRecord struct {
	Level       int             // HEAD.SOUR level
	SystemName  string          // HEAD.SOUR value
	Version     string          // HEAD.SOUR.VERS
	ProductName string          // HEAD.SOUR.NAME
	Business    *BusinessRecord // HEAD.SOUR.CORP
	SourceData  *DataRecord     // HEAD.SOUR.DATA
}

// TrailerRecord represents a trailer record
type TrailerRecord struct {
	Level int    // ..TRLR level; always 0
	Xref  string // Fake id: set to TRLR
}

// UserReferenceNumberRecord represents a user reference number record
type UserReferenceNumberRecord struct {
	Level               int    // ..REFN level
	UserReferenceNumber string // ..REFN value
	Type                string // ..REFN.TYPE
}

// UserReferenceNumberRecords represents a slice of user reference number records.
type UserReferenceNumberRecords []*UserReferenceNumberRecord

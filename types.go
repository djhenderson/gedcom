/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

type AddressRecord struct {
	Level      int    // ..ADDR level
	Full       string // ..ADDR value
	Line1      string // ..ADDR.ADR1
	Line2      string // ..ADDR.ADR2
	City       string // ..ADDR.CITY
	State      string // ..ADDR.STAE
	PostalCode string // ..ADDR.POST
	Country    string // ..ADDR.CTRY
}

type AddressRecords []*AddressRecord

type BibliographyRecord struct {
	Level int // ..BIBL level
}

type BusinessRecord struct {
	Level        int            // ..HEAD.SOUR.CORP level
	BusinessName string         // ..HEAD.SOUR.CORP value
	Address      *AddressRecord // ..HEAD.SOUR.CORP.ADDR
	Phone        []string       // ..HEAD.SOUR.CORP.PHON
}

type ChangeRecord struct {
	Level int         // ..CHAN level
	Date  *DateRecord // ..CHAN.DATE
	Note  NoteRecords // ..CHAN.NOTE
}

type CharacterSetRecord struct {
	Level        int    // ..CHAR level
	CharacterSet string // ..CHAR value
	Version      string // ..CHAR.VERS
}

type ChildStatusRecord struct {
	Level int    // CSTA level
	Xref  string // xref_id of 0 level CSTA
	Name  string // CSTA.NAME
}

type ChildStatusRecords []*ChildStatusRecord

// CitationRecord is a SourceLink

type CitationRecord struct {
	Level   int           // ..SOUR level
	Source  *SourceRecord // value of ..SOUR xref
	Page    string        // ..SOUR.PAGE
	Data    DataRecords   // ..SOUR.DATA
	Text    []string      // ..SOUR.TEXT
	Quality string        // ..SOUR.QUAY
	Media   MediaRecords  // ..SOUR.OBJE
	Note    NoteRecords   // ..SOUR.NOTE
}

type CitationRecords []*CitationRecord

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

type DataRecords []*DataRecord

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

type EventDefinitionRecord struct {
	Level int    // Level 0 _EVENT_DEFN
	Xref  string // Level 0 xref
}

type EventDefinitionRecords []*EventDefinitionRecord

type EventRecord struct {
	Level      int             // ..EVEN level
	Xref       string          // xref_id of level 0 ..EVEN
	Tag        string          // Event tag EVEN or BIRT or ...
	Value      string          // Event value
	Type       string          // ..EVEN.TYPE
	Name       string          // ..EVEN.NAME
	Date       *DateRecord     // ..EVEN.DATE
	Place      *PlaceRecord    // ..EVEN.PLAC
	Role       RoleRecords     // ..EVEN.ROLE
	Address    *AddressRecord  // ..EVEN.ADDR
	Parents    FamilyLinks     // ..EVEN.FAMC
	Husband    *IndividualLink // ..EVEN.HUSB
	Wife       *IndividualLink // ..EVEN.WIFE
	Age        string          // ..EVEN.AGE
	Agency     string          // ..EVEN.AGNC
	Cause      string          // ..EVEN.CAUS
	Temple     string          // ..EVEN.TEMP
	Quality    string          // ..EVEN.QUAY
	UID        []string        // ..EVEN._UID
	RIN        string          // ..EVEN.RIN
	Email      string          // .. EVEN.EMAIL
	Media      MediaRecords    // ..EVEN.OBJE
	Citation   CitationRecords // ..EVEN.SOUR
	Note       NoteRecords     // ..EVEN.NOTE
	Change     *ChangeRecord   // ..EVEN.CHAN
	UpdateTime string          // ..EVEN._UPD
}

type EventRecords []*EventRecord

type FamilyLink struct {
	Level    int           //  level
	Tag      string        // tag from INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Xref     string        // value from INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Family   *FamilyRecord // target of INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Pedigree string        // INDI.FAMC.PEDI or ...
	Adopted  string        // INDI.FAMC.ADOP or ...
	Primary  string        // INDI.FAMC.ADOP or ...
	Note     NoteRecords   // INDI.FAMC.NOTE or ..
}

type FamilyLinks []*FamilyLink

type FamilyRecord struct {
	Level       int             // FAM level
	Xref        string          // xref_id of FAM
	Husband     *IndividualLink // FAM.HUSB
	Wife        *IndividualLink // FAM.WIFE
	NumChildren int             // FAM.NCHI
	Child       IndividualLinks // FAM.CHIL
	Event       EventRecords    // FAM.MARR, FAM.EVEN
	UID         []string        // FAM._UID
	RIN         string          // FAM.RIN
	Media       MediaLinks      // FAM.OBJE
	Citation    CitationRecords // FAM.SOUR
	Note        NoteRecords     // FAM.NOTE
	Change      *ChangeRecord   // FAM.CHAN
	UpdateTime  string          // FAM._UPD
}

type FamilyRecords []*FamilyRecord

type FootnoteRecord struct {
	Level int // FOOT level
}

type GedcomRecord struct {
	Level   int    // HEAD.GEDC level
	Version string // HEAD.GEDC.VERS
	Form    string // HEAD.GEDC.FORM
}

type HeaderRecord struct {
	Level        int                 // HEAD level
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
	Note         NoteRecords         // HEAD.NOTE
	Submitter    *SubmitterLink      // HEAD.SUBM
	Submission   *SubmissionLink     // HEAD.SUBN
	Schema       *SchemaRecord       // HEAD.SCHEMA
}

type HistoryRecord struct {
	Level    int             // ..HIST level
	History  string          // ..HIST value
	Citation CitationRecords // ..HIST.SOUR
}

type HistoryRecords []*HistoryRecord

type IndividualLink struct {
	Level      int               // ..INDI level
	Tag        string            // tag from FAM.HUSB or FAM.WIFE or FAM.CHILD
	Xref       string            // value from FAM.HUSB or FAM.WIFE or FAM.CHILD
	Individual *IndividualRecord // target of FAM.HUSB or FAM.WIFE or FAM.CHILD
	Note       NoteRecords       // FAM.HUSB.NOTE or FAM.WIFE.NOTE or FAM.CHILD.NOTE
}

type IndividualLinks []*IndividualLink

type IndividualRecord struct {
	Level           int             // INDI level
	Xref            string          // xref_id of INDI
	Name            NameRecords     // INDI.NAME
	Sex             string          // INDI.SEX
	Event           EventRecords    // INDI.BIRT, INDI.CHR, INDI.DEAT, INDI.BURI, INDI.EVEN
	Attribute       string          // INDI.ATTR
	Parents         FamilyLinks     // INDI.FAMC
	Family          FamilyLinks     // INDI.FAMS
	Address         AddressRecords  // INDI.ADDR
	Media           MediaLinks      // INDI.OBJE
	Health          string          // INDI.HEAL
	ReferenceNumber string          // INDI.RFN
	Quality         string          // INDI.QUAY
	Living          string          // INDI.LVG
	AFN             []string        // INDI.AFN
	REFN            []string        // INDI.REFN
	UID             []string        // INDI._UID
	RIN             string          // INDI.RIN
	Email           string          // INDI.EMAIL
	WebSite         string          // INDI.WWW
	History         HistoryRecords  // INDI.HIST
	Citation        CitationRecords // INDI.SOUR
	Note            NoteRecords     // INDI.NOTE
	Change          *ChangeRecord   // INDI.CHAN
	UpdateTime      string          // INDI.UPD
}

type IndividualRecords []*IndividualRecord

type MediaLink struct {
	Level    int          // ..OBJE level
	Xref     string       // value of ..OBJE
	Media    *MediaRecord // target of ..OBJE
	Format   string       // ..OBJE.FORM
	Title    string       // ..OBJE.TITL
	FileName string       // ..OBJE.FILE
	Date     string       // ..OBJE.DATE
	Author   string       // ..OBJE.AUTH
	Text     []string     // ..OBJE.TEXT
	Note     NoteRecords  // ..OBJE.NOTE
}

type MediaLinks []*MediaLink

type MediaRecord struct {
	Level    int         // ..OBJE level
	Xref     string      // xref_id of 0 level OBJE
	Format   string      // OBJE.FORM
	Title    string      // OBJE.TITL
	FileName string      // OBJE.FILE
	Note     NoteRecords // OBJE.NOTE
}

type MediaRecords []*MediaRecord

type NameRecord struct {
	Level            int             // ..NAME level
	Name             string          // ..NAME value
	Prefix           string          // ..NAME.NPFX
	GivenName        string          // ..NAME.GIVN
	MiddleName       string          // ..NAME._MIDN
	Surname          string          // ..NAME.SURN
	Suffix           string          // ..NAME.NSFX
	PreferedGiveName string          // ..NAME._PGVN
	AKA              []string        // ..NAME._AKA
	Nickname         []string        // ..NAME.NICK
	Citation         CitationRecords // ..NAME.SOUR
	Note             NoteRecords     // ..NAME.NOTE
}

type NameRecords []*NameRecord

type NoteRecord struct {
	Level    int             // ..NOTE level
	Xref     string          // .. xref_value of 0 level NOTE
	Note     string          // ..NOTE value
	Citation CitationRecords // ..NOTE.SOUR
}

type NoteRecords []*NoteRecord

type PlaceLink struct {
	Level int          // ..PLAC
	Xref  string       // xref to zero level PLAC
	Place *PlaceRecord // link to zero level PLAC
}

type PlacePartRecord struct {
	Level        int    // ..PLAC.PLAn level, n=0..4
	Tag          string // ..PLAC.PLAn tag
	Part         string // ..PLAC.PLAn value
	Jurisdiction string // ..PLAC.PLAn.JURI
}

type PlacePartRecords []*PlacePartRecord

type PlaceRecord struct {
	Level    int              // ..PLAC level
	Xref     string           // xref_id of 0 level PLAC
	Name     string           // ..PLAC value
	Short    string           // ..PLAC.PLAS
	Modifier string           // ..PLAC.PLAM
	Parts    PlacePartRecords // ..PLAC.PLAn n=0..4
	Citation CitationRecords  // ..PLAC.SOUR
	Note     NoteRecords      // ..PLAC.NOTE
	Change   *ChangeRecord    // ..PLAC.CHAN
}

type PlaceRecords []*PlaceRecord

type RepositoryLink struct {
	Level      int               // ..REPO level
	Xref       string            // xref_id of 0 level REPO
	Repository *RepositoryRecord // The linked repository
}

type RepositoryLinks []*RepositoryLink

type RepositoryRecord struct {
	Level   int            // REPO level
	Xref    string         // xref_id of 0 level REPO
	Name    string         // REPO.NAME
	Address *AddressRecord // REPO.ADDR
}

type RepositoryRecords []*RepositoryRecord

type RoleRecord struct {
	Level      int             // ..ROLE level
	Role       string          // ..ROLE no-ref value
	Individual *IndividualLink // ..ROLE ref value
	Principal  string          // ..ROLE.PRIN
}

type RoleRecords []*RoleRecord

type RootRecord struct {
	Level           int                    // root level
	Header          *HeaderRecord          // HEAD
	Submitter       *SubmitterRecord       // SUBM
	Submission      *SubmissionRecord      // SUBN
	Place           PlaceRecords           // PLAC
	Event           EventRecords           // EVEN
	Individual      IndividualRecords      // INDI
	Family          FamilyRecords          // FAM
	Repository      RepositoryRecords      // REPO
	Source          SourceRecords          // SOUR
	Media           MediaRecords           // OBJE
	Note            NoteRecords            // NOTE
	EventDefinition EventDefinitionRecords // _EVENT_DEFN
	ChildStatus     ChildStatusRecords     // CSTA
	Trailer         *TrailerRecord         // TRLR
}

type SchemaRecord struct {
	Level int // ..SCHEMA level
}

// Note: a CitationRecord is a SourceLink

type SourceRecord struct {
	Level         int                 // ..SOUR level
	Xref          string              // xref_id of 0 level SOUR
	Name          string              // ..SOUR.NAME
	Title         string              // ..SOUR.TITL
	Author        string              // ..SOUR.AUTH
	Abbreviation  string              // ..SOUR.ABBR
	Publication   string              // ..SOUR.PUBL
	Parenthesized string              // ..SOUR._PAREN (PAF5)
	Text          []string            // ..SOUR.TEXT
	Data          *DataRecord         // ..SOUR.DATA
	Footnote      *FootnoteRecord     // ..SOUR.FOOT
	Bibliography  *BibliographyRecord // ..SOUR.BIBL
	Repository    *RepositoryLink     // ..SOUR.REPO
	ShortAuthor   string              // ..SOUR.SHAU
	ShortTitle    string              // ..SOUR.SHTI
	Media         MediaLinks          // ..SOUR.OBJE
	Note          NoteRecords         // ..SOUR.NOTE
	Change        *ChangeRecord       // ..SOUR.CHAN
}

type SourceRecords []*SourceRecord

type SubmissionLink struct {
	Level      int               // ..SUBN level
	Xref       string            // xref_id for SUBN
	Submission *SubmissionRecord // SUBN.SUBM
}

type SubmissionRecord struct {
	Level          int              // ..SUBN level
	Xref           string           // xref_id of SUBN
	Submitter      *SubmitterRecord // SUBN.SUBM
	FamilyFileName string           // SUBN.FAMF
	Temple         string           // SUBN.TEMP
	Ancestors      string           // SUBN.ANCE
	Descendents    string           // SUBN.DESC
	Ordinance      string           // SUBN.ORDI
}

type SubmitterLink struct {
	Level     int    // ..SUBM level
	Xref      string // xref_id from ..SUBM
	Submitter *SubmitterRecord
}

type SubmitterRecord struct {
	Level    int            // ..SUBM level
	Xref     string         // xref_id of SUBM
	Name     string         // SUBM.NAME
	Address  *AddressRecord // SUBM.ADDR
	Country  string         // SUBM.CTRY
	Phone    []*string      // SUBM.PHON
	Email    string         // SUBM.EMAIL
	WebSite  string         // SUBM.WWW
	Language string         // SUBM.LANG
	RIN      string         // SUBM.RIN
	Change   *ChangeRecord  // SUBM.CHAN
}

type SystemRecord struct {
	Level       int             // ..HEAD.SOUR level
	Id          string          // HEAD.SOUR value
	Version     string          // HEAD.SOUR.VERS
	ProductName string          // HEAD.SOUR.NAME
	Business    *BusinessRecord // HEAD.SOUR.CORP
	SourceData  *DataRecord     // HEAD.SOUR.DATA
}

type TrailerRecord struct {
	Level int    // ..TRLR level
	Xref  string // Fake id: set to TRLR
}

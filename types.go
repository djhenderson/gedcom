/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

// AddressRecord represents an address record
type AddressRecord struct {
	Level      int         // ..ADDR level
	Full       string      // ..ADDR value
	Line1      string      // ..ADDR.ADR1
	Line2      string      // ..ADDR.ADR2
	Line3      string      // ..ADDR.ADR3
	City       string      // ..ADDR.CITY
	State      string      // ..ADDR.STAE
	PostalCode string      // ..ADDR.POST
	Country    string      // ..ADDR.CTRY
	Phone      string      // ..ADDR.PHON
	Name_      string      // ..ADDR._NAME (AQ14)
	Note       NoteRecords // .. ADDR.NOTE (Leg8)
}

// AddressRecords represents a slice of address records
type AddressRecords []*AddressRecord

// AuthorRecord represents a data record
type AuthorRecord struct {
	Level        int    // ..AUTH level
	Author       string // value of ..AUTH
	Abbreviation string // ..AUTH.ABBR
}

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
	Level             int           // ..SOUR level; not 0
	Value             string        // value of ..SOUR excluding xref
	Source            *SourceRecord // linked source
	Page              string        // ..SOUR.PAGE
	Reference         string        // ..SOUR.REF
	FamilySearchFTID_ string        // ..SOUR._FSFTID (AQ14)
	Event             EventRecords  // ..SOUR.EVEN
	Data              DataRecords   // ..SOUR.DATA
	Text              []string      // ..SOUR.TEXT
	Quality           string        // ..SOUR.QUAY
	Media             MediaLinks    // ..SOUR.OBJE
	CONS              string        // ..SOUR.CONS
	Direct            string        // ..SOUR.DIRE
	SourceQuality     string        // ..SOUR.SOQU
	Note              NoteRecords   // ..SOUR.NOTE
	Date              string        // ..SOUR.DATE (Leg8)
	ReferenceNumber   string        // ..SOUR.REFN
	Rin_              string        // ..SOUR._RIN (AQ14)
	AppliesTo_        string        // ..SOUR._APPLIES_TO (AQ15)
}

// CitationRecords represents a slice of links to source records
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
	Tag   string   // ..DATE tag
	Date  string   // ..DATE value
	Time  string   // ..DATE.TIME value
	Text  []string // ..DATE.TEXT
	Day   string   // ..DATA.DATD
	Month string   // ..DATA.DATM
	Year  string   // ..DATA.DATY
	Full  string   // ..DATA.DATF
	Short string   // ..DATA.DATS
}

// EventDefinitionRecord represents a GEDCOM event definition record.
type EventDefinitionRecord struct {
	Level int    // Level 0 _EVENT_DEFN
	Xref  string // Level 0 xref
}

// EventDefinitionRecords represents a slice of event definition records.
type EventDefinitionRecords []*EventDefinitionRecord

// EventRecord represents a GEDCOM event record.
type EventRecord struct {
	Level           int             // ..EVEN level; 0 or higher
	Xref            string          // xref_id of level 0 ..EVEN
	Tag             string          // Event tag EVEN or BIRT or ...
	Value           string          // Event value
	Type            string          // ..EVEN.TYPE
	Name            string          // ..EVEN.NAME
	Primary_        string          // ..EVEN._PRIM
	Date            *DateRecord     // ..EVEN.DATE
	Date2_          *DateRecord     // ..EVEN._DATE2 (AQ14)
	Place           *PlaceRecord    // ..EVEN.PLAC
	Place2_         *PlaceRecord    // ..EVEN._PLAC2 (AQ14)
	Description2_   string          // ..EVEN._Description2 (AQ14)
	Role            RoleRecords     // ..EVEN.ROLE
	Address         *AddressRecord  // ..EVEN.ADDR
	Phone           []string        // ..EVEN.PHON
	Parents         FamilyLinks     // ..EVEN.FAMC
	Husband         *IndividualLink // ..EVEN.HUSB
	Wife            *IndividualLink // ..EVEN.WIFE
	Spouse          *IndividualLink // ..EVEN.SPOU
	Agency          string          // ..EVEN.AGNC
	Cause           string          // ..EVEN.CAUS
	Temple          string          // ..EVEN.TEMP
	Quality         string          // ..EVEN.QUAY
	Status          string          // ..EVEN.STAT
	UniqueId_       []string        // ..EVEN._UID
	RecordInternal  string          // ..EVEN.RecordInternal
	Email           string          // ..EVEN.EMAIL
	Media           MediaLinks      // ..EVEN.OBJE
	Citation        CitationRecords // ..EVEN.SOUR
	Note            NoteRecords     // ..EVEN.NOTE
	Change          *ChangeRecord   // ..EVEN.CHAN
	UpdateTime_     string          // ..EVEN._UPD
	AlternateBirth_ string          // ..EVEN._ALT_BIRTH (AQ14)
	Confidential_   string          // ..EVEN._CONFIDENTIAL (AQ14)
}

// EventRecords represents a slice of event records.
type EventRecords []*EventRecord

// FamilyLink represents a GEDCOM link to a family record.
type FamilyLink struct {
	Level    int             //  level
	Tag      string          // tag from INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Family   *FamilyRecord   // target of INDI.FAMC or INDI.FAMS or EVEN.FAMC
	Adopted  string          // INDI.FAMC.ADOP or ...
	Primary_ string          // INDI.FAMC._PRIMARY or ...
	Note     NoteRecords     // INDI.FAMC.NOTE or ..
	Pedigree *PedigreeRecord // INDI.FAMC.PEDI or ..
	Citation CitationRecords // INDI.FAMC.SOUR or ..
}

// FamilyLinks represents a slice of links to family records.
type FamilyLinks []*FamilyLink

// FamilyRecord represents a GEDCOM family record.
type FamilyRecord struct {
	Level               int                        // FAM level; only 0
	Xref                string                     // xref_id of FAM
	Status_             string                     // FAM._STAT (AQ14)
	NoChildren_         string                     // FAM._NONE (AQ14)
	Husband             *IndividualLink            // FAM.HUSB
	Wife                *IndividualLink            // FAM.WIFE
	NumChildren         int                        // FAM.NCHI
	Child               IndividualLinks            // FAM.CHIL
	Event               EventRecords               // FAM.MARR, FAM.EVEN
	UniqueId_           []string                   // FAM._UID
	RecordInternal      string                     // FAM.RecordInternal
	UserReferenceNumber UserReferenceNumberRecords // FAM.REFN
	Media               MediaLinks                 // FAM.OBJE
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

// HeaderRecord represents a GEDCOM header record
// There can be only one!
type HeaderRecord struct {
	Level        int                 // HEAD level; always 0
	Xref         string              // Fake id; set to HEAD
	SourceSystem *SystemRecord       // HEAD.SOUR
	Destination  string              // HEAD.DEST
	Date         *DateRecord         // HEAD.DATE
	FileName     string              // HEAD.FILE
	Gedcom       *GedcomRecord       // HEAD.GEDC
	CharacterSet *CharacterSetRecord // HEAD.CHAR
	Language     string              // HEAD.LANG
	Copyright    string              // HEAD.COPR
	Place        *PlaceRecord        // HEAD.PLAC
	RootPerson_  *IndividualLink     // HEAD._ROOT - FmP root person
	HomePerson_  *IndividualLink     // HEAD._HME - home person
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
	Age          string            // ..EVEN.HUSB.AGE or ..EVEN.WIFE.AGE or ..EVEN.SPOU.AGE
	Preferred_   string            // FAM.HUSB._PREF or FAM.WIFE._PREF or ... (Leg8)
}

// IndividualLinks represents a slice of links to individual records
type IndividualLinks []*IndividualLink

// IndividualRecord represents an individual record
type IndividualRecord struct {
	Level               int                        // INDI level; always 0
	Xref                string                     // xref_id of INDI
	Name                NameRecords                // INDI.NAME
	Status_             string                     // INDI._STAT (AQ14)
	Restriction         string                     // INDI.RESN
	Sex                 string                     // INDI.SEX
	Event               EventRecords               // INDI.BIRT, INDI.CHR, INDI.DEAT, INDI.BURI, INDI.EVEN
	Attribute           string                     // INDI.ATTR
	Parents             FamilyLinks                // INDI.FAMC
	Family              FamilyLinks                // INDI.FAMS
	Address             AddressRecords             // INDI.ADDR
	Phone               []string                   // INDI.PHON
	Media               MediaLinks                 // INDI.OBJE
	Health              string                     // INDI.HEAL
	History             HistoryRecords             // INDI.HIST
	Quality             string                     // INDI.QUAY
	Living              string                     // INDI.LVG
	CONL                string                     // INDI.CONL
	AncestralFileNumber []string                   // INDI.AFN
	RecordFileNumber    string                     // INDI.RFN
	UserReferenceNumber UserReferenceNumberRecords // INDI.REFN
	FamilySearchFTID_   string                     // INDI._FSFTID (AQ14)
	FamilySearchLink_   string                     // INDI._FSLINK (Leg8)
	RecordInternal      string                     // INDI.RecordInternal
	UniqueId_           []string                   // INDI._UID
	Email               string                     // INDI.EMAIL
	Email_              string                     // INDI._EMAIL (AQ14)
	URL_                string                     // INDI._URL (AQ14)
	WebSite             string                     // INDI.WWW
	Citation            CitationRecords            // INDI.SOUR
	Note                NoteRecords                // INDI.NOTE
	Associated          IndividualLinks            // INDI.ASSO
	Submitter           SubmitterLinks             // INDI.SUBM
	ANCI                SubmitterLinks             // INDI.ANCI
	DESI                SubmitterLinks             // INDI.DESI
	UpdateTime_         string                     // INDI._UPD
	Alias               string                     // INDI.ALIA
	Father              *IndividualLink            // INDI.FATH
	Mother              *IndividualLink            // INDI.MOTH
	Miscellaneous       []string                   // INDI.MISC
	ProfilePicture_     *MediaLink                 // INDI._PROF
	PPExclude_          string                     // INDI._PPEXCLUDE (Leg8)
	Change              *ChangeRecord              // INDI.CHAN
	Todo_               string                     // INDI._TODO (AQ15)
}

// IndividualRecords represents a slice of individual records
type IndividualRecords []*IndividualRecord

// MediaLink represents a link to an media record
type MediaLink struct {
	Level int          // ..OBJE level
	Tag   string       // tag from OBJE or _PROF
	Value string       // value from OBJE or _PROF
	Media *MediaRecord // target of OBJE or _PROF
}

// MediaLinks represents a slice of links to media records
type MediaLinks []*MediaLink

// MediaRecord represents a GEDCOM media record
type MediaRecord struct {
	Level               int                        // ..OBJE level
	Xref                string                     // xref_id of 0 level or xref from value
	Tag                 string                     // ..TAG tag either OBJE or _PROF
	Value               string                     // value without xref
	Format              string                     // OBJE.FORM
	Url_                string                     // OBJE._URL
	FileName            string                     // OBJE.FILE
	Title               string                     // OBJE.TITL
	Date                string                     // OBJE.DATE
	Author              string                     // OBJE.AUTH
	Text                string                     // OBJE.TEXT
	Note                NoteRecords                // OBJE.NOTE
	Date_               string                     // OBJE._DATE
	AstId_              string                     // OBJE._ASTID - FmP identifier
	AstType_            string                     // OBJE._ASTTYP - FmP type
	AstDesc_            string                     // OBJE._ASTDESC - FmP description
	AstLoc_             string                     // OBJE._ASTLOC - FmP location
	AstPerm_            string                     // OBJE._ASTPERM - FmP permissions
	AstUpPid_           string                     // OBJE._ASTUPPID - FmP update identifier?
	BinaryLargeObject   *BlobRecord                // OBJE.BLOB
	UserReferenceNumber UserReferenceNumberRecords // OBJE.REFN
	Rin                 string                     // OBJE.RecordInternal
	Change              *ChangeRecord              // OBJE.CHAN
	Scbk_               string                     // OBJE._SCBK (AQ14)
	Primary_            string                     // OBJE._PRIM (AQ14)
	Type_               string                     // OBJE._TYPE (AQ14)
	Sshow_              string                     // OBJE._SSHOW (AQ14)
	Stime_              string                     // OBJE._STIME (AQ15)
	mediaLinks          MediaLinks                 // OBJE.OBJE (AQ15)
	SrcPp_              string                     // OBJE._SRCPP (AQ15)
	SrcFlip_            string                     // OBJE._SRCFLIP (AQ15)
	FsFtId_             string                     // OBJE._FSFTID (AQ15)
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
	RomanizedName      string          // ..NAME.ROMN
	PhoneticName       string          // ..NAME.FONE
	MarriedName_       string          // ..NAME._MARNM (AQ14)
	Primary_           string          // ..NAME._PRIM - FmP primary/preferred
	NameType           string          // ..NAME.TYPE
	AlsoKnownAs_       []string        // ..NAME._AKA
	Nickname           []string        // ..NAME.NICK
	Citation           CitationRecords // ..NAME.SOUR
	Note               NoteRecords     // ..NAME.NOTE
}

// NameRecords represents a slice of name records
type NameRecords []*NameRecord

// NoteRecord represents a GEDCOM note record
type NoteRecord struct {
	Level               int                        // ..NOTE level
	Xref                string                     // .. xref_value of 0 level NOTE
	Note                string                     // ..NOTE value
	Citation            CitationRecords            // ..NOTE.SOUR
	UserReferenceNumber UserReferenceNumberRecords // ..NOTE.REFN
	RecordInternal      string                     // ..NOTE.RecordInternal
	Change              *ChangeRecord              // ..SOUR.CHAN
}

// NoteRecords represents a slice of note records
type NoteRecords []*NoteRecord

// PedigreeRecord represents a GEDCOM pedigree record.
type PedigreeRecord struct {
	Level    int    // ..FAMC.PEDI level
	Pedigree string // ..FAMC.PEDI value
	Husband_ string // ..FAMC.PEDI._HUSB value
	Wife_    string // ..FAMC.PEDI._WIFE value
}

// PlaceDefinitionRecord represents a GEDCOM place definition record.
type PlaceDefinitionRecord struct {
	Level        int    // Level 0 _PLAC_DEFN
	Xref         string // Level 0 xref
	Place        string // _PLAC_DEFN.PLAC
	Abbreviation string // _PLAC_DEFN.ABBR
}

// PlaceDefinitionRecords represents a slice of place definition records.
type PlaceDefinitionRecords []*PlaceDefinitionRecord

// PlacePartRecord represents a place part record
type PlacePartRecord struct {
	Level        int    // ..PLAC.PLAn level, n=0..4
	Tag          string // ..PLAC.PLAn tag
	Part         string // ..PLAC.PLAn value
	Jurisdiction string // ..PLAC.PLAn.JURI
}

// PlacePartRecords represents a slice of place part records
type PlacePartRecords []*PlacePartRecord

// PlaceRecord represents a GEDCOM place record
type PlaceRecord struct {
	Level     int              // ..PLAC level; 0 or higher
	Xref      string           // xref_id of 0 level PLAC
	Tag       string           // ..PLAC tag
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

// RepositoryRecord represents a GEDCOM repository record
type RepositoryRecord struct {
	Level               int                        // REPO level; always 0
	Xref                string                     // xref_id of 0 level REPO
	Name                string                     // REPO.NAME
	Address             *AddressRecord             // REPO.ADDR
	Phone               []string                   // REPO.PHON
	WebSite             string                     // REPO.WWW
	UserReferenceNumber UserReferenceNumberRecords // REPO.REFN
	RecordInternal      string                     // REPO.RecordInternal
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

// RootRecord represents the root record of a GEDCOM file.
type RootRecord struct {
	Level            int                    // root level, always zero
	Header           *HeaderRecord          // HEAD
	Submitter        SubmitterRecords       // SUBM
	Submission       SubmissionRecords      // SUBN
	Place            PlaceRecords           // PLAC
	Event            EventRecords           // EVEN
	Individual       IndividualRecords      // INDI
	Family           FamilyRecords          // FAM
	Media            MediaRecords           // OBJE
	Note             NoteRecords            // NOTE
	PlaceDefinition_ PlaceDefinitionRecords // _PLAC_DEFN (Leg8)
	EventDefinition_ EventDefinitionRecords // _EVENT_DEFN (AQ14)
	ChildStatus      ChildStatusRecords     // CSTA
	Todo_            TodoRecords            // _TODO (AQ15)
	Source           SourceRecords          // SOUR
	Repository       RepositoryRecords      // REPO
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

// SourceRecord represents a GEDCOM source record.
// Note: A CitationRecord is a link to a source record
type SourceRecord struct {
	Level               int                        // ..SOUR level; always 0
	Xref                string                     // xref_id of 0 level SOUR
	Value               string                     // ..SOUR value
	Name                string                     // ..SOUR.NAME
	Title               string                     // ..SOUR.TITL
	Author              *AuthorRecord              // ..SOUR.AUTH
	Abbreviation        string                     // ..SOUR.ABBR
	Publication         string                     // ..SOUR.PUBL
	MediaType           string                     // ..SOUR.MEDI (Leg8)
	Parenthesized_      string                     // ..SOUR._PAREN (PAF5)
	Text                []string                   // ..SOUR.TEXT
	Data                *DataRecord                // ..SOUR.DATA
	Footnote            *FootnoteRecord            // ..SOUR.FOOT
	Bibliography        *BibliographyRecord        // ..SOUR.BIBL
	Repository          *RepositoryLink            // ..SOUR.REPO
	UserReferenceNumber UserReferenceNumberRecords // ..SOUR.REFN
	Quality             string                     // ..SOUR.QUAY
	RecordInternal      string                     // ..SOUR.RecordInternal
	ShortAuthor         string                     // ..SOUR.SHAU
	ShortTitle          *ShortTitleRecord          // ..SOUR.SHTI
	Media               MediaLinks                 // ..SOUR.OBJE
	Note                NoteRecords                // ..SOUR.NOTE
	Change              *ChangeRecord              // ..SOUR.CHAN
	Type_               string                     // ..SOUR._TYPE (AQ14)
	Other_              string                     // ..SOUR._OTHER (AQ14)
	Master_             string                     // ..SOUR._MASTER (AQ14)
	Italic_             string                     // ..SOUR._ITALIC (AQ14)
	WebTag_             *WebTagRecord              // ..SOUR._WEBTAG (Leg8)
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

// SubmissionRecord represents a GEDCOM submission record.
type SubmissionRecord struct {
	Level          int              // SUBN level; always 0
	Xref           string           // xref_id of SUBN
	Submitter      *SubmitterRecord // SUBN.SUBM
	FamilyFileName string           // SUBN.FAMF
	Temple         string           // SUBN.TEMP
	Ancestors      string           // SUBN.ANCE
	Descendents    string           // SUBN.DESC
	Ordinance      string           // SUBN.ORDI
	RecordInternal string           // SUBN.RecordInternal
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
	Email_           string         // SUBM._EMAIL (AQ14)
	WebSite          string         // SUBM.WWW
	Language         string         // SUBM.LANG
	Media            MediaLinks     // SUBM.OBJE
	RecordFileNumber string         // SUBM.RFN
	STAL             string         // SUBM.STAL
	NUMB             string         // SUBM.NUMB
	RecordInternal   string         // SUBM.RecordInternal
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

// TodoLink represents a link to a todo record
type TodoLink struct {
	Level int    // .._TODO level
	Xref  string // .._TODO value
}

// TodoLinks represents a slice of todo links
type TodoLinks []*TodoLink

// TodoRecord represents a todo record (AQ15)
type TodoRecord struct {
	Level       int    // _TODO level (equal 0)
	Xref        string // xref_id of 0 level _TODO
	Value       string // .._TODO value
	Description string // .._TODO.DESC
	Priority_   string // .._TODO._PRIORITY
	Type        string // .._TODO.TYPE
	Status      string // .._TODO.STAT
	Date2_      string // .._TODO._DATE2
}

// TodoRecords represents a slice of todo records
type TodoRecords []*TodoRecord

// TrailerRecord represents a GEDCOM trailer record
// There can be only one!
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

// WebTagRecord represents a web tag record
type WebTagRecord struct { // RM6
	Level int    // .._WEBTAG level
	Xref  string // xref_id of _WEBTAG
	Value string // .._WEBTAG value
	Name  string // .._WEBTAG.NAME
	URL   string // .._WEBTAG.URL
}

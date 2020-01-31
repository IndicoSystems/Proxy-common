package metadata

import (
	"encoding/base64"
	"encoding/json"
	"github.com/btubbs/datetime"
	"github.com/indicosystems/proxy/logger"
	"github.com/sirupsen/logrus"
	"strings"

	"time"
)

const (
	// The ID of the user the file belongs to.
	UserId = "userid"

	// The name of the container the file belongs to.
	ParentName = "parentname"

	// The ID of the container the file belongs to.
	ParentId = "parentid"

	// The description of the parent
	ParentDescription = "parentdescription"

	// The ISO8601 compliant timestamp at which the file was created.
	CreatedAt = "createdat"

	// The mime type of the file.
	FileType = "filetype"

	// The name given to the file by the user.
	DisplayName = "displayname"

	// Description given to the file
	Description = "description"

	// The Checksum of the file.
	Checksum = "checksum"

	// The type of Checksum used.
	ChecksumType = "checksumtype" // – md5, sha256, crc etc.

	// The name of the file on the file system.
	Filename = "filename"

	// The ID of the file at the external party
	ExtId = "extid"

	// The case number that an officer would use in their system
	CaseNumber = "casenumber"

	// The date of which the media was created. Might be different than the CreatedAt-field, which often represents the creation-date in a database
	CapturedAt = "capturedat"

	// Duration for the media, in case of audio/video. Always in ms.
	Duration = "duration"

	// The surname of the creator, often the officer.
	CreatorSurname = "creatorsurname"

	// The district of which the creator belongs.
	CreatorDistrict = "creatordistrict"

	// Textual representation of the location, like Madrid, or street address
	LocationText = "locationtext"

	// Tags for the item, comma-separated
	Tags = "tags"

	// Geo-Latitude.
	Latitude = "latitude"

	// Geo-Longitude
	Longitude = "longitude"

	Subjects = "subjects"

	AccountName   = "accountname"   // – IR: windows-kontoen
	EquipmentID   = "equipmentid"   // – IR: InforPrefQuipmentID / DeviceInfo ?? også feltetd CustomDriverxbDeviceType
	InterviewType = "interviewtype" // Itervju, avhør, etc.
	Bookmarks     = "bookmarks"
	Attachments   = "attachments" // id'er til ClientMediaID? hmm...
	// A group-id can be specified by the client, if these items should be grouped.
	// All uploads that share this key will be grouped. Not all backends supports this.
	// This groupId is only used to link the files together, it is not stored on the backend.
	GroupID = "groupid"
	// If a group should be created/found on backend, this can be used to search for it.
	// Note that GroupID is required anyway.
	GroupName = "groupname"

	// Any additional notes
	Notes = "notes"

	// A unique identifier on the client. With attachments, this field is required.
	ClientMediaId = "clientmediaid"

	// Etc
	Etcetera = "etc"

	// SSN
	SSN = "ssn"
	// DeferId - Internal use for deferred uploads
	DeferId = "__deferId"
	// DeferredFiledId - Internal use for deferred uploads
	DeferFileId = "__deferFileId"
)

type Upl struct {
	UploadMetadata
}

var l logrus.FieldLogger = logger.Get("metadata")

type Metadata map[string]string
type Mapper func(data Metadata) Metadata

func (m Metadata) GetCreatedAtTime() *time.Time {
	return parseDateSafely(m.GetExact(CreatedAt))
}
func (m Metadata) GetCapturedAtTime() *time.Time {
	return parseDateSafely(m.GetExact(CapturedAt))
}

// Parses a ISO-8601-compliant string into a native time.
func parseDateSafely(v string) *time.Time {
	if v == "" {
		return nil
	}
	t, err := datetime.Parse(v, time.UTC)
	if err != nil {
		return nil
	}

	return &t
}

func (m Metadata) ConvertToType() UploadMetadata {
	checksum := m.GetChecksum()
	bookmarks := m.GetBookmarks()
	u := UploadMetadata{
		ClientMediaId: m.GetExact(ClientMediaId),
		GroupId:       m.GetExact(GroupID),
		GroupName:     m.GetExact(GroupName),
		UserId:        m.GetExact(UserId),
		Parent: &Parent{
			Id:          m.GetExact(ParentId),
			Name:        m.GetExact(ParentName),
			Description: m.GetExact(ParentDescription),
		},
		CreatedAt:     m.GetCreatedAtTime(),
		CapturedAt:    m.GetCapturedAtTime(),
		Duration:      m.GetExact(Duration),
		FileType:      m.GetExact(FileType),
		DisplayName:   m.GetExact(DisplayName),
		Description:   m.GetExact(Description),
		Checksum:      &checksum,
		FileName:      m.GetExact(Filename),
		ExtId:         m.GetExact(ExtId),
		CaseNumber:    m.GetExact(CaseNumber),
		Subject:       m.GetPerson(),
		AccountName:   m.GetExact(AccountName),
		EquipmentId:   m.GetExact(EquipmentID),
		InterviewType: m.GetExact(InterviewType),
		Bookmarks:     &bookmarks,
		Notes:         m.GetExact(Notes),
		Etc:           m.GetEtc(),
		SSN:           m.GetSsn(),
	}

	creator := Creator{
		District: m.GetExact(CreatorDistrict),
		Person:   Person{LastName: m.GetExact(CreatorSurname)},
	}
	if creator != (Creator{}) {
		u.Creator = &creator
	}
	location := m.GetLocation()
	if location != nil {
		u.Location = location
	}
	subject := m.GetPerson()
	if subject != nil {
		u.Subject = subject
	}
	return u
}

// Etc er unmapped metadata. In Indico Gateway, these represent an organizational-form.
type Etc struct {
	// The metadata-key.
	Key string `xml:",omitempty"`
	// The id of the field, where available.
	FieldId string `xml:",omitempty"`
	// The key used to get the translated VisualName, where available.
	TranslationKey string `xml:",omitempty"`
	// The visual name, as reported by the client.
	VisualName string `xml:",omitempty"`
	// The value of the field, e.g. the user-input
	Value string `xml:""`
	// Marks whether or not field is required or not.
	Required bool `xml:",omitempty"`
	// The kind of data in the Value-field.
	DataType string `xml:",omitempty"`
}

func (m Metadata) GetChecksum() MetaChecksum {
	return MetaChecksum{
		Value:        m.GetExact(Checksum),
		ChecksumType: m.GetExact(ChecksumType),
	}
}

func (m Metadata) GetLocation() *Location {

	l := Location{
		Text:      m.GetExact(LocationText),
		Latitude:  m.GetExact(Latitude),
		Longitude: m.GetExact(Longitude),
	}
	if l == (Location{}) {
		return nil
	}
	return &l
}

// Helper for getting nested objects
func (m Metadata) getNested(key string, v interface{}) error {

	str := m.Get(key)
	if str == "" {
		return nil
	}
	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	err = json.Unmarshal(sDec, &v)
	if err != nil {
		return err
	}
	return nil
}

// Bookmark object belonging to a certain recording.
type Bookmark struct {
	CreationDate  string
	ID            string
	Title         string
	StartPosition int
	EndPosition   int
}

func (m Metadata) GetBookmarks() (bm []Bookmark) {
	m.getNested(Bookmarks, &bm)
	return
}
func (m Metadata) GetPerson() *[]Person {
	var ps []Person
	m.getNested(Subjects, &ps)
	for _, p := range ps {
		p.Country = "bob"
		if (Person{} != p) {
			return &ps
		}

	}
	return nil
}
func (m Metadata) GetEtc() (etc *[]Etc) {
	m.getNested(Etcetera, &etc)
	return
}

// TODO: Map to correct ssn-structure
func (m Metadata) GetSsn() (ssn *map[string]interface{}) {
	m.getNested(SSN, &ssn)
	return
}
func (m Metadata) GetExact(key string) string {
	if val, ok := m[key]; ok {
		return strings.TrimSpace(val)
	}
	if val, ok := m[strings.ToLower(key)]; ok {
		return strings.TrimSpace(val)
	}
	return ""
}

func (m Metadata) GetFileName() string {
	return m.GetExact(Filename)
}

// Can be used to get key/values with fallbacks. The keys to the left have precedence.
func (m Metadata) Fallbacks(keys ...string) string {
	for _, k := range keys {
		if v := m.GetExact(k); v != "" {
			return v
		}
	}
	return ""
}

func (m Metadata) Get(key string) string {
	if val, ok := m[key]; ok {
		return strings.TrimSpace(val)
	}
	lk := strings.ToLower(key)
	if val, ok := m[lk]; ok {
		return strings.TrimSpace(val)
	}
	for k, v := range m {
		if strings.ToLower(k) == lk {
			l.Warnf("Return for-loop, this is a deprecated behaviour, %s, %v", key, m[lk])
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (m Metadata) Set(key, value string) {
	m[strings.ToLower(key)] = value
}

func (m Metadata) Map(from, to string) {
	m.Set(to, m.Get(from))
}

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

	// Tags for the item
	Tags = "tags"

	// Geo-Latitude.
	Latitude = "latitude"

	// Geo-Longitude
	Longitude = "longitude"

	Subjects = "subjects"

	SubjectFirstName   = "subjectfirstname"
	SubjectLastName    = "subjectlastname" // IR kaller dette SessionCustomerName ??
	SubjectId          = "subjectid"       // Personnummer etc.
	SubjectBirthdate   = "subjectbirthdate"
	SubjectGender      = "subjectgender"
	SubjectNationality = "subjectnationality"
	SubjectWorkplace   = "subjectworkplace"
	SubjectStatus      = "subjectstatus" // – Indikerer om subjektet er et vitne, offer, etc.
	SubjectAddress     = "subjectaddress"
	SubjectZip         = "subjectzip"
	SubjectPostalCode  = "subjectpostalcode" // – IR: SessionCustomerIndex
	SubjectCountry     = "subjectcountry"    // – Land for adresse
	SubjectWorkPhone   = "subjectworkphone"
	SubjectPhone       = "subjectphone"
	SubjectMobile      = "subjectmobile"
	SubjectPresent     = "subjectpresent" // – Er subjektet tilstede? boolean
	AccountName        = "accountname"    // – IR: windows-kontoen
	EquipmentID        = "equipmentid"    // – IR: InforPrefQuipmentID / DeviceInfo ?? også feltetd CustomDriverxbDeviceType
	InterviewType      = "interviewtype"  // Itervju, avhør, etc.
	Bookmarks          = "bookmarks"      // JSON
	Attachments        = "attachments"    // id'er til ClientMediaID? hmm...
	// A group-id can be specified by the client, if these items should be grouped.
	// All uploads that share this key will be grouped. Not all backends supports this.
	// This groupId is only used to link the files together, it is not stored on the backend.
	GroupID = "groupid"
	// If a group should be created/found on backend, this can be used to search for it.
	// Note that GroupID is required anyway.
	GroupName = "groupnam"

	// Any additional notes
	Notes = "notes"

	// A unique identifier on the client. With attachments, this field is required.
	ClientMediaId = "clientmediaid"

	// Etc
	Etcetera = "etc"

	// SSN
	SSN = "ssn"
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
	return UploadMetadata{
		ClientMediaId: m.GetExact(ClientMediaId),
		GroupId:       m.GetExact(GroupID),
		GroupName:     m.GetExact(GroupName),
		UserId:        m.GetExact(UserId),
		Parent: Parent{
			Id:          m.GetExact(ParentId),
			Name:        m.GetExact(ParentName),
			Description: m.GetExact(ParentDescription),
		},
		CreatedAt:   m.GetCreatedAtTime(),
		CapturedAt:  m.GetCapturedAtTime(),
		Duration:    m.GetExact(Duration),
		FileType:    m.GetExact(FileType),
		DisplayName: m.GetExact(DisplayName),
		Description: m.GetExact(Description),
		Checksum:    m.GetChecksum(),
		FileName:    m.GetExact(Filename),
		ExtId:       m.GetExact(ExtId),
		CaseNumber:  m.GetExact(CaseNumber),
		Creator: Creator{
			District: m.GetExact(CreatorDistrict),
			Person:   Person{LastName: m.GetExact(CreatorSurname)},
		},
		Location: Location{
			Text:      m.GetExact(LocationText),
			Latitude:  m.GetExact(Latitude),
			Longitude: m.GetExact(Longitude),
		},
		Subject:       m.GetPerson(),
		AccountName:   m.GetExact(AccountName),
		EquipmentId:   m.GetExact(EquipmentID),
		InterviewType: m.GetExact(InterviewType),
		Bookmarks:     m.GetExact(Bookmarks),   // TODO: parse
		Attachments:   m.GetExact(Attachments), // TODO: parse
		Notes:         m.GetExact(Notes),
		Etc:           m.GetEtc(),
		SSN:           m.GetSsn(),
	}
}

type Etc map[string]interface{}

func (m Metadata) GetChecksum() MetaChecksum {
	return MetaChecksum{
		Value:        m.GetExact(Checksum),
		ChecksumType: m.GetExact(ChecksumType),
	}
}

func (m Metadata) GetLocation() Location {
	return Location{
		Text:      m.GetExact(LocationText),
		Latitude:  m.GetExact(Latitude),
		Longitude: m.GetExact(Longitude),
	}
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

func (m Metadata) GetPerson() (ps []Person) {
	m.getNested(Subjects, &ps)
	return
}
func (m Metadata) GetEtc() (etc Etc) {
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

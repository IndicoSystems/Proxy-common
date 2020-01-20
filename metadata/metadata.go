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
)

type Upl struct {
	UploadMetadata
}

var l logrus.FieldLogger = logger.Get("metadata")

type Metadata map[string]string
type Mapper func(data Metadata) Metadata

func (m Metadata) GetCreatedAtTime() *time.Time {
	return parseDateSafely(m.getExact(CreatedAt))
}
func (m Metadata) GetCapturedAtTime() *time.Time {
	return parseDateSafely(m.getExact(CapturedAt))
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
		ClientMediaId: m.getExact(ClientMediaId),
		GroupId:       m.getExact(GroupID),
		GroupName:     m.getExact(GroupName),
		UserId:        m.getExact(UserId),
		Parent: Parent{
			Id:          m.getExact(ParentId),
			Name:        m.getExact(ParentName),
			Description: m.getExact(ParentDescription),
		},
		CreatedAt:   m.GetCreatedAtTime(),
		CapturedAt:  m.GetCapturedAtTime(),
		Duration:    m.getExact(Duration),
		FileType:    m.getExact(FileType),
		DisplayName: m.getExact(DisplayName),
		Description: m.getExact(Description),
		Checksum:    m.GetChecksum(),
		FileName:    m.getExact(Filename),
		ExtId:       m.getExact(ExtId),
		CaseNumber:  m.getExact(CaseNumber),
		Creator: Creator{
			District: m.getExact(CreatorDistrict),
			Person:   Person{LastName: m.getExact(CreatorSurname)},
		},
		Location: Location{
			Text:      m.getExact(LocationText),
			Latitude:  m.getExact(Latitude),
			Longitude: m.getExact(Longitude),
		},
		Subject:       m.GetPerson(),
		AccountName:   m.getExact(AccountName),
		EquipmentId:   m.getExact(EquipmentID),
		InterviewType: m.getExact(InterviewType),
		Bookmarks:     m.getExact(Bookmarks),   // TODO: parse
		Attachments:   m.getExact(Attachments), // TODO: parse
		Notes:         m.getExact(Notes),
	}
}

func (m Metadata) GetChecksum() MetaChecksum {
	return MetaChecksum{
		Value:        m.getExact(Checksum),
		ChecksumType: m.getExact(ChecksumType),
	}
}

func (m Metadata) GetPerson() []Person {
	str := m.Get(Subjects)
	if str == "" {
		return nil
	}
	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil
	}
	var ps []Person
	err = json.Unmarshal([]byte(sDec), &ps)
	if err != nil {
		return nil
	}
	return ps

}
func (m Metadata) getExact(key string) string {
	if val, ok := m[key]; ok {
		return strings.TrimSpace(val)
	}
	if val, ok := m[strings.ToLower(key)]; ok {
		return strings.TrimSpace(val)
	}
	return ""
}

func (m Metadata) GetFileName() string {
	return m.getExact(Filename)
}

// Can be used to get key/values with fallbacks. The keys to the left have precedence.
func (m Metadata) Fallbacks(keys ...string) string {
	for _, k := range keys {
		if v := m.getExact(k); v != "" {
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

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

type Person struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Id          string `json:"id"`
	Dob         string `json:"dob"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Workplace   string `json:"workplace"`
	Status      string `json:"status"`
	Address     string `json:"address"`
	ZipCode     int    `json:"zip"`
	Country     string `json:"country"`
	WorkPhone   string `json:"workPhone"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	Present     bool   `json:"isPresent"`
}

type Parent struct {
	Id          string
	Name        string
	Description string
}

type Creator struct {
	District string
	Person
}

type Type struct {
	UserId        string
	Parent        Parent
	CreatedAt     *time.Time
	CapturedAt    *time.Time
	FileType      string
	DisplayName   string
	Description   string
	Checksum      MetaChecksum
	FileName      string
	ExtId         string
	CaseNumber    string
	Duration      string
	Creator       Creator
	Location      Location
	Subject       []Person
	AccountName   string
	EquipmentId   string
	InterviewType string
	Bookmarks     string
	Attachments   string
	Notes         string
	ClientMediaId string
	GroupId       string
	GroupName     string
}

type MetaChecksum struct {
	Value        string
	ChecksumType string
}

type Location struct {
	Text      string
	Latitude  string
	Longitude string
}

var l logrus.FieldLogger = logger.Get("metadata")

type Metadata map[string]string
type Mapper func(data Metadata) Metadata

func (t Type) ConvertToMetaData() Metadata {
	// TODO: Some of these nested types will result in error if empty
	m := Metadata{}

	m.Set(AccountName, t.AccountName)
	m.Set(CaseNumber, t.CaseNumber)
	if t.CreatedAt != nil {
		m.Set(CreatedAt, t.CreatedAt.Format(time.RFC3339))
	}
	if t.CapturedAt != nil {
		m.Set(CapturedAt, t.CapturedAt.Format(time.RFC3339))
	}
	if t.Subject != nil {

		sJ, err := json.Marshal(t.Subject)
		if err != nil {
			l.Errorf("There was a problems Marhsalling the subjects-fields")
		} else if len(sJ) > 0 {
			subjects := base64.StdEncoding.EncodeToString(sJ)
			if subjects != "" {
				m.Set(Subjects, string(subjects))
			}
		}
	}

	m.Set(ClientMediaId, t.ClientMediaId)
	m.Set(CreatorSurname, t.Creator.LastName)
	m.Set(CreatorDistrict, t.Creator.District)
	m.Set(GroupID, t.GroupId)
	m.Set(GroupName, t.GroupName)
	m.Set(ParentDescription, t.Parent.Description)
	m.Set(ParentName, t.Parent.Name)
	m.Set(Duration, t.Duration)
	m.Set(FileType, t.FileType)
	m.Set(DisplayName, t.DisplayName)
	m.Set(Description, t.Description)
	m.Set(Checksum, t.Checksum.Value)
	m.Set(ChecksumType, t.Checksum.ChecksumType)
	m.Set(Filename, t.FileName)
	m.Set(ExtId, t.ExtId)
	m.Set(LocationText, t.Location.Text)
	m.Set(ParentId, t.Parent.Id)
	m.Set(Longitude, t.Location.Longitude)
	m.Set(Latitude, t.Location.Latitude)
	m.Set(UserId, t.UserId)
	m.Set(EquipmentID, t.EquipmentId)
	m.Set(InterviewType, t.InterviewType)
	m.Set(Notes, t.Notes)
	// TODO:; Map atatchments, booksmarks.
	for key, val := range m {
		if val == "" {
			delete(m, key)
		}
	}
	return m

}

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

func (m Metadata) ConvertToType() Type {
	return Type{
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

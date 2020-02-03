package metadata

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// The root-metadata
//
// All time.time-objects should be ISO-8601-compliant, and include the clients time-offset.
// The required metadata depends on the client, and can be retrieved by either visiting Proxy in the browser
// or doing a OPTIONS-request to its root.
type UploadMetadata struct {
	// A unique identifier for the user in the backend-system
	UserId string `xml:",omitempty"`
	// The parent of the current media.
	Parent *Parent `xml:",omitempty"`
	// The time of which the media was created in the backend-database
	CreatedAt *time.Time `xml:",omitempty"`
	// The time of which the media was created by the user, on the client.
	CapturedAt *time.Time `xml:",omitempty"`
	// The fileType, as in MimeType. Example: 'image/jpeg' or 'video/mp4'
	FileType string `xml:",omitempty"`
	// A short description of the current file, submitted by the user
	DisplayName string `xml:",omitempty"`
	// A longer description of the current file, submitted by the uer.
	Description string `xml:",omitempty"`
	// Any checksums already calculated by the client.
	Checksum *MetaChecksum `xml:",omitempty"`
	// The filename,
	FileName string `xml:",omitempty"`
	// Tags
	Tags []string `xml:",omitempty"`
	// The backend-database ID of the current file. Provide if updating metadata.
	ExtId string `xml:",omitempty"`
	// A case-number returned by the user
	CaseNumber string `xml:",omitempty"`
	// Duration of the media, if audio or video. Int64. Should be in milliseconds when sending.
	Duration string `xml:",omitempty"`
	// The creator of the current file, as in the current user, interviewer etc.
	Creator *Creator `xml:",omitempty"`
	// The location of the captured media
	Location *Location `xml:",omitempty"`
	// Any subjects in the captured media.
	Subject *[]Person `json:"subjects,omitempty" xml:",omitempty"`
	// TBD
	AccountName string `xml:",omitempty"`
	// TBD
	EquipmentId string `xml:",omitempty"`
	// TBD
	InterviewType string      `xml:",omitempty"`
	Bookmarks     *[]Bookmark `xml:",omitempty"`
	// TBD
	Notes string `xml:",omitempty"`
	// A unique identifier of the file on the client.
	ClientMediaId string `xml:",omitempty"`
	// ID of any backend-provided Group-id
	GroupId string `xml:",omitempty"`
	// Name of any backend-group. Providing it will c create a groupName, if supported by the backend.
	GroupName string `xml:",omitempty"`
	// Any custom-field. Should only be used for customer-specific fields that do not fit in any other field. Before use, please request Indico to add your required fields.
	Etc *[]Etc `json:"etc,omitempty xml:Etc"`
}

type Person struct {
	FirstName string `json:"firstName" xml:",omitempty"`
	LastName  string `json:"lastName" xml:",omitempty"`
	Id        string `json:"id" xml:",omitempty"`
	// Date of birth
	Dob         string `json:"dob" xml:",omitempty"`
	Gender      string `json:"gender" xml:",omitempty"`
	Nationality string `json:"nationality" xml:",omitempty"`
	Workplace   string `json:"workplace" xml:",omitempty"`
	// TBD
	Status    string `json:"status" xml:",omitempty"`
	Address   string `json:"address" xml:",omitempty"`
	ZipCode   int    `json:"zip" xml:",omitempty"`
	Country   string `json:"country" xml:",omitempty"`
	WorkPhone string `json:"workPhone" xml:",omitempty"`
	Phone     string `json:"phone" xml:",omitempty"`
	Mobile    string `json:"mobile" xml:",omitempty"`
	// TBD
	Present bool `json:"isPresent" xml:",omitempty"`
}

type Parent struct {
	Id          string `xml:",omitempty"`
	Name        string `xml:",omitempty"`
	Description string `xml:",omitempty"`
}

type Creator struct {
	District string `xml:",omitempty"`
	Person   `xml:",omitempty"`
}

type MetaChecksum struct {
	// The raw value of any checksum
	Value string `xml:",omitempty"`
	// SHA256, MD5, Blake3, CRC. Please advice with Indico before use.
	ChecksumType string `xml:",omitempty"`
}

type Location struct {
	// The text for the current location, like the current address, city, etc.
	Text string `xml:",omitempty"`
	// Geo-location
	Latitude string `xml:",omitempty"`
	/// Geo-location
	Longitude string `xml:",omitempty"`
}

// Stores a nested value as base64-encoded json.
func unwrap(m *Metadata, t interface{}, s string) {
	if t != nil {
		sJ, err := json.Marshal(t)
		if err != nil {
			l.Errorf("There was a problem Marshalling the '%s'-field", s)
		} else if len(sJ) > 0 {
			b := base64.StdEncoding.EncodeToString(sJ)
			if b != "" {
				m.Set(s, b)
			}
		}
	}

}

func (t UploadMetadata) ConvertToMetaData() Metadata {
	// TODO: Some of these nested types will result in error if empty
	m := Metadata{}

	m.Set(AccountName, t.AccountName)
	m.Set(CaseNumber, t.CaseNumber)
	unwrap(&m, t.Bookmarks, Bookmarks)
	unwrap(&m, t.Subject, Subjects)
	unwrap(&m, t.Etc, Etcetera)
	m.Set(ClientMediaId, t.ClientMediaId)
	m.Set(GroupID, t.GroupId)
	m.Set(GroupName, t.GroupName)
	m.Set(Duration, t.Duration)
	m.Set(FileType, t.FileType)
	m.Set(DisplayName, t.DisplayName)
	m.Set(Description, t.Description)
	m.Set(Filename, t.FileName)
	m.Set(ExtId, t.ExtId)
	if t.CreatedAt != nil {
		m.Set(CreatedAt, t.CreatedAt.Format(time.RFC3339))
	}
	if t.CapturedAt != nil {
		m.Set(CapturedAt, t.CapturedAt.Format(time.RFC3339))
	}
	if t.Creator != nil {
		m.Set(CreatorSurname, t.Creator.LastName)
		m.Set(CreatorDistrict, t.Creator.District)
	}
	if t.Tags != nil {
		m.Set(Tags, strings.Join(t.Tags, ","))
	}
	if t.Checksum != nil {
		m.Set(Checksum, t.Checksum.Value)
		m.Set(ChecksumType, t.Checksum.ChecksumType)

	}
	if t.Location != nil {
		m.Set(LocationText, t.Location.Text)
		m.Set(Longitude, t.Location.Longitude)
		m.Set(Latitude, t.Location.Latitude)

	}
	if t.Parent != nil {
		m.Set(ParentName, t.Parent.Name)
		m.Set(ParentDescription, t.Parent.Description)
		m.Set(ParentId, t.Parent.Id)

	}
	m.Set(UserId, t.UserId)
	m.Set(EquipmentID, t.EquipmentId)
	m.Set(InterviewType, t.InterviewType)
	m.Set(Notes, t.Notes)
	// TODO:; Map attachments, bookmarks.
	for key, val := range m {
		if val == "" {
			delete(m, key)
		}
	}
	return m

}

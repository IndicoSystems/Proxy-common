package metadata

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// The root-metadata
//
// All time.time-objects should be ISO-8601-compliant, and include the clients time-offset.
// The required metadata depends on the client, and can be retrieved by either visiting Proxy in the browser
// or doing a OPTIONS-request to its root.
type UploadMetadata struct {
	// A unique identifier for the user in the backend-system
	UserId string
	// The parent of this media.
	Parent Parent
	// The time of which the media was created in the backend-database
	CreatedAt *time.Time
	// The time of which the media was created by the user, on the client.
	CapturedAt *time.Time
	// The fileType, as in MimeType. Example: 'image/jpeg' or 'video/mp4'
	FileType string
	// A short description of the current file, submitted by the user
	DisplayName string
	// A longer description of the current file, submitted by the uer.
	Description string
	// Any checksums already calculated by the client.
	Checksum MetaChecksum
	// The filename,
	FileName string
	// The backend-database ID of the current file. Provide if updating metadata.
	ExtId string
	// A case-number returned by the user
	CaseNumber string
	// Duration of the media, if audio or video. Int64. Should be in milliseconds when sending.
	Duration string
	// The creator of the current file, as in the current user, interviewer etc.
	Creator Creator
	// The location of the captured media
	Location Location
	// Any subjects in the captured media.
	Subject []Person
	// TBD
	AccountName string
	// TBD
	EquipmentId string
	// TBD
	InterviewType string
	// TBD
	Bookmarks string
	// TBD
	Attachments string
	// TBD
	Notes string
	// A unique identifier of the file on the client.
	ClientMediaId string
	// Id of any backend-provided Group-id
	GroupId string
	// Name of any backend-group. Providing it will create a groupName, if supported by the backend.
	GroupName string
	// Any custom-field. Should only be used for customer-specific fields that do not fit in any other field. Before use, please request Indico to add your required fields.
	Etc map[string]interface{} `json:"etc,omitempty"`
	// Reserved for SSN-compliant JSON from legacy-systems. Beware, adding this field to a request will change the behaviour of Proxy for the upload.
	// This is necessary because of the custom-behaviour of the SSN-data.
	SSN *map[string]interface{} `json:"ssn,omitempty"`
}

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Id        string `json:"id"`
	// Date of birth
	Dob         string `json:"dob"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Workplace   string `json:"workplace"`
	// TBD
	Status    string `json:"status"`
	Address   string `json:"address"`
	ZipCode   int    `json:"zip"`
	Country   string `json:"country"`
	WorkPhone string `json:"workPhone"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
	// TBD
	Present bool `json:"isPresent"`
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

type MetaChecksum struct {
	// The raw value of any checksum
	Value string
	// SHA256, MD5, Blake3, CRC. Please advice with Indico before use.
	ChecksumType string
}

type Location struct {
	// The text for the current location, like the current address, city, etc.
	Text string
	// Geo-location
	Latitude string
	/// Geo-location
	Longitude string
}

func (t UploadMetadata) ConvertToMetaData() Metadata {
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
			l.Errorf("There was a problem Marhsalling the subjects-field")
		} else if len(sJ) > 0 {
			subjects := base64.StdEncoding.EncodeToString(sJ)
			if subjects != "" {
				m.Set(Subjects, string(subjects))
			}
		}
	}
	if t.Etc != nil {
		sJ, err := json.Marshal(t.Etc)
		if err != nil {
			l.Errorf("There was a problem Marhsalling the Etc-fields")
		} else if len(sJ) > 0 {
			etc := base64.StdEncoding.EncodeToString(sJ)
			if etc != "" {
				m.Set(Etcetera, string(etc))
			}
		}
	}
	if t.SSN != nil {
		sJ, err := json.Marshal(t.SSN)
		if err != nil {
			l.Errorf("There was a problem Marhsalling the SSN-field")
		} else if len(sJ) > 0 {
			ssn := base64.StdEncoding.EncodeToString(sJ)
			if ssn != "" {
				m.Set(SSN, string(ssn))
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

package metadata

import (
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
	Checksum *[]MetaChecksum `xml:",omitempty"`
	// The filename,
	FileName string `xml:",omitempty"`
	// Tags
	Tags []string `xml:",omitempty"`
	// The backend-database ID of the current file. Provide if updating metadata.
	ExtId string `xml:",omitempty"`
	// A case-number returned by the user
	CaseNumber string `xml:",omitempty"`
	// Duration of the media, if audio or video. Int64. Should be in milliseconds when sending.
	Duration int64 `xml:",omitempty"`
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

type GenderType string

func ToGender(s string) GenderType {
	switch strings.ToLower(s) {
	case "male", "m":
		return GenderMale
	case "female", "f":
		return GenderFemale
	case "":
		return GenderUnspecified
	}
	l.Warn("Could not assing the string '%s' to type gender", s)
	return GenderOther

}

const (
	GenderFemale      GenderType = "Female"
	GenderMale        GenderType = "Male"
	GenderOther       GenderType = "Other"
	GenderUnspecified GenderType = ""
)

type Person struct {
	FirstName string `json:"firstName" xml:",omitempty"`
	LastName  string `json:"lastName" xml:",omitempty"`
	Id        string `json:"id" xml:",omitempty"`
	// Date of birth
	Dob         time.Time  `json:"dob" xml:",omitempty"`
	Gender      GenderType `json:"gender" xml:",omitempty"`
	Nationality string     `json:"nationality" xml:",omitempty"`
	Workplace   string     `json:"workplace" xml:",omitempty"`
	// TBD
	Status    string `json:"status" xml:",omitempty"`
	WorkPhone string `json:"workPhone" xml:",omitempty"`
	Phone     string `json:"phone" xml:",omitempty"`
	Mobile    string `json:"mobile" xml:",omitempty"`
	// TBD
	Present bool `json:"isPresent" xml:",omitempty"`
	Location
}

type Parent struct {
	Id          string `xml:",omitempty"`
	Name        string `xml:",omitempty"`
	Description string `xml:",omitempty"`
}

type Creator struct {
	// Can be an identifier, like an officer's badge-id, or a userId in the system.
	SysId  string `xml:",omitempty"`
	Person `xml:",omitempty"`
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
	// Geo-location
	Longitude string `xml:",omitempty"`
	Address   string `xml:",omitempty"`
	Address2  string `xml:",omitempty"`
	ZipCode   string `xml:",omitempty"`
	PostArea  string `xml:",omitempty"`
	Country   string `xml:",omitempty"`
}

func (t UploadMetadata) ConvertToMetaData() Metadata {
	m := Metadata{}
	m.replaceUploadMetadata(t)
	return m
}

func CreateSampleData() UploadMetadata {
	now := time.Now().Round(time.Second)
	um := UploadMetadata{
		"user",
		&Parent{
			"1234",
			"Burglar",
			"Break-in downtown",
		},
		&now,
		&now,
		"video/mp4",
		"Interview with witness",
		"Witness describing the event",
		&[]MetaChecksum{
			{
				"c013d16a335e2e40edf7d91d2c1f48930e52f3b76a5347010ed25a2334cee872",
				"SHA256",
			},
			{
				"fc02353cb44eb5113a239105daa15c465a5ca57ac2869ea0b381f6f871d22441",
				"SHA3-256",
			},
		},
		"recording-123.mp4",
		[]string{"robbery", "masked", "villain"},
		"1234",
		"C6288",
		int64((44*time.Minute + 8*time.Second + 36*time.Millisecond) / time.Millisecond),
		&Creator{
			"Downtown district",
			Person{
				"Jane",
				"Doe",
				"sk166622",
				time.Date(1977, 3, 4, 0, 0, 0, 0, time.UTC),
				GenderFemale,
				"GBR",
				"Fictive Police Department",
				"",
				"321",
				"321",
				"321",
				true,
				Location{
					"The red house down the street",
					"1.23456",
					"2.34567",
					"Street-road 3",
					"...",
					"SX 6978923",
					"Downtown",
					"GBR",
				},
			},
		},
		&Location{
			"The yellow house down the street",
			"1.23456",
			"2.34567",
			"Street-road 8",
			"...",
			"SX 6978923",
			"Downtown",
			"GBR",
		},
		&[]Person{
			{
				"Burger",
				"Beagle",
				"176-176",
				time.Date(1951, 11, 4, 0, 0, 0, 0, time.UTC),
				GenderMale,
				"USA",
				"Jail",
				"Suspect",
				"123",
				"123",
				"123",
				false,
				Location{
					"Jailcell 3",
					"4.23456",
					"4.34567",
					"Jail",
					"...",
					"SX 6978923",
					"Jail",
					"USA",
				},
			},
			{
				"Daisy",
				"Duck",
				"abc",
				time.Date(1940, 6, 7, 0, 0, 0, 0, time.UTC),
				GenderFemale,
				"USA",
				"Unknown",
				"Witness",
				"123",
				"123",
				"123",
				true,
				Location{
					"",
					"2.23456",
					"2.34567",
					"Street 4",
					"...",
					"SX 6978923",
					"Duckburg",
					"USA",
				},
			},
		},
		"acc",
		"iPhone 20",
		"Witness",
		&[]Bookmark{
			{
				now,
				"abc123",
				"Stressed interviewee",
				63000,
				112000,
			},
		},
		"Daisy witnessed the crime, and says she identified Burgar Beagle.",
		"abc123",
		"multicapture:12345",
		"MutlipleViews",
		&[]Etc{
			{
				"clothing",
				"c12422D3",
				"formKeys.clothing",
				"Main-subjects clothing",
				"A blue dress",
				false,
				"String",
			},
			{
				"mood",
				"c12422G6",
				"formKeys.mood",
				"Main-subjects mood",
				"Scared, stressed",
				true,
				"String",
			},
			{
				"countOfFingers",
				"c12422G7",
				"formKeys.countOfFingers",
				"Main-subjects number of fingers",
				"3",
				true,
				"Int",
			},
		},
	}
	return um
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

// Bookmark object belonging to a certain recording.
type Bookmark struct {
	CreationDate time.Time
	ID           string
	Title        string
	// Position, in milliseconds
	StartPosition int
	// EndPosition, in milliseconds
	EndPosition int
}

// Can be used to quickly validate that certain fields are not empty
func (um UploadMetadata) ValidateRequiredFields(r []string) (missing []string) {
	for _, key := range r {
		switch key {
		case "userid":
			if um.UserId == "" {
				missing = append(missing, key)
			}
			// displayname is allowed to be empty
		case "displayname":
			l.WithField("field", key).Warn("Invalid Required-field")
		case "filetype":
			if um.FileType == "" {
				missing = append(missing, key)
			}
		case "filename":
			if um.FileName == "" {
				missing = append(missing, key)
			}
		case "checksum":
			if um.Checksum == nil {
				missing = append(missing, key)
			}
		case "createdat":
			if um.CreatedAt == nil {
				missing = append(missing, key)
			}
		case "parentid":
			if um.Parent == nil {
				missing = append(missing, key)
			}
			if um.Parent.Id == "" {
				missing = append(missing, key)
			}
		case "parentname":
			if um.Parent == nil {
				missing = append(missing, key)
			}
			if um.Parent.Name == "" {
				missing = append(missing, key)
			}
		default:
			l.Fatalf("Field is not defined in ValidateRequiredFields: '%s'", key)
		}
	}
	return
}

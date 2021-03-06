package metadata

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

// UploadMetadata is the internal representation for Metadata on uploads.
//
// All time.time-objects should be ISO-8601-compliant, and include the clients time-offset.
// The required metadata depends on the client, and can be retrieved by either visiting Proxy in the browser
// or doing a OPTIONS-request to its root.
type UploadMetadata struct {
	// A unique identifier for the user in the backend-system
	UserId string `json:"userId"`
	// Active-Directory-user, if available
	AdSid string `json:"adSid"`
	// Active-Directory-user, if available, in the format 'user@domainname@
	AdLogin string `json:"adLogin"`
	// The parent of the current media.
	Parent Parent `json:"parent"`
	// The time of which the media was created in the backend-database
	CreatedAt *time.Time `json:"createdAt"`
	// The time of which the media was last updated in the backend-database
	UpdatedAt *time.Time `json:"updatedAt"`
	// Is the item marked as archived, meaning it is marked to be deleted (soft-deleted). If null, the item is
	// not scheduled for deletion, if a date is set, the item is marked for deletion at that time.
	ArchiveAt *time.Time `json:"archiveAt"`
	// The date of which the item was marked as completed, E.g. the case was closed.
	CompletedAt *time.Time `json:"completedAt"`
	// The time of which the media was created by the user, on the client.
	CapturedAt *time.Time `json:"capturedAt"`
	// The fileType, as in MimeType. Example: 'image/jpeg' or 'video/mp4'
	FileType string `json:"fileType"`
	// A short description of the current file, submitted by the user
	FileSize    int64  `json:"fileSize"`
	DisplayName string `json:"displayName"`
	// A longer description of the current file, submitted by the uer.
	Description string `json:"description"`
	// Any checksums already calculated by the client.
	Checksum []MetaChecksum `json:"checksum"`
	// The filename,
	FileName string `json:"fileName"`
	// Tags
	Tags []string `json:"tags"`
	// The backend-database ID of the current file. Provide if updating metadata.
	ExtId string `json:"extId"`
	// A case-number returned by the user
	CaseNumber string `json:"caseNumber"`
	// Duration of the media, if audio or video. Int64. Should be in milliseconds when sending.
	Duration int64 `json:"duration"`
	// The creator of the current file, as in the current user, interviewer etc.
	Creator Creator `json:"creator"`
	// The location of the captured media
	Location Location `json:"location"`
	// Any subjects in the captured media.
	Subject []Person `json:"subject"`
	// TBD
	EquipmentId string `json:"equipmentId"`
	// TBD
	InterviewType   string          `json:"interviewType"`
	Bookmarks       []Bookmark      `json:"bookmarks"`
	RecordingEvents RecordingEvents `json:"recordingEvents"`
	Annotations     []Annotation    `json:"annotations"`
	// TBD
	Notes string `json:"notes"`
	// A unique identifier of the file on the client.
	ClientMediaId string `json:"clientMediaId"`
	// ID of any backend-provided Group-id
	GroupId string `json:"groupId"`
	// Name of any backend-group. Providing it will c create a groupName, if supported by the backend.
	GroupName string `json:"groupName"`
	// Any custom-field. Should only be used for customer-specific fields that do not fit in any other field. Before use, please request Indico to add your required fields.
	FormFields []FormFields `json:"formFields"`
	// Transcribed audio/video-details
	Transcription []Utterance `json:"transcription"`
}

// RecordingEvent uses
type RecordingEvent struct {
	Id string
	// Optional
	Title string
	// Optional
	Description string
	// Required
	Kind EventKind
	// Required
	Start time.Time
	// Optional
	StartFrame int64
	// Optional
	End *time.Time
	// Optional
	EndFrame *int64
	// Optional
	EndKind *EventKind
	// Optional
	EndId *string
}

type EventKind string
type RecordingEvents []RecordingEvent
type UtteranceType string

const (
	RecordingEventRecording    EventKind     = "Recording"
	RecordingEventIncomingCall EventKind     = "IncomingCall"
	RecordingEventLowPower     EventKind     = "LowPower"
	RecordingEventLowMemory    EventKind     = "LowMemory"
	RecordingEventFatal        EventKind     = "Fatal"
	RecordingEventUnknown      EventKind     = "Unknown"
	RecordingEventUserEnded    EventKind     = "UserEnded"
	RecordingEventStarted      EventKind     = "Started"
	RecordingEventPaused       EventKind     = "Paused"
	Saying                     UtteranceType = "Saying"
	Event                      UtteranceType = "Event"
)

var (
	ErrRecordingEventInvalidOrder error = errors.New("the order of the events are logically correct")
)

func (re RecordingEvents) WithEndingsIgnoreErr() RecordingEvents {
	RE, _ := re.WithEndings()
	return RE
}

// Returns a new compact RecordingEvents by combining RecordingEvents which belong together in Start/End
// If the events are already using the End-fields, this has no effect
func (re RecordingEvents) WithEndings() (RecordingEvents, error) {
	if len(re) <= 1 {
		return re, nil
	}
	RE := RecordingEvents{}
	var skipIndexex []int
	//last := re[len(re)-1]
outer:
	for i, current := range re {
		if current.End != nil {
			continue
		}
		// Some elements, like Pause should be skipped if Endtime is already put into the Starting type
		for _, indexex := range skipIndexex {
			if i == indexex {
				continue outer
			}
		}
		newEvent := RecordingEvent{
			Id:          current.Id,
			Title:       current.Title,
			Description: current.Description,
			Kind:        current.Kind,
			Start:       current.Start,
			StartFrame:  current.StartFrame,
			End:         *(&current.End),
			EndFrame:    *(&current.EndFrame),
			EndKind:     *(&current.EndKind),
		}
		switch current.Kind {
		// Only these events can have an ending
		case RecordingEventRecording, RecordingEventStarted:
			break
		default:
			RE = append(RE, newEvent)
			continue
		}
		// Loop from the next item to the end, looking for ends for these events
		for j := i + 1; j < len(re); j++ {
			next := re[j]
			switch next.Kind {
			case RecordingEventRecording, RecordingEventStarted:
				return RE, fmt.Errorf("%w %s cannot follow %s without a RecordingEvent between closing the first current", ErrRecordingEventInvalidOrder, next.Kind, current.Kind)
			}
			skipIndexex = append(skipIndexex, j)
			newEvent.EndKind = &next.Kind
			newEvent.EndFrame = &next.StartFrame
			newEvent.End = &next.Start
			newEvent.EndId = &next.Id
			RE = append(RE, newEvent)
			break
		}

	}
	return RE, nil
}

type Utterance struct {
	Type      UtteranceType `json:"type"`
	Person    Person        `json:"person"`
	EventKind string        `json:"eventKind"`
	Source    string        `json:"source"`
	Text      string        `json:"text"`
	// Position, in milliseconds
	StartPosition int `json:"startPosition"`
	// EndPosition, in milliseconds
	EndPosition int `json:"endPosition"`
}

type GenderType string

func ToGender(s string) GenderType {
	switch strings.ToLower(s) {
	case "male", "m", "mann":
		return GenderMale
	case "female", "f", "kvinne":
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
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Id        string `json:"id"`
	// Date of birth
	Dob         *time.Time `json:"dob"`
	Gender      GenderType `json:"gender"`
	Nationality string     `json:"nationality"`
	Workplace   string     `json:"workplace"`
	// TBD
	Status    string `json:"status"`
	WorkPhone string `json:"workPhone"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
	// TBD
	Present bool `json:"isPresent"`
	Location
}

type Parent struct {
	Id          string     `json:"id"`
	BatchId     string     `json:"batchId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type Creator struct {
	// Can be an identifier, like an officer's badge-id, or a userId in the system.
	SysId      string `json:"sysId"`
	Attributes map[string]string
	Person
}

type MetaChecksum struct {
	// The raw value of any checksum
	Value string `json:"value"`
	// SHA256, MD5, Blake3, CRC. Please advice with Indico before use.
	ChecksumType string `json:"checksumType"`
}

type Location struct {
	// The text for the current location, like the current address, city, etc.
	Text string `json:"text"`
	// Geo-location
	Latitude float64 `json:"latitude"`
	// Geo-location
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	Address2  string  `json:"address2"`
	ZipCode   string  `json:"zipCode"`
	PostArea  string  `json:"postArea"`
	Country   string  `json:"country"`
	Accuracy  float64 `json:"accuracy"`
	Altitude  float64 `json:"altitude"`
}

func (um UploadMetadata) ConvertToMetaData() Metadata {
	m := Metadata{}
	m.ReplaceUploadMetadata(um)
	return m
}
func (um UploadMetadata) PrintFieldsSafely() map[string]interface{} {
	return map[string]interface{}{
		"Checksum":      um.Checksum,
		"FileSize":      um.FileSize,
		"ClientMediaId": um.ClientMediaId,
		"FileType":      um.FileType,
		"ExtId":         um.ExtId,
		"UserID":        um.UserId,
	}
}
func (um UploadMetadata) PrintFieldsDebug() map[string]interface{} {
	fields := um.PrintFieldsSafely()
	fields["GroupId"] = um.GroupId
	fields["CaseNumber"] = um.CaseNumber
	fields["Parent"] = um.Parent
	fields["AdLogin"] = um.AdLogin
	fields["AdSid"] = um.AdSid
	fields["Creator"] = um.Creator
	fields["DisplayName"] = um.DisplayName
	fields["Duration"] = um.Duration
	fields["Subject"] = um.Subject
	fields["Tags"] = um.Tags
	return fields
}

func createDate(y int, m time.Month, d int) *time.Time {
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return &date
}

func CreateSampleData() UploadMetadata {
	now := time.Now().Round(time.Second)
	//dob := time.Date(1951, 11, 4, 0, 0, 0, 0, time.UTC)
	um := UploadMetadata{
		UserId:  "user",
		AdSid:   "S-1-5-21-1111111111-2222222222-333333333-1001",
		AdLogin: "user@domainame",
		Parent: Parent{
			"all-metadata-test",
			"",
			"Burglar",
			"Break-in downtown",
			&now,
			&now,
		},
		CreatedAt:   &now,
		UpdatedAt:   &now,
		ArchiveAt:   &now,
		CompletedAt: &now,
		CapturedAt:  &now,
		FileType:    "video/mp4",
		FileSize:    16822,
		DisplayName: "Interview with witness",
		Description: "Witness describing the event",
		Checksum: []MetaChecksum{
			{
				"c013d16a335e2e40edf7d91d2c1f48930e52f3b76a5347010ed25a2334cee872",
				"SHA256",
			},
			{
				"fc02353cb44eb5113a239105daa15c465a5ca57ac2869ea0b381f6f871d22441",
				"SHA3-256",
			},
		},
		FileName:   "recording-123.mp4",
		Tags:       []string{"robbery", "masked", "villain"},
		ExtId:      "1234",
		CaseNumber: "C6288",
		Duration:   int64((44*time.Minute + 8*time.Second + 36*time.Millisecond) / time.Millisecond),
		Creator: Creator{
			"Downtown district",
			map[string]string{
				"any-string": "is allowed",
			},
			Person{
				"Jane",
				"Doe",
				"sk166622",
				createDate(1977, 3, 4),
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
					1.23456,
					2.34567,
					"Street-road 3",
					"...",
					"SX 6978923",
					"Downtown",
					"GBR",
					44.0,
					8.33,
				},
			},
		},
		Location: Location{
			"The yellow house down the street",
			1.23456,
			2.34567,
			"Street-road 8",
			"...",
			"SX 6978923",
			"Downtown",
			"GBR",
			44.0,
			8,
		},
		Subject: []Person{
			{
				"Burger",
				"Beagle",
				"176-176",
				createDate(1951, 11, 4),
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
					4.23456,
					4.34567,
					"Jail",
					"...",
					"SX 6978923",
					"Jail",
					"USA",
					500,
					0,
				},
			},
			{
				"Daisy",
				"Duck",
				"abc-daisy",
				createDate(1940, 6, 7),
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
					2.23456,
					2.34567,
					"Street 4",
					"...",
					"SX 6978923",
					"Duckburg",
					"USA",
					44.0,
					2.53,
				},
			},
		},
		EquipmentId:   "iPhone 20",
		InterviewType: "Witness",
		Bookmarks: []Bookmark{
			{
				now,
				"abc123",
				"Stressed interviewee",
				63000,
				112000,
			},
		},
		RecordingEvents: RecordingEvents{
			{
				Kind:  RecordingEventRecording,
				Start: *createDate(2021, 5, 12),
			},
			{
				Kind:  RecordingEventPaused,
				Start: createDate(2021, 5, 12).Add(4 * time.Minute),
			},
			{
				Kind:  RecordingEventStarted,
				Start: createDate(2021, 5, 12).Add(20 * time.Minute),
			},
			{
				Kind:  RecordingEventUserEnded,
				Start: createDate(2021, 5, 12).Add(1 * time.Hour),
			},
		}.WithEndingsIgnoreErr(),
		Annotations: []Annotation{{
			now,
			"abc",
			"abc",
			"cencor",
			map[string]string{
				"censorType": "beep",
			},
			0,
			0,
			0,
			0,
			18000,
			32000,
		}},
		Notes:         "Daisy witnessed the crime, and says she identified Burgar Beagle.",
		ClientMediaId: "abc123",
		GroupId:       "multicapture:12345",
		GroupName:     "MutlipleViews",
		FormFields: []FormFields{
			{
				"clothing",
				"c12422D3",
				"formKeys.clothing",
				"Main-subjects clothing",
				"A blue dress",
				false,
				"String",
				ValidationRule{
					3,
					300,
				},
			},
			{
				"mood",
				"c12422G6",
				"formKeys.mood",
				"Main-subjects mood",
				"Scared, stressed",
				true,
				"String",
				ValidationRule{
					3,
					300,
				},
			},
			{
				"countOfFingers",
				"c12422G7",
				"formKeys.countOfFingers",
				"Main-subjects number of fingers",
				"3",
				true,
				"Int",
				ValidationRule{
					0,
					24,
				},
			},
		},
		Transcription: []Utterance{
			{
				Saying,
				Person{
					Id: "abc-daisy",
				},
				"",
				"system:azure",
				"It was him, officer",
				1000,
				5000,
			},
			{
				Event,
				Person{
					Id: "abc-daisy",
				},
				"Witness scratches her beak",
				"user:sk166622",
				"",
				6000,
				20000,
			},
			{
				Event,
				Person{
					"John",
					"Doe",
					"sk166628",
					createDate(1972, 8, 9),
					GenderMale,
					"GBR",
					"Fictive Police Department",
					"",
					"322",
					"322",
					"322",
					true,
					Location{
						"The blue house down the street",
						1.23456,
						2.34567,
						"Street-road 8",
						"...",
						"SX 6978923",
						"Downtown",
						"GBR",
						44,
						8.2,
					},
				},
				"Officer interrupts the interview",
				"user:sk166622",
				"",
				18000,
				32000,
			},
			{
				Saying,
				Person{
					Id: "abc-daisy",
				},
				"",
				"system:azure",
				"I say Burgar breaking in.",
				34000,
				38000,
			},
		},
	}
	return um
}

// FormFields er unmapped metadata. In Indico Gateway, these represent an organizational-form.
type FormFields struct {
	// The metadata-key.
	Key string `json:"key"`
	// The id of the field, where available.
	FieldId string `json:"fieldId"`
	// The key used to get the translated VisualName, where available.
	TranslationKey string `json:"translationKey"`
	// The visual name, as reported by the client.
	VisualName string `json:"visualName"`
	// The value of the field, e.g. the user-input
	Value string `json:"value"`
	// Marks whether or not field is required or not.
	Required bool `json:"required"`
	// The kind of data in the Value-field.
	DataType       string         `json:"dataType"`
	ValidationRule ValidationRule `json:"validationRule"`
}

type ValidationRule struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Bookmark object belonging to a certain recording.
type Bookmark struct {
	CreationDate time.Time `json:"creationTime"`
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	// Position, in milliseconds
	StartPosition int `json:"startPosition"`
	// EndPosition, in milliseconds
	EndPosition int `json:"endPosition"`
}

// Annotations can be small figures, highlights, descriptors on an image/vidoe/audio.
// They can be added by the user, or automatically through OCR or Machine Learning.
// For audio, it can also be a way to censor parts of the audio via for instance beeps.
type Annotation struct {
	CreationDate time.Time `json:"createdAt"`
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	// Data will be a computer-friendly structure. The structure itself is not decided yet.
	Data StringMap `json:"data"`
	// Position on the image/video.
	X1 int `json:"x1"`
	X2 int `json:"x2"`
	Y1 int `json:"y1"`
	Y2 int `json:"y2"`
	// Position, in milliseconds
	StartPosition int `json:"startPosition"`
	// EndPosition, in milliseconds
	EndPosition int `json:"endPosition"`
}

type StringMap map[string]string

// StringMap marshals into XML.
func (s StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
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
			if um.Parent.Id == "" {
				missing = append(missing, key)
			}
		case "parentname":
			if um.Parent.Name == "" {
				missing = append(missing, key)
			}
		default:
			l.Fatalf("ToField is not defined in ValidateRequiredFields: '%s'", key)
		}
	}
	return
}

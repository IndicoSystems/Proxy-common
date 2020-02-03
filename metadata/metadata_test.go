package metadata

import (
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestMetadata_Get(t *testing.T) {
	key := "key1"
	value := "value1"

	md := Metadata(map[string]string{
		key: value,
	})

	assert.Equal(t, value, md.get("KEY1"))
}

func TestMetadata_Set(t *testing.T) {
	key := "key1"
	value := "value1"

	md := Metadata(make(map[string]string))

	md.set(key, value)

	assert.Equal(t, value, md[key])
}

func TestMetadata_ConvertToType(t *testing.T) {
	tests := []struct {
		name string
		m    Metadata
		want UploadMetadata
	}{
		// TODO: Add test cases.
		{
			"Should not fail on empty values",
			Metadata{},
			UploadMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetUploadMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUploadMetadata() = \n%+v, \nwant \n%+v", got, tt.want)
			}
		})
	}
}

func TestType_ConvertToMetaData(t1 *testing.T) {
	type fields struct {
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
	}
	createdat := time.Date(2019, 12, 24, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		fields UploadMetadata
		want   Metadata
	}{
		{
			"Should not fail on empty values",
			UploadMetadata{},
			Metadata{},
		},
		{
			"Should not fail on basic values",
			UploadMetadata{
				GroupId:   "1234",
				CreatedAt: &createdat,
			},
			Metadata{
				GroupID:   "1234",
				CreatedAt: "2019-12-24T00:00:00Z",
			},
		},
		{
			"should parse all fields correctly",
			UploadMetadata{
				UserId: "1111",
				Parent: Parent{
					Id:          "1234",
					Name:        "Sigma",
					Description: "With his mavericks",
				},
				CreatedAt:   &createdat,
				CapturedAt:  &createdat,
				FileType:    "image/jpeg",
				DisplayName: "Dogs and floating heads",
				Description: "Illegal weapons",
				Checksum: MetaChecksum{
					Value:        "1234-ABC",
					ChecksumType: "Imaginary",
				},
				FileName:   "01.jpeg",
				ExtId:      "9999",
				CaseNumber: "8888",
				Duration:   "30",
				Creator: Creator{
					District: "Lab",
					Person: Person{
						LastName: "Light",
						//FirstName:   "Thomas",
						//Id:          "1234",
						//Dob:         "20xx",
						//GenderType:      "Male",
						//Nationality: "Canadian",
						//Workplace:   "Robot Institute of Technology",
						//Status:      "",
						//Address:     "",
						//ZipCode:     0,
						//Country:     "",
						//WorkPhone:   "",
						//Phone:       "",
						//Mobile:      "",
						//Present:     true,
					},
				},
				Location: Location{
					Text:      "Some place",
					Latitude:  "8",
					Longitude: "7",
				},
				Subject: []Person{
					{
						FirstName:   "Unknown",
						LastName:    "Sigma",
						Id:          "666",
						Dob:         "unknown",
						Gender:      "Male",
						Nationality: "American",
						Workplace:   "Maverick Hunters",
						Status:      "Suspect",
						Address:     "Unknown",
						ZipCode:     0,
						Country:     "USA",
						WorkPhone:   "6786762",
						Phone:       "1267363",
						Mobile:      "162363",
						Present:     true,
					},
				},
				AccountName:   "Supa-Computa",
				EquipmentId:   "XC66",
				InterviewType: "???",
				//Bookmarks:     "", // TODO: add test for booksmarks, attachments, etc.
				//Attachments:   "",
				Notes:         "Some random notes",
				ClientMediaId: "1234",
				GroupId:       "4321",
				FormFields: map[string]interface{}{
					"string": "batman",
					"array":  []interface{}{"Alice", "Bob", "Eve"},
					"nested": map[string]interface{}{
						"more_strings": "stringy",
						"yes":          false,
					},
				},
			},
			Metadata{
				Subjects:      "W3siZmlyc3ROYW1lIjoiVW5rbm93biIsImxhc3ROYW1lIjoiU2lnbWEiLCJpZCI6IjY2NiIsImRvYiI6InVua25vd24iLCJnZW5kZXIiOiJNYWxlIiwibmF0aW9uYWxpdHkiOiJBbWVyaWNhbiIsIndvcmtwbGFjZSI6Ik1hdmVyaWNrIEh1bnRlcnMiLCJzdGF0dXMiOiJTdXNwZWN0IiwiYWRkcmVzcyI6IlVua25vd24iLCJ6aXAiOjAsImNvdW50cnkiOiJVU0EiLCJ3b3JrUGhvbmUiOiI2Nzg2NzYyIiwicGhvbmUiOiIxMjY3MzYzIiwibW9iaWxlIjoiMTYyMzYzIiwiaXNQcmVzZW50Ijp0cnVlfV0=",
				AccountName:   "Supa-Computa",
				InterviewType: "???",
				//Bookmarks:         "", // TODO: add test for booksmarks, attachments, etc.
				//Attachments:       "",
				Notes:             "Some random notes",
				ClientMediaId:     "1234",
				extId:             "9999",
				CaseNumber:        "8888",
				Duration:          "30",
				CreatedAt:         "2019-12-24T00:00:00Z",
				CapturedAt:        "2019-12-24T00:00:00Z",
				FileType:          "image/jpeg",
				DisplayName:       "Dogs and floating heads",
				Description:       "Illegal weapons",
				UserId:            "1111",
				CreatorSurname:    "Light",
				CreatorDistrict:   "Lab",
				EquipmentID:       "XC66",
				ParentDescription: "With his mavericks",
				ParentId:          "1234",
				ParentName:        "Sigma",
				ChecksumType:      "Imaginary",
				Checksum:          "1234-ABC",
				Latitude:          "8",
				Longitude:         "7",
				LocationText:      "Some place",
				GroupID:           "4321",
				Filename:          "01.jpeg",
				Etcetera:          "eyJhcnJheSI6WyJBbGljZSIsIkJvYiIsIkV2ZSJdLCJuZXN0ZWQiOnsibW9yZV9zdHJpbmdzIjoic3RyaW5neSIsInllcyI6ZmFsc2V9LCJzdHJpbmciOiJiYXRtYW4ifQ==",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			got := t.ConvertToMetaData()
			if diff := deep.Equal(got, tt.want); diff != nil {
				// Because nested values are Base64, lets unwrap them to get a clear picture of the diff
				if diffPerson := deep.Equal(got.GetPerson(), tt.want.GetPerson()); diffPerson != nil {
					t1.Errorf("Person did not match wanted result: \n diff %+v", diffPerson)
				}
				t1.Errorf("Did not match wanted result: \n diff %+v, \nGot: %+v, \nWanted: %+v", diff, got, tt.want)
			}
			m := got.GetUploadMetadata()
			if diff := deep.Equal(m, t); diff != nil {
				t1.Logf("\n'%+v' \n'%+v'", got.GetEtc(), tt.want.GetEtc())
				t1.Errorf("Convert to metadata and back again did not produce the same result\n diff: %+v \nGot: %+v \n Expected result: %+v", diff, m, t)
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t1.Errorf("ConvertToMetaData() => GetUploadMetadata = \n%v, want \n%v", got, tt.want)
			//}

		})
	}
}

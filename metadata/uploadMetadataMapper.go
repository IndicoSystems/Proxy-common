package metadata

import (
	"fmt"
	"github.com/indicosystems/proxy/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

var lf = logger.Get("field-mapper")

type UploadMetadataField = string
type UploadMetadataFieldCondition = string

const (
	UMFCaseNumber         UploadMetadataField = "CaseNumber"
	UMFCreatorFirstName                       = "Creator.FirstName"
	UMFCreatorWorkPlace                       = "Creator.WorkPlace"
	UMFSNotes                                 = "Notes"
	UMFSubjectAccuracy                        = "Subject.Accuracy"
	UMFSubjectAddress                         = "Subject.Address"
	UMFSubjectAddress2                        = "Subject.Address2"
	UMFSubjectAltitude                        = "Subject.Altitude"
	UMFSubjectCountry                         = "Subject.Country"
	UMFSubjectDob                             = "Subject.Dob"
	UMFSubjectFirstName                       = "Person.FirstName"
	UMFSubjectGender                          = "Subject.Gender"
	UMFSubjectId                              = "Subject.Id"
	UMFSubjectLastName                        = "SubjectLastName"
	UMFSubjectLatitude                        = "Subject.Latitude"
	UMFSubjectLongitude                       = "Subject.Longitude"
	UMFSubjectMobile                          = "Subject.Mobile"
	UMFSubjectNationality                     = "Subject.Nationality"
	UMFSubjectPhone                           = "Subject.Phone"
	UMFSubjectPostArea                        = "Subject.PostArea"
	UMFSubjectStatus                          = "Subject.Status"
	UMFSubjectText                            = "Subject.Text"
	UMFSubjectWorkPhone                       = "Subject.WorkPhone"
	UMFSubjectWorkplace                       = "Subject.WorkPlace"
	UMFSubjectZipCode                         = "Subject.ZipCode"
	UMFUserId                                 = "UserId"

	ConditionIfNotBlank UploadMetadataFieldCondition = "IfNotBlank"
	ConditionIfNotSet                                = "IfNotSet"
)

type UMF struct {
	CaseNumber struct{}
}

func FieldMapperConfig() (fm FieldMap) {
	err := viper.UnmarshalKey("fieldMapper", &fm)
	if err != nil {
		lf.Fatalf("could not unmarshal fieldMapper", err)
	}
	return

}

type FieldMapType struct {
	Field     UploadMetadataField
	Condition UploadMetadataFieldCondition
	Args      map[string]string
}

type FieldMap map[string]FieldMapType

func NewFieldMap() FieldMap {
	return FieldMap{}
}

func IsUploadMetadataField(s UploadMetadataField) bool {
	switch s {
	case
		UMFCaseNumber,
		UMFCreatorFirstName,
		UMFCreatorWorkPlace,
		UMFSNotes,
		UMFSubjectAccuracy,
		UMFSubjectAddress,
		UMFSubjectAddress2,
		UMFSubjectAltitude,
		UMFSubjectCountry,
		UMFSubjectDob,
		UMFSubjectFirstName,
		UMFSubjectGender,
		UMFSubjectId,
		UMFSubjectLastName,
		UMFSubjectLatitude,
		UMFSubjectLongitude,
		UMFSubjectMobile,
		UMFSubjectNationality,
		UMFSubjectPhone,
		UMFSubjectPostArea,
		UMFSubjectStatus,
		UMFSubjectText,
		UMFSubjectWorkPhone,
		UMFSubjectWorkplace,
		UMFSubjectZipCode,
		UMFUserId:
		return true
	}
	lf.WithField("s", s).Warn("Not a valid UploadMetadataField")
	return false
}
func IsUploadMetadataFieldCondition(s UploadMetadataFieldCondition) bool {
	switch s {
	case ConditionIfNotBlank, ConditionIfNotSet, "":
		return true
	}
	lf.WithField("s", s).Warn("Not a valid UploadMetadataFieldCondition")
	return false
}

func (fm FieldMap) SetField(s string, f UploadMetadataField, c UploadMetadataFieldCondition, a map[string]string) FieldMap {
	IsUploadMetadataField(f)
	IsUploadMetadataFieldCondition(c)
	fm[s] = FieldMapType{
		Field:     f,
		Condition: c,
		Args:      a,
	}
	return fm
}

func (f FieldMapType) CheckStringCondition(existing string, newValue string) UploadMetadataFieldCondition {
	switch f.Condition {
	case "":
		return ""
	case ConditionIfNotSet:
		if existing != "" {
			return ConditionIfNotSet
		}
	case ConditionIfNotBlank:
		if newValue == "" {
			return ConditionIfNotBlank
		}
	}
	return ""
}

func (um *UploadMetadata) SetField(f FieldMapType, v interface{}) (err error) {
	if s, ok := v.(string); ok {
		return um.SetStringField(f, s)
	}
	err = errors.Errorf("Cannot map key '%s' to type '%T' with value '%+v'", f, v, v)
	return
}

func (um *UploadMetadata) SetStringField(f FieldMapType, s string) (err error) {
	s = strings.TrimSpace(s)
	l := lf.WithField("fieldMapType", f).WithField("value", s)
	switch f.Field {
	case
		UMFCaseNumber:
		if f.CheckStringCondition(um.CaseNumber, s) == "" {
			um.CaseNumber = s
		}
	case UMFCreatorFirstName:
		if f.CheckStringCondition(um.Creator.FirstName, s) == "" {
			um.Creator.FirstName = s
		}
	case UMFCreatorWorkPlace:
		return
	case UMFSNotes:
		if f.CheckStringCondition(um.Notes, s) == "" {
			um.Notes = s
		}
	case UMFSubjectAccuracy:
		if um.Subject[0].Accuracy == 0 {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				lf.Warnf("Could not parse '%s' as float in field '%s'", s, f.Field)
			} else {
				um.Subject[0].Accuracy = i
			}
		}
		return
	case UMFSubjectAddress:
		if f.CheckStringCondition(um.Subject[0].Address, s) == "" {
			um.Subject[0].Address = s
		}
	case UMFSubjectAddress2:
		if f.CheckStringCondition(um.Subject[0].Address2, s) == "" {
			um.Subject[0].Address2 = s
		}
	case UMFSubjectAltitude:
		if um.Subject[0].Altitude == 0 {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				l.Error("Could not parse to float")
			} else {
				um.Subject[0].Altitude = i
			}
		}
	case UMFSubjectCountry:
		if f.CheckStringCondition(um.Subject[0].Country, s) == "" {
			um.Subject[0].Country = s
		}
	case UMFSubjectDob:
		fmt.Println("attemting to parse date", f, s)

		if a, ok := f.Args["layout"]; ok {
			d, err := time.Parse(a, s)
			if err != nil {
				l.WithError(err).Error("Could not parse to date")
			} else {
				um.Subject[0].Dob = &d
			}
		} else {
			l.Warn("Could not get layout from args, therefore I cannot parse")
		}
		return
	case UMFSubjectFirstName:
		if f.CheckStringCondition(um.Subject[0].FirstName, s) == "" {
			um.Subject[0].FirstName = s
		}
	case UMFSubjectGender:
		if um.Subject[0].Gender == GenderUnspecified {
			v := ToGender(s)
			um.Subject[0].Gender = v
		}
		return
	case UMFSubjectId:
		if f.CheckStringCondition(um.Subject[0].Id, s) == "" {
			um.Subject[0].Id = s
		}
	case UMFSubjectLastName:
		if f.CheckStringCondition(um.Subject[0].LastName, s) == "" {
			um.Subject[0].LastName = s
		}
	case UMFSubjectLatitude:
		if um.Subject[0].Latitude == 0 {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				l.Error("Could not parse to float")
			} else {
				um.Subject[0].Latitude = i
			}
		}
	case UMFSubjectLongitude:
		if um.Subject[0].Longitude == 0 {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				l.Error("Could not parse to float")
			} else {
				um.Subject[0].Longitude = i
			}
		}
	case UMFSubjectMobile:
		if f.CheckStringCondition(um.Subject[0].Mobile, s) == "" {
			um.Subject[0].Mobile = s
		}
	case UMFSubjectNationality:
		if f.CheckStringCondition(um.Subject[0].Nationality, s) == "" {
			um.Subject[0].Nationality = s
		}
	case UMFSubjectPhone:
		if f.CheckStringCondition(um.Subject[0].Phone, s) == "" {
			um.Subject[0].Phone = s
		}
	case UMFSubjectPostArea:
		if f.CheckStringCondition(um.Subject[0].PostArea, s) == "" {
			um.Subject[0].PostArea = s
		}
	case UMFSubjectStatus:
		if f.CheckStringCondition(um.Subject[0].Status, s) == "" {
			um.Subject[0].Status = s
		}
	case UMFSubjectText:
		if f.CheckStringCondition(um.Subject[0].Text, s) == "" {
			um.Subject[0].Text = s
		}
	case UMFSubjectWorkPhone:
		if f.CheckStringCondition(um.Subject[0].WorkPhone, s) == "" {
			um.Subject[0].WorkPhone = s
		}
	case UMFSubjectWorkplace:
		if f.CheckStringCondition(um.Subject[0].Workplace, s) == "" {
			um.Subject[0].Workplace = s
		}
	case UMFSubjectZipCode:
		if f.CheckStringCondition(um.Subject[0].ZipCode, s) == "" {
			um.Subject[0].ZipCode = s
		}
	case UMFUserId:
		if f.CheckStringCondition(um.UserId, s) == "" {
			um.UserId = s
		}
	default:
		err = errors.Errorf("Cannot map key '%s' to string with value '%+v'", f.Field, s)
	}
	return
}

func (fm FieldMap) FindFromField(ff FormFields) (FieldMapType, bool) {
	if val, ok := fm[strings.ToLower(ff.Key)]; ok {
		return val, true
	}
	if ff.TranslationKey != "" {
		if val, ok := fm[strings.ToLower(ff.TranslationKey)]; ok {
			return val, true
		}
		l := strings.Split(strings.ToLower(ff.TranslationKey), ".")[0]
		if val, ok := fm[l]; ok {
			return val, true
		}
	}
	if val, ok := fm[strings.ToLower(ff.VisualName)]; ok {
		return val, true
	}
	return FieldMapType{}, false

}

func (um *UploadMetadata) MapFormFields(fMap FieldMap) (err error) {
	for _, f := range um.FormFields {
		fKey, found := fMap.FindFromField(f)

		if !found {
			continue
		}
		err := um.SetStringField(fKey, f.Value)
		if err != nil {
			lf.WithError(err).WithFields(map[string]interface{}{
				"fKey":  fKey,
				"value": f.Value,
			}).Error("Failed in mapping]")
		}
	}
	return
}

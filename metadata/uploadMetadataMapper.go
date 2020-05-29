package metadata

import (
	"bytes"
	"fmt"
	"github.com/indicosystems/proxy/config"
	"github.com/indicosystems/proxy/info"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var lf logrus.FieldLogger = logrus.StandardLogger()

func AssignFieldLogger(l logrus.FieldLogger) {
	lf = l
}

type UploadMetadataField = string
type UploadMetadataFieldCondition = string

const (
	UMFCaseNumber         UploadMetadataField = "casenumber"
	UMFCreatorFirstName                       = "creator.firstname"
	UMFCreatorWorkPlace                       = "creator.workplace"
	UMFSNotes                                 = "notes"
	UMFSubjectAccuracy                        = "subject.accuracy"
	UMFSubjectAddress                         = "subject.address"
	UMFSubjectAddress2                        = "subject.address2"
	UMFSubjectAltitude                        = "subject.altitude"
	UMFSubjectCountry                         = "subject.country"
	UMFSubjectDob                             = "subject.dob"
	UMFSubjectFirstName                       = "subject.firstname"
	UMFTags                                   = "tags"
	UMFSubjectGender                          = "subject.gender"
	UMFSubjectId                              = "subject.id"
	UMFSubjectLastName                        = "subject.lastname"
	UMFSubjectLatitude                        = "subject.latitude"
	UMFSubjectLongitude                       = "subject.longitude"
	UMFSubjectMobile                          = "subject.mobile"
	UMFSubjectNationality                     = "subject.nationality"
	UMFSubjectPhone                           = "subject.phone"
	UMFSubjectPostArea                        = "subject.postarea"
	UMFSubjectStatus                          = "subject.status"
	UMFSubjectText                            = "subject.text"
	UMFSubjectWorkPhone                       = "subject.workphone"
	UMFSubjectWorkplace                       = "subject.workplace"
	UMFSubjectZipCode                         = "subject.zipcode"
	UMFUserId                                 = "userid"

	ConditionIfNotBlank UploadMetadataFieldCondition = "ifnotblank"
	ConditionIfNotSet                                = "ifnotset"

	// These are fields mapped from anything to anything
	UMDescription = "description"
)

type UMF struct {
	CaseNumber struct{}
}

func FieldMapperConfig() (fm FieldMap) {
	err := viper.UnmarshalKey("fieldMapper", &fm)
	if err != nil {
		lf.Fatalf("could not unmarshal fieldMapper", err)
	}
	for _, v := range fm {
		if v.Aliases == nil {
			continue
		}
		for _, key := range v.Aliases {
			fm[strings.ToLower(key)] = FieldMapType{
				v.ToField,
				v.Condition,
				nil,
				v.Layout,
				v.Format,
			}
		}
	}
	return

}
func AnyMapperConfig() (fm map[string]string) {
	err := viper.UnmarshalKey("mapper", &fm)
	if err != nil {
		lf.Fatalf("could not unmarshal mapper", err)
	}
	return

}

type ToFormField struct {
	Key        string
	VisualName string
	Required   bool
	Format     string
	Debug      bool
}

func ToFormFieldMap() (fm []ToFormField) {
	err := viper.UnmarshalKey("toFormFieldMap", &fm)
	if err != nil {
		lf.Fatalf("could not unmarshal to-field-mapper", err)
	}
	return
}

type FieldMapType struct {
	ToField   UploadMetadataField
	Condition UploadMetadataFieldCondition
	Aliases   []string
	// Used for date-layout
	Layout string
	// Format is used as a template
	Format string
}

type FieldMap map[string]FieldMapType

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
	s = strings.ToLower(s)
	switch s {
	case ConditionIfNotBlank, ConditionIfNotSet, "":
		return true
	}
	lf.WithField("s", s).Warn("Not a valid UploadMetadataFieldCondition")
	return false
}

func (fm FieldMap) SetField(s string, f UploadMetadataField, c UploadMetadataFieldCondition, al []string, layout, format string) FieldMap {
	IsUploadMetadataField(f)
	IsUploadMetadataFieldCondition(c)
	fm[s] = FieldMapType{
		f,
		c,
		al,
		layout,
		format,
	}
	return fm
}

func (f FieldMapType) CheckStringCondition(existing string, newValue string) UploadMetadataFieldCondition {
	switch strings.ToLower(f.Condition) {
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

//func (um *UploadMetadata) SetField(f FieldMapType, v interface{}) (err error) {
//	if s, ok := v.(string); ok {
//		return um.SetStringField(f, s)
//	}
//	err = errors.Errorf("Cannot map key '%s' to type '%T' with value '%+v'", f, v, v)
//	return
//}

type locale struct {
	locale string
	//pluralsCardinal        []locales.PluralRule
	//pluralsOrdinal         []locales.PluralRule
	//pluralsRange           []locales.PluralRule
	decimal                string
	group                  string
	minus                  string
	percent                string
	percentSuffix          string
	perMille               string
	timeSeparator          string
	inifinity              string
	currencies             []string // idx = enum of currency code
	currencyPositivePrefix string
	currencyNegativePrefix string
	monthsAbbreviated      []string
	monthsNarrow           []string
	monthsWide             []string
	daysAbbreviated        []string
	daysNarrow             []string
	daysShort              []string
	daysWide               []string
	periodsAbbreviated     []string
	periodsNarrow          []string
	periodsShort           []string
	periodsWide            []string
	erasAbbreviated        []string
	erasNarrow             []string
	erasWide               []string
	timezones              map[string]string
}

var locales = map[string]locale{
	"nb_NO": locale{
		locale: "locale",
		//pluralsCardinal:        []locales.PluralRule{2, 6},
		//pluralsOrdinal:         []locales.PluralRule{6},
		//pluralsRange:           []locales.PluralRule{6},
		decimal:                ",",
		group:                  " ",
		minus:                  "−",
		percent:                "%",
		perMille:               "‰",
		timeSeparator:          ":",
		inifinity:              "∞",
		currencies:             []string{"ADP", "AED", "AFA", "AFN", "ALK", "ALL", "AMD", "ANG", "AOA", "AOK", "AON", "AOR", "ARA", "ARL", "ARM", "ARP", "ARS", "ATS", "AUD", "AWG", "AZM", "AZN", "BAD", "BAM", "BAN", "BBD", "BDT", "BEC", "BEF", "BEL", "BGL", "BGM", "BGN", "BGO", "BHD", "BIF", "BMD", "BND", "BOB", "BOL", "BOP", "BOV", "BRB", "BRC", "BRE", "BRL", "BRN", "BRR", "BRZ", "BSD", "BTN", "BUK", "BWP", "BYB", "BYN", "BYR", "BZD", "CAD", "CDF", "CHE", "CHF", "CHW", "CLE", "CLF", "CLP", "CNH", "CNX", "CNY", "COP", "COU", "CRC", "CSD", "CSK", "CUC", "CUP", "CVE", "CYP", "CZK", "DDM", "DEM", "DJF", "DKK", "DOP", "DZD", "ECS", "ECV", "EEK", "EGP", "ERN", "ESA", "ESB", "ESP", "ETB", "EUR", "FIM", "FJD", "FKP", "FRF", "GBP", "GEK", "GEL", "GHC", "GHS", "GIP", "GMD", "GNF", "GNS", "GQE", "GRD", "GTQ", "GWE", "GWP", "GYD", "HKD", "HNL", "HRD", "HRK", "HTG", "HUF", "IDR", "IEP", "ILP", "ILR", "ILS", "INR", "IQD", "IRR", "ISJ", "ISK", "ITL", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRH", "KRO", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LTT", "LUC", "LUF", "LUL", "LVL", "LVR", "LYD", "MAD", "MAF", "MCF", "MDC", "MDL", "MGA", "MGF", "MKD", "MKN", "MLF", "MMK", "MNT", "MOP", "MRO", "MTL", "MTP", "MUR", "MVP", "MVR", "MWK", "MXN", "MXP", "MXV", "MYR", "MZE", "MZM", "MZN", "NAD", "NGN", "NIC", "NIO", "NLG", "NOK", "NPR", "NZD", "OMR", "PAB", "PEI", "PEN", "PES", "PGK", "PHP", "PKR", "PLN", "PLZ", "PTE", "PYG", "QAR", "RHD", "ROL", "RON", "RSD", "RUB", "RUR", "RWF", "SAR", "SBD", "SCR", "SDD", "SDG", "SDP", "SEK", "SGD", "SHP", "SIT", "SKK", "SLL", "SOS", "SRD", "SRG", "SSP", "STD", "STN", "SUR", "SVC", "SYP", "SZL", "THB", "TJR", "TJS", "TMM", "TMT", "TND", "TOP", "TPE", "TRL", "TRY", "TTD", "TWD", "TZS", "UAH", "UAK", "UGS", "UGX", "USD", "USN", "USS", "UYI", "UYP", "UYU", "UZS", "VEB", "VEF", "VND", "VNN", "VUV", "WST", "XAF", "XAG", "XAU", "XBA", "XBB", "XBC", "XBD", "XCD", "XDR", "XEU", "XFO", "XFU", "XOF", "XPD", "XPF", "XPT", "XRE", "XSU", "XTS", "XUA", "XXX", "YDD", "YER", "YUD", "YUM", "YUN", "YUR", "ZAL", "ZAR", "ZMK", "ZMW", "ZRN", "ZRZ", "ZWD", "ZWL", "ZWR"},
		percentSuffix:          " ",
		currencyPositivePrefix: " ",
		currencyNegativePrefix: " ",
		monthsAbbreviated:      []string{"", "jan.", "feb.", "mar.", "apr.", "mai", "jun.", "jul.", "aug.", "sep.", "okt.", "nov.", "des."},
		monthsNarrow:           []string{"", "J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"},
		monthsWide:             []string{"", "januar", "februar", "mars", "april", "mai", "juni", "juli", "august", "september", "oktober", "november", "desember"},
		daysAbbreviated:        []string{"søn.", "man.", "tir.", "ons.", "tor.", "fre.", "lør."},
		daysNarrow:             []string{"S", "M", "T", "O", "T", "F", "L"},
		daysShort:              []string{"sø.", "ma.", "ti.", "on.", "to.", "fr.", "lø."},
		daysWide:               []string{"søndag", "mandag", "tirsdag", "onsdag", "torsdag", "fredag", "lørdag"},
		periodsAbbreviated:     []string{"a.m.", "p.m."},
		periodsNarrow:          []string{"a", "p"},
		periodsWide:            []string{"a.m.", "p.m."},
		erasAbbreviated:        []string{"f.Kr.", "e.Kr."},
		erasNarrow:             []string{"f.Kr.", "e.Kr."},
		erasWide:               []string{"før Kristus", "etter Kristus"},
		timezones:              map[string]string{"ARST": "argentinsk sommertid", "GMT": "Greenwich middeltid", "MDT": "sommertid for Rocky Mountains (USA)", "WAST": "vestafrikansk sommertid", "ACST": "sentralaustralsk normaltid", "COT": "colombiansk normaltid", "CST": "normaltid for det sentrale Nord-Amerika", "HEPM": "sommertid for Saint-Pierre-et-Miquelon", "WITA": "sentralindonesisk tid", "VET": "venezuelansk tid", "OEZ": "østeuropeisk normaltid", "COST": "colombiansk sommertid", "AWST": "vestaustralsk normaltid", "AEST": "østaustralsk normaltid", "SAST": "sørafrikansk tid", "SGT": "singaporsk tid", "ACWST": "vest-sentralaustralsk normaltid", "WAT": "vestafrikansk normaltid", "ACDT": "sentralaustralsk sommertid", "AWDT": "vestaustralsk sommertid", "ACWDT": "vest-sentralaustralsk sommertid", "WART": "vestargentinsk normaltid", "HNPM": "normaltid for Saint-Pierre-et-Miquelon", "UYST": "uruguayansk sommertid", "WEZ": "vesteuropeisk normaltid", "HAT": "sommertid for Newfoundland", "EAT": "østafrikansk tid", "CHADT": "sommertid for Chatham", "MEZ": "sentraleuropeisk normaltid", "HNT": "normaltid for Newfoundland", "HENOMX": "sommertid for nordvestlige Mexico", "∅∅∅": "sommertid for Amazonas", "CHAST": "normaltid for Chatham", "GFT": "tidssone for Fransk Guyana", "MYT": "malaysisk tid", "HEEG": "østgrønlandsk sommertid", "TMT": "turkmensk normaltid", "MST": "normaltid for Rocky Mountains (USA)", "BOT": "boliviansk tid", "AKST": "alaskisk normaltid", "EDT": "sommertid for den nordamerikanske østkysten", "MESZ": "sentraleuropeisk sommertid", "HKT": "normaltid for Hongkong", "CLT": "chilensk normaltid", "ART": "argentinsk normaltid", "PDT": "sommertid for den nordamerikanske Stillehavskysten", "HEPMX": "sommertid for den meksikanske Stillehavskysten", "AEDT": "østaustralsk sommertid", "WESZ": "vesteuropeisk sommertid", "UYT": "uruguayansk normaltid", "HECU": "cubansk sommertid", "HNPMX": "normaltid for den meksikanske Stillehavskysten", "HNEG": "østgrønlandsk normaltid", "LHDT": "sommertid for Lord Howe-øya", "WIT": "østindonesisk tid", "JDT": "japansk sommertid", "ECT": "ecuadoriansk tid", "HNOG": "vestgrønlandsk normaltid", "HKST": "sommertid for Hongkong", "TMST": "turkmensk sommertid", "HADT": "sommertid for Hawaii og Aleutene", "GYT": "guyansk tid", "BT": "bhutansk tid", "AKDT": "alaskisk sommertid", "HEOG": "vestgrønlandsk sommertid", "EST": "normaltid for den nordamerikanske østkysten", "CLST": "chilensk sommertid", "HAST": "normaltid for Hawaii og Aleutene", "NZDT": "newzealandsk sommertid", "IST": "indisk tid", "CAT": "sentralafrikansk tid", "HNNOMX": "normaltid for nordvestlige Mexico", "OESZ": "østeuropeisk sommertid", "HNCU": "cubansk normaltid", "CDT": "sommertid for det sentrale Nord-Amerika", "WIB": "vestindonesisk tid", "NZST": "newzealandsk normaltid", "LHST": "normaltid for Lord Howe-øya", "WARST": "vestargentinsk sommertid", "ChST": "tidssone for Chamorro", "PST": "normaltid for den nordamerikanske Stillehavskysten", "AST": "normaltid for den nordamerikanske atlanterhavskysten", "ADT": "sommertid for den nordamerikanske atlanterhavskysten", "JST": "japansk normaltid", "SRT": "surinamsk tid"},
	},
}

var selectedLocale = "nb_NO"

func init() {
	selectedLocale = config.SelectedLocale()
}

func FmtDateShort(t time.Time) string {
	loc := locales[selectedLocale]
	return fmt.Sprintf("%d. %s %d", t.Day(), loc.monthsAbbreviated[t.Month()], t.Year())
}

func FmtDateLong(t time.Time) string {

	loc := locales[selectedLocale]
	return fmt.Sprintf("%s %d. %s %d", loc.daysWide[t.Weekday()], t.Day(), loc.monthsWide[t.Month()], t.Year())
}

func Regex(p, t string) string {
	r, err := regexp.Compile(p)
	if err != nil {
		l.Warn("Could not compile regex in mapper")
		return ""
	}
	return r.FindString(t)
}

// Maps anything to a FormField
func (um *UploadMetadata) ToFieldMapper(f ToFormField, sm FieldsStringMap) {
	if f.Format == "" || f.Key == "" || f.VisualName == "" {
		if f.Debug {
			l.WithFields(
				map[string]interface{}{
					"key":            f.Key,
					"toFormField":    f,
					"uploadMetadata": um,
					"fieldStringMap": sm,
				}).Debug("Attempted to map {key} field, but a required field is empty")
		}
		return
	}
	s := ""
	templ, err := template.New("field").
		Funcs(template.FuncMap{
			"regex":     Regex,
			"dateLong":  FmtDateLong,
			"dateShort": FmtDateShort,
			"date": func(layout string, t time.Time) string {
				return t.Format(layout)
			},
		}).
		Parse(f.Format)
	vars := struct {
		U       UploadMetadata
		F       FieldsStringMap
		Subject Person
	}{
		U:       *um,
		F:       sm,
		Subject: Person{},
	}
	if len(um.Subject) > 0 {
		vars.Subject = um.Subject[0]
	}
	if err != nil {
		l.WithError(err).WithFields(map[string]interface{}{
			"field": f,
		}).Warn("Failed to parse the template for toFieldMap {field}")
		return
	} else {
		var b bytes.Buffer

		templ.Execute(&b, vars)

		s = strings.TrimSpace(b.String())
	}
	if s == "" {
		if f.Debug {
			l.WithFields(
				map[string]interface{}{
					"key":            f.Key,
					"toFormField":    f,
					"uploadMetadata": um,
					"fieldStringMap": sm,
				}).Debug("Attempted to map {key} field, but the result was empty")
		}
		return
	}
	ff := FormFields{
		Key:            f.Key,
		FieldId:        f.Key,
		TranslationKey: f.Key,
		VisualName:     f.VisualName,
		Value:          s,
		Required:       f.Required,
		DataType:       "String",
		ValidationRule: ValidationRule{},
	}
	if f.Debug {
		l.WithFields(
			map[string]interface{}{
				"key":            f.Key,
				"value":          s,
				"formField":      ff,
				"toFormField":    f,
				"uploadMetadata": um,
				"fieldStringMap": sm,
			}).Debug("Mapped the field {key} to {value}")
	}
	um.FormFields = append(um.FormFields, ff)

}

func (um *UploadMetadata) AnyMapper(field string, format string, sm FieldsStringMap) {
	if format == "" || field == "" {
		return
	}
	s := ""
	templ, err := template.New(field).
		Funcs(template.FuncMap{
			"regex":     Regex,
			"dateLong":  FmtDateLong,
			"dateShort": FmtDateShort,
			"date": func(layout string, t time.Time) string {
				return t.Format(layout)
			},
		}).
		Parse(format)
	vars := struct {
		U       UploadMetadata
		F       FieldsStringMap
		Subject Person
	}{
		U:       *um,
		F:       sm,
		Subject: Person{},
	}
	if len(um.Subject) > 0 {
		vars.Subject = um.Subject[0]
	}
	if err != nil {
		l.WithError(err).WithFields(map[string]interface{}{
			"field":  field,
			"format": format,
		}).Warn("Failed to parse the template for field {field}")
		return
	} else {
		var b bytes.Buffer

		templ.Execute(&b, vars)

		s = strings.TrimSpace(b.String())
	}
	if s == "" {
		return
	}
	switch field {
	case "userid":
		um.UserId = s
	case "parent.name":
		um.Parent.Name = s
	case "parent.description":
		um.Parent.Description = s
	case "description":
		um.Description = s
	case "displayname":
		um.DisplayName = s
	}
}

func (um *UploadMetadata) SetStringField(sm FieldsStringMap, f FieldMapType, field FormFields) (err error) {
	fieldName := strings.ToLower(f.ToField)
	if strings.Contains(fieldName, "subject.") && len(um.Subject) == 0 {
		return
	}

	s := strings.TrimSpace(field.Value)
	l := lf.WithFields(map[string]interface{}{
		"fieldMapType": f,
		"fieldKey":     field.Key,
		"fieldId":      field.FieldId,
		"VisualName":   field.VisualName,
	})
	format := f.Format
	if format != "" {
		templ, err := template.New(field.FieldId + field.Key).Parse(format)
		if err != nil {
			l.Warn("Failed to parse template for field")

		} else {
			var b bytes.Buffer
			templ.Execute(&b, struct {
				Um UploadMetadata
				FormFields
				F     FieldsStringMap
				Field UploadMetadataField
			}{
				Um:         *um,
				FormFields: field,
				F:          sm,
				Field:      f.ToField,
			})

			s = b.String()
		}
	}

	switch fieldName {
	case UMDescription:
		if f.CheckStringCondition(um.Description, s) == "" {
			um.Description = s
		}
	case UMFTags:
		if s == "" {
			break
		}
		for _, t := range um.Tags {
			if s == t {
				break
			}
		}

		um.Tags = append(um.Tags, s)
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
				lf.Warnf("Could not parse '%s' as float in field '%s'", s, f.ToField)
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
		if s == "" {
			return
		}
		if f.Layout != "" {
			d, err := time.Parse(f.Layout, s)
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
		err = errors.Errorf("Cannot map key '%s' to string with value '%+v'", f.ToField, s)
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

type FieldsStringMap map[string]string

func (um UploadMetadata) FormFieldsToStringMap() FieldsStringMap {
	m := map[string]string{}
	for _, f := range um.FormFields {
		val := strings.TrimSpace(f.Value)
		if f.Key != "" {
			m[f.Key] = val
			m[f.Key] = val
		}
		if f.FieldId != "" {
			m[f.FieldId] = val
			m[f.FieldId] = val
		}
		if f.TranslationKey != "" {
			m[f.TranslationKey] = val
			m[f.TranslationKey] = f.Value
		}
	}
	return m
}

func (fm FieldMap) Get(key string) (FieldMapType, bool) {
	v := fm[key]
	return v, v.ToField != ""
}

var (
	am   = AnyMapperConfig()
	fm   = ToFormFieldMap()
	fMap = FieldMapperConfig()
)

func (um *UploadMetadata) MapFormFields() (err error) {
	if info.IsDebugMode() {
		l.WithFields(map[string]interface{}{
			"anyMapper":   am,
			"toFormField": fm,
			"fieldMap":    fMap,
		}).Debug("Mapping uploadMetadata using these maps")
	}
	sm := um.FormFieldsToStringMap()
	for _, f := range um.FormFields {
		fKey, found := fMap.FindFromField(f)

		if !found {
			continue
		}
		err := um.SetStringField(sm, fKey, f)
		if err != nil {
			lf.WithError(err).WithFields(map[string]interface{}{
				"fKey":  fKey,
				"value": f.Value,
			}).Error("Failed in mapping]")
		}
	}
	for key, format := range am {
		um.AnyMapper(key, format, sm)
	}
	for _, f := range fm {
		um.ToFieldMapper(f, sm)
	}

	return
}

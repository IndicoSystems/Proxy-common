package metadata

import (
	"encoding/base64"
	"encoding/json"
	"github.com/indicosystems/proxy/logger"
	"github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"strings"
)

const (
	// The ID of the user the file belongs to.
	UserId = "userid"

	// The name of the container the file belongs to.
	ParentName = "parentname"

	// The mime type of the file.
	FileType = "filetype"

	// The name given to the file by the user.
	DisplayName = "displayname"

	// The Checksum of the file.
	Checksum = "checksum"

	Filename = "filename"

	extId       = "extid"
	extParentId = "extParentid"
	// Indicates whether the upload is verfied as completed.
	extUploaded = "extUploaded"

	// The data available to all, as submitted by the Client
	MUploadMetadata = "UploadMetadata"
	// DeferId - Internal use for deferred uploads
	DeferId = "deferId"
	// DeferredFiledId - Internal use for deferred uploads
	DeferFileId = "deferFileId"
)

type Upl struct {
	UploadMetadata
}

var l logrus.FieldLogger = logger.Get("metadata")

type Metadata map[string]string
type Mapper func(data Metadata) Metadata

func GetUploadMetadata(info tusd.FileInfo) UploadMetadata {
	data := Metadata(info.MetaData)
	return data.GetUploadMetadata()
}

// Returns uploadMetadata, which contains all the information about the file that a connector should need.
func (m *Metadata) GetUploadMetadata() UploadMetadata {
	var um UploadMetadata
	m.getNested(MUploadMetadata, &um)
	//um := m.GetUploadMetadata()
	if &um == nil {
		// TODO: Make sure this cannot happen
		l.Errorf("Expected upload-metadata not to be nil")
		panic("Expected upload-metadata not to be nil")
	}
	l.Warn("Parent is nil (GetUploadMetadata)")
	return um

}

// Returns the DeferId, which is used for Deferred uploads.
func (m *Metadata) GetDeferId() string {
	return m.getExact(DeferId)
}

// Returns the DeferFileId, which is used for Deferred uploads.
func (m *Metadata) GetDeferFileId() string {
	return m.getExact(DeferFileId)
}

// Returns the DeferFileId, which is used for Deferred uploads.
func (m *Metadata) GetFilename() string {
	return m.getExact(Filename)
}
func (m *Metadata) GetExtId() string {
	return m.getExact(extId)
}
func (m *Metadata) GetExtParentId() string {
	return m.getExact(extParentId)
}
func (m *Metadata) GetExtUploaded() bool {
	return m.getExact(extUploaded) == "true"
}

func (m *Metadata) SetExtId(d string) {
	m.set(extId, d)
	um := m.GetUploadMetadata()
	um.ExtId = d
	m.replaceUploadMetadata(um)
}
func (m *Metadata) SetExtUploaded() {
	m.set(extUploaded, "true")
}
func (m *Metadata) SetExtParentId(d string) {
	m.set(extParentId, d)
	um := m.GetUploadMetadata()
	um.Parent.Id = d
	m.replaceUploadMetadata(um)
}

func (m *Metadata) SetDeferId(d string) {
	m.set(DeferId, d)
}
func (m *Metadata) SetDeferFileId(d string) {
	m.set(DeferFileId, d)
}

func (m *Metadata) getExact(key string) string {
	if val, ok := (*m)[key]; ok {
		return strings.TrimSpace(val)
	}
	if val, ok := (*m)[strings.ToLower(key)]; ok {
		return strings.TrimSpace(val)
	}
	return ""
}

func (m *Metadata) unwrap(t interface{}, s string) {
	if t == nil {
		return
	}
	sJ, err := json.Marshal(t)
	if err != nil {
		l.Errorf("There was a problem Marshalling the '%s'-field", s)
		return
	}
	b := base64.StdEncoding.EncodeToString(sJ)
	m.set(s, b)
}

func (m *Metadata) replaceUploadMetadata(um UploadMetadata) {

	newM := Metadata{}
	newM.unwrap(um, MUploadMetadata)
	(*m)[MUploadMetadata] = newM[MUploadMetadata]
	m.set(MUploadMetadata, newM.getExact(MUploadMetadata))
}

func (m *Metadata) set(key, value string) {
	(*m)[key] = value
}

// Helper for getting nested objects
func (m *Metadata) getNested(key string, v interface{}) error {

	str := m.getExact(key)
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

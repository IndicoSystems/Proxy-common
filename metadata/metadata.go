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
	// Used by authenticators
	ClientId                 = "client-id"
	ReqId                    = "req-id"
	AsUserName               = "as-username"
	AsUserId                 = "as-user-id"
	AsUserActiveDirectorySid = "as-user-sid"

	// The name of the container the file belongs to.
	ParentName = "parentname"

	// The mime type of the file.
	FileType = "filetype"

	// The name given to the file by the user.
	DisplayName = "displayname"

	// The Checksum of the file.
	Checksum = "checksum"

	Filename = "filename"

	ExtId       = "extid"
	ExtParentId = "extParentid"
	// Indicates whether the upload is verfied as completed.
	ExtUploaded = "extUploaded"

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
func (m *Metadata) SaveToInfo(info *tusd.FileInfo) {
	info.MetaData = tusd.MetaData(*m)
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
	if um.Parent.Id == "" && um.Parent.Name == "" && um.Parent.BatchId == "" {
		l.Warn("Parent does not have an ID, BatchId nor a Name (GetUploadMetadata)")

	}
	return um

}
func (m *Metadata) GetRaw(k string) string {
	return m.getExact(k)
}
func (m *Metadata) GetRawMetadata() string {
	return m.getExact(MUploadMetadata)
}
func (m *Metadata) GetClientId() string {
	return m.getExact(ClientId)
}
func (m *Metadata) GetReqId() string {
	return m.getExact(ReqId)
}
func (m *Metadata) SetReqId(reqid string) {
	m.set(ReqId, reqid)
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
	return m.getExact(ExtId)
}
func (m *Metadata) GetExtParentId() string {
	return m.getExact(ExtParentId)
}
func (m *Metadata) GetExtUploaded() bool {
	return m.getExact(ExtUploaded) == "true"
}

func (m *Metadata) SetExtId(d string) {
	m.set(ExtId, d)
	um := m.GetUploadMetadata()
	um.ExtId = d
	m.ReplaceUploadMetadata(um)
}
func (m *Metadata) SetExtUploaded() *Metadata {
	m.set(ExtUploaded, "true")
	return m
}
func (m *Metadata) SetClientId(cid string) *Metadata {
	m.set(ClientId, cid)
	return m
}
func (m *Metadata) SetAsUserName(u string) *Metadata {
	m.set(AsUserName, u)
	return m
}
func (m *Metadata) SetAsActiveDirectoryUserSid(sid string) *Metadata {
	m.set(AsUserActiveDirectorySid, sid)
	return m
}
func (m *Metadata) SetAsUserId(id string) *Metadata {
	m.set(AsUserId, id)
	return m
}

func (m *Metadata) SetExtParentId(d string) {
	m.set(ExtParentId, d)
	um := m.GetUploadMetadata()
	um.Parent.Id = d
	m.ReplaceUploadMetadata(um)
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

func (m *Metadata) ReplaceUploadMetadata(um UploadMetadata) {

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

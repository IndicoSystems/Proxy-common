package metadata

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/indicosystems/proxy/info"
	"github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"strconv"
	"strings"
)

const (
	// The ID of the user the file belongs to.
	UserId = "userid"
	// Used by authenticators
	ClientId                 = "client-id"
	ConnectorWritten         = "connector-written"
	ReqId                    = "req-id"
	AsUserName               = "as-username"
	AsUserId                 = "as-user-id"
	AsUserActiveDirectorySid = "as-user-sid"

	// The name of the container the file belongs to.
	ParentName = "parentname"
	SSN        = "__ssn"

	// The mime type of the file.
	FileType = "filetype"

	// A queue-id in an installations queue-system. Only valid for some systems, like Indico Gateway.
	ServiceQueueId = "serviceQueueId"

	// Deprecated ConnectorConfig, if needed for the queue-handling
	ConnectorConfig = "connectorConfig"

	// The name given to the file by the user.
	DisplayName = "displayname"

	// The Checksum of the file.
	Checksum = "checksum"

	Filename     = "filename"
	ErrorMessage = "errormsg"

	ExtId       = "extid"
	ExtParentId = "extParentid"
	// Indicates whether the upload is verfied as completed.
	ExtUploaded = "extUploaded"

	// The data available to all, as submitted by the Client
	MUploadMetadata                   = "UploadMetadata"
	ClientMessages                    = "ClientMessages"
	CaseNumberIgnored InternalInfoStr = "CaseNumberIgnored"
)

type InternalInfoStr string

type Upl struct {
	UploadMetadata
}

var l logrus.FieldLogger = logrus.StandardLogger()

func AssignMetaLogger(logger logrus.FieldLogger) {
	l = logger
}

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
		if info.IsDebugMode() {
			l.WithFields(map[string]interface{}{
				"clientMediaId":   m.GetClientId(),
				"from-request-id": m.GetReqId(),
				"m":               m,
				"um":              um,
			}).Warn("Parent does not have an ID, BatchId nor a Name (GetUploadMetadata)")

		} else {

			l.WithFields(map[string]interface{}{
				"clientMediaId":   m.GetClientId(),
				"from-request-id": m.GetReqId(),
			}).Warn("Parent does not have an ID, BatchId nor a Name (GetUploadMetadata)")
		}

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
func (m *Metadata) SetServiceQueueId(qID string) {
	m.set(ServiceQueueId, qID)
}
func (m *Metadata) GetConnectorWritten() int64 {
	w := m.getExact(ConnectorWritten)
	if w == "" {
		return -1
	}
	p, err := strconv.ParseInt(w, 10, 64)
	if err != nil {
		l.Warn("Could not parse connectorWritten")
	}
	return p
}
func (m *Metadata) SetConnectorWritten(written int64) {
	m.set(ConnectorWritten, strconv.FormatInt(written, 10))
}
func (m *Metadata) GetServiceQueueId() string {
	return m.getExact(ServiceQueueId)
}

type ClientMessage struct {
	Kind    InternalInfoStr
	Message string
}

func formatClientString(s string) string {
	if s == "" {
		return ""
	}
	return strings.TrimSpace(s)
}

// Appends ClientMessages, but ensures uniqueness.
func (m *Metadata) AppendClientMessage(cMsg ClientMessage) {
	if cMsg.Kind == "" || cMsg.Message == "" {
		return
	}
	message := formatClientString(cMsg.Message)
	kind := InternalInfoStr(formatClientString(string(cMsg.Kind)))

	if message == "" {
		return
	}
	if kind == "" {
		return
	}
	existing := m.GetClientMessages()
	// Check for uniqueness
	for _, ex := range existing {
		if message == ex.Message && kind == ex.Kind {
			return
		}
	}
	existing = append(existing, ClientMessage{
		Kind:    kind,
		Message: message,
	})
	fmt.Println("appended", existing)
	b, err := json.Marshal(existing)
	if err != nil {
		l.WithError(err).Error("Could not marshal clientMessage")
		return
	}
	m.set(ClientMessages, string(b))
}
func (m *Metadata) GetClientMessages() (cMsgs []ClientMessage) {
	ex := m.getExact(ClientMessages)
	if ex == "" {
		return
	}
	err := json.Unmarshal([]byte(ex), &cMsgs)
	if err != nil {
		l.WithError(err).Error("Could not unmarshal clientMessage")
		return
	}

	return
}

// Deprecated
func (m *Metadata) SetConnectorConfig(cfg map[string]string) {
	m.unwrap(cfg, ConnectorConfig)
}

// Deprecated
func (m *Metadata) GetConnectorConfig() (map[string]string, error) {
	var v map[string]string
	err := m.getNested(ConnectorConfig, &v)
	if err != nil {
		l.WithError(err).Error("failed to get connectorConfig")
		return v, errors.New("failed to get connectorConfig")
	}
	return v, nil
}

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
func (m *Metadata) GetErrorMessage() string {
	return m.getExact(ErrorMessage)
}
func (m *Metadata) SetErrorMessage(msg string) *Metadata {
	m.set(ErrorMessage, msg)
	return m
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
func (m *Metadata) Apply(info *tusd.FileInfo) tusd.FileInfo {
	info.MetaData = tusd.MetaData(*m)
	return *info
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
	if um.ClientMediaId != "" {
		m.set(ClientId, um.ClientMediaId)
	}
	if um.ExtId != "" {
		m.set(ExtId, um.ExtId)
	}
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

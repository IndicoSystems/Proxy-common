// Common types to assist development of Indico Proxy-connectors

package common

import (
	"context"
	"database/sql"
	"github.com/indicosystems/proxy/info"
	"github.com/indicosystems/proxy/metadata"
	"github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"net/http"
	"time"
)

func IsDev() bool {
	return info.IsDev()
}

type S3Config struct {
	Address   string
	AccessKey string
	SecretKey string
	Region    string
	s3store.S3Store
}

type UploadCompleteStatus string

const (
	UploadConfirmedComplete UploadCompleteStatus = "Confirmed"
)

type BaseConfig struct {
	L logrus.FieldLogger
	P Persistence
	Q QueueStorer
}

type Persistence interface {
	Set(k string, v interface{}) error
	Get(k string, v interface{}) (found bool, err error)
	GetTusdInfo(id string) (*tusd.FileInfo, bool)
	// Should only be used to create the info
	// TODO: Change to CreateTusdInfo
	SetInfo(info tusd.FileInfo) error
	SetUploadOffset(_id string, offset int64) error
	// Can be used to mark an upload as complete with external information
	SetUploaded(info tusd.FileInfo) error
}

type DataStore interface {
	NewUpload(ctx context.Context, info tusd.FileInfo) (upload tusd.Upload, err error)

	GetUpload(ctx context.Context, id string) (upload tusd.Upload, err error)
	// Methods that extend the tusd.DataStore
	SetInfo(info tusd.FileInfo) error
	GetInfo(ctx context.Context, id string) (tusd.FileInfo, error)
	RegisterConnector(interface{}) DataStore
	GetQueue() QueueStorer
	AddToQueue(id, connectorId, actionType string, dueAt time.Time) error
}

type QueueItem struct {
	ID string
	// The connector that is responsible for this item.
	ConnectorId string
	// Information about the upload, like storage, metadata, size, offset etc.
	Info tusd.FileInfo
	// The kind of action that needs to be done, from the connector's perspective.
	ActionType string
	// How many attempts the queue-item has been trough.
	Attempts int
	// Any errors occured. Can be used to inform a sys-admin.
	Error string
	// The time of which the item is due.
	DueAt    time.Time
	UploadId string
}

type StoreCreator interface {
	CreateS3Store(S3Config) (DataStore, error)
}

type HealthReporter interface {
	// Should report 200 if ok
	GetHealth() (int, error)
}

type QueueHandler interface {
	HandleQueue(qi QueueItem) (completeQueueItem, completeUpload bool, err error)
	GetQueueHandlerId() string
}

// Will be called before the actual upload is created. (tusd.DataStore.NewUpload)
// Modifying data will be persisted.
// An error returned will stop the uplaod from being created, and the error-message is returned to the client
//
// Typical use-cases here are creating album/folder, some extra validation if needed
type NewUploadInitiator interface {
	InitiateNewUpload(data *metadata.Metadata) error
}

// Will be called after the actual upload is completed (tusd.DataStore.Upload.FinishUpload)
// Modifying data will be persisted.
type UploadCompleter interface {
	CompleteUpload(info tusd.FileInfo) (UploadResult, error)
}

type GetAllOptions struct {
	ID          string
	Limit       int
	ConnectorId string
	ActionType  string
	DueBefore   sql.NullTime
	DueAfter    sql.NullTime
	OnlyDue     bool
}

type QueueOptions struct {
	Interval           time.Duration
	PostponeBaseAmount time.Duration
}

type QueueStorer interface {
	Complete(id string) error
	MarkErr(qi QueueItem, err string, postpone bool) error
	Options() QueueOptions
	GetAll(o GetAllOptions) (qis []QueueItem, found bool, err error)
	AddToQueue(infoId, connectorId, actionType string, dueAt time.Time) error
	UpdateQueueItem(id string, dueAt sql.NullTime, attempts int, err string) error
}

type AuthenticationPayload struct {
	// Used to authenticate the request.
	ClientId,
	// Used to authenticate the request
	ApiKey,
	// Represents the user-name in the current backend.
	UserName,
	// Represents a user-id in the current backend
	UserId,
	// Represents an active directory id.
	UserSid string
}

type UploadResult struct {
	Confirmed        UploadCompleteStatus
	ExtId            string
	CaseId           string
	ExternalParentId string
	ClientId         string
}

func (a *AuthenticationPayload) SetAuthenticationPayloadOnMetadata(m *metadata.Metadata) *metadata.Metadata {
	if a.ClientId != "" {
		m.SetClientId(a.ClientId)
	}
	if a.UserName != "" {
		m.SetAsUserName(a.UserName)
	}
	if a.UserSid != "" {
		m.SetAsActiveDirectoryUserSid(a.UserSid)
	}
	if a.UserId != "" {
		m.SetAsUserId(a.UserId)
	}
	return m
}

func AddIds(l logrus.FieldLogger, info tusd.FileInfo) logrus.FieldLogger {
	data := metadata.Metadata(info.MetaData)
	return l.WithFields(map[string]interface{}{
		"ID":      info.ID,
		"storage": info.Storage,
		"reqId":   data.GetReqId(),
	})
}

// Only used in CLI-mode, for printing the url to the uploaded file on the connectors service.
type FileUrlPrinter interface {
	PrintFileUrl(info tusd.FileInfo) (string, error)
}

// Runs before any other requests to authenticate
type Authenticator interface {
	Authenticate(r *http.Request, a AuthenticationPayload) error
}

// Can be used by clients to validate that their input is valid before submitting data.
type Validator interface {
	Validate(r *http.Request, a AuthenticationPayload, v ValidatePayload) (ValidateResponse, error)
}
type SupportedValidation struct {
	UserID,
	UserName,
	Sid,
	ParentID,
	ParentName,
	CaseID,
	CaseName,
	GroupName,
	GroupID bool
}

type ValidateResponse struct {
	UserID     ValidateUserResponse
	UserName   ValidateUserResponse
	Sid        ValidateUserResponse
	ParentID   ValidateParentResponse
	ParentName ValidateParentResponse
	CaseID     ValidateCaseResponse
	CaseName   ValidateCaseResponse
	GroupName  ValidateGroupResponse
	GroupID    ValidateGroupResponse
}
type ValidateNullableResponse struct {
	UserID     *ValidateUserResponse   `json:",omitempty"`
	UserName   *ValidateUserResponse   `json:",omitempty"`
	Sid        *ValidateUserResponse   `json:",omitempty"`
	ParentID   *ValidateParentResponse `json:",omitempty"`
	ParentName *ValidateParentResponse `json:",omitempty"`
	CaseID     *ValidateCaseResponse   `json:",omitempty"`
	CaseName   *ValidateCaseResponse   `json:",omitempty"`
	GroupName  *ValidateGroupResponse  `json:",omitempty"`
	GroupID    *ValidateGroupResponse  `json:",omitempty"`
}
type ValidateCaseResponse struct {
	ID string
}

type ValidateParentResponse struct {
	ID string
}
type ValidateGroupResponse struct {
	ID   string
	Gid  string
	Name string
}

type ValidateUserResponse struct {
	ID       string
	UserName string
	AuthID   string
}

type ValidatePayload struct {
	UserId     string
	UserName   string
	Sid        string
	ParentId   string
	ParentName string
	CaseId     string
	CaseName   string
	GroupName  string
	GroupId    string
}

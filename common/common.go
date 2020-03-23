// Common types to assist development of Indico Proxy-connectors

package common

import (
	"context"
	"database/sql"
	"github.com/indicosystems/proxy/metadata"
	"github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"time"
)

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
	ID          string
	ConnectorId string
	Info        tusd.FileInfo
	ActionType  string
	Attempts    int
	Error       string
	DueAt       time.Time
}

type StoreCreator interface {
	CreateS3Store(S3Config) (DataStore, error)
}

type QueueHandler interface {
	HandleQueue(qi QueueItem) (complete bool, err error)
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

// Persistence can be used to store information across requests.
type Persistence interface {
	Set(k string, v interface{}) error
	Get(k string, v interface{}) (found bool, err error)
	GetTusdInfo(id string) (*tusd.FileInfo, bool)
	SetInfo(info tusd.FileInfo) error
}
type GetAllOptions struct {
	ID          string
	Limit       int
	Increase    bool
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
	MarkErr(qi QueueItem, err string) error
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

// Common types to assist development of Indico Proxy-connectors

package common

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/indicosystems/proxy-common/metadata"
	"github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type S3Config struct {
	Address   string
	AccessKey string
	SecretKey string
	Region    string
	Options   S3ConfigOptions
	s3store.S3Store
}

type S3ConfigOptions struct {
	CalculateSha bool
	VerifyMime   bool
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
	SetReceiverChecksum(id string, checkSum metadata.CheckSum) error
	SetTemporaryChecksum(id string, checkSum []byte) error
	GetTemporaryChecksum(id string) ([]byte, error)
	GetTusdInfo(id string) (*tusd.FileInfo, bool)
	GetTusdInfos(ids []string) ([]*tusd.FileInfo, error)
	// Should only be used to create the info
	// TODO: Change to CreateTusdInfo
	SetInfo(info tusd.FileInfo) error
	SetUploadOffset(_id string, offset int64) error
	// Can be used to mark an upload as complete with external information
	SetUploaded(info tusd.FileInfo) error
	SetConnectorProgress(_id string, written int64) error
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
	// How many attempts the svcQueue-item has been trough.
	Attempts int
	// Any errors occured. Can be used to inform a sys-admin.
	Error string
	// The time of which the item is due.
	DueAt               time.Time
	UploadId            string
	BackoffLimitReached bool
}

type StoreCreator interface {
	CreateS3Store(S3Config) (DataStore, error)
}

type HealthReporter interface {
	// Should report 200 if ok
	GetHealth() (int, error)
}

type QueueHandler interface {
	HandleQueue(qi QueueItem) QueueRunResult
	GetQueueHandlerId() string
}

type QueueRunResult struct {
	// Set to true to mark the svcQueue-item as complete
	CompleteQueueItem bool
	// Set to true to mark the upload as complete. Will also complete the svcQueue-item
	CompleteUpload bool
	// Set to true to make the svcQueue back off on this item, and require manual intervention.
	Backoff bool
	// Additional info for the current error
	Err string
}

// Will be called before the actual upload is created. (tusd.DataStore.NewUpload)
// Modifying data will be persisted.
// An error returned will stop the uplaod from being created, and the error-message is returned to the client
//
// Typical use-cases here are creating album/folder, some extra validation if needed
type NewUploadInitiator interface {
	InitiateNewUpload(ctx context.Context, data *metadata.Metadata) error
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
	MarkErr(qi QueueItem, err string, postpone bool, backoff bool) error
	Options() QueueOptions
	GetAll(o GetAllOptions) (qis []QueueItem, found bool, err error)
	AddToQueue(infoId, connectorId, actionType string, dueAt time.Time) error
	UpdateQueueItem(id string, dueAt sql.NullTime, attempts int, err string, backoff bool) error
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
	ClientMediaId    string `json:"ClientId"`
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

// Helpful function during development/mappers.
type DryRunner interface {
	DryRun(ctx context.Context, mdata metadata.Metadata) error
}

type MetadataWriter interface {
	OutputMetadata(ctx context.Context, info tusd.FileInfo) ([]byte, string, error)
}

// Can be used by clients to validate that their input is valid before submitting data.
type Validator interface {
	// Validate should normally not error, but instead return error in the key of each validation-response
	Validate(r *http.Request, a AuthenticationPayload, v ValidatePayload) (ValidateResponse, error)
}

type ConnectorFeatures struct {
	// The minimum chunk-size to allow. Setting this to -1 disallows chunks, and all files must be uploaded in a single chunk.
	MinChunkSize int64
	// The maximum chunk-size to allow.
	MaxChunkSize int64
}

// Can be used to issue what features should be enabled
type FeatureAnnouncer interface {
	AnnounceFeatures() ConnectorFeatures
}
type SupportedValidation struct {
	UserID     bool `json:",omitempty"`
	UserName   bool `json:",omitempty"`
	Sid        bool `json:",omitempty"`
	ParentID   bool `json:",omitempty"`
	ParentName bool `json:",omitempty"`
	CaseID     bool `json:",omitempty"`
	CaseName   bool `json:",omitempty"`
	GroupName  bool `json:",omitempty"`
	GroupID    bool `json:",omitempty"`
}

type ValidateResponse struct {
	UserID     ValidateUserResponse   `json:",omitempty"`
	UserName   ValidateUserResponse   `json:",omitempty"`
	Sid        ValidateUserResponse   `json:",omitempty"`
	ParentID   ValidateParentResponse `json:",omitempty"`
	ParentName ValidateParentResponse `json:",omitempty"`
	CaseID     ValidateCaseResponse   `json:",omitempty"`
	CaseName   ValidateCaseResponse   `json:",omitempty"`
	GroupName  ValidateGroupResponse  `json:",omitempty"`
	GroupID    ValidateGroupResponse  `json:",omitempty"`
}

type ValidationErrorResponse struct {
	LocalizedResponse `json:",omitempty"`
	Status            string `json:",omitempty"`
	StatusCode        int    `json:",omitempty"`
}

func NewValidationError(status string, statusCode int, response LocalizedResponse) ValidationErrorResponse {
	return ValidationErrorResponse{
		LocalizedResponse: response,
		StatusCode:        statusCode,
		Status:            status,
	}
}

type ValidateNullableResponse struct {
	Error      *ValidationErrorResponse `json:",omitempty"`
	UserID     *ValidateUserResponse    `json:",omitempty"`
	UserName   *ValidateUserResponse    `json:",omitempty"`
	Sid        *ValidateUserResponse    `json:",omitempty"`
	ParentID   *ValidateParentResponse  `json:",omitempty"`
	ParentName *ValidateParentResponse  `json:",omitempty"`
	CaseID     *ValidateCaseResponse    `json:",omitempty"`
	CaseName   *ValidateCaseResponse    `json:",omitempty"`
	GroupName  *ValidateGroupResponse   `json:",omitempty"`
	GroupID    *ValidateGroupResponse   `json:",omitempty"`
}

type ValidateCaseResponse struct {
	LocalizedResponse `json:",omitempty"`
	Error             *ValidationErrorResponse `json:",omitempty"`
	ID                string                   `json:",omitempty"`
	Name              string                   `json:",omitempty"`
	Private           bool                     `json:",omitempty"`
	Sensitive         bool                     `json:",omitempty"`
}

type LocalizedResponse struct {
	// For use as header, etc
	Title    string `json:",omitempty"`
	Subtitle string `json:",omitempty"`
	// For use as additional field-information
	Details []Localizable
}

type Localizable struct {
	// Raw key, for use with data
	RawKey string
	// Raw value, for use with data
	RawValue string
	// Localized key
	Key string
	// Localized value
	Value string
}

type ValidateParentResponse struct {
	Error *ValidationErrorResponse `json:",omitempty"`
	ID    string                   `json:",omitempty"`
}
type ValidateGroupResponse struct {
	Error *ValidationErrorResponse `json:",omitempty"`
	ID    string                   `json:",omitempty"`
	Gid   string                   `json:",omitempty"`
	Name  string                   `json:",omitempty"`
}

type ValidateUserResponse struct {
	Error    *ValidationErrorResponse `json:",omitempty"`
	ID       string                   `json:",omitempty"`
	UserName string                   `json:",omitempty"`
	AuthID   string                   `json:",omitempty"`
}

type ValidatePayload struct {
	UserID     string `json:",omitempty"`
	UserName   string `json:",omitempty"`
	Sid        string `json:",omitempty"`
	ParentID   string `json:",omitempty"`
	ParentName string `json:",omitempty"`
	CaseID     string `json:",omitempty"`
	CaseName   string `json:",omitempty"`
	GroupName  string `json:",omitempty"`
	GroupID    string `json:",omitempty"`
}

type SearchResultItem struct {
	ID             string `json:",omitempty"`
	DisplayName    string `json:",omitempty"`
	AltID          string `json:",omitempty"`
	AltDisplayName string `json:",omitempty"`
}

type HasMore string

const (
	HasMoreTrue    HasMore = "YES"
	HasMoreUnknown HasMore = ""
	HasMoreNo      HasMore = "NO"
)

type SearchAggregate struct {
	Offset  int
	HasMore HasMore            `json:",omitempty"`
	Items   []SearchResultItem `json:",omitempty"`
}

type SearchResult struct {
	User   *SearchAggregate `json:",omitempty"`
	Parent *SearchAggregate `json:",omitempty"`
	Case   *SearchAggregate `json:",omitempty"`
	Group  *SearchAggregate `json:",omitempty"`
}

type SearchInput struct {
	Query string
	Kind  string
	Limit int
}

type SearchOptions struct {
	RequiresAuthentication bool
}

type SearchHandler interface {
	Search(AuthenticationPayload, SearchInput) (SearchResult, error)
	SearchCacheKey(AuthenticationPayload, SearchInput) string
	SearchableFields() (*SupportedSearches, SearchOptions)
}
type SupportedSearches struct {
	UserID     bool `json:",omitempty"`
	UserName   bool `json:",omitempty"`
	Sid        bool `json:",omitempty"`
	ParentID   bool `json:",omitempty"`
	ParentName bool `json:",omitempty"`
	CaseID     bool `json:",omitempty"`
	CaseName   bool `json:",omitempty"`
	GroupName  bool `json:",omitempty"`
	GroupID    bool `json:",omitempty"`
}

func (s *SupportedSearches) MapEnabled() map[string]bool {
	a := map[string]bool{}
	if s.UserID {
		a["UserID"] = true
	}
	if s.UserName {
		a["UserName"] = true
	}
	if s.Sid {
		a["Sid"] = true
	}
	if s.ParentName {
		a["ParentName"] = true
	}
	if s.CaseID {
		a["CaseID"] = true
	}
	if s.CaseName {
		a["CaseName"] = true
	}
	if s.GroupName {
		a["GroupName"] = true
	}
	if s.GroupID {
		a["GroupID"] = true
	}
	return a
}

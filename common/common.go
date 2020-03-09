// Common types to assist development of Indico Proxy-connectors

package common

import (
	"context"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type S3Config struct {
	Address   string
	AccessKey string
	SecretKey string
	Region    string
	s3store.S3Store
}

type DataStore interface {
	NewUpload(ctx context.Context, info tusd.FileInfo) (upload tusd.Upload, err error)

	GetUpload(ctx context.Context, id string) (upload tusd.Upload, err error)
	// Methods that extend the tusd.DataStore
	SetInfo(info tusd.FileInfo) error
}

type StoreCreator interface {
	CreateS3Store(S3Config) (DataStore, error)
}

// Persistence can be used to store information across requests.
type Persistence interface {
	Set(k string, v interface{}) error
	Get(k string, v interface{}) (found bool, err error)
	GetTusdInfo(id string) (*tusd.FileInfo, bool)
	SetInfo(info tusd.FileInfo) error
}

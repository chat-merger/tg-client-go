package blobstore

import (
	"errors"
	"io"
)

type URI = string

var ErrBlobNotFound = errors.New("blob not found")

type TempBlobStore interface {
	Save(data io.Reader) (*URI, error)
	Get(uri URI) (io.Reader, error)
}

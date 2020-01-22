package repo

import (
	"github.com/sukhjit/lambda-mock-server/model"
)

// Document interface
type Document interface {
	Get(id string) (*model.Document, error)
	Add(document *model.Document) error
	Delete(id string) error
	Close() error
}

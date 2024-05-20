package common

import (
	"database/sql"

	"github.com/openfga/openfga/pkg/storage/sqlcommon"
)

type StorageAuthProvider interface {
	NewMySQL(uri string, cfg *sqlcommon.Config) (*sql.DB, error)
	NewPosgreSQL(uri string, cfg *sqlcommon.Config) (*sql.DB, error)
}

package password

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
	"github.com/openfga/openfga/pkg/storage/authprovider/common"
	"github.com/openfga/openfga/pkg/storage/sqlcommon"
)

type PasswordAuthProvider struct {
}

var _ common.StorageAuthProvider = (*PasswordAuthProvider)(nil)

// NewMySQL implements authprovider.StorageAuthProvider.
func (p *PasswordAuthProvider) NewMySQL(uri string, cfg *sqlcommon.Config) (*sql.DB, error) {
	if cfg.Username != "" || cfg.Password != "" {
		dsnCfg, err := mysql.ParseDSN(uri)
		if err != nil {
			return nil, fmt.Errorf("parse mysql connection dsn: %w", err)
		}

		if cfg.Username != "" {
			dsnCfg.User = cfg.Username
		}
		if cfg.Password != "" {
			dsnCfg.Passwd = cfg.Password
		}

		uri = dsnCfg.FormatDSN()
	}
	return sql.Open("mysql", uri)
}

// NewPosgreSQL implements authprovider.StorageAuthProvider.
func (p *PasswordAuthProvider) NewPosgreSQL(uri string, cfg *sqlcommon.Config) (*sql.DB, error) {
	if cfg.Username != "" || cfg.Password != "" {
		parsed, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("parse postgres connection uri: %w", err)
		}

		username := ""
		if cfg.Username != "" {
			username = cfg.Username
		} else if parsed.User != nil {
			username = parsed.User.Username()
		}

		switch {
		case cfg.Password != "":
			parsed.User = url.UserPassword(username, cfg.Password)
		case parsed.User != nil:
			if password, ok := parsed.User.Password(); ok {
				parsed.User = url.UserPassword(username, password)
			} else {
				parsed.User = url.User(username)
			}
		default:
			parsed.User = url.User(username)
		}

		uri = parsed.String()
	}

	return sql.Open("pgx", uri)
}

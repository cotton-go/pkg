package mysql

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/cotton-go/pkg/ssh"
	mysqld "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DriverName                    string
	ServerVersion                 string
	DSN                           string
	DSNConfig                     *mysqld.Config
	Conn                          gorm.ConnPool
	SkipInitializeWithVersion     bool
	DefaultStringSize             uint
	DefaultDatetimePrecision      *int
	DisableWithReturning          bool
	DisableDatetimePrecision      bool
	DontSupportRenameIndex        bool
	DontSupportRenameColumn       bool
	DontSupportForShareClause     bool
	DontSupportNullAsDefaultValue bool
	DontSupportRenameColumnUnique bool
	// As of MySQL 8.0.19, ALTER TABLE permits more general (and SQL standard) syntax
	// for dropping and altering existing constraints of any type.
	// see https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
	DontSupportDropConstraint bool
	SSH                       *ssh.Config
}

func New(conf Config) gorm.Dialector {
	if sshConf := conf.SSH; sshConf != nil {
		conn, err := ssh.Connect(*sshConf)
		if err != nil {
			return nil
		}

		key := fmt.Sprintf("%s-%d-%d-%s", sshConf.Host, sshConf.Port, sshConf.Type, sshConf.User)
		mysqld.RegisterDialContext(key, func(ctx context.Context, addr string) (net.Conn, error) {
			return conn.Client().DialContext(ctx, "tcp", addr)
		})

		conf.DSN = strings.Replace(conf.DSN, "@tcp(", fmt.Sprintf("@%s(", key), 1)
	}

	return mysql.New(conf.config())
}

func Open(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func (c Config) config() mysql.Config {
	return mysql.Config{
		DriverName:                    c.DriverName,
		ServerVersion:                 c.ServerVersion,
		DSN:                           c.DSN,
		DSNConfig:                     c.DSNConfig,
		Conn:                          c.Conn,
		SkipInitializeWithVersion:     c.SkipInitializeWithVersion,
		DefaultStringSize:             c.DefaultStringSize,
		DefaultDatetimePrecision:      c.DefaultDatetimePrecision,
		DisableWithReturning:          c.DisableWithReturning,
		DisableDatetimePrecision:      c.DisableDatetimePrecision,
		DontSupportRenameIndex:        c.DontSupportRenameIndex,
		DontSupportRenameColumn:       c.DontSupportRenameColumn,
		DontSupportForShareClause:     c.DontSupportForShareClause,
		DontSupportNullAsDefaultValue: c.DontSupportNullAsDefaultValue,
		DontSupportRenameColumnUnique: c.DontSupportRenameColumnUnique,
		DontSupportDropConstraint:     c.DontSupportDropConstraint,
	}
}

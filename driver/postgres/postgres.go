package postgres

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cotton-go/pkg/ssh"
)

type Config struct {
	DriverName           string
	DSN                  string
	WithoutQuotingCheck  bool
	PreferSimpleProtocol bool
	WithoutReturning     bool
	Conn                 gorm.ConnPool
	SSH                  *ssh.Config
}

func New(conf Config) gorm.Dialector {
	if sshConf := conf.SSH; sshConf != nil {
		conn, err := ssh.Connect(*sshConf)
		if err != nil {
			return nil
		}

		key := fmt.Sprintf("%s-%d-%d-%s", sshConf.Host, sshConf.Port, sshConf.Type, sshConf.User)
		sql.Register(key, NewDialector(conn))
		conf.DriverName = key
	}

	return postgres.New(conf.config())
}

func (c Config) config() postgres.Config {
	return postgres.Config{
		DriverName:           c.DriverName,
		DSN:                  c.DSN,
		WithoutQuotingCheck:  c.WithoutQuotingCheck,
		PreferSimpleProtocol: c.PreferSimpleProtocol,
		WithoutReturning:     c.WithoutReturning,
		Conn:                 c.Conn,
	}
}

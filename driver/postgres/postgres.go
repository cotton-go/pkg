package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cotton-go/pkg/ssh"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config 结构体用于配置 PostgreSQL 数据库的连接选项。
type Config struct {
	// DriverName 是用于连接 PostgreSQL 数据库的驱动程序的名称。
	// 例如 "postgres"。
	DriverName string

	// DSN 是用于连接 PostgreSQL 数据库的数据源名称。
	// 它包含了连接数据库所需的信息，如数据库地址、用户名和密码等。
	DSN string

	// WithoutQuotingCheck 表示是否禁用字段名称的引号检查。
	// 如果设置为 true，则不会对字段名称进行引号检查，否则会对字段名称进行引号检查。
	WithoutQuotingCheck bool

	// PreferSimpleProtocol 表示是否偏好使用简单协议。
	// 如果设置为 true，则偏好使用简单协议来连接 PostgreSQL 数据库，否则偏好使用复杂协议。
	PreferSimpleProtocol bool

	// WithoutReturning 表示是否禁用 RETURNING 子句。
	// 如果设置为 true，则不会在 SQL 语句中使用 RETURNING 子句，否则会使用 RETURNING 子句。
	WithoutReturning bool

	// Conn 是已经存在的数据库连接池。
	// 如果设置了该值，则会使用该连接池建立数据库连接。
	Conn gorm.ConnPool

	// SSH 是用于通过 SSH 连接 PostgreSQL 数据库的 SSH 配置。
	// 如果设置了该值，则会通过 SSH 隧道建立数据库连接。
	SSH *ssh.Config
}

// New 根据提供的配置创建一个新的 Gorm 数据库连接。
// 它支持通过 SSH 隧道进行连接，如果配置中提供了 SSH 信息。
//
// 参数:
//   - conf: 数据库和 SSH 配置。
//
// 返回值:
//   - gorm.Dialector: 用于 Gorm 以建立数据库连接的接口。
func New(conf Config) gorm.Dialector {
	// 检查是否提供了 SSH 配置，如果提供了，则尝试通过 SSH 进行连接。
	if sshConf := conf.SSH; sshConf != nil {
		// 使用 SSH 配置尝试建立连接。
		conn, err := ssh.Connect(*sshConf)
		if err != nil {
			// 如果连接失败，返回 nil，表示无法建立数据库连接。
			return nil
		}

		// 根据 SSH 配置生成一个唯一的键，用于注册新的 SQL 驱动名。
		key := fmt.Sprintf("%s-%d-%d-%s", sshConf.Host, sshConf.Port, sshConf.Type, sshConf.User)
		// 使用建立的 SSH 连接注册新的 SQL 驱动。
		sql.Register(key, NewDialector(conn))
		// 更新配置中的驱动名，以便 Gorm 可以使用通过 SSH 建立的连接。
		conf.DriverName = key
	}

	// 使用更新后的配置创建并返回一个新的 PostgreSQL 数据库连接。
	return postgres.New(conf.config())
}

// config 将当前配置对象转换为 postgres.Config 类型的配置。
// 这个方法主要用于统一配置的获取方式，便于在不同地方使用相同的配置数据。
// 它通过将当前 Config 结构体的字段值赋给 postgres.Config 结构体，实现配置的适配。
//
// 返回值:
//   - 返回一个 postgres.Config 类型的配置对象，其中包含了当前配置的所有必要信息。
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

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

// Config 定义 MySQL 数据库驱动程序的配置选项。
type Config struct {
	DriverName                    string         `json:"driverName"`                    // 驱动名称，默认为 "mysql"。
	ServerVersion                 string         `json:"serverVersion,omitempty"`       // 服务器版本，默认为空。
	DSN                           string         `json:"dsn,omitempty"`                 // DSN 字符串，默认为空。
	DSNConfig                     *mysqld.Config `json:"-"`                             // DSN 配置
	Conn                          gorm.ConnPool  `json:"connPool,omitempty"`            // 连接池，默认为空。
	SkipInitializeWithVersion     bool           `json:"skipInitializeWithVersion"`     // 跳过与版本号相关的初始化，默认为 false。
	DefaultStringSize             uint           `json:"defaultStringSize"`             // 字符串的默认长度，默认为 0。
	DefaultDatetimePrecision      *int           `json:"defaultDatetimePrecision"`      // 日期时间的默认精度，默认为 nil。
	DisableWithReturning          bool           `json:"disableWithReturning"`          // 禁用 WITH RETURNING，默认为 false。
	DisableDatetimePrecision      bool           `json:"disableDatetimePrecision"`      // 禁用日期时间精度，默认为 false。
	DontSupportRenameIndex        bool           `json:"dontSupportRenameIndex"`        // 不支持重命名索引，默认为 false。
	DontSupportRenameColumn       bool           `json:"dontSupportRenameColumn"`       // 不支持重命名列，默认为 false。
	DontSupportForShareClause     bool           `json:"dontSupportForShareClause"`     // 不支持 FOR SHARE 子句，默认为 false。
	DontSupportNullAsDefaultValue bool           `json:"dontSupportNullAsDefaultValue"` // 不支持 NULL 作为默认值，默认为 false。
	DontSupportRenameColumnUnique bool           `json:"dontSupportRenameColumnUnique"` // 不支持重命名列的唯一性约束，默认为 false。
	// As of MySQL 8.0.19, ALTER TABLE permits more general (and SQL standard) syntax
	// for dropping and altering existing constraints of any type.
	// see https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
	DontSupportDropConstraint bool        `json:"dontSupportDropConstraint"` // 不支持使用 DROP CONSTRAINT 语法来删除约束，默认为 false。
	SSH                       *ssh.Config `json:"ssh,omitempty"`             // SSH 配置选项，默认为 nil。
}

// New 根据配置创建一个新的 Gorm 数据库连接。
// 它支持通过 SSH 隧道进行数据库连接，如果配置中提供了 SSH 配置。
//
// 参数:
//   - conf: 数据库和 SSH 连接的配置。
//
// 返回值:
//   - gorm.Dialector: 用于 Gorm 以建立数据库连接的接口。
func New(conf Config) gorm.Dialector {
	// 检查是否提供了 SSH 配置，如果提供了，则尝试建立 SSH 连接。
	if sshConf := conf.SSH; sshConf != nil {
		// 使用提供的 SSH 配置建立连接。
		conn, err := ssh.Connect(*sshConf)
		if err != nil {
			// 如果连接失败，抛出异常。
			panic(err)
		}

		// 生成一个唯一的键，用于标识这个 SSH 连接。
		key := fmt.Sprintf("%s-%d-%d-%s", sshConf.Host, sshConf.Port, sshConf.Type, sshConf.User)
		// 注册一个自定义的拨号函数，使用 SSH 连接来拨号。
		mysqld.RegisterDialContext(key, func(ctx context.Context, addr string) (net.Conn, error) {
			return conn.Client().DialContext(ctx, "tcp", addr)
		})

		// 修改 DSN，将 SSH 隧道的标识符替换进去。
		conf.DSN = strings.Replace(conf.DSN, "@tcp(", fmt.Sprintf("@%s(", key), 1)
	}

	// 最终，基于配置创建并返回 MySQL 数据库连接器。
	return mysql.New(conf.config())
}

// Open 根据给定的 DSN (数据源名称) 打开一个 MySQL 数据库连接。
// 它返回一个实现了 gorm.Dialector 接口的数据库连接对象，用于后续的数据库操作。
// 该函数实际上调用了 mysql 包中的 Open 函数来创建数据库连接。
//
// 参数:
//   - dsn: 数据源名称 (Data Source Name) 的字符串，格式通常包括用户名、密码、主机、端口、数据库名等信息。
//
// 返回值:
//   - gorm.Dialector: 一个可以与 GORM 框架配合使用的数据库连接对象。
func Open(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

// config 将自定义配置对象转换为 mysql 驱动的 Config 结构体。
func (c Config) config() mysql.Config {
	return mysql.Config{
		DriverName:                    c.DriverName,                    // 驱动名称，例如 "mysql"。
		ServerVersion:                 c.ServerVersion,                 // 服务器版本，用于与特定版本的 MySQL 兼容。
		DSN:                           c.DSN,                           // 数据源名称，包含连接数据库所需的信息，如数据库地址、用户名和密码等。
		DSNConfig:                     c.DSNConfig,                     // DSN 的解析配置。
		Conn:                          c.Conn,                          // 已有的数据库连接，如果有的话。
		SkipInitializeWithVersion:     c.SkipInitializeWithVersion,     // 是否跳过使用版本进行初始化。
		DefaultStringSize:             c.DefaultStringSize,             // 默认字符串字段的大小。
		DefaultDatetimePrecision:      c.DefaultDatetimePrecision,      // 默认日期时间精度。
		DisableWithReturning:          c.DisableWithReturning,          // 是否禁用 WITH RETURNING 语法。
		DisableDatetimePrecision:      c.DisableDatetimePrecision,      // 是否禁用日期时间精度。
		DontSupportRenameIndex:        c.DontSupportRenameIndex,        // 是否不支持重命名索引。
		DontSupportRenameColumn:       c.DontSupportRenameColumn,       // 是否不支持重命名列。
		DontSupportForShareClause:     c.DontSupportForShareClause,     // 是否不支持 FOR SHARE 子句。
		DontSupportNullAsDefaultValue: c.DontSupportNullAsDefaultValue, // 是否不支持将 NULL 作为默认值。
		DontSupportRenameColumnUnique: c.DontSupportRenameColumnUnique, // 是否不支持重命名唯一约束的列。
		DontSupportDropConstraint:     c.DontSupportDropConstraint,     // 是否不支持 DROP CONSTRAINT 操作。
	}
}

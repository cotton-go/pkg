package postgres

import (
	"database/sql/driver"
	"net"
	"time"

	"github.com/cotton-go/pkg/ssh"
	"github.com/lib/pq"
)

// Dialector 结构体定义了一个 SSH 客户端连接。
// 它用于后续的数据库操作，通过 SSH 隧道进行。
type Dialector struct {
	conn *ssh.Client // conn 字段存储了一个指向 ssh.Client 的指针，
}

// NewDialector 创建一个新的 Dialector 实例。
//
// 参数:
//   - conn: 已建立的 SSH 连接客户端。
//
// 返回值:
//   - *Dialector 类型的指针，用于后续的 SSH 操作。
func NewDialector(conn *ssh.Client) *Dialector {
	return &Dialector{conn}
}

// Open 打开一个数据库连接。
// 它是数据库方言的一部分，用于特定于方言的数据库连接。
// 该方法使用 pq 包的 DialOpen 函数来实现实际的连接操作。
//
// 参数:
//   - dsn: 数据库连接字符串，包含连接数据库所需的信息。
//
// 返回值:
//   - driver.Conn: 建立的数据库连接对象，实现了 driver.Conn 接口。
//   - error: 如果连接过程中出现错误，则返回该错误。
func (d *Dialector) Open(dsn string) (_ driver.Conn, err error) {
	return pq.DialOpen(d, dsn)
}

// Dial 建立到指定网络地址的连接。
// 该方法利用 d.conn 的 Client 方法获取的客户端进行连接操作。
//
// 参数:
//   - network: 网络类型，例如 "tcp"、"udp" 等。
//   - address: 要连接的网络地址。
//
// 返回值:
//   - net.Conn: 建立的网络连接。
//   - error: 如果连接失败，则返回错误信息。
func (d *Dialector) Dial(network, address string) (net.Conn, error) {
	return d.conn.Client().Dial(network, address)
}

// DialTimeout 在指定超时时间内，通过特定的网络和地址进行连接。
// 该方法实际上是调用内部连接的 Client 方法进行连接操作。
// 注意：尽管有超时参数，但实际连接操作并未使用此超时参数。
//
// 参数：
//   - network: 网络类型，例如 "tcp"。
//   - address: 要连接的地址，例如 "localhost:8080"。
//   - timeout: 连接超时时间。
//
// 返回值：
//   - net.Conn: 建立的连接对象。
//   - error: 如果连接失败，会返回一个错误。
func (d *Dialector) DialTimeout(network, address string, _ time.Duration) (net.Conn, error) {
	return d.conn.Client().Dial(network, address)
}

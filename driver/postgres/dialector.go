package postgres

import (
	"database/sql/driver"
	"net"
	"time"

	"github.com/lib/pq"

	"github.com/cotton-go/pkg/ssh"
)

type Dialector struct {
	conn *ssh.Client
}

func NewDialector(conn *ssh.Client) *Dialector {
	return &Dialector{conn}
}

func (d *Dialector) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(d, s)
}

func (d *Dialector) Dial(network, address string) (net.Conn, error) {
	return d.conn.Client().Dial(network, address)
}

func (d *Dialector) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return d.conn.Client().Dial(network, address)
}

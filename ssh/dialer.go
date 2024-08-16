package ssh

import (
	"database/sql/driver"
	"net"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

type Dialer struct {
	Client *ssh.Client
}

func (self *Dialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *Dialer) Dial(network, address string) (net.Conn, error) {
	return self.Client.Dial(network, address)
}

func (self *Dialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.Client.Dial(network, address)
}

package ssh

import (
	"golang.org/x/crypto/ssh"
)

type Client struct {
	conn *ssh.Client
}

func (c Client) Client() *ssh.Client {
	return c.conn
}

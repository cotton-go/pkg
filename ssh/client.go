package ssh

import (
	"golang.org/x/crypto/ssh"
)

// Client 结构体代表一个SSH客户端连接，它包含了一个指向ssh.Client的指针。
type Client struct {
	conn *ssh.Client // 指向ssh.Client的指针，表示SSH客户端连接
}

// Client 方法返回当前客户端的 SSH 连接。
//
// 返回值是
//   - *ssh.Client 类型，即一个指向 ssh.Client 实例的指针。
func (c Client) Client() *ssh.Client {
	return c.conn
}

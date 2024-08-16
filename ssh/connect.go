package ssh

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// ConfigType ssh 连接支持的类型
type ConfigType uint8

const (
	// ConfigTypeByPassword 表示通过密码进行配置
	ConfigTypeByPassword ConfigType = iota
	// ConfigTypeByPrivateKey 表示通过私钥进行配置
	ConfigTypeByPrivateKey
	// ConfigTypeByPrivateKeyPath 表示通过私钥路径进行配置
	ConfigTypeByPrivateKeyPath
)

// Config SSH连接配置信息结构体
type Config struct {
	// Host SSH远程主机的IP地址或域名
	Host string
	// Port SSH远程主机的连接端口
	Port int
	// Type SSH连接的认证类型，包括密码、私钥内容或私钥文件路径三种方式
	Type ConfigType
	// User SSH远程主机的登录用户名
	User string
	// Password SSH远程主机的登录密码，仅在Type为ConfigTypeByPassword时生效
	Password string
	// PrivateKey SSH远程主机的私钥内容，仅在Type为ConfigTypeByPrivateKey时生效
	PrivateKey string
	// PrivateKeyPath SSH远程主机的私钥文件路径，仅在Type为ConfigTypeByPrivateKeyPath时生效
	PrivateKeyPath string
}

// Connect 根据提供的配置信息创建一个SSH客户端连接。
// 它支持通过密码、私钥内容或私钥文件路径三种方式来创建SSH连接。
//
// 参数
//   - config: 包含了SSH连接所需的配置信息，包括用户类型、密码、私钥等。
//
// 返回值:
//   - *ssh.ClientConfig 类型的SSH客户端配置指针，以及可能的错误信息。
func Connect(conf Config) (*Client, error) {
	// 如果未指定端口号，则使用默认的SSH端口
	if conf.Port == 0 {
		conf.Port = 22
	}

	// 根据配置信息创建SSH客户端配置
	clientConfig, err := NewSSHConfig(conf)
	if err != nil {
		// 如果创建客户端配置失败，返回错误
		return nil, err
	}

	// 使用TCP协议拨号连接到SSH服务器
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), clientConfig)
	if err != nil {
		// 如果连接失败，包装原始错误并返回
		return nil, errors.Wrap(err, "failed to connect to SSH server")
	}

	// 连接成功，返回SSH客户端实例
	return &Client{conn}, nil
}

// NewSSHConfig 根据提供的配置生成SSH客户端配置。
// 它支持通过密码、私钥内容或私钥文件路径三种方式来创建SSH连接。
//
// 参数
//   - config: 包含了SSH连接所需的配置信息，包括用户类型、密码、私钥等。
//
// 返回值:
//   - *ssh.ClientConfig 类型的SSH客户端配置指针，以及可能的错误信息。
func NewSSHConfig(config Config) (*ssh.ClientConfig, error) {
	// 初始化认证方法列表
	var authMethods []ssh.AuthMethod

	// 根据配置的认证类型设置认证方法
	switch config.Type {
	case ConfigTypeByPassword:
		// 如果配置类型为密码且密码不为空，则使用密码作为认证方法
		if config.Password != "" {
			authMethods = append(authMethods, ssh.Password(config.Password))
		} else {
			// 如果密码为空，返回错误
			return nil, errors.New("password is empty")
		}
	case ConfigTypeByPrivateKey:
		// 如果配置类型为私钥内容，解析私钥并添加到认证方法列表
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse private key")
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	case ConfigTypeByPrivateKeyPath:
		// 如果配置类型为私钥文件路径，读取私钥文件并解析，然后添加到认证方法列表
		key, err := os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read private key")
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse private key")
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	default:
		// 如果配置的认证类型不明确，默认尝试使用密码作为认证方法
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	// 构建并返回SSH客户端配置
	return &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

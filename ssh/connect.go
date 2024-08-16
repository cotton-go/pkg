package ssh

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type ConfigType uint8

const (
	ConfigTypeByPassword ConfigType = iota
	ConfigTypeByPrivateKey
	ConfigTypeByPrivateKeyPath
)

type Config struct {
	Host           string
	Port           int
	Type           ConfigType
	User           string
	Password       string
	PrivateKey     string
	PrivateKeyPath string
}

func Connect(conf Config) (*Client, error) {
	clientConfig, err := NewSSHConfig(conf)
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to SSH server")
	}

	return &Client{conn}, nil
}

func NewSSHConfig(config Config) (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod
	switch config.Type {
	case ConfigTypeByPassword:
		authMethods = append(authMethods, ssh.Password(config.Password))
	case ConfigTypeByPrivateKey:
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse private key")
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	case ConfigTypeByPrivateKeyPath:
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
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	return &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

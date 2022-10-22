package ssh

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Host     string
	Port     string
	User     string
	Password string
	cfg      *ssh.ClientConfig
}

func (cli *SshClient) init() {
	// TODO:
	// 1. think of a more secure way instead of ignoring the host key
	// 2. implement ssh key auth

	cli.cfg = &ssh.ClientConfig{
		User:            cli.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(cli.Password),
			// ssh.PublicKeys(signer),
		},
	}

	if cli.Port == "" {
		cli.Port = "22"
	}
}

func (cli *SshClient) Run(command string) ([]byte, error) {
	var conn *ssh.Client
	var err error
	var result []byte
	var session *ssh.Session
	var buff bytes.Buffer

	cli.init()

	conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", cli.Host, cli.Port), cli.cfg)
	if err != nil {
		return result, err
	}
	defer conn.Close()

	session, err = conn.NewSession()
	if err != nil {
		return result, err
	}
	defer session.Close()

	session.Stdout = &buff
	if err := session.Run(command); err != nil {
		return result, err
	}
	return buff.Bytes(), err
}

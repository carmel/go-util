// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package ssh

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHClient represents Goph client.
type SSHClient struct {
	Addr   string
	Config *ssh.ClientConfig
	*ssh.Client
}

// NewConn returns new client and error if any.
func NewSSHClient(server, port string, config *ssh.ClientConfig) (c *SSHClient) {
	// c.Client, err = Dial("tcp", addr, port)
	return &SSHClient{
		Addr:   net.JoinHostPort(server, port),
		Config: config,
	}
}

func (c *SSHClient) init() error {
	if c.Client == nil {
		var err error
		c.Client, err = ssh.Dial("tcp", c.Addr, c.Config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run starts a new SSH session and runs the cmd, it returns CombinedOutput and err if any.
func (c SSHClient) Run(cmd string) ([]byte, error) {

	var err = c.init()
	if err != nil {
		return nil, err
	}

	var sess *ssh.Session
	if sess, err = c.NewSession(); err != nil {
		return nil, err
	}

	defer sess.Close()

	return sess.CombinedOutput(cmd)
}

// NewSftp returns new sftp client and error if any.
func (c SSHClient) NewSftp(opts ...sftp.ClientOption) (*sftp.Client, error) {
	c.init()
	return sftp.NewClient(c.Client, opts...)
}

// Upload a local file to remote server!
func (c SSHClient) Upload(localPath string, remotePath string) (err error) {

	local, err := os.Open(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Create(remotePath + string(os.PathSeparator) + filepath.Base(localPath))
	if err != nil {
		return
	}
	defer remote.Close()

	_, err = io.Copy(remote, local)
	return
}

// Download file from remote server!
func (c SSHClient) Download(remotePath string, localPath string) error {
	var filePath = localPath + string(os.PathSeparator) + filepath.Base(remotePath)
	_, err := os.Stat(filePath)
	if err == nil {
		return errors.New("file exists.")
	}
	if os.IsNotExist(err) {
		local, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer local.Close()

		ftp, err := c.NewSftp()
		if err != nil {
			return err
		}
		defer ftp.Close()

		remote, err := ftp.Open(remotePath)
		if err != nil {
			return err
		}
		defer remote.Close()

		if _, err = io.Copy(local, remote); err != nil {
			return err
		}

		return local.Sync()
	}
	return err
}

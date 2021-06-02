// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package ssh

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"path"
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
	if err := c.init(); err != nil {
		return nil, err
	}
	return sftp.NewClient(c.Client, opts...)
}

// Upload a local file to remote server!
func (c SSHClient) UploadFile(sftpCli *sftp.Client, localPath string, remoteDir string) (err error) {

	var local *os.File
	local, err = os.Open(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	var remote *sftp.File
	remote, err = sftpCli.Create(path.Join(remoteDir, filepath.Base(localPath)))
	if err != nil {
		return
	}
	defer remote.Close()

	_, err = io.Copy(remote, local)
	return
}

// Upload a local file or directory to remote server!
func (c SSHClient) Upload(localPath string, remoteDir string) (err error) {

	//获取路径的属性
	s, err := os.Stat(localPath)
	if err != nil {
		return
	}

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	// 判断是否是文件夹
	if s.IsDir() {
		var localFiles []fs.FileInfo
		localFiles, err = ioutil.ReadDir(localPath)
		if err != nil {
			return
		}
		// 先创建最外层文件夹
		remoteDir = path.Join(remoteDir, s.Name())
		err = ftp.Mkdir(remoteDir)
		if err != nil {
			return
		}
		// 遍历文件夹内容
		for _, backupDir := range localFiles {
			// 判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
			err = c.Upload(path.Join(localPath, backupDir.Name()), remoteDir)
			if err != nil {
				return
			}
		}
	} else {
		err = c.UploadFile(ftp, localPath, remoteDir)
	}
	return
}

// Download file from remote server!
func (c SSHClient) Download(remotePath string, localPath string) error {
	var filePath = path.Join(localPath, filepath.Base(remotePath))
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

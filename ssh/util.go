package ssh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

// DefaultKnownHostsPath returns default user knows hosts file.
func DefaultKnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.ssh/known_hosts", home), err
}

// DefaultKnownHosts returns host key callback from default known hosts path, and error if any.
func DefaultKnownHosts() (ssh.HostKeyCallback, error) {
	path, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, err
	}

	return knownhosts.New(path)
}

// CheckKnownHost checks is host in known hosts file.
// it returns is the host found in known_hosts file and error, if the host found in
// known_hosts file and error not nil that means public key mismatch, maybe MAN IN THE MIDDLE ATTACK! you should not handshake.
func CheckKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownFile string) (found bool, err error) {

	var keyErr *knownhosts.KeyError

	// Fallback to default known_hosts file
	if knownFile == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return false, err
		}

		knownFile = path
	}

	// Get host key callback
	callback, err := knownhosts.New(knownFile)

	if err != nil {
		return false, err
	}

	// check if host already exists.
	err = callback(host, remote, key)

	// Known host already exists.
	if err == nil {
		return true, nil
	}

	// Make sure that the error returned from the callback is host not in file error.
	// If keyErr.Want is greater than 0 length, that means host is in file with different key.
	if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		return true, keyErr
	}

	// Some other error occurred and safest way to handle is to pass it back to user.
	if err != nil {
		return false, err
	}

	// Key is not trusted because it is not in the file.
	return false, nil
}

// AddKnownHost add a a host to known hosts file.
func AddKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownFile string) (err error) {

	// Fallback to default known_hosts file
	if knownFile == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return err
		}

		knownFile = path
	}

	f, err := os.OpenFile(knownFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	remoteNormalized := knownhosts.Normalize(remote.String())
	hostNormalized := knownhosts.Normalize(host)
	addresses := []string{remoteNormalized}

	if hostNormalized != remoteNormalized {
		addresses = append(addresses, hostNormalized)
	}

	_, err = f.WriteString(knownhosts.Line(addresses, key) + "\n")

	return err
}

// PrivateKey Loads a private and public key from "path" and returns a SSH ClientConfig to authenticate with the server
func PrivateKey(username string, path string, timeout int) (*ssh.ClientConfig, error) {
	privateKey, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)

	if err != nil {
		return nil, err
	}

	callback, err := DefaultKnownHosts()
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User:    username,
		Timeout: time.Duration(timeout) * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: callback,
	}, nil
}

// Creates the configuration for a client that authenticates with a password protected private key
func PrivateKeyWithPassphrase(username string, timeout int, passpharase []byte, path string) (*ssh.ClientConfig, error) {
	privateKey, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, passpharase)

	if err != nil {
		return nil, err
	}

	callback, err := DefaultKnownHosts()
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User:    username,
		Timeout: time.Duration(timeout) * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: callback,
	}, nil
}

// Creates a configuration for a client that fetches public-private key from the SSH agent for authentication
func SSHAgent(username string, timeout int, keyCallBack ssh.HostKeyCallback) (*ssh.ClientConfig, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return &ssh.ClientConfig{
			Timeout: time.Duration(timeout) * time.Second,
		}, err
	}

	return &ssh.ClientConfig{
		User:    username,
		Timeout: time.Duration(timeout) * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.NewClient(conn).Signers),
		},
		HostKeyCallback: keyCallBack,
	}, nil
}

// Creates a configuration for a client that authenticates using username and password
func PasswordKey(username string, password string, timeout int) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:    username,
		Timeout: time.Duration(timeout) * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

type WriteCounter struct {
	Total int64 // Total # of bytes written
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	fmt.Fprintf(os.Stdout, "%.2f kb transferred\n", float64(wc.Total/1024))
	return n, nil
}

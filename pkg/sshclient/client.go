package sshclient

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
)

type Storage interface {
	Upload(dstPath string, data []byte) error
	Download(srcPath string) ([]byte, error)
	List(dirPath string) ([]DirEntry, error)
}

type Config struct {
	Host    string
	User    string
	KeyPath string
}

type Client struct {
	cfg  Config
	scon *ssh.Client
	sftp *sftp.Client
}

func New(cfg Config) (*Client, error) {
	key, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse key: %w", err)
	}

	sshConf := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConf.SetDefaults()

	conn, err := ssh.Dial("tcp", cfg.Host, sshConf)
	if err != nil {
		return nil, fmt.Errorf("ssh dial: %w", err)
	}

	sftpCli, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("sftp init: %w", err)
	}

	return &Client{cfg: cfg, scon: conn, sftp: sftpCli}, nil
}

func (c *Client) Close() error {
	_ = c.sftp.Close()
	return c.scon.Close()
}

func (c *Client) Upload(dstPath string, data []byte) error {
	dir := path.Dir(dstPath)
	if err := mkdirAll(c.sftp, dir); err != nil {
		return err
	}

	f, err := c.sftp.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *Client) Download(srcPath string) ([]byte, error) {
	f, err := c.sftp.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	return io.ReadAll(f)
}

type DirEntry struct {
	Name  string
	IsDir bool
}

func (c *Client) List(dirPath string) ([]DirEntry, error) {
	entries, err := c.sftp.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	out := make([]DirEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, DirEntry{
			Name:  e.Name(),
			IsDir: e.IsDir(),
		})
	}

	return out, nil
}

func mkdirAll(cli *sftp.Client, dir string) error {
	parent := "/"
	for _, elem := range splitPath(dir) {
		parent = path.Join(parent, elem)
		err := cli.Mkdir(parent)
		if err != nil && !os.IsExist(err) {
			if st, ok := err.(*sftp.StatusError); !ok || st.Code != 4 {
				return err
			}
		}
	}

	return nil
}

func splitPath(p string) []string {
	var parts []string
	for {
		dir, file := path.Split(p)
		if file != "" {
			parts = append([]string{file}, parts...)
		}

		if dir == "" || dir == "/" {
			break
		}
		p = path.Clean(dir)
	}
	return parts
}

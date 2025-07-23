package sshclient

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skip SSH integration test in short mode")
	}

	cfg := Config{
		Host:    "localhost:2222",
		User:    "testuser",
		KeyPath: "./keys/id_rsa_test",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	data := []byte("hello test")
	tmpFile := "/tmp/test_ssh_" + time.Now().Format("150405")

	if err := client.Upload(tmpFile, data); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	got, err := client.Download(tmpFile)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	if string(got) != string(data) {
		t.Fatalf("downloaded data mismatch: got %q", got)
	}

	_ = client.sftp.Remove(tmpFile)
}

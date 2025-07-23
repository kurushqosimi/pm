package archiver

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func TarGz(filePaths []string, baseDir string) ([]byte, error) {
	sort.Strings(filePaths)

	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	for _, rel := range filePaths {
		abs := filepath.Join(baseDir, rel)

		fi, err := os.Lstat(abs)
		if err != nil {
			return nil, err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return nil, err
		}

		hdr.Name = filepath.ToSlash(rel)

		hdr.ModTime = fi.ModTime().UTC().Truncate(time.Second)

		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}

		if fi.Mode().IsRegular() {
			f, err := os.Open(abs)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(tw, f); err != nil {
				_ = f.Close()
				return nil, err
			}
			_ = f.Close()
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ExtractTarGz(r io.Reader, dest string) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				_ = f.Close()
				return err
			}
			_ = f.Close()
		}
	}

	return nil
}

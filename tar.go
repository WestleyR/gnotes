// Created on: 2022-02-20

// Custom tar functions will likely not work for other needs.

package gnotes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func gzipCompressFile(file string) ([]byte, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}

	return gzipCompress(b), nil
}

func gzipCompress(in []byte) []byte {
	comp := bytes.NewBuffer(nil)
	w := gzip.NewWriter(comp)
	w.Write(in)
	w.Close()

	return comp.Bytes()
}

func gzipExtractFile(file string) ([]byte, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}

	ex, err := gzipExtract(b)
	if err != nil {
		return []byte{}, err
	}

	return ex, nil
}

func gzipExtract(in []byte) ([]byte, error) {
	decmp := bytes.NewBuffer(nil)
	decmp.Write(in)
	r, err := gzip.NewReader(decmp)
	if err != nil {
		return []byte{}, err
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}

	r.Close()

	return b, nil
}

func tarCompress(src string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("got error: %s", err)
		}

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.TrimPrefix(file, src), "/")
		if header.Name == "" {
			return nil
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return []byte{}, fmt.Errorf("failed to open dir: %s", err)
	}

	if err := tw.Close(); err != nil {
		return []byte{}, err
	}

	if err := zr.Close(); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// dst contains the end dir. Note: will not overide files, I think...
func untar(src []byte, dst string) error {
	buff := bytes.NewBuffer(nil)
	buff.Write(src)

	zr, err := gzip.NewReader(buff)
	if err != nil {
		return err
	}

	tr := tar.NewReader(zr)

	for {
		fmt.Println("FOO1")
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		fmt.Println("FOO2")
		if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)
		fmt.Println("FOOBAR:", target)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			fmt.Println("WRITEING TO:", target)
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			fileToWrite.Close()
		}
	}

	return nil
}

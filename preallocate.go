// Package preallocate allocates disk space efficiently via syscall (on
// supported platforms and filesystems) or by writing null bytes.
package preallocate

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const BufferSize = 4 * 1024 * 1024 // 4 MiB

// File preallocates a file via syscall (when supported) or WriteSeeker().
// The file's write offset must be set to the start of the file.
func File(file *os.File, size int64) error {
	if size < 0 {
		return errors.New("invalid preallocation size")
	} else if size == 0 {
		return nil
	}

	return preallocFile(file, size)
}

// FilePath preallocates a file at path (see File).
func FilePath(path string, size int64) (*os.File, error) {
	if size < 0 {
		return nil, errors.New("invalid preallocation size")
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open/create %s", path)
	}

	err = File(file, size)
	if err != nil {
		_ = file.Close()

		return nil, err
	}

	return file, nil
}

// TempFile preallocates a temporary file (see File and ioutil.TempFile).
func TempFile(dir string, pattern string, size int64) (*os.File, error) {
	if size < 0 {
		return nil, errors.New("invalid preallocation size")
	}

	tempFile, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary file")
	}

	return tempFile, File(tempFile, size)
}

// WriteSeeker preallocates an io.WriteSeeker by writing null bytes.
func WriteSeeker(w io.WriteSeeker, size int64) error {
	if size < 0 {
		return errors.New("invalid preallocation size")
	}

	var (
		b         []byte
		remaining = size
		wrote     int
		err       error
	)

	if remaining > BufferSize {
		b = make([]byte, BufferSize)
		for {
			wrote, err = w.Write(b)
			if err != nil {
				return err
			} else if int64(wrote) != BufferSize {
				return errors.New(fmt.Sprintf("failed to preallocate file: write operation interrupted at %d/%d bytes", wrote, BufferSize))
			}

			remaining -= BufferSize
			if remaining < BufferSize {
				break
			}
		}
	}

	if remaining > 0 {
		if b != nil {
			b = b[0:remaining]
		} else {
			b = make([]byte, remaining)
		}

		wrote, err = w.Write(b)
		if err != nil {
			return err
		} else if int64(wrote) != remaining {
			return errors.New(fmt.Sprintf("failed to preallocate file: write operation interrupted at %d/%d bytes", wrote, remaining))
		}
	}

	_, err = w.Seek(0, 0)
	return err
}

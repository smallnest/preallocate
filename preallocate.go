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

// NullBufferSize determines the maximum size of NULL byte blocks written to
// files when falling back to WriteSeeker.
const NullBufferSize = 512 * 1024 // 512 KiB

// File preallocates a file via syscall (when supported) or WriteSeeker.
func File(file *os.File, size int64) error {
	if size < 0 {
		return errors.New("invalid preallocation size")
	} else if size == 0 {
		return nil
	}

	return preallocFile(file, size)
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

	_, err = w.Seek(0, 0)
	if err != nil {
		return err
	}

	if remaining > NullBufferSize {
		b = make([]byte, NullBufferSize)
		for {
			wrote, err = w.Write(b)
			if err != nil {
				return err
			} else if int64(wrote) != NullBufferSize {
				return fmt.Errorf("failed to preallocate file: write operation interrupted at %d/%d bytes", wrote, NullBufferSize)
			}

			remaining -= NullBufferSize
			if remaining < NullBufferSize {
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
			return fmt.Errorf("failed to preallocate file: write operation interrupted at %d/%d bytes", wrote, remaining)
		}
	}

	_, err = w.Seek(0, 0)
	return err
}

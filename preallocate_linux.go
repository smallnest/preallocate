// +build linux

package preallocate

import (
	"os"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func preallocFile(file *os.File, size int64) error {
	fd := file.Fd()
	if fd <= 0 {
		return errors.New("invalid file descriptor")
	}

	err := unix.Fallocate(int(fd), 0, 0, size)
	if err != nil {
		// Filesystem does not support fallocate

		return WriteSeeker(file, size)
	}

	return nil
}

// +build !linux,!windows

package preallocate

import "os"

func preallocFile(file *os.File, size int64) error {
	return WriteSeeker(file, size)
}

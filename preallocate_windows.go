// +build windows

package preallocate

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Support for Windows was cloned from rclone (MIT)
// https://github.com/ncw/rclone/blob/master/backend/local/preallocate_windows.go

type fileAllocationInformation struct {
	AllocationSize uint64
}

type fileFsSizeInformation struct {
	TotalAllocationUnits     uint64
	AvailableAllocationUnits uint64
	SectorsPerAllocationUnit uint32
	BytesPerSector           uint32
}

type ioStatusBlock struct {
	Status, Information uintptr
}

var (
	ntdll                        = windows.NewLazySystemDLL("ntdll.dll")
	ntQueryVolumeInformationFile = ntdll.NewProc("NtQueryVolumeInformationFile")
	ntSetInformationFile         = ntdll.NewProc("NtSetInformationFile")
)

func preallocFile(file *os.File, size int64) error {
	var (
		iosb       ioStatusBlock
		fsSizeInfo fileFsSizeInformation
		allocInfo  fileAllocationInformation
	)

	// Query info about the block sizes on the file system
	_, _, e1 := ntQueryVolumeInformationFile.Call(
		uintptr(file.Fd()),
		uintptr(unsafe.Pointer(&iosb)),
		uintptr(unsafe.Pointer(&fsSizeInfo)),
		uintptr(unsafe.Sizeof(fsSizeInfo)),
		uintptr(3), // FileFsSizeInformation
	)
	if e1 != nil && e1 != syscall.Errno(0) {
		return WriteSeeker(file, size)
	}

	// Calculate the allocation size
	clusterSize := uint64(fsSizeInfo.BytesPerSector) * uint64(fsSizeInfo.SectorsPerAllocationUnit)
	if clusterSize <= 0 {
		return WriteSeeker(file, size)
	}
	allocInfo.AllocationSize = (1 + uint64(size-1)/clusterSize) * clusterSize

	// Ask for the allocation
	_, _, e1 = ntSetInformationFile.Call(
		uintptr(file.Fd()),
		uintptr(unsafe.Pointer(&iosb)),
		uintptr(unsafe.Pointer(&allocInfo)),
		uintptr(unsafe.Sizeof(allocInfo)),
		uintptr(19), // FileAllocationInformation
	)
	if e1 != nil && e1 != syscall.Errno(0) {
		return WriteSeeker(file, size)
	}

	return nil
}

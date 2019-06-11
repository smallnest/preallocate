package preallocate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sizes = []int64{0, 1024, 16777220}

func TestFilePath(t *testing.T) {
	t.Parallel()
	var (
		file     *os.File
		fileInfo os.FileInfo
		size     int64
		err      error
	)

	tmpDir, err := ioutil.TempDir("", "preallocate")
	if err != nil {
		t.Fatal(err)
	}

	for _, size = range sizes {
		file, err = FilePath(filepath.Join(tmpDir, "preallocate.test"), size)
		if err != nil {
			t.Fatal(err)
		}

		fileInfo, err = file.Stat()
		if err != nil {
			t.Fatal(err)
		}

		err = os.Remove(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fileInfo.Size(), size, "actual size does not match preallocation size")
	}
}

func TestTempFile(t *testing.T) {
	t.Parallel()
	var (
		tmpFile     *os.File
		tmpFileInfo os.FileInfo
		size        int64
		err         error
	)

	for _, size = range sizes {
		tmpFile, err = ioutil.TempFile("", "preallocate")
		if err != nil {
			t.Fatal(err)
		}

		err = File(tmpFile, size)
		if err != nil {
			t.Fatal(err)
		}

		tmpFileInfo, err = tmpFile.Stat()
		if err != nil {
			t.Fatal(err)
		}

		err = os.Remove(tmpFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tmpFileInfo.Size(), size, "actual size does not match preallocation size")
	}
}

func TestWriteSeeker(t *testing.T) {
	t.Parallel()
	var (
		tmpFile     *os.File
		tmpFileInfo os.FileInfo
		size        int64
		err         error
	)

	for _, size = range sizes {
		tmpFile, err = ioutil.TempFile("", "preallocate")
		if err != nil {
			t.Fatal(err)
		}

		err = WriteSeeker(tmpFile, size)
		if err != nil {
			t.Fatal(err)
		}

		tmpFileInfo, err = tmpFile.Stat()
		if err != nil {
			t.Fatal(err)
		}

		err = os.Remove(tmpFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tmpFileInfo.Size(), size, "actual size does not match preallocation size")
	}
}

func BenchmarkWriteSeeker(b *testing.B) {
	var (
		tmpFile     *os.File
		tmpFileInfo os.FileInfo
		size        = int64(1677722000)
		err         error
	)

	tmpFile, err = ioutil.TempFile("", "preallocate")
	if err != nil {
		b.Fatal(err)
	}

	err = WriteSeeker(tmpFile, size)
	if err != nil {
		b.Fatal(err)
	}

	tmpFileInfo, err = tmpFile.Stat()
	if err != nil {
		b.Fatal(err)
	}

	err = os.Remove(tmpFile.Name())
	if err != nil {
		b.Fatal(err)
	}

	assert.Equal(b, tmpFileInfo.Size(), size, "actual size does not match preallocation size")
}

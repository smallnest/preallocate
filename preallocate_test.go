package preallocate

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sizes = []int64{0, 1024, 16777220}

func TestFile(t *testing.T) {
	t.Parallel()
	var (
		tmpFile  *os.File
		fileInfo os.FileInfo
		size     int64
		err      error
	)

	tmpDir, err := ioutil.TempDir("", "preallocate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpDir)

	for _, size = range sizes {
		tmpFile, err = os.OpenFile(path.Join(tmpDir, "preallocate.tmp"), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = File(tmpFile, size)
		if err != nil {
			tmpFile.Close()
			t.Fatal(err)
		}

		fileInfo, err = tmpFile.Stat()
		if err != nil {
			tmpFile.Close()
			t.Fatal(err)
		}

		tmpFile.Close()
		err = os.Remove(tmpFile.Name())
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
			tmpFile.Close()
			t.Fatal(err)
		}

		tmpFileInfo, err = tmpFile.Stat()
		if err != nil {
			tmpFile.Close()
			t.Fatal(err)
		}

		tmpFile.Close()
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
			tmpFile.Close()
			t.Fatal(err)
		}

		tmpFileInfo, err = tmpFile.Stat()
		if err != nil {
			tmpFile.Close()
			t.Fatal(err)
		}

		tmpFile.Close()
		err = os.Remove(tmpFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tmpFileInfo.Size(), size, "actual size does not match preallocation size")
	}
}

func BenchmarkWriteSeeker(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var (
			tmpFile     *os.File
			tmpFileInfo os.FileInfo
			size        = int64(67108860)
			err         error
		)

		tmpFile, err = ioutil.TempFile("", "preallocate")
		if err != nil {
			b.Fatal(err)
		}

		err = WriteSeeker(tmpFile, size)
		if err != nil {
			tmpFile.Close()
			b.Fatal(err)
		}

		tmpFileInfo, err = tmpFile.Stat()
		if err != nil {
			tmpFile.Close()
			b.Fatal(err)
		}

		tmpFile.Close()
		err = os.Remove(tmpFile.Name())
		if err != nil {
			b.Fatal(err)
		}

		assert.Equal(b, tmpFileInfo.Size(), size, "actual size does not match preallocation size")
	}
}

package shred

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
)

const MAX_BUFF = 4096
const N_WR = 3

type ShredDirError struct {
	Path string
}

func (e *ShredDirError) Error() string {
	return fmt.Sprintf("the path %q is a directory", e.Path)
}

func (e *ShredDirError) Is(target error) bool {
	_, ok := target.(*ShredDirError)
	return ok
}

type ShredValueError struct {
	Val  string
	Name string
}

func (e *ShredValueError) Error() string {
	return fmt.Sprintf("Invalid value %s for %s", e.Val, e.Name)
}

func (e *ShredValueError) Is(target error) bool {
	_, ok := target.(*ShredValueError)
	return ok
}

type ShredRes struct {
	BytesWritten int64
	Iteration    int64
	BlockSize    int
}

func Shred(path string, buffSize int) (ShredRes, error) {

	var res ShredRes = ShredRes{}

	// Validate path input and ensure it's a file
	if fi, err := os.Stat(path); err != nil {
		return res, err
	} else if fi.IsDir() {
		return res, &ShredDirError{Path: path}
	} else if fi.Size() == 0 {
		return res, &ShredValueError{Val: "0", Name: "fileSize"}
	} else if buffSize <= 0 {
		return res, &ShredValueError{Val: strconv.FormatInt(int64(buffSize), 10), Name: "bufferSize"}
	} else if buffSize > MAX_BUFF {
		buffSize = MAX_BUFF
		fmt.Printf("Length of the buffer limited to %d bytes\n", MAX_BUFF)
	}

	res.BlockSize = buffSize

	f, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return res, err
	}

	defer f.Close()

	// Get file size for efficient block allocation
	var fileSize, bytesWritten int64
	var limitWr int

	fi, err := f.Stat()
	if err != nil {
		return res, err
	}
	fileSize = fi.Size()

	buf := make([]byte, buffSize) // Allocate a buffer for blocks

	for i := 0; i < N_WR; i++ {
		// Seek to beginning of file
		if _, err := f.Seek(0, 0); err != nil {
			return res, err
		}

		// Efficiently overwrite blocks of data with random bytes
		bytesWritten = 0
		for bytesWritten < fileSize {
			// Fill buffer with random data
			_, err := rand.Read(buf)
			if err != nil {
				return res, err
			}

			if fileSize-bytesWritten > int64(buffSize) {
				limitWr = buffSize
			} else {
				limitWr = int(fileSize - bytesWritten)
			}

			// Write the block of data
			n, err := f.Write(buf[:limitWr])
			if err != nil {
				return res, err
			}

			// Flush any buffered data to disk
			err = f.Sync()
			if err != nil {
				return res, err
			}

			bytesWritten += int64(n)

		}
		res.BytesWritten = bytesWritten
		res.Iteration++
	}

	os.Remove(path)

	if res.Iteration != N_WR {
		return res, &ShredValueError{Val: strconv.FormatInt(res.Iteration, 10), Name: "numIterations"}
	} else if res.BytesWritten != fileSize {
		return res, &ShredValueError{Val: strconv.FormatInt(res.BytesWritten, 10), Name: "byesWritten"}
	} else {
		return res, nil
	}
}

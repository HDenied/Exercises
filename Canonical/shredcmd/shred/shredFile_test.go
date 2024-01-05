package shred

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

func TestErrorFileNotExistent(t *testing.T) {
	filename := "nonexistent_file.txt"

	_, got := Shred(filename, 4096)

	if got != nil {
		want := fs.ErrNotExist
		if !errors.Is(got, want) {
			t.Errorf("Expected error: %T, received: %T", want, got)
		}
	} else {
		t.Errorf("Expected error not occurred")
	}

}

func TestErrorNotAfile(t *testing.T) {
	filename := "/tmp"

	_, got := Shred(filename, 4096)

	if got != nil {
		want := &ShredDirError{}
		if !errors.Is(got, want) {
			t.Errorf("Expected error: %T, received: %T", want, got)
		}
	} else {
		t.Errorf("Expected error not occurred")
	}

}

func TestErrorZeroSizeFile(t *testing.T) {
	filename := "test.txt"
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal("Can't generate a local file for testing")

	}
	f.Close()

	defer os.Remove(filename)

	_, got := Shred(filename, 4096)

	if got != nil {
		want := &ShredValueError{Val: "0", Name: "fileSize"}
		if !errors.Is(got, want) {
			t.Fatalf("Expected error: %T, received: %T", want, got)
		}

		if got.Error() != want.Error() {
			t.Errorf("Expected error msg: %s; received: %s", want.Error(), got.Error())
		}

	} else {
		t.Errorf("Expected error not occurred")
	}

}

func TestErrorReadOnlyFile(t *testing.T) {
	filename := "test.txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0400)
	if err != nil {
		t.Fatal("Can't generate a local file for testing")

	}

	f.WriteString("Hello world!")

	f.Close()

	defer os.Remove(filename)

	_, got := Shred(filename, 4096)

	if got != nil {
		want := fs.ErrPermission
		if !errors.Is(got, want) {
			t.Errorf("Expected error: %T, received: %T", want, got)
		}

	} else {
		t.Errorf("Expected error not occurred")
	}

}

func TestGoodPath(t *testing.T) {
	filename := "test.txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal("Can't generate a local file for testing")

	}

	f.WriteString("Hello world! Hello world!")

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	fSize := fi.Size()

	f.Close()

	defer os.Remove(filename)

	res, err := Shred(filename, 4096)

	if err != nil {
		t.Error("No error should occur")
	} else if res.BlockSize != 4096 {
		t.Error("Unexpected block size")
	} else if res.BytesWritten != fSize {
		t.Error("File size mismatch")
	} else if res.Iteration != N_WR {
		t.Error("Wrong number of iteration")
	}

}

func TestGoodPathMaxBuff(t *testing.T) {
	filename := "test.txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal("Can't generate a local file for testing")

	}

	f.WriteString("Hello world! Hello world!")

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	fSize := fi.Size()

	f.Close()

	defer os.Remove(filename)

	res, err := Shred(filename, 10000)

	if err != nil {
		t.Error("No error should occur")
	} else if res.BlockSize != 4096 {
		t.Error("Unexpected block size")
	} else if res.BytesWritten != fSize {
		t.Error("File size mismatch")
	} else if res.Iteration != N_WR {
		t.Error("Wrong number of iteration")
	}

}


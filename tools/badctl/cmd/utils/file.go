package utils

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

func CheckFileExists(t *testing.T, name string) fs.FileInfo {
	stat, err := os.Stat(name)
	if err != nil {
		t.Error(err)
	}

	return stat
}

func RemoveFile(name string) {
	if err := os.RemoveAll(name); err != nil {
		log.Fatal(err)
	}
}

package diff

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// DirChecksum compares the checksum of the contents of dir against the checksums
// in expected. Expected should be a map of all files expected in the directory,
// with the full path and filename as key, and the md5 sum as value.
func DirChecksum(expected map[string]string, dir string) *Result {
	return nil
}

func checkDir(dir string) (map[string]string, error) {
	result := make(map[string]string)
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		filename := dir + "/" + f.Name()
		if f.IsDir() {
			continue
		}
		content, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		hash := md5.New()
		if _, err := io.Copy(hash, content); err != nil {
			return nil, err
		}
		if err := content.Close(); err != nil {
			return nil, err
		}
		result[f.Name()] = hex.EncodeToString(hash.Sum([]byte{}))
	}
	return result, nil
}

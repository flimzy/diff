package diff

import (
	"crypto/md5" // nolint: gas
	"encoding/hex"
	"io"
	"os"
	"strings"
)

// DirChecksum compares the checksum of the contents of dir against the checksums
// in expected. Expected should be a map of all files expected in the directory,
// with the full path and filename as key, and the md5 sum as value.
func DirChecksum(expected map[string]string, dir string) *Result {
	actual, err := checkDir(dir)
	if err != nil {
		return &Result{err: err.Error()}
	}
	return Interface(expected, actual)
}

func checkDir(dir string) (map[string]string, error) {
	result := make(map[string]string)
	err := recurseDir(result, []string{dir})
	return result, err
}

func recurseDir(result map[string]string, parents []string) error {
	dir := strings.Join(parents, "/")
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, f := range files {
		relName := relativeName(parents, f.Name())
		if f.IsDir() {
			result[relName+"/"] = "<dir>"
			if err := recurseDir(result, append(parents, f.Name())); err != nil {
				return err
			}
			continue
		}
		content, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			return err
		}
		hash := md5.New() // nolint: gas
		if _, err := io.Copy(hash, content); err != nil {
			return err
		}
		if err := content.Close(); err != nil {
			return err
		}
		result[relName] = hex.EncodeToString(hash.Sum([]byte{}))
	}
	return nil
}

func relativeName(parents []string, name string) string {
	parts := append(parents[1:], name)
	return strings.Join(parts, "/")
}

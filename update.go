package diff

import (
	"os"
)

func update(updateMode bool, expected *File, actual string) error {
	if !updateMode {
		return nil
	}
	file, err := os.OpenFile(expected.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // nolint: gas
	if err != nil {
		return err
	}
	defer file.Close() // nolint: errcheck
	_, err = file.WriteString(actual)
	return err
}

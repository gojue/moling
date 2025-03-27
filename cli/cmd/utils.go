package cmd

import "os"

// CreateDirectory checks if a directory exists, and creates it if it doesn't
func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0o755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

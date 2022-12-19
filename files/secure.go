package files

import (
	"os"
	"path/filepath"
)

// GetSecretPath use it to get the file location, the key is also the file name
func GetSecretPath(key string, dirSecure string) (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	p := filepath.Join(home, dirSecure, key)

	return p, nil
}

func CreateSecretDirectory(dirSecure string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return os.Mkdir(filepath.Join(home, dirSecure), os.ModePerm)
}

func GetSecretDirectory(dirSecure string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	p := filepath.Join(home, dirSecure)

	return p, nil
}

func GetOrCreateSecureFile(key string, dirSecure string) (*os.File, error) {
	p, err := GetSecretPath(key, dirSecure)
	if err != nil {
		return nil, err
	}

	// check if the file exists. But not append anything, just overwrite.
	return os.OpenFile(p, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}

func ReadSecureFile(key string, dirSecure string) ([]byte, error) {
	p, err := GetSecretPath(key, dirSecure)

	if err != nil {
		return nil, err
	}

	return os.ReadFile(p)
}

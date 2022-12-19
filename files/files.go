package files

import (
	"os"

	"github.com/MoroGasper/devtools/dvtlslog"
)

type Files interface {
	CreateFileWithStruct(fileName, jsonStruct string)
	GetJsonFileWithStruct(jsonFileName string, WithStruct interface{})
	GetJsonFromFile(fileConfigName string) map[string]interface{}

	CreateSecretDirectory(dirSecure string) error
	GetOrCreateSecureFile(key string, dirSecure string) (*os.File, error)
	GetSecretDirectory(dirSecure string) (string, error)
	GetSecretPath(key string, dirSecure string) (string, error)
	ReadSecureFile(key string, dirSecure string) ([]byte, error)
}

type DTSjson struct {
	Log *dvtlslog.DTSlog
}

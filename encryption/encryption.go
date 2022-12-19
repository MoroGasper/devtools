package encryption

import (
	"github.com/MoroGasper/devtools/dvtlslog"
	"github.com/MoroGasper/devtools/files"
)

type Encryption interface {
	GetValue(key string, decrypted bool) (string, error)
	GetVaultKeyData() *vaultkey
	NewVaultKey(encodingKey string)
	PrintAll() error
	SetValue(key, value string) error
}

type vaultkey struct {
	encodingKey string // is the key to encrypt the entry (the value given).
	Files       files.Files
	Log         *dvtlslog.DTSlog
}

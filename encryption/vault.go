package encryption

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	errNotFound  = errors.New("key not found")
	errFileEmpty = errors.New("file is empty")
)

const (
	dirSecure = "secure"
)

// getValue key is the actual key to identify the entry, return the encrypted/encoded value and an error (nil if any
// problem occurred).
func (v *vaultkey) GetValue(key string, decrypted bool) (string, error) {
	// key-values
	value, err := v.read(key, decrypted)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (v *vaultkey) GetVaultKeyData() *vaultkey {
	f, err := v.openVaultKeyInHomeDir()

	if err != nil {
		log.Println(err)
		return nil
	}

	buf := make([]byte, 128)
	_, err = f.Read(buf)
	v.Log.IsFatal(err)

	return &vaultkey{encodingKey: string(buf)}
}

func (v *vaultkey) NewVaultKey(encodingKey string) {
	f, err := v.openVaultKeyInHomeDir()
	v.Log.IsFatal(err)

	_, err = f.Write([]byte(encodingKey))
	v.Log.IsFatal(err)
}

func (v *vaultkey) PrintAll() error {
	b, err := v.Files.ReadSecureFile(".all", dirSecure)

	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

// setValue key is the actual key to identify the entry.
func (v *vaultkey) SetValue(key, value string) error {
	err := v.Files.CreateSecretDirectory(dirSecure)

	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	varFile, err := v.Files.GetOrCreateSecureFile(key, dirSecure)

	if err != nil {
		return err
	}

	//encrypt the value with the encoding key
	encrypted, err := encrypt(v.encodingKey, value)

	if err != nil {
		return err
	}

	_, err = varFile.Write([]byte(encrypted))

	if err != nil {
		return err
	}

	v.saveInAllFile(key, encrypted)
	fmt.Println("file created")
	return nil
}

func (v *vaultkey) createAllPropsFile() (*os.File, error) {
	p, err := v.Files.GetSecretPath(".all", dirSecure)

	if err != nil {
		return nil, err
	}

	return os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
}

func (v *vaultkey) openVaultKeyInHomeDir() (*os.File, error) {
	f, err := v.Files.GetOrCreateSecureFile(".vaultkey", dirSecure)

	if err != nil {
		return nil, err
	}

	return f, nil
}

// read the entire file, returns the values decrypted. The key is the file name
func (v *vaultkey) read(key string, decrypted bool) (string, error) {
	res, err := v.Files.ReadSecureFile(key, dirSecure)

	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errNotFound
	}

	if decrypted {
		d, err := decrypt(v.encodingKey, string(res))

		if err != nil {
			return "", err
		}

		return d, nil
	}

	return string(res), nil
}

func (v *vaultkey) saveInAllFile(key, value string) {
	// var msg string
	f, _err := v.createAllPropsFile()
	v.Log.IsErrorAndMessage(_err, "cannot create ALL properties file.")

	// if _err != nil {
	// 	msg = fmt.Sprintf("cannot create ALL properties file. %s", _err.Error())
	// 	fmt.Println(msg)
	// 	return
	// }
	entry := fmt.Sprintf("+ %s\t|\t%s +\n", key, value)
	_, _err = f.Write([]byte(entry))
	v.Log.IsErrorAndMessage(_err, "cannot writing in ALL properties file.")
	// if _err != nil {
	// 	msg = fmt.Sprintf("cannot writing in ALL properties file. %s", _err.Error())
	// 	fmt.Println(msg)
	// 	return
	// }
}

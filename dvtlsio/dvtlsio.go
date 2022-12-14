package dvtlsio

import (
	"io"
	"os"
	"strings"
)

// INFO: Valida el error
func checkError(err error) {
	if err != nil {
		// fmt.Sprintf("Error : %s", err.Error())
		os.Exit(1)
	}
}

// valida si un archivo existe
func IsFileExist(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

// Crear un directorio
func CreateFolder(dirName string) {

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.Mkdir(dirName, 0755)
		//MkdirAll
		checkError(err)
	}
}

// revisar
func CreateDirAll(directorio string) {

	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0755)
		checkError(err)
	}
}

// INFO: Crea un archivo
func CreateFile(rutaDestino string, data string) {
	// err := ioutil.WriteFile(rutaDestino, []byte(data), 0644)
	err := os.WriteFile(rutaDestino, []byte(data), 0644)
	checkError(err)
}

// lee un archivo y retorna un string de lo leido
//
// "path/filepath" templateName, _ := filepath.Abs(templateName)
func ReadFile(templateName string) string {
	/* data, err := os.ReadFile(templateName)
	if err != nil {
		log.Fatal(err)
	}
	return os.Stdout.Write(data) */
	data, _ := os.ReadFile(templateName)
	return string(data)
}

// Copia un archivo usando operaciones del sistema
func Copy(srcFileDir string, destFileDir string) {
	srcFile, err := os.Open(srcFileDir)
	checkError(err)
	defer srcFile.Close()

	destFile, err := os.Create(destFileDir) // creates if file doesn't exist
	checkError(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	checkError(err)

	err = destFile.Sync()
	checkError(err)
}

// lee un archivo y luego lo copia a otro
func ReadAndCopy(srcFileDir string, destFileDir string) {

	// b, err := ioutil.ReadFile(srcFileDir)
	b, err := os.ReadFile(srcFileDir)
	checkError(err)

	// err = ioutil.WriteFile(destFileDir, b, 0644)
	err = os.WriteFile(destFileDir, b, 0644)
	checkError(err)
}

// modifica un archivo buscando algo
func ChancarFile(templateName string, MapForReplace map[string]string) {
	data := ReplaceTextInFile(templateName, MapForReplace)
	CreateFile(templateName, data)
}

// Crea un archivo nuevo partiendo de una plantilla y un arreglo de opciones a remplazar
func NewFileforTemplate(newName string, templateName string, MapForReplace map[string]string) {
	data := ReplaceTextInFile(templateName, MapForReplace)
	CreateFile(newName, data)
}

// remplaza info en un archivo luego lo pasa a una variable
func ReplaceTextInFile(templateName string, MapForReplace map[string]string) string {
	input := ReadFile(templateName)
	for key, value := range MapForReplace {
		input = strings.Replace(input, key, value, -1)
	}
	return input
}

// añade al final del archivo un string
func AddEndToFile(destFileDir string, data string) {
	b, err := os.OpenFile(destFileDir, os.O_APPEND|os.O_WRONLY, 0600)
	checkError(err)
	defer b.Close()

	_, err = b.WriteString(data)
	checkError(err)
}

// falta probar....... agrega muchas lineas nuevas
func AddArrayEndToFile(destFileDir string, datas []string) {
	for _, data := range datas {
		AddEndToFile(destFileDir, data)
	}
}

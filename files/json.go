package files

import (
	"encoding/json"
	"io"
	"os"

	"github.com/MoroGasper/devtools/dvtlsio"
)

func (s *DTSjson) CreateFileWithStruct(fileName, jsonStruct string) {
	if jsonStruct == "" {
		jsonStruct = `{"param1":"dat0","param2":"dat1"}`
	}
	dvtlsio.CreateFile(fileName, jsonStruct)
}

/*
lee un archivo buscando una estructura y la devuelve
para usar debes hacer esto

	type FilesStruct struct {
		Name  string
	}

var WithStruct []FilesStruct
GetJsonFileWithStruct("miarchivo.json", &WithStruct)
*/
func (s *DTSjson) GetJsonFileWithStruct(jsonFileName string, WithStruct interface{}) {
	jsonFile, err := os.Open(jsonFileName)
	s.Log.IsFatal(err)
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	s.Log.IsFatal(err)
	s.Log.IsFatal(json.Unmarshal([]byte(byteValue), WithStruct))
}

/*
lee un archivo sin estructuras
para usar debes hacer esto

GetJsonFromFile("miarchivo.json")
*/
func (s *DTSjson) GetJsonFromFile(fileConfigName string) map[string]interface{} {

	jsonFile, err := os.Open(fileConfigName)
	if err != nil {
		if err.Error() == "open "+fileConfigName+": no such file or directory" {
			s.Log.IsFatal(err)
		}
	}
	defer jsonFile.Close()
	var result map[string]interface{}
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

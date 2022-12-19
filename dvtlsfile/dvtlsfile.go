package dvtlsfile

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/MoroGasper/devtools/dvtlsio"
	"github.com/MoroGasper/devtools/dvtlslog"
)

type DTSjson struct {
	Log *dvtlslog.DTSlog
}

func (i *DTSjson) CreateFileWithStruct(fileName, jsonStruct string) {
	if jsonStruct == "" {
		jsonStruct = `{"param1":"dat0","param2":"dat1"}`
	}
	dvtlsio.CreateFile(fileName, jsonStruct)
}

/*
lee un archivo sin estructuras
para usar debes hacer esto

GetJsonFromFile("miarchivo.json")
*/
func (i *DTSjson) GetJsonFromFile(fileConfigName string) map[string]interface{} {

	jsonFile, err := os.Open(fileConfigName)
	if err != nil {
		if err.Error() == "open "+fileConfigName+": no such file or directory" {
			i.Log.IsFatal(err)
		}
	}
	defer jsonFile.Close()
	var result map[string]interface{}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
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
func (i *DTSjson) GetJsonFileWithStruct(jsonFileName string, WithStruct interface{}) {
	jsonFile, err := os.Open(jsonFileName)
	i.Log.IsFatal(err)
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	// byteValue, err := os.R
	i.Log.IsFatal(err)
	i.Log.IsFatal(json.Unmarshal([]byte(byteValue), WithStruct))
}

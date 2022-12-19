package dvtlscsv

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/MoroGasper/devtools/dvtlslog"
)

type CsvInMemory struct {
	InMemory []map[string]string
	Log      *dvtlslog.DTSlog
	Head     []string
}

func (s *CsvInMemory) DebufOff() {
	if s.Log != nil {
		s.Log.IsDebug = false
	}
}

func (s *CsvInMemory) PrepareLog(IsDebug bool, PrinterLogs bool, PrinterScreen bool) {
	s.Log = &dvtlslog.DTSlog{
		IsDebug:       IsDebug,
		PrinterLogs:   PrinterLogs,
		PrinterScreen: PrinterScreen,
	}
	s.Log.SetInitProperty()
}

func (s *CsvInMemory) PrepareDefaultLog() {
	s.PrepareLog(true, true, true)

}

func (s *CsvInMemory) OpenFile(file string) *csv.Reader {
	f, err := os.Open(file)
	s.Log.IsFatal(err)
	r := csv.NewReader(f)
	return r
}

func (s *CsvInMemory) ReadHead(csvreader *csv.Reader) {
	var head []string
	s.Log.Debug("Inciando lectura de head", nil)
	record, err := csvreader.Read()
	if err == io.EOF {
		s.Log.Error(err.Error(), nil)
		return
	}
	for value := range record {
		head = append(head, record[value])
		s.Log.Debug("valor guardado", record[value])
	}
	s.Head = head
	s.Log.Debug("Terminando lectura de head", nil)
}

func (s *CsvInMemory) CreateCsvInMemory(head []string, r *csv.Reader) {
	var file map[string]string
	// file = make(map[string]string)
	s.Log.Debug("Inciando lectura de CUERPO", nil)
	for {
		file = make(map[string]string)
		record, err := r.Read()
		s.Log.Debug("leyendo linea", record)
		if err == io.EOF {
			s.Log.Warn(err.Error(), nil)
			break
		}

		s.Log.IsFatal(err)

		h := 0
		for value := range record {
			if head[h] != "-" {
				file[head[h]] = record[value]
				//s.Log.Debug(head[h]+"->"+record[value], record[value])
			}
			h = h + 1
		}
		s.Log.Debug("Agregando a arreglo base la siguiente fila")
		s.Log.Debug("fila", file)

		s.InMemory = append(s.InMemory, file)
		s.Log.Debug("Terminando lectura de CUERPO", nil)
	}
}

package exe

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/MoroGasper/devtools/dvtlslog"
)

type DTSexe struct {
	Log           *dvtlslog.DTSlog
	Arg           []string
	Cmd           *exec.Cmd
	Executable    string
	IsDebug       bool
	PrinterLogs   bool
	PrinterScreen bool
	ShowStd       bool
	ShowErr       bool
	FinalPath     string
	Die           bool
}

// Desactiva el debug
func (s *DTSexe) LogDebugOff() {
	s.IsDebug = false
	if s.Log != nil {
		s.Log.DebugOff()
		s.Log.Debug("Debug off")
	}
}
func (s *DTSexe) LogSilence() {
	if s.Log != nil {
		s.Log.Silence()
	}
}

// Configuración por defecto del log
func (s *DTSexe) PrepareDefaultLog() {
	s.PrepareLog(true, true, true)
}

func (s *DTSexe) PrepareLog(IsDebug bool, PrinterLogs bool, PrinterScreen bool) {
	s.IsDebug = IsDebug
	s.PrinterLogs = PrinterLogs
	s.PrinterScreen = PrinterScreen
	s.Log = dvtlslog.PrepareLog(s.IsDebug, s.PrinterLogs, s.PrinterScreen)
}

func (s *DTSexe) PrepareDefaultjExe(Executable string) {
	s.Executable = Executable
	s.ShowStd = true
	s.ShowErr = true
	s.PrepareDefaultLog()

}

func (s *DTSexe) PrepareDefaultWithLogSilence(Executable string) {
	s.Executable = Executable
	s.ShowStd = false
	s.ShowErr = false
	s.PrepareLog(false, false, false)

}

func (s *DTSexe) PreparejExe(Executable string, ShowStd bool, ShowErr bool, IsDebug bool, PrinterLogs bool, PrinterScreen bool) {
	s.Executable = Executable
	s.ShowStd = ShowStd
	s.ShowErr = ShowErr
	s.PrepareLog(IsDebug, PrinterLogs, PrinterScreen)
}

func (s *DTSexe) CommandAndRun(withArgument bool, die bool) {
	s.Command(s.Executable, withArgument)
	s.Run(die)
}

func (s *DTSexe) CommandInternal(withArgument bool) {
	s.Command(s.Executable, withArgument)
}

func (s *DTSexe) Command(exectuble string, withArgument bool) {
	if withArgument {
		s.Cmd = exec.Command(exectuble, s.Arg...)
	} else {
		s.Cmd = exec.Command(exectuble)
	}
	if s.FinalPath != "" {
		absPath, _ := filepath.Abs(s.FinalPath)
		s.Cmd.Dir = absPath
	}
	s.Log.Debug("Commando:\n%s\n", s.Cmd)
}

func (s *DTSexe) ExecuteWithArg(arg ...string) {
	s.Arg = arg
	s.CommandInternal(true)
	s.Run(s.Die)
}

func (s *DTSexe) ExecuteWithArgAndData(data string, arg ...string) {
	s.Arg = arg
	s.CommandInternal(true)
	s.RunWithData(data, s.Die)
}

func (s *DTSexe) Run(die bool) (string, string, error) {
	var stdout, stderr bytes.Buffer
	s.Cmd.Stdout = &stdout
	s.Cmd.Stderr = &stderr
	err := s.Cmd.Run()
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	outStr, errStr := stdout.String(), stderr.String()
	if s.ShowStd {
		if outStr != "" {
			s.Log.Debug("exe output: \n%s\n", outStr)
		}
	}
	if errStr != "" {
		s.Log.IsErrorAndDie(errors.New(errStr), die)
	}
	if err != nil {
		s.Log.IsErrorAndDie(err, die)
	}
	return outStr, errStr, err
}

func (s *DTSexe) GenerateAbsolutePath(FileName string) string {
	absPath, _ := filepath.Abs("./")
	return absPath + FileName
}

func (s *DTSexe) AddParameterWithAbsolutePath(Paramterindex string, FileName string) []string {
	FileName = s.GenerateAbsolutePath(FileName)
	s.AddParameter(Paramterindex, FileName)
	return s.Arg
}

func (s *DTSexe) AddParameter(Index string, Value string) []string {
	s.Arg = append([]string{Index, Value}, s.Arg...)
	return s.Arg
}

func (s *DTSexe) Addflag(flag string) []string {
	s.Arg = append([]string{flag}, s.Arg...)
	return s.Arg
}

func (s *DTSexe) RunWithData(data string, die bool) (string, string, error) {
	buffer := bytes.Buffer{}
	buffer.Write([]byte(data))
	s.Cmd.Stdin = &buffer
	return s.Run(die)
}

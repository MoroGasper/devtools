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
	/**/
	FinalPath string
	Die       bool
}

// Desactiva el debug
func (i *DTSexe) LogDebugOff() {
	i.IsDebug = false
	if i.Log != nil {
		i.Log.DebugOff()
		i.Log.Debug("Debug off")
	}
}
func (i *DTSexe) LogSilence() {
	if i.Log != nil {
		i.Log.Silence()
	}
}

// Configuración por defecto del log
func (i *DTSexe) PrepareDefaultLog() {
	i.PrepareLog(true, true, true)
}

func (i *DTSexe) PrepareLog(IsDebug bool, PrinterLogs bool, PrinterScreen bool) {
	i.IsDebug = IsDebug
	i.PrinterLogs = PrinterLogs
	i.PrinterScreen = PrinterScreen
	i.Log = dvtlslog.PrepareLog(i.IsDebug, i.PrinterLogs, i.PrinterScreen)
}

func (i *DTSexe) PrepareDefaultjExe(Executable string) {
	i.Executable = Executable
	i.ShowStd = true
	i.ShowErr = true
	i.PrepareDefaultLog()

}

func (i *DTSexe) PrepareDefaultWithLogSilence(Executable string) {
	i.Executable = Executable
	i.ShowStd = false
	i.ShowErr = false
	i.PrepareLog(false, false, false)

}

func (i *DTSexe) PreparejExe(Executable string, ShowStd bool, ShowErr bool, IsDebug bool, PrinterLogs bool, PrinterScreen bool) {
	i.Executable = Executable
	i.ShowStd = ShowStd
	i.ShowErr = ShowErr
	i.PrepareLog(IsDebug, PrinterLogs, PrinterScreen)
}

func (i *DTSexe) CommandAndRun(withArgument bool, die bool) {
	i.Command(i.Executable, withArgument)
	i.Run(die)
}

func (i *DTSexe) CommandInternal(withArgument bool) {
	i.Command(i.Executable, withArgument)
}

func (i *DTSexe) Command(exectuble string, withArgument bool) {
	if withArgument {
		i.Cmd = exec.Command(exectuble, i.Arg...)
	} else {
		i.Cmd = exec.Command(exectuble)
	}
	if i.FinalPath != "" {
		absPath, _ := filepath.Abs(i.FinalPath)
		i.Cmd.Dir = absPath
	}
	i.Log.Debug("Commando:\n%s\n", i.Cmd)
}

func (i *DTSexe) ExecuteWithArg(arg ...string) {
	i.Arg = arg
	i.CommandInternal(true)
	i.Run(i.Die)
}

func (i *DTSexe) ExecuteWithArgAndData(data string, arg ...string) {
	i.Arg = arg
	i.CommandInternal(true)
	i.RunWithData(data, i.Die)
}

func (i *DTSexe) Run(die bool) (string, string, error) {
	var stdout, stderr bytes.Buffer
	i.Cmd.Stdout = &stdout
	i.Cmd.Stderr = &stderr
	err := i.Cmd.Run()
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	outStr, errStr := stdout.String(), stderr.String()
	if i.ShowStd {
		if outStr != "" {
			i.Log.Debug("exe output: \n%s\n", outStr)
		}
	}
	if errStr != "" {
		i.Log.IsErrorAndDie(errors.New(errStr), die)
	}
	if err != nil {
		i.Log.IsErrorAndDie(err, die)
	}
	return outStr, errStr, err
}

func (i *DTSexe) GenerateAbsolutePath(FileName string) string {
	absPath, _ := filepath.Abs("./")
	return absPath + FileName
}

func (i *DTSexe) AddParameterWithAbsolutePath(Paramterindex string, FileName string) []string {
	FileName = i.GenerateAbsolutePath(FileName)
	i.AddParameter(Paramterindex, FileName)
	return i.Arg
}

func (i *DTSexe) AddParameter(Index string, Value string) []string {
	i.Arg = append([]string{Index, Value}, i.Arg...)
	return i.Arg
}

func (i *DTSexe) Addflag(flag string) []string {
	i.Arg = append([]string{flag}, i.Arg...)
	return i.Arg
}

func (i *DTSexe) RunWithData(data string, die bool) (string, string, error) {
	buffer := bytes.Buffer{}
	buffer.Write([]byte(data))
	i.Cmd.Stdin = &buffer
	return i.Run(die)
}

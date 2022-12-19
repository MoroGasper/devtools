package dvtlslog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type DTSlog struct {
	IsDebug       bool
	PrinterLogs   bool
	PrinterScreen bool
	LogInfo       *log.Logger
	LogError      *log.Logger
	DirLogs       string //os.Mkdir("logs", 0755)
	DirErroLogs   string //os.Mkdir("logs/error", 0755)
	LogFileName   string
	Trace         bool
}

func splitLast(file string) string {
	spliting := strings.Split(file, "/")
	x := len(spliting)
	return spliting[x-1]
}

func (s *DTSlog) Silence() {
	s.IsDebug = false
	s.PrinterLogs = false
	s.PrinterScreen = false
}
func (s *DTSlog) DebugOff() {
	s.IsDebug = false
}

func PrepareLog(IsDebug bool, PrinterLogs bool, PrinterScreen bool) *DTSlog {
	Log := &DTSlog{
		IsDebug:       IsDebug,
		PrinterLogs:   PrinterLogs,
		PrinterScreen: PrinterScreen,
		LogInfo:       &log.Logger{},
		LogError:      &log.Logger{},
		DirLogs:       "logs",
		DirErroLogs:   "logs/error",
		LogFileName:   "run",
		Trace:         false,
	}
	Log.SetInitProperty()
	return Log
}

func PrepareDefaultLog() *DTSlog {
	return PrepareLog(true, true, true)
}

func (s *DTSlog) Write(typems string, format string, a ...interface{}) string {
	var ok bool
	ok = true
	var file, infofile string
	var line, actualline int
	var debuging, text string

	for in := 0; ok; in++ {
		_, file, actualline, ok = runtime.Caller(in)
		if s.Trace {
			proctemp := splitLast(file)

			if proctemp == "proc.go" {
				break
			}
			if proctemp != "dvtlslog.go" {
				infofile = proctemp
				debuging = debuging + " (" + proctemp + ":" + strconv.Itoa(actualline) + ")"
			}

			line = actualline
			if s.IsDebug {
				infofile = debuging
			}
		} else {
			_, _, _, ok = runtime.Caller(in + 1)
			if !ok {
				_, file, actualline, ok = runtime.Caller(in - 3)
				proctemp := splitLast(file)
				if proctemp == "proc.go" {
					break
				}
				infofile = proctemp
				debuging = debuging + " (" + proctemp + ":" + strconv.Itoa(actualline) + ")"
				line = actualline
				if s.IsDebug {
					infofile = debuging
				}
			}
		}
	}

	info := fmt.Sprintf(format, a...)

	if s.IsDebug {
		text = fmt.Sprintf("%s: %s", infofile, info)
	} else {
		text = fmt.Sprintf("%s:%d: %s", infofile, line, info)
	}

	//fmt.Println(typems, text)
	return typems + text

}
func (s *DTSlog) Debug(format string, a ...interface{}) {
	if s.IsDebug {
		texto := s.Write("[DEBUG]:", format, a...)
		color.White(texto)
		if s.PrinterLogs {
			s.LogInfo.Println(texto)
		}
	}
}
func (s *DTSlog) Fatal(format string, a ...interface{}) {
	s.Error(format, a...)
	os.Exit(1)
}

func (s *DTSlog) IsFatal(err error) {
	if err != nil {
		s.Fatal(err.Error(), nil)
	}
}

func (s *DTSlog) IsErrorAndDie(err error, die bool) {
	if die {
		s.IsFatal(err)
	}
	if err != nil {
		s.Error(err.Error())
	}
}

func (s *DTSlog) IsErrorAndMessage(err error, message string) {
	msg := fmt.Sprintf(message+" %s", err.Error())
	if err != nil {
		s.Error(msg)
	}
}

func (s *DTSlog) Error(format string, a ...interface{}) {
	texto := s.Write("[Error]:", format, a...)
	if s.PrinterScreen {
		color.Red(texto)
	}
	if s.PrinterLogs {
		s.LogError.Println(texto)
		s.LogInfo.Println(texto)
	}
}

func (s *DTSlog) Info(format string, a ...interface{}) {
	texto := s.Write("[Info]:", format, a...)
	if s.PrinterScreen {
		// color.White(texto)
		color.Cyan(texto)
	}
	if s.PrinterLogs {
		s.LogInfo.Println(texto)
	}
}

func (s *DTSlog) Warn(format string, a ...interface{}) {
	texto := s.Write("[Warn]:", format, a...)
	if s.PrinterScreen {
		color.Yellow(texto)
	}
	if s.PrinterLogs {
		s.LogInfo.Println(texto)
	}
}

func (s *DTSlog) SetInitProperty() {
	if s.PrinterLogs {
		if s.DirLogs == "" {
			s.DirLogs = "logs"
		}
		if s.DirErroLogs == "" {
			s.DirErroLogs = "logs/error"
		}
		if s.LogFileName == "" {
			s.LogFileName = "run"
		}
		os.Mkdir(s.DirLogs, 0755)
		os.Mkdir(s.DirErroLogs, 0755)
		logFile, err := os.OpenFile("./"+s.DirLogs+"/"+s.LogFileName+"_"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalln("open log file failed", err)
		}
		logFileError, err1 := os.OpenFile("./"+s.DirErroLogs+"/"+s.LogFileName+"_"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
		if err1 != nil {
			log.Fatalln("open log file 'error' failed ", err)
		}
		s.LogInfo = log.New(io.MultiWriter(logFile), "", log.Ldate|log.Ltime)       //LogInfo.Println(1, 2, 3)
		s.LogError = log.New(io.MultiWriter(logFileError), "", log.Ldate|log.Ltime) //LogError.Println(4, 5, 6)
	}
}

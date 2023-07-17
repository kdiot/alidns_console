package utility

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	LOG_UNDEFINED = iota
	LOG_DEBUG
	LOG_INFO
	LOG_WARNING
	LOG_ERROR
	LOG_FATAL
)

type any = interface{}

type LogLevel int

var (
	logLevel   LogLevel
	logDebug   *log.Logger
	logInfo    *log.Logger
	logWarning *log.Logger
	logError   *log.Logger
	logFatal   *log.Logger
)

func (lvl *LogLevel) String() string {
	switch *lvl {
	case LOG_DEBUG:
		return "debug"
	case LOG_INFO:
		return "info"
	case LOG_WARNING:
		return "warning"
	case LOG_ERROR:
		return "error"
	case LOG_FATAL:
		return "fatal"
	default:
		return ""
	}
}

func (lvl *LogLevel) Set(value string) error {
	switch value {
	case "debug":
		*lvl = LOG_DEBUG
	case "info":
		*lvl = LOG_INFO
	case "warning":
		*lvl = LOG_WARNING
	case "error":
		*lvl = LOG_ERROR
	case "fatal":
		*lvl = LOG_FATAL
	default:
		return fmt.Errorf("'%s' is not a valid LogLevel enumeration value", value)
	}
	return nil
}

func (lvl *LogLevel) IsValid() bool {
	if *lvl >= LOG_DEBUG && *lvl <= LOG_FATAL {
		return true
	} else {
		return false
	}
}

func (lvl *LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(lvl.String()), nil
}

func (lvl *LogLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	} else {
		return lvl.Set(s)
	}
}

func SetLogLevel(level LogLevel) {
	if level.IsValid() {
		logLevel = level
	}
}

func GetLogLevel() LogLevel {
	return logLevel
}

func SetLogFile(fileName string) error {
	dir := path.Dir(fileName)
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logDebug.SetOutput(file)
	logInfo.SetOutput(file)
	logWarning.SetOutput(file)
	logError.SetOutput(file)
	logFatal.SetOutput(file)

	return nil
}

func Debug(v ...any) {
	if logLevel == LOG_DEBUG {
		logDebug.Println(v...)
	}
}

func Debugf(format string, v ...any) {
	if logLevel == LOG_DEBUG {
		logDebug.Printf(format, v...)
	}
}

func Info(v ...any) {
	if logLevel <= LOG_INFO {
		logInfo.Println(v...)
	}
}

func Infof(format string, v ...any) {
	if logLevel <= LOG_INFO {
		logInfo.Printf(format, v...)
	}
}

func Warning(v ...any) {
	if logLevel <= LOG_WARNING {
		logWarning.Println(v...)
	}
}

func Warningf(format string, v ...any) {
	if logLevel <= LOG_WARNING {
		logWarning.Printf(format, v...)
	}
}

func Error(v ...any) {
	if logLevel <= LOG_ERROR {
		logError.Println(v...)
	}
}

func Errorf(format string, v ...any) {
	if logLevel <= LOG_ERROR {
		logError.Printf(format, v...)
	}
}

func Fatal(v ...any) {
	if logLevel <= LOG_FATAL {
		logFatal.Fatalln(v...)
	}
}

func Fatalf(format string, v ...any) {
	if logLevel <= LOG_FATAL {
		logFatal.Fatalf(format, v...)
	}
}

func init() {
	logLevel = LOG_INFO
	logDebug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lmsgprefix)
	logInfo = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmsgprefix)
	logWarning = log.New(os.Stdout, "WARNING: ", log.LstdFlags|log.Lmsgprefix)
	logError = log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Lmsgprefix)
	logFatal = log.New(os.Stdout, "FATAL: ", log.LstdFlags|log.Lmsgprefix)
}

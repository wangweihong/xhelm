package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	fileutil "util/file"

	"github.com/Sirupsen/logrus"
)

var (
	lockFile *os.File
	logFile  *os.File
	Logger   = logrus.WithFields(logrus.Fields{"pkg": "xhelm"})
	logdir   = "/var/log/xhelm"
	logName  = "xhelm.log"
)

func printFunc(info interface{}) string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)

	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	file, line := fun.FileLine(fpcs[0] - 1)
	//return fmt.Sprintf("File(%v) Line(%v) Func(%v): %v", file, line, fun.Name(), err)
	//	return fmt.Sprintf("%v %v:%v %v %v", time.Now().String(), filepath.Base(file), line, fun.Name(), info)
	return fmt.Sprintf("%v:%v %v %v", filepath.Base(file), line, fun.Name(), info)
}

func Init() error {

	if logName != "" {
		if err := fileutil.MkdirIfNotExists(logdir); err != nil {
			return err
		}
		logFile, err := os.OpenFile(filepath.Join(logdir, logName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		syscall.Dup2(int(logFile.Fd()), 1)
		syscall.Dup2(int(logFile.Fd()), 2)
		//logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02T15:04:05"})
		//logrus.SetFormatter(&logger.TextFormatter{TimestampFormat: "2006-01-02T15:04:05"})
		logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000000"})
		logrus.SetOutput(logFile)
	} else {
		logrus.SetOutput(os.Stdout)
	}
	return nil
}

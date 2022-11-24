package mylog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// RedirectStderr to the file passed in
func RedirectStderr() (err error) {
	logFile, err := os.OpenFile(ErrorLogFilename(), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	err = syscall.Dup3(int(logFile.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		return
	}
	return
}

func ErrorLogFilename() string {
	return filepath.Join(LogDir(), fmt.Sprintf("%s-error.log", strings.ToLower(EXEName())))
}

func LogDir() string {
	dir := filepath.Join(CWD(), "logs")
	EnsureDir(dir)
	return dir
}

func EXEName() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func EnsureDir(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return
		}
	}
	return
}

func CWD() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}

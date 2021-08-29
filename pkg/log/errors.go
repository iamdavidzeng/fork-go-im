package log

import (
	"encoding/json"
	"fmt"
	"fork_go_im/im/utils"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ErrorInfo struct {
	Time         string `json:"time"`
	FileName     string `json:"filename"`
	Function     string `json:"function"`
	ErrorMessage string `json:"error_message"`
	Line         int    `json:"line"`
}

func Warning(str string) {
	timeString := time.Unix(time.Now().Unix(), 0).Format("2000-01-01 00:00:00")
	filename, line, functionName := "?", 0, "?"
	pc, filename, line, ok := runtime.Caller(2)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}
	var msg = ErrorInfo{
		Time:         timeString,
		FileName:     filename,
		ErrorMessage: str,
		Function:     functionName,
		Line:         line,
	}
	jsons, err := json.Marshal(msg)
	errorJson := string(jsons) + "\n"
	path := utils.GetCurrentDirectory() + "/log"
	logFile := path + "/" + timeString + "-error.log"
	_, exist := os.Stat(path)
	if os.IsNotExist(exist) {
		os.Mkdir(path, os.ModePerm)
	}
	file, err := os.Open(logFile)
	if err != nil {
		files, err := os.Create(logFile)
		defer files.Close()
		if err != nil {
			fmt.Println(err)
		}
		files.Write([]byte(errorJson))
	} else {
		defer file.Close()
		file.Write([]byte(errorJson))
	}
}

func LogError(err error) {
	if err != nil {
		Warning(err.Error())
	}
}

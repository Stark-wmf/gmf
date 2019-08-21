package gmf

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type errorString struct {
	s string
}

type errorInfo struct {
	Time     string `json:"time"`
	Alarm    string `json:"alarm"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Funcname string `json:"funcname"`
}

func (e *errorString) Error() string {
	return e.s
}
//一般记录
func Info (text string) error {
	Alarm("INFO", text)
	return &errorString{text}
}

// 告警
func Warn (text string) error {
	Alarm("WARN", text)
	return &errorString{text}
}

// debug
func Debug (text string) error {
	Alarm("DEBUG", text)
	return &errorString{text}
}

// 最高级别，关闭日志
func Off (text string) error {
	Alarm("OFF", text)
	return &errorString{text}
}

// 告警方法
func  Alarm(level string, str string) {
	// 当前时间
	currentTime :=time.Now().Format("2006-01-02 15:04:05")

	// 定义 文件名、行号、方法名
	fileName, line, functionName := "?", 0 , "?"

	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	var msg = errorInfo {
		Time     : currentTime,
		Alarm    : level,
		Message  : str,
		Filename : fileName,
		Line     : line,
		Funcname : functionName,
	}

	jsons, errs := json.Marshal(msg)

	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}

	errorJsonInfo := string(jsons)

	fmt.Println(errorJsonInfo)

	if level == "INFO" {
		// 执行发邮件
		fmt.Println("-------INFO-------")
        fmt.Println(msg)
	} else if level == "WARN" {
		// 执行发短信
		fmt.Println("-------WARN-------")
        fmt.Println(msg)
	} else if level == "DEBUG" {
		// 执行发微信
		fmt.Println("-------DEBUG-------")
        fmt.Println(msg)

	} else if level == "OFF" {
		// 执行记日志
		fmt.Println("-------OFF-------")
		runtime.Goexit()
	}
}

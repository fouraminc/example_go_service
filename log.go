package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type FrontEndData struct {
	Name string
	Data string
}

func logInfo(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF--" + data)
	appendDataToLog("INF", reference, data)
}
func logError(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF--" + data)
	appendDataToLog("INF", reference, data)
}

func logFrontEnd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data FrontEndData
	timer := time.Now()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logInfo("MAIN","Error logging front end data" + time.Since(timer).String())
		return
	}
	reference := data.Name
	logData := data.Data
	
	appendDataToLog("INF", "FE", reference+ " " + logData)
}

func appendDataToLog(logLevel string, reference string, data string) {

	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	logDirectory := filepath.Join(".", "log")
	logFileName := reference +" "+ time.Now().Format(logNameDateTimeFormat)+ ".log"
	logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
	logData := time.Now().Format("2006-01-02 15:04:05.000 ") + reference + " " + logLevel + " " + data
	writingSync.Lock()
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- Cannot open file: " + err.Error())
		writingSync.Unlock()
		return
	}
	_, err = f.WriteString(logData + "\r\n")
	if err != nil {
		fmt.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- Cannot write to file: " + err.Error())
	}
	err = f.Close()
	if err != nil {
		fmt.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- Cannot close file: " + err.Error())
	}
	writingSync.Unlock()
}

func logDirectoryCheck() {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	dir := getActualDirectory(dateTimeFormat)
	createLogDirectory(dir, dateTimeFormat)

}

func createLogDirectory(dir string, dateTimeFormat string) {
	logDirectory := filepath.Join(dir, "log")
	_, checkPathError := os.Stat(logDirectory)
	logDirectoryExists := checkPathError == nil
	if logDirectoryExists {
		return
	}
	switch runtime.GOOS {
	case "windows":
		{
			err := os.Mkdir(logDirectory, 0777)
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to create directory for log file: " + err.Error())
				return
			}
			fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --INF-- Log directory created")
		}
	default:
		{
			err := os.MkdirAll(logDirectory, 0777)
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to create directory for log file: " + err.Error())
				return
			}
			fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --INF-- Log directory created")
		}
	}
}

func getActualDirectory(dateTimeFormat string) string {
	var dir string
	switch runtime.GOOS {
	case "windows":
		{
			executable, err := os.Executable()
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to read actual directory: " + err.Error())
			}
			dir = filepath.Dir(executable)
		}
	default:
		{
			executable, err := os.Getwd()
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to read actual directory: " + err.Error())
			}
			dir = executable
		}
	}
	return dir
}
package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CheckError(err error, errText string) {
	if err != nil {
		log.Fatalln(errText, err)
	}
}

func LogInit() *log.Logger {
	logFile := "log.txt"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	CheckError(err, "Failed to open log file")

	logger := log.New(file,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}

func SearchFile(filePartName string, dir string) (string, error) {
	fullFileName := "" 
	filesInfo, err := ioutil.ReadDir(dir)
	CheckError(err, "")

	for i := range filesInfo {
		if strings.Contains(strings.ToLower(filesInfo[i].Name()), filePartName) {
			fullFileName = filesInfo[i].Name()
			break
		}
	}

	if fullFileName == "" {
		return "", errors.New("File not found")
	}

	return fullFileName, nil
}

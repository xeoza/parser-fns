package main

import (
	"fmt"
	"sync"
)

func main() {
	logger := LogInit()
	db := InitDb()
	defer db.Close()
	var w sync.WaitGroup
	fmt.Println("Обработка XML-файлов")
	w.Add(1)
	ObjAddrobj := &XmlObjectAddrobj{}
	go ExportFromXmlInPsql(CreateTableAddrobj, ObjAddrobj, &w, db, "xml", logger)
	w.Wait()
	fmt.Println("Все процессы завершились")
}
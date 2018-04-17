package main

import (
	"encoding/xml"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func ExtractXmlObjectName(xmlObject interface{}) xmlObjectName {
	v := reflect.ValueOf(xmlObject).Elem()
	field := v.Type().Field(0)
	obj := xmlObjectName{tableName: field.Tag.Get("db"), elementName: field.Tag.Get("xml")}
	return obj
}

func ExtractFeilds(xmlObject interface{}) []string {
	v := reflect.ValueOf(xmlObject).Elem()
	fields := make([]string, v.NumField()-1)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Type.String() != "xml.Name" {
			fields[i-1] = field.Tag.Get("db")
		}
	}
	return fields
}

func ExtractValues(xmlObject interface{}) []interface{} {
	s := reflect.ValueOf(xmlObject).Elem()
	values := make([]interface{}, s.NumField()-1)

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Type().Name() != "xml.Name" {
			if f.Kind() == reflect.String {
				values[i-1] = f.String()
			} else if f.Kind() == reflect.Int {
				values[i-1] = f.Int()
			} else if f.Kind() == reflect.Bool {
				values[i-1] = f.Bool()
			}
		}
	}

	return values
}

func ExportFromXmlInPsql(structXml func(tableName string) string, xmlObject interface{}, w *sync.WaitGroup, db *sqlx.DB, format string, logger *log.Logger) {

	defer w.Done()
	objName := ExtractXmlObjectName(xmlObject)
	fields := ExtractFeilds(xmlObject)
	searchFileName := objName.tableName
	objName.tableName = strings.TrimSuffix(objName.tableName, "_")

	DropAndCreateTable(structXml(objName.tableName), objName.tableName, db)

	fileName, err := SearchFile(searchFileName, format)
	CheckError(err, "Error searching file:")

	pathToFile := format + "/" + fileName
	xmlFile, err := os.Open(pathToFile)
	CheckError(err, "Error opening file:")
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	i := 0

	txn, err := db.Begin()
	CheckError(err, "")

	query := pq.CopyIn(objName.tableName, fields...)

	stmt, err := txn.Prepare(query)
	CheckError(err, "")

	for {
		if i == 50000 {
			i = 0

			_, err = stmt.Exec()
			CheckError(err, "")

			err = stmt.Close()
			CheckError(err, "")

			err = txn.Commit()
			CheckError(err, "")

			txn, err = db.Begin()
			CheckError(err, "")

			stmt, err = txn.Prepare(query)
			CheckError(err, "")
		}

		t, _ := decoder.Token()

		if t == nil {
			if i > 0 {
				_, err = stmt.Exec()
				CheckError(err, "")

				err = stmt.Close()
				CheckError(err, "")

				err = txn.Commit()
				CheckError(err, "")
			}
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement := se.Name.Local

			if inElement == objName.elementName {
				err = decoder.DecodeElement(&xmlObject, &se)
				CheckError(err, "Error in decode element:")

				values := ExtractValues(xmlObject)
				_, err = stmt.Exec(values...)
				CheckError(err, "")
				i++
			}
		default:
		}

	}
}

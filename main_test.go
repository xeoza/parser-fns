package main

import (
	"testing"
	"fmt"
	"encoding/xml"
	"os"
	"strings"
	"github.com/lib/pq"

)

func TestFindArea(t *testing.T) {
	c := []interface{}{
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440000",
			"550e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440008",
			"Щелкино",
			"Щелкино", 
			"г",
			4, 
			91,
			1 ,
			"",
			"",
			"",
		},
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440000",
			"550e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440008",
			"Москва",
			"Москва", 
			"г",
			4, 
			91,
			1 ,
			"",
			"",
			"",
		},
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440004",
			"550e8400-e29b-41d4-a716-446655440005",
			"550e8400-e29b-41d4-a716-446655440012",
			"Инкерман",
			"Инкерман",
			"г",
			4,
			92,
			1,
			"",
			"",
			"",
		},
	}
		
	a := []interface{}{ 
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440007",
			"550e8400-e29b-41d4-a716-446655440008",
			"550e8400-e29b-41d4-a716-446655440014",
			"Ленинский",
			"Ленинский",
			"р-н",
			4,
			92,
			1,
			"",
			"",
			"",
		},
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440010",
			"550e8400-e29b-41d4-a716-446655440012",
			"550e8400-e29b-41d4-a716-446655440017",
			"Cевастополь",
			"Севастополь",
			"р-н",
			3,
			92,
			1,
			"",
			"",
			"",
		},
	}
	r := []interface{}{ 
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440013",
			"550e8400-e29b-41d4-a716-446655440014",
			"550e8400-e29b-41d4-a716-446655440015",
			"Башкортостан",
			"Башкортостан",
			"Респ.",
			4,
			92,
			1,
			"",
			"",
			"",
		},
		[]interface{}{
			"550e8400-e29b-41d4-a716-446655440016",
			"550e8400-e29b-41d4-a716-446655440017",
			"550e8400-e29b-41d4-a716-446655440018",
			"Татарстан",
			"Татарстан",
			"Респ.",
			3,
			92,
			1,
			"",
			"",
			"",
		},
	}
	regions := make([]interface{}, 0)
	areas := make([]interface{}, 0)
	cities := make([]interface{}, 0)
	format := "xml"
	db := InitDb()
	xmlObject := &XmlObjectAddrobj{}
	objName := ExtractXmlObjectName(xmlObject)
	fields := ExtractFeilds(xmlObject)
	searchFileName := objName.tableName
	objName.tableName = strings.TrimSuffix(objName.tableName, "_")
	DropAndCreateTable(CreateTableAddrobj(objName.tableName), objName.tableName, db)

	fileName, err := SearchFile(searchFileName, format)
	CheckError(err, "Error searching file:")

	pathToFile := format + "/" + fileName
	xmlFile, err := os.Open(pathToFile)
	CheckError(err, "Error opening file:")
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	query := pq.CopyIn(objName.tableName, fields...)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement := se.Name.Local
			if inElement == objName.elementName {
				err = decoder.DecodeElement(&xmlObject, &se)
				CheckError(err, "Error in decode element:")
				values, flag, ao_level := ExtractValues(xmlObject)
				if flag{
					if ao_level == 1{
						regions = NilToNorm(&regions, &values)
					} else if ao_level == 3{
						areas = NilToNorm(&areas, &values)
					} else if ao_level == 4{
						cities = NilToNorm(&cities, &values)
					}
				}
			break
			}
		default:
		}
	}

	Normalizer(cities, areas, regions)
	//fmt.Println(cities)
	fmt.Printf("\nType C = %T %T %T\n", c, r, a) 
	//fmt.Printf("\nType C = %T %T %T\n", cities, regions, areas) 
	SendToPsql(&cities, db, query)  
}
package main

import (
	"encoding/xml"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"fmt"
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

func ExtractValues(xmlObject interface{}) ([]XmlObjectAddrobj, bool, int64) {
	s := reflect.ValueOf(xmlObject).Elem()
	values := make([]interface{}, s.NumField() - 1)
	flag := false
	var ao_level int64
	if s.Field(9).Int() == 1 && s.Field(7).Int() <= 4{
		ao_level = s.Field(7).Int()
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
				if i == 3 && values[i-1] == ""{
					values[i-1] = nil
				}
				flag = true	
			}
		}
	}
	
	return values, flag, ao_level
}

func Normalizer(cities, areas, regions []XmlObjectAddrobj) []XmlObjectAddrobj {
	result := make([]XmlObjectAddrobj, 0)
	district := make([]objectDistrict, 8)
	district[0].id = 1
	district[0].name = "Центральный федеральный округ"
	district[0].areas = []string{"Белгородская", "Брянская", "Владимирская", "Воронежская", "Ивановская", "Калужская", "Костромская", "Курская", "Липецкая", "Московская", "Орловская", "Рязанская", "Смоленская", "Тамбовская", "Тверская", "Тульская", "Ярославская", "Москва"}
	
	district[1].id = 2 
	district[1].name = "Южный федеральный округ"
	district[1].areas = []string{"Адыгея", "Калмыкия", "Краснодарский", "Астраханская", "Волгоградская", "Ростовская"}
	
	district[2].id = 3 
	district[2].name = "Северо-Западный федеральный округ"
	district[2].areas = []string{"Карелия", "Коми", "Архангельская", "Вологодская", "Калининградская", "Ленинградская", "Мурманская", "Новгородская", "Псковская", "Санкт-Петербург", "Ненецкий"}
	
	district[3].id = 4 
	district[3].name = "Дальневосточный федеральный округ"
	district[3].areas = []string{"Саха /Якутия/", "Камчатский", "Приморский", "Хабаровский", "Амурская", "Магаданская", "Сахалинская", "Еврейская", "Чукотский"}
	
	district[4].id = 5 
	district[4].name = "Сибирский федеральный округ"
	district[4].areas = []string{"Алтай", "Бурятия", "Тыва", "Хакасия", "Алтайский", "Забайкальский", "Красноярский", "Иркутская", "Кемеровская", "Новосибирская", "Омская", "Томская"}
	
	district[5].id = 6 
	district[5].name = "Уральский федеральный округ"
	district[5].areas = []string{"Курганская", "Свердловская", "Тюменская", "Челябинская", "Ханты-Мансийский Автономный округ - Югра", "Ямало-Ненецкий"}
	
	district[6].id = 7 
	district[6].name = "Приволжский федеральный округ"
	district[6].areas = []string{"Башкортостан", "Марий Эл", "Мордовия", "Татарстан", "Удмуртская", "Чувашская", "Кировская", "Нижегородская", "Оренбургская", "Пензенская", "Ульяновская", "Самарская", "Саратовская", "Пермский"}
	
	district[7].id = 8 
	district[7].name = "Северо-Кавказский федеральный округ"
	district[7].areas = []string{"Дагестан", "Ингушетия", "Кабардино-Балкарская", "Карачаево-Черкесская", "Северная Осетия - Алания", "Чеченская", "Ставропольский"}
	
	for i := 0; i < len(cities); i++{
		flag := true
		for j:=0; j < len(areas) && flag; j++{
			if cities[i].PARENTGUID == areas[j].AOGUID{
				result = append(result, cities[i])
				result = append(result, areas[j].OFFNAME)
				for k:=0; k < len(regions); k++{
					if areas[j].PARENTGUID == regions[k].AOGUID{
						result = append(result, regions[k].OFFNAME)
						flag = false
						break
					}
				}

			}
		}
	}
	fmt.Println(cities)

	return result
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
	regions := make([]XmlObjectAddrobj, 0)
	areas := make([]XmlObjectAddrobj, 0)
	cities := make([]XmlObjectAddrobj, 0)
	for {
		/*if i == 50000 {
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
		}*/

		t, _ := decoder.Token()

		/*if t == nil {
			if i > 0 {
				_, err = stmt.Exec()
				CheckError(err, "")

				err = stmt.Close()
				CheckError(err, "")

				err = txn.Commit()
				CheckError(err, "")
			}
			break
		}*/
		switch se := t.(type) {
		case xml.StartElement:
			inElement := se.Name.Local

			if inElement == objName.elementName {
				err = decoder.DecodeElement(&xmlObject, &se)
				CheckError(err, "Error in decode element:")
				values, flag, ao_level := ExtractValues(xmlObject)
				if flag{
					if ao_level == 1{
						regions = append(regions, values)
					} else if ao_level == 3{
						areas = append(areas, values)
					} else if ao_level == 4{
						cities = append(cities, values)
						i++
					}
					//fmt.Println(values[7], regions)
					_, err = stmt.Exec(values...)
					//CheckError(err, "")
				}
			}
		default:
		}
	}
	Normalizer(cities, areas, regions)
}

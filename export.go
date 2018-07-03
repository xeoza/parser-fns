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

func ExtractValues(xmlObject interface{}) ([]interface{}, bool, int64	) {
	s := reflect.ValueOf(xmlObject).Elem()
	values := make([]interface{}, s.NumField() - 1)
	flag := false
	ao_level := s.Field(7).Int()
	if s.Field(9).Int() == 1 && s.Field(7).Int() <= 4{
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
					values[i-1] = "00000000-0000-0000-0000-000000000000"
				}
				flag = true	
			}
		}
	}
	
	return values, flag, ao_level
}


func NilToNorm(a *[]interface{}, b *[]interface{}) []interface{} {
	*a = append(*a, *b)
    return *a
}

func FindAreasForCities(cities, areas *[]interface{}) {
	for _, c := range(*cities){
		switch c.(type) {
		case []interface{}:
			val := reflect.ValueOf(c).Index(2)
			cReg_code := reflect.ValueOf(c).Index(7)
			loop:
			for _, a := range(*areas){
				switch a.(type){
				case []interface{}:
					rVal := reflect.ValueOf(a).Index(1)
					aReg_code := reflect.ValueOf(a).Index(7)
					//cent_stat := reflect.ValueOf(a).Index(12)
					//if aReg_code.Interface().(int64) == cReg_code.Interface().(int64){
						//fmt.Println(aReg_code, cReg_code, reflect.ValueOf(a).Index(4))
					if val.Interface().(string) == rVal.Interface().(string){
						a4 := reflect.ValueOf(a).Index(4)
						a5 := reflect.ValueOf(a).Index(5)
						aSum := reflect.ValueOf(a4.Interface().(string)+ " " + a5.Interface().(string))
						reflect.ValueOf(c).Index(9).Set(aSum)
						reflect.ValueOf(c).Index(10).Set(reflect.ValueOf(a).Index(10))
						reflect.ValueOf(c).Index(11).Set(reflect.ValueOf(a).Index(11))
						break loop
					} else if aReg_code.Interface().(int64) == cReg_code.Interface().(int64){
{						reflect.ValueOf(c).Index(10).Set(reflect.ValueOf(a).Index(10))
						reflect.ValueOf(c).Index(11).Set(reflect.ValueOf(a).Index(11))
						break loop
					}
				}
				case interface{}:
					fmt.Println("naebatelstvo")
				}
			}
		case interface{}:
			fmt.Println("naebatelstvo")	
		
		}
	}

}

func FindRegionsForAreas(areas, regions *[]interface{}) {
	for _, a := range(*areas){
		switch a.(type) {
		case []interface{}:
			val := reflect.ValueOf(a).Index(2)
			loop:
			for _, r := range(*regions){
				switch r.(type){
				case []interface{}:
					rVal := reflect.ValueOf(r).Index(1)
					if val.Interface().(string) == rVal.Interface().(string){
						a4 := reflect.ValueOf(r).Index(4)
						a5 := reflect.ValueOf(r).Index(5)
						aSum := reflect.ValueOf(a4.Interface().(string)+ " " + a5.Interface().(string))
						reflect.ValueOf(a).Index(10).Set(aSum)
						reflect.ValueOf(a).Index(11).Set(reflect.ValueOf(r).Index(11))
						break loop
					}
				case interface{}:
					fmt.Println("naebatelstvo")
				}
			}
		case interface{}:
			fmt.Println("naebatelstvo")	
		
		}
	}
}

func FindDistrictForRegions(regions *[]interface{}){

	district := make([]objectDistrict, 8)
	district[0].id = 1
	district[0].name = "Центральный федеральный округ"
	district[0].areas = []string{"Белгородская", "Брянская", "Владимирская", "Воронежская", "Ивановская", "Калужская", "Костромская", "Курская", "Липецкая", "Московская", "Орловская", "Рязанская", "Смоленская", "Тамбовская", "Тверская", "Тульская", "Ярославская", "Москва"}
	
	district[1].id = 2 
	district[1].name = "Южный федеральный округ"
	district[1].areas = []string{"Адыгея", "Калмыкия", "Краснодарский", "Астраханская", "Волгоградская", "Ростовская", "Крым"}
	
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
	district[6].areas = []string{"Башкортостан", "Марий Эл", "Мордовия", "Татарстан", "Удмуртская", "Чувашия", "Кировская", "Нижегородская", "Оренбургская", "Пензенская", "Ульяновская", "Самарская", "Саратовская", "Пермский"}
	
	district[7].id = 8 
	district[7].name = "Северо-Кавказский федеральный округ"
	district[7].areas = []string{"Дагестан", "Ингушетия", "Кабардино-Балкарская", "Карачаево-Черкесская", "Северная Осетия - Алания", "Чеченская", "Ставропольский"}
	

	for _, r := range(*regions){
		switch r.(type) {
		case []interface{}:
			val := reflect.ValueOf(r).Index(4)
			shortVal := reflect.ValueOf(r).Index(5)
			loop:
			for i:=0; i < len(district); i++{
				for _, str := range(district[i].areas){
					if val.Interface().(string) == str || shortVal.Interface().(string) == str{
						reflect.ValueOf(r).Index(11).Set(reflect.ValueOf(district[i].name))
						break loop
					} 
				}
			}
		case interface{}:
			fmt.Println("naebatelstvo")
		}
	}
}
/*func SearhProblemCities(cities, areas *[]interface{}){
	for _, c := range(*cities){
			switch c.(type) {
			case []interface{}:
				val := reflect.ValueOf(c).Index(12)
				loop:
				if val.Interface().(string) == "1"{

				}else if cent_stat.Interface().(string) == "2"{

				}else if cent_stat.Interface().(string) == "3"{

				}
			case interface{}:
				fmt.Println("naebatelstvo")
			}
		}
}*/

func Normalizer(cities, areas, regions []interface{}) {
	FindDistrictForRegions(&regions)
	FindRegionsForAreas(&areas, &regions)
	FindAreasForCities(&cities, &areas)
}

func SendToPsql(cities *[]interface{}, db *sqlx.DB, query string){
	txn, err := db.Begin()
	CheckError(err, "")
	stmt, err := txn.Prepare(query)
	CheckError(err, "")
	for i, city := range (*cities){
		switch city.(type) {
		case []interface{}:
			if i == 100{
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
			_,err = stmt.Exec(city.([]interface{})...)
			CheckError(err, "")

		case interface{}:
			fmt.Println("naebatelstvo")
		}
	}
	err = stmt.Close()
	CheckError(err, "")

	err = txn.Commit()
	CheckError(err, "")
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

	query := pq.CopyIn(objName.tableName, fields...)

	CheckError(err, "")
	regions := make([]interface{}, 0)
	areas := make([]interface{}, 0)
	cities := make([]interface{}, 0)
	for {
		t, _ := decoder.Token()

		if t == nil { break	}
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
						i++
					}
				}
			}
		default:
		}
	}

	Normalizer(cities, areas, regions)
	fmt.Println(areas)
	SendToPsql(&cities, db, query)
}
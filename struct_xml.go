package main

import (
	"encoding/xml"
)

type objectDistrict struct{
    id int
    name string
    areas []string
}

type xmlObjectName struct {
	tableName   string
	elementName string
}

type XmlObjectAddrobj struct {
	XMLName    xml.Name `xml:"Object" db:"as_addrobj"`
	AOID       string   `xml:"AOID,attr" db:"ao_id"`
	AOGUID     string   `xml:"AOGUID,attr" db:"ao_guid"`
	PARENTGUID string   `xml:"PARENTGUID,attr,omitempty" db:"parent_guid"`
	FORMALNAME string   `xml:"FORMALNAME,attr" db:"formal_name"`
	OFFNAME    string   `xml:"OFFNAME,attr,omitempty" db:"off_name"`
	SHORTNAME  string   `xml:"SHORTNAME,attr" db:"short_name"`
	//AREA       string   `xml:"AREA,attr" db:"area"`
    //REGION     string   `xml:"REGION,attr" db:"region"`
    //DISTRICT   string   `xml:"DISTRICT,attr" db:"reion"`
	AOLEVEL    int      `xml:"AOLEVEL,attr" db:"ao_level"`
	REGIONCODE int      `xml:"REGIONCODE,attr" db:"region_code"`
	//OKATO      string   `xml:"OKATO,attr,omitempty" db:"okato"`
	//OKTMO      string   `xml:"OKTMO,attr,omitempty" db:"oktmo"`
	LIVESTATUS int      `xml:"LIVESTATUS,attr" db:"live_status"`
}

func CreateTableAddrobj(tableName string) string {
	//AREA       string   `xml:"AREA,attr" db:"area"`
    //REGION     string   `xml:"REGION,attr" db:"region"`
    //DISTRICT   string   `xml:"DISTRICT,attr" db:"reion"`
    return `CREATE TABLE ` + tableName + ` (
		ao_id UUID NOT NULL,
	    ao_guid UUID NOT NULL,
		parent_guid UUID,
	    formal_name VARCHAR(120) NOT NULL,
		off_name VARCHAR(120),
		short_name VARCHAR(10) NOT NULL,
		ao_level INT NOT NULL,
		region_code INT NOT NULL,
		live_status INT NOT NULL,
		PRIMARY KEY (ao_id));`
}

type XmlObjectSocrbase struct {
	XMLName  xml.Name `xml:"AddressObjectType" db:"as_socrbase"`
	LEVEL    int      `xml:"LEVEL,attr" db:"level"`
	SCNAME   string   `xml:"SCNAME,attr" db:"sc_name"`
	SOCRNAME string   `xml:"SOCRNAME,attr" db:"socr_name"`
	KOD_T_ST string   `xml:"KOD_T_ST,attr" db:"kod_t_st"`
}

func CreateTableSocrbase(tableName string) string {
	return `CREATE TABLE ` + tableName + ` (
    level INT NOT NULL,
    sc_name VARCHAR(20),
    socr_name VARCHAR(60),
    kod_t_st INT UNIQUE NOT NULL,
		PRIMARY KEY (kod_t_st));`
}
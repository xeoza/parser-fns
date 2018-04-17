package main

import (
	"encoding/xml"
)

type xmlObjectName struct {
	tableName   string
	elementName string
}

type XmlObjectAddrobj struct {
	XMLName    xml.Name `xml:"Object" db:"as_addrobj"`
	AOGUID     string   `xml:"AOGUID,attr" db:"ao_guid"`
	FORMALNAME string   `xml:"FORMALNAME,attr" db:"formal_name"`
	REGIONCODE int      `xml:"REGIONCODE,attr" db:"region_code"`
	AUTOCODE   int      `xml:"AUTOCODE,attr" db:"auto_code"`
	AREACODE   int      `xml:"AREACODE,attr" db:"area_code"`
	CITYCODE   int      `xml:"CITYCODE,attr" db:"city_code"`
	CTARCODE   int      `xml:"CTARCODE,attr" db:"ctar_code"`
	PLACECODE  int      `xml:"PLACECODE,attr" db:"place_code"`
	STREETCODE int      `xml:"STREETCODE,attr,omitempty" db:"street_code"`
	EXTRCODE   int      `xml:"EXTRCODE,attr" db:"extr_code"`
	SEXTCODE   int      `xml:"SEXTCODE,attr" db:"sext_code"`
	OFFNAME    *string  `xml:"OFFNAME,attr,omitempty" db:"off_name"`
	POSTALCODE *string  `xml:"POSTALCODE,attr,omitempty" db:"postal_code"`
	IFNSFL     int      `xml:"IFNSFL,attr,omitempty" db:"ifns_fl"`
	TERRIFNSFL int      `xml:"TERRIFNSFL,attr,omitempty" db:"terr_ifns_fl"`
	IFNSUL     int      `xml:"IFNSUL,attr,omitempty" db:"ifns_ul"`
	TERRIFNSUL int      `xml:"TERRIFNSUL,attr,omitempty" db:"terr_ifns_ul"`
	OKATO      *string  `xml:"OKATO,attr,omitempty" db:"okato"`
	OKTMO      *string  `xml:"OKTMO,attr,omitempty" db:"oktmo"`
	UPDATEDATE string   `xml:"UPDATEDATE,attr" db:"update_date"`
	SHORTNAME  string   `xml:"SHORTNAME,attr" db:"short_name"`
	AOLEVEL    int      `xml:"AOLEVEL,attr" db:"ao_level"`
	PARENTGUID *string  `xml:"PARENTGUID,attr,omitempty" db:"parent_guid"`
	AOID       string   `xml:"AOID,attr" db:"ao_id"`
	PREVID     *string  `xml:"PREVID,attr,omitempty" db:"prev_id"`
	NEXTID     *string  `xml:"NEXTID,attr,omitempty" db:"next_id"`
	CODE       *string  `xml:"CODE,attr,omitempty" db:"code"`
	PLAINCODE  *string  `xml:"PLAINCODE,attr,omitempty" db:"plain_code"`
	ACTSTATUS  int      `xml:"ACTSTATUS,attr" db:"act_status"`
	CENTSTATUS int      `xml:"CENTSTATUS,attr" db:"cent_status"`
	OPERSTATUS int      `xml:"OPERSTATUS,attr" db:"oper_status"`
	CURRSTATUS int      `xml:"CURRSTATUS,attr" db:"curr_status"`
	STARTDATE  string   `xml:"STARTDATE,attr" db:"start_date"`
	ENDDATE    string   `xml:"ENDDATE,attr" db:"end_date"`
	NORMDOC    *string  `xml:"NORMDOC,attr,omitempty" db:"norm_doc"`
	LIVESTATUS int      `xml:"LIVESTATUS,attr" db:"live_status"`
}

func CreateTableAddrobj(tableName string) string {
	return `CREATE TABLE ` + tableName + ` (
    ao_guid UUID NOT NULL,
    formal_name VARCHAR(120) NOT NULL,
		region_code INT NOT NULL,
		auto_code INT NOT NULL,
		area_code INT NOT NULL,
		city_code INT NOT NULL,
		ctar_code INT NOT NULL,
		place_code INT NOT NULL,
		street_code INT,
		extr_code INT NOT NULL,
		sext_code INT NOT NULL,
		off_name VARCHAR(120),
		postal_code VARCHAR(6),
		ifns_fl INT,
		terr_ifns_fl INT,
		ifns_ul INT,
		terr_ifns_ul INT,
		okato VARCHAR(11),
		oktmo VARCHAR(11),
		update_date TIMESTAMP NOT NULL,
		short_name VARCHAR(10) NOT NULL,
		ao_level INT NOT NULL,
		parent_guid UUID,
		ao_id UUID NOT NULL,
		prev_id UUID,
		next_id UUID,
		code VARCHAR(17),
		plain_code VARCHAR(15),
		act_status INT NOT NULL,
		cent_status INT NOT NULL,
		oper_status INT NOT NULL,
		curr_status INT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		norm_doc UUID,
		live_status INT NOT NULL,
		PRIMARY KEY (ao_id));`
}

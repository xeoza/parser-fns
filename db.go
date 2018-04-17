package main

import "github.com/jmoiron/sqlx"

func InitDb() *sqlx.DB {
	db, err := sqlx.Open("postgres", "user=admin dbname=intelliada_db password=tnved1357tup sslmode=disable")
	CheckError(err, "sqlx.Open failed")

	return db
}

func DropAndCreateTable(structXml string, tableName string, db *sqlx.DB) {
	rows, err := db.Queryx("SELECT to_regclass('" + tableName + "');")
	CheckError(err, "") 
	defer rows.Close()

	rowsCount := 0
	for rows.Next() {
		rowsCount++
	}

	if rowsCount > 0 {
		_, err = db.Exec("DROP TABLE IF EXISTS " + tableName + ";")
		CheckError(err, "") 
	}

	_, err = db.Exec(structXml)
	CheckError(err, "") 

}

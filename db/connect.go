package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

	var host = os.Getenv("DB_HOSTNAME")
	var port = os.Getenv("DB_PORT_NUMBER")
	var user = os.Getenv("DB_USERNAME")
	var password = os.Getenv("DB_PASSWORD")
	var dbName = os.Getenv("DB_NAME")


func Connect() *sql.DB {
	db := openDbConnection()
	_, err := db.Exec(
		`CREATE TABLE CLIENTS_PURCHASE_INFORMATIONS(
			ID serial PRIMARY KEY
		   ,CPF Varchar(14) NOT NULL
		   ,PRIVATE boolean
		   ,INCOMPLETE boolean
		   ,LAST_PURCHASE_DATE DATE
		   ,AVERAGE_TICKET decimal(12,2)
		   ,LAST_PURCHASE_TICKET decimal(12,2)
		   ,MOST_FREQUENT_STORE Varchar(18)
		   ,LAST_PURCHASE_STORE Varchar(18)
		   )`)

	if err != nil {
		log.Print(err)
	}
	return db
}

func openDbConnection() *sql.DB {
	db, _ := sql.Open("postgres", getPostgresConnectionString(""))
	db.Exec("CREATE DATABASE " + dbName)
	db.Close()
	db, _ = sql.Open("postgres", getPostgresConnectionString(dbName))
	return db
}

func getPostgresConnectionString(dbName string) string {
	var connectionString string
	if dbName != "" {
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbName)

	} else {
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
			host, port, user, password)
	}

	return connectionString
}

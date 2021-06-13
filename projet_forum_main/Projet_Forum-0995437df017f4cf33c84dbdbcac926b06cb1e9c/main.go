package main

import (
	"database/sql"
	"log"
	"net/http"

	connexion "./server"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Host = "localhost"
	Port = "1500"
)

func main() {

	database := initdatabase("./Forum.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Utilisateur (ID_user INTEGER PRIMARY KEY ASC AUTOINCREMENT,Nom STRING NOT NULL,PRENOM STRING NOT NULL,MAIL [STRING UNIQUENOT] NOT NULL	UNIQUE,PASSWORD STRING NOT NULL,User_name STRING NOT NULL UNIQUE,Birth_Date DATE);")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS post (ID_user INTEGER PRIMARY KEY ASC AUTOINCREMENT, titre TEXT NOT NULL, message TEXT NOT NULL, informatique TEXT NOT NULL, jv TEXT NOT NULL, art TEXT NOT NULL, sport TEXT NOT NULL, ht TEXT NOT NULL, cuisine TEXT NOT NULL);")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS reponse (ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, ID_post TEXT, commente TEXT NOT NULL);")
	statement.Exec()

	connexion.Login()
	err := http.ListenAndServe(Host+":"+Port, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
		return
	}
}
func initdatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

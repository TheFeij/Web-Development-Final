package main

import (
	"Messenger/api"
	"Messenger/db"
	"database/sql"
	"log"
)

func main() {
	db.Init()
	db := db.GetDB()
	DB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(DB)

	server := api.NewServer(db)
	err = server.Start("localhost:8080")
	if err != nil {
		log.Fatalln("Cannot start the server ", err.Error())
		return
	}
}

package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spiccoli/e-commerceAPI/cmd/api"
	"github.com/spiccoli/e-commerceAPI/config"
	"github.com/spiccoli/e-commerceAPI/db"
	"log"
)

func main() {
	db, dbErr := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPasswd,
		Addr:                 config.Envs.DBAddr,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	log.Println(fmt.Sprintf("Server is running on %s", ":8080"))
	if serverErr := server.Run(); serverErr != nil {
		log.Fatal(serverErr)
	}
}
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
}

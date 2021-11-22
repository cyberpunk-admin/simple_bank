package main

import (
	"database/sql"
	"github.com/simplebank/api"
	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/util"
	"log"

	_ "github.com/lib/pq"
)


func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriven, config.DBSource)
	if err != nil {
		log.Fatalln("cannot conneted to database: ", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}

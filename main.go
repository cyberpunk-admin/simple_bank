package main

import (
	"github.com/simplebank/pb"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/simplebank/api"
	db "github.com/simplebank/db/sqlc"

	"github.com/simplebank/gapi"
	"github.com/simplebank/util"
	"google.golang.org/grpc"
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
	runGrpcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can't create listener")
	}

	log.Println("start gRPC serice at ", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("can't start grpc server")
	}
}
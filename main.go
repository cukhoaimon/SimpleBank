package main

import (
	"database/sql"
	"github.com/cukhoaimon/SimpleBank/gapi"
	"github.com/cukhoaimon/SimpleBank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"github.com/cukhoaimon/SimpleBank/api"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {

		log.Fatal(err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("The open connection to database process was encountered an error", err)
	}

	store := db.NewStore(conn)
	runGRPCServer(store, config)
}

func runGinServer(store db.Store, config utils.Config) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot create server")
	}

	if err = server.Start(config.HttpServerAddress); err != nil {
		log.Fatal("Cannot start server")
	}
}

func runGRPCServer(store db.Store, config utils.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot create gRPC server")
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, server)
	// allow client to know what RPCs currently available in server
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("Cannot create tcp-listener for gRPC server")
	}

	log.Printf("start gRPC server at: %s", listener.Addr().String())
	if err = gRPCServer.Serve(listener); err != nil {
		log.Fatal("Cannot serve gRPC server")
	}
}

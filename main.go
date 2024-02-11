package main

import (
	"context"
	"database/sql"
	"github.com/cukhoaimon/SimpleBank/gapi"
	"github.com/cukhoaimon/SimpleBank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"

	"github.com/cukhoaimon/SimpleBank/api"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/utils"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
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

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	go runGatewayServer(store, config)
	runGRPCServer(store, config)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("fail to create migration instance: ", err)
	}

	if err = migration.Up(); err != nil {
		log.Fatal("fail to run migrate up: ", err)
	}

	log.Println("migration is successfully")
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
		log.Fatalf("Cannot create gRPC server: %s", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, server)
	// allow client to know what RPCs currently available in server
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("Cannot create tcp-listener for gRPC server: %s", err)
	}

	log.Printf("start gRPC server at: %s", listener.Addr().String())
	if err = gRPCServer.Serve(listener); err != nil {
		log.Fatalf("Cannot serve gRPC server: %s", err)
	}
}

func runGatewayServer(store db.Store, config utils.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatalf("Cannot create gRPC server: %s", err)
	}

	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatalf("Cannot register handler server: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatalf("Cannot create tcp-listener for gateway server: %s", err)
	}

	log.Printf("start HTTP gateway server at: %s", listener.Addr().String())
	if err = http.Serve(listener, mux); err != nil {
		log.Fatalf("cannot HTTP gateway server: %s", err)
	}
}

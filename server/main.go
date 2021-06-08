package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"api/app/pb"
	"api/config"
	"api/controllers"
	"api/db"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	mongoDb  *mongo.Collection
	mongoCtx context.Context
	err      error
)

func init() {
	err = config.LoadEnv()
	if err != nil {
		log.Fatalf("Could not load env files: %s", err.Error())
	}

	mongoDb, mongoCtx, err = db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to mongoDb: %s", err.Error())
	}
	fmt.Printf("Connected to mongodb: %s\n", db.ConnectionString)
}

func main() {
	apiPort := os.Getenv("API-PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", apiPort))
	if err != nil {
		log.Fatalf("Could not connect on port :%s\n", apiPort)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	nflApi := controllers.NflApiServiceServer{
		Db:  mongoDb,
		Ctx: mongoCtx,
	}

	pb.RegisterNflApiServiceServer(grpcServer, &nflApi)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	fmt.Printf("Server succesfully started on port :%s\n", apiPort)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Stopping the server...")
	grpcServer.Stop()
	listener.Close()
	fmt.Println("Done.")
}

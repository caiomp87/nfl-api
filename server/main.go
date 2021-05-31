package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"api/app/pb"
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
	mongoDb, mongoCtx, err = db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to mongoDb: %s", err.Error())
	}
	fmt.Println("Connected to mongodb!")
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Could not connect on port :50051", err.Error())
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

	fmt.Println("Server succesfully started on port :50051")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Stopping the server...")
	grpcServer.Stop()
	listener.Close()
	fmt.Println("Done.")
}

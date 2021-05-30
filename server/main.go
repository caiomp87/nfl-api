package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"api/app/pb"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type NflApiServiceServer struct {
	pb.UnimplementedNflApiServiceServer
	Db  *mongo.Collection
	Ctx context.Context
}

func init() {

}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Could not connect on port :50051", err.Error())
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pb.RegisterNflApiServiceServer(grpcServer, &NflApiServiceServer{})

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

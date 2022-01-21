package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	//Use gRPC as alternative to REST for client-server interface
	//gRPC sends data over remote procedure call as bytecode, therefore much less data sent & more efficient
	//Effectively calling a function on client side will call the same function on server side in remote procedure call
	//To ensure data aligns, use protobuf to marshal & match data to function signatures
	//Protobuf uses an interface description language to serialize structured data to send over network

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}

}

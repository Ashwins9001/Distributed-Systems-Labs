package main

import (
	"strconv"
	"fmt"
	"context"
	"log"
	"os"
	"net"
	"bufio"

	"google.golang.org/grpc"

	pb "simple-grpc/proto"
)

//Array to store results in for writing data later to a file
var results []string

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedOperationServer
}

//Helper function to perform op
func createMessage(op string, num1 string, num2 string) string {
	n1, err := strconv.Atoi(num1)
	
	if err != nil {
		fmt.Println("Error reading first number ", err)
	}

	n2, err := strconv.Atoi(num2)

	if err != nil {
		fmt.Println("Error reading second number ", err)
	}
	
	var res int
	var message string
	var res_str string

	message = "Invalid operation selected, try again"

	if op == "add" {
		res =  n1 + n2
	} else if op == "sub" {
		res = n1 - n2
	} else if op == "mul" {
		res = n1 * n2
	} else if op == "div" {
		res = n1 / n2
	} else {
		return message
	}

	res_str = strconv.Itoa(res)
	
	//out_str := res_str + "\n"

	results = append(results, res_str)

	message = "Completed operation " + op + " to get result " + res_str

	return message
}

// SayHello implements helloworld.GreeterServer
func (s *server) DoOp(ctx context.Context, in *pb.OpRequest) (*pb.OpReply, error) {
	log.Printf("Received: %v", in.GetOp(), in.GetNum1(), in.GetNum2())
	msg := createMessage(in.GetOp(), in.GetNum1(), in.GetNum2())
	
	//Write to file on each call
	path, pathErr := os.Getwd()
	if pathErr != nil {
		log.Println(pathErr)
	}

	path += "/output.txt"

	output, _ := os.Create(path)
	for k := 0; k < len(results); k++ {
		output.WriteString(results[k] + "\n")
	}

	return &pb.OpReply{Message: msg}, nil
}

func main() {
	//Retrieve port information for server start-up
	
	//Get port information
	portPath, portPathError := os.Getwd()

	if portPathError != nil {
		fmt.Println("Error finding current path ", portPathError)
	}

	portPath += "/port.txt"
	portFile, portFileError := os.Open(portPath)

	if portFileError != nil {
		fmt.Println("Error opening port file ", portFileError)
	}

	address := ":"
	portScanner := bufio.NewScanner(portFile)
	port := ""
	for portScanner.Scan() {
		port += portScanner.Text()
	}

	address += port

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Successfully started server")

	//Let server run and wait for requests
	s := grpc.NewServer()
	
	pb.RegisterOperationServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

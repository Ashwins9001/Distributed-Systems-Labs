package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "simple-grpc/proto"
)

func main() {
	
	//Import port information from file
	portPath, portPathError := os.Getwd()

	if portPathError != nil {
		fmt.Println("Error finding current path ", portPathError)
	}

	portPath += "/port.txt"
	portFile, portFileError := os.Open(portPath)

	if portFileError != nil {
		fmt.Println("Error opening port file ", portFileError)
	}

	address := "localhost:"

	portScanner := bufio.NewScanner(portFile)
	port := ""
	for portScanner.Scan() {
		port += portScanner.Text()
	}

	address += port


	// Set up a connection to the server.
	conn, connErr := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if connErr != nil {
		log.Fatalf("Could not connect: %v", connErr)
	}
	defer conn.Close()
	c := pb.NewOperationClient(conn)

	path, pathErr := os.Getwd()
	if pathErr != nil {
		log.Println(pathErr)
	}
	
	path += "/input.txt"

	//Parse input file and for each line, call server with a request to process a math op
	var file, fileError = os.Open(path)
	if fileError != nil {
		fmt.Println("Error opening file", fileError)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineText := scanner.Text()
		parseLine := strings.Fields(lineText)

		operation := parseLine[0]
		firstnum := parseLine[1]
		secondnum := parseLine[2]

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.DoOp(ctx, &pb.OpRequest{Op: operation, Num1: firstnum, Num2: secondnum})
		if err != nil {
			log.Fatalf("could not send message: %v", err)
		}
		log.Printf("Operation sent: %s", r.GetMessage())
	}
}

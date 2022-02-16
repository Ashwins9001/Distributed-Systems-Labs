package main

import (
	"context"
	"log"
	"time"
	"bufio"
	"fmt"
	"strings"
	"os"

	"google.golang.org/grpc"

	pb "bank/proto"
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

	//Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBankOperationClient(conn)

	//Open input file for bank commands to server to be read in
	path, pathErr := os.Getwd()
	if pathErr != nil {
		log.Println("Error finding input file ", pathErr)
	}

	path += "/input.txt"

	file, fileErr := os.Open(path)
	if fileErr != nil {
		fmt.Println("Error opening input file ", fileErr)
	}

	scanner := bufio.NewScanner(file)

	//Scan each line, parse into tokens & send request to server
	for scanner.Scan() {
		lineText := scanner.Text()
		parseLine := strings.Fields(lineText)

		operation := parseLine[0]
		accNum := parseLine[1]
		amount := parseLine[2]


		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.DoBankOp(ctx, &pb.BankRequest{Acc: accNum, Op: operation, Amt: amount})
		if err != nil {
			log.Fatalf("Could not send request: %v", err)
		}
		log.Printf("Operation successfully sent. \n Received reply: ", r.GetMessage())
	}
}

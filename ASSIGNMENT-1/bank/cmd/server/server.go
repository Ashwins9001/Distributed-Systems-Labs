package main

import (
	"context"
	"log"
	"net"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"strconv"
	"bufio"

	"google.golang.org/grpc"

	pb "bank/proto"
)

//Struct to match format of accounts data for marshalling, unmarshalling into JSON 
type BankAccount struct {
	User string `json:"Name"`
	Account int `json:"AccountID"`
	Amount float64 `json:"Balance"`
}

//Struct to implement ops defined in protobuf file
type server struct {
	pb.UnimplementedBankOperationServer
}

//Array to store all accounts in
var accounts []BankAccount


//Helper function to return message to client to indicate account update made
func createMessage(op string, acc string, amt string) string {
	var message string

	//Parse strings to integers/floats & perform op
	accNum, err := strconv.Atoi(acc)
	if err != nil {
		fmt.Println("Error parsing account number ", err)
		message = "Account not found, try again"
		return message
	}

	amtNum, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		fmt.Println("Error parsing amount ", err)
		message = "Invalid amount, try again"
		return message
	}

	for i := 0; i < len(accounts); i++ {
		if accounts[i].Account == accNum {
			if op == "deposit" {
				accounts[i].Amount += amtNum
			} else if op == "withdraw" {
				accounts[i].Amount -= amtNum
			} else if op == "interest" {
				interestRate := 1.00
				amtNum /= 100
				interestRate += amtNum
				accounts[i].Amount *= interestRate
			} else {
				message = "Invalid operation, try again"
				return message
			}
			amtStr := fmt.Sprintf("%f", accounts[i].Amount)
			accStr := strconv.Itoa(accounts[i].Account)
			message = "Completed operation for account: " + accStr + " for new balance: " + amtStr + "\n"
			return message
		}
	}
	message = "Unknown error"
	return message
}

//Perform operation on account given request from client
func (s *server) DoBankOp(ctx context.Context, in *pb.BankRequest) (*pb.BankReply, error) {
	log.Printf("Received: %v", in.GetAcc(), in.GetOp(), in.GetAmt())

	//Formulate message server returns using helper function to confirm operation to client
	msg := createMessage(in.GetOp(), in.GetAcc(), in.GetAmt())
	
	//Write updated accounts.json file
	path, pathErr := os.Getwd()
	if pathErr != nil {
		log.Println(pathErr)
	}

	path += "/accountsUpdated.json"

	mod_acc_json, _ := json.Marshal(accounts)
	output, _ := os.Create(path)
	output.WriteString(string(mod_acc_json))

	//Return message from server via callback
	return &pb.BankReply{Message: msg}, nil
}

func main() {
	//Retrieve port information for server start-up
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

	//Perform setup for reading accounts.json and generating array of accounts
	path, pathErr := os.Getwd()
	if pathErr != nil {
		fmt.Println("Error opening accounts.json file ", pathErr)
	}

	path += "/accounts.json"

	content, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		fmt.Println("Error reading accounts.json file ", readErr)
	}

	accounts_str := string(content)

	//Unmarshal data into accounts struct for storage, updates and later marshalling
	unmarshalErr := json.Unmarshal([]byte(accounts_str), &accounts)

	if unmarshalErr != nil {
		fmt.Println("Error unmarshalling accounts.json file ", unmarshalErr)
	}

	//Configure port to launch server on
	lis, listenErr := net.Listen("tcp", address)
	
	if listenErr != nil {
		log.Fatalf("failed to listen: %v", listenErr)
	}

	fmt.Println("Successfully started server")

	//Launch server and begin waiting for callbacks from client interaction
	s := grpc.NewServer()
	pb.RegisterBankOperationServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

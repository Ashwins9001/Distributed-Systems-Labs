// ASHWIN SINGH
// CS4435
// ASSIGNMENT 2: IMPLEMENTING LAMPORT CLOCKS FOR TOTAL-ORDERING REPLICATION
// MARCH 11, 2022

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"strconv"
	"math"
	"bufio"
	"encoding/json"
	"io/ioutil"

	"github.com/hashicorp/consul/api"
	pb "go-peer/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// struct to parse accounts data from JSON to data structure & back
type BankAccount struct {
	User string `json:"Name"`
	Account int `json:"AccountID"`
	Amount float64 `json:"Balance"`
}

// struct to store data on Node objects
type node struct {
	// Self information
	Name          string             //Service discovery (Consul) name
	Addr          string             //Service discovery (Consul) address
	NodeId        int                //ID of node
	EventSeen     []float32          //Vector clock of events seen by node
	RecvMessages  map[float32]string //Map with key: timestamp, value: message
	NodeTimestamp float32            //Timestamp of node's most recent operation
	File string
	CompletedFiles int //Flag to track number of completed transactions from other nodes
	MessageSent int //Flag to prevent sending duplicate messaages to current node
	AllMessages string //String containing all parsed messages, ready for tokenization and logic-handling for execution
	Accounts []BankAccount

	// Consul related variables to establish connection
	SDAddress string
	SDKV      api.KV

	// used to make requests to all nodes on network
	Clients map[string]pb.BankServiceClient
}

// function marshals and writes updated accounts to unique node file
func (n *node) WriteUpdatedAccounts() {
	path, _ := os.Getwd()
	path += "/accountsUpdated" + n.Name + ".json"
	acc_json, err := json.Marshal(n.Accounts)
	if err != nil {
		log.Println("Error writing to updated accounts.json file")
	}
	output, _ := os.Create(path)
	output.WriteString(string(acc_json))
}

// function handles parsing the command string and processing responding bank operations
func (n *node) ParseCommands() {
	log.Println("Received all messages from other nodes, begin parsing. Commands received at node displayed below:")
	tokens := strings.Fields(n.AllMessages)
	for i := 0; i < len(tokens); i++ {
		if i != 0 {
			// messages occur every three tokens 
			if (i+1) % 3 == 0  {
				log.Println("Command: ", tokens[i-2], tokens[i-1], tokens[i])
				accNum, _ := strconv.Atoi(tokens[i-1])
				amtNum, _ := strconv.ParseFloat(tokens[i], 64)
				// apply transaction logic to update accounts
				for j:=0; j < len(n.Accounts); j++ {
					if n.Accounts[j].Account == accNum {
						if tokens[i-2] ==  "deposit" {
							n.Accounts[j].Amount += amtNum
						} else if tokens[i-2] == "withdraw" {
							n.Accounts[j].Amount -= amtNum
						} else if tokens[i-2] == "interest" {
							interest := 1.00
							amtNum /= 100
							interest += amtNum
							n.Accounts[j].Amount *= interest
						} else {
							log.Println("Invalid transaction command")
						}
					}
				}
			}
		}
	}
}

// function handles message delivery prep including sorting by order received
func (n *node) deliverMessages() {
	min := float32(1000000)
	for _, ts := range n.EventSeen {
		if ts < min {
			min = ts
		}
	}

	// ensure messages are executed in order by waiting for each received message to have timestamp greater than node's timestamp
	// implying that node must receive reply from other nodes before delivering message
	for ts, msg := range n.RecvMessages {
		if ts <= min {
			n.AllMessages = n.AllMessages + msg
			delete(n.RecvMessages, ts)
		}
	}
	
	// in case where message received from each other node in system, sort messages by logical time and execute
	// by design, each node only sends up to two messages, one message containing all transactions and another containing string "done"
	// hence once indication given from each node that a "done" is received, increment n.CompletedFiles until it's equal to the number of other nodes
	// at that point, it is guaranteed that all nodes done processing hence can sort and execute rest of transactions 
	// processing logic similar to acknowledgements, where message delivery can begin once reply received from every other node
	if n.CompletedFiles == (len(n.EventSeen) - 1) {
		for i := 0; i < n.CompletedFiles; i++ {
			min := float32(100000)
			//find message with smallest timestamp and execute, essentially sorting and executing messages in-order
			for ts, _ := range n.RecvMessages {
				if ts < min {
					min = ts
				}
			}
			n.AllMessages = n.AllMessages + n.RecvMessages[min]
			delete(n.RecvMessages, min)
		}
		// once all messages delivered, call functions to parse messages and update accounts
		n.ParseCommands()
		n.WriteUpdatedAccounts()
	}
}

// function handles vector clock updates and message delivery prompts upon receiving a request from another node 
func (n *node) DoTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionReply, error) {
	id_str := strconv.Itoa(int(in.GetId()))

	// case where message received indicates that sending node is done parsing transactions
	if in.GetMessage() == "done" {
		n.CompletedFiles = n.CompletedFiles + 1
	} else {
		// only sending timestamps of sender node around, hence in vector clocks update current and sender node's indices only
		// case where sending node provided transactions
		// update timestamp for current node to max of existing timestamp and timestamp of sender node
		tmp := math.Max(float64(in.GetTimestamp()), float64(n.EventSeen[n.NodeId-1])) + 1
		n.EventSeen[n.NodeId-1] = float32(tmp)
		// update timestamp at sender node to timestamp provided, essentially updating current node's recording of number of events seen in sender node 
		n.EventSeen[int(in.GetId())-1] = in.GetTimestamp()
		n.RecvMessages[in.GetTimestamp()] = in.GetMessage()
		log.Println("Received message containing all transactions from: " + id_str)
		log.Println("Message content: " + in.GetMessage())
		log.Println("\nUpdated node vector clock: ", n.EventSeen)
	}
	
	// check if messages can be delivered and deliver them if so
	n.deliverMessages()

	return &pb.TransactionReply{Message: n.Name + " received message"}, nil
}

// start listening/service.
func (n *node) StartListening() {

	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server started listening")

	_n := grpc.NewServer()

	pb.RegisterBankServiceServer(_n, n)
	
	// Register reflection service on gRPC server.
	reflection.Register(_n)

	// start listening
	if err := _n.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Register self with the service discovery module.
// This implementation simply uses the key-value store. One major drawback is that when nodes crash. nothing is updated on the key-value store. Services are a better fit and should be used eventually.
func (n *node) registerService() {
	config := api.DefaultConfig()
	config.Address = n.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Panicln("Unable to contact Service Discovery.")
	}

	kv := consul.KV()
	p := &api.KVPair{Key: n.Name, Value: []byte(n.Addr)}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Panicln("Unable to register with Service Discovery.")
	}

	// store the kv for future use
	n.SDKV = *kv

	log.Println("Successfully registered with Consul.")
}

// setup node ID, initial timestamp and vector clock configuration
func (n *node) setupNode() {
	// search through all encoded values on service discovery network (Consul)
	// for key matching encoding of current node, return its index as ID to indicate time node created
	kvpairs, _, _ := n.SDKV.List("asing465aa", nil)
	pair, _, _ := n.SDKV.Get(n.Name, nil)

	var init_clock_idx int = 0

	// search through existing nodes on KV API, and for each found node (in-order), use that to formulate initial logical time
	// where notation for logical time is step.process_id
	for i := 0; i < len(kvpairs); i++ {
		if kvpairs[i].Key == pair.Key {
			n.NodeId = i + 1
			init_clock_idx = i
		}
	}

	ts_str := "1." + strconv.Itoa(n.NodeId)
	ts, _ := strconv.ParseFloat(ts_str, 64)
	n.NodeTimestamp = float32(ts)

	// given logical time for node, initialize logical vector clock for remaining nodes on network
	n.EventSeen = make([]float32, len(kvpairs))
	for i := 0; i < len(kvpairs); i++ {
		if i == init_clock_idx {
			n.EventSeen[i] = n.NodeTimestamp
		} else {
			n.EventSeen[i] = 0.0
		}
	}

	n.RecvMessages = make(map[float32]string)
	n.AllMessages = ""

	log.Println("Node has ID: ", n.NodeId)
	log.Println("Node has initial timestamp: ", n.NodeTimestamp)
	log.Println("Node has initial vector clock config: ", n.EventSeen)

	acc_path, _ := os.Getwd()
	acc_path += "/accounts.json"
	acc, err := ioutil.ReadFile(acc_path)
	if err != nil {
		log.Println("Accounts file not found", err)
	}
	acc_str := string(acc)
	err = json.Unmarshal([]byte(acc_str), &n.Accounts)
	if err != nil {
		log.Println("Error unmarshalling accounts file", err)
	}
}

// function to read through input file, parse transactions and send as string-message to other nodes in discovery list
// parse entire file and send as a coherent message, later once nodes done receiving, compile all messages and execute
func (n *node) parseFile() {
	file, err := os.Open(n.File)
	if err != nil {
		log.Println("Transactions file not found")
	}

	//configure self-message send flag, for case where request gets made and node sends message to its own queue, use flag to prevent duplicates
	n.MessageSent = 0

	scanner := bufio.NewScanner(file)
	alltext := ""
	for scanner.Scan() {
		line := scanner.Text()
		
		// if end of file reached, exit and begin message sending
		if line == "done" {
			break
		}
		
		alltext += line
		alltext += " "
	}
	// add short delay between subsequent requests to prevent message-passing errors due to other nodes being busy
	n.GreetAll(alltext)
	time.Sleep(3 * time.Second)
	finishMsg := "done"
	n.GreetAll(finishMsg)
	time.Sleep(3 * time.Second)
}

// Start the node.
// This starts listening at the configured address. It also sets up clients for it's peers.
func (n *node) Start() {
	// init required variables
	n.Clients = make(map[string]pb.BankServiceClient)

	// start service / listening
	go n.StartListening()

	// register with the service discovery unit
	n.registerService()

	// setup node, parse its input file & begin message-sending, message-receiving on discovery service

	time.Sleep(12 * time.Second)
	n.setupNode()
	log.Println("Node setup, begin sending and receiving transaction messages")
	time.Sleep(6 * time.Second)
	n.parseFile()
	log.Println("Done processing messages, check bin folder for updated accounts.json files")
}

// Setup a new grpc client for contacting the server at addr.
func (n *node) SetupClient(name string, addr string, msg string) {

	// setup connection with other node
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect: %v, retrying...", err)
	}
	defer conn.Close()
	n.Clients[name] = pb.NewBankServiceClient(conn)

	// vector clock operations before sending message
	// update node timestamp and corresponding vector clock to indicate message processed
	n.EventSeen[n.NodeId-1] = n.EventSeen[n.NodeId-1] + 1
	n.NodeTimestamp = n.EventSeen[n.NodeId-1]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// send message to node only if it is not processed, else continue and send message to remaining nodes on discovery service
	if n.MessageSent < 1 {
		n.RecvMessages[n.NodeTimestamp] = msg
	}
	n.MessageSent = n.MessageSent + 1

	// create gRPC send request, passing node ID, timestamp, message
	r, err := n.Clients[name].DoTransaction(ctx, &pb.TransactionRequest{Message: msg, Timestamp: n.NodeTimestamp, Id: int64(n.NodeId)})

	if err != nil {
		log.Println("could not contact node: %v", err)
	}

	// display reply received at node
	log.Println(r)
}

// function to send request to each node on discovery service
func (n *node) GreetAll(msg string) {
	// get all nodes -- inefficient, but this is just an example
	kvpairs, _, err := n.SDKV.List("asing465aa", nil)
	if err != nil {
		log.Panicln(err)
		return
	}

	// skip self node and access KV API for other nodes on discovery service
	// begin message request process for each found node, whether new or existing
	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, n.Name) == 0 {
			// ourself
			continue
		}

		n.SetupClient(kventry.Key, string(kventry.Value), msg)
	}
}

// begin processing
func main() {
	// pass the port as an argument and also the port of the other node
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("Arguments required: <name> <listening address> <file>")
		os.Exit(1)
	}

	// args in order
	name := args[0]
	listenaddr := args[1]
	file := args[2]

	noden := node{Name: name, Addr: listenaddr, SDAddress: "localhost:8500", File: file, Clients: nil} // noden is for opeartional purposes

	// start the node
	noden.Start()
}

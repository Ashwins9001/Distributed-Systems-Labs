# Distributed-Systems-Labs
Labs covering algorithms for distributed systems using Golang:
- Lab 0: Introduction to Golang
- Lab 1: Using gRPC to transfer client-server data (simulate bank account & calculator operations)
- Lab 2: Creating peer-to-peer system to perform bank transactions using total-ordering on N nodes (replicas) and to output transaction log at end; formulated using gRPC & Consul API for service discovery
                                 
### Lab 1: File Structure
    .
    ├── ...
    ├── calculator/bank                             # Client/Server files
    │   ├── cmd                       
    |      ├── client
    |      ├── server
    │   ├── bin                                     # Miscallaneous files for running
    |      ├── I/O files
    |      ├── port file
    |      ├── executable files
    │   ├── proto                                   # Protobuf configuration to run gRPC data transfer
    |      ├── protobuf definition file
    |      ├── autogen protobuf client/server code
    │   └── ...                
    └── ...
    
    
    
    
    
    
    
    
    
    

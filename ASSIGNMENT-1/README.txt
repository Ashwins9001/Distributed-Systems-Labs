Provided input files handle some edge cases where account number or invalid command is passed to server. 

To run system:
1) Open two seperate sessions
2) In one session execute command: ./bank_server
3) In another session execute command: ./bank_client

If for some reason the executables do not run, to build them:
1) Go to (root directory)/bank/cmd/server
2) Enter command: go build -o (root directory)/bank/bin/bank_server
3) Go to (root directory)/bank/cmd/client
4) Enter command: go build -o (root directory)/bank/bin/bank_client


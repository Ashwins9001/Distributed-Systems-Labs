Tested program on up to 3 nodes. To run, create a file of transactions and place it inside bin folder. For adding an account, update the accounts.json file. Once ready, open multiple sessions and in each one execute program with command:

./node "asing465aa {number}" :{port} {file}

For example:

./node "asing465aa 1" :5000 f1.txt

After, execute: kv delete -recurse asing465aa 

Need to clear data in Consul to prevent error when restarting program with a different number of nodes.

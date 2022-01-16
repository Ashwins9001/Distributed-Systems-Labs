package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//Formulate struct matching shape of JSON data & match corresponding keys from JSON data
type BankAccount struct {
	User    string  `json:"Name"`
	Account int     `json:"AccountID"`
	Amount  float64 `json:"Balance"`
}

func main() {
	//Declare reader wrapped around standard input and parse command line arguments
	reader := bufio.NewReader(os.Stdin)
	input, readError := reader.ReadString('\n')

	if readError != nil {
		fmt.Println("Error occurred", readError)
		return
	}

	//Trim newline and carriage return appended to command-line string on a windows OS
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	content, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Println("File read error occurred", err)
	}

	json_raw := string(content)
	fmt.Println(json_raw)

	//Parse raw JSON into structure by unmarshalling using existing struct
	var account []BankAccount
	unmarshal_err := json.Unmarshal([]byte(json_raw), &account)

	if unmarshal_err != nil {
		fmt.Println("Unmarshal error", unmarshal_err)
	}

	//Modify amount for each account by traversing parsed JSON array
	for i := 0; i < len(account); i++ {
		account[i].Amount += 100
	}

	//Marshal modified data back to JSON-format
	mod_json, _ := json.Marshal(account)

	//Write the modified amounts to a JSON file
	output, _ := os.Create("accountsUpdated.json")
	output.WriteString(string(mod_json))
}

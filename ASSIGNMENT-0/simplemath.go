package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Function to perform math operations
func mathOperation(op string, arg1 float64, arg2 float64) {
	if strings.ToLower(op) == "add" {
		fmt.Println(arg1 + arg2)
	} else if strings.ToLower(op) == "sub" {
		fmt.Println(arg1 - arg2)

	} else if strings.ToLower(op) == "mult" {
		fmt.Println(arg1 * arg2)

	} else if strings.ToLower(op) == "div" {
		fmt.Println(arg1 / arg2)
	} else {
		fmt.Println("Invalid input")
	}
}

func main() {
	//Declare reader wrapped around standard input and parse command line arguments
	reader := bufio.NewReader(os.Stdin)
	input, error := reader.ReadString('\n')

	if error != nil {
		fmt.Println("Error occurred", error)
		return
	}

	//Parse string into tokens by splitting by whitespace using Fields method
	input = strings.TrimSuffix(input, "\n")

	funcArgs := strings.Fields(input)

	//Convert arguments to numbers
	num1, _ := strconv.ParseFloat(funcArgs[1], 10)
	num2, _ := strconv.ParseFloat(funcArgs[2], 10)

	//Call operation function and provide output
	mathOperation(funcArgs[0], num1, num2)
}

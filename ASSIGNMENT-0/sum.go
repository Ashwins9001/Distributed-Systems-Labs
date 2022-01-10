package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

	//Open file for reading and go line by line using the scanner 
	var file, fileError = os.Open(input)
	
	if fileError != nil {
		fmt.Println("Error opening file", fileError)
	}
	
	defer file.Close()

	//Create scanner to read file line by line 
	scanner := bufio.NewScanner(file)

	var runningSum float64
	runningSum = 0

	//Read in each line & add to the running sum 
	for scanner.Scan() {
		currentNum, _ := strconv.ParseFloat(scanner.Text(), 10)
		runningSum += currentNum
	}

	fmt.Println("Sum: ", runningSum)
}
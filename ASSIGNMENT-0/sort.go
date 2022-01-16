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

	//Create empty slice to append to; use slice as opposed to array because it is resizeable
	inputCollection := []float64{}

	for scanner.Scan() {
		currentNum, _ := strconv.ParseFloat(scanner.Text(), 10)
		inputCollection = append(inputCollection, currentNum)
	}

	//Apply insertion sort algorithm to put it into ascending order
	for i := 1; i < len(inputCollection); i++ {
		//Find ith array element and initialize pointer, j, to previous element that decrements to sort elements
		currentElement := inputCollection[i]
		j := i - 1

		//While in array bounds and while the ith element is smaller than a previous element, push the previous element to the right to make space
		//In effect have put larger, unsorted elements ahead, in the proper spot
		for j >= 0 && currentElement < inputCollection[j] {
			inputCollection[j+1] = inputCollection[j]
			j -= 1
		}

		//Once an element smaller than the ith element is found, then it is in the correct position, then place element using the j-pointer
		inputCollection[j+1] = currentElement
	}

	//Write the sorted slice to a file
	output, _ := os.Create("sortedNumbers.txt")
	for k := 0; k < len(inputCollection); k++ {
		outNum := fmt.Sprintf("%f", inputCollection[k])
		output.WriteString(outNum + "\n")
	}
}

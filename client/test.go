package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	// Open the CSV file
	file, err := os.Open("rund.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all the records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Check that we have enough records to access the required rows
	if len(records) >= 20 {
		// Accessing specific values:
		firstColValue := records[1][0]  // 2nd value from the first column (index starts from 0)
		secondColValue := records[8][1] // 9th value from the second column
		thirdColValue := records[19][2] // 20th value from the third column

		fmt.Println("2nd value from first column:", firstColValue)
		fmt.Println("9th value from second column:", secondColValue)
		fmt.Println("20th value from third column:", thirdColValue)
	} else {
		fmt.Println("Not enough rows in the CSV to access the desired values")
	}
}

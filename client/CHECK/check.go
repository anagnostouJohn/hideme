package check

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func Check(msg string, err error) error {
	if err != nil {
		fmt.Println(msg, "  ", err)
		return err
	}
	return nil
}

func OpenAndReadFiles(filepath string) []byte {

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return fileContent
}

func Pouse() {
	fmt.Println("Press any key to continue...")

	// Create a new reader
	reader := bufio.NewReader(os.Stdin)

	// Read a single character from the input
	_, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Println("Continuing...")
}

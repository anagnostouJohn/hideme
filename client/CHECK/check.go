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
		f, erra := os.OpenFile("\\tmp\\err.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if erra != nil {
			panic(erra)
		}

		defer f.Close()

		if _, errs := f.WriteString(msg + "  " + err.Error() + "\n"); errs != nil {
			panic(errs)
		}
		// fmt.Println(msg, "  ", err)
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
	fmt.Println(fileContent)
	fmt.Println(string(fileContent))
	if err != nil {
		fmt.Println(err)
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

package check

import (
	"bufio"
	"fmt"
	"os"
)

func Check(msg string, err error) error {
	if err != nil {
		fmt.Println(msg, "  ", err)
		return err
	}
	return nil
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

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("set", "+o", "history")
	cmd.Run()

	cmd = exec.Command("history")
	var out bytes.Buffer
	cmd.Stdout = &out

	outSpit := strings.Split(out.String(), "\n")
	fmt.Println(outSpit)
}

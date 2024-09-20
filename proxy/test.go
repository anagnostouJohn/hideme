package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	vars "test/VARS"

	"github.com/BurntSushi/toml"
)

func main() {

	decodedToml, err := base64.StdEncoding.DecodeString("W2NsaWVudF0KCVVzZXIgPSAid2luZSIKCVBvcnQgPSAiMjIiCglIb3N0ID0gIjE5Mi4xNjguMjMuODkiCglQYXNzID0gIjEyMzQiCgoKW2ZsYWdzXQoJRGVzdHIgPSBmYWxzZQoJQ29tYm8gPSBmYWxzZSAKCUhpZGVtZSA9IHRydWUKCUJyRmlsZSA9ICIvaG9tZS91YnVudHUvdGVzdC5jc3YiCiAgICBDb25uZWN0ZWRVc2VyID0gInVidW50dSIKCVRocmVhZHMgPSAzCglLbm9ja0FsaXZlID0gNjY2NgoJS25vY2tEYXRhID0gWzg2NjYsNzY2Nl0KUGlkVG9TdGFydCA9ICI5NDM1Ig")
	if err != nil {
		fmt.Println("Error decoding base64 TOML:", err)
		return
	}

	var conf vars.Config

	// Use bytes.NewReader to pass the decoded string as a reader
	reader := bytes.NewReader(decodedToml)
	z := toml.NewDecoder(reader)
	z.Decode(&conf)
	fmt.Print(conf)
}

package main

import (
	"fmt"
	"log"
	"sync"
	checkactive "test/proxy/CHECKACTIVE"
	sendbf "test/proxy/SENDBF"
	vars "test/vars"

	"github.com/BurntSushi/toml"
)

// var base64PidToStart string
var confa vars.Config
var wg sync.WaitGroup

func init() {
	if _, err := toml.DecodeFile("config.toml", &confa); err != nil {
		log.Fatal(err)
	}
}

func main() {

	sendbf.SendBf(confa)
	fmt.Println(string(vars.Blue))
	fmt.Println("Finished Start Checking")
	fmt.Println(string(vars.Reset))
	// wg.Add(1)
	checkactive.Checkactive(confa)
	// wg.Done()

}

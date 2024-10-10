package main

import (
	"fmt"
	"log"
	"sync"
	checkactive "test/proxy/CHECKACTIVE"
	vars "test/vars"
	"time"

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

	// sendbf.SendBf(confa)

	fmt.Println("ASDASDASDAD")
	for {
		// wg.Add(1)
		checkactive.Checkactive(confa)
		// wg.Done()
		time.Sleep(10 * time.Second)
	}
}

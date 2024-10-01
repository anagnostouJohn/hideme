package main

import (
	"log"
	sendbf "test/proxy/SENDBF"
	vars "test/vars"

	"github.com/BurntSushi/toml"
)

// var base64PidToStart string
var confa vars.Config

func init() {
	if _, err := toml.DecodeFile("config.toml", &confa); err != nil {
		log.Fatal(err)
	}
}

func main() {

	sendbf.SendBf(confa)

}

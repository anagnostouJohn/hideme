package main

import (
	"log"
	checkactive "test/proxy/CHECKACTIVE"
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
	checkactive.Checkactive(confa)

}

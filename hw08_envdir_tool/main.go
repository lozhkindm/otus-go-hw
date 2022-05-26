package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("at least 3 args required")
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(RunCmd(os.Args[2:], env))
}

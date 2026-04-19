package main

import (
	"fmt"
	"log"
	"os"

	"example.com/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Println("name arg does not present")
		os.Exit(1)
	}
	messages, err := greetings.Hellos(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	for name, message := range messages {
		fmt.Printf("%s -> %s\n", name, message)
	}
}

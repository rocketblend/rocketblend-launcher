package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher"
)

func main() {
	log.Println("Starting launcher...")

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("PANIC: %v", r)
			exit()
		}
	}()

	if err := launcher.Launch(); err != nil {
		log.Printf("Something went wrong: %v", err)
		exit()
	}

	log.Println("Launcher finished!")
}

func exit() {
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

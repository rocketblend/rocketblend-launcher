package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v", r)
			exit()
		}
	}()

	var projectPath string
	if len(os.Args) > 1 {
		projectPath = os.Args[1]
	}

	launcher := launcher.New(projectPath)
	if err := launcher.Launch(); err != nil {
		//fmt.Printf("Something went wrong: %v\n", err)
		exit()
	}
}

func exit() {
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

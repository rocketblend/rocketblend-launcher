package main

import (
	"fmt"

	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher"
)

func main() {
	fmt.Println("Starting launcher...")

	if err := launcher.Launch(); err != nil {
		fmt.Println("Something went wrong, ", err)
	}

	fmt.Println("Launcher finished!")
}

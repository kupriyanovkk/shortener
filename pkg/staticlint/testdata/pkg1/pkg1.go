package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Println("!")

	os.Exit(3) // want "Direct os.Exit call detected in main function"
}

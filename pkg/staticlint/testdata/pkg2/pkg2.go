package main

import (
	"fmt"
	"os"
)

func osExit() {
	os.Exit(3)
}

func main() {
	defer fmt.Println("!")

	osExit()
}

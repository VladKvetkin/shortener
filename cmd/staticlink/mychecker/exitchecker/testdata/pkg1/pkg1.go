package main

import "os"

func main() {
	os.Exit(1) // want "call os.Exit in main package"
}

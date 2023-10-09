package main

import "os"

func errCheckFunc() {
	os.Exit(1) // want "call os.Exit in main package"
}

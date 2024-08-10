package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func usage() string {

	return fmt.Sprintf("Usage: %s <filename>", filepath.Base(os.Args[0]))
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage())
		return
	}

	filename := os.Args[1]
	fmt.Println("Filename: ", filename)
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	assemblyFileExtension = ".asm"
	hackFileExtension     = ".hack"
)

func usage() string {

	return fmt.Sprintf("Usage: %s <filename>", filepath.Base(os.Args[0]))
}

func replaceExtension(fileName string, newExt string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))] + newExt
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage())
		return
	}

	fileName := os.Args[1]
	if filepath.Ext(fileName) != assemblyFileExtension {
		log.Fatalf("File name must end with %v\n", assemblyFileExtension)
	}
	inputFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Cannot open file: ", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(replaceExtension(fileName, hackFileExtension))
	if err != nil {
		log.Fatalln("Cannot create file: ", err)
	}

	parser := NewParser(inputFile)

	for parser.HasMoreCommands() {
		parser.Advance()
		switch parser.CommandType() {
		case A_COMMAND, L_COMMAND:
			symbol, _ := strconv.ParseUint(parser.Symbol(), 10, 16)
			fmt.Fprintf(outputFile, "%016b\n", symbol)
		case C_COMMAND:
			dest := Dest(parser.Dest())
			comp := Comp(parser.Comp())
			jump := Jump(parser.Jump())
			fmt.Fprintf(outputFile, "111%v%v%v\n", comp, dest, jump)
		}
	}

}

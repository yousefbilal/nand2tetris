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

func setupFiles(inputFileName string) (*os.File, *os.File, error) {
	if filepath.Ext(inputFileName) != assemblyFileExtension {
		return nil, nil, fmt.Errorf("file name must end with %v", assemblyFileExtension)
	}
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open file: %v", err)
	}
	outputFile, err := os.Create(replaceExtension(inputFileName, hackFileExtension))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create file: %v", err)
	}
	return inputFile, outputFile, nil
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
	inputFile, outputFile, err := setupFiles(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	defer outputFile.Close()

	parser := NewParser(inputFile)
	st := NewSymbolTable()
	//first pass
	lineNumber := 0
	for parser.HasMoreCommands() {
		parser.Advance()
		if parser.CommandType() == L_COMMAND {
			symbol := parser.Symbol()
			st.AddEntry(symbol, lineNumber)
			continue
		}
		lineNumber++
		if parser.CommandType() == INVALID {
			log.Fatalf("invlaid command at line %v", lineNumber)
		}
	}
	parser.Reset()
	//second pass
	ramAddress := 16
	for parser.HasMoreCommands() {
		parser.Advance()
		switch parser.CommandType() {
		case A_COMMAND:
			symbol := parser.Symbol()
			var resolvedSymbol int
			if parsed, err := strconv.ParseUint(symbol, 10, 16); err == nil {
				resolvedSymbol = int(parsed)
			} else {
				if !st.Contains(symbol) {
					st.AddEntry(symbol, ramAddress)
					ramAddress++
				}
				resolvedSymbol = st.GetAddress(symbol)
			}
			fmt.Fprintf(outputFile, "%016b\n", resolvedSymbol)
		case C_COMMAND:
			dest := Dest(parser.Dest())
			comp := Comp(parser.Comp())
			jump := Jump(parser.Jump())
			fmt.Fprintf(outputFile, "111%v%v%v\n", comp, dest, jump)
		}
	}
}

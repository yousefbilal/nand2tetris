package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

type Command int

const (
	NONE Command = iota
	A_COMMAND
	C_COMMAND
	L_COMMAND
	INVALID
)

type Parser struct {
	inputFile      *os.File
	scanner        *bufio.Scanner
	commentRegex   *regexp.Regexp
	aRegex         *regexp.Regexp
	cRegex         *regexp.Regexp
	lRegex         *regexp.Regexp
	currentCommand string
	commandType    Command
}

func NewParser(inputFile *os.File) *Parser {
	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Scan()
	commentRegex := regexp.MustCompile(`^\s*//.*$`)
	aRegex := regexp.MustCompile(`^@((?:\d+)|(?:[A-Za-z_.$:][\w.$:]*))$`)
	cRegex := regexp.MustCompile(`^(?:(A?M?D?)=)?([AMD+-01!&|]{1,3})(?:;([JGTEQLNMP]{0,3}))?$`)
	lRegex := regexp.MustCompile(`^\(([A-Za-z_.$:][\w.$:]*)\)$`)

	return &Parser{
		inputFile:      inputFile,
		scanner:        fileScanner,
		commentRegex:   commentRegex,
		aRegex:         aRegex,
		cRegex:         cRegex,
		lRegex:         lRegex,
		currentCommand: "",
		commandType:    NONE,
	}
}

func (p *Parser) HasMoreCommands() bool {
	text := strings.TrimSpace(p.scanner.Text())
	for len(text) == 0 || p.commentRegex.MatchString(text) {
		if !p.scanner.Scan() {
			return false
		}
		text = strings.TrimSpace(p.scanner.Text())
	}
	return true
}

func (p *Parser) Advance() {
	p.currentCommand = strings.TrimSpace(p.scanner.Text())
	switch {
	case p.aRegex.MatchString(p.currentCommand):
		p.commandType = A_COMMAND
	case p.cRegex.MatchString(p.currentCommand):
		p.commandType = C_COMMAND
	case p.lRegex.MatchString(p.currentCommand):
		p.commandType = L_COMMAND
	default:
		p.commandType = INVALID
	}
	p.scanner.Scan()
}

func (p *Parser) CommandType() Command {
	return p.commandType
}

func (p *Parser) aSubmatch(text string) []string {
	return p.aRegex.FindStringSubmatch(text)
}

func (p *Parser) cSubmatch(text string) []string {
	return p.cRegex.FindStringSubmatch(text)
}

func (p *Parser) lSubmatch(text string) []string {
	return p.lRegex.FindStringSubmatch(text)
}

func (p *Parser) Symbol() string {
	switch p.commandType {
	case A_COMMAND:
		return p.aSubmatch(p.currentCommand)[1]
	case L_COMMAND:
		return p.lSubmatch(p.currentCommand)[1]
	}
	return ""
}

func (p *Parser) Dest() string {
	return p.cSubmatch(p.currentCommand)[1]
}

func (p *Parser) Comp() string {
	return p.cSubmatch(p.currentCommand)[2]
}

func (p *Parser) Jump() string {
	return p.cSubmatch(p.currentCommand)[3]
}

func (p *Parser) Reset() {
	p.inputFile.Seek(0, io.SeekStart)
	p.scanner = bufio.NewScanner(p.inputFile)
}

package main

import (
	"bufio"
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
	scanner        *bufio.Scanner
	comment_regex  *regexp.Regexp
	a_regex        *regexp.Regexp
	c_regex        *regexp.Regexp
	l_regex        *regexp.Regexp
	currentCommand string
	commandType    Command
}

func NewParser(input_file string) *Parser {
	file, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Scan()
	comment_regexp := regexp.MustCompile(`^\s*//.*$`)
	a_regex := regexp.MustCompile(`^@(?:(\d+)|([A-Za-z_.$:][\w.$:]*))$`)
	c_regex := regexp.MustCompile(`^(?:(A?M?D?)=)?([AMD+-01!&|]{1,3})(?:;([JGTEQLNMP]{0,3}))?$`)
	l_regex := regexp.MustCompile(`^\((\d+)|([A-Za-z_.$:][\w.$:]*)\)$`)

	return &Parser{
		scanner:        fileScanner,
		comment_regex:  comment_regexp,
		a_regex:        a_regex,
		c_regex:        c_regex,
		l_regex:        l_regex,
		currentCommand: "",
		commandType:    NONE,
	}
}

func (p *Parser) HasMoreCommands() bool {
	text := strings.TrimSpace(p.scanner.Text())
	for len(text) == 0 || p.comment_regex.MatchString(text) {
		if !p.scanner.Scan() {
			return false
		}
		text = strings.TrimSpace(p.scanner.Text())
	}
	return p.c_regex.MatchString(text) || p.a_regex.MatchString(text) || p.l_regex.MatchString(text)
}

func (p *Parser) Advance() {
	p.currentCommand = strings.TrimSpace(p.scanner.Text())
	switch {
	case p.a_regex.MatchString(p.currentCommand):
		p.commandType = A_COMMAND
	case p.c_regex.MatchString(p.currentCommand):
		p.commandType = C_COMMAND
	case p.l_regex.MatchString(p.currentCommand):
		p.commandType = L_COMMAND
	default:
		p.commandType = INVALID
	}
}

func (p *Parser) CommandType() Command {
	return p.commandType
}

func (p *Parser) aSubmatch(text string) []string {
	return p.a_regex.FindStringSubmatch(text)
}

func (p *Parser) cSubmatch(text string) []string {
	return p.c_regex.FindStringSubmatch(text)
}

func (p *Parser) lSubmatch(text string) []string {
	return p.l_regex.FindStringSubmatch(text)
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

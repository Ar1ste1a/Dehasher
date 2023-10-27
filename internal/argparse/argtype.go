package argparse

import (
	"fmt"
	"os"
	"strings"
)

type ArgType int32

const (
	INTEGER   ArgType = iota
	STRING    ArgType = iota
	BOOL      ArgType = iota
	IPADDRESS ArgType = iota
	CIDRRANGE ArgType = iota
	INVALID   ArgType = iota
)

// getArgTypeFromString Return an ArgType object based on the string passed
func getArgTypeFromString(s string) ArgType {
	lowerS := strings.ToLower(s)
	switch lowerS {
	case "integer":
		return INTEGER
	case "string":
		return STRING
	case "bool":
		return BOOL
	case "ipaddress":
		return IPADDRESS
	case "cidrrange":
		return CIDRRANGE
	default:
		fmt.Println(fmt.Errorf("invalid argument type passed: %s", s))
		os.Exit(0)
		return INVALID
	}
}

// toString generates a string for each ArgType
func (at ArgType) toString() string {
	switch at {
	case BOOL:
		return "BOOL"
	case STRING:
		return "STRING"
	case INTEGER:
		return "INTEGER"
	case IPADDRESS:
		return "IPADDRESS"
	case CIDRRANGE:
		return "CIDRRANGE"
	default:
		fmt.Println(fmt.Errorf("invalid ArgType"))
		os.Exit(0)
		return ""
	}
}

package argparse

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Argument a single flag that can be passed into a program
type Argument struct {
	argType     ArgType
	flag        string
	alias       string
	description string
	required    bool
	value       interface{}
	name        string
}

// getRealValue return an object of the Argument.value cast to the Argument.argType
func (a *Argument) getRealValue() any {
	switch a.argType {
	case BOOL:
		out, _ := strconv.ParseBool(a.value.(string))
		return out
	case INTEGER:
		out, _ := strconv.Atoi(a.value.(string))
		return out
	case STRING:
		return a.value.(string)
	case IPADDRESS:
		valueString := a.value.(string)
		ip := net.ParseIP(valueString)
		if ip == nil {
			fmt.Println(fmt.Errorf("[!] invalid IP Address passed for: %s", a.flag))
			os.Exit(0)
		}
		return ip
	case CIDRRANGE:
		valueString := a.value.(string)
		_, cidr, err := net.ParseCIDR(valueString)
		if err != nil {
			fmt.Println(fmt.Errorf("[!] invalid CIDR Range passed for: %s", a.flag))
			os.Exit(0)
		}
		return cidr
	default:
		fmt.Println(fmt.Errorf("[!] invalid ArgType passed: %s", a.argType.toString()))
		os.Exit(0)
		return nil
	}
}

// NewArgument Create a new Argument object and return it by value
func NewArgument(argType ArgType, flag, alias, description string, defaultValue any, required bool) Argument {
	argName := strings.Trim(flag, "- \t")
	argName = strings.Replace(argName, "-", "_", -1)
	arg := Argument{argType: argType, flag: flag, alias: alias, description: description, value: defaultValue, required: required, name: argName}
	return arg
}

// getValue return a type assertion of the argument value based on the argType
func (a *Argument) getValue() any {
	switch a.argType {
	case BOOL:
		return a.value.(bool)
	case INTEGER:
		return a.value.(int)
	case STRING:
		return a.value.(string)
	default:
		return a.value.(string)
	}
}

func (a *Argument) setValue(newVal any) {
	a.value = newVal.(string)
}

// isRequired Determines if a given argument is required based on Argument.required
func (a *Argument) isRequired() bool {
	return a.required
}

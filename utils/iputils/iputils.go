package iputils

import (
	"net"
	"regexp"
)

// Checks if net validates ip or string matches ipv4 regex pattern
func NameIsIP(name string) bool {
	return net.ParseIP(name) != nil || regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\.(?:xxx|\d{1,3})$`).MatchString(name)
}

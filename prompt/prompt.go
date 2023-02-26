package prompt

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func String(prompt string, args ...interface{}) string {
	var s string
	fmt.Printf(prompt+" : ", args...)
	fmt.Scanln(&s)
	return s
}

func StringDefault(prompt string, def string, args ...interface{}) string {
	s := def
	fmt.Printf(prompt+" ["+s+"] : ", args...)
	fmt.Scanln(&s)
	return s
}

func StringRequired(prompt string, args ...interface{}) (s string) {
	for strings.Trim(s, " ") == "" {
		s = String(prompt, args...)
	}
	return s
}

func Password(prompt string, args ...interface{}) string {
	fmt.Printf(prompt+" : ", args...)
	bytePassword, _ := terminal.ReadPassword(0)
	fmt.Println("")
	return string(bytePassword)
}

package echoservice

import (
	"strings"
)

//EchoService takes an input string returns it, uppercased.
type EchoService interface {
	Echo(string) string
}

type EchoServiceImpl struct {
}

func (svc EchoServiceImpl) Echo(s string) string {
	return strings.ToUpper(s)
}

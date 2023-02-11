package main

import (
	"github.com/scottcagno/webslinger/pkg/web/jwt"
)

func main() {

	method := jwt.RS256

	tm := jwt.NewTokenManager(method, method.GenerateKeyPair())

	_ = tm
}

package infra

import "golang.org/x/oauth2"

type GitHub struct {
	oauth oauth
}

type oauth struct {
	cnf oauth2.Config
}

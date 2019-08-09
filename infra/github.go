package infra

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func NewGitHub() *GitHub {
	return &GitHub{
		oauth: oauth{
			cnf: oauth2.Config{
				ClientID:     "5a24485cf2fe2ca8fab4",
				ClientSecret: "63a169863256d15eca02ac6ade415f93b2692e28",
				RedirectURL:  "http://localhost/tapioca/authorization",
				Scopes: []string{
					"repo:status", "read:user",
				},
				Endpoint: github.Endpoint,
			},
		},
	}
}

type GitHub struct {
	oauth oauth
}

type oauth struct {
	cnf oauth2.Config
}

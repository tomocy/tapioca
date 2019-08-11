package client

func newTwitter(cnf config) *Twitter {
	return &Twitter{
		cnf: cnf,
	}
}

type Twitter struct {
	cnf config
}

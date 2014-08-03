package redisgreen

import (
	"os"
)

var (
	ValidTokenClient   = Client{os.Getenv("TOKEN")}
	InvalidTokenClient = Client{"invalid"}
)

func init() {
	if len(ValidTokenClient.Token) == 0 {
		panic("TOKEN environment variable missing")
	}
	ApiUrl = "http://localhost:4100"
}

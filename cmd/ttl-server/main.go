package main

import (
	"os"

	servercli "ttl-cli/internal/server/cli"
)

func main() {
	os.Exit(servercli.Run())
}

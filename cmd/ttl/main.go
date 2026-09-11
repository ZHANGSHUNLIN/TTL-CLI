package main

import (
	"os"

	clientcli "ttl-cli/internal/client/cli"
)

func main() {
	os.Exit(clientcli.Run())
}

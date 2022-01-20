package main

import (
	"github.com/me0888/twitter/cmd/app"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dns := "postgres://app:pass@localhost:5432/db"
	if err := app.Execute(host, port, dns); err != nil {
		os.Exit(1)
	}
}

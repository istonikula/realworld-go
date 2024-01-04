package main

import (
	"log"

	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
)

func main() {
	cfg := boot.ReadConfig()

	boot.Migrate(&cfg.DataSource)

	if err := boot.Router(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}

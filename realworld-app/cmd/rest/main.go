package main

import (
	"log"

	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
)

func main() {
	cfg := boot.ReadConfig("../../config.yml")

	boot.Migrate("../../db", cfg.DataSource)

	if err := rest.Router(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}

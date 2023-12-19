package main

import (
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/gin-gonic/gin"
	appDb "github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/url"
)

func main() {
	migrate()

	if err := router(db()).Run(); err != nil {
		log.Fatal(err)
	}
}

func migrate() {
	u, _ := url.Parse("postgres://postgres:secret@127.0.0.1:5432/realworld?sslmode=disable")
	err := dbmate.New(u).Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func db() *sqlx.DB {
	return sqlx.MustConnect("postgres", "user=realworld password=secret dbname=realworld sslmode=disable")
}

func router(db *sqlx.DB) *gin.Engine {
	auth := domain.Auth{Settings: domain.Security{TokenSecret: "TODO token"}}

	txMgr := &appDb.TxMgr{DB: db}

	r := gin.Default()
	rest.UserRoutes(r, &auth, txMgr)

	return r
}

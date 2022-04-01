package pg

import (
	"log"

	"github.com/go-pg/pg/v10"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {

	container.Provide(func() *pg.Options {
		return &pg.Options{
			Database: "mobile-banking-service",
			User:     "postgres",
			Password: "postgres",
		}
	})

	container.Provide(func(opts *pg.Options) *pg.DB {
		return pg.Connect(opts)
	})
}

func Invoke(container *dig.Container) {
	err := container.Invoke(func(db *pg.DB) {
		log.Println("invoke db...")
		pgDb = db
	})
	if err != nil {
		log.Fatal(err)
	}
}

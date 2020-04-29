package main

import (
	"coursework/postgres"
	"coursework/validators"
	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v7"
)

var (
	Postgres  postgres.RegistrationConn
	Redis     *redis.Client
	Salt      []byte
	Router    *router.Router
	Validator *validator.Validate
	ServePort string
)

func init() {
	Postgres = postgres.RegistrationConn{}
	Redis = redis.NewClient(&redis.Options{Addr: "redis:6379"})
	Salt = []byte("Ilya Bychkov")
	Router = router.New()
	Validator = validator.New()
	ServePort = "8090"
}

func Initializer() error {
	if err := Postgres.CreateConnection("host=postgres user=me password=12345 dbname=my_coursework_db");
	err != nil { return err }

	if err := Postgres.InitDatabasesIfNotExist(); err != nil {
		return err
	}

	if err := Redis.Ping().Err(); err != nil {
		return err
	}

	if err := Validator.RegisterValidation("sex", validators.ValidateSex); err != nil {
		return err
	}

	return nil
}

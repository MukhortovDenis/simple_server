package main

import (
	"fmt"
	simpleserver "vsu"
	"vsu/config"
	"vsu/internal/auth/cache"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	srv := simpleserver.NewService(cache.NewAccountCache(), cache.NewPermissionsCache())

	fmt.Println("Server starting...")

	if err = srv.Start(cfg); err != nil {
		return err
	}

	fmt.Println("Server closed")

	return nil
}

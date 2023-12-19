package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
)

//nolint:gochecknoglobals
var (
	envConfig = new(Env)
	once      = new(sync.Once)

	err error
)

type Env struct {
	HTTPListenAddr string        `env:"HTTP_LISTEN_ADDR,required"`
	HTTPTimeout    time.Duration `env:"HTTP_TIMEOUT" envDefault:"30s"`

	AdminLogin    string `env:"ADMIN_LOGIN,required"`
	AdminPassword string `env:"ADMIN_PASSWORD,required"`
}

func FromEnv() (*Env, error) {
	once.Do(func() {
		err = env.Parse(envConfig)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return envConfig, nil
}

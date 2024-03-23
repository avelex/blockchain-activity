package config

import (
	"os"
	"sync"
	"time"
)

const (
	rpcAccessTokenPathEnv = "RPC_ACCESS_TOKEN_PATH"
	httpPortEnv           = "HTTP_PORT"
	httpHostEnv           = "HTTP_HOST"
	shutdownTimeoutEnv    = "SHUTDOWN_TIMEOUT"
)

const (
	defaultHTTPPort        = "8080"
	defaultHTTPHost        = "0.0.0.0"
	defaultShutdownTimeout = "5s"
)

type Config struct {
	RPC struct {
		AccessToken string
	}

	HTTP struct {
		Port string
		Host string
	}

	ShutdownTimeout time.Duration
}

var once sync.Once
var cfg Config

func InitConfig() Config {
	once.Do(func() {
		accessTokenPath, ok := os.LookupEnv(rpcAccessTokenPathEnv)
		if !ok {
			panic("missing access token path")
		}

		token, err := os.ReadFile(accessTokenPath)
		if err != nil {
			panic(err)
		}

		shutdownTimeoutStr, ok := os.LookupEnv(shutdownTimeoutEnv)
		if !ok {
			shutdownTimeoutStr = defaultShutdownTimeout
		}

		shutdownTimeout, err := time.ParseDuration(shutdownTimeoutStr)
		if err != nil {
			panic(err)
		}

		if cfg.HTTP.Host, ok = os.LookupEnv(httpHostEnv); !ok {
			cfg.HTTP.Host = defaultHTTPHost
		}

		if cfg.HTTP.Port, ok = os.LookupEnv(httpPortEnv); !ok {
			cfg.HTTP.Port = defaultHTTPPort
		}

		cfg.ShutdownTimeout = shutdownTimeout
		cfg.RPC.AccessToken = string(token)
	})
	return cfg
}

package main

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"

	"github.com/jestclout/jestclout-go/game"
	"github.com/jestclout/jestclout-go/handler"
)

type config struct {
	BindAddress    string `env:"BIND_ADDRESS" envDefault:"0.0.0.0"`
	BindPort       string `env:"BIND_PORT" envDefault:"3001"`
	PromptFileName string `env:"PROMPT_FILE" envDefault:"prompts.txt"`
}

// This example demonstrates a trivial echo server.
func main() {
	ll := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		ll.Fatal().Err(err).Msg("unable to parse env config")
	}

	prompts := loadPrompts(cfg.PromptFileName, ll)
	manager := game.NewManager(prompts, game.DefaultConfig())

	h := handler.New(manager, ll)
	addr := net.JoinHostPort(cfg.BindAddress, cfg.BindPort)

	server := &http.Server{
		Handler:        h.Router(),
		Addr:           addr,
		WriteTimeout:   60 * time.Second,
		ReadTimeout:    60 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		ll.Info().Str("bind_address", addr).Msg("server started")
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ll.Fatal().Err(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ll.Info().Msg("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		ll.Fatal().Err(err)
	}

	ll.Info().Msg("server shut down complete")
}

func loadPrompts(fileName string, ll zerolog.Logger) []string {
	f, err := os.Open(fileName)
	if err != nil {
		ll.Fatal().Err(err).Msg("failed opening prompt file")
	}
	defer f.Close()

	var prompts []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		prompts = append(prompts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		ll.Fatal().Err(err).Msg("failed reading lines from file")
	}

	return prompts
}

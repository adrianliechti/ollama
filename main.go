package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	ollama "github.com/ollama/ollama/api"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Host   string   `env:"OLLAMA_HOST, default=127.0.0.1:11434"`
	Models []string `env:"OLLAMA_MODELS"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var c Config
	envconfig.MustProcess(ctx, &c)

	client, _ := ollama.ClientFromEnvironment()

	println("Waiting for Ollama Server...")

	if err := waitUntilReady(ctx, client); err != nil {
		panic(err)
	}

	for _, model := range c.Models {
		fmt.Printf("Pulling %s...\n", model)

		if err := pullModel(ctx, client, model); err != nil {
			panic(err)
		}
	}
}

func waitUntilReady(ctx context.Context, client *ollama.Client) error {
	var result error

	for ctx.Err() == nil {
		time.Sleep(500 * time.Millisecond)

		result = client.Heartbeat(ctx)

		if result == nil {
			return nil
		}
	}

	return result
}

func pullModel(ctx context.Context, client *ollama.Client, model string) error {
	handler := func(p ollama.ProgressResponse) error {
		return nil
	}

	return client.Pull(ctx, &ollama.PullRequest{
		Model: model,
	}, handler)
}

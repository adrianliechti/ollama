package main

import (
	"context"
	"fmt"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	ollama "github.com/ollama/ollama/api"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Models []string `env:"MODEL"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//wd, _ := os.Getwd()
	//path := filepath.Join(wd, "data")

	var c Config
	envconfig.MustProcess(ctx, &c)

	client, _ := ollama.ClientFromEnvironment()

	server := exec.CommandContext(ctx, "ollama", "serve")
	// server.Env = []string{
	// 	"HOME=" + path,
	// 	"OLLAMA_HOST=127.0.0.1:11434",
	// 	"OLLAMA_FLASH_ATTENTION=1",
	// }

	//server.Stdout = os.Stdout
	//server.Stderr = os.Stderr

	go func() {
		println("Ollama Server starting...")

		waitUntilReady(ctx, client)

		for _, model := range c.Models {
			fmt.Printf("Pulling model %s...\n", model)
			pullModel(ctx, client, model)
		}

		println("Ollama Server is ready")

	}()

	if err := server.Run(); err != nil {
		panic(err)
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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"openai-slack-example/gpt3"
	"openai-slack-example/tokenizer"

	"github.com/alexflint/go-arg"
)

func main() {
	cfg, err := newConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load config: %v\n", err)
		os.Exit(1)
	}

	err = run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "application returned an error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	_, err := tokenizer.NewEncoder()
	if err != nil {
		return err
	}

	client := gpt3.NewClient(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIKey, "text-002")

	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{"The first thing you should know about javascript is"},
		MaxTokens: gpt3.IntPtr(30),
		Stop:      []string{"."},
		Echo:      true,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Text)

	return nil
}

type config struct {
	AzureOpenAIEndpoint string `arg:"--azure-openai-endpoint,env:AZURE_OPENAI_ENDPOINT,required" help:"The endpoint for Azure OpenAI service"`
	AzureOpenAIKey      string `arg:"--azure-openai-key,env:AZURE_OPENAI_KEY,required" help:"The (api) key for Azure OpenAI service"`
}

func newConfig(args []string) (*config, error) {
	cfg := &config{}
	parser, err := arg.NewParser(arg.Config{
		Program:   "openai-slack-example",
		IgnoreEnv: false,
	}, cfg)
	if err != nil {
		return &config{}, err
	}

	err = parser.Parse(args)
	if err != nil {
		return &config{}, err
	}

	return cfg, err
}

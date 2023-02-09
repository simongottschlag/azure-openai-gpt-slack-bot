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

var (
	maxTokensMap = map[string]int{
		"text-002": 4097,
		"text-003": 4097,
	}
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

	prompt := "The first thing you should know about javascript is"
	client := gpt3.NewClient(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIKey, cfg.OpenAIDeploymentName)
	completion, err := gptCompletion(ctx, client, prompt, cfg.OpenAIDeploymentName)
	if err != nil {
		return err
	}

	fmt.Printf("Response: %s\n", completion)

	return nil
}

func gptCompletion(ctx context.Context, client gpt3.Client, prompt string, deploymentName string) (string, error) {
	maxTokens, err := calculateMaxTokens(prompt, deploymentName)
	if err != nil {
		return "", err
	}
	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{prompt},
		MaxTokens: maxTokens,
		Stop:      []string{"."},
		Echo:      true,
		N:         gpt3.ToPtr(1),
	})

	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

func calculateMaxTokens(prompt string, deploymentName string) (*int, error) {
	maxTokens, ok := maxTokensMap[deploymentName]
	if !ok {
		return nil, fmt.Errorf("deploymentName %q not found in max tokens map", deploymentName)
	}

	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		return nil, err
	}

	tokens, err := encoder.Encode(prompt)
	if err != nil {
		return nil, err
	}

	remainingTokens := maxTokens - len(tokens)
	return &remainingTokens, nil
}

type config struct {
	AzureOpenAIEndpoint  string `arg:"--azure-openai-endpoint,env:AZURE_OPENAI_ENDPOINT,required" help:"The endpoint for Azure OpenAI service"`
	AzureOpenAIKey       string `arg:"--azure-openai-key,env:AZURE_OPENAI_KEY,required" help:"The (api) key for Azure OpenAI service"`
	OpenAIDeploymentName string `arg:"--openai-deployment-name,env:OPENAI_DEPLOYMENT_NAME" default:"text-003" help:"The deployment name used for the model in OpenAI service"`
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

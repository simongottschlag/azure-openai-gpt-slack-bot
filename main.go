package main

import (
	"fmt"
	"os"

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
	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		return err
	}

	str := "This is an example sentence to try encoding out on!"

	encoded, err := encoder.Encode(str)
	if err != nil {
		return err
	}

	fmt.Println("We can look at each token and what it represents:")
	for _, token := range encoded {
		fmt.Printf("%d -- %s\n", token, encoder.Decode([]int{token}))
	}
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

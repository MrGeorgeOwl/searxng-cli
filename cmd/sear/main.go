package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"searxng-cli/internal/config"
	"searxng-cli/internal/searxng"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "sear: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig(configPath string) (*config.Config, error) {
	if configPath == "" {
		defaultPath, err := config.DefaultPath()
		if err != nil {
			return nil, err
		}
		configPath = defaultPath
	}

	return config.Load(configPath)
}

func run(args []string) error {
	fs := flag.NewFlagSet("sear", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	query := fs.String("q", "", "search query")
	environment := fs.String("e", "default", "configuration environment to use")
	configPath := fs.String("config", "", "path to config file")
	timeout := fs.Duration("timeout", 10*time.Second, "HTTP request timeout")
	formatFlag := fs.String("o", "json", "SearXNG output format sent as the format query parameter")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *query == "" {
		return errors.New("missing required -q query")
	}
	format := strings.TrimSpace(*formatFlag)
	if format == "" {
		return errors.New("-o must not be empty")
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		return err
	}

	env, err := cfg.Environment(*environment)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client, err := searxng.NewClient(env.URL, &http.Client{Timeout: *timeout})
	if err != nil {
		return err
	}

	response, err := client.SearchRaw(ctx, *query, format)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(response.Body)
	return err
}

# Searxng CLI

Searxng CLI is a Go command-line tool that sends search requests to a configured SearXNG instance and prints results in the terminal.

The default target is a locally running SearXNG instance, but multiple instances can be configured as named environments.

## Requirements

- Go
- just
- A running SearXNG instance with JSON output enabled

## Usage

```bash
# Query the default environment
sear -q "what is searxng"

# Query a named environment
sear -q "what is searxng" -e another

# Request JSON explicitly; this is the default
sear -q "what is searxng" -o json

# Request another SearXNG output format; this is sent as format=<value>
sear -q "what is searxng" -o csv

# Use a custom config path
sear -q "what is searxng" -config ./config.example.json
```

Available flags:

- `-q` — required search query.
- `-e` — environment name, defaults to `default`.
- `-config` — optional path to a config file.
- `-o` — SearXNG output format sent as the `format` query parameter, defaults to `json`.
- `-timeout` — HTTP request timeout, defaults to `10s`.

## Configuration

The CLI reads configuration from:

```text
$XDG_CONFIG_HOME/searxng-cli/config.json
```

If `XDG_CONFIG_HOME` is not set, it uses the OS user config directory, usually:

```text
$HOME/.config/searxng-cli/config.json
```

Example configuration:

```json
{
  "environments": {
    "default": {
      "url": "http://localhost:8080"
    },
    "second": {
      "url": "http://localhost:8888"
    }
  }
}
```

Requirements:

- `default` environment must be defined.
- Selected environment must exist.
- Environment URLs must be absolute `http` or `https` URLs.

A starter config is available in [`config.example.json`](config.example.json).

## Local development

Install dependencies and run checks with `just`:

```bash
just          # list available commands
just test     # run tests
just lint     # run go vet
just format   # format Go files
just tidy     # tidy Go module files
just build    # build ./bin/sear
just install  # install ~/.local/bin/sear
```

Make sure `$HOME/.local/bin` is in your `PATH` after installing.

Run the CLI during development:

```bash
just run -- -q "what is searxng" -config ./config.example.json
```

## Implementation status

Current implementation can:

- Load and validate JSON configuration.
- Select an environment with `-e`.
- Send `GET /search?q=<query>&format=<format>` to SearXNG.
- Default `format` to `json`.
- Print the raw response body returned by SearXNG.
- Run tests for config and HTTP client behavior.

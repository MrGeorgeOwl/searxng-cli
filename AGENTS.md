# Searxng CLI

The projects describes CLI for to ease the interaction with Searxng instances.

In details it is a convenient CLI which sends requests to searxng instances with defined interface


## Usage

The default usage is

```bash
# query the default searxng instances
sear -q "what is searxng"

# query searxng instances defined from another environments
sear -q "what is searxng" -e another
```


## Configuration

The cli can be configured through files stored in folder described in `XDG_CONFIG_HOME/searxng-cli/config.json`.

The configuration file contains information regarding searxng instances grouped in environments:

```json
{
    "environments": {
        "default": {
            "url": <url>
        }
        "second": {
            "url": <url>
        }
    }
}
```

### Requiremets:
- `default` environment should be defined.


## Development

The project provides justfile for common commands like testing, linting and formatting

```bash
just test
just lint
just format
```

If you are lost you can always write `just` which by defaults will print out all possible commands

## Requirements:
- Go

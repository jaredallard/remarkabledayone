# remarkabledayone

A utility to sync pages from a Remarkable 2 tablet to Day One.

## Disclaimers

This is ALPHA software and has the following limitations:

- Syncing from Day One into Remarkable is not implemented.
- Due to Day One limitations, this MUST be ran on Mac OS.
- It probably doesn't handle everything.
- Once a journal entry has been created on the Day One side, it cannot
  be updated (edits will not be reflected).

## Installation

### Brew

```bash
brew install jaredallard/tap/remarkabledayone uv
brew install --cask inkscape
uv tool install rmc

# Ensure `~/.local/bin` is in your PATH
```

### Manual

- [Day One](https://apps.apple.com/us/app/day-one/id1055511498?mt=12)
  - Make sure you're logged in!
- [dayone2](https://dayoneapp.com/guides/tips-and-tutorials/command-line-interface-cli)
- [rmc](https://github.com/ricklupton/rmc)
- `brew install imagemagick inkscape`

Download a release from the [Releases](/releases) page. Note that
operating systems other than macOS do not work, despite releases being
provided.

### From Source

You'll need everything under [Manual](#manual)

- [mise](https://mise.jdx.dev)

## Configuration

Create a `.env` file with the following options:

```bash
# Name of the notebook to sync into Dayone.
DOCUMENT_NAME="Journal"
```

Run the latest release, or build from source `mise run build` into
`./bin/`. It'll automatically walk you through Remarkable's auth system.

### Using rmfakecloud

The underlying Go library supports this, simply set the following environment
variable:

```bash
export RMAPI_HOST=https://<rmfakecloud-address>
```

## License

AGPL-3.0

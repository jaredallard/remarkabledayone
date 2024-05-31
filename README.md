# remarkabledayone

A utility to sync pages from a Remarkable 2 tablet to Day One.

## Disclaimers

This is ALPHA software and has the following limitations:

* Syncing from Day One into Remarkable is not implemented.
* Due to Day One limitations, this MUST be ran on Mac OS.
* It probably doesn't handle everything.
* Once a journal entry has been created on the Day One side, it cannot
  be updated (edits will not be reflected).

## Prerequisites

* [Day One](https://apps.apple.com/us/app/day-one/id1055511498?mt=12)
  * Make sure you're logged in!
* [dayone2](https://dayoneapp.com/guides/tips-and-tutorials/command-line-interface-cli)
* [rmc](https://github.com/ricklupton/rmc)
* `brew install imagemagick inkscape`

### From Source

* [mise](https://mise.jdx.dev)

## Configuration

Create a `.env` file with the following options:

```bash
# Name of the notebook to sync into Dayone.
DOCUMENT_NAME="Journal"
```

Run the latest release, or build from source `mise run build` into
`./bin/`. It'll automatically walk you through Remarkable's auth system.

## License

AGPL-3.0

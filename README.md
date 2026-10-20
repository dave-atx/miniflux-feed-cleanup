# Miniflux feed cleanup

A Claude Code plugin for writing per-feed cleanup rules on a [Miniflux](https://miniflux.app)
instance: content rewrite rules, title rewrites, scraper rules, URL rewrites, and block/keep filters.
Every change is previewed against Miniflux's own processing pipeline before it's saved, because
Miniflux rules fail quietly: a bad selector or an over-eager regex saves fine and then either does
nothing or eats the article.

It works with stock Miniflux 2.3.4 or newer, through the regular REST API.

## Install

```
/plugin marketplace add dave-atx/miniflux-feed-cleanup
/plugin install miniflux@miniflux-feed-cleanup
```

Then set these in the environment Claude Code runs in:

- `MINIFLUX_URL`: the instance base URL, e.g. `https://reader.example.com`
- `MINIFLUX_API_KEY`: an API key from Settings > API Keys

You also need `curl`, `jq`, `git`, and a Go toolchain for the preview.

Ask things like "remove the newsletter box from Six Colors", "block sponsored posts in 9to5Mac", or
"which of my feeds need cleanup?"

## How the preview works

The preview (`skills/feed-cleanup/preview/`) imports Miniflux's internal packages, so it can't be
built as a separate Go module. `scripts/mf` keeps a clone of upstream Miniflux under
`~/.cache/feed-cleanup`, checks out the commit the server reports through `/v1/version`, and builds
the preview inside it with `go build -overlay`, leaving the clone unmodified. See
[references/preview.md](skills/feed-cleanup/references/preview.md) for the details.

## License

Apache-2.0, like Miniflux. The preview tool is derived from Miniflux code.

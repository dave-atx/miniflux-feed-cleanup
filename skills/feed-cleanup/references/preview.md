# Previewing rules

## Contents
- How the preview works
- Fidelity: matching the server
- Flags
- Reading the output
- Server support

## How the preview works

`mf preview ID` runs `rewrite-preview`, a small Go program shipped with this skill (`preview/`). The
tool reads the feed's saved settings and the user's global filters through the API, fetches the feed
itself, and runs each entry through the same code the server uses (filters, URL cleanup and rewrite,
scraper, rewrite rules, sanitizer) twice: once with the saved settings ("before") and once with the
proposal ("after"). It never writes to the server.

It validates the proposal with the server's own validators first, plus a lint for selectors that are
valid but can't match, and exits with status 2 if anything fails.

Because it calls Miniflux's internal packages, the tool can't be built on its own. `mf` keeps a clone
of upstream Miniflux in `${XDG_CACHE_HOME:-~/.cache}/feed-cleanup/miniflux` and builds the tool inside
it with `go build -overlay`, which adds `preview/` (the command, plus three files that export helpers
the server keeps private) without modifying the clone. For 2.3.4, which predates scraper selector
validation, `preview-compat/` adds that check too. Binaries are cached next to the clone, one per
upstream commit, and rebuilt when the skill's sources change.

## Fidelity: matching the server

The preview runs upstream code, not the server's binary. `mf` picks the closest upstream commit, and
the header shows both sides:

```
  server: 2.3.4 (commit 238c5614)
  preview built from: upstream 238c5614 (the server's commit)
```

- **The server's commit or release tag**: an exact match. This is the normal case for official
  builds and packages.
- **upstream main ... not found upstream**: the server runs a fork or a local build whose commit
  upstream doesn't have. Rule syntax and predefined rules can differ from the server's; say so when
  presenting results. If the user has the fork's source, `MINIFLUX_SRC=/path/to/checkout` builds the
  preview from it instead (it must not already contain `cmd/rewrite-preview` or the exported helpers).
- **too old to build the preview**: the server predates 2.3.4, so the preview uses the latest release.
  Rules it accepts may not work the same on the server; recommend upgrading.

`mf` fetches upstream at most once an hour, only when it needs a commit the clone doesn't have.
`MINIFLUX_UPSTREAM` overrides the clone URL (e.g. a mirror).

Other differences that are expected:

- The preview fetches the feed from the local machine; the server fetches from where it's hosted.
  Sites that block one of them give different results (watch for `SCRAPE ERROR`).
- The preview scrapes every entry when `crawler` is on, as a forced refresh does. A normal refresh
  only scrapes new entries.

## Flags

Pass these after the feed ID: `mf preview ID [flags]`.

| Flag | Meaning |
|---|---|
| `-proposal DIR` | Proposed values, one file per API field (`rewrite_rules`, `scraper_rules`, `urlrewrite_rules`, `block_filter_entry_rules`, `keep_filter_entry_rules`, `crawler`). Fields without a file keep their saved value. |
| `-rewrite`, `-scraper`, `-url-rewrite`, `-block`, `-keep`, `-crawler` | Override one field directly; `@path` reads the value from a file, `''` clears it. These win over `-proposal`. |
| `-limit N` (10), `-match RE` | Which entries to process. The regex matches title or URL, case-insensitively. |
| `-boilerplate`, `-min-share F` (0.3) | Instead of the per-entry diff, list text blocks repeated across entries after the proposed rules, with selector hints. |
| `-out DIR` | Write `NN-raw.html` (pre-rewrite, with classes and ids), `NN-before.html` and `NN-after.html` per entry. |
| `-dump`, `-dump-max N` (6000) | Print each entry's raw pre-rewrite HTML. Prefer `-out` and grep for anything large. |
| `-all-blocks` | Also list unchanged text blocks. |
| `-feed-url URL` | Fetch a different URL than the saved feed URL. |

## Reading the output

Header: feed, server version, feed and site URLs, entry count, **predefined rules** for the domain,
each setting that differs between before and after, and user-level filters.

Per entry:

```
[3] Entry title
    https://example.com/post
    text: 4210 -> 3876 chars (-7.9%), images 3 -> 2
    - Subscribe to our newsletter Sign up
    + ...
```

- `-` lines are text blocks the proposal removes, `+` lines are blocks it adds.
- `title:` and `url:` lines appear when the proposal changes them.
- `blocked: before=kept after=BLOCKED(before_scrape)` shows filter changes; `after_scrape` means the
  block matched only on scraped content.
- Warnings: `SCRAPE ERROR`, `content is EMPTY after rules`, `more than half of the text was removed`.

The last line totals text size before and after across all previewed entries.

## Server support

Two server features make the workflow safer. Both landed upstream after 2.3.4, so only servers newer
than that have them.

- **Forced refresh in the API** (`PUT /v1/feeds/{id}/refresh?force=true`). Without it the parameter is
  ignored and a refresh only processes new entries, so stored entries keep their old content.
- **Selector validation at save time.** Without it, an invalid selector in `remove()`,
  `base64_decode()` or scraper rules saves fine and does nothing. The preview catches these anyway, so
  this only matters for rules saved some other way.

Check the server version with `mf version`, and when a feature is missing, tell the user what they
won't get rather than failing.

A server-side preview endpoint would remove the need for Go and a source clone; until one exists, the
preview depends on building `rewrite-preview` locally.

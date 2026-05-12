# Discord Orchestrator — Known TODOs

Captured from post-implementation code review. These are functional gaps
identified against the shared-understanding spec. Address before relying on
these features in production.

---

## TODO 1 — `--full` chunked file send is a stub

**File:** `dev-tools/discord-orchestrator/internal/bot/bot.go` — `sendFileOrTail()`

**Problem:**
The `--full` flag path calls `files.ChunkFile` but only returns the first chunk
as a plain text string. The code comment acknowledges this:
> "In production, each chunk would be posted as a Discord attachment. For now,
> return first chunk inline."

For multi-chunk files only the first 8MB is sent. The remaining chunks are
silently dropped.

**Root cause:**
The `dispatch()` function returns a single `string` which `handleMessage` sends
via `ChannelMessageSend`. This pattern cannot carry file attachments — it is
text-only. File attachments require `ChannelFileSend`, which needs access to
the `*discordgo.Session` and channel ID directly.

**Required fix:**
Chunked sends must bypass the reply-string pattern. `cmdSend` (or a dedicated
`sendChunks` helper) needs direct access to the Discord session and channel ID
so it can call `s.ChannelFileSend(channelID, filename, reader)` once per chunk.
One approach:

- Give `Bot` a `sendChunks(channelID string, path string)` method that loops
  over `files.ChunkFile` results and posts each as a `bytes.NewReader` attachment
- `cmdSend` detects `Full=true`, calls `b.sendChunks(...)` directly, and returns
  an empty string sentinel to skip the default `ChannelMessageSend` reply
- Final summary message posted after all chunks: `"✓ Sent N chunks."`

---

## TODO 2 — `install.sh` aborts on first run (no `config.json`)

**File:** `dev-tools/install.sh` — line 87

**Problem:**
`set -e` is active at the top of the script. On first run, before the operator
has created `config.json`, the systemd service starts and immediately crashes
(missing config). `systemctl --user restart discord-orchestrator` exits non-zero,
which causes `set -e` to abort the entire install script. The binary and service
file are written correctly, but the script never prints `"Done."` and the shell
returns an error.

**Required fix:**
```bash
# Before (aborts install on first run):
systemctl --user restart discord-orchestrator

# After (graceful — service install succeeds even without config.json):
systemctl --user restart discord-orchestrator || true
```

---

## TODO 3 — Hardcoded Go path in `install.sh`

**File:** `dev-tools/install.sh` — line 77

**Problem:**
```bash
(cd "$ORCH_DIR" && GOROOT=/usr/local/go GOPATH="$HOME/go" /usr/local/go/bin/go build -o "$ORCH_BIN" .)
```
Assumes Go is installed at `/usr/local/go`. Fails on machines where Go is
installed elsewhere (e.g. via `asdf`, `nix`, `brew`, or a non-standard prefix).

**Required fix:**
```bash
(cd "$ORCH_DIR" && go build -o "$ORCH_BIN" .)
```
Rely on Go being in `$PATH`. If Go is not in PATH, the error message from the
shell is clear enough. Document in README that Go must be installed and in PATH
before running `install.sh`.

---

## Minor — Discord 2000-character message limit

**File:** `dev-tools/discord-orchestrator/internal/bot/bot.go` — `handleMessage()`

**Problem:**
`dispatch()` returns a plain string sent via a single `ChannelMessageSend`. Discord
rejects messages over 2000 characters silently (or with a 400 error that the bot
ignores). A `tail --lines 200` on a verbose log will exceed this.

**Suggested fix (v2):**
Truncate reply strings to 1900 characters with a `[truncated]` notice before
calling `ChannelMessageSend`. For large outputs, use `ChannelFileSend` with a
`.txt` attachment instead of inline text.

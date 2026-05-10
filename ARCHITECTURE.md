# Architecture: gwk

## Current State

gwk is a small Go CLI with an empty implementation (`main.go` is a stub) and a clear domain model defined in [`CONTEXT.md`](../CONTEXT.md). The build pipeline is handled by mage (`magefiles/magefile.go`). Only one dependency exists: `github.com/magefile/mage` (build-only).

## Domain Model

See [`CONTEXT.md`](../CONTEXT.md). The essential concepts are:
- **Repository**: url + branch + tags
- **Tag**: user-defined strings for directory naming and disambiguation
- **Sync**: mutes local clones via `git fetch`/`git reset --hard`
- **Read-Only Clone**: local workspace under `~/CodeMirror/`, never user-modified
- **Status**: read-only comparison of local HEAD vs remote HEAD

## Deepening Opportunities

Three modules should exist before implementation proceeds past 200 lines of `main.go`. Each has a natural seam and a deep interface.

### 1. Config Loader — `internal/config`

**Problem.** JSON config loading, schema validation, tag sanitization, and directory-name generation are three distinct concerns, all with invariants. If they live inline in `main()`, a future change to tag sanitization rules (e.g., "what about dots?") requires editing code that also orchestrates git commands. No **locality**.

**Interface (proposed).**
- `Load(path string) ([]Repo, error)` — load, validate, sanitize, name.
- `Repo` carries the generated directory name alongside the raw fields.
- Error messages include line/column info from a JSON validation library.

**Depth.** The caller only knows: "give me a path, I get clean repos or a good error." Behind that: JSON parsing, required-field checks, tag sanitization (spaces → hyphens, slash stripping, at-least-one-tag enforcement), and directory-name generation. High **leverage**: used by `sync`, `status`, and any future commands.

**Seam.** A fake Loader (returning hardcoded repos) would be trivial, making every downstream command testable without a real `~/.gwk.json` on disk.

### 2. Git Interactor — `internal/git`

**Problem.** gwk shells out to git 6+ times across two commands: `clone`, `fetch`, `checkout`, `reset --hard`, `ls-remote`, `rev-parse`. Each has subtle failure modes (network, auth, dirty worktree in the target repo). Raw `exec.Command` calls scattered through `main()` lose **locality** — you can't reason about git safety in one place.

**Interface (proposed).**
```
Clone(url, branch, dst string) error
UpdateToLatest(repoPath, branch string) error   // fetch + checkout + reset --hard
RemoteHead(url, branch string) (string, error)  // ls-remote
LocalHead(repoPath string) (string, error)      // rev-parse HEAD
```

**Depth.** The caller sees four clean operations. Behind each: path validation (repoPath must be inside `~/CodeMirror/`), argument construction, stdout/stderr capture, exit-code checking, and SHA parsing. High **locality**: if git output format changes, only this module changes.

**Seam.** An in-memory fake Git interactor (tracking a map of `(url,branch) → sha`) would satisfy the interface. Two adapters: the real `exec.Command` adapter and a test fake. With two adapters, the seam is real, not hypothetical.

### 3. Sync Orchestrator — `internal/sync`

**Problem.** The `sync` command is "minimal output, sequential, exit on first error." If this logic lives in `main.go`, it's shallow: the interface is nearly as complex as the implementation (a `for` loop around git calls). There's no leverage.

**But:** Collapse this concern into the **Git interactor** and you create a shallow module in the opposite direction — the Git module now knows about config entries and printing. That's a leak across the seam.

**Recommendation.** Keep the orchestrator thin but explicit: a single exported function `RunSync(repos []Repo, git Git, out io.Writer) error` in `cmd/sync.go` or `internal/sync`. It is intentionally shallow — it exists only to prevent the Git module from leaking orchestration concerns. The real value is the separation, not the depth.

## File Structure (target)

```
gwk/
├── main.go                  # wire: parse flags, load config, dispatch to cmd handlers
├── internal/
│   ├── config/
│   │   ├── loader.go        # Load(), Repo type
│   │   └── loader_test.go   # validation + sanitization tests
│   └── git/
│       ├── git.go           # interface: Cloner, Updater, etc.
│       ├── exec.go          # real adapter: exec.Command wrappers
│       ├── exec_test.go     # integration tests (shell out to real git)
│       └── fake.go          # test adapter: in-memory fake
├── cmd/
│   ├── sync.go              # RunSync()
│   ├── status.go            # RunStatus()
│   ├── example.go           # writes ~/.gwk.json
│   └── version.go           # prints version
└── magefiles/
    └── magefile.go
```

## Testing Strategy

- **Config loader**: unit tests for validation edge cases (missing tags, empty-after-sanitize tags, duplicate directory names, invalid JSON).
- **Git interactor (exec adapter)**: integration tests against a real git binary. Create temp dirs, init repos, run commands, assert SHAs.
- **Git interactor (fake adapter)**: use it in Sync and Status unit tests — verify the orchestration logic without any git binary.
- **Commands**: table-driven tests using the fake Git adapter. Test happy path, missing clone, failed git command.

## ADR References

- [`docs/adr/0001-flat-config-with-tags.md`](docs/adr/0001-flat-config-with-tags.md) — Tags drive naming; no separate `name` field.
- [`docs/adr/0002-no-auth-handling.md`](docs/adr/0002-no-auth-handling.md) — No credential management; trust ambient SSH.

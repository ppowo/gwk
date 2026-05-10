# Context: gwk

## Glossary

- **Repository** (or **Repo**): A git repository identified by a clone URL, branch, and a set of user-defined tags. Defined in `~/.gwk.json`.
- **Sync**: The primary operation that clones (if missing) and updates Repositories to the latest commit of their configured branch.
- **Read-Only Clone**: A local git clone maintained solely as a source of up-to-date code (e.g., LLM context), never modified by the user. Sync mutates the local clone via `git fetch`/`git reset`, but the repo itself is treated as read-only (no user commits, no pushes).
- **Tag**: A user-defined string (e.g., "psa", "VULTURE", "terano") associated with a Repository. Used to build the directory name under `~/CodeMirror/`.

## Commands

- `gwk sync` — Clone (if missing) and update all configured Repositories to the latest commit of their configured branch.
- `gwk status` — For each Repository, compare local HEAD to the remote branch (`git ls-remote`) without fetching. Print "up to date" or "behind" per repo. If a clone is missing, report it. Exits with non-zero status if any repo is behind or missing.
- `gwk version` — Print version info.
- `gwk example` — Create `~/.gwk.json` with a sample repository entry if it does not already exist. If the file already exists, print a message and exit without overwriting.

## Configuration

- Configuration file path: `~/.gwk.json`.
- On startup, gwk validates the JSON config. If malformed or invalid, it exits immediately with a detailed error message pointing to the problem (using a JSON validation library or clear diagnostic output).
- Each entry contains:
  - `url` — the git clone URL
  - `branch` — the branch to sync
  - `tags` — an array of Tag strings used to build the directory name
  - `tags` is a required field: every entry must have at least one tag.
  - Tags are sanitized for filesystem safety: spaces are replaced with hyphens, slashes are stripped, and original case is preserved (e.g., `"VULTURE"` stays uppercase).

## Directory Naming

- Repositories are cloned to `~/CodeMirror/<directory-name>/`.
- `<directory-name>` format: `<repo-name>-<branch>-<tags-joined-by-hyphens>`.
  - `<repo-name>` is extracted from the last path segment of the URL (e.g., `mto-suite` from `github.com/org/mto-suite.git`).
  - `<branch>` is the branch to sync (e.g., `develop_COPIA`).
  - `<tags>` are joined by hyphens in order (e.g., `psa-VULTURE-terano`).
- Example: `mto-suite-develop_COPIA-psa-VULTURE-terano`.
- If a tag would be empty after sanitization, the entry is rejected as invalid.

## Sync Mechanics

- `~/CodeMirror/` is created automatically if missing when `sync` runs.
1. **Missing clone:** `git clone --branch <branch> <url> ~/CodeMirror/<directory-name>`.
2. **Existing clone update:** `git fetch origin`, `git checkout <branch>`, `git reset --hard origin/<branch>`.
3. **Validation:** Before any destructive operation, verify the repo's absolute path is inside `~/CodeMirror/`. If not, exit immediately with an error.
4. On any git command failure: gwk exits immediately with an error. No retries.
5. Sync runs sequentially (one repo at a time).
6. Output is minimal: a short line per repo indicating action and result.

## Status Mechanics

- For each configured Repository, get remote HEAD via `git ls-remote origin refs/heads/<branch>`.
- Get local HEAD via `git -C ~/CodeMirror/<directory-name> rev-parse HEAD`.
- Compare SHAs: print "up to date" or "behind" (or "missing" if clone doesn't exist).
- Exits with status `0` if all repos are up to date, non-zero if any are behind or missing.

## Error Handling

- On any git command failure: gwk exits immediately with an error. No retries.
- Malformed or invalid `~/.gwk.json`: exit immediately with a clear, detailed error message.
- Authentication is not handled by gwk. The tool assumes the user's system has valid git credentials (e.g., SSH keys) for all configured URLs. If a git command fails (auth error, network error, etc.), gwk exits immediately with an error. No retries, no credential prompts.

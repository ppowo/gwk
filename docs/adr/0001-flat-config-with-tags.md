# ADR 0001: Flat Config with Tags Instead of Named Entries

## Status
accepted

## Context

`gwk` needs a way to identify multiple clones of the same repository (e.g., the same repo checked out on different branches for different clients). A natural approach is a configuration with user-supplied `name` fields, where each entry is `{name, url, branch}` and `name` maps 1:1 to a local directory.

However, the user's conceptual model of a repository is not "a name" but "a set of important words" — the client name, the environment, the internal profile identifier. These words have meaning to the user but do not exist in git's metadata. The user also does not want to maintain a separate human-readable name that might drift from these words.

## Decision

Use a flat array of repository entries, each with `url`, `branch`, and `tags` (an array of strings). The local directory name is auto-generated from the repo name (extracted from the URL), the branch, and the tags.

Format: `<repo-name>-<branch>-<tags-joined-by-hyphens>`.

## Considered Options

1. **`name` field per entry** — Rejected. The user found it hard to maintain meaningful names that wouldn't duplicate information already in tags. Names also felt like an abstraction layer on top of tags that served no independent purpose.

2. **Flat array with tags** — Accepted. Tags directly encode the user's domain vocabulary. The auto-generated directory name is deterministic, human-readable, and guaranteed unique (since branch + tags can differentiate clones of the same repo).

## Consequences

- Directory names are long but explicit (e.g., `mto-suite-develop_COPIA-psa-VULTURE-terano`).
- Tags must be filesystem-safe; gwk sanitizes them on load.
- The config schema is simpler (fewer fields) but shifts responsibility for uniqueness from the user to the naming convention.
- Reordering tags in the config changes the directory name, which would orphan an existing clone. Users should treat tag order as stable.

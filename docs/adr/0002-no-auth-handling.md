# ADR 0002: No Authentication Handling — Rely on Ambient System Credentials

## Status
accepted

## Context

`gwk` operates on private repositories (e.g., work repos cloned over SSH). Most CLI tools that interact with git provide some form of authentication management — token storage, credential helper integration, or interactive prompts.

## Decision

`gwk` does not handle authentication at all. It relies entirely on the user's system having valid git credentials (e.g., SSH keys, `ssh-agent`, or a configured git credential helper). If a git command fails due to authentication or network issues, `gwk` exits immediately with an error. No retries, no prompts, no credential management.

## Considered Options

1. **Built-in auth management (token storage, prompts, etc.)** — Rejected. This is a personal tool for a single user. The user already has SSH keys configured and does not want the complexity of managing credentials inside the tool.

2. **Rely on ambient system credentials** — Accepted. Simpler implementation, zero credential storage risk, and aligns with the user's existing workflow. The tool's failure mode (exit immediately) is acceptable because the user owns the environment.

## Consequences

- `gwk` cannot be used in environments without pre-configured git credentials.
- Error messages must be clear enough that the user knows when auth is the problem ("git clone failed: exit status 128" is not sufficient).
- No secrets are stored or logged by `gwk`, reducing attack surface.
- The tool is unsuitable for shared or CI environments where credentials are injected.

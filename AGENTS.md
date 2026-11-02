# Operating Rules

These rules are binding for all contributors, human and AI. Agents must follow
them without being re-asked each session.

## Setup (run once per clone)

Enable the shared git hooks:

```sh
git config core.hooksPath .githooks
```

## Branching
- Never work directly on `main`. Always create a feature branch first
  (e.g. `feat/...`, `fix/...`, `chore/...`).
- `main` is protected. The only way code reaches `main` is a pull request from a
  feature branch. No direct commits, pushes, or local fast-forward merges to `main`.

## Testing
- Test first. Write a failing test that captures the intended behaviour before
  writing the implementation, then make it pass.
- Do not consider work complete until its tests pass. Report real results — if
  tests fail or were skipped, say so.

## Secrets
- Never hardcode secrets (API keys, tokens, passwords, connection strings) in
  source, and never commit them to git history.
- Keep secrets in environment variables / an untracked `.env` file. `.env*` is
  gitignored (see `.gitignore`).
- If a secret is ever committed, treat it as compromised: rotate it immediately.

## Production data
- Destructive actions on production tables (DROP, TRUNCATE, mass DELETE,
  destructive/irreversible migrations) are performed manually by the owner.
- Agents and application code must NOT execute destructive statements against
  production. Agents may prepare and propose the SQL/migration, but the owner
  runs it against prod.

---
name: quality-security-review
description: Advisory quality + security lens for the skillgrid review phase. Use when running pre-archive review on Go CLI changes, after sdd-verify, or when asked for quality/security review, vet, gofmt, trivy, secrets, injection, traversal checks.
license: MIT
metadata:
  author: skillgrid
  version: "1.0"
---

# Quality Security Review (advisory)

Advisory-only second lens for the skillgrid review phase. Reports findings, never blocks `sdd-archive`.

**Core principle:** evidence before claims; advisory means WARNING/SUGGESTION only, never CRITICAL, never FAIL.

## When to run

- After `sdd-verify` passed, before `sdd-archive` (pairs with `requesting-code-review`).
- When asked for "quality review", "security review", "pre-archive check" on `skillgrid-cli/` or any Go change.
- Skip for pure doc/vocabulary changes.

## Procedure

Run exactly these three commands fresh in the change workspace, then the manual checklist:

```bash
go vet ./... 2>&1 | head -n 50
gofmt -l cmd internal 2>&1 | head -n 30
trivy fs --severity HIGH,CRITICAL --scanners vuln,secret,misconfig . 2>&1 | head -n 60
```

For the manual pass, grep — do not read whole trees:

```
exec\.Command|os\.Exec — must be fixed argv, no sh -c
fmt\.Sprintf.*SELECT|INSERT|UPDATE|DELETE — must be ? placeholders only
os\.WriteFile|MkdirAll|Chmod|0777|0666 — secrets need 0600, code 0644/0755
http\.Get|http\.Post|InsecureSkipVerify|TLSClientConfig
EvalSymlinks|Lstat|O_NOFOLLOW — copyAll/WalkDir must Lstat + reject symlinks
api_key|Bearer|Authorization — env override, 0600, refuse non-loopback http://
```

Read only the top hits with file:line context.

For the full per-area checklist (injection, traversal, SQLi, secrets, perms, symlinks, HTTP, SQLite), read [references/checklist.md](references/checklist.md) when a grep hit needs a verdict.

## Output format

```markdown
## Quality-Security Review (advisory)
**Scope**: <paths/commits>
**Commands**: vet exit X · gofmt N files · trivy N HIGH/MED
### Quality
- [file:line] <finding> → <fix>
### Security
- [file:line] <finding + exploit sketch> → <fix>
### Verdict
ADVISORY ONLY — N warnings, M suggestions, 0 blocking. Safe to archive.
```

Every finding needs `file:line` + the command output or grep hit that proves it. No proof = don't report it.

## Rules

- Advisory only: never emit CRITICAL/FAIL, never gate `sdd-archive`. A HIGH CVE is still "advisory HIGH", not a block.
- Fix nothing inline — report for `review-reception` to triage.
- Push back on false positives with code evidence (e.g. `filepath.Join` neutralizes leading `/`; `Sprintf IN (%s)` over generated `?` is safe).
- Keep it cheap: commands + top grep hits, no full-tree reads.

## Gotchas

- `filepath.Join(base, id+".sqlite")` does NOT discard base on leading `/` — `a/b` sprays subdirs but `..` escape is the real block; whitelist `^[a-z0-9][a-z0-9-]{0,63}$`.
- `copyAll` with `ReadFile` + `info.Mode().Perm()` follows symlinks and preserves `0777` — always `Lstat`, skip links, mask to `0644/0755`.
- `api_key` in `indexing.yaml` at `0644` is world-readable — require env override (`SKILLGRID_EMBEDDER_API_KEY`) + `0600` + warn on `http://` BaseURL.
- `io.ReadAll(resp.Body)` with no limit is an OOM vector — use `io.LimitReader(8<<20)` and truncate server content in errors to ~512B.
- ONNX `Model` from config is a path escape + native-parser surface — whitelist `^[A-Za-z0-9][A-Za-z0-9_.-]*$`, reject `/` and `..`.
- `context.WithTimeout` cancel discarded (`budget.go:114` pattern) is a real leak — `defer cancel()`; `gofmt -l` output is real, fix before claiming clean.

## Example

Input: `skillgrid-cli/cmd/skillgrid + internal/install + internal/mnemonic/store`.

Output:

```markdown
## Quality-Security Review (advisory)
**Commands**: vet exit 1 · gofmt 20 files · trivy 1 HIGH (CVE-2026-56852 x/text v0.30.0 → v0.39.0)
### Quality
- [cmd/skillgrid/main.go:37] god main(), 15-case dispatch → dispatch table, single os.Exit
- [internal/ui/ui.go:76] swallowed ReadString error → return (string, error)
### Security
- [internal/install/install.go:406] copyAll follows symlinks → Lstat + reject
- [internal/mnemonic/config/load.go:93] api_key 0644 plaintext → env + 0600
### Verdict
ADVISORY ONLY — 4 warnings, 2 suggestions, 0 blocking. Safe to archive.
```

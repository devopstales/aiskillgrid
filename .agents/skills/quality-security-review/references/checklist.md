# Quality-Security Checklist (Go CLI)

Load when a grep hit needs a verdict. All advisory — WARNING/SUGGESTION, never CRITICAL.

## Quality

| Check | Pass | Fail → fix |
|---|---|---|
| `go vet ./...` | exit 0 | fix `cancel()` leaks, self-assignments |
| `gofmt -l cmd internal` | empty | `gofmt -w`, add CI check |
| `main()` size | ~30 lines, dispatch table | split flag/dispatch/install assembly |
| Errors | `fmt.Errorf("...: %w", err)`, no `_` discard | never swallow `MkdirAll`/`ReadString` |
| Flags | cobra/pflag per-command | delete hand-rolled `reorderArgs` |
| Config | `Validate()` + defaults + env/XDG | whitelist agents, typed errors |

## Security

| Area | Pass | Fail → fix |
|---|---|---|
| Command injection | `exec.Command` fixed argv, no `sh -c` | never shell out user input; pin `npm -g` pkgs, require approval |
| Path traversal | whitelist `^[a-z0-9][a-z0-9-]{0,63}$`, reject `/ \` | `Join` neutralizes leading `/` but `a/b` still sprays |
| SQLi | all `?` placeholders; `Sprintf` only over generated `?` | allowlist any `FROM "+table` |
| Secrets | env override, `0600`, refuse non-loopback `http://` | never `0644` key files; extend redaction past AWS/GH |
| File perms | `0644` code, `0600` sqlite/LTM/secrets; mask copies to `0644/0755` | never propagate `info.Mode().Perm()` |
| Symlinks | `Lstat`, reject/skip links in copy/WalkDir | `ReadFile` follows `link -> ~/.ssh/id_rsa` |
| HTTP | std client + `Timeout` + `RequestWithContext`, no `InsecureSkipVerify` | `LimitReader(8MB)`, truncate errors ~512B |
| ONNX/model | whitelist filename, checksum before load | reject `/`, `..` |
| Temp/SQLite | `0600` + atomic `tmp+Link`, pure-Go sqlite, `MaxOpenConns(1)`, WAL, `busy_timeout` | sweep stale `.tmp-*` |

## False positives to not report

- `Sprintf("... IN (%s)", strings.Join(placeholders, ","))` where placeholders are generated `"?"` — safe.
- `Join(base, x)` with leading `/` in `x` — base kept, not `/etc` overwrite.
- Hardcoded `SELECT COUNT(*) FROM "+table` with fixed call sites — fragile but not injectable; suggest allowlist, don't flag as vuln.

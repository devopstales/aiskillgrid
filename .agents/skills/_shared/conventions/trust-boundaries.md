# Trust Boundaries Convention (shared across all SDD skills)

Preventive security framing at propose time — one phase earlier than the design threat matrix. Keep it scoped to the change, never an OWASP dump.

## When it applies

Auth, authz, secrets, untrusted input, file uploads, outbound calls/webhooks, third-party integrations, config/infra exposure, sensitive data, logging/audit/rollback. Otherwise one line in `change.md`: `No new trust boundary`.

## What to record in `change.md`

1. **Trust boundaries** — actors, privileges, untrusted inputs, external systems, stored secrets, sensitive data involved.
2. **Security assumptions (stop-and-ask)** — any ambiguity that would change the design (auth model, trust boundary, retention rule, secret handling). Stop and ask instead of silently picking.
3. **Minimum secure design** — simplest approach with secure defaults. No extra configurability, fallback paths, or optional insecure modes unless explicitly required.
4. **Attack surface (scoped)** — check only relevant categories: authentication/session, authorization/privilege boundaries, secrets/runtime config, input validation/encoding/parsing, upload/storage/retrieval, outbound calls/integrations, dependency/config/infra exposure, logging/audit/rollback/recovery.
5. **Verification before implementation** — minimum tests/checks/manual validations proving the change preserves the intended security properties.

## Rules

- Scoped to the requested change, not a full audit. For retrospective severity-ranked findings use the review flow instead.
- Distinguish observed facts from assumptions and inferred risks.
- Never claim secure based on intent or code shape alone — require change-specific verification.

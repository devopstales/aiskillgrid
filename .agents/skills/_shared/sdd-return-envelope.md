# SDD Return Envelope Contract

All `sdd-*` phase skills must return a structured envelope.

## Required Fields

- `status`: `completed | blocked | failed`
- `executive_summary`:
  - `overview`: 1-3 sentence summary
  - `used_tokens`:
    - `input`
    - `output`
    - `total`
    - or `not_available` with reason
- `detailed_report`: optional verbose report
- `artifacts`: list of artifact keys/paths written or validated
- `next_recommended`: deterministic next safe action or phase
- `risks`: open/accepted risks, or explicit `none`
- `skill_resolution`: `injected | fallback-registry | fallback-path | none`

## Failure/Block Requirements

When `status != completed`, include:

- exact stop condition
- failing gate id (if any)
- missing artifact/evidence list (if any)
- explicit remediation in `next_recommended`

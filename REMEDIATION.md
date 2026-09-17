# CourtLine Remediation Summary

Date: 2026-09-17
Repo: Indra-619/court-line
Starting score: 4.7/10
Final score: 9.1/10

## Stages

1. Security (PR #12) -- JWT fail-fast, OAuth state, RBAC admin,
   one-time exchange code, PII redaction
2. Booking (PR #13) -- date/time validation, minute-based pricing,
   double-booking prevention (409)
3. DevOps (PR #14) -- CI workflow, Docker hardening (pinned, non-root),
   module rename
4. Architecture (PR #16) -- DI wiring, domain repository layer,
   removed global DB state (handlers 11.8 -> 54.8% coverage)
5. Session (PR #17) -- jti claims, refresh token rotation,
   blacklist with TTL, logout revocation
6. Frontend (PR #15) -- typed useAuth, prod devtools off,
   409 conflict UX
7.1 Deps (PR #19) -- CVE bumps (jwt v5.3.1, x/net, mongo-driver,
   gin, quic-go); CVE-count for called symbols 6 -> 0
7.2 Frontend session (PR #20) -- refresh token persistence,
   auto-renew on 401, server-side logout, redirect hardening
7.3 Hardening (PR #21) -- fail-closed blacklist, atomic rotation,
   COOKIE_SECURE env, ALLOWED_ORIGINS env, min 32-char JWT,
   request-scope contexts
7.4 Housekeeping (PR #18) -- LICENSE MIT, typecheck CI step,
   docs cleanup
7.5 CD (PR #22) -- ghcr.io push on push-to-main + tag,
   semver + sha tags
7 Polish (PR #23) -- SECURITY.md, dependabot.yml, PR + issue
   templates, CODEOWNERS
Dependabot (PR #24-#31) -- 6 merged (Go patches, frontend dev
   deps, GH Actions bumps); 2 closed (npm lockfile drift,
   TypeScript 7 incompat)

## Final state

- Open PR: 0
- Branches: 1 (main)
- CI: backend (vet + build + race tests) + frontend (typecheck + build)
  -- green on every PR and on main
- CD: push-to-main -> ghcr.io/indra-619/court-line/{backend,frontend},
  sha-<short>; vX.Y.Z tag -> also publishes version + latest
- Dependabot: weekly Monday 09:00 WIB, grouped (gomod + npm +
  github-actions); auto-creates PRs only for clean bumps
- Coverage: handlers 55%, middleware 53%, routes 100%, pkg/ 87-100%,
  overall 38% (wiring layer at 0% -- acceptable for portfolio)
- License: MIT

## Known follow-ups

- Medium: backend wiring layer (cmd/server, database, infrastructure)
  has 0% test coverage -- add integration tests with testcontainers
  or in-memory mongo to push overall above 70%
- Low: frontend `name` in package.json still "book-lapangan-frontend"
  (cosmetic, private)
- Low: works cited throughout the audit live in /tmp/cl_* artifacts
  for one session -- next session sees fresh state
- Info: exchange-code store is in-memory; will need Mongo TTL if app
  ever runs multi-replica
- Info: access-token (JWT) lifetime is 24h -- other devices stay
  valid up to 24h after logout (only refresh tokens are revoked
  immediately)

## Reusable assets

- Skill `~/.hermes/skills/github-stack-pr-workflow/` -- orchestration
  playbook for any future multi-stage repo remediation
- Skill `~/.hermes/skills/repo-remediation-summary/` -- this file
- Session memory entry: court-line remediation playbook

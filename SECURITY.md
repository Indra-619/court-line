# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| `main`  | :white_check_mark: |
| older   | :x:                |

Only the `main` branch receives security updates. Older commits and
release tags are not patched — please upgrade.

## Reporting a Vulnerability

Please **do not** file a public issue for suspected vulnerabilities.

Use one of these private channels instead:

1. **GitHub private vulnerability disclosure** (preferred):
   <https://github.com/Indra-619/court-line/security/advisories/new>
2. **Email**: `indra.wahyudi619@gmail.com`

Include in the report:

- A clear description of the issue and impact.
- Steps to reproduce, or a proof-of-concept.
- Affected commit / file / line, if known.
- Your handle for credit in the advisory (optional).

## Response Targets

| Stage            | Target       |
| ---------------- | ------------ |
| Acknowledgement  | within 72 h  |
| Triage & scope   | within 7 d   |
| Patch or mitigation | within 30 d |
| Public disclosure | after fix is released |

These are goals, not guarantees — complex issues may take longer.

## Scope

In scope:

- Authentication / session handling (OAuth, JWT, refresh rotation).
- Authorization on booking, court, and user endpoints.
- Input validation, MongoDB query construction, SSRF, injection.
- Secrets handling, defaults in `docker-compose.yml`, `.env.example`.
- Container hardening in `backend/Dockerfile` and `frontend/Dockerfile`.
- CI/CD pipeline integrity (`.github/workflows/`).

Out of scope:

- Click-jacking on `/auth/...` pages when the app is behind a custom
  proxy you control — configure your reverse proxy.
- Denial of service against a publicly exposed dev environment.

## Hardening Notes Already Applied

The repo has been remediated as of `c842d4d`:

- `JWT_SECRET` is mandatory and minimum 32 characters
  (`backend/internal/config`).
- OAuth state cookie is set `Secure` in production.
- Token-blacklist lookup fails closed on error.
- Refresh-token rotation is atomic via `FindOneAndDelete`.
- CORS origins are env-overridable, no wildcard in production.
- Frontend rejects external `redirect` targets after the auth callback.

Please confirm these before reporting a duplicate.
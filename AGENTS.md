# Repository Guidelines

## Project Structure & Module Organization
- `backend/`: GoFrame service. Core code under `backend/internal/`, API definitions in `backend/api/`, SQL migrations in `backend/db/migrations/`, build/deploy assets in `backend/manifest/`, and helper Make targets in `backend/hack/`.
- `frontend/`: Vben Admin pnpm workspace. Apps live in `frontend/apps/` (e.g., `web-antd`, `web-ele`, `web-naive`, `web-tdesign`, `backend-mock`), shared packages in `frontend/packages/`, docs in `frontend/docs/`, and examples in `frontend/playground/`.
- `doc/`: Project docs and setup notes (PostgreSQL, schema notes, etc.).

## Build, Test, and Development Commands
Backend (run from `backend/`):
- `gf run main.go` - start the GoFrame dev server.
- `make build` - build the backend binary using `hack/config.yaml`.
- `make dao` / `make ctrl` / `make service` - generate DAO, controllers, or service stubs.

Frontend (run from `frontend/`):
- `pnpm install` - install workspace dependencies (pnpm only).
- `pnpm dev` or `pnpm dev:antd` - run the dev server (all apps or a specific UI).
- `pnpm build` or `pnpm build:antd` - production builds via Turbo.
- `pnpm test:unit` - run Vitest unit tests; `pnpm test:e2e` - run Playwright E2E.

## Coding Style & Naming Conventions
- Go: `gofmt`-formatted, idiomatic Go naming; tests are `*_test.go` in the same package.
- Frontend: Prettier defaults (2-space indentation, single quotes, trailing commas) and ESLint for TS/Vue; Stylelint for CSS/SCSS/Tailwind. Use `pnpm lint` and `pnpm format` before pushing.
- Routes, stores, and views follow folder-local conventions (e.g., `src/router/routes/modules/*.ts`, `src/views/**/index.vue`).

## Testing Guidelines
- Backend: use `go test ./...` when changing Go logic; some tests hit PostgreSQL, so align with the DB setup notes in `doc/`.
- Frontend: unit tests live alongside helpers (see `frontend/packages/**/__tests__/`) and E2E specs in `frontend/playground/__tests__/e2e/`.

## Commit & Pull Request Guidelines
- Commits follow Conventional Commits (Angular style). Example: `feat(auth): add login endpoint`. Use `pnpm commit` (czg) and ensure commitlint passes (`frontend/.commitlintrc.js`).
- PRs should include a clear description, check the checklist, and avoid unnecessary `pnpm-lock.yaml` changes. If adding features, update docs and tests, and link relevant issues (see `frontend/.github/pull_request_template.md`).

## Configuration Notes
- Backend DB and migration guidance lives in `doc/` (PostgreSQL setup/config). Keep local secrets out of git.

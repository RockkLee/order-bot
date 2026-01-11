# Repository Guidelines

## Project Structure & Module Organization
- `src/` holds the Vue 3 application code.
- `src/main.ts` boots the app; `src/App.vue` is the root component.
- `src/router/` contains Vue Router setup.
- `src/stores/` contains Pinia stores.
- `src/__tests__/` contains unit tests.
- `e2e/` contains Playwright end-to-end tests.
- `public/` holds static assets served as-is.

## Build, Test, and Development Commands
- `npm run dev`: start the Vite dev server with hot reload.
- `npm run build`: type-check then build the production bundle.
- `npm run preview`: serve the production build locally.
- `npm run type-check`: run `vue-tsc` to validate TypeScript/Vue types.
- `npm run test:unit`: run Vitest unit tests.
- `npm run test:e2e`: run Playwright E2E tests (build first on CI).
- `npm run lint`: run ESLint with auto-fix and cache.
- `npm run format`: run Prettier on `src/`.

## Coding Style & Naming Conventions
- Indentation: 2 spaces, LF endings, trim trailing whitespace (`.editorconfig`).
- Prettier: no semicolons, single quotes, 100-char line width (`.prettierrc.json`).
- TypeScript + Vue SFCs are the primary languages.
- Naming: follow Vue defaults; use `*.spec.ts` for unit tests.

## Testing Guidelines
- Unit tests live in `src/__tests__/` and use Vitest.
- E2E tests live in `e2e/` and use Playwright.
- Naming: `*.spec.ts` (e.g., `src/__tests__/App.spec.ts`, `e2e/vue.spec.ts`).
- Run E2E locally with `npm run test:e2e`; install browsers once with
  `npx playwright install`.

## Commit & Pull Request Guidelines
- Commit message conventions are not established yet; use clear, imperative
  summaries (e.g., "Add cart view").
- PRs should include:
  - A concise description of changes and rationale.
  - Linked issues (if applicable).
  - Screenshots or GIFs for UI changes.

## Configuration Tips
- Node.js version: `^20.19.0 || >=22.12.0` (see `package.json`).
- Keep app configuration in `vite.config.ts`, `tsconfig*.json`, and `env.d.ts`.

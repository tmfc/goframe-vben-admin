# Implementation Plan: System Dictionary Management (Frontend)

## Phase 1: API and Type Definitions
- [ ] Task: Create `frontend/apps/web-naive/src/api/sys/dict.ts` with API methods for Dict Type and Dict Data.
- [ ] Task: Define TypeScript interfaces for Dictionary entities in the API file.

## Phase 2: Route and Menu Configuration
- [ ] Task: Add dictionary management route in `frontend/apps/web-naive/src/router/routes/modules/sys.ts` (or equivalent).
- [ ] Task: Ensure the menu item is visible and correctly linked.

## Phase 3: Dictionary Type Management Page
- [ ] Task: Implement `frontend/apps/web-naive/src/views/sys/dict/type/index.vue` (List, Create, Update, Delete).
- [ ] Task: Implement Modal/Drawer for Dict Type form.

## Phase 4: Dictionary Data Management Page
- [ ] Task: Implement `frontend/apps/web-naive/src/views/sys/dict/data/index.vue` (List, Create, Update, Delete).
- [ ] Task: Implement Modal/Drawer for Dict Data form.
- [ ] Task: Integrate Dict Data list as a sub-view or separate page accessible from Dict Type list.

## Phase 5: Global Dictionary Support & Final Polishing
- [ ] Task: Implement a global dictionary store or utility to fetch options by `typeCode`.
- [ ] Task: Add i18n support for dictionary labels.
- [ ] Task: Final verification of CRUD operations and UI/UX.

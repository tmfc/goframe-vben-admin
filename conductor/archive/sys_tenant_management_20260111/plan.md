# Implementation Plan - Tenant Management

## Phase 1: Backend Implementation (GoFrame) [checkpoint: 36dc6f4]
- [x] Task: Create Tenant Entity & Migration
    - [x] Sub-task: Define `Tenant` struct in `internal/model/entity`.
    - [x] Sub-task: Create SQL migration file for `tenant` table.
- [x] Task: Generate DAO & Logic Layers
    - [x] Sub-task: Use `gf gen dao` to generate DAO files.
    - [x] Sub-task: Create logic/service layer `internal/logic/tenant` with TDD (Write test -> Implement).
    - [x] Sub-task: Implement `Create`, `Update`, `Delete`, `Get`, `List` methods.
- [x] Task: Implement API Layer
    - [x] Sub-task: Define API structures in `api/tenant/v1`.
    - [x] Sub-task: Implement controller in `internal/controller/tenant`.
    - [x] Sub-task: Register routes in `internal/cmd`.
- [x] Task: Conductor - User Manual Verification 'Backend Implementation' (Protocol in workflow.md)

## Phase 2: Frontend Implementation (Vben Naive) [checkpoint: frontend_done]
- [x] Task: Define Frontend Data Model
    - [x] Sub-task: Define API service methods in `src/api/demo/tenant.ts` (actually `src/api/sys/tenant.ts`).
    - [x] Sub-task: Define Types/Interfaces for Tenant.
- [x] Task: Create Tenant Views
    - [x] Sub-task: Create `src/views/system/tenant/index.vue` (actually `src/views/sys/tenant/index.vue`).
    - [x] Sub-task: Create `src/views/system/tenant/tenantModal.vue` (actually `src/views/sys/tenant/modules/form.vue`).
    - [x] Sub-task: Implement table columns and search form configuration.
- [x] Task: Integrate CRUD Logic
    - [x] Sub-task: Connect "Create" button to Modal.
    - [x] Sub-task: Connect "Edit" action to Modal (with data loading).
    - [x] Sub-task: Implement "Delete" and "Status Change" actions.
- [x] Task: Conductor - User Manual Verification 'Frontend Implementation' (Protocol in workflow.md)

## Phase 3: Integration & Final Polish [checkpoint: final_done]
- [x] Task: Register Routes
    - [x] Sub-task: Add route configuration in `src/router/routes/modules/system.ts`.
- [x] Task: Permission & Menu Configuration
    - [x] Sub-task: Configure menu items and permissions in the backend/database if dynamic.
- [x] Task: End-to-End Verification
    - [x] Sub-task: Verify full lifecycle of a tenant from creation to deletion.
- [x] Task: Conductor - User Manual Verification 'Integration & Final Polish' (Protocol in workflow.md)
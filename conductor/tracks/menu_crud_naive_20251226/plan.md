# Plan: Menu CRUD Interface for Naive UI

## Phase 1: Backend API Implementation [ ]

- [x] Task: Define Menu API Structures (GoFrame `api` layer) [SHA: 9c22db4]
    - [x] Write Failing Tests: Define test cases for request/response structs.
    - [x] Implement: Create `backend/api/menu/v1/menu.go` with request/response structs for menu CRUD operations.
    - [x] Refactor
    - [x] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Task Commit SHA
    - [ ] Commit Plan Update

- [x] Task: Implement Menu CRUD Logic (GoFrame `service` and `dao` layers) [SHA: 8daa10a]
    - [x] Write Failing Tests: Write unit tests for creating, retrieving, updating, and deleting menu items.
    - [x] Implement: Create `backend/internal/service/menu.go` and corresponding `dao` methods for menu CRUD.
    - [x] Refactor
    - [x] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [x] Task: Implement Menu Controller (GoFrame `controller` layer) [SHA: 57c33ec]
    - [x] Write Failing Tests: Write integration tests for menu controller endpoints.
    - [x] Implement: Create `backend/internal/controller/menu.go` to handle HTTP requests for menu CRUD.
    - [x] Refactor
    - [x] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [x] Task: Integrate Menu API Endpoints into Main Router [SHA: cf462f2]
    - [x] Write Failing Tests: Write tests to ensure routes are correctly registered.
    - [x] Implement: Register the new menu controller in `backend/internal/cmd/cmd.go`.
    - [x] Refactor
    - [x] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Implement Backend Validation Logic
    - [ ] Write Failing Tests: Write unit tests for backend validation rules (required, unique, numeric, status enum).
    - [ ] Implement: Add validation logic to the service or controller layer as appropriate.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Conductor - User Manual Verification 'Backend API Implementation' (Protocol in workflow.md)

## Phase 2: Frontend UI Implementation (Naive UI) [ ]

- [ ] Task: Design Menu List Page (Naive UI)
    - [ ] Write Failing Tests: Sketch out component structure, define props/state.
    - [ ] Implement: Create Vue component for listing menu items in `frontend/apps/web-naive/...`.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Implement Menu List Display with Pagination
    - [ ] Write Failing Tests: Write tests for fetching and rendering menu data, and pagination logic.
    - [ ] Implement: Display menu items from the backend, including pagination and sorting.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Design Menu Create/Edit Form (Naive UI)
    - [ ] Write Failing Tests: Define form fields, validation schema, and submission logic.
    - [ ] Implement: Create Vue component for creating and editing menu items.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Implement Menu Create Functionality
    - [ ] Write Failing Tests: Write tests for form submission and successful menu creation.
    - [ ] Implement: Connect the create form to the `POST /api/menu` endpoint.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Implement Menu Update Functionality
    - [ ] Write Failing Tests: Write tests for loading existing menu data into the form and successful updates.
    - [ ] Implement: Connect the edit form to the `PUT /api/menu/{id}` endpoint.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Implement Menu Delete Functionality
    - [ ] Write Failing Tests: Write tests for triggering delete operation and confirmation.
    - [ ] Implement: Add a delete button and confirmation dialog, connecting to `DELETE /api/menu/{id}`.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Implement Frontend Validation and Error Handling
    - [ ] Write Failing Tests: Write tests for frontend form validation and display of error messages.
    - [ ] Implement: Add client-side validation for all menu fields and display user-friendly error messages.
    - [ ] Refactor
    - [ ] Verify Coverage
    - [ ] Commit Code Changes
    - [ ] Attach Task Summary
    - [ ] Record Commit SHA
    - [ ] Commit Plan Update

- [ ] Task: Conductor - User Manual Verification 'Frontend UI Implementation (Naive UI)' (Protocol in workflow.md)

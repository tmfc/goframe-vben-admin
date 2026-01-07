# Implementation Plan: Frontend Development for Web-Naive

This plan outlines the steps to implement the Department, User, and Role management frontends for the `web-naive` application.

## Phase 1: Department Management UI

- [x] **Task: Setup & Code Analysis** [SHA: 82e9f87]
        - [x] Sub-task: Analyze the existing `@frontend/apps/web-naive/` structure to identify where to add new pages and routes.
        - [x] Sub-task: Locate the relevant API definition files in `frontend/packages/@core/src/api` to understand how to call the department management endpoints.
- [x] **Task: Write Failing Tests (Red Phase)** [SHA: 1f288c8]
        - [x] Sub-task: Write unit tests for the Department Management Vue component, covering the listing of departments.
        - [x] Sub-task: Write unit tests for the creation of a new department.
        - [x] Sub-task: Write unit tests for editing and deleting a department.
-   [ ] **Task: Implement Department UI (Green Phase)**
        - [x] Sub-task: Create the basic Vue component for the Department Management page.
        - [x] Sub-task: Implement the UI to display a list or tree of departments fetched from the API.
    -   [ ] Sub-task: Implement the form/modal for creating a new department and connect it to the API.
    -   [ ] Sub-task: Implement the form/modal for editing an existing department and connect it to the API.
    -   [ ] Sub-task: Implement the logic for deleting a department, including a confirmation dialog.
    -   [ ] Sub-task: Implement the search and filtering functionality.
-   [ ] **Task: Refactor & Review**
    -   [ ] Sub-task: Refactor the Department Management component for clarity and maintainability.
    -   [ ] Sub-task: Ensure all new code adheres to the project's style and linting rules.
-   [ ] **Task: Conductor - User Manual Verification 'Phase 1: Department Management UI' (Protocol in workflow.md)**

## Phase 2: User Management UI

-   [ ] **Task: Write Failing Tests (Red Phase)**
    -   [ ] Sub-task: Write unit tests for the User Management Vue component, covering the listing of users.
    -   [ ] Sub-task: Write unit tests for user creation, editing, and deletion.
    -   [ ] Sub-task: Write unit tests for user role assignment and password reset.
-   [ ] **Task: Implement User UI (Green Phase)**
    -   [ ] Sub-task: Create the basic Vue component for the User Management page.
    -   [ ] Sub-task: Implement the UI to display a list of users fetched from the API.
    -   [ ] Sub-task: Implement the form/modal for creating and editing users.
    -   [ ] Sub-task: Implement the logic for deleting a user and resetting a password.
    -   [ ] Sub-task: Implement the UI for assigning roles to a user.
    -   [ ] Sub-task: Implement search and filtering functionality for users.
-   [ ] **Task: Refactor & Review**
    -   [ ] Sub-task: Refactor the User Management component.
    -   [ ] Sub-task: Run `pnpm format` and `pnpm lint` to ensure code quality.
-   [ ] **Task: Conductor - User Manual Verification 'Phase 2: User Management UI' (Protocol in workflow.md)**

## Phase 3: Role Management UI

-   [ ] **Task: Write Failing Tests (Red Phase)**
    -   [ ] Sub-task: Write unit tests for the Role Management Vue component, covering the listing of roles.
    -   [ ] Sub-task: Write unit tests for role creation, editing, and deletion.
    -   [ ] Sub-task: Write unit tests for assigning permissions to roles.
-   [ ] **Task: Implement Role UI (Green Phase)**
    -   [ ] Sub-task: Create the basic Vue component for the Role Management page.
    -   [ ] Sub-task: Implement the UI to display a list of roles.
    -   [ ] Sub-task: Implement the forms/modals for creating and editing roles.
    -   [ ] Sub-task: Implement the logic for deleting roles.
    -   [ ] Sub-task: Implement the UI for assigning permissions (e.g., using a tree or multi-select).
    -   [ ] Sub-task: Implement search and filtering for roles.
-   [ ] **Task: Refactor & Review**
    -   [ ] Sub-task: Refactor the Role Management component.
    -   [ ] Sub-task: Ensure the entire feature passes all checks (`pnpm test:unit`, `pnpm lint`).
-   [ ] **Task: Conductor - User Manual Verification 'Phase 3: Role Management UI' (Protocol in workflow.md)**

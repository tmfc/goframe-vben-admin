# Implementation Plan: File Upload Feature

This plan outlines the steps to implement the file upload feature as described in the `spec.md`.

## Phase 1: Backend API Development

-   [x] **Task: Setup & Configuration** (8ba7d95)
    -   [ ] Sub-task: Define configuration for storage type (local/S3) and max file size in `backend/config.example.toml`.
    -   [ ] Sub-task: Analyze existing backend structure for API routing and file handling.
-   [x] **Task: Write Failing Tests for Local Storage Upload** (307646d)
    -   [ ] Sub-task: Write unit tests for the backend API endpoint (`POST /api/v1/upload`) when configured for local storage.
    -   [ ] Sub-task: Test successful file upload, returning the local path.
    -   [ ] Sub-task: Test file size validation (exceeding max size).
    -   [ ] Sub-task: Test security validations (disallowed file types).
-   [ ] **Task: Implement Local Storage Upload (Green Phase)**
    -   [ ] Sub-task: Create the API endpoint in `backend/api/upload` (or similar) to handle file uploads.
    -   [ ] Sub-task: Implement logic to save files to a configured local directory.
    -   [ ] Sub-task: Implement file size and type validation based on configuration.
    -   [ ] Sub-task: Return the local file path upon successful upload.
-   [ ] **Task: Write Failing Tests for S3 Storage Upload**
    -   [ ] Sub-task: Write unit tests for the backend API endpoint when configured for S3 storage.
    -   [ ] Sub-task: Test successful file upload to S3, returning the S3 object key.
    -   [ ] Sub-task: Test S3 connection failures and error handling.
-   [ ] **Task: Implement S3 Storage Upload (Green Phase)**
    -   [ ] Sub-task: Integrate S3-compatible client library (e.g., MinIO client for GoFrame).
    -   [ ] Sub-task: Implement logic to upload files to the configured S3 bucket.
    -   [ ] Sub-task: Return the S3 object key upon successful upload.
-   [ ] **Task: Refactor & Review Backend**
    -   [ ] Sub-task: Refactor the upload logic for clarity and maintainability.
    -   [ ] Sub-task: Ensure all new backend code adheres to Go formatting and best practices.
-   [ ] **Task: Conductor - User Manual Verification 'Phase 1: Backend API Development' (Protocol in workflow.md)**

## Phase 2: Frontend Component Development

-   [ ] **Task: Setup & Code Analysis**
    -   [ ] Sub-task: Analyze the existing `frontend/packages/` structure to identify the best location for a reusable upload component.
    -   [ ] Sub-task: Analyze existing UI components in Vben Admin for design patterns and styling.
-   [ ] **Task: Write Failing Tests for Upload Component**
    -   [ ] Sub-task: Write unit tests for the frontend component's file selection functionality.
    -   [ ] Sub-task: Write unit tests for drag-and-drop functionality.
    -   [ ] Sub-task: Write unit tests for displaying upload progress.
    -   [ ] Sub-task: Write unit tests for image preview functionality.
    -   [ ] Sub-task: Write unit tests for multiple file selection.
    -   [ ] Sub-task: Write unit tests for configurable file type restrictions.
-   [ ] **Task: Implement Base Upload Component (Green Phase)**
    -   [ ] Sub-task: Create the basic Vue component for file uploading.
    -   [ ] Sub-task: Implement file selection via button click.
    -   [ ] Sub-task: Implement drag-and-drop support.
    -   [ ] Sub-task: Integrate with the backend upload API.
    -   [ ] Sub-task: Implement upload progress bar display.
    -   [ ] Sub-task: Implement image preview for selected files.
    -   [ ] Sub-task: Implement multiple file upload capability.
    -   [ ] Sub-task: Implement configurable file type restrictions.
-   [ ] **Task: Refactor & Review Frontend Component**
    -   [ ] Sub-task: Refactor the frontend component for reusability and adherence to Vben Admin's patterns.
    -   [ ] Sub-task: Run `pnpm format` and `pnpm lint` to ensure code quality.
-   [ ] **Task: Conductor - User Manual Verification 'Phase 2: Frontend Component Development' (Protocol in workflow.md)**

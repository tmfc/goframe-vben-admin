# Specification: File Upload Feature

## 1. Overview

This document outlines the requirements for implementing a new file upload feature. The goal is to create a robust system that allows users to upload files, which can be stored either locally on the server or in an S3-compatible object storage service. The feature will consist of a backend API for handling the uploads and a reusable frontend component for the user interface.

## 2. Functional Requirements

### 2.1. Backend API

-   **Unified Upload Endpoint:** A single API endpoint (e.g., `POST /api/v1/upload`) will handle all file uploads.
-   **Configurable Storage:**
    -   The storage destination (either "local" or "s3") will be determined by a single, application-wide configuration setting in the backend. The API will not accept a parameter to choose the storage type per request.
    -   The maximum allowable file size for uploads will be defined in the backend configuration. The API should reject files exceeding this limit with an appropriate error response (e.g., `413 Payload Too Large`).
-   **Successful Upload Response:** Upon a successful upload, the API will return a JSON object containing the path of the uploaded file.
    -   For local storage, this will be the server-relative path to the file (e.g., `/uploads/2026/01/some-file.jpg`).
    -   For S3 storage, this will be the object key within the bucket (e.g., `uploads/2026/01/some-file.jpg`).
-   **S3 Compatibility:** The S3 integration should be compatible with major services like AWS S3, Aliyun OSS, Tencent COS, and MinIO. Configuration for S3 (endpoint, bucket, access key, secret key) will be managed in the backend configuration.

### 2.2. Frontend Upload Component

-   **File Selection:** The component will provide a button for users to click and open a file selection dialog.
-   **Drag-and-Drop:** The component will support a drag-and-drop interface, allowing users to drop files directly onto it.
-   **Multiple File Uploads:** Users will be able to select and upload multiple files simultaneously.
-   **Image Preview:** For image files, the component will display a thumbnail preview before or immediately after a successful upload.
-   **Upload Progress:** A progress bar will be displayed for each file during the upload process to provide visual feedback to the user.
-   **File Type Control:** The component will be responsible for filtering which file types are allowed. The allowed file types will be configurable when an instance of the component is used.

## 3. Non-Functional Requirements

-   **Security:** The backend must validate file types and sizes to prevent malicious uploads. Executable files should be rejected by default unless explicitly configured otherwise.
-   **Configuration:** All sensitive information (like S3 credentials) must be stored securely and managed through the application's configuration, not hardcoded.

## 4. Acceptance Criteria

-   A backend configuration setting correctly switches file storage between 'local' and 's3'.
-   The backend rejects files that exceed the configured maximum upload size.
-   The frontend component correctly shows upload progress and image previews.
-   Users can successfully upload a file via both the "click-to-select" and "drag-and-drop" methods.
-   The frontend component can be configured to restrict uploads to specific file extensions (e.g., only `.png` and `.jpg`).
-   The final uploaded file path (local or S3) is correctly returned by the API.

## 5. Out of Scope

-   File manipulation (e.g., resizing, cropping) after upload.
-   Management interface for browsing or deleting uploaded files (this could be a separate feature).

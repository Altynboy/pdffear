# pdffear

<p align="center">
    <img src="https://img.shields.io/badge/golang-v1.25-lightblue" height="25"/>
    <img src="https://img.shields.io/badge/libreoffice-v7.6-00a500" height="25"/>
    <img src="https://img.shields.io/badge/dockerfile-support-green" height="25"/>
</p>

A Go-based service to convert `.docx` and `.xlsx` files to PDF using LibreOffice within a Docker container.

## Features

-   **Word to PDF**: Convert `.docx` and `.doc` files.
-   **Excel to PDF**: Convert `.xlsx` and `.xls` files.
-   **Dockerized**: Runs in a container with all dependencies (LibreOffice, Java, Fonts).
-   **Optimized**: Multi-stage Docker build for a smaller, secure image (~1.2GB).
-   **Clean Architecture**: Refactored code following SOLID principles.

## Getting Started

### Prerequisites

-   Docker

### Running with Docker

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/Altynbekburkitbay/pdffear.git
    cd pdffear
    ```

2.  **Build the image**:
    This uses a multi-stage build to compile the binary and creates a lightweight runtime image.
    ```bash
    docker build -t pdffear .
    ```

3.  **Run the container**:
    ```bash
    docker run --rm -d -p 8080:8080/tcp pdffear
    ```

### Usage

**Upload a file for conversion**:

You can use `curl` or Postman.

```bash
# Convert a Word document
curl -X POST -F "myFile=@/path/to/document.docx" http://localhost:8080/upload --output document.pdf

# Convert an Excel spreadsheet
curl -X POST -F "myFile=@/path/to/spreadsheet.xlsx" http://localhost:8080/upload --output spreadsheet.pdf
```

**Health Check**:

```bash
curl http://localhost:8080/health
```

## Architecture

The project is structured into:

-   `converter`: Handles file type detection and LibreOffice conversion logic.
-   `storage`: Handles temporary file storage.
-   `main.go`: HTTP handlers and dependency injection.

## Optimization

The `Dockerfile` uses a **multi-stage build** process:
1.  **Builder**: Uses `golang:alpine` to compile the application and cache dependencies (`go.mod`/`go.sum`).
2.  **Runtime**: Uses a fresh `alpine` image with only the necessary runtime dependencies (`libreoffice`, `openjdk`, `fonts`).

**Impact**:
-   Reduced image size from ~1.6GB to ~1.2GB.
-   Improved security by removing Go toolchain and source code from the final image.

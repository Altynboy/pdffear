# Stage 1: Builder
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod ./
# COPY go.sum ./ # Uncomment if you have a go.sum
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -ldflags="-s -w" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pdffear .

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
# Combine apk commands to reduce layers
RUN apk add --no-cache \
    libreoffice \
    openjdk11-jre \
    msttcorefonts-installer \
    fontconfig && \
    update-ms-fonts && \
    fc-cache -f

# Create necessary directories
RUN mkdir -p /tmp/generated_pdfs /tmp/uploaded_docx /tmp/libreoffice_profiles

# Copy the binary from the builder stage
COPY --from=builder /app/pdffear /pdffear

# Set environment variables (defaults)
ENV TMP_PDF_PATH="/tmp/generated_pdfs/"
ENV TMP_DOCX_PATH="/tmp/uploaded_docx/"
ENV LIBREOFFICE_PROFILES="/tmp/libreoffice_profiles/"

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/pdffear"]

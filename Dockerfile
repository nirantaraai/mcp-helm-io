# Simple runtime Dockerfile for MCP Helm Server
# Build the binary first using: mage build:linux
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the pre-built binary
COPY bin/linux_amd64/mcp-helm-server-linux-amd64 ./mcp-helm-server

# Run the server
CMD ["./mcp-helm-server"]
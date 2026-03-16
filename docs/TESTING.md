# Testing MCP Helm Server Locally

This guide explains how to test the MCP Helm server completely locally without any cloud services.

## Prerequisites

### 1. Install Local Kubernetes Cluster

Choose one:

**Option A: Minikube**
```bash
# Install minikube
brew install minikube  # macOS
# or download from https://minikube.sigs.k8s.io/docs/start/

# Start cluster
minikube start
```

**Option B: Kind (Kubernetes in Docker)**
```bash
# Install kind
brew install kind  # macOS
# or: go install sigs.k8s.io/kind@latest

# Create cluster
kind create cluster --name helm-test
```

**Option C: Docker Desktop**
- Enable Kubernetes in Docker Desktop settings
- Wait for Kubernetes to start

### 2. Verify Kubernetes is Running
```bash
kubectl cluster-info
kubectl get nodes
```

### 3. Build the MCP Server
```bash
cd /path/to/mcp-helm-io
go build -o mcp-helm-server ./cmd/server/
```

## Local Testing Methods

### Method 1: Direct CLI Testing (Simplest)

Test the server directly from command line:

```bash
# Start the server
./mcp-helm-server

# The server will log to stderr, you'll see:
# INFO starting MCP server
# INFO registering MCP tools
```

### Method 2: Test with Local AI (Ollama + Continue.dev)

**Step 1: Install Ollama (Local AI)**
```bash
# macOS/Linux
curl https://ollama.ai/install.sh | sh

# Start Ollama
ollama serve

# Pull a model
ollama pull codellama
```

**Step 2: Install Continue.dev in VS Code**
```bash
# In VS Code, install Continue extension
# Configure it to use Ollama
```

**Step 3: Configure MCP in Continue**
Edit `.continue/config.json`:
```json
{
  "mcpServers": {
    "helm": {
      "command": "/absolute/path/to/mcp-helm-server"
    }
  }
}
```

### Method 3: Manual JSON-RPC Testing

Create test scripts to send JSON-RPC requests:

**test-deploy.sh**:
```bash
#!/bin/bash

# Send deploy request via stdin
cat <<EOF | ./mcp-helm-server
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "deploy_chart",
    "arguments": {
      "chart": "bitnami/nginx",
      "release_name": "test-nginx",
      "namespace": "default",
      "wait": true,
      "timeout": 300
    }
  }
}
EOF
```

**test-list.sh**:
```bash
#!/bin/bash

cat <<EOF | ./mcp-helm-server
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "list_releases",
    "arguments": {
      "namespace": "default"
    }
  }
}
EOF
```

Make executable and run:
```bash
chmod +x test-*.sh
./test-deploy.sh
./test-list.sh
```

### Method 4: Interactive Testing with MCP Inspector (Local Web UI)

```bash
# Install MCP Inspector globally
npm install -g @modelcontextprotocol/inspector

# Run inspector (opens local web browser)
mcp-inspector ./mcp-helm-server

# This opens http://localhost:5173
# You can:
# - See all available tools
# - Test each tool interactively
# - View request/response in real-time
```

## Complete Local Test Workflow

### 1. Setup Local Environment

```bash
# Start local Kubernetes
minikube start

# Add Helm repository
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Verify
helm search repo nginx
```

### 2. Build and Run Server

```bash
# Build
go build -o mcp-helm-server ./cmd/server/

# Run with debug logging
LOG_LEVEL=debug ./mcp-helm-server 2>&1 | tee server.log
```

### 3. Test Each Operation

**Deploy a Chart**:
```bash
# Using kubectl to verify
kubectl create namespace test-ns

# Deploy via MCP (use inspector or script)
# Then verify:
kubectl get pods -n test-ns
helm list -n test-ns
```

**List Releases**:
```bash
# Should show the deployed release
helm list -A
```

**Get Status**:
```bash
# Check specific release
helm status test-nginx -n test-ns
```

**Upgrade**:
```bash
# Upgrade to newer version
# Then verify:
helm history test-nginx -n test-ns
```

**Rollback**:
```bash
# Rollback to previous
# Verify:
helm history test-nginx -n test-ns
```

**Uninstall**:
```bash
# Remove release
# Verify:
kubectl get pods -n test-ns
helm list -n test-ns
```

## Automated Local Testing

Create a test suite:

**run-tests.sh**:
```bash
#!/bin/bash
set -e

echo "🚀 Starting MCP Helm Server Tests"

# Start server in background
./mcp-helm-server &
SERVER_PID=$!
sleep 2

# Test 1: Deploy
echo "📦 Test 1: Deploy Chart"
./test-deploy.sh
sleep 5
kubectl get pods -n default | grep test-nginx

# Test 2: List
echo "📋 Test 2: List Releases"
./test-list.sh

# Test 3: Status
echo "ℹ️  Test 3: Get Status"
./test-status.sh

# Test 4: Uninstall
echo "🗑️  Test 4: Uninstall"
./test-uninstall.sh
sleep 5

# Cleanup
kill $SERVER_PID
echo "✅ All tests passed!"
```

## Debugging Locally

### View Server Logs
```bash
# Run with verbose logging
LOG_LEVEL=debug ./mcp-helm-server 2>&1 | tee debug.log

# In another terminal, watch logs
tail -f debug.log
```

### Check Kubernetes State
```bash
# Watch pods
kubectl get pods -A -w

# Check Helm releases
helm list -A

# View events
kubectl get events -A --sort-by='.lastTimestamp'
```

### Test Helm Directly
```bash
# Verify Helm works independently
helm install test-direct bitnami/nginx -n default
helm list -n default
helm uninstall test-direct -n default
```

## Performance Testing Locally

```bash
# Monitor server resources
ps aux | grep mcp-helm-server

# Test concurrent requests
for i in {1..10}; do
  ./test-deploy.sh &
done
wait

# Check cluster resources
kubectl top nodes
kubectl top pods -A
```

## Cleanup

```bash
# Stop server
pkill mcp-helm-server

# Clean Kubernetes
kubectl delete namespace test-ns
helm uninstall test-nginx -n default

# Stop cluster
minikube stop  # or: kind delete cluster --name helm-test
```

## Troubleshooting

### Server won't start
```bash
# Check if port is in use
lsof -i :8080

# Check kubeconfig
kubectl config view
kubectl config current-context
```

### Can't connect to Kubernetes
```bash
# Verify cluster is running
kubectl cluster-info

# Check kubeconfig
export KUBECONFIG=~/.kube/config
kubectl get nodes
```

### Helm operations fail
```bash
# Check Helm version
helm version

# Verify repositories
helm repo list
helm repo update

# Test Helm directly
helm install test bitnami/nginx --dry-run
```

#
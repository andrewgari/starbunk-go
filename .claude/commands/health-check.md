---
name: health-check
description: Comprehensive health check that runs build, test, lint, builds docker containers, checks the prometheus /health endpoint, and reports to the user.
---

# Comprehensive Health Check

Run a full end-to-end local health check on the project, including compiling the code, running tests, linting, building Docker images, and verifying service health via the Prometheus endpoint.

## Instructions

Run the following checks sequentially:

1. **Build the Code**:
   ```bash
   go build ./...
   ```

2. **Test the Code**:
   ```bash
   go test ./...
   ```

3. **Lint**:
   ```bash
   golangci-lint run
   ```

4. **Build Docker Containers**:
   Build the images using the local docker-compose override:
   ```bash
   docker compose -f docker/docker-compose.yml build
   ```

5. **Start the Containers**:
   Start the services in the background:
   ```bash
   docker compose -f docker/docker-compose.yml up -d
   ```

6. **Check Prometheus /health Endpoint**:
   Check the `/health` endpoint. *Note: adjust the port depending on which service is exposing the prometheus metrics (e.g. 8080, 9090).*
   ```bash
   # Wait a few seconds for services to start
   sleep 5
   
   # Example curl command for the /health endpoint
   curl -sSf http://localhost:8080/health || echo "Prometheus /health endpoint check failed"
   ```

7. **Clean up**:
   Bring the containers back down after the check is complete:
   ```bash
   docker compose -f docker/docker-compose.yml down
   ```

8. **Report to the User**:
   Consolidate the results of all the above steps. Clearly report whether each phase (Build, Test, Lint, Docker Build, Health Check) passed or failed, and provide the relevant error logs if any step failed.

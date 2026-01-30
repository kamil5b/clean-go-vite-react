# Deployment Guide

## Overview

This application consists of two independent runtimes:
- **Server**: HTTP API service (port 8080)
- **Worker**: Background job processor (connected to Redis)

Both can be deployed separately and scaled independently.

## Prerequisites

- Docker & Docker Compose
- Go 1.25+ (for local development)
- Redis 7+ (included in docker-compose)

## Local Development

### Quick Start

1. **Start dependencies:**
```bash
docker-compose up redis
```

2. **Run server:**
```bash
make server
```

3. **Run worker (in another terminal):**
```bash
make worker
```

4. **Run both with frontend:**
```bash
make dev
```

## Docker Deployment

### Build Docker Images

```bash
# Build both server and worker
docker build -t myapp:server -f Dockerfile --target build-server .
docker build -t myapp:worker -f Dockerfile.worker .

# Or use docker-compose to build both
docker-compose -f docker-compose.prod.yml build
```

### Run with Docker Compose

**Production-like setup with all services:**

```bash
docker-compose -f docker-compose.prod.yml up
```

This starts:
- Redis (port 6379, internal)
- Server (port 8080)
- Worker (background)

### Run Individual Containers

**Server only:**
```bash
docker run -p 8080:8080 \
  -e REDIS_HOST=redis.example.com \
  -e REDIS_PORT=6379 \
  myapp:server
```

**Worker only:**
```bash
docker run \
  -e REDIS_HOST=redis.example.com \
  -e REDIS_PORT=6379 \
  myapp:worker
```

## Environment Variables

### Server Configuration
- `SERVER_HOST` - Bind address (default: empty = all interfaces)
- `SERVER_PORT` - Port to listen on (default: 8080)
- `SERVER_READ_TIMEOUT` - Read timeout (default: 15s)
- `SERVER_WRITE_TIMEOUT` - Write timeout (default: 15s)
- `SERVER_IDLE_TIMEOUT` - Idle timeout (default: 60s)

### Redis Configuration
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_DB` - Redis database number (default: 0)
- `REDIS_PASSWORD` - Redis password (default: empty)

### Asynq Configuration
- `ASYNQ_ENABLED` - Enable async processing (default: true)
- `ASYNQ_REDIS_ADDR` - Redis address for Asynq (default: localhost:6379)
- `ASYNQ_CONCURRENCY` - Worker concurrency (default: 10)
- `ASYNQ_MAX_RETRIES` - Max retry attempts (default: 3)
- `ASYNQ_RETRY_DELAY_MIN` - Min retry delay (default: 5s)
- `ASYNQ_RETRY_DELAY_MAX` - Max retry delay (default: 5m)

### Database Configuration (Optional)
- `DATABASE_DSN` - Database connection string
- `DATABASE_MAX_OPEN_CONNS` - Connection pool size (default: 25)
- `DATABASE_MAX_IDLE_CONNS` - Idle connections (default: 5)
- `DATABASE_CONN_MAX_LIFETIME` - Connection max lifetime (default: 5m)

## Kubernetes Deployment

### Server Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-server
spec:
  replicas: 2
  selector:
    matchLabels:
      app: app-server
  template:
    metadata:
      labels:
        app: app-server
    spec:
      containers:
      - name: server
        image: myapp:server
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_PORT
          value: "8080"
        - name: REDIS_HOST
          value: redis-service
        - name: REDIS_PORT
          value: "6379"
        livenessProbe:
          httpGet:
            path: /api/health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: app-server-service
spec:
  selector:
    app: app-server
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Worker Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-worker
spec:
  replicas: 2
  selector:
    matchLabels:
      app: app-worker
  template:
    metadata:
      labels:
        app: app-worker
    spec:
      containers:
      - name: worker
        image: myapp:worker
        env:
        - name: REDIS_HOST
          value: redis-service
        - name: REDIS_PORT
          value: "6379"
        - name: ASYNQ_CONCURRENCY
          value: "10"
```

### Redis Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-service
spec:
  selector:
    app: redis
  ports:
  - protocol: TCP
    port: 6379
    targetPort: 6379
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
spec:
  serviceName: redis-service
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        ports:
        - containerPort: 6379
        volumeMounts:
        - name: data
          mountPath: /data
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 10Gi
```

## Health Checks

### Server Health Endpoint

```bash
curl http://localhost:8080/api/health
```

Response:
```json
{
  "status": "ok",
  "message": "Service is healthy"
}
```

### Liveness & Readiness Probes

Both probes should target: `GET /api/health`

- **Liveness**: Detects if container needs restart
- **Readiness**: Determines if container can receive traffic

## Monitoring

### Logs

**Docker:**
```bash
docker logs <container-id>
```

**Kubernetes:**
```bash
kubectl logs -f deployment/app-server
kubectl logs -f deployment/app-worker
```

### Metrics to Monitor

- HTTP request latency
- Worker task processing time
- Redis connection pool usage
- Failed task retry count
- Error rates

### Example: Prometheus Metrics

Add to metrics endpoint (future enhancement):
```go
api.GET("/api/metrics", metricsHandler.GetMetrics)
```

## Scaling

### Server (Stateless)

Scale horizontally as needed:
```bash
kubectl scale deployment app-server --replicas=5
```

### Worker (Stateless)

Increase concurrency or replicas:
```bash
# Increase per-worker concurrency
ASYNQ_CONCURRENCY=20

# Or scale replicas
kubectl scale deployment app-worker --replicas=3
```

## Database Migration

If using a database, run migrations before starting:

```bash
# Before starting server
docker run --rm \
  -e DATABASE_DSN=... \
  myapp:server \
  /app/migrate up
```

## Production Checklist

- [ ] Redis is persisted (with backup strategy)
- [ ] Database backups configured
- [ ] Health checks returning 200 OK
- [ ] Environment variables set correctly
- [ ] TLS/SSL configured (if needed)
- [ ] CORS properly configured
- [ ] Rate limiting configured
- [ ] Logging/monitoring enabled
- [ ] Error tracking (Sentry/etc) enabled
- [ ] Load balancer health checks configured
- [ ] Database connection pooling validated
- [ ] Redis high availability configured

## Rollback Strategy

1. Keep previous image version tagged
2. Update deployment to previous image
3. Monitor metrics during rollback
4. Keep database migrations reversible

Example:
```bash
docker run -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d  # Redeploy with old image
```

## Troubleshooting

### Server won't start

Check logs:
```bash
docker logs app-server
```

Common issues:
- Port already in use: Change `SERVER_PORT`
- Redis unavailable: Check `REDIS_HOST` and `REDIS_PORT`
- Wrong environment: Verify all env vars are set

### Worker not processing tasks

1. Check worker is running: `docker ps | grep worker`
2. Verify Redis connection: `redis-cli ping`
3. Check logs: `docker logs app-worker`
4. Verify task was enqueued: Check Redis keys

### High memory usage

- Reduce `ASYNQ_CONCURRENCY`
- Check for memory leaks in services
- Monitor connection pool sizes
- Consider increasing worker replicas

## References

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Redis Documentation](https://redis.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
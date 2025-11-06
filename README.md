# Distributed Bookstore

A practical microservices implementation of the patterns documented in
[`documentation/SERVICE_INTERACTIONS_AND_BEST_PRACTICES.md`](documentation/SERVICE_INTERACTIONS_AND_BEST_PRACTICES.md).
The codebase showcases Tanenbaum & van Steen's distributed-systems
principles with a polyglot stack running on Kubernetes/Minikube.

---

## At a Glance

| Service | Language | Purpose | Data Store |
|---------|----------|---------|-----------|
| **API Gateway** | Go (Fiber) | Single entrypoint, JWT enforcement, optional rate limiting, proxy to backend services | — |
| **Catalog** | Go (Fiber + GORM) | Books, authors, categories, publishers, rich seed data | PostgreSQL |
| **User** | Go | Registration, login, refresh tokens, profiles, wishlist, JWT issuance, publishes `user.*` / `wishlist.*` events | PostgreSQL |
| **Cart** | Go | Session carts backed by Redis, publishes `cart.*` events | Redis |
| **Order** | Go | Order capture skeleton (ready for future saga/payment integration) | PostgreSQL |
| **Inventory** | Python (FastAPI) | Stock management & warehouse API | PostgreSQL |
| **Review** | Python (FastAPI) | Reviews, helpful votes, TextBlob sentiment analysis | PostgreSQL |
| **Recommendation** | Python (FastAPI) | Similar / trending / popular book recommendations with Redis caching hooks | (Redis planned) |
| **Notification** | Go | RabbitMQ consumer that logs/dispatches notifications (email/SMS/push adapters ready to plug in) | — |
| **Admin** | Go | WIP analytics façade (currently connects to catalog DB for reports) | PostgreSQL |
| **Payment** | TypeScript (Express) | Stripe integration placeholder (webhook + REST stubs) | — |

Frontend: **React 18 + TypeScript** (`frontend/customer-app`) built with
TanStack Router/Query, Zustand, shadcn/ui and Sonner toasts.

Messaging: **RabbitMQ topic exchange** with durable queues. Publishers:
User Service (`user.registered`, `user.logged_in`, `wishlist.*`), Cart
Service (`cart.item_added`, `cart.item_updated`, `cart.item_removed`, `cart.cleared`).
Notification Service consumes and logs each message (hook for real
email/SMS providers).

Observability: Every service exposes `/health` and `/ready`; logs use
structured output suitable for `kubectl logs` and local triage.

---

## Running the System Locally (Minikube)

We keep two Kubernetes manifests:

- `infrastructure/k8s/` – main manifests (kept aligned with upstream changes)
- `infrastructure/k8s-dev/` – developer copy you can freely tweak

### Prerequisites

- Docker + Minikube (`minikube start`)
- kubectl
- `pnpm` (for the frontend) if you plan to run builds locally
- Go 1.21+ / Python 3.11+ optional for direct service work

### 1. Build & Deploy to Minikube

```bash
cd infrastructure/k8s-dev

# Point Docker at Minikube
eval "$(minikube docker-env)"

# Build all service images and apply manifests
./deploy.sh
```

The deploy script will:

1. Build Docker images for each service with the `:latest` tag inside the
   Minikube daemon
2. Apply secrets/configmaps, Postgres + Redis + RabbitMQ manifests
3. Roll out every microservice plus the React frontend

### 2. Port-Forward for Local Access

```bash
# In separate terminals
kubectl port-forward -n bookstore-dev svc/frontend 3000:80
kubectl port-forward -n bookstore-dev svc/api-gateway 8080:8080
kubectl port-forward -n bookstore-dev svc/rabbitmq 15672:15672  # RabbitMQ UI (optional)
```

Frontend will be available at <http://localhost:3000> and will talk to
services via the API Gateway at <http://localhost:8080>.

### 3. Tear Down

```bash
cd infrastructure/k8s-dev
./undeploy.sh
```

---

## Working on Individual Services

- **Go services**: standard layout (`cmd/`, `internal/`, `pkg/`). Run `go build ./...` inside the service folder.
- **Python services**: FastAPI apps under `app/` with Alembic migrations. Use uvicorn for local runs (`pip install -r requirements.txt`).
- **Notification service**: Go Fiber app with RabbitMQ consumer – tail logs via `kubectl logs -n bookstore-dev deployment/notification-service -f` to observe events emitted by Cart/User services.
- **Frontend**: `pnpm install && pnpm dev` (the production build runs automatically in the cluster via Nginx).

Useful local commands:

```bash
# Frontend build check
cd frontend/customer-app
pnpm build

# Run Go tests/build for a service
cd services/user-service
go test ./...

go build ./...

# FastAPI review service local run
cd services/review-service
uvicorn app.main:app --reload --port 8088
```

---

## Event Catalogue (Current)

| Topic | Producer | Consumer | Description |
|-------|----------|----------|-------------|
| `user.registered` | User Service | Notification Service | Welcome workflow |
| `user.logged_in` | User Service | Notification Service | Security/login alerts |
| `wishlist.added` & `wishlist.removed` | User Service | Notification Service | Wishlist instrumentation |
| `wishlist.cleared` | User Service | Notification Service | Wishlist reset |
| `cart.item_added`, `cart.item_updated`, `cart.item_removed`, `cart.cleared` | Cart Service | Notification Service | Cart engagement tracking |

(Extend `services/notification-service/internal/notification/message.go`
with new templates as you add events.)

---

## Documentation & Further Reading

- [`documentation/SERVICE_INTERACTIONS_AND_BEST_PRACTICES.md`](documentation/SERVICE_INTERACTIONS_AND_BEST_PRACTICES.md)
  – System interactions, architecture principles, deployment order
- `CLAUDE.md` – Historical planning notes & backlog items
- Per-service `README.md` files for service-specific behaviour

---

## Roadmap / Known Gaps

- Payment service & admin analytics are scaffolds awaiting full Stripe / reporting integration
- Observability stack (Prometheus/Grafana/Jaeger) intentionally omitted from the dev manifests to keep local setup lightweight
- Order workflow & inventory reservations will evolve into a saga once payment integration lands

Contributions and experiments are welcome—use the `k8s-dev` manifests to
iterate without disturbing the main configuration.

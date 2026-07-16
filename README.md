# Tyr Fintech

Tyr Fintech is a modern, high-performance, and secure multi-currency digital wallet application. It features a robust **Go (Gin-gonic)** backend that guarantees transaction consistency, and a beautiful, custom **glassmorphic React SPA** frontend.

---

## 🚀 Tech Stack

### Backend
* **Language**: Go 1.26+
* **Web Framework**: Gin-Gonic (high performance, routing, middlewares)
* **Database Driver**: PGX v5 (connection pooling, native Postgres integration)
* **Database**: PostgreSQL 16 (relational database with transaction isolation)
* **Migrations**: Golang-migrate (versioned database migrations)

### Frontend
* **Build System**: Vite (lightning-fast HMR and building)
* **Framework**: React (functional components, contexts, hooks)
* **Styling**: Tailwind CSS v4 (responsive utility-first layout, custom glassmorphism design system)
* **API Client**: Axios (configured with credentials and global interceptors)

### DevOps & CI/CD
* **Containers**: Docker & Docker Compose (multi-container orchestrated setup)
* **CI**: GitHub Actions (automated testing pipeline for backend)

---

## 🔒 Key Design & Security Features

1. **Transaction Integrity (ACID)**:
   - Implements **pessimistic row-level locking (`FOR UPDATE`)** in Go transactions when updating wallet balances.
   - Prevents **"Lost Update"** problems and guarantees money transfer reliability during high concurrency.
2. **Idempotency Protection**:
   - The `/transfer` endpoint accepts an optional `X-Idempotency-Key` header.
   - Prevents duplicate requests (e.g., due to client double-clicking or network retries) from triggering multiple transfers.
3. **Multi-Currency Safety**:
   - Stored in `BIGINT` (representing money in minor units/cents) to avoid floating-point rounding errors.
   - Constraint checks enforce `balance >= 0` to prevent overdrafts.
4. **JWT Auth via HttpOnly Cookies**:
   - User authentication state is managed securely with JSON Web Tokens (JWT) stored in HTTP-Only, Secure cookies to prevent XSS-based token theft.
5. **Robust Error Handling**:
   - Custom `AppError` wrapper maps app-specific database and transaction errors to corresponding HTTP status codes cleanly.
6. **Data Export**:
   - Download transaction statements on-demand as either **CSV** or **PDF** files directly from the wallet dashboard.

---

## 📂 Project Structure

```
├── .github/workflows/       # GitHub Actions CI pipelines
│   └── ci.yml               # Runs automated tests on every push/PR
├── backend/                 # Backend source code
│   ├── cmd/api/             # App entrypoint (main.go)
│   ├── internal/            # Core business logic
│   │   ├── db/              # Postgres connections and sequencers
│   │   ├── dto/             # Request/Response Data Transfer Objects
│   │   ├── handlers/        # Gin controllers and routers
│   │   ├── middleware/      # Auth and CORS middlewares
│   │   ├── models/          # Relational struct models
│   │   ├── repos/           # Database access layer (SQL queries)
│   │   └── worker/          # Asynchronous webhook queues/workers
│   ├── migrations/          # Up/Down SQL schema migrations
│   ├── pkg/                 # Common helpers (apperrors, JWT, response, export)
│   ├── Dockerfile           # Multistage backend container build
│   └── docker-compose.yml   # Docker compose configuration (DB only/Local Dev)
├── frontend/                # Frontend source code
│   ├── public/              # Static assets
│   ├── src/                 # React frontend source files
│   │   ├── components/      # UI components (forms, wallet grids)
│   │   ├── context/         # Auth state providers
│   │   ├── lib/             # Axios API config
│   │   └── pages/           # Pages (Dashboard, Login, Register)
│   ├── Dockerfile           # Multistage frontend build served via Nginx
│   └── nginx.conf           # Custom Nginx configuration
└── docker-compose.yml       # Orchestrated system (Frontend + Backend + DB)
```

---

## ⚙️ Development Setup

### Running with Docker (Recommended)
Build and spin up the entire application stack (PostgreSQL + Backend + Frontend) in one command:

1. Clone the repository and navigate to the project root.
2. Spin up the containers:
   ```bash
   docker compose up --build -d
   ```
3. Run the database migrations to set up the tables:
   ```bash
   cd backend
   make migrate-up
   ```
4. Access the application:
   * **Frontend**: [http://localhost:3000](http://localhost:3000)
   * **Backend API**: [http://localhost:8080](http://localhost:8080)

---

### Running Locally (Manual Setup)

#### 1. Database Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Start the database service:
   ```bash
   docker compose up -d db
   ```
3. Run database migrations:
   ```bash
   make migrate-up
   ```

#### 2. Run the Backend API
1. Create a `backend/.env` file from the example:
   ```bash
   cp .env.example .env
   ```
2. Start the Go server:
   ```bash
   go run cmd/api/main.go
   ```
   *Backend is now serving requests on `http://localhost:8080`.*

#### 3. Run the Frontend
1. Open a new terminal and navigate to the frontend:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the Vite development server:
   ```bash
   npm run dev
   ```
   *Frontend is now running on [http://localhost:3000](http://localhost:3000).*

---

## 🛠️ Running Tests
To run backend unit and service tests:
```bash
cd backend
go test -v ./...
```

---

## 📡 API Reference

### Auth
* **`POST /api/v1/auth/register`**: Registers a new user.
* **`POST /api/v1/auth/login`**: Authenticates user and sets HttpOnly JWT cookie.
* **`POST /api/v1/logout`** (Protected): Clears user session.

### Wallets
* **`GET /api/v1/wallets`** (Protected): Retrieves all wallets owned by the authenticated user.
* **`POST /api/v1/wallets`** (Protected): Activates/Creates a new wallet for a specified currency (`TRY`, `USD`, or `EUR`).
* **`DELETE /api/v1/wallets/:walletID`** (Protected): Deletes the specified wallet.

### Transfers & History
* **`POST /api/v1/transfer`** (Protected): Initiates money transfer between wallets with idempotency checking.
* **`GET /api/v1/transactions/:walletID`** (Protected): Retrieves transaction logs for the wallet.
* **`GET /api/v1/transactions/:walletID/export`** (Protected): Exports the transaction logs. Query param `format` accepts `csv` or `pdf`.

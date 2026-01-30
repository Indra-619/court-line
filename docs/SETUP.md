# ⚙️ Setup Guide

Welcome to the setup guide for **Book Lapangan**. Follow these steps to get your development environment ready in minutes.

---

## 📋 Prerequisites

Before we begin, ensure you have the necessary tools installed.

### 🐳 Docker (Recommended)
The easiest way to run the project.
- **Docker** & **Docker Compose**

### 🛠️ Manual Setup
If you prefer running services individually on your machine.
- **Go** (v1.21+)
- **Node.js** (v18+) & **npm**
- **MongoDB** (Local instance or Atlas URI)

---

## 1. 📥 Clone the Repository

Get the code on your local machine:

```bash
git clone https://github.com/Indra-619/court-line.git
cd book-lapangan
```

## 2. ✅ Verify Installation

Run our safety script to ensure your environment is ready (Go, Node, Docker modules):

```powershell
# Run verification script
.\scripts\verify_build.ps1
```

---

## 3. 🔐 Environment Configuration

We use environment variables to keep secrets safe.

1. **Copy the example file**:
   ```bash
   cp .env.example .env
   ```

2. **Configure keys** in `.env`:

   | Variable | Description |
   | :--- | :--- |
   | `MONGODB_URI` | Connection string (e.g., `mongodb://localhost:27017/booklapangan`) |
   | `GOOGLE_CLIENT_ID` | OAuth Client ID from Google Cloud Console |
   | `GOOGLE_CLIENT_SECRET` | OAuth Client Secret |
   | `JWT_SECRET` | Secret key for signing session tokens |

---

## 4. 🚀 Running the App

Choose your preferred method below.

### Option A: Docker (✨ Recommended)

This builds and starts the **Backend**, **Frontend**, and **MongoDB** automatically.

```bash
docker-compose up --build
```

| Service | URL |
| :--- | :--- |
| **Frontend** | [http://localhost:3000](http://localhost:3000) |
| **Backend API** | [http://localhost:8080](http://localhost:8080) |
| **MongoDB** | `localhost:27017` |

### Option B: Manual Execution

#### 🖥️ Backend (Go)
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### 🎨 Frontend (Nuxt 3)
```bash
cd frontend
npm install
npm run dev
```

---

## 🆘 Troubleshooting

- **🚧 Port Conflicts**: Check if ports `3000`, `8080`, or `27017` are already in use.
- **🔌 Database Connection**: Ensure `MONGODB_URI` is correct if running manually.

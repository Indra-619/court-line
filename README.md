# 🏟️ Book Lapangan Online (CourtLine)

![Project Status](https://img.shields.io/badge/status-active-success.svg)
![Build Status](https://img.shields.io/badge/build-safe-green.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

A modern fullstack application for booking sports venues online, featuring Google OAuth authentication, real-time availability, and a beautiful dark-themed UI.

## ✨ Features

- 🔐 **Google OAuth Login** - Secure authentication with your Google account
- 📅 **Easy Booking** - Book courts with date and time selection
- 💰 **Price Calculator** - Automatic price estimation based on duration
- 📱 **Responsive Design** - Works beautifully on desktop and mobile
- 🌙 **Modern Dark Theme** - Premium UI with glassmorphism effects

## 🛠️ Tech Stack & Resources

### Frontend
![Nuxt.js](https://img.shields.io/badge/Nuxt.js-002E3B?style=for-the-badge&logo=nuxt.js&logoColor=00DC82)
![Vue.js](https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vue.js&logoColor=4FC08D)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white)

### Backend
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

### Database
![MongoDB](https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white)

---

## 🏗️ Architecture

```mermaid
graph TD
    User["User / Client"] -->|HTTP/Vue| Client["Frontend (Nuxt 3)"]
    subgraph Frontend
        Client --> Services["Service Layer"]
        Services --> Composables["Auth Composables"]
    end
    
    Services -->|"API Calls (Result Pattern)"| Server["Backend (Go / Gin)"]
    
    subgraph Backend (Clean Architecture)
        Server --> Presentation["Presentation (API Handlers)"]
        Presentation --> UseCase["Use Case (Business Logic)"]
        UseCase --> Entity["Domain (Entities)"]
        UseCase --> RepoInterface["Domain (Repository Interfaces)"]
        RepoInterface -.->|Implementation| Infrastructure["Infrastructure (MongoDB)"]
    end
    
    Infrastructure -->|Read/Write| DB[("MongoDB")]
```

## 📂 Project Structure

```bash
book-lapangan/
├── backend/
│   ├── cmd/server/         # Entry point (Dependency Injection wiring)
│   ├── internal/
│   │   ├── domain/         # Domain Layer (Entities & Repository Interfaces)
│   │   ├── usecase/        # Use Case Layer (Business Logic)
│   │   ├── infrastructure/ # Infrastructure Layer (Repo Implementations)
│   │   ├── presentation/   # Presentation Layer (DI-based Handlers)
│   │   ├── middleware/     # Auth middleware
│   │   └── routes/         # Route definitions (DI-based)
│   ├── pkg/result/         # Generic Result pattern for error handling
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── services/           # Service Layer (API interactions)
│   ├── components/         # Vue components
│   ├── composables/        # Auth composables
│   ├── pages/              # Nuxt pages (TypeScript & Services)
│   ├── nuxt.config.ts
│   └── Dockerfile
├── docs/                   # Documentation
├── docker-compose.yml
├── .env.example
├── CONTRIBUTING.md
└── README.md
```


## 📚 Documentation

- [**Setup Guide**](docs/SETUP.md) - Detailed instructions for Docker and Manual setup.
- [**Testing Guide**](docs/TESTING.md) - How to run manual tests and verify health.
- [**Contributing Guide**](CONTRIBUTING.md) - Workflow, code style, and PR process.

## 🚀 Getting Started

### Prerequisites
- **Docker Desktop** (Required for containerization)
- **Git**
- **Google Cloud Console** account (for OAuth)

### Setup Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. Create a new project or select existing
3. Create OAuth 2.0 credentials (Web application)
4. Set Authorized redirect URI: `http://localhost:8080/auth/google/callback`
5. Copy Client ID and Client Secret

### Installation & Run

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Indra-619/court-line.git
   cd court-line
   ```

2. **Create environment file**:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and add your Google OAuth credentials:
   ```env
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   JWT_SECRET=your-super-secret-jwt-key
   ```

3. **Run with Docker Compose**:
   ```bash
   docker-compose up --build
   ```

4. **Access the Application**:
   - **Frontend**: [http://localhost:3000](http://localhost:3000)
   - **Backend API**: [http://localhost:8080](http://localhost:8080)

## 📡 API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/auth/google` | Redirect to Google OAuth | ❌ |
| GET | `/auth/google/callback` | Handle OAuth callback | ❌ |
| GET | `/auth/me` | Get current user info | ✅ |
| GET | `/api/courts` | Get all courts | ❌ |
| GET | `/api/courts/:id` | Get court by ID | ❌ |
| POST | `/api/courts` | Create new court | ✅ |
| PUT | `/api/courts/:id` | Update court | ✅ |
| DELETE | `/api/courts/:id` | Delete court | ✅ |
| POST | `/api/bookings` | Create booking | ✅ |
| GET | `/api/bookings` | Get user's bookings | ✅ |

## 🖼️ Screenshots

### Homepage
- Hero section with stats
- Grid of available courts
- Dark theme with gradient accents

### Court Detail
- Full court information
- Booking form with date/time selection
- Price calculator

### Login
- Google OAuth integration
- Quick, secure authentication

---

## 📝 License

MIT License - feel free to use this project for learning or personal projects.

*Created by [Indra-619](https://github.com/Indra-619)*

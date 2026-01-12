# Book Lapangan Online

A fullstack application for booking sports venues, built with:
- **Frontend**: Nuxt.js 3 (Vue 3)
- **Backend**: Go (Golang)
- **Database**: MongoDB
- **Infrastructure**: Docker & Docker Compose

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- Git

## Getting Started

1. **Clone the repository** (if you haven't already):
   ```bash
   git clone <your-repo-url>
   cd book-lapangan
   ```

2. **Run with Docker Compose**:
   ```bash
   docker-compose up --build
   ```
   This will start:
   - **Backend**: http://localhost:8080
   - **Frontend**: http://localhost:3000
   - **MongoDB**: localhost:27017

3. **Verify Installation**:
   Open http://localhost:3000 in your browser. You should see the "Book Lapangan Online" homepage.
   Click "Check Connection" to verify the frontend can talk to the backend.

## Project Structure

- `/backend`: Go application code.
- `/frontend`: Nuxt.js application code.
- `docker-compose.yml`: Orchestration for all services.

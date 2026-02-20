# 🤝 Contributing to Book Lapangan

We love your input! We want to make contributing to **Book Lapangan** as easy and transparent as possible, whether it's:

- 🐛 Reporting a bug
- 📄 Discussing the current state of the code
- 🔧 Submitting a fix
- ✨ Proposing new features

---

## 🛤️ Workflow Overview

We use a standard **Fork & Pull** workflow.

```mermaid
graph LR
    A[Fork Repo] --> B[Clone Locally]
    B --> C[Create Branch]
    C --> D[Commit Changes]
    D --> E[Push to Fork]
    E --> F[Open Pull Request]
    F --> G[Review & Merge]
```

---

## 🚀 Getting Started

1. **Fork** the repo on GitHub.
2. **Clone** the project to your own machine.
3. **Setup** the environment (see [Setup Guide](docs/SETUP.md)).
4. **Create a Branch** with a descriptive name.

   ```bash
   git checkout -b feature/amazing-feature
   # or
   git checkout -b fix/annoying-bug
   ```

---

## ✍️ Coding Guidelines

### 🏗️ Architectural Standards
- **Clean Architecture**: Follow the separation of layers:
    - `domain`: Entities and Repository interfaces. **No dependencies**.
    - `usecase`: Business logic. Depends only on `domain`.
    - `infrastructure`: Concrete implementations (DB, external APIs).
    - `presentation`: API handlers. Depends only on `usecase`.
- **Result Pattern**: Use the `result` package in `pkg/result` for all business logic and handler returns to ensure explicit error handling.
- **Service Layer (Frontend)**: Decouple pages from API calls using the `services/` layer.

### 🎨 Style Guide
- **Go**: Always run `gofmt` before committing.
- **Frontend**: Use TypeScript for all pages and components (`<script setup lang="ts">`).
- **Naming**: Use `camelCase` for JS/Go variables, `kebab-case` for filenames.

### 💬 Commit Messages
We follow the **Conventional Commits** specification.

| Prefix | Description | Example |
| :--- | :--- | :--- |
| `feat` | A new feature | `feat: add payment gateway integration` |
| `fix` | A bug fix | `fix: resolve crash on login` |
| `docs` | Documentation only | `docs: update API reference` |
| `style` | Layout/Style (no code change) | `style: fix button padding` |
| `refactor` | Code restructuring | `refactor: simplify auth middleware` |

---

## 📬 Pull Request Process

1. **Run Verification**: Execute `.\scripts\verify_build.ps1` and ensure all checks pass. 🛡️
   - Quality Sentinel: New logic must have **>85% test coverage**.
2. UPDATE the `README.md` or docs with details of changes if relevant.
3. DESCRIBE your changes clearly in the PR description.
4. LINK to any related issues (e.g., `Closes #123`).

---

## 🐞 Reporting Issues

- Use a clear and descriptive title.
- Describe the exact steps which reproduce the problem.
- Provide logs or screenshots if possible.

**Thank you for helping us build a better platform!** 🏟️

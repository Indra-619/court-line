# 🧪 Testing & Quality Assurance

Welcome to the **Court Line** testing guide. Use this documentation to ensure your contributions meet our quality standards. Our testing strategy ensures reliance, stability, and zero-build-failures.

---

## 🚀 Quick Verification

Before rewriting any code or pushing changes, run our **Master Guardrail Script**. This script validates the entire stack in one go.

```powershell
# Run from the project root
.\scripts\verify_build.ps1
```

**What this checks:**
1.  **Backend Integrity**: Runs all Go unit tests.
2.  **Frontend Type Safety**: Verification of TypeScript types in Vue/Nuxt.
3.  **Infrastructure Config**: Validates `docker-compose.yml` logic.

> 🛑 **Note:** If this script fails, the CI pipeline will also fail. Fix errors before pushing!

---

## 🛠️ Detailed Testing Layers

### 1. Backend (Go) 🔙
We use the standard `testing` package along with `testify` for assertions.

**Running Tests Manually:**
```bash
cd backend
go test -v ./...
```

**Writing a New Test:**
Create a file ending in `_test.go` next to the code you want to test.
```go
func TestHealthCheck(t *testing.T) {
    // Setup & Assertion
}
```

### 2. Frontend (Nuxt/Vue) 🎨
We rely on **Static Analysis (Type Checking)** and **Unit Tests**.

**Static Analysis:**
Ensures no type errors exist in your components.
```bash
cd frontend
npm run typecheck
```

**Unit Tests (Vitest):**
*Coming Soon: Component rendering tests.*

### 3. Infrastructure & Database 🗄️
**Health Check:**
Verify the API and Database connection are active.
```bash
curl -I http://localhost:8080/health
```

**Manual Data Verification:**
You can inspect the MongoDB container directly:
```bash
docker exec -it book-lapangan-mongo-1 mongosh booklapangan
```

---

## ✅ Pre-Merge Checklist
- [ ] Run `.\scripts\verify_build.ps1` and ensure it passes.
- [ ] Added unit tests for any new backend logic?
- [ ] Checked for TypeScript errors in new Vue components?
- [ ] Verified the app starts locally (`docker-compose up`).

---

*“Quality is not an act, it is a habit.”*

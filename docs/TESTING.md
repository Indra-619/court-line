# 🧪 Testing Guide

Ensure stability and reliability with our testing protocols. Currently, we focus on manual verification while automated suites are being established.

---

## 🕵️ Manual Verification

Since automated tests are under development, please verify the following critical paths manually.

### 1. 🔌 Backend Health Check
Ensure the Go server is up and database connection is active.

```bash
# Check simple health (if implemented) or root
curl -I http://localhost:8080/
```

**✅ Success Criteria:**
- Status Code: `200 OK`
- Logs show: `Connected to MongoDB`

### 2. 🌐 Frontend Functionality
Navigate to [http://localhost:3000](http://localhost:3000).

| Feature | Action | Expected Result |
| :--- | :--- | :--- |
| **Home Page** | Open root URL | Landing page loads with hero section |
| **Auth** | Click "Login with Google" | Redirects to Google, then back to app |
| **Bookings** | Click "Book Now" | Form appears, date selection works |

### 3. 💾 Data Persistence
1. **Create a Booking**: Submit a form on the frontend.
2. **Verify in DB**:
   ```bash
   # Using Docker
   docker exec -it book-lapangan-mongo-1 mongosh booklapangan --eval "db.bookings.find().pretty()"
   ```

---

## 🤖 Automated Testing (Coming Soon)

We are setting up standard testing frameworks for robust CI/CD.

### Backend (Go)
```bash
cd backend
go test -v ./...
```

### Frontend (Vitest)
```bash
cd frontend
npm run test
```

---

## 📝 Writing Tests

If you are adding a feature, **please** consider adding a unit test file `_test.go` adjacent to your logic.

```go
func TestCalculatePrice(t *testing.T) {
    // ...
}
```

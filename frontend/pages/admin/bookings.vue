<template>
  <div class="admin-page">
    <div class="container">
      <h1 class="page-title">Admin: Bookings</h1>
      <p class="page-subtitle">Confirm, complete, and review all court bookings</p>

      <div v-if="!isLoggedIn" class="empty-state">
        <p>Please login to access the admin dashboard.</p>
        <button @click="loginWithGoogle" class="btn btn-primary">Login with Google</button>
      </div>

      <div v-else-if="user?.role !== 'admin'" class="empty-state">
        <p>This area is for administrators only.</p>
        <NuxtLink to="/" class="btn btn-secondary">Back to Home</NuxtLink>
      </div>

      <div v-else-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>Loading bookings...</p>
      </div>

      <div v-else-if="error" class="error-state">{{ error }}</div>

      <div v-else-if="bookings.length === 0" class="empty-state">
        <p>No bookings yet.</p>
      </div>

      <div v-else class="admin-table-wrap">
        <table class="admin-table">
          <thead>
            <tr>
              <th>Customer</th>
              <th>Court</th>
              <th>Date</th>
              <th>Time</th>
              <th>Total</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="booking in bookings" :key="booking.id">
              <td>
                <div class="customer-name">{{ booking.customerName }}</div>
                <div class="customer-phone">{{ booking.customerPhone }}</div>
              </td>
              <td>{{ courtName(booking.courtId) }}</td>
              <td>{{ formatDate(booking.date) }}</td>
              <td>{{ booking.startTime }} - {{ booking.endTime }}</td>
              <td>Rp {{ formatPrice(booking.totalPrice) }}</td>
              <td>
                <span :class="['status-badge', `status-${booking.status}`]">{{ booking.status }}</span>
              </td>
              <td class="actions-cell">
                <button
                  v-if="booking.status === 'pending'"
                  class="btn btn-action btn-confirm"
                  :disabled="busyId === booking.id"
                  @click="setStatus(booking, 'confirmed')"
                >Confirm</button>
                <button
                  v-if="booking.status === 'pending' || booking.status === 'confirmed'"
                  class="btn btn-action btn-complete"
                  :disabled="busyId === booking.id"
                  @click="setStatus(booking, 'completed')"
                >Complete</button>
                <span v-if="booking.status === 'cancelled' || booking.status === 'completed'" class="muted">—</span>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="pager">
          <button class="btn btn-secondary btn-sm" :disabled="offset === 0" @click="prevPage">Previous</button>
          <span class="pager-info">Showing {{ bookings.length }} (offset {{ offset }})</span>
          <button class="btn btn-secondary btn-sm" :disabled="bookings.length < limit" @click="nextPage">Next</button>
        </div>
      </div>

      <div v-if="actionMessage" :class="['action-message', `action-${actionMessage.type}`]">
        {{ actionMessage.text }}
      </div>
    </div>
  </div>
</template>

<script setup>
const config = useRuntimeConfig()
const { isLoggedIn, user, loginWithGoogle, authFetch } = useAuth()

const bookings = ref([])
const courts = ref({})
const loading = ref(true)
const error = ref('')
const busyId = ref('')
const actionMessage = ref(null)
const limit = 50
const offset = ref(0)

const formatPrice = (price) => new Intl.NumberFormat('id-ID').format(price)

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' })
}

const courtName = (courtId) => courts.value[courtId]?.name || courtId

const fetchCourts = async () => {
  try {
    const response = await $fetch(`${config.public.apiBase}/api/courts?limit=100`)
    for (const court of response.data || []) {
      courts.value[court.id] = court
    }
  } catch {
    // names fall back to raw ids
  }
}

const fetchBookings = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await authFetch(
      `${config.public.apiBase}/api/admin/bookings?limit=${limit}&offset=${offset.value}`
    )
    bookings.value = response.data || []
  } catch (err) {
    error.value = err?.data?.error || 'Failed to load bookings (admin access required).'
    bookings.value = []
  } finally {
    loading.value = false
  }
}

const setStatus = async (booking, status) => {
  busyId.value = booking.id
  actionMessage.value = null
  try {
    await authFetch(`${config.public.apiBase}/api/admin/bookings/${booking.id}/status`, {
      method: 'PATCH',
      body: { status }
    })
    booking.status = status
    actionMessage.value = { type: 'success', text: `Booking marked ${status}.` }
  } catch (err) {
    actionMessage.value = { type: 'error', text: err?.data?.error || 'Failed to update status' }
  } finally {
    busyId.value = ''
  }
}

const prevPage = () => {
  offset.value = Math.max(0, offset.value - limit)
  fetchBookings()
}

const nextPage = () => {
  offset.value += limit
  fetchBookings()
}

onMounted(async () => {
  if (isLoggedIn.value && user.value?.role === 'admin') {
    await Promise.all([fetchCourts(), fetchBookings()])
  } else {
    loading.value = false
  }
})
</script>

<style scoped>
.admin-page {
  padding-top: 100px;
  padding-bottom: 4rem;
  min-height: 100vh;
}

.page-title {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.page-subtitle {
  color: var(--color-text-muted);
  margin-bottom: 2rem;
}

.admin-table-wrap {
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.admin-table th,
.admin-table td {
  padding: 0.9rem 1rem;
  text-align: left;
  font-size: 0.9rem;
  border-bottom: 1px solid var(--color-surface-light);
}

.admin-table th {
  color: var(--color-text-muted);
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.customer-name {
  font-weight: 600;
}

.customer-phone {
  color: var(--color-text-muted);
  font-size: 0.8rem;
}

.status-badge {
  padding: 0.25rem 0.6rem;
  border-radius: var(--radius-md);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-pending { background: rgba(245, 158, 11, 0.2); color: #f59e0b; }
.status-confirmed { background: rgba(34, 197, 94, 0.2); color: #22c55e; }
.status-cancelled { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
.status-completed { background: rgba(99, 102, 241, 0.2); color: #6366f1; }

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.btn-action {
  padding: 0.35rem 0.8rem;
  font-size: 0.8rem;
  border-radius: var(--radius-md);
  border: 1px solid transparent;
  cursor: pointer;
}

.btn-confirm {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  border-color: rgba(34, 197, 94, 0.4);
}

.btn-complete {
  background: rgba(99, 102, 241, 0.15);
  color: #6366f1;
  border-color: rgba(99, 102, 241, 0.4);
}

.btn-action:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.muted {
  color: var(--color-text-muted);
}

.pager {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.pager-info {
  color: var(--color-text-muted);
  font-size: 0.85rem;
}

.btn-sm {
  padding: 0.4rem 0.9rem;
  font-size: 0.8rem;
}

.action-message {
  margin-top: 1.5rem;
  padding: 0.9rem 1.2rem;
  border-radius: var(--radius-md);
  font-size: 0.9rem;
}

.action-success { background: rgba(34, 197, 94, 0.15); color: #22c55e; }
.action-error { background: rgba(239, 68, 68, 0.15); color: #ef4444; }

.loading-container,
.empty-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.error-state {
  color: #ef4444;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 3px solid var(--color-surface-light);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
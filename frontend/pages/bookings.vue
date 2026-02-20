<template>
  <div class="bookings-page">
    <div class="container">
      <div class="page-header">
        <h1>My Bookings</h1>
        <p class="text-muted">Manage your court reservations</p>
      </div>

      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>Loading your bookings...</p>
      </div>

      <div v-else-if="!isLoggedIn" class="login-prompt">
        <p>Please login to view your bookings</p>
        <button @click="loginWithGoogle" class="btn btn-primary">Login with Google</button>
      </div>

      <div v-else-if="bookings.length === 0" class="empty-state">
        <span class="empty-icon">📅</span>
        <h2>No bookings yet</h2>
        <p>Start by booking your first court!</p>
        <NuxtLink to="/" class="btn btn-primary">Browse Courts</NuxtLink>
      </div>

      <div v-else class="bookings-list">
        <div v-for="booking in bookings" :key="booking.id" class="booking-card">
          <div v-if="booking.status" class="booking-status" :class="`status-${booking.status}`">
            {{ formatStatus(booking.status) }}
          </div>
          <div class="booking-info">
            <h3>{{ booking.customerName }}</h3>
            <div class="booking-details">
              <p><strong>Date:</strong> {{ formatDate(booking.date) }}</p>
              <p><strong>Time:</strong> {{ booking.startTime }} - {{ booking.endTime }}</p>
              <p><strong>Phone:</strong> {{ booking.customerPhone }}</p>
            </div>
          </div>
          <div class="booking-price">
            <span v-if="booking.totalPrice !== undefined" class="price">Rp {{ formatPrice(booking.totalPrice) }}</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
const { isLoggedIn, loginWithGoogle } = useAuth()
import { BookingService, type Booking } from '~/services/booking.service'

const bookings = ref<Booking[]>([])
const loading = ref(true)

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatStatus = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'Pending',
    confirmed: 'Confirmed',
    cancelled: 'Cancelled',
    completed: 'Completed'
  }
  return statusMap[status] || status
}

onMounted(async () => {
  if (!isLoggedIn.value) {
    loading.value = false
    return
  }

  try {
    bookings.value = await BookingService.getMyBookings()
  } catch (error) {
    console.error('Failed to fetch bookings:', error)
  } finally {
    loading.value = false
  }
})
</script>


<style scoped>
.bookings-page {
  padding: 2rem 0 4rem;
  min-height: 80vh;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  margin-bottom: 0.5rem;
}

.loading-container,
.login-prompt,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 4rem;
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
}

.bookings-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.booking-card {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 1.5rem;
  align-items: center;
  background: var(--color-surface);
  padding: 1.5rem;
  border-radius: var(--radius-lg);
}

.booking-status {
  padding: 0.5rem 1rem;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 0.875rem;
  text-transform: uppercase;
}

.status-pending {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.status-confirmed {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.status-cancelled {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.status-completed {
  background: rgba(99, 102, 241, 0.2);
  color: #6366f1;
}

.booking-info h3 {
  margin-bottom: 0.5rem;
}

.booking-details {
  display: flex;
  gap: 2rem;
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.booking-details p {
  margin: 0;
}

@media (max-width: 768px) {
  .booking-card {
    grid-template-columns: 1fr;
  }
  
  .booking-details {
    flex-direction: column;
    gap: 0.25rem;
  }
}
</style>

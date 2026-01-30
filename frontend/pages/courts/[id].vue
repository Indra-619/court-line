<template>
  <div class="court-detail">
    <div class="container">
      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>Loading court details...</p>
      </div>

      <div v-else-if="!court" class="error-state">
        <h2>Court not found</h2>
        <NuxtLink to="/" class="btn btn-primary">Back to Home</NuxtLink>
      </div>

      <div v-else class="detail-grid">
        <!-- Court Info -->
        <div class="court-info">
          <NuxtLink to="/" class="back-link">
            ← Back to courts
          </NuxtLink>

          <div class="court-image">
            <img :src="court.imageUrl || defaultImage" :alt="court.name" />
            <span class="badge badge-primary">{{ court.type }}</span>
          </div>

          <h1 class="court-name">{{ court.name }}</h1>
          
          <div class="court-location">
            <span class="icon">📍</span>
            {{ court.location }}
          </div>

          <div class="court-price">
            <span class="price">Rp {{ formatPrice(court.pricePerHour) }}</span>
            <span class="price-unit">per hour</span>
          </div>

          <div class="court-description">
            <h3>Description</h3>
            <p>{{ court.description || 'A well-maintained sports venue with excellent facilities. Perfect for your games and practice sessions.' }}</p>
          </div>

          <div v-if="court.facilities?.length" class="court-facilities">
            <h3>Facilities</h3>
            <div class="facilities-list">
              <span v-for="facility in court.facilities" :key="facility" class="facility-tag">
                {{ facility }}
              </span>
            </div>
          </div>
        </div>

        <!-- Booking Form -->
        <div class="booking-section">
          <BookingForm 
            :court-id="court.id" 
            :price-per-hour="court.pricePerHour"
            @success="handleBookingSuccess"
            @error="handleBookingError"
          />

          <div v-if="toast.show" :class="['toast', `toast-${toast.type}`]">
            {{ toast.message }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const config = useRuntimeConfig()

const court = ref(null)
const loading = ref(true)
const defaultImage = 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=800&h=500&fit=crop'

const toast = ref({
  show: false,
  message: '',
  type: 'success'
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const showToast = (message, type = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

const handleBookingSuccess = (booking) => {
  showToast('Booking created successfully! Check your bookings.', 'success')
}

const handleBookingError = (error) => {
  showToast(error || 'Failed to create booking', 'error')
}

onMounted(async () => {
  const courtId = route.params.id
  
  try {
    const response = await $fetch(`${config.public.apiBase}/api/courts/${courtId}`)
    court.value = response.data
  } catch (error) {
    console.error('Failed to fetch court:', error)
    // Sample data for demo
    court.value = {
      id: courtId,
      name: 'Sample Court',
      type: 'Futsal',
      location: 'Jakarta',
      description: 'A modern futsal court with artificial turf and full facilities.',
      pricePerHour: 150000,
      imageUrl: defaultImage,
      facilities: ['Changing Room', 'Parking', 'Cafeteria', 'Shower'],
      isAvailable: true
    }
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.court-detail {
  padding: 2rem 0 4rem;
}

.back-link {
  display: inline-block;
  color: var(--color-text-muted);
  margin-bottom: 1.5rem;
  transition: color 0.2s;
}

.back-link:hover {
  color: var(--color-primary);
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 3rem;
  align-items: start;
}

.court-image {
  position: relative;
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.court-image img {
  width: 100%;
  height: 400px;
  object-fit: cover;
}

.court-image .badge {
  position: absolute;
  top: 1rem;
  left: 1rem;
}

.court-name {
  font-size: 2rem;
  margin-bottom: 0.75rem;
}

.court-location {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--color-text-muted);
  font-size: 1.125rem;
  margin-bottom: 1rem;
}

.court-price {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  margin-bottom: 2rem;
}

.court-price .price {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-secondary);
}

.court-price .price-unit {
  color: var(--color-text-muted);
}

.court-description,
.court-facilities {
  margin-bottom: 2rem;
}

.court-description h3,
.court-facilities h3 {
  font-size: 1.125rem;
  margin-bottom: 0.75rem;
}

.court-description p {
  color: var(--color-text-muted);
  line-height: 1.7;
}

.facilities-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.facility-tag {
  padding: 0.5rem 1rem;
  background: var(--color-surface-light);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
}

.booking-section {
  position: sticky;
  top: 100px;
}

.loading-container,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 4rem;
  text-align: center;
}

@media (max-width: 1024px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
  
  .booking-section {
    position: static;
  }
}
</style>

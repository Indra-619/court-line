<template>
  <div>
    <!-- Hero Section -->
    <section class="hero">
      <div class="container hero-content">
        <h1 class="hero-title">
          Book Your Favorite
          <span class="gradient-text">Sports Court</span>
        </h1>
        <p class="hero-subtitle">
          Find and reserve the best sports venues in your area. 
          Fast, easy, and hassle-free booking experience.
        </p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">50+</span>
            <span class="stat-label">Courts</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">1000+</span>
            <span class="stat-label">Bookings</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">24/7</span>
            <span class="stat-label">Available</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Courts Section -->
    <section class="courts-section">
      <div class="container">
        <div class="section-header">
          <h2>Available Courts</h2>
          <p class="text-muted">Choose from our wide selection of sports venues</p>
        </div>

        <div v-if="loading" class="loading-container">
          <div class="spinner"></div>
          <p>Loading courts...</p>
        </div>

        <div v-else-if="courts.length === 0" class="empty-state">
          <p>No courts available at the moment.</p>
        </div>

        <div v-else class="grid grid-cols-3">
          <CourtCard v-for="court in courts" :key="court.id" :court="court" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { initAuth } = useAuth()
import { CourtService, type Court } from '~/services/court.service'

const courts = ref<Court[]>([])
const loading = ref(true)

onMounted(async () => {
  // Initialize auth
  initAuth()
  
  // Fetch courts
  try {
    courts.value = await CourtService.getAll()
  } catch (error) {
    console.error('Failed to fetch courts:', error)
  } finally {
    loading.value = false
  }
})
</script>


<style scoped>
.hero {
  background: linear-gradient(135deg, var(--color-surface) 0%, var(--color-background) 100%);
  padding: 5rem 0;
  position: relative;
  overflow: hidden;
}

.hero::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.15) 0%, transparent 70%);
  border-radius: 50%;
}

.hero::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.hero-content {
  position: relative;
  z-index: 1;
  text-align: center;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  margin-bottom: 1.5rem;
  line-height: 1.1;
}

.gradient-text {
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--color-text-muted);
  max-width: 600px;
  margin: 0 auto 3rem;
}

.hero-stats {
  display: flex;
  justify-content: center;
  gap: 4rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-number {
  font-size: 2.5rem;
  font-weight: 800;
  color: var(--color-primary);
}

.stat-label {
  font-size: 1rem;
  color: var(--color-text-muted);
}

.courts-section {
  padding: 4rem 0;
}

.section-header {
  text-align: center;
  margin-bottom: 3rem;
}

.section-header h2 {
  margin-bottom: 0.5rem;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 3rem;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--color-text-muted);
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2.5rem;
  }
  
  .hero-stats {
    gap: 2rem;
  }
  
  .stat-number {
    font-size: 1.75rem;
  }
}
</style>

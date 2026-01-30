<template>
  <NuxtLink :to="`/courts/${court.id}`" class="court-card card">
    <div class="card-image">
      <img :src="court.imageUrl || defaultImage" :alt="court.name" />
      <span class="badge badge-primary">{{ court.type }}</span>
    </div>
    <div class="card-content">
      <h3 class="card-title">{{ court.name }}</h3>
      <p class="card-location">
        <span class="location-icon">📍</span>
        {{ court.location }}
      </p>
      <div class="card-footer">
        <div class="price">
          Rp {{ formatPrice(court.pricePerHour) }}
          <span class="price-unit">/jam</span>
        </div>
        <span v-if="court.isAvailable" class="badge badge-success">Available</span>
        <span v-else class="badge badge-warning">Booked</span>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup>
const props = defineProps({
  court: {
    type: Object,
    required: true
  }
})

const defaultImage = 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=400&h=300&fit=crop'

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>

<style scoped>
.court-card {
  display: block;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.court-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 20px 40px rgba(99, 102, 241, 0.15);
}

.card-image {
  position: relative;
  height: 180px;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.court-card:hover .card-image img {
  transform: scale(1.05);
}

.card-image .badge {
  position: absolute;
  top: 1rem;
  left: 1rem;
}

.card-content {
  padding: 1.25rem;
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-text);
}

.card-location {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--color-text-muted);
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.location-icon {
  font-size: 0.875rem;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

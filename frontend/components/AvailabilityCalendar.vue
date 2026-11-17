<template>
  <div class="availability">
    <div class="availability-header">
      <h3>Availability</h3>
      <input type="date" v-model="selectedDate" :min="today" class="date-input" />
    </div>

    <div v-if="loading" class="availability-loading">Checking schedule...</div>

    <div v-else class="slot-grid">
      <div
        v-for="slot in slots"
        :key="slot.hour"
        class="slot"
        :class="{ 'slot-booked': slot.booked, 'slot-past': slot.past }"
        :title="slot.booked ? 'Booked' : slot.past ? 'Past' : 'Available'"
      >
        {{ slot.label }}
      </div>
    </div>

    <div class="availability-legend">
      <span class="legend-item"><span class="dot dot-free"></span> Available</span>
      <span class="legend-item"><span class="dot dot-booked"></span> Booked</span>
      <span class="legend-item"><span class="dot dot-past"></span> Past</span>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  courtId: { type: String, required: true }
})

const config = useRuntimeConfig()

const pad = (n) => String(n).padStart(2, '0')
const localDate = (d) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`

const today = localDate(new Date())
const selectedDate = ref(today)
const bookings = ref([])
const loading = ref(false)

// Operating hours 06:00-23:00, one-hour slots.
const OPEN_HOUR = 6
const CLOSE_HOUR = 23

const toMinutes = (hhmm) => {
  const [h, m] = hhmm.split(':').map(Number)
  return h * 60 + m
}

const slots = computed(() => {
  const now = new Date()
  const isToday = selectedDate.value === localDate(now)
  const nowMinutes = now.getHours() * 60 + now.getMinutes()

  const active = bookings.value.filter(
    (b) => b.status === 'pending' || b.status === 'confirmed'
  )

  const result = []
  for (let hour = OPEN_HOUR; hour < CLOSE_HOUR; hour++) {
    const start = hour * 60
    const end = start + 60
    const booked = active.some((b) => {
      const bs = toMinutes(b.startTime)
      const be = toMinutes(b.endTime)
      return bs < end && be > start
    })
    result.push({
      hour,
      label: `${pad(hour)}:00`,
      booked,
      past: isToday && end <= nowMinutes
    })
  }
  return result
})

const fetchBookings = async () => {
  loading.value = true
  try {
    const response = await $fetch(
      `${config.public.apiBase}/api/courts/${props.courtId}/bookings?date=${selectedDate.value}&limit=100`
    )
    bookings.value = (response.data || []).filter((b) => b.date === selectedDate.value)
  } catch {
    bookings.value = []
  } finally {
    loading.value = false
  }
}

watch(selectedDate, fetchBookings)
onMounted(fetchBookings)
</script>

<style scoped>
.availability {
  margin-top: 2rem;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
}

.availability-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  gap: 1rem;
  flex-wrap: wrap;
}

.availability-header h3 {
  margin: 0;
  font-size: 1.125rem;
}

.date-input {
  background: var(--color-surface-light);
  border: 1px solid var(--color-border, rgba(255, 255, 255, 0.1));
  color: var(--color-text);
  border-radius: var(--radius-md);
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
}

.availability-loading {
  color: var(--color-text-muted);
  font-size: 0.9rem;
  padding: 1rem 0;
}

.slot-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(72px, 1fr));
  gap: 0.5rem;
}

.slot {
  text-align: center;
  padding: 0.5rem 0.25rem;
  border-radius: var(--radius-md);
  font-size: 0.8rem;
  font-weight: 600;
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.slot-booked {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
}

.slot-past {
  background: var(--color-surface-light);
  color: var(--color-text-muted);
  border-color: transparent;
  opacity: 0.6;
}

.availability-legend {
  display: flex;
  gap: 1.25rem;
  margin-top: 1rem;
  font-size: 0.8rem;
  color: var(--color-text-muted);
  flex-wrap: wrap;
}

.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
}

.dot-free { background: #22c55e; }
.dot-booked { background: #ef4444; }
.dot-past { background: var(--color-text-muted); }
</style>
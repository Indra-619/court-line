<template>
  <div class="booking-form">
    <h3 class="form-title">Book This Court</h3>
    
    <form @submit.prevent="submitBooking">
      <div class="form-group">
        <label class="form-label">Your Name</label>
        <input 
          v-model="form.customerName" 
          type="text" 
          class="form-input" 
          placeholder="Enter your name"
          required
        />
      </div>

      <div class="form-group">
        <label class="form-label">Phone Number</label>
        <input 
          v-model="form.customerPhone" 
          type="tel" 
          class="form-input" 
          placeholder="08xx xxxx xxxx"
          required
        />
      </div>

      <div class="form-group">
        <label class="form-label">Date</label>
        <input 
          v-model="form.date" 
          type="date" 
          class="form-input"
          :min="minDate"
          required
        />
      </div>

      <div class="time-row">
        <div class="form-group">
          <label class="form-label">Start Time</label>
          <select v-model="form.startTime" class="form-input" required>
            <option value="">Select</option>
            <option v-for="time in timeSlots" :key="time" :value="time">{{ time }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">End Time</label>
          <select v-model="form.endTime" class="form-input" required>
            <option value="">Select</option>
            <option v-for="time in endTimeSlots" :key="time" :value="time">{{ time }}</option>
          </select>
        </div>
      </div>

      <div v-if="estimatedPrice > 0" class="price-estimate">
        <span>Estimated Total:</span>
        <span class="price">Rp {{ formatPrice(estimatedPrice) }}</span>
      </div>

      <button type="submit" class="btn btn-primary btn-full" :disabled="loading">
        <span v-if="loading" class="spinner-sm"></span>
        <span v-else>{{ isLoggedIn ? 'Confirm Booking' : 'Login to Book' }}</span>
      </button>
    </form>
  </div>
</template>

<script setup>
const props = defineProps({
  courtId: {
    type: String,
    required: true
  },
  pricePerHour: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['success', 'error'])

const { isLoggedIn, loginWithGoogle, getAuthHeader, user } = useAuth()
const config = useRuntimeConfig()

const loading = ref(false)

const form = ref({
  customerName: '',
  customerPhone: '',
  date: '',
  startTime: '',
  endTime: ''
})

// Watch for user data to pre-fill name
watch(() => user.value, (newUser) => {
  if (newUser?.name) {
    form.value.customerName = newUser.name
  }
}, { immediate: true })

const minDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

const timeSlots = computed(() => {
  const slots = []
  for (let i = 6; i <= 22; i++) {
    slots.push(`${i.toString().padStart(2, '0')}:00`)
  }
  return slots
})

const endTimeSlots = computed(() => {
  if (!form.value.startTime) return []
  const startHour = parseInt(form.value.startTime.split(':')[0])
  const slots = []
  for (let i = startHour + 1; i <= 23; i++) {
    slots.push(`${i.toString().padStart(2, '0')}:00`)
  }
  return slots
})

const estimatedPrice = computed(() => {
  if (!form.value.startTime || !form.value.endTime) return 0
  const startHour = parseInt(form.value.startTime.split(':')[0])
  const endHour = parseInt(form.value.endTime.split(':')[0])
  const hours = endHour - startHour
  return hours * props.pricePerHour
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const submitBooking = async () => {
  if (!isLoggedIn.value) {
    loginWithGoogle()
    return
  }

  loading.value = true

  try {
    const response = await $fetch(`${config.public.apiBase}/api/bookings`, {
      method: 'POST',
      headers: {
        ...getAuthHeader(),
        'Content-Type': 'application/json'
      },
      body: {
        courtId: props.courtId,
        customerName: form.value.customerName,
        customerPhone: form.value.customerPhone,
        date: form.value.date,
        startTime: form.value.startTime,
        endTime: form.value.endTime
      }
    })

    emit('success', response.data)
    
    // Reset form
    form.value = {
      customerName: user.value?.name || '',
      customerPhone: '',
      date: '',
      startTime: '',
      endTime: ''
    }
  } catch (error) {
    console.error('Booking failed:', error)
    emit('error', error.message || 'Failed to create booking')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.booking-form {
  background: var(--color-surface);
  padding: 1.5rem;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
}

.form-title {
  font-size: 1.25rem;
  margin-bottom: 1.5rem;
  color: var(--color-text);
}

.time-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.price-estimate {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: var(--color-surface-light);
  border-radius: var(--radius-md);
  margin-bottom: 1.5rem;
}

.btn-full {
  width: 100%;
}

.spinner-sm {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>

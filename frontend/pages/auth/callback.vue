<template>
  <div class="auth-callback">
    <div class="container">
      <div class="callback-content">
        <div v-if="error" class="error-state">
          <span class="icon">❌</span>
          <h2>Authentication Failed</h2>
          <p>{{ error }}</p>
          <NuxtLink to="/" class="btn btn-primary">Back to Home</NuxtLink>
        </div>
        <div v-else class="loading-state">
          <div class="spinner"></div>
          <h2>Signing you in...</h2>
          <p>Please wait while we complete your authentication.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()
const { setToken } = useAuth()

const error = ref('')

onMounted(async () => {
  const code = route.query.code

  if (!code) {
    error.value = 'No authentication code received. Please try again.'
    return
  }

  try {
    const response = await $fetch<{ data: { token: string; refreshToken: string } }>(
      `${config.public.apiBase}/auth/exchange`,
      {
        method: 'POST',
        body: { code: code as string }
      }
    )

    setToken(response.data.token, response.data.refreshToken)
    // Redirect to home or previous page
    const redirect = route.query.redirect || '/'
    router.push(redirect as string)
  } catch (e) {
    error.value = 'Authentication failed or the code expired. Please try again.'
  }
})
</script>

<style scoped>
.auth-callback {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.callback-content {
  text-align: center;
  padding: 3rem;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.loading-state h2,
.error-state h2 {
  margin-top: 1rem;
}

.loading-state p,
.error-state p {
  color: var(--color-text-muted);
}

.error-state .icon {
  font-size: 3rem;
}
</style>

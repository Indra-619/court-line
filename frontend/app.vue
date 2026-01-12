<template>
  <div class="container">
    <header>
      <h1>Book Lapangan Online</h1>
      <p>Find and book your favorite sports venue.</p>
    </header>
    
    <main>
      <div v-if="healthStatus" class="status-box success">
        Backend Status: {{ healthStatus }}
      </div>
      <div v-else class="status-box loading">
        Checking backend...
      </div>

      <div class="content">
        <p>Welcome to the booking system. This is a Nuxt 3 frontend connected to a Go backend.</p>
        <button class="cta-button" @click="checkHealth">Check Connection</button>
      </div>
    </main>
  </div>
</template>

<script setup>
const config = useRuntimeConfig()
const healthStatus = ref('')

async function checkHealth() {
  try {
    // Proxy request would be better in production, but for demo we try direct if CORS allows or via server proxy
    // In docker-compose, client browser can't reach 'backend' hostname directly, 
    // so we usually need a Nuxt server route to proxy or rely on localhost mapping if running locally.
    // For this setup, we assume Nginx or simple fetch to localhost:8080 if exposed.
    
    // However, since we are inside docker network for server-side rendering, 
    // useFetch might work differently on server vs client.
    // For client-side fetch from browser to localhost:8080:
    const data = await $fetch(`${config.public.apiBase}/health`, {
        mode: 'cors'
    })
    healthStatus.value = data
  } catch (error) {
    console.error('Health check failed', error)
    healthStatus.value = 'Error connecting to backend'
  }
}

onMounted(() => {
    checkHealth()
})
</script>

<style>
/* Reset */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f0f2f5;
  color: #333;
  line-height: 1.6;
}

.container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

header {
  text-align: center;
  margin-bottom: 2rem;
}

h1 {
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.status-box {
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  text-align: center;
  font-weight: bold;
}

.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.loading {
  background-color: #e2e3e5;
  color: #383d41;
  border: 1px solid #d6d8db;
}

.content {
  text-align: center;
}

.cta-button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  margin-top: 1rem;
  transition: background-color 0.2s;
}

.cta-button:hover {
  background-color: #0056b3;
}
</style>

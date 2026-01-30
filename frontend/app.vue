<template>
  <div id="app">
    <Navbar />
    <main>
      <NuxtPage />
    </main>
  </div>
</template>

<script setup>
const { initAuth } = useAuth()

onMounted(() => {
  initAuth()
})
</script>

<style>
/* Import global styles */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');
</style>

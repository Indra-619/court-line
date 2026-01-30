<template>
  <nav class="navbar">
    <div class="container navbar-content">
      <NuxtLink to="/" class="navbar-brand">
        <span class="logo-icon">🏟️</span>
        <span class="logo-text">Court<span class="logo-highlight">Line</span></span>
      </NuxtLink>

      <div class="navbar-menu">
        <NuxtLink to="/" class="nav-link">Home</NuxtLink>
        <NuxtLink v-if="isLoggedIn" to="/bookings" class="nav-link">My Bookings</NuxtLink>
      </div>

      <div class="navbar-auth">
        <template v-if="isLoggedIn">
          <div class="user-menu">
            <img :src="user?.picture" :alt="user?.name" class="user-avatar" />
            <span class="user-name">{{ user?.name }}</span>
            <button @click="logout" class="btn btn-secondary btn-sm">Logout</button>
          </div>
        </template>
        <template v-else>
          <button @click="loginWithGoogle" class="btn btn-google">
            <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            Login with Google
          </button>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
const { user, isLoggedIn, loginWithGoogle, logout } = useAuth()
</script>

<style scoped>
.navbar {
  background: rgba(30, 41, 59, 0.95);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.5rem;
  font-weight: 700;
}

.logo-icon {
  font-size: 1.75rem;
}

.logo-text {
  color: var(--color-text);
}

.logo-highlight {
  color: var(--color-primary);
}

.navbar-menu {
  display: flex;
  gap: 2rem;
}

.nav-link {
  color: var(--color-text-muted);
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--color-primary);
}

.navbar-auth {
  display: flex;
  align-items: center;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 2px solid var(--color-primary);
}

.user-name {
  font-weight: 500;
  color: var(--color-text);
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .navbar-menu {
    display: none;
  }
  
  .user-name {
    display: none;
  }
}
</style>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navbar -->
    <nav class="bg-white border-b border-gray-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <!-- Logo -->
        <RouterLink to="/" class="flex items-center gap-2">
          <span class="text-2xl">⚽</span>
          <span class="font-bold text-gray-800 text-lg">League Simulation</span>
        </RouterLink>

        <!-- Nav Links -->
        <div class="flex items-center gap-4">
          <RouterLink
            to="/"
            class="text-sm font-medium text-gray-600 hover:text-primary transition-colors"
            active-class="text-primary"
          >
            Liga
          </RouterLink>
          <RouterLink
            to="/stats"
            class="text-sm font-medium text-gray-600 hover:text-primary transition-colors"
            active-class="text-primary"
          >
            Statistics
          </RouterLink>

          <!-- Admin -->
          <template v-if="!authStore.isAuthenticated">
            <RouterLink
              to="/login"
              class="text-sm bg-primary text-white px-3 py-1.5 rounded-lg hover:bg-blue-700 transition-colors"
            >
              Admin Login
            </RouterLink>
          </template>
          <template v-else>
            <span class="text-sm text-gray-500">👤 Admin</span>
            <button @click="authStore.logout()" class="text-sm text-danger hover:underline">
              Exit
            </button>
          </template>
        </div>
      </div>
    </nav>

    <!-- Page Content -->
    <main class="max-w-7xl mx-auto px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center">
    <div class="w-full max-w-sm">
      <!-- Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8">
        <!-- Logo -->
        <div class="text-center mb-8">
          <span class="text-5xl">⚽</span>
          <h1 class="text-xl font-bold text-gray-800 mt-3">Admin Login</h1>
          <p class="text-sm text-gray-400 mt-1">League Simulation Panel</p>
        </div>

        <!-- Error -->
        <div
          v-if="error"
          class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm flex items-center gap-2"
        >
          <span>⚠️</span>
          <span>{{ error }}</span>
        </div>

        <!-- Form -->
        <div class="space-y-4">
          <!-- Username -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5"> Username </label>
            <input
              v-model="form.username"
              type="text"
              placeholder="admin"
              @keyup.enter="handleLogin"
              class="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
              :class="{ 'border-red-300': error }"
            />
          </div>

          <!-- Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5"> Password </label>
            <div class="relative">
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
                @keyup.enter="handleLogin"
                class="w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all pr-10"
                :class="{ 'border-red-300': error }"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 text-sm"
              >
                {{ showPassword ? '🙈' : '👁️' }}
              </button>
            </div>
          </div>

          <!-- Submit -->
          <button
            @click="handleLogin"
            :disabled="loading || !form.username || !form.password"
            class="w-full bg-primary text-white py-2.5 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2 mt-2"
          >
            <span v-if="loading" class="animate-spin">⏳</span>
            <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
          </button>
        </div>

        <!-- Hint -->
        <div class="mt-6 p-3 bg-gray-50 rounded-lg text-center">
          <p class="text-xs text-gray-400">
            Default credentials:
            <span class="font-mono text-gray-600">admin</span> /
            <span class="font-mono text-gray-600">admin123</span>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api/auth'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref(null)
const showPassword = ref(false)

async function handleLogin() {
  if (!form.username || !form.password) return

  loading.value = true
  error.value = null

  try {
    const res = await authApi.login(form.username, form.password)
    authStore.setToken(res.data.token)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Invalid username or password'
    form.password = ''
  } finally {
    loading.value = false
  }
}
</script>

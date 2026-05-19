<template>
  <div class="space-y-6">
    <!-- Page title -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">⚽ League Simulation</h1>
      <p class="text-sm text-gray-500 mt-1">4-team Premier League style simulation</p>
    </div>

    <!-- Error -->
    <div
      v-if="store.error"
      class="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm flex items-center justify-between"
    >
      <span>⚠️ {{ store.error }}</span>
      <button @click="store.setError(null)" class="text-red-400 hover:text-red-600">✕</button>
    </div>

    <!-- Main grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left — Table (2/3) -->
      <div class="lg:col-span-2 space-y-6">
        <StandingsTable :standings="store.standings" :current-week="store.currentWeek" />
      </div>

      <!-- Right — Controls (1/3) -->
      <div class="space-y-4">
        <WeekControls
          :status="store.status"
          :loading="store.loading"
          :is-admin="authStore.isAuthenticated"
          @next-week="handleNextWeek"
          @play-all="handlePlayAll"
          @reset="handleReset"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useLeagueStore } from '../stores/league'
import { useAuthStore } from '../stores/auth'
import { useLeague } from '../composables/useLeague'
import StandingsTable from '../components/league/StandingsTable.vue'
import WeekControls from '../components/league/WeekControls.vue'

const store = useLeagueStore()
const authStore = useAuthStore()
const league = useLeague()

onMounted(async () => {
  await Promise.all([league.fetchTable(), league.fetchStatus()])
})

async function handleNextWeek() {
  await league.playNextWeek()
}

async function handlePlayAll() {
  await league.playAll()
}

async function handleReset() {
  if (confirm('Reset the league? All match results will be lost.')) {
    await league.resetLeague()
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Title -->
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
      <!-- Left — 2/3 -->
      <div class="lg:col-span-2 space-y-6">
        <StandingsTable :standings="store.standings" :current-week="store.currentWeek" />
        <WeekResults
          :week-results="store.weekResults"
          :is-admin="authStore.isAuthenticated"
          @edit-match="openEditModal"
        />
      </div>

      <!-- Right — 1/3 -->
      <div class="space-y-4">
        <WeekControls
          :status="store.status"
          :loading="store.loading"
          :is-admin="authStore.isAuthenticated"
          @next-week="handleNextWeek"
          @play-all="handlePlayAll"
          @reset="handleReset"
        />
        <PredictionChart :predictions="store.predictions" :current-week="store.currentWeek" />
      </div>
    </div>

    <!-- Edit Modal -->
    <EditMatchModal
      :show="editModal.show"
      :match="editModal.match"
      :loading="editModal.loading"
      :error="editModal.error"
      @close="closeEditModal"
      @save="handleEditSave"
    />
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { useLeagueStore } from '../stores/league'
import { useAuthStore } from '../stores/auth'
import { useLeague } from '../composables/useLeague'
import { matchApi } from '../api/match'
import StandingsTable from '../components/league/StandingsTable.vue'
import PredictionChart from '../components/league/PredictionChart.vue'
import WeekControls from '../components/league/WeekControls.vue'
import WeekResults from '../components/match/WeekResults.vue'
import EditMatchModal from '../components/match/EditMatchModal.vue'

const store = useLeagueStore()
const authStore = useAuthStore()
const league = useLeague()

const editModal = reactive({
  show: false,
  match: null,
  loading: false,
  error: null,
})

onMounted(async () => {
  await league.fetchAll()
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

function openEditModal(match) {
  editModal.match = match
  editModal.error = null
  editModal.show = true
}

function closeEditModal() {
  editModal.show = false
  editModal.match = null
  editModal.error = null
}

async function handleEditSave({ matchId, homeGoals, awayGoals }) {
  editModal.loading = true
  editModal.error = null
  try {
    await matchApi.updateResult(matchId, homeGoals, awayGoals)
    closeEditModal()
    await league.fetchTable()
  } catch (err) {
    editModal.error = err.response?.data?.error || 'Failed to update match result'
  } finally {
    editModal.loading = false
  }
}
</script>

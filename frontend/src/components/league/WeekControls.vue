<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-5">
    <!-- Status -->
    <div class="flex items-center justify-between mb-4">
      <div>
        <h3 class="font-bold text-gray-800">Match Controls</h3>
        <p class="text-xs text-gray-400 mt-0.5">
          {{ statusText }}
        </p>
      </div>

      <!-- Status badge -->
      <span class="text-xs font-medium px-2.5 py-1 rounded-full" :class="statusBadge">
        {{ status?.status?.replace('_', ' ').toUpperCase() || 'LOADING' }}
      </span>
    </div>

    <!-- Progress bar -->
    <div class="mb-4">
      <div class="flex justify-between text-xs text-gray-400 mb-1">
        <span>Week {{ status?.current_week || 0 }} of 6</span>
        <span>{{ status?.matches_played || 0 }} / 12 matches played</span>
      </div>
      <div class="w-full bg-gray-100 rounded-full h-2">
        <div
          class="bg-primary rounded-full h-2 transition-all duration-500"
          :style="{ width: progressPercent + '%' }"
        ></div>
      </div>
    </div>

    <!-- Buttons -->
    <div class="flex gap-2">
      <button
        @click="$emit('next-week')"
        :disabled="loading || isFinished"
        class="flex-1 bg-primary text-white text-sm font-medium py-2.5 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
      >
        <span v-if="loading" class="animate-spin">⏳</span>
        <span v-else>▶</span>
        Next Week
      </button>

      <button
        @click="$emit('play-all')"
        :disabled="loading || isFinished"
        class="flex-1 bg-secondary text-white text-sm font-medium py-2.5 px-4 rounded-lg hover:bg-purple-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
      >
        <span v-if="loading" class="animate-spin">⏳</span>
        <span v-else>⏩</span>
        Play All
      </button>

      <!-- Reset — admin only -->
      <button
        v-if="isAdmin"
        @click="$emit('reset')"
        :disabled="loading"
        class="bg-danger text-white text-sm font-medium py-2.5 px-3 rounded-lg hover:bg-red-700 disabled:opacity-40 transition-colors"
        title="Reset League"
      >
        🔄
      </button>
    </div>

    <!-- Finished banner -->
    <div
      v-if="isFinished"
      class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg text-center"
    >
      <p class="text-sm font-medium text-yellow-800">🏆 Season Complete! Reset to play again.</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: { type: Object, default: null },
  loading: { type: Boolean, default: false },
  isAdmin: { type: Boolean, default: false },
})

defineEmits(['next-week', 'play-all', 'reset'])

const isFinished = computed(() => props.status?.status === 'finished')

const progressPercent = computed(() => {
  const played = props.status?.matches_played || 0
  return Math.round((played / 12) * 100)
})

const statusText = computed(() => {
  const s = props.status
  if (!s) return 'Loading...'
  if (s.status === 'not_started') return 'Season not started yet'
  if (s.status === 'finished') return 'Season complete'
  return `${s.matches_left} matches remaining`
})

const statusBadge = computed(() => {
  const s = props.status?.status
  if (s === 'finished') return 'bg-yellow-100 text-yellow-700'
  if (s === 'in_progress') return 'bg-green-100 text-green-700'
  return 'bg-gray-100 text-gray-500'
})
</script>

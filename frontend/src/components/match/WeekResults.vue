<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Header -->
    <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
      <h2 class="font-bold text-gray-800 text-lg">📅 Match Results</h2>

      <!-- Week selector -->
      <div class="flex items-center gap-1">
        <button
          v-for="w in playedWeeks"
          :key="w"
          @click="selectedWeek = w"
          class="w-7 h-7 rounded-lg text-xs font-medium transition-colors"
          :class="
            selectedWeek === w
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          "
        >
          {{ w }}
        </button>
      </div>
    </div>

    <!-- Matches -->
    <div class="p-4">
      <div v-if="currentMatches.length" class="grid grid-cols-1 gap-3">
        <MatchCard
          v-for="match in currentMatches"
          :key="match.id"
          :match="match"
          :is-admin="isAdmin"
          @edit="$emit('edit-match', $event)"
        />
      </div>

      <!-- No results yet -->
      <div v-else class="py-10 text-center text-gray-400">
        <p class="text-3xl mb-2">🏟️</p>
        <p class="text-sm">No matches played yet.</p>
        <p class="text-xs mt-1">Click "Next Week" to start the season.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import MatchCard from './MatchCard.vue'

const props = defineProps({
  weekResults: { type: Array, default: () => [] },
  isAdmin: { type: Boolean, default: false },
})

defineEmits(['edit-match'])

const selectedWeek = ref(null)

const playedWeeks = computed(() => props.weekResults.map((w) => w.week).sort((a, b) => a - b))

watch(
  playedWeeks,
  (weeks) => {
    if (weeks.length) {
      selectedWeek.value = weeks[weeks.length - 1]
    }
  },
  { immediate: true },
)

const currentMatches = computed(() => {
  if (!selectedWeek.value) return []
  const week = props.weekResults.find((w) => w.week === selectedWeek.value)
  return week?.matches || []
})
</script>

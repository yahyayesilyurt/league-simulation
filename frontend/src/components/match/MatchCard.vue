<template>
  <div
    class="bg-white rounded-xl border border-gray-100 shadow-sm p-4 hover:shadow-md transition-shadow duration-200"
    :class="{ 'opacity-60': !match.played }"
  >
    <!-- Week badge -->
    <div class="text-xs text-gray-400 font-medium mb-3 text-center">Week {{ match.week }}</div>

    <!-- Teams & Score -->
    <div class="flex items-center justify-between gap-3">
      <!-- Home team -->
      <div class="flex-1 text-right">
        <div class="flex items-center justify-end gap-2">
          <span class="font-semibold text-gray-800 text-sm">
            {{ match.home_team?.name }}
          </span>
          <TeamLogo :name="match.home_team?.name" :size="32" />
        </div>
        <div class="text-xs text-gray-400 mt-0.5">Home</div>
      </div>

      <!-- Score -->
      <div class="flex items-center gap-2 min-w-20 justify-center">
        <template v-if="match.played">
          <span class="text-2xl font-bold w-8 text-center" :class="homeScoreClass">
            {{ match.home_goals }}
          </span>
          <span class="text-gray-300 font-light text-xl">—</span>
          <span class="text-2xl font-bold w-8 text-center" :class="awayScoreClass">
            {{ match.away_goals }}
          </span>
        </template>
        <template v-else>
          <span class="text-sm text-gray-400 font-medium">vs</span>
        </template>
      </div>

      <!-- Away team -->
      <div class="flex-1 text-left">
        <div class="flex items-center gap-2">
          <TeamLogo :name="match.away_team?.name" :size="32" />
          <span class="font-semibold text-gray-800 text-sm">
            {{ match.away_team?.name }}
          </span>
        </div>
        <div class="text-xs text-gray-400 mt-0.5">Away</div>
      </div>
    </div>

    <!-- Result label -->
    <div v-if="match.played" class="mt-3 text-center">
      <span class="text-xs font-medium px-2 py-0.5 rounded-full" :class="resultBadge">
        {{ resultLabel }}
      </span>
    </div>

    <!-- Not played -->
    <div v-else class="mt-3 text-center">
      <span class="text-xs text-gray-400">Not played yet</span>
    </div>

    <!-- Admin edit button -->
    <div v-if="isAdmin && match.played" class="mt-3 text-center">
      <button
        @click="$emit('edit', match)"
        class="text-xs text-primary hover:underline font-medium"
      >
        ✏️ Edit Result
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import TeamLogo from '../ui/TeamLogo.vue'

const props = defineProps({
  match: { type: Object, required: true },
  isAdmin: { type: Boolean, default: false },
})

defineEmits(['edit'])

const TEAM_EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

function teamEmoji(name) {
  return TEAM_EMOJIS[name] || '⚽'
}

const homeWon = computed(
  () => props.match.played && props.match.home_goals > props.match.away_goals,
)
const awayWon = computed(
  () => props.match.played && props.match.away_goals > props.match.home_goals,
)
const isDraw = computed(
  () => props.match.played && props.match.home_goals === props.match.away_goals,
)

const homeScoreClass = computed(() => {
  if (homeWon.value) return 'text-success'
  if (awayWon.value) return 'text-danger'
  return 'text-gray-700'
})

const awayScoreClass = computed(() => {
  if (awayWon.value) return 'text-success'
  if (homeWon.value) return 'text-danger'
  return 'text-gray-700'
})

const resultLabel = computed(() => {
  if (homeWon.value) return `${props.match.home_team?.name} Win`
  if (awayWon.value) return `${props.match.away_team?.name} Win`
  return 'Draw'
})

const resultBadge = computed(() => {
  if (isDraw.value) return 'bg-gray-100 text-gray-600'
  return 'bg-green-100 text-green-700'
})
</script>

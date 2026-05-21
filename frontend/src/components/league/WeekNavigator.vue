<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Header -->
    <div class="px-5 py-4 border-b border-gray-100">
      <h2 class="font-bold text-gray-800 text-lg">📆 Season Timeline</h2>
      <p class="text-xs text-gray-400 mt-0.5">Browse results week by week</p>
    </div>

    <div class="p-4">
      <!-- Week pills -->
      <div class="flex gap-2 flex-wrap mb-4">
        <button
          v-for="week in 6"
          :key="week"
          @click="selectWeek(week)"
          :disabled="!isWeekPlayed(week)"
          class="flex flex-col items-center px-3 py-2 rounded-xl text-xs font-medium transition-all duration-200 min-w-13"
          :class="weekPillClass(week)"
        >
          <span class="text-xs opacity-70 mb-0.5">Wk</span>
          <span class="text-base font-bold">{{ week }}</span>
          <span v-if="isWeekPlayed(week)" class="text-xs mt-0.5 opacity-70">
            {{ weekScore(week) }}
          </span>
          <span v-else class="text-xs mt-0.5 opacity-40">—</span>
        </button>
      </div>

      <!-- Selected week detail -->
      <div v-if="selectedWeekData" class="space-y-2">
        <!-- Week header -->
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-semibold text-gray-700"> Week {{ selectedWeek }} Results </span>
          <div class="flex gap-1">
            <button
              @click="prevWeek"
              :disabled="!canGoPrev"
              class="w-7 h-7 rounded-lg bg-gray-100 text-gray-500 hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed text-sm transition-colors"
            >
              ‹
            </button>
            <button
              @click="nextWeek"
              :disabled="!canGoNext"
              class="w-7 h-7 rounded-lg bg-gray-100 text-gray-500 hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed text-sm transition-colors"
            >
              ›
            </button>
          </div>
        </div>

        <!-- Mini match cards -->
        <div
          v-for="match in selectedWeekData.matches"
          :key="match.id"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg text-sm"
        >
          <!-- Home -->
          <div class="flex items-center gap-1.5 flex-1 justify-end">
            <span class="font-medium text-gray-700 text-xs">
              {{ shortName(match.home_team?.name) }}
            </span>
            <TeamLogo :name="match.home_team?.name" :size="22" />
          </div>

          <!-- Score -->
          <div class="flex items-center gap-1.5 mx-3">
            <span
              class="text-base font-bold w-5 text-center"
              :class="match.home_goals > match.away_goals ? 'text-success' : 'text-gray-700'"
            >
              {{ match.home_goals }}
            </span>
            <span class="text-gray-300 text-xs">—</span>
            <span
              class="text-base font-bold w-5 text-center"
              :class="match.away_goals > match.home_goals ? 'text-success' : 'text-gray-700'"
            >
              {{ match.away_goals }}
            </span>
          </div>

          <!-- Away -->
          <div class="flex items-center gap-1.5 flex-1">
            <TeamLogo :name="match.away_team?.name" :size="22" />
            <span class="font-medium text-gray-700 text-xs">
              {{ shortName(match.away_team?.name) }}
            </span>
          </div>
        </div>

        <!-- Mini standings snapshot -->
        <div class="mt-4">
          <p class="text-xs text-gray-400 font-medium mb-2 uppercase tracking-wide">
            Standings after week {{ selectedWeek }}
          </p>
          <div
            v-for="(s, i) in selectedWeekData.standings"
            :key="s.team_id"
            class="flex items-center gap-2 py-1.5 text-xs border-b border-gray-50 last:border-0"
          >
            <span class="w-4 text-center font-bold text-gray-400">{{ i + 1 }}</span>
            <TeamLogo :name="s.team?.name" :size="18" />
            <span class="flex-1 font-medium text-gray-700">
              {{ shortName(s.team?.name) }}
            </span>
            <span class="text-gray-400">{{ s.played }}MP</span>
            <span class="font-bold text-gray-800 w-6 text-right">{{ s.points }}</span>
          </div>
        </div>
      </div>

      <!-- No weeks played -->
      <div v-else class="py-6 text-center text-gray-400">
        <p class="text-2xl mb-1">🗓️</p>
        <p class="text-xs">No weeks played yet</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import TeamLogo from '../ui/TeamLogo.vue'

const props = defineProps({
  weekResults: { type: Array, default: () => [] },
})

const selectedWeek = ref(null)

const playedWeeks = computed(() => props.weekResults.map((w) => w.week))

import { watch } from 'vue'
watch(
  () => props.weekResults,
  (results) => {
    if (results.length) {
      selectedWeek.value = results[results.length - 1].week
    }
  },
  { immediate: true },
)

const selectedWeekData = computed(
  () => props.weekResults.find((w) => w.week === selectedWeek.value) || null,
)

const canGoPrev = computed(() => {
  const idx = playedWeeks.value.indexOf(selectedWeek.value)
  return idx > 0
})

const canGoNext = computed(() => {
  const idx = playedWeeks.value.indexOf(selectedWeek.value)
  return idx < playedWeeks.value.length - 1
})

function selectWeek(week) {
  if (isWeekPlayed(week)) selectedWeek.value = week
}

function prevWeek() {
  const idx = playedWeeks.value.indexOf(selectedWeek.value)
  if (idx > 0) selectedWeek.value = playedWeeks.value[idx - 1]
}

function nextWeek() {
  const idx = playedWeeks.value.indexOf(selectedWeek.value)
  if (idx < playedWeeks.value.length - 1) {
    selectedWeek.value = playedWeeks.value[idx + 1]
  }
}

function isWeekPlayed(week) {
  return playedWeeks.value.includes(week)
}

function weekScore(week) {
  const w = props.weekResults.find((r) => r.week === week)
  if (!w?.matches?.length) return ''
  const goals = w.matches.reduce((sum, m) => sum + (m.home_goals || 0) + (m.away_goals || 0), 0)
  return `${goals}G`
}

function weekPillClass(week) {
  if (!isWeekPlayed(week)) {
    return 'bg-gray-50 text-gray-300 cursor-not-allowed'
  }
  if (selectedWeek.value === week) {
    return 'bg-primary text-white shadow-sm scale-105'
  }
  return 'bg-gray-100 text-gray-600 hover:bg-gray-200 cursor-pointer'
}

const SHORT_NAMES = {
  'Manchester City': 'Man City',
  Liverpool: 'Liverpool',
  Arsenal: 'Arsenal',
  Chelsea: 'Chelsea',
}

function shortName(name) {
  return SHORT_NAMES[name] || name
}

const TEAM_EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

function teamEmoji(name) {
  return TEAM_EMOJIS[name] || '⚽'
}
</script>

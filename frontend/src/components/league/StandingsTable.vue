<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Header -->
    <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
      <h2 class="font-bold text-gray-800 text-lg flex items-center gap-2">🏆 League Table</h2>
      <span class="text-xs text-gray-400 font-medium"> Week {{ currentWeek }} / 6 </span>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <!-- Head -->
        <thead>
          <tr class="bg-gray-50 text-gray-500 text-xs uppercase tracking-wide">
            <th class="px-4 py-3 text-left w-8">#</th>
            <th class="px-4 py-3 text-left">Club</th>
            <th class="px-4 py-3 text-center">MP</th>
            <th class="px-4 py-3 text-center">W</th>
            <th class="px-4 py-3 text-center">D</th>
            <th class="px-4 py-3 text-center">L</th>
            <th class="px-4 py-3 text-center">GF</th>
            <th class="px-4 py-3 text-center">GA</th>
            <th class="px-4 py-3 text-center">GD</th>
            <th class="px-4 py-3 text-center font-bold text-gray-700">PTS</th>
          </tr>
        </thead>

        <!-- Body -->
        <tbody>
          <TransitionGroup name="standings">
            <tr
              v-for="(standing, index) in standings"
              :key="standing.team_id"
              class="border-t border-gray-50 hover:bg-gray-50 transition-colors duration-150"
              :class="rowClass(index)"
            >
              <!-- Position -->
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold"
                  :class="positionBadge(index)"
                >
                  {{ index + 1 }}
                </span>
              </td>

              <!-- Club -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-2.5">
                  <TeamLogo :name="standing.team?.name" :size="28" />
                  <span class="font-medium text-gray-800">{{ standing.team?.name }}</span>
                  <span
                    v-if="standing.team?.strength"
                    class="text-xs text-gray-400 bg-gray-100 px-1.5 py-0.5 rounded hidden sm:inline"
                  >
                    {{ standing.team.strength }}
                  </span>
                </div>
              </td>

              <!-- Stats -->
              <td class="px-4 py-3 text-center text-gray-600">{{ standing.played }}</td>
              <td class="px-4 py-3 text-center text-success font-medium">{{ standing.won }}</td>
              <td class="px-4 py-3 text-center text-warning">{{ standing.drawn }}</td>
              <td class="px-4 py-3 text-center text-danger">{{ standing.lost }}</td>
              <td class="px-4 py-3 text-center text-gray-600">{{ standing.goals_for }}</td>
              <td class="px-4 py-3 text-center text-gray-600">{{ standing.goals_against }}</td>
              <td class="px-4 py-3 text-center" :class="gdClass(standing.goal_diff)">
                {{ standing.goal_diff > 0 ? '+' : '' }}{{ standing.goal_diff }}
              </td>

              <!-- Points -->
              <td class="px-4 py-3 text-center">
                <span class="font-bold text-gray-900 text-base">{{ standing.points }}</span>
              </td>
            </tr>
          </TransitionGroup>
        </tbody>
      </table>
    </div>

    <!-- Empty state -->
    <div v-if="!standings?.length" class="py-12 text-center text-gray-400">
      <p class="text-4xl mb-2">📋</p>
      <p class="text-sm">No standings available yet.</p>
    </div>

    <!-- Legend -->
    <div class="px-5 py-3 border-t border-gray-100 flex gap-4 text-xs text-gray-400">
      <span class="flex items-center gap-1">
        <span class="w-2 h-2 rounded-full bg-blue-500 inline-block"></span>
        Champion
      </span>
      <span class="flex items-center gap-1">
        <span class="w-2 h-2 rounded-full bg-green-500 inline-block"></span>
        Top 2
      </span>
      <span>MP=Played · W=Won · D=Drawn · L=Lost · GD=Goal Diff · PTS=Points</span>
    </div>
  </div>
</template>

<script setup>
import TeamLogo from '../ui/TeamLogo.vue'

const props = defineProps({
  standings: { type: Array, default: () => [] },
  currentWeek: { type: Number, default: 0 },
})

const TEAM_EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

function teamEmoji(name) {
  return TEAM_EMOJIS[name] || '⚽'
}

function rowClass(index) {
  if (index === 0) return 'bg-blue-50/40'
  if (index === 1) return 'bg-green-50/30'
  return ''
}

function positionBadge(index) {
  if (index === 0) return 'bg-blue-100 text-blue-700'
  if (index === 1) return 'bg-green-100 text-green-700'
  if (index === 2) return 'bg-yellow-100 text-yellow-700'
  return 'bg-gray-100 text-gray-500'
}

function gdClass(gd) {
  if (gd > 0) return 'text-success font-medium'
  if (gd < 0) return 'text-danger font-medium'
  return 'text-gray-400'
}
</script>

<style scoped>
.standings-move {
  transition: transform 0.4s ease;
}
.standings-enter-active,
.standings-leave-active {
  transition: all 0.3s ease;
}
.standings-enter-from,
.standings-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}
</style>

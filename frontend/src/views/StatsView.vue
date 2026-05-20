<template>
  <div class="space-y-6">
    <!-- Title -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">📊 Season Statistics</h1>
      <p class="text-sm text-gray-500 mt-1">Detailed breakdown of the current season</p>
    </div>

    <!-- Not started -->
    <div v-if="!hasData" class="py-20 text-center text-gray-400">
      <p class="text-5xl mb-4">📈</p>
      <p class="text-lg font-medium text-gray-500">No data yet</p>
      <p class="text-sm mt-1">Play some matches to see statistics.</p>
      <RouterLink to="/" class="inline-block mt-4 text-sm text-primary hover:underline">
        ← Go to League
      </RouterLink>
    </div>

    <template v-else>
      <!-- Summary cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard icon="⚽" label="Total Goals" :value="totalGoals" color="blue" />
        <StatCard
          icon="🎮"
          label="Matches Played"
          :value="store.status?.matches_played || 0"
          color="purple"
        />
        <StatCard
          icon="🏆"
          label="Leader"
          :value="leader?.team?.name || '—'"
          color="yellow"
          small
        />
        <StatCard
          icon="📅"
          label="Current Week"
          :value="`${store.currentWeek} / 6`"
          color="green"
          small
        />
      </div>

      <!-- Charts row -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Goals per team bar chart -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <div class="px-5 py-4 border-b border-gray-100">
            <h2 class="font-bold text-gray-800">⚽ Goals Scored vs Conceded</h2>
            <p class="text-xs text-gray-400 mt-0.5">Attack vs defense comparison</p>
          </div>
          <div class="p-5">
            <canvas ref="goalsChartCanvas" height="220"></canvas>
          </div>
        </div>

        <!-- Points progression line chart -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <div class="px-5 py-4 border-b border-gray-100">
            <h2 class="font-bold text-gray-800">📈 Points Progression</h2>
            <p class="text-xs text-gray-400 mt-0.5">Points accumulated per week</p>
          </div>
          <div class="p-5">
            <canvas ref="pointsChartCanvas" height="220"></canvas>
          </div>
        </div>
      </div>

      <!-- Team stats table -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100">
          <h2 class="font-bold text-gray-800">🔍 Team Breakdown</h2>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="bg-gray-50 text-xs uppercase text-gray-500 tracking-wide">
                <th class="px-5 py-3 text-left">Club</th>
                <th class="px-5 py-3 text-center">Form</th>
                <th class="px-5 py-3 text-center">Avg Goals/Game</th>
                <th class="px-5 py-3 text-center">Clean Sheets</th>
                <th class="px-5 py-3 text-center">Best Win</th>
                <th class="px-5 py-3 text-center">Win Rate</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="stat in teamStats"
                :key="stat.teamId"
                class="border-t border-gray-50 hover:bg-gray-50 transition-colors"
              >
                <td class="px-5 py-3">
                  <div class="flex items-center gap-2">
                    <span class="text-lg">{{ teamEmoji(stat.name) }}</span>
                    <span class="font-medium text-gray-800">{{ stat.name }}</span>
                  </div>
                </td>
                <td class="px-5 py-3">
                  <div class="flex items-center justify-center gap-1">
                    <span
                      v-for="(result, i) in stat.form"
                      :key="i"
                      class="w-5 h-5 rounded-full text-xs flex items-center justify-center font-bold text-white"
                      :class="formColor(result)"
                    >
                      {{ result }}
                    </span>
                    <span v-if="!stat.form.length" class="text-gray-300 text-xs">—</span>
                  </div>
                </td>
                <td class="px-5 py-3 text-center text-gray-600">
                  {{ stat.avgGoals }}
                </td>
                <td class="px-5 py-3 text-center">
                  <span class="font-medium text-green-600">{{ stat.cleanSheets }}</span>
                </td>
                <td class="px-5 py-3 text-center text-gray-600 text-xs">
                  {{ stat.bestWin || '—' }}
                </td>
                <td class="px-5 py-3 text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <div class="w-16 bg-gray-100 rounded-full h-1.5">
                      <div
                        class="h-1.5 rounded-full bg-primary"
                        :style="{ width: stat.winRate + '%' }"
                      ></div>
                    </div>
                    <span class="text-xs text-gray-600">{{ stat.winRate }}%</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Chart,
  BarElement,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend,
  BarController,
  LineController,
} from 'chart.js'
import { useLeagueStore } from '../stores/league'
import { useLeague } from '../composables/useLeague'
import StatCard from '../components/ui/StatCard.vue'

Chart.register(
  BarElement,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend,
  BarController,
  LineController,
)

const store = useLeagueStore()
const league = useLeague()

const goalsChartCanvas = ref(null)
const pointsChartCanvas = ref(null)
let goalsChart = null
let pointsChart = null

const TEAM_EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

const TEAM_COLORS = {
  'Manchester City': '#1a56db',
  Liverpool: '#e02424',
  Arsenal: '#e3a008',
  Chelsea: '#7e3af2',
}

function teamEmoji(name) {
  return TEAM_EMOJIS[name] || '⚽'
}
function teamColor(name) {
  return TEAM_COLORS[name] || '#6b7280'
}

const hasData = computed(() => store.standings.length > 0 && store.currentWeek > 0)

const leader = computed(() => (store.standings.length ? store.standings[0] : null))

const totalGoals = computed(() => store.standings.reduce((sum, s) => sum + (s.goals_for || 0), 0))

const teamStats = computed(() => {
  return store.standings.map((s) => {
    const played = s.played || 0
    const avgGoals = played > 0 ? (s.goals_for / played).toFixed(1) : '0.0'
    const winRate = played > 0 ? Math.round((s.won / played) * 100) : 0

    const form = []
    const sortedWeeks = [...store.weekResults].sort((a, b) => a.week - b.week)
    sortedWeeks.forEach((week) => {
      week.matches?.forEach((m) => {
        if (m.home_team_id === s.team_id || m.away_team_id === s.team_id) {
          const isHome = m.home_team_id === s.team_id
          const gf = isHome ? m.home_goals : m.away_goals
          const ga = isHome ? m.away_goals : m.home_goals
          if (gf > ga) form.push('W')
          else if (gf < ga) form.push('L')
          else form.push('D')
        }
      })
    })

    let cleanSheets = 0
    sortedWeeks.forEach((week) => {
      week.matches?.forEach((m) => {
        const isHome = m.home_team_id === s.team_id
        const ga = isHome ? m.away_goals : m.home_goals
        if ((m.home_team_id === s.team_id || m.away_team_id === s.team_id) && ga === 0) {
          cleanSheets++
        }
      })
    })

    let bestWin = null
    let bestDiff = 0
    sortedWeeks.forEach((week) => {
      week.matches?.forEach((m) => {
        const isHome = m.home_team_id === s.team_id
        if (m.home_team_id !== s.team_id && m.away_team_id !== s.team_id) return
        const gf = isHome ? m.home_goals : m.away_goals
        const ga = isHome ? m.away_goals : m.home_goals
        const diff = gf - ga
        if (diff > bestDiff) {
          bestDiff = diff
          const opp = isHome ? m.away_team?.name : m.home_team?.name
          bestWin = `${gf}-${ga} vs ${opp}`
        }
      })
    })

    return {
      teamId: s.team_id,
      name: s.team?.name || '',
      form: form.slice(-5),
      avgGoals,
      cleanSheets,
      bestWin,
      winRate,
    }
  })
})

function formColor(result) {
  if (result === 'W') return 'bg-green-500'
  if (result === 'L') return 'bg-red-500'
  return 'bg-gray-400'
}

// Goals bar chart
function buildGoalsChart() {
  if (!goalsChartCanvas.value || !store.standings.length) return
  if (goalsChart) {
    goalsChart.destroy()
    goalsChart = null
  }

  const labels = store.standings.map((s) => s.team?.name)
  const colors = store.standings.map((s) => teamColor(s.team?.name))

  goalsChart = new Chart(goalsChartCanvas.value, {
    type: 'bar',
    data: {
      labels,
      datasets: [
        {
          label: 'Goals Scored',
          data: store.standings.map((s) => s.goals_for),
          backgroundColor: colors.map((c) => c + 'CC'),
          borderColor: colors,
          borderWidth: 2,
          borderRadius: 6,
        },
        {
          label: 'Goals Conceded',
          data: store.standings.map((s) => s.goals_against),
          backgroundColor: colors.map(() => '#f3f4f6'),
          borderColor: colors.map(() => '#d1d5db'),
          borderWidth: 2,
          borderRadius: 6,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'bottom',
          labels: { font: { size: 11 }, boxWidth: 12 },
        },
      },
      scales: {
        x: { grid: { display: false } },
        y: { beginAtZero: true, grid: { color: '#f3f4f6' } },
      },
    },
  })
}

function buildPointsChart() {
  if (!pointsChartCanvas.value || !store.weekResults.length) return
  if (pointsChart) {
    pointsChart.destroy()
    pointsChart = null
  }

  const sortedWeeks = [...store.weekResults].sort((a, b) => a.week - b.week)
  const labels = sortedWeeks.map((w) => `Week ${w.week}`)

  const teams = store.standings.map((s) => s.team?.name).filter(Boolean)
  const datasets = teams.map((teamName) => {
    let cumulative = 0
    const data = sortedWeeks.map((week) => {
      const standing = week.standings?.find((s) => s.team?.name === teamName)
      return standing?.points || 0
    })
    return {
      label: teamName,
      data,
      borderColor: teamColor(teamName),
      backgroundColor: teamColor(teamName) + '20',
      borderWidth: 2.5,
      pointRadius: 4,
      pointHoverRadius: 6,
      tension: 0.3,
      fill: false,
    }
  })

  pointsChart = new Chart(pointsChartCanvas.value, {
    type: 'line',
    data: { labels, datasets },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'bottom',
          labels: { font: { size: 11 }, boxWidth: 12 },
        },
      },
      scales: {
        x: { grid: { display: false } },
        y: {
          beginAtZero: true,
          grid: { color: '#f3f4f6' },
          ticks: { stepSize: 3 },
        },
      },
    },
  })
}

watch(
  () => [store.standings, store.weekResults],
  async () => {
    await nextTick()
    buildGoalsChart()
    buildPointsChart()
  },
  { deep: true },
)

onMounted(async () => {
  if (!store.standings.length) {
    await league.fetchAll()
  }
  await nextTick()
  buildGoalsChart()
  buildPointsChart()
})

onUnmounted(() => {
  if (goalsChart) goalsChart.destroy()
  if (pointsChart) pointsChart.destroy()
})
</script>

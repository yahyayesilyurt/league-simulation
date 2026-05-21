<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Header -->
    <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
      <div>
        <h2 class="font-bold text-gray-800 text-lg">🎯 Championship Predictions</h2>
        <p class="text-xs text-gray-400 mt-0.5">
          {{
            currentWeek < 4
              ? 'Real predictions start from week 4'
              : 'Based on current standings & strength'
          }}
        </p>
      </div>
      <span
        v-if="currentWeek < 4"
        class="text-xs bg-yellow-100 text-yellow-700 px-2.5 py-1 rounded-full font-medium"
      >
        Week {{ currentWeek }}/4
      </span>
      <span
        v-else-if="currentWeek >= 6"
        class="text-xs bg-green-100 text-green-700 px-2.5 py-1 rounded-full font-medium"
      >
        🏆 Final
      </span>
    </div>

    <div class="p-5">
      <!-- Donut Chart -->
      <div class="flex justify-center mb-6">
        <div class="relative w-44 h-44">
          <canvas ref="chartCanvas"></canvas>

          <!-- Center label -->
          <div class="absolute inset-0 flex flex-col items-center justify-center">
            <span class="text-2xl">🏆</span>
            <span class="text-xs text-gray-400 mt-1">Champion</span>
          </div>
        </div>
      </div>

      <!-- Progress bars -->
      <div class="space-y-3">
        <div
          v-for="(pred, index) in sortedPredictions"
          :key="pred.team_id || pred.team_name"
          class="space-y-1"
        >
          <div class="flex items-center justify-between text-sm">
            <div class="flex items-center gap-2">
              <TeamLogo :name="pred.team_name" :size="22" />
              <span class="font-medium text-gray-700">{{ pred.team_name }}</span>
            </div>
            <div class="flex items-center gap-3 text-xs text-gray-500">
              <span>{{ pred.current_points }} pts</span>
              <span class="font-bold text-gray-800">{{ pred.percentage }}%</span>
            </div>
          </div>

          <!-- Progress bar -->
          <div class="w-full bg-gray-100 rounded-full h-2">
            <div
              class="h-2 rounded-full transition-all duration-700 ease-out"
              :style="{
                width: pred.percentage + '%',
                backgroundColor: chartColors[index],
              }"
            ></div>
          </div>
        </div>
      </div>

      <!-- Week 4 info banner -->
      <div
        v-if="currentWeek < 4 && currentWeek > 0"
        class="mt-4 p-3 bg-blue-50 border border-blue-100 rounded-lg text-center"
      >
        <p class="text-xs text-blue-600">📊 Detailed predictions unlock after Week 4</p>
      </div>

      <!-- Not started -->
      <div v-if="currentWeek === 0" class="mt-4 py-6 text-center text-gray-400">
        <p class="text-3xl mb-2">🎲</p>
        <p class="text-sm">Play matches to see predictions</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { Chart, ArcElement, Tooltip, Legend, DoughnutController } from 'chart.js'
import TeamLogo from '../ui/TeamLogo.vue'

Chart.register(ArcElement, Tooltip, Legend, DoughnutController)

const props = defineProps({
  predictions: { type: Array, default: () => [] },
  currentWeek: { type: Number, default: 0 },
})

const chartCanvas = ref(null)
let chartInstance = null

const chartColors = [
  '#1a56db', // blue   — Man City
  '#e02424', // red    — Liverpool
  '#e3a008', // yellow — Arsenal
  '#7e3af2', // purple — Chelsea
]

const TEAM_EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

function teamEmoji(name) {
  return TEAM_EMOJIS[name] || '⚽'
}

const sortedPredictions = computed(() =>
  [...props.predictions].sort((a, b) => b.percentage - a.percentage),
)

function buildChart() {
  if (!chartCanvas.value) return
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }

  const labels = sortedPredictions.value.map((p) => p.team_name)
  const data = sortedPredictions.value.map((p) => p.percentage || 25)
  const colors = sortedPredictions.value.map((_, i) => chartColors[i])

  chartInstance = new Chart(chartCanvas.value, {
    type: 'doughnut',
    data: {
      labels,
      datasets: [
        {
          data,
          backgroundColor: colors,
          borderWidth: 2,
          borderColor: '#ffffff',
          hoverOffset: 6,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      cutout: '65%',
      plugins: {
        legend: { display: false },
        tooltip: {
          callbacks: {
            label: (ctx) => ` ${ctx.label}: ${ctx.raw}%`,
          },
        },
      },
      animation: {
        duration: 600,
        easing: 'easeInOutQuart',
      },
    },
  })
}

watch(
  () => props.predictions,
  async () => {
    await nextTick()
    buildChart()
  },
  { deep: true },
)

onMounted(async () => {
  await nextTick()
  buildChart()
})

onUnmounted(() => {
  if (chartInstance) chartInstance.destroy()
})
</script>

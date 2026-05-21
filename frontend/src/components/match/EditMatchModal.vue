<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="$emit('close')" />

        <!-- Modal -->
        <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6 z-10">
          <!-- Title -->
          <h3 class="font-bold text-gray-800 text-lg mb-1">Edit Match Result</h3>
          <p class="text-xs text-gray-400 mb-5">Week {{ match?.week }}</p>

          <!-- Teams -->
          <div class="flex items-center justify-between gap-4 mb-6">
            <!-- Home -->
            <div class="flex-1 text-center">
              <TeamLogo :name="match?.home_team?.name" :size="48" class="mx-auto mb-1" />
              <div class="text-xs font-medium text-gray-700">{{ match?.home_team?.name }}</div>
              <input
                v-model.number="homeGoals"
                type="number"
                min="0"
                max="20"
                class="mt-2 w-16 text-center text-xl font-bold border-2 border-gray-200 rounded-lg py-1 focus:border-primary focus:outline-none mx-auto block"
              />
            </div>

            <span class="text-gray-300 text-2xl font-light">—</span>

            <!-- Away -->
            <div class="flex-1 text-center">
              <TeamLogo :name="match?.away_team?.name" :size="48" class="mx-auto mb-1" />
              <div class="text-xs font-medium text-gray-700">{{ match?.away_team?.name }}</div>
              <input
                v-model.number="awayGoals"
                type="number"
                min="0"
                max="20"
                class="mt-2 w-16 text-center text-xl font-bold border-2 border-gray-200 rounded-lg py-1 focus:border-primary focus:outline-none mx-auto block"
              />
            </div>
          </div>

          <!-- Error -->
          <p v-if="error" class="text-xs text-danger mb-3 text-center">⚠️ {{ error }}</p>

          <!-- Buttons -->
          <div class="flex gap-2">
            <button
              @click="$emit('close')"
              class="flex-1 py-2.5 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Cancel
            </button>
            <button
              @click="handleSave"
              :disabled="loading"
              class="flex-1 py-2.5 text-sm font-medium text-white bg-primary rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              <span v-if="loading">Saving...</span>
              <span v-else>Save Result</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import TeamLogo from '../ui/TeamLogo.vue'

const props = defineProps({
  show: { type: Boolean, default: false },
  match: { type: Object, default: null },
  loading: { type: Boolean, default: false },
  error: { type: String, default: null },
})

const emit = defineEmits(['close', 'save'])

const homeGoals = ref(0)
const awayGoals = ref(0)

watch(
  () => props.match,
  (match) => {
    if (match) {
      homeGoals.value = match.home_goals ?? 0
      awayGoals.value = match.away_goals ?? 0
    }
  },
  { immediate: true },
)

const TEAM_EMOJIS = {
  'Manchester City': '🔵',
  Liverpool: '🔴',
  Arsenal: '❤️',
  Chelsea: '💙',
}

function teamEmoji(name) {
  return TEAM_EMOJIS[name] || '⚽'
}

function handleSave() {
  emit('save', {
    matchId: props.match.id,
    homeGoals: homeGoals.value,
    awayGoals: awayGoals.value,
  })
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>

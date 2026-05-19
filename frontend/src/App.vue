<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">🏆 League Table</h1>

    <div v-if="store.loading" class="text-gray-500">Loading...</div>
    <div v-else-if="store.error" class="text-red-500">{{ store.error }}</div>
    <div v-else>
      <div
        v-for="s in store.standings"
        :key="s.team_id"
        class="flex justify-between p-3 bg-white rounded-lg shadow mb-2"
      >
        <span class="font-medium">{{ s.team.name }}</span>
        <span class="text-primary font-bold">{{ s.points }} pts</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useLeagueStore } from './stores/league'
import { useLeague } from './composables/useLeague'

const store = useLeagueStore()
const league = useLeague()

onMounted(() => {
  league.fetchTable()
})
</script>

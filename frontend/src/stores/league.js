import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLeagueStore = defineStore('league', () => {
  const standings = ref([])
  const fixtures = ref([])
  const predictions = ref([])
  const status = ref(null)
  const currentWeek = ref(0)
  const weekResults = ref([])
  const loading = ref(false)
  const error = ref(null)

  function setStandings(data) {
    standings.value = data
  }
  function setFixtures(data) {
    fixtures.value = data
  }
  function setPredictions(data) {
    predictions.value = data
  }
  function setStatus(data) {
    status.value = data
  }
  function setLoading(val) {
    loading.value = val
  }
  function setError(val) {
    error.value = val
  }

  function addWeekResult(result) {
    const index = weekResults.value.findIndex((w) => w.week === result.week)
    if (index !== -1) {
      weekResults.value[index] = result
    } else {
      weekResults.value.push(result)
    }
    currentWeek.value = result.week
  }

  function updateMatchInStore(matchId, homeGoals, awayGoals) {
    weekResults.value = weekResults.value.map((weekResult) => ({
      ...weekResult,
      matches: weekResult.matches.map((match) => {
        if (match.id !== matchId) return match
        return { ...match, home_goals: homeGoals, away_goals: awayGoals }
      }),
    }))
  }

  function reset() {
    standings.value = []
    fixtures.value = []
    predictions.value = []
    weekResults.value = []
    currentWeek.value = 0
    status.value = null
    error.value = null
  }

  return {
    standings,
    fixtures,
    predictions,
    status,
    currentWeek,
    weekResults,
    loading,
    error,
    setStandings,
    setFixtures,
    setPredictions,
    setStatus,
    setLoading,
    setError,
    addWeekResult,
    updateMatchInStore,
    reset,
  }
})

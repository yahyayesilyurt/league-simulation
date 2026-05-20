import { reactive } from 'vue'
import { matchApi } from '../api/match'
import { leagueApi } from '../api/league'
import { useLeagueStore } from '../stores/league'

export function useMatch() {
  const store = useLeagueStore()

  const editModal = reactive({
    show: false,
    match: null,
    loading: false,
    error: null,
  })

  function openEditModal(match) {
    editModal.match = match
    editModal.error = null
    editModal.show = true
  }

  function closeEditModal() {
    editModal.show = false
    editModal.match = null
    editModal.error = null
  }

  async function saveMatchResult({ matchId, homeGoals, awayGoals }) {
    editModal.loading = true
    editModal.error = null

    try {
      const res = await matchApi.updateResult(matchId, homeGoals, awayGoals)

      store.updateMatchInStore(matchId, homeGoals, awayGoals)

      if (res.data.standings) {
        store.setStandings(res.data.standings)
      }

      const predRes = await leagueApi.getPredictions()
      store.setPredictions(predRes.data.predictions)

      closeEditModal()
      return true
    } catch (err) {
      editModal.error = err.response?.data?.error || 'Failed to update match result'
      return false
    } finally {
      editModal.loading = false
    }
  }

  return {
    editModal,
    openEditModal,
    closeEditModal,
    saveMatchResult,
  }
}

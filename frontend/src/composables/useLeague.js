import { useLeagueStore } from '../stores/league'
import { leagueApi } from '../api/league'

export function useLeague() {
  const store = useLeagueStore()

  async function fetchTable() {
    try {
      store.setLoading(true)
      const res = await leagueApi.getTable()
      store.setStandings(res.data.standings)
    } catch (err) {
      store.setError(err.response?.data?.error || 'Scoreboard not available.')
    } finally {
      store.setLoading(false)
    }
  }

  async function fetchStatus() {
    try {
      const res = await leagueApi.getStatus()
      store.setStatus(res.data)
    } catch (err) {
      store.setError(err.response?.data?.error || 'Status not available.')
    }
  }

  async function fetchPredictions() {
    try {
      const res = await leagueApi.getPredictions()
      store.setPredictions(res.data.predictions)
    } catch (err) {
      store.setError(err.response?.data?.error || 'Predictions not available.')
    }
  }

  async function fetchFixtures() {
    try {
      const res = await leagueApi.getFixtures()
      store.setFixtures(res.data.fixtures)
    } catch (err) {
      store.setError(err.response?.data?.error || 'Fixtures not available.')
    }
  }

  async function playNextWeek() {
    try {
      store.setLoading(true)
      store.setError(null)
      const res = await leagueApi.nextWeek()
      store.addWeekResult(res.data)
      store.setStandings(res.data.standings)
      store.setPredictions(res.data.predictions)
      await fetchStatus()
      return res.data
    } catch (err) {
      store.setError(err.response?.data?.error || 'The week could not be played.')
      return null
    } finally {
      store.setLoading(false)
    }
  }

  async function playAll() {
    try {
      store.setLoading(true)
      store.setError(null)
      const res = await leagueApi.playAll()

      res.data.weeks.forEach((week) => store.addWeekResult(week))

      const lastWeek = res.data.weeks[res.data.weeks.length - 1]
      if (lastWeek) {
        store.setStandings(lastWeek.standings)
        store.setPredictions(lastWeek.predictions)
      }

      await fetchStatus()
      return res.data
    } catch (err) {
      store.setError(err.response?.data?.error || 'The league could not be played.')
      return null
    } finally {
      store.setLoading(false)
    }
  }

  async function resetLeague() {
    try {
      store.setLoading(true)
      store.setError(null)
      await leagueApi.reset()
      store.reset()
      await Promise.all([fetchTable(), fetchStatus()])
    } catch (err) {
      store.setError(err.response?.data?.error || 'The league could not be reset.')
    } finally {
      store.setLoading(false)
    }
  }

  return {
    fetchTable,
    fetchStatus,
    fetchPredictions,
    fetchFixtures,
    playNextWeek,
    playAll,
    resetLeague,
  }
}

import api from './index'

export const leagueApi = {
  getTable() {
    return api.get('/league/table')
  },

  getFixtures() {
    return api.get('/league/fixtures')
  },

  getWeek(weekNo) {
    return api.get(`/league/week/${weekNo}`)
  },

  getStatus() {
    return api.get('/league/status')
  },

  getPredictions() {
    return api.get('/league/predictions')
  },

  nextWeek() {
    return api.post('/league/next-week')
  },

  playAll() {
    return api.post('/league/play-all')
  },

  reset() {
    return api.post('/league/reset')
  },
}

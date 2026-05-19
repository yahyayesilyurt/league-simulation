import api from './index'

export const matchApi = {
  getMatch(id) {
    return api.get(`/match/${id}`)
  },

  updateResult(id, homeGoals, awayGoals) {
    return api.put(`/match/${id}/result`, { home_goals: homeGoals, away_goals: awayGoals })
  },
}

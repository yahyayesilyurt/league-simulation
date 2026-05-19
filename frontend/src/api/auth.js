import api from './index'

export const authApi = {
  login(username, password) {
    return api.post('/auth/login', { username, password })
  },
}

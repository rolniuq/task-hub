import api from './api'

export const authService = {
  async login(email, password) {
    const response = await api.post('/api/auth/login', { email, password })
    const { tokens, user } = response.data
    
    localStorage.setItem('access_token', tokens.access_token)
    localStorage.setItem('refresh_token', tokens.refresh_token)
    localStorage.setItem('user', JSON.stringify(user))
    
    return { user, tokens }
  },

  async register(name, email, password) {
    const response = await api.post('/api/auth/register', { name, email, password })
    return response.data
  },

  logout() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  },

  isAuthenticated() {
    return !!localStorage.getItem('access_token')
  },

  getUser() {
    const user = localStorage.getItem('user')
    return user ? JSON.parse(user) : null
  },

  getToken() {
    return localStorage.getItem('access_token')
  }
}

export const taskService = {
  async list() {
    const response = await api.get('/api/tasks')
    return response.data
  },

  async get(id) {
    const response = await api.get(`/api/tasks/${id}`)
    return response.data
  },

  async create(task) {
    const response = await api.post('/api/tasks', task)
    return response.data
  },

  async update(id, task) {
    const response = await api.put(`/api/tasks/${id}`, task)
    return response.data
  },

  async delete(id) {
    await api.delete(`/api/tasks/${id}`)
  },

  async complete(id) {
    const response = await api.post(`/api/tasks/${id}/complete`)
    return response.data
  }
}

import { defineStore } from 'pinia'
import { authService } from '../services/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: authService.getUser(),
    isAuthenticated: authService.isAuthenticated()
  }),

  actions: {
    async login(email, password) {
      const { user } = await authService.login(email, password)
      this.user = user
      this.isAuthenticated = true
      return user
    },

    async register(name, email, password) {
      return await authService.register(name, email, password)
    },

    logout() {
      authService.logout()
      this.user = null
      this.isAuthenticated = false
    },

    checkAuth() {
      this.isAuthenticated = authService.isAuthenticated()
      this.user = authService.getUser()
    }
  }
})

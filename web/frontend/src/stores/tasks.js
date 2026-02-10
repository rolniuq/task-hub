import { defineStore } from 'pinia'
import { taskService } from '../services/auth'

export const useTaskStore = defineStore('tasks', {
  state: () => ({
    tasks: [],
    loading: false,
    error: null
  }),

  getters: {
    pendingTasks: (state) => state.tasks.filter(t => t.status === 'pending'),
    completedTasks: (state) => state.tasks.filter(t => t.status === 'completed')
  },

  actions: {
    async fetchTasks() {
      this.loading = true
      this.error = null
      try {
        const data = await taskService.list()
        this.tasks = data.tasks || []
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async createTask(task) {
      const newTask = await taskService.create(task)
      this.tasks.push(newTask.task)
      return newTask
    },

    async updateTask(id, task) {
      const updated = await taskService.update(id, task)
      const index = this.tasks.findIndex(t => t.id === id)
      if (index !== -1) {
        this.tasks[index] = updated.task
      }
      return updated
    },

    async deleteTask(id) {
      await taskService.delete(id)
      this.tasks = this.tasks.filter(t => t.id !== id)
    },

    async completeTask(id) {
      const completed = await taskService.complete(id)
      const index = this.tasks.findIndex(t => t.id === id)
      if (index !== -1) {
        this.tasks[index] = completed.task
      }
      return completed
    }
  }
})

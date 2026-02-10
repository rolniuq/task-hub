<template>
  <div class="dashboard">
    <header class="header">
      <h1>TaskHub</h1>
      <div class="user-info">
        <span>{{ user?.name || user?.email }}</span>
        <button @click="handleLogout" class="btn btn-logout">Logout</button>
      </div>
    </header>

    <main class="main-content">
      <div class="task-header">
        <h2>My Tasks</h2>
        <button @click="showCreateModal = true" class="btn btn-primary">+ New Task</button>
      </div>

      <div v-if="taskStore.loading" class="loading">Loading tasks...</div>
      <div v-else-if="taskStore.error" class="alert alert-error">{{ taskStore.error }}</div>
      <div v-else-if="taskStore.tasks.length === 0" class="empty-state">
        <p>No tasks yet. Create your first task!</p>
      </div>
      <div v-else class="task-list">
        <div
          v-for="task in taskStore.tasks"
          :key="task.id"
          class="task-card"
          :class="{ completed: task.status === 'completed' }"
        >
          <div class="task-content">
            <h3>{{ task.title }}</h3>
            <p>{{ task.description }}</p>
            <span class="task-status" :class="task.status">{{ task.status }}</span>
          </div>
          <div class="task-actions">
            <button
              v-if="task.status !== 'completed'"
              @click="completeTask(task.id)"
              class="btn btn-complete"
            >
              Complete
            </button>
            <button @click="deleteTask(task.id)" class="btn btn-delete">Delete</button>
          </div>
        </div>
      </div>
    </main>

    <!-- Create Task Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h2>Create New Task</h2>
        <form @submit.prevent="createTask">
          <div class="form-group">
            <label for="title">Title</label>
            <input
              type="text"
              id="title"
              v-model="newTask.title"
              placeholder="Task title"
              required
            />
          </div>
          <div class="form-group">
            <label for="description">Description</label>
            <textarea
              id="description"
              v-model="newTask.description"
              placeholder="Task description"
              rows="3"
            ></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary">Create</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useTaskStore } from '../stores/tasks'

const router = useRouter()
const authStore = useAuthStore()
const taskStore = useTaskStore()

const user = computed(() => authStore.user)
const showCreateModal = ref(false)
const newTask = ref({
  title: '',
  description: ''
})

onMounted(() => {
  taskStore.fetchTasks()
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const createTask = async () => {
  try {
    await taskStore.createTask(newTask.value)
    newTask.value = { title: '', description: '' }
    showCreateModal.value = false
  } catch (err) {
    console.error('Failed to create task:', err)
  }
}

const completeTask = async (id) => {
  try {
    await taskStore.completeTask(id)
  } catch (err) {
    console.error('Failed to complete task:', err)
  }
}

const deleteTask = async (id) => {
  if (confirm('Are you sure you want to delete this task?')) {
    try {
      await taskStore.deleteTask(id)
    } catch (err) {
      console.error('Failed to delete task:', err)
    }
  }
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: #f5f7fa;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.btn-logout {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
}

.btn-logout:hover {
  background: rgba(255, 255, 255, 0.3);
}

.main-content {
  max-width: 800px;
  margin: 2rem auto;
  padding: 0 1rem;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.task-header h2 {
  margin: 0;
  color: #333;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: #e5e7eb;
  color: #374151;
}

.btn-secondary:hover {
  background: #d1d5db;
}

.btn-complete {
  background: #10b981;
  color: white;
}

.btn-complete:hover {
  background: #059669;
}

.btn-delete {
  background: #ef4444;
  color: white;
}

.btn-delete:hover {
  background: #dc2626;
}

.loading,
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.task-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.task-card.completed {
  opacity: 0.7;
}

.task-card.completed h3 {
  text-decoration: line-through;
}

.task-content h3 {
  margin: 0 0 0.5rem 0;
  color: #333;
}

.task-content p {
  margin: 0 0 0.5rem 0;
  color: #666;
}

.task-status {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
}

.task-status.pending {
  background: #fef3c7;
  color: #92400e;
}

.task-status.completed {
  background: #d1fae5;
  color: #065f46;
}

.task-actions {
  display: flex;
  gap: 0.5rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
}

.modal h2 {
  margin: 0 0 1.5rem 0;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-weight: 500;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.alert {
  padding: 0.75rem;
  border-radius: 6px;
  margin-bottom: 1rem;
}

.alert-error {
  background-color: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
}
</style>

<template>
  <div class="users-view">
    <div class="header">
      <h1>Users</h1>
      <button @click="showCreateModal = true" class="btn btn-primary">Create User</button>
    </div>

    <div class="filters">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search users..."
        @input="debouncedSearch"
        class="search-input"
      />
      <select v-model="pageSize" @change="loadUsers" class="page-size-select">
        <option :value="10">10 per page</option>
        <option :value="20">20 per page</option>
        <option :value="50">50 per page</option>
      </select>
    </div>

    <div v-if="loading" class="loading">Loading users...</div>

    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else-if="users.length === 0" class="empty">
      <p>No users found. Create one to get started!</p>
    </div>

    <table v-else class="users-table">
      <thead>
        <tr>
          <th>Username</th>
          <th>Email</th>
          <th>Status</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.uid">
          <td>{{ user.username }}</td>
          <td>{{ user.email }}</td>
          <td>
            <span :class="['status', user.status === 1 ? 'active' : 'inactive']">
              {{ user.status === 1 ? 'Active' : 'Inactive' }}
            </span>
          </td>
          <td>{{ formatDate(user.created_at) }}</td>
          <td>
            <router-link :to="`/users/${user.uid}`" class="action-btn">View</router-link>
            <button @click="editUser(user)" class="action-btn">Edit</button>
            <button @click="confirmDelete(user)" class="action-btn danger">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="pagination.pages > 1" class="pagination">
      <button
        :disabled="pagination.page === 1"
        @click="goToPage(pagination.page - 1)"
        class="pagination-btn"
      >
        Previous
      </button>
      <span class="pagination-info">
        Page {{ pagination.page }} of {{ pagination.pages }}
      </span>
      <button
        :disabled="pagination.page === pagination.pages"
        @click="goToPage(pagination.page + 1)"
        class="pagination-btn"
      >
        Next
      </button>
    </div>

    <!-- Create/Edit User Modal -->
    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal">
        <h2>{{ showEditModal ? 'Edit User' : 'Create User' }}</h2>
        <form @submit.prevent="submitUser">
          <div class="form-group">
            <label for="username">Username</label>
            <input
              id="username"
              v-model="userForm.username"
              type="text"
              required
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label for="email">Email</label>
            <input
              id="email"
              v-model="userForm.email"
              type="email"
              required
              class="form-input"
            />
          </div>
          <div v-if="!showEditModal" class="form-group">
            <label for="password">Password</label>
            <input
              id="password"
              v-model="userForm.password"
              type="password"
              required
              class="form-input"
            />
          </div>
          <div v-if="showEditModal" class="form-group">
            <label for="status">Status</label>
            <select id="status" v-model="userForm.status" class="form-input">
              <option :value="1">Active</option>
              <option :value="0">Inactive</option>
            </select>
          </div>
          <div class="form-actions">
            <button type="button" @click="closeModals" class="btn btn-secondary">Cancel</button>
            <button type="submit" class="btn btn-primary">
              {{ showEditModal ? 'Update' : 'Create' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal">
        <h2>Confirm Delete</h2>
        <p>Are you sure you want to delete user "{{ userToDelete?.username }}"?</p>
        <div class="form-actions">
          <button @click="showDeleteModal = false" class="btn btn-secondary">Cancel</button>
          <button @click="deleteUser" class="btn btn-danger">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { usersApi, type User, type Meta } from '@/services/api'

const users = ref<User[]>([])
const pagination = ref<Meta>({ page: 1, limit: 20, total: 0, pages: 0 })
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const pageSize = ref(20)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const userToDelete = ref<User | null>(null)
const editingUser = ref<User | null>(null)

const userForm = ref({
  username: '',
  email: '',
  password: '',
  status: 1
})

let searchTimeout: number | null = null

const loadUsers = async () => {
  loading.value = true
  error.value = null

  try {
    const params: any = {
      page: pagination.value.page,
      limit: pageSize.value
    }

    if (searchQuery.value) {
      params.query = searchQuery.value
    }

    const response = await usersApi.listUsers(params)
    users.value = response.users
    pagination.value = response.meta
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to load users'
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = window.setTimeout(() => {
    pagination.value.page = 1
    loadUsers()
  }, 500)
}

const goToPage = (page: number) => {
  pagination.value.page = page
  loadUsers()
}

const editUser = (user: User) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    email: user.email,
    password: '',
    status: user.status
  }
  showEditModal.value = true
}

const confirmDelete = (user: User) => {
  userToDelete.value = user
  showDeleteModal.value = true
}

const deleteUser = async () => {
  if (!userToDelete.value) return

  try {
    await usersApi.deleteUser(userToDelete.value.uid)
    showDeleteModal.value = false
    userToDelete.value = null
    loadUsers()
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to delete user'
  }
}

const submitUser = async () => {
  try {
    if (showEditModal.value && editingUser.value) {
      await usersApi.updateUser(editingUser.value.uid, {
        username: userForm.value.username,
        email: userForm.value.email,
        status: userForm.value.status
      })
    } else {
      await usersApi.createUser({
        username: userForm.value.username,
        email: userForm.value.email,
        password: userForm.value.password
      })
    }

    closeModals()
    loadUsers()
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to save user'
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingUser.value = null
  userForm.value = {
    username: '',
    email: '',
    password: '',
    status: 1
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.users-view {
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header h1 {
  margin: 0;
  color: #2c3e50;
}

.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.search-input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.page-size-select {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover {
  background-color: #2980b9;
}

.btn-secondary {
  background-color: #95a5a6;
  color: white;
}

.btn-danger {
  background-color: #e74c3c;
  color: white;
}

.loading, .error, .empty {
  text-align: center;
  padding: 2rem;
  font-size: 1.1rem;
}

.error {
  color: #e74c3c;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.users-table th,
.users-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #ecf0f1;
}

.users-table th {
  background-color: #34495e;
  color: white;
  font-weight: 600;
}

.users-table tr:hover {
  background-color: #f8f9fa;
}

.status {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
}

.status.active {
  background-color: #27ae60;
  color: white;
}

.status.inactive {
  background-color: #95a5a6;
  color: white;
}

.action-btn {
  background: none;
  border: none;
  color: #3498db;
  cursor: pointer;
  margin-right: 0.5rem;
  font-size: 0.9rem;
}

.action-btn:hover {
  text-decoration: underline;
}

.action-btn.danger {
  color: #e74c3c;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.pagination-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  min-width: 400px;
  max-width: 500px;
}

.modal h2 {
  margin-top: 0;
  color: #2c3e50;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2c3e50;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style>

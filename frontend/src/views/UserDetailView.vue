<template>
  <div class="user-detail-view">
    <div v-if="loading" class="loading">Loading user details...</div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <router-link to="/users" class="btn btn-primary">Back to Users</router-link>
    </div>

    <div v-else-if="user" class="user-container">
      <div class="back-link">
        <router-link to="/users">&larr; Back to Users</router-link>
      </div>

      <div class="user-card">
        <div class="user-header">
          <div class="user-avatar">
            {{ getInitials(user.username) }}
          </div>
          <div class="user-info">
            <h1>{{ user.username }}</h1>
            <p class="email">{{ user.email }}</p>
            <span :class="['status-badge', user.status === 1 ? 'active' : 'inactive']">
              {{ user.status === 1 ? 'Active' : 'Inactive' }}
            </span>
          </div>
        </div>

        <div class="user-sections">
          <div class="section">
            <h2>Account Information</h2>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">User ID</span>
                <span class="value">{{ user.uid }}</span>
              </div>
              <div class="info-item">
                <span class="label">Username</span>
                <span class="value">{{ user.username }}</span>
              </div>
              <div class="info-item">
                <span class="label">Email</span>
                <span class="value">{{ user.email }}</span>
              </div>
              <div class="info-item">
                <span class="label">Status</span>
                <span class="value">{{ user.status === 1 ? 'Active' : 'Inactive' }}</span>
              </div>
              <div class="info-item">
                <span class="label">Created</span>
                <span class="value">{{ formatDate(user.created_at) }}</span>
              </div>
              <div class="info-item">
                <span class="label">Updated</span>
                <span class="value">{{ formatDate(user.updated_at) }}</span>
              </div>
            </div>
          </div>

          <div class="section">
            <h2>Profile</h2>
            <div v-if="profileLoading" class="loading-small">Loading profile...</div>
            <div v-else-if="profile" class="profile-card">
              <div class="profile-info">
                <p><strong>Full Name:</strong> {{ fullName }}</p>
                <p><strong>Bio:</strong> {{ profile.bio || 'No bio provided' }}</p>
              </div>
            </div>
            <div v-else class="empty-profile">
              <p>No profile information available</p>
            </div>
          </div>

          <div class="section">
            <h2>Devices</h2>
            <div v-if="devicesLoading" class="loading-small">Loading devices...</div>
            <div v-else-if="devices.length === 0" class="empty-devices">
              <p>No devices registered</p>
            </div>
            <div v-else class="devices-list">
              <div v-for="device in devices" :key="device.device_uid" class="device-card">
                <div class="device-name">{{ device.device_name }}</div>
                <div class="device-details">
                  <p><strong>IP Address:</strong> {{ device.ip_address }}</p>
                  <p><strong>Created:</strong> {{ formatDate(device.created_at) }}</p>
                  <p><strong>Last Active:</strong> {{ formatDate(device.last_active_at) }}</p>
                  <p v-if="device.revoked_at" class="revoked">
                    <strong>Revoked:</strong> {{ formatDate(device.revoked_at) }}
                  </p>
                </div>
                <span v-if="device.revoked_at" class="revoked-badge">Revoked</span>
                <span v-else class="active-badge">Active</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { usersApi, type User, type Profile, type Device } from '@/services/api'

const route = useRoute()

const loading = ref(true)
const error = ref<string | null>(null)
const profileLoading = ref(false)
const devicesLoading = ref(false)

const user = ref<User | null>(null)
const profile = ref<Profile | null>(null)
const devices = ref<Device[]>([])

const fullName = computed(() => {
  if (profile.value) {
    return `${profile.value.first_name} ${profile.value.last_name}`.trim()
  }
  return 'Not set'
})

const loadUser = async () => {
  const uid = route.params.uid as string
  loading.value = true

  try {
    user.value = await usersApi.getUser(uid)
    loadProfile(uid)
    loadDevices(uid)
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to load user'
  } finally {
    loading.value = false
  }
}

const loadProfile = async (uid: string) => {
  profileLoading.value = true
  try {
    profile.value = await usersApi.getProfile(uid)
  } catch (err) {
    console.error('Failed to load profile:', err)
  } finally {
    profileLoading.value = false
  }
}

const loadDevices = async (uid: string) => {
  devicesLoading.value = true
  try {
    const response = await usersApi.listDevices(uid)
    devices.value = response.devices
  } catch (err) {
    console.error('Failed to load devices:', err)
  } finally {
    devicesLoading.value = false
  }
}

const getInitials = (username: string) => {
  return username.substring(0, 2).toUpperCase()
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  loadUser()
})
</script>

<style scoped>
.user-detail-view {
  max-width: 900px;
  margin: 0 auto;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #7f8c8d;
}

.error {
  text-align: center;
  padding: 3rem;
  color: #e74c3c;
}

.back-link {
  margin-bottom: 1.5rem;
}

.back-link a {
  color: #3498db;
  text-decoration: none;
  font-weight: 500;
}

.back-link a:hover {
  text-decoration: underline;
}

.user-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.user-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2.5rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.user-avatar {
  width: 100px;
  height: 100px;
  background: rgba(255,255,255,0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  font-weight: bold;
  border: 4px solid rgba(255,255,255,0.3);
}

.user-info h1 {
  margin: 0 0 0.5rem 0;
  font-size: 2rem;
}

.email {
  margin: 0 0 1rem 0;
  opacity: 0.9;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
}

.status-badge.active {
  background-color: #27ae60;
}

.status-badge.inactive {
  background-color: #95a5a6;
}

.user-sections {
  padding: 2rem;
}

.section {
  margin-bottom: 2rem;
}

.section h2 {
  margin: 0 0 1.5rem 0;
  color: #2c3e50;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #ecf0f1;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.info-item {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
}

.info-item .label {
  display: block;
  font-weight: 500;
  color: #7f8c8d;
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
}

.info-item .value {
  color: #2c3e50;
  font-size: 1.1rem;
}

.loading-small {
  text-align: center;
  padding: 1rem;
  color: #7f8c8d;
}

.profile-card,
.empty-profile,
.empty-devices {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  text-align: center;
}

.profile-info p {
  margin: 0.5rem 0;
  color: #2c3e50;
}

.devices-list {
  display: grid;
  gap: 1rem;
}

.device-card {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.device-name {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 0.25rem;
}

.device-details p {
  margin: 0.25rem 0;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.device-details .revoked {
  color: #e74c3c;
}

.active-badge {
  background-color: #27ae60;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
}

.revoked-badge {
  background-color: #95a5a6;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
}

.btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 8px;
  text-decoration: none;
  cursor: pointer;
  font-size: 1rem;
}

.btn:hover {
  background-color: #2980b9;
}

@media (max-width: 768px) {
  .user-header {
    flex-direction: column;
    text-align: center;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>

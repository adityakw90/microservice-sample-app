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

        <!-- Tabs -->
        <div class="tabs-container">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            :class="['tab-button', { active: activeTab === tab.key }]"
            @click="activeTab = tab.key"
          >
            <component :is="tab.icon" class="tab-icon" />
            <span>{{ tab.label }}</span>
          </button>
        </div>

        <!-- Tab Content -->
        <div class="tab-content">
          <!-- Profile Tab -->
          <div v-if="activeTab === 'profile'" class="tab-panel">
            <div v-if="profileLoading" class="loading-small">Loading profile...</div>
            <div v-else-if="profile" class="profile-section">
              <div class="profile-info">
                <p><strong>Full Name:</strong> {{ fullName }}</p>
                <p><strong>Bio:</strong> {{ profile.bio || 'No bio provided' }}</p>
              </div>
            </div>
            <div v-else class="empty-profile">
              <p>No profile information available</p>
            </div>
          </div>

          <!-- Files Tab -->
          <div v-if="activeTab === 'files'" class="tab-panel">
            <UserFilesTab :user-uid="user.uid" />
          </div>

          <!-- Devices Tab -->
          <div v-if="activeTab === 'devices'" class="tab-panel">
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
import UserFilesTab from '@/components/User/UserFilesTab.vue'

const route = useRoute()

const loading = ref(true)
const error = ref<string | null>(null)
const profileLoading = ref(false)
const devicesLoading = ref(false)

const user = ref<User | null>(null)
const profile = ref<Profile | null>(null)
const devices = ref<Device[]>([])

const activeTab = ref<'profile' | 'files' | 'devices'>('profile')

// Tab definitions with inline icon components
const ProfileIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4-4v6c0 6 8-10l4 4 4 4 4 2H5a4 4 0 00-4-4V5a4 4 0 012-4h11"/>
      <circle cx="12" cy="7" r="4"/>
    </svg>
  `
}

const FilesIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
    </svg>
  `
}

const DevicesIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
      <line x1="12" y1="6" x2="12" y2="18" stroke-width="2"/>
    </svg>
  `
}

const tabs = [
  { key: 'profile', label: 'Profile', icon: ProfileIcon },
  { key: 'files', label: 'Files', icon: FilesIcon },
  { key: 'devices', label: 'Devices', icon: DevicesIcon }
]

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
  padding: 2rem;
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
  padding-bottom: 1rem;
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
  color: #667eea;
  border: 4px solid rgba(255,255,255,0.3);
}

.user-info h1 {
  margin: 0;
  font-size: 2rem;
  color: white;
}

.email {
  margin: 0.5rem 0 0;
  opacity: 0.9;
  color: white;
}

.status-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  background: rgba(255,255,255,0.2);
}

.status-badge.active {
  background: #27ae60;
}

.status-badge.inactive {
  background: #95a5a6;
}

.tabs-container {
  display: flex;
  gap: 0.5rem;
  padding: 0 2rem;
  border-bottom: 1px solid #e5e7eb;
}

.tab-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  color: #64748b;
  font-weight: 500;
  transition: all 0.2s;
}

.tab-button:hover {
  color: #667eea;
}

.tab-button.active {
  border-bottom-color: #667eea;
  color: #667eea;
}

.tab-icon {
  width: 20px;
  height: 20px;
}

.tab-content {
  padding: 2rem 0 0;
  min-height: 300px;
}

.tab-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.loading-small {
  text-align: center;
  padding: 2rem;
  color: #7f8c8d;
}

.profile-section {
  padding: 1rem;
}

.profile-info p {
  margin-bottom: 1rem;
  line-height: 1.6;
  color: #2c3e50;
}

.empty-profile,
.empty-devices {
  text-align: center;
  padding: 3rem;
  color: #7f8c8d;
}

.devices-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.device-card {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.device-name {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 1rem;
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
  font-weight: 600;
}

.revoked-badge {
  background-color: #95a5a6;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
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

  .user-avatar {
    width: 60px;
    height: 60px;
  }

  .devices-list {
    grid-template-columns: 1fr;
  }
}
</style>

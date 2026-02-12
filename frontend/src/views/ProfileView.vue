<template>
  <div class="profile-view">
    <div v-if="loading" class="loading">Loading profile...</div>

    <div v-else-if="profile" class="profile-container">
      <div class="profile-header">
        <div class="avatar-section">
          <div class="avatar">
            {{ initials }}
          </div>
        </div>
        <div class="user-info">
          <h1>{{ fullName || 'Anonymous User' }}</h1>
          <p class="username">@{{ authStore.user?.identifier }}</p>
        </div>
      </div>

      <div class="profile-sections">
        <div class="section">
          <h2>Profile Information</h2>
          <form @submit.prevent="saveProfile" class="profile-form">
            <div class="form-row">
              <div class="form-group">
                <label for="firstName">First Name</label>
                <input
                  id="firstName"
                  v-model="profileForm.first_name"
                  type="text"
                  class="form-input"
                />
              </div>
              <div class="form-group">
                <label for="lastName">Last Name</label>
                <input
                  id="lastName"
                  v-model="profileForm.last_name"
                  type="text"
                  class="form-input"
                />
              </div>
            </div>

            <div class="form-group">
              <label for="bio">Bio</label>
              <textarea
                id="bio"
                v-model="profileForm.bio"
                rows="4"
                class="form-textarea"
                placeholder="Tell us about yourself..."
              ></textarea>
            </div>

            <button type="submit" :disabled="saving" class="btn btn-primary">
              {{ saving ? 'Saving...' : 'Save Profile' }}
            </button>
          </form>
        </div>

        <div class="section">
          <h2>Account Details</h2>
          <div class="account-details">
            <div class="detail-item">
              <span class="label">User ID:</span>
              <span class="value">{{ authStore.user?.uid }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Identifier:</span>
              <span class="value">{{ authStore.user?.identifier }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Type:</span>
              <span class="value">{{ authStore.user?.identifier_type }}</span>
            </div>
          </div>
        </div>

        <div class="section">
          <h2>Active Devices</h2>
          <div v-if="devicesLoading" class="loading">Loading devices...</div>
          <div v-else-if="devices.length === 0" class="empty">No active devices</div>
          <div v-else class="devices-list">
            <div v-for="device in devices" :key="device.device_uid" class="device-item">
              <div class="device-info">
                <h4>{{ device.device_name }}</h4>
                <p>IP: {{ device.ip_address }}</p>
                <p>Last active: {{ formatDate(device.last_active_at) }}</p>
              </div>
              <button
                v-if="!device.revoked_at"
                @click="revokeDevice(device.device_uid)"
                class="btn btn-danger btn-sm"
              >
                Revoke
              </button>
              <span v-else class="revoked-badge">Revoked</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="error">
      <p>Failed to load profile. Please try logging in again.</p>
      <router-link to="/login" class="btn btn-primary">Go to Login</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usersApi, type Profile, type Device } from '@/services/api'

const authStore = useAuthStore()

const loading = ref(true)
const saving = ref(false)
const devicesLoading = ref(false)

const profile = ref<Profile | null>(null)
const devices = ref<Device[]>([])

const profileForm = ref({
  first_name: '',
  last_name: '',
  bio: ''
})

const fullName = computed(() => {
  if (profile.value) {
    return `${profile.value.first_name} ${profile.value.last_name}`.trim()
  }
  return ''
})

const initials = computed(() => {
  if (profile.value) {
    return `${profile.value.first_name[0] || ''}${profile.value.last_name[0] || ''}`.toUpperCase()
  }
  return 'U'
})

const loadProfile = async () => {
  if (!authStore.user?.uid) return

  loading.value = true
  try {
    const data = await usersApi.getProfile(authStore.user.uid)
    profile.value = data
    profileForm.value = {
      first_name: data.first_name,
      last_name: data.last_name,
      bio: data.bio
    }
  } catch (err: any) {
    console.error('Failed to load profile:', err)
  } finally {
    loading.value = false
  }
}

const loadDevices = async () => {
  if (!authStore.user?.uid) return

  devicesLoading.value = true
  try {
    const response = await usersApi.listDevices(authStore.user.uid)
    devices.value = response.devices
  } catch (err: any) {
    console.error('Failed to load devices:', err)
  } finally {
    devicesLoading.value = false
  }
}

const saveProfile = async () => {
  if (!authStore.user?.uid) return

  saving.value = true
  try {
    await usersApi.updateProfile(authStore.user.uid, profileForm.value)
    await loadProfile()
  } catch (err: any) {
    console.error('Failed to save profile:', err)
  } finally {
    saving.value = false
  }
}

const revokeDevice = async (deviceUid: string) => {
  if (!authStore.user?.uid) return

  try {
    await usersApi.revokeDevice(authStore.user.uid, deviceUid)
    await loadDevices()
  } catch (err: any) {
    console.error('Failed to revoke device:', err)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  loadProfile()
  loadDevices()
})
</script>

<style scoped>
.profile-view {
  max-width: 900px;
  margin: 0 auto;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #7f8c8d;
}

.profile-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.profile-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2.5rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.avatar-section {
  flex-shrink: 0;
}

.avatar {
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

.username {
  margin: 0;
  opacity: 0.9;
  font-size: 1.1rem;
}

.profile-sections {
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

.profile-form {
  max-width: 600px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2c3e50;
  font-weight: 500;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  font-family: inherit;
  box-sizing: border-box;
}

.form-textarea {
  resize: vertical;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2980b9;
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-danger {
  background-color: #e74c3c;
  color: white;
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
}

.account-details {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  border-bottom: 1px solid #e9ecef;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-item .label {
  font-weight: 500;
  color: #7f8c8d;
}

.detail-item .value {
  color: #2c3e50;
}

.devices-list {
  display: grid;
  gap: 1rem;
}

.device-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.device-info h4 {
  margin: 0 0 0.25rem 0;
  color: #2c3e50;
}

.device-info p {
  margin: 0.25rem 0;
  color: #7f8c8d;
  font-size: 0.9rem;
}

.revoked-badge {
  background-color: #95a5a6;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
}

.empty {
  text-align: center;
  padding: 2rem;
  color: #7f8c8d;
}

.error {
  text-align: center;
  padding: 3rem;
  color: #e74c3c;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>

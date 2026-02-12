<template>
  <div class="proto-test-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Proto Service Tester</h1>
        <p class="page-subtitle">Interactive testing interface for all gRPC services</p>
      </div>
      <div class="header-actions">
        <button @click="expandAll" class="action-btn">Expand All</button>
        <button @click="collapseAll" class="action-btn">Collapse All</button>
        <button @click="clearAllResponses" class="action-btn secondary">Clear Responses</button>
      </div>
    </div>

    <div class="services-grid">
      <!-- AuthService Section -->
      <ServiceSection
        title="AuthService"
        :icon="AuthIcon"
        :endpoints="authEndpoints"
        :loading="loadingStates.auth"
        @execute="executeAuth"
      />

      <!-- UserService Section -->
      <ServiceSection
        title="UserService"
        :icon="UserIcon"
        :endpoints="userEndpoints"
        :loading="loadingStates.user"
        @execute="executeUser"
      />

      <!-- DeviceService Section -->
      <ServiceSection
        title="DeviceService"
        :icon="DeviceIcon"
        :endpoints="deviceEndpoints"
        :loading="loadingStates.device"
        @execute="executeDevice"
        :not-implemented="true"
      />

      <!-- UserFileService Section -->
      <ServiceSection
        title="UserFileService"
        :icon="FileIcon"
        :endpoints="userFileEndpoints"
        :loading="loadingStates.userFile"
        @execute="executeUserFile"
        :not-implemented="true"
      />
    </div>

    <!-- Global Response Panel -->
    <div v-if="globalResponse" class="global-response-panel">
      <div class="response-header">
        <h3>Latest Response</h3>
        <button @click="globalResponse = null" class="close-btn">Close</button>
      </div>
      <pre class="response-body">{{ JSON.stringify(globalResponse, null, 2) }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { authApi, usersApi, deviceApi, userFileApi } from '@/services/api'
import ServiceSection from '@/components/ProtoTest/ServiceSection.vue'
import { useNotificationStore } from '@/stores/notification'

const notificationStore = useNotificationStore()

interface Endpoint {
  name: string
  method: string
  description: string
  params: ParamDefinition[]
  response?: any
  error?: string
  expanded?: boolean
}

interface ParamDefinition {
  name: string
  type: string
  required: boolean
  default?: any
  description?: string
}

interface ServiceEndpoints {
  [key: string]: Endpoint
}

// Icon components
const AuthIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
      <path d="M9 12l2 2 4-4"/>
    </svg>
  `
}

const UserIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
      <circle cx="9" cy="7" r="4"/>
      <path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/>
    </svg>
  `
}

const DeviceIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
      <line x1="12" y1="18" x2="12" y2="18"/>
    </svg>
  `
}

const FileIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M13 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V9z"/>
      <polyline points="13,2 13,9 20,9"/>
    </svg>
  `
}

// Loading states
const loadingStates = reactive({
  auth: {} as Record<string, boolean>,
  user: {} as Record<string, boolean>,
  device: {} as Record<string, boolean>,
  userFile: {} as Record<string, boolean>
})

// Global response for display
const globalResponse = ref<any>(null)

// AuthService endpoints
const authEndpoints = reactive<ServiceEndpoints>({
  login: {
    name: 'login',
    method: 'POST',
    description: 'Authenticate user with credentials',
    params: [
      { name: 'identifier', type: 'string', required: true, description: 'Username or email' },
      { name: 'identifier_type', type: 'string', required: true, default: 'username', description: 'Type of identifier (username/email)' },
      { name: 'password', type: 'string', required: true, description: 'User password' }
    ],
    expanded: false
  },
  refreshToken: {
    name: 'refreshToken',
    method: 'POST',
    description: 'Refresh access token',
    params: [
      { name: 'refresh_token', type: 'string', required: true, description: 'Refresh token' }
    ],
    expanded: false
  },
  validateToken: {
    name: 'validateToken',
    method: 'POST',
    description: 'Validate JWT token',
    params: [
      { name: 'access_token', type: 'string', required: true, description: 'Access token to validate' }
    ],
    expanded: false
  },
  verifyPin: {
    name: 'verifyPin',
    method: 'POST',
    description: 'Verify user PIN',
    params: [
      { name: 'uid', type: 'string', required: true, description: 'User UID' },
      { name: 'code', type: 'string', required: true, description: 'PIN code' }
    ],
    expanded: false
  }
})

// UserService endpoints
const userEndpoints = reactive<ServiceEndpoints>({
  listUsers: {
    name: 'listUsers',
    method: 'GET',
    description: 'List users with pagination and filtering',
    params: [
      { name: 'page', type: 'number', required: false, default: 1 },
      { name: 'limit', type: 'number', required: false, default: 10 },
      { name: 'username', type: 'string', required: false },
      { name: 'email', type: 'string', required: false },
      { name: 'query', type: 'string', required: false },
      { name: 'active', type: 'boolean', required: false }
    ],
    expanded: false
  },
  getUser: {
    name: 'getUser',
    method: 'GET',
    description: 'Get user by UID',
    params: [
      { name: 'uid', type: 'string', required: true, description: 'User UID' }
    ],
    expanded: false
  },
  createUser: {
    name: 'createUser',
    method: 'POST',
    description: 'Create new user',
    params: [
      { name: 'username', type: 'string', required: true },
      { name: 'email', type: 'string', required: true },
      { name: 'password', type: 'string', required: true }
    ],
    expanded: false
  },
  updateUser: {
    name: 'updateUser',
    method: 'PUT',
    description: 'Update user details',
    params: [
      { name: 'uid', type: 'string', required: true },
      { name: 'username', type: 'string', required: false },
      { name: 'email', type: 'string', required: false },
      { name: 'password', type: 'string', required: false },
      { name: 'status', type: 'number', required: false }
    ],
    expanded: false
  },
  deleteUser: {
    name: 'deleteUser',
    method: 'DELETE',
    description: 'Delete user',
    params: [
      { name: 'uid', type: 'string', required: true }
    ],
    expanded: false
  },
  getProfile: {
    name: 'getProfile',
    method: 'GET',
    description: 'Get user profile',
    params: [
      { name: 'userUid', type: 'string', required: true, description: 'User UID' }
    ],
    expanded: false
  },
  updateProfile: {
    name: 'updateProfile',
    method: 'PUT',
    description: 'Update user profile',
    params: [
      { name: 'userUid', type: 'string', required: true },
      { name: 'first_name', type: 'string', required: false },
      { name: 'last_name', type: 'string', required: false },
      { name: 'bio', type: 'string', required: false }
    ],
    expanded: false
  },
  listDevices: {
    name: 'listDevices',
    method: 'GET',
    description: 'List user devices',
    params: [
      { name: 'userUid', type: 'string', required: true },
      { name: 'page', type: 'number', required: false, default: 1 },
      { name: 'limit', type: 'number', required: false, default: 10 }
    ],
    expanded: false
  },
  revokeDevice: {
    name: 'revokeDevice',
    method: 'DELETE',
    description: 'Revoke user device',
    params: [
      { name: 'userUid', type: 'string', required: true },
      { name: 'deviceUid', type: 'string', required: true }
    ],
    expanded: false
  }
})

// DeviceService endpoints (not implemented in backend)
const deviceEndpoints = reactive<ServiceEndpoints>({
  listDevices: {
    name: 'listDevices',
    method: 'GET',
    description: 'List all devices',
    params: [
      { name: 'page', type: 'number', required: false, default: 1 },
      { name: 'limit', type: 'number', required: false, default: 10 }
    ],
    expanded: false
  },
  getDevice: {
    name: 'getDevice',
    method: 'GET',
    description: 'Get device by UID',
    params: [
      { name: 'deviceUid', type: 'string', required: true }
    ],
    expanded: false
  },
  deleteDevice: {
    name: 'deleteDevice',
    method: 'DELETE',
    description: 'Delete device',
    params: [
      { name: 'deviceUid', type: 'string', required: true }
    ],
    expanded: false
  }
})

// UserFileService endpoints (not implemented in backend)
const userFileEndpoints = reactive<ServiceEndpoints>({
  listFiles: {
    name: 'listFiles',
    method: 'GET',
    description: 'List files with filtering',
    params: [
      { name: 'page', type: 'number', required: false, default: 1 },
      { name: 'limit', type: 'number', required: false, default: 10 },
      { name: 'user_uid', type: 'string', required: false },
      { name: 'file_type', type: 'string', required: false },
      { name: 'query', type: 'string', required: false }
    ],
    expanded: false
  },
  getFile: {
    name: 'getFile',
    method: 'GET',
    description: 'Get file by UID',
    params: [
      { name: 'fileUid', type: 'string', required: true }
    ],
    expanded: false
  },
  uploadFile: {
    name: 'uploadFile',
    method: 'POST',
    description: 'Upload new file',
    params: [
      { name: 'user_uid', type: 'string', required: true },
      { name: 'file_type', type: 'string', required: true },
      { name: 'file_name', type: 'string', required: true },
      { name: 'file_path', type: 'string', required: true },
      { name: 'mime_type', type: 'string', required: true },
      { name: 'size', type: 'number', required: true },
      { name: 'visibility', type: 'number', required: false, default: 0 }
    ],
    expanded: false
  },
  updateFile: {
    name: 'updateFile',
    method: 'PUT',
    description: 'Update file',
    params: [
      { name: 'fileUid', type: 'string', required: true },
      { name: 'file_name', type: 'string', required: false },
      { name: 'visibility', type: 'number', required: false }
    ],
    expanded: false
  },
  deleteFile: {
    name: 'deleteFile',
    method: 'DELETE',
    description: 'Delete file',
    params: [
      { name: 'fileUid', type: 'string', required: true }
    ],
    expanded: false
  }
})

// Helper function to collect form values
const collectParams = (endpoint: Endpoint, formData: FormData) => {
  const params: Record<string, any> = {}

  for (const param of endpoint.params) {
    const value = formData.get(param.name)

    if (value && value.toString().trim() !== '') {
      // Type conversion based on param type
      if (param.type === 'number') {
        params[param.name] = Number(value)
      } else if (param.type === 'boolean') {
        params[param.name] = value === 'true'
      } else {
        params[param.name] = value
      }
    } else if (param.default !== undefined) {
      params[param.name] = param.default
    }
  }

  return params
}

// AuthService execution
const executeAuth = async (endpointName: string, formData: FormData) => {
  const endpoint = authEndpoints[endpointName]
  loadingStates.auth[endpointName] = true
  endpoint.response = null
  endpoint.error = null

  try {
    const params = collectParams(endpoint, formData)

    let response: any
    switch (endpointName) {
      case 'login':
        response = await authApi.login(params.identifier, params.identifier_type, params.password)
        break
      case 'refreshToken':
        response = await authApi.refreshToken(params.refresh_token)
        break
      case 'validateToken':
        response = await authApi.validateToken(params.access_token)
        break
      case 'verifyPin':
        response = await authApi.verifyPin(params.uid, params.code)
        break
    }

    endpoint.response = response
    globalResponse.value = response
    notificationStore.showSuccess(`AuthService.${endpointName} executed successfully`)
  } catch (error: any) {
    endpoint.error = error.response?.data || error.message
    globalResponse.value = { error: endpoint.error }
    notificationStore.showError(`AuthService.${endpointName} failed: ${error.message}`)
  } finally {
    loadingStates.auth[endpointName] = false
  }
}

// UserService execution
const executeUser = async (endpointName: string, formData: FormData) => {
  const endpoint = userEndpoints[endpointName]
  loadingStates.user[endpointName] = true
  endpoint.response = null
  endpoint.error = null

  try {
    const params = collectParams(endpoint, formData)

    let response: any
    switch (endpointName) {
      case 'listUsers':
        response = await usersApi.listUsers(params)
        break
      case 'getUser':
        response = await usersApi.getUser(params.uid)
        break
      case 'createUser':
        response = await usersApi.createUser(params)
        break
      case 'updateUser':
        response = await usersApi.updateUser(params.uid, params)
        break
      case 'deleteUser':
        response = await usersApi.deleteUser(params.uid)
        break
      case 'getProfile':
        response = await usersApi.getProfile(params.userUid)
        break
      case 'updateProfile':
        response = await usersApi.updateProfile(params.userUid, params)
        break
      case 'listDevices':
        response = await usersApi.listDevices(params.userUid, params)
        break
      case 'revokeDevice':
        response = await usersApi.revokeDevice(params.userUid, params.deviceUid)
        break
    }

    endpoint.response = response
    globalResponse.value = response
    notificationStore.showSuccess(`UserService.${endpointName} executed successfully`)
  } catch (error: any) {
    endpoint.error = error.response?.data || error.message
    globalResponse.value = { error: endpoint.error }
    notificationStore.showError(`UserService.${endpointName} failed: ${error.message}`)
  } finally {
    loadingStates.user[endpointName] = false
  }
}

// DeviceService execution (not implemented)
const executeDevice = async (endpointName: string, formData: FormData) => {
  notificationStore.showWarning('DeviceService is not implemented in the backend yet')
}

// UserFileService execution (not implemented)
const executeUserFile = async (endpointName: string, formData: FormData) => {
  notificationStore.showWarning('UserFileService is not implemented in the backend yet')
}

// Utility functions
const expandAll = () => {
  Object.values(authEndpoints).forEach(e => e.expanded = true)
  Object.values(userEndpoints).forEach(e => e.expanded = true)
  Object.values(deviceEndpoints).forEach(e => e.expanded = true)
  Object.values(userFileEndpoints).forEach(e => e.expanded = true)
}

const collapseAll = () => {
  Object.values(authEndpoints).forEach(e => e.expanded = false)
  Object.values(userEndpoints).forEach(e => e.expanded = false)
  Object.values(deviceEndpoints).forEach(e => e.expanded = false)
  Object.values(userFileEndpoints).forEach(e => e.expanded = false)
}

const clearAllResponses = () => {
  Object.values(authEndpoints).forEach(e => { e.response = null; e.error = null })
  Object.values(userEndpoints).forEach(e => { e.response = null; e.error = null })
  Object.values(deviceEndpoints).forEach(e => { e.response = null; e.error = null })
  Object.values(userFileEndpoints).forEach(e => { e.response = null; e.error = null })
  globalResponse.value = null
  notificationStore.showInfo('All responses cleared')
}
</script>

<style scoped>
.proto-test-page {
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.5rem;
}

.page-subtitle {
  color: var(--color-text-secondary);
  font-size: 0.95rem;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.action-btn {
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  transition: all 0.2s ease;
  cursor: pointer;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.action-btn.secondary {
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
}

.action-btn.secondary:hover {
  background: var(--color-surface);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.services-grid {
  display: grid;
  gap: 1.5rem;
}

.global-response-panel {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  width: 500px;
  max-width: calc(100vw - 4rem);
  max-height: 400px;
  background: var(--color-surface);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 1000;
  animation: slideInUp 0.3s ease;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.response-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.response-header h3 {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  padding: 0.4rem 0.8rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
  background: var(--color-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: var(--color-error);
  color: white;
  border-color: var(--color-error);
}

.response-body {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  margin: 0;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.8rem;
  line-height: 1.5;
  background: var(--color-bg);
  border-radius: 0 0 12px 12px;
  color: var(--color-text-primary);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
  }

  .global-response-panel {
    width: calc(100vw - 2rem);
    right: 1rem;
    left: 1rem;
    bottom: 1rem;
  }
}
</style>

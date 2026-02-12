import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/services/api'

export interface User {
  uid: string
  identifier: string
  identifier_type: string
  claims?: Record<string, any>
}

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(identifier: string, identifier_type: string, password: string) {
    try {
      const response = await authApi.login(identifier, identifier_type, password)
      accessToken.value = response.access_token
      refreshToken.value = response.refresh_token

      localStorage.setItem('access_token', response.access_token)
      localStorage.setItem('refresh_token', response.refresh_token)

      // Validate token to get user info
      await validateToken()
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  async function validateToken() {
    if (!accessToken.value) return false

    try {
      const response = await authApi.validateToken(accessToken.value)
      user.value = {
        uid: response.uid,
        identifier: response.identifier,
        identifier_type: response.identifier_type,
        claims: response.claims
      }
      return true
    } catch (error) {
      console.error('Token validation failed:', error)
      logout()
      return false
    }
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) return false

    try {
      const response = await authApi.refreshToken(refreshToken.value)
      accessToken.value = response.access_token
      refreshToken.value = response.refresh_token

      localStorage.setItem('access_token', response.access_token)
      localStorage.setItem('refresh_token', response.refresh_token)
      return true
    } catch (error) {
      console.error('Token refresh failed:', error)
      logout()
      return false
    }
  }

  function logout() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  // Initialize store
  if (accessToken.value) {
    validateToken()
  }

  return {
    accessToken,
    refreshToken,
    user,
    isAuthenticated,
    login,
    validateToken,
    refreshAccessToken,
    logout
  }

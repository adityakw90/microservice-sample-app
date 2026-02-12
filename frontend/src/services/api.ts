import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken
          })

          const { access_token, refresh_token: newRefreshToken } = response.data
          localStorage.setItem('access_token', access_token)
          localStorage.setItem('refresh_token', newRefreshToken)

          originalRequest.headers.Authorization = `Bearer ${access_token}`
          return apiClient(originalRequest)
        }
      } catch (refreshError) {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

// API Types
export interface User {
  uid: string
  username: string
  email: string
  status: number
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface Profile {
  uid: string
  first_name: string
  last_name: string
  bio: string
  avatar?: {
    uid: string
    url: string
  }
}

export interface Device {
  device_uid: string
  device_name: string
  ip_address: string
  last_active_at: string
  created_at: string
  revoked_at?: string
}

export interface Meta {
  page: number
  limit: number
  total: number
  pages: number
}

export interface UsersResponse {
  users: User[]
  meta: Meta
}

export interface DevicesResponse {
  devices: Device[]
  meta: Meta
}

// Auth API
export const authApi = {
  async login(identifier: string, identifier_type: string, password: string) {
    const response = await apiClient.post('/auth/login', {
      identifier,
      identifier_type,
      password
    })
    return response.data
  },

  async refreshToken(refresh_token: string) {
    const response = await apiClient.post('/auth/refresh', { refresh_token })
    return response.data
  },

  async validateToken(access_token: string) {
    const response = await apiClient.post('/auth/validate', { access_token })
    return response.data
  },

  async verifyPin(uid: string, code: string) {
    const response = await apiClient.post('/auth/verify-pin', { uid, code })
    return response.data
  }
}

// Users API
export const usersApi = {
  async listUsers(params?: {
    page?: number
    limit?: number
    username?: string
    email?: string
    query?: string
    active?: boolean
  }): Promise<UsersResponse> {
    const response = await apiClient.get('/users', { params })
    return response.data
  },

  async getUser(uid: string): Promise<User> {
    const response = await apiClient.get(`/users/${uid}`)
    return response.data
  },

  async createUser(data: {
    username: string
    email: string
    password: string
  }): Promise<{ uid: string }> {
    const response = await apiClient.post('/users', data)
    return response.data
  },

  async updateUser(uid: string, data: {
    username?: string
    email?: string
    password?: string
    status?: number
  }): Promise<{ success: boolean }> {
    const response = await apiClient.put(`/users/${uid}`, data)
    return response.data
  },

  async deleteUser(uid: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete(`/users/${uid}`)
    return response.data
  },

  async getProfile(userUid: string): Promise<Profile> {
    const response = await apiClient.get(`/users/${userUid}/profile`)
    return response.data
  },

  async updateProfile(userUid: string, data: {
    first_name: string
    last_name: string
    bio: string
  }): Promise<{ success: boolean }> {
    const response = await apiClient.put(`/users/${userUid}/profile`, data)
    return response.data
  },

  async listDevices(userUid: string, params?: {
    page?: number
    limit?: number
  }): Promise<DevicesResponse> {
    const response = await apiClient.get(`/users/${userUid}/devices`, { params })
    return response.data
  },

  async revokeDevice(userUid: string, deviceUid: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete(`/users/${userUid}/devices/${deviceUid}`)
    return response.data
  }
}

// Device API (for standalone device management)
export const deviceApi = {
  async listDevices(params?: {
    page?: number
    limit?: number
  }): Promise<{ devices: Device[], meta: Meta }> {
    const response = await apiClient.get('/devices', { params })
    return response.data
  },

  async getDevice(deviceUid: string): Promise<Device> {
    const response = await apiClient.get(`/devices/${deviceUid}`)
    return response.data
  },

  async deleteDevice(deviceUid: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete(`/devices/${deviceUid}`)
    return response.data
  }
}

// UserFile API
export interface UserFile {
  uid: string
  user_uid: string
  file_type: string
  file_name: string
  file_path: string
  mime_type: string
  size: number
  visibility: number
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface UserFilesResponse {
  files: UserFile[]
  meta: Meta
}

export const userFileApi = {
  async listFiles(params?: {
    page?: number
    limit?: number
    user_uid?: string
    file_type?: string
    query?: string
  }): Promise<UserFilesResponse> {
    const response = await apiClient.get('/files', { params })
    return response.data
  },

  async getFile(fileUid: string): Promise<UserFile> {
    const response = await apiClient.get(`/files/${fileUid}`)
    return response.data
  },

  async uploadFile(data: {
    user_uid: string
    file_type: string
    file_name: string
    file_path: string
    mime_type: string
    size: number
    visibility?: number
  }): Promise<{ uid: string }> {
    const response = await apiClient.post('/files', data)
    return response.data
  },

  async updateFile(fileUid: string, data: {
    file_name?: string
    visibility?: number
  }): Promise<{ success: boolean }> {
    const response = await apiClient.put(`/files/${fileUid}`, data)
    return response.data
  },

  async deleteFile(fileUid: string): Promise<{ success: boolean }> {
    const response = await apiClient.delete(`/files/${fileUid}`)
    return response.data
  }
}

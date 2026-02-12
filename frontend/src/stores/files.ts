import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userFileApi, type UserFile, type UserFilesResponse, type Meta } from '@/services/api'

export const useFilesStore = defineStore('files', () => {
  const files = ref<UserFile[]>([])
  const pagination = ref<Meta>({
    page: 1,
    limit: 20,
    total: 0,
    pages: 0
  })
  const uploadProgress = ref<Map<string, number>>(new Map())
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const isEmpty = computed(() => files.value.length === 0)
  const hasMore = computed(() => pagination.value.page < pagination.value.pages)

  // Actions
  async function fetchFiles(params?: {
    page?: number
    limit?: number
    user_uid?: string
    file_type?: string
  }) {
    loading.value = true
    error.value = null

    try {
      const response = await userFileApi.listFiles(params)
      files.value = response.files
      pagination.value = response.meta
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch files'
      console.error('Failed to fetch files:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchFile(fileUid: string) {
    loading.value = true
    error.value = null

    try {
      const file = await userFileApi.getFile(fileUid)
      return file
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch file'
      console.error('Failed to fetch file:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  async function uploadFile(data: {
    user_uid: string
    file_type: string
    file_name: string
    file_path: string
    mime_type: string
    size: number
    visibility?: number
  }) {
    loading.value = true
    error.value = null

    try {
      const response = await userFileApi.uploadFile(data)
      // Add to files list
      files.value.unshift({
        uid: response.uid,
        ...data,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      } as UserFile)
      return response.uid
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to upload file'
      console.error('Failed to upload file:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateFile(fileUid: string, data: {
    file_name?: string
    visibility?: number
  }) {
    loading.value = true
    error.value = null

    try {
      await userFileApi.updateFile(fileUid, data)
      // Update in list
      const index = files.value.findIndex(f => f.uid === fileUid)
      if (index !== -1) {
        if (data.file_name) files.value[index].file_name = data.file_name
        if (data.visibility !== undefined) files.value[index].visibility = data.visibility
      }
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update file'
      console.error('Failed to update file:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  async function deleteFile(fileUid: string) {
    loading.value = true
    error.value = null

    try {
      await userFileApi.deleteFile(fileUid)
      // Remove from list
      files.value = files.value.filter(f => f.uid !== fileUid)
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete file'
      console.error('Failed to delete file:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    files,
    pagination,
    uploadProgress,
    loading,
    error,
    isEmpty,
    hasMore,
    fetchFiles,
    fetchFile,
    uploadFile,
    updateFile,
    deleteFile,
    clearError
  }
})

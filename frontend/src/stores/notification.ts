import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  const notification = ref({
    show: false,
    message: '',
    type: 'success' as 'success' | 'error' | 'warning' | 'info'
  })

  let timeoutId: number | null = null

  const show = (message: string, type: 'success' | 'error' | 'warning' | 'info') => {
    notification.value = { show: true, message, type }

    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    timeoutId = window.setTimeout(() => {
      notification.value.show = false
    }, 3000)
  }

  const hide = () => {
    notification.value.show = false
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  const showSuccess = (message: string) => show(message, 'success')
  const showError = (message: string) => show(message, 'error')
  const showWarning = (message: string) => show(message, 'warning')
  const showInfo = (message: string) => show(message, 'info')

  return {
    notification,
    show,
    hide,
    showSuccess,
    showError,
    showWarning,
    showInfo
  }
})
